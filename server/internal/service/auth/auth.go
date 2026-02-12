package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jokeoa/goigaming/internal/core/domain"
	"github.com/jokeoa/goigaming/internal/core/ports"
	"github.com/jokeoa/goigaming/internal/repository/postgres"
	"github.com/shopspring/decimal"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	pool       *pgxpool.Pool
	userRepo   ports.UserRepository
	userFn     func(db postgres.DBTX) ports.UserRepository
	walletFn   func(db postgres.DBTX) ports.WalletRepository
	jwtSecret  []byte
	tokenTTL   time.Duration
}

func NewService(
	pool *pgxpool.Pool,
	userRepo ports.UserRepository,
	userFn func(db postgres.DBTX) ports.UserRepository,
	walletFn func(db postgres.DBTX) ports.WalletRepository,
	jwtSecret string,
	tokenTTL time.Duration,
) *Service {
	return &Service{
		pool:      pool,
		userRepo:  userRepo,
		userFn:    userFn,
		walletFn:  walletFn,
		jwtSecret: []byte(jwtSecret),
		tokenTTL:  tokenTTL,
	}
}

func (s *Service) Register(ctx context.Context, username, email, password string) (domain.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return domain.User{}, fmt.Errorf("AuthService.Register hash: %w", err)
	}

	var user domain.User

	err = postgres.RunInTx(ctx, s.pool, func(tx pgx.Tx) error {
		userRepo := s.userFn(tx)
		walletRepo := s.walletFn(tx)

		u, createErr := userRepo.Create(ctx, domain.User{
			Username:     username,
			Email:        email,
			PasswordHash: string(hash),
		})
		if createErr != nil {
			return createErr
		}

		_, createErr = walletRepo.Create(ctx, domain.Wallet{
			UserID:  u.ID,
			Balance: decimal.Zero,
		})
		if createErr != nil {
			return fmt.Errorf("create wallet: %w", createErr)
		}

		user = u
		return nil
	})
	if err != nil {
		return domain.User{}, fmt.Errorf("AuthService.Register: %w", err)
	}

	return user, nil
}

func (s *Service) Login(ctx context.Context, email, password string) (domain.TokenPair, error) {
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return domain.TokenPair{}, domain.ErrInvalidCredentials
		}
		return domain.TokenPair{}, fmt.Errorf("AuthService.Login: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return domain.TokenPair{}, domain.ErrInvalidCredentials
	}

	accessToken, err := s.generateAccessToken(user)
	if err != nil {
		return domain.TokenPair{}, err
	}

	refreshToken, err := s.generateRefreshToken(user)
	if err != nil {
		return domain.TokenPair{}, err
	}

	user.RefreshToken = &refreshToken
	_, err = s.userRepo.Update(ctx, user)
	if err != nil {
		return domain.TokenPair{}, fmt.Errorf("save refresh token: %w", err)
	}

	return domain.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(s.tokenTTL.Seconds()),
	}, nil
}

func (s *Service) ValidateToken(tokenString string) (domain.TokenClaims, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return s.jwtSecret, nil
	})
	if err != nil {
		return domain.TokenClaims{}, domain.ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return domain.TokenClaims{}, domain.ErrInvalidToken
	}

	sub, ok := claims["sub"].(string)
	if !ok {
		return domain.TokenClaims{}, domain.ErrInvalidToken
	}

	userID, err := uuid.Parse(sub)
	if err != nil {
		return domain.TokenClaims{}, domain.ErrInvalidToken
	}

	username, _ := claims["username"].(string)

	return domain.TokenClaims{
		UserID:   userID,
		Username: username,
	}, nil
}

func (s *Service) RefreshToken(ctx context.Context, refreshToken string) (domain.TokenPair, error) {
	claims, err := s.validateRefreshToken(refreshToken)
	if err != nil {
		return domain.TokenPair{}, domain.ErrInvalidToken
	}

	user, err := s.userRepo.FindByID(ctx, claims.UserID)
	if err != nil {
		return domain.TokenPair{}, domain.ErrUserNotFound
	}

	if user.RefreshToken == nil || *user.RefreshToken != refreshToken {
		return domain.TokenPair{}, domain.ErrInvalidToken
	}

	newAccessToken, err := s.generateAccessToken(user)
	if err != nil {
		return domain.TokenPair{}, err
	}

	newRefreshToken, err := s.generateRefreshToken(user)
	if err != nil {
		return domain.TokenPair{}, err
	}

	user.RefreshToken = &newRefreshToken
	_, err = s.userRepo.Update(ctx, user)
	if err != nil {
		return domain.TokenPair{}, err
	}

	return domain.TokenPair{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(s.tokenTTL.Seconds()),
	}, nil
}

func (s *Service) validateRefreshToken(tokenString string) (domain.TokenClaims, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return s.jwtSecret, nil
	})
	if err != nil {
		return domain.TokenClaims{}, domain.ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return domain.TokenClaims{}, domain.ErrInvalidToken
	}

	sub, ok := claims["sub"].(string)
	if !ok {
		return domain.TokenClaims{}, domain.ErrInvalidToken
	}

	userID, err := uuid.Parse(sub)
	if err != nil {
		return domain.TokenClaims{}, domain.ErrInvalidToken
	}

	username, _ := claims["username"].(string)

	return domain.TokenClaims{
		UserID:   userID,
		Username: username,
	}, nil
}

func (s *Service) generateAccessToken(user domain.User) (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"sub":      user.ID.String(),
		"username": user.Username,
		"iat":      now.Unix(),
		"exp":      now.Add(15 * time.Minute).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

func (s *Service) generateRefreshToken(user domain.User) (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"sub":      user.ID.String(),
		"username": user.Username,
		"iat":      now.Unix(),
		"exp":      now.Add(7 * 24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}
