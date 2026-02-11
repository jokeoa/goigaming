package ports

import (
	"context"

	"github.com/google/uuid"
	"github.com/jokeoa/goigaming/internal/core/domain"
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

type PokerTableRepository interface {
	Create(ctx context.Context, table domain.PokerTable) (domain.PokerTable, error)
	FindByID(ctx context.Context, id uuid.UUID) (domain.PokerTable, error)
	FindActive(ctx context.Context) ([]domain.PokerTable, error)
	Update(ctx context.Context, table domain.PokerTable) (domain.PokerTable, error)
}

type PokerPlayerRepository interface {
	Create(ctx context.Context, player domain.PokerPlayer) (domain.PokerPlayer, error)
	FindByID(ctx context.Context, id uuid.UUID) (domain.PokerPlayer, error)
	FindByTableID(ctx context.Context, tableID uuid.UUID) ([]domain.PokerPlayer, error)
	FindByTableAndUser(ctx context.Context, tableID, userID uuid.UUID) (domain.PokerPlayer, error)
	UpdateStack(ctx context.Context, playerID uuid.UUID, stack decimal.Decimal) error
	UpdateStatus(ctx context.Context, playerID uuid.UUID, status domain.PlayerStatus) error
	Delete(ctx context.Context, playerID uuid.UUID) error
	CountByTableID(ctx context.Context, tableID uuid.UUID) (int, error)
}

type PokerHandRepository interface {
	Create(ctx context.Context, hand domain.PokerHand) (domain.PokerHand, error)
	FindByID(ctx context.Context, id uuid.UUID) (domain.PokerHand, error)
	Update(ctx context.Context, hand domain.PokerHand) (domain.PokerHand, error)
	FindLatestByTableID(ctx context.Context, tableID uuid.UUID) (domain.PokerHand, error)
	CreateHandPlayer(ctx context.Context, hp domain.PokerHandPlayer) (domain.PokerHandPlayer, error)
	FindHandPlayers(ctx context.Context, handID uuid.UUID) ([]domain.PokerHandPlayer, error)
	UpdateHandPlayer(ctx context.Context, hp domain.PokerHandPlayer) error
	CreateAction(ctx context.Context, action domain.PokerAction) (domain.PokerAction, error)
	FindActionsByHandID(ctx context.Context, handID uuid.UUID) ([]domain.PokerAction, error)
}
