package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jokeoa/igaming/internal/core/domain"
	"github.com/shopspring/decimal"
)

type WalletRepository struct {
	db DBTX
}

func NewWalletRepository(db DBTX) *WalletRepository {
	return &WalletRepository{db: db}
}

func (r *WalletRepository) Create(ctx context.Context, wallet domain.Wallet) (domain.Wallet, error) {
	query := `
		INSERT INTO wallets (user_id, balance, version)
		VALUES ($1, $2, 1)
		RETURNING user_id, balance, version, updated_at
	`

	var w domain.Wallet
	err := r.db.QueryRow(ctx, query, wallet.UserID, wallet.Balance).Scan(
		&w.UserID, &w.Balance, &w.Version, &w.UpdatedAt,
	)
	if err != nil {
		return w, fmt.Errorf("WalletRepository.Create: %w", err)
	}

	return w, nil
}

func (r *WalletRepository) FindByUserID(ctx context.Context, userID uuid.UUID) (domain.Wallet, error) {
	query := `
		SELECT user_id, balance, version, updated_at
		FROM wallets
		WHERE user_id = $1
	`

	var w domain.Wallet
	err := r.db.QueryRow(ctx, query, userID).Scan(
		&w.UserID, &w.Balance, &w.Version, &w.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return w, domain.ErrWalletNotFound
		}
		return w, fmt.Errorf("WalletRepository.FindByUserID: %w", err)
	}

	return w, nil
}

func (r *WalletRepository) UpdateBalance(ctx context.Context, userID uuid.UUID, newBalance decimal.Decimal, currentVersion int) (domain.Wallet, error) {
	query := `
		UPDATE wallets
		SET balance = $1, version = version + 1, updated_at = NOW()
		WHERE user_id = $2 AND version = $3
		RETURNING user_id, balance, version, updated_at
	`

	var w domain.Wallet
	err := r.db.QueryRow(ctx, query, newBalance, userID, currentVersion).Scan(
		&w.UserID, &w.Balance, &w.Version, &w.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return w, domain.ErrOptimisticLock
		}
		return w, fmt.Errorf("WalletRepository.UpdateBalance: %w", err)
	}

	return w, nil
}
