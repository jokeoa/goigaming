package ports

import (
	"context"

	"github.com/google/uuid"
	"github.com/jokeoa/igaming/internal/core/domain"
	"github.com/shopspring/decimal"
)

type UserRepository interface {
	Create(ctx context.Context, user domain.User) (domain.User, error)
	FindByID(ctx context.Context, id uuid.UUID) (domain.User, error)
	FindByEmail(ctx context.Context, email string) (domain.User, error)
	FindByUsername(ctx context.Context, username string) (domain.User, error)
	Update(ctx context.Context, user domain.User) (domain.User, error)
}

type WalletRepository interface {
	Create(ctx context.Context, wallet domain.Wallet) (domain.Wallet, error)
	FindByUserID(ctx context.Context, userID uuid.UUID) (domain.Wallet, error)
	UpdateBalance(ctx context.Context, userID uuid.UUID, newBalance decimal.Decimal, currentVersion int) (domain.Wallet, error)
}

type TransactionRepository interface {
	Create(ctx context.Context, tx domain.Transaction) (domain.Transaction, error)
	FindByWalletID(ctx context.Context, filter domain.TransactionFilter) ([]domain.Transaction, error)
}
