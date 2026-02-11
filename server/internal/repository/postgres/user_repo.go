package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jokeoa/goigaming/internal/core/domain"
)

type UserRepository struct {
	db DBTX
}

func NewUserRepository(db DBTX) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user domain.User) (domain.User, error) {
	query := `
		INSERT INTO users (username, email, password_hash)
		VALUES ($1, $2, $3)
		RETURNING id, username, email, password_hash, created_at, updated_at
	`

	var u domain.User
	err := r.db.QueryRow(ctx, query, user.Username, user.Email, user.PasswordHash).Scan(
		&u.ID, &u.Username, &u.Email, &u.PasswordHash, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return u, domain.ErrUserAlreadyExists
		}
		return u, fmt.Errorf("UserRepository.Create: %w", err)
	}

	return u, nil
}

query := `
	SELECT id, username, email, password_hash, refresh_token, created_at, updated_at
	FROM users
	WHERE id = $1
`
var u domain.User
err := r.db.QueryRow(ctx, query, id).Scan(
	&u.ID, &u.Username, &u.Email, &u.PasswordHash, &u.RefreshToken, &u.CreatedAt, &u.UpdatedAt,
)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return u, domain.ErrUserNotFound
		}
		return u, fmt.Errorf("UserRepository.FindByID: %w", err)
	}

	return u, nil
}

func (r *UserRepository) Update(ctx context.Context, user domain.User) (domain.User, error) {
	query := `
		UPDATE users
		SET username = $1, email = $2, password_hash = $3, refresh_token = $4, updated_at = NOW()
		WHERE id = $5
		RETURNING id, username, email, password_hash, refresh_token, created_at, updated_at
	`
	var u domain.User
	err := r.db.QueryRow(ctx, query, user.Username, user.Email, user.PasswordHash, user.RefreshToken, user.ID).Scan(
		&u.ID, &u.Username, &u.Email, &u.PasswordHash, &u.RefreshToken, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return u, domain.ErrUserNotFound
		}
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return u, domain.ErrUserAlreadyExists
		}
		return u, fmt.Errorf("UserRepository.Update: %w", err)
	}
	return u, nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	query := `
		SELECT id, username, email, password_hash, refresh_token, created_at, updated_at
		FROM users
		WHERE email = $1
	`
	var u domain.User
	err := r.db.QueryRow(ctx, query, email).Scan(
		&u.ID, &u.Username, &u.Email, &u.PasswordHash, &u.RefreshToken, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return u, domain.ErrUserNotFound
		}
		return u, fmt.Errorf("UserRepository.FindByEmail: %w", err)
	}
	return u, nil
}

func (r *UserRepository) FindByUsername(ctx context.Context, username string) (domain.User, error) {
	query := `
		SELECT id, username, email, password_hash, refresh_token, created_at, updated_at
		FROM users
		WHERE username = $1
	`
	var u domain.User
	err := r.db.QueryRow(ctx, query, username).Scan(
		&u.ID, &u.Username, &u.Email, &u.PasswordHash, &u.RefreshToken, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return u, domain.ErrUserNotFound
		}
		return u, fmt.Errorf("UserRepository.FindByUsername: %w", err)
	}
	return u, nil
}
