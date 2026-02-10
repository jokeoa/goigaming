package ports

import (
	"context"

	"github.com/google/uuid"
	"github.com/jokeoa/igaming/internal/core/domain"
)

type AuthService interface {
	Register(ctx context.Context, username, email, password string) (domain.User, error)
	Login(ctx context.Context, email, password string) (domain.TokenPair, error)
	ValidateToken(token string) (domain.TokenClaims, error)
}

type UserService interface {
	GetProfile(ctx context.Context, userID uuid.UUID) (domain.UserProfile, error)
	GetByID(ctx context.Context, userID uuid.UUID) (domain.User, error)
}

type WalletService interface {
	CreateWallet(ctx context.Context, userID uuid.UUID) (domain.Wallet, error)
	GetBalance(ctx context.Context, userID uuid.UUID) (domain.Wallet, error)
	Deposit(ctx context.Context, userID uuid.UUID, amount string) (domain.Wallet, error)
	Withdraw(ctx context.Context, userID uuid.UUID, amount string) (domain.Wallet, error)
	GetTransactions(ctx context.Context, userID uuid.UUID, limit, offset int) ([]domain.Transaction, error)
}
