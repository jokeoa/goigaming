package ports

import (
	"context"

	"github.com/google/uuid"
	"github.com/jokeoa/goigaming/internal/core/domain"
)

type AuthService interface {
	Register(ctx context.Context, username, email, password string) (domain.User, error)
	Login(ctx context.Context, email, password string) (domain.TokenPair, error)
	ValidateToken(token string) (domain.TokenClaims, error)
	RefreshToken(ctx context.Context, refreshToken string) (domain.TokenPair, error)
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

type PokerService interface {
	CreateTable(ctx context.Context, table domain.PokerTable) (domain.PokerTable, error)
	GetTable(ctx context.Context, tableID uuid.UUID) (domain.PokerTable, error)
	ListTables(ctx context.Context) ([]domain.PokerTable, error)
	JoinTable(ctx context.Context, tableID, userID uuid.UUID, seatNumber int, buyIn string) (domain.PokerPlayer, error)
	LeaveTable(ctx context.Context, tableID, userID uuid.UUID) error
	GetTableState(ctx context.Context, tableID uuid.UUID) (domain.WSTableState, error)
}
