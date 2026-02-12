package ports

import (
    "context"
    "time"
    
    "github.com/google/uuid"
    "github.com/jokeoa/goigaming/internal/core/domain"
    "github.com/shopspring/decimal"
)

// ... существующие UserRepository, WalletRepository, TransactionRepository ...

// GameSessionRepository - репозиторий игровых сессий
type GameSessionRepository interface {
    Create(ctx context.Context, session *domain.GameSession) error
    GetByID(ctx context.Context, id uuid.UUID) (*domain.GameSession, error)
    UpdateState(ctx context.Context, id uuid.UUID, state []byte) error
    ListActive(ctx context.Context, gameType string) ([]*domain.GameSession, error)
    Close(ctx context.Context, id uuid.UUID, result *domain.SessionResult) error
}

// SportEventRepository - репозиторий спортивных событий
type SportEventRepository interface {
    Create(ctx context.Context, event *domain.SportEvent) (*domain.SportEvent, error)
    GetByID(ctx context.Context, id uuid.UUID) (*domain.SportEvent, error)
    List(ctx context.Context, filter domain.SportEventFilter) ([]*domain.SportEvent, int64, error)
    Update(ctx context.Context, event *domain.SportEvent) (*domain.SportEvent, error)
    Delete(ctx context.Context, id uuid.UUID) error
    UpdateStatus(ctx context.Context, id uuid.UUID, status string) error
    SetResult(ctx context.Context, id uuid.UUID, homeScore, awayScore int) error
}

// SportBetRepository - репозиторий ставок на спорт
type SportBetRepository interface {
    Create(ctx context.Context, bet *domain.SportBet) (*domain.SportBet, error)
    GetByID(ctx context.Context, id uuid.UUID) (*domain.SportBet, error)
    ListByUser(ctx context.Context, userID uuid.UUID, filter domain.BetFilter) ([]*domain.SportBet, error)
    ListByEvent(ctx context.Context, eventID uuid.UUID) ([]*domain.SportBet, error)
    Cancel(ctx context.Context, id uuid.UUID) error
    Settle(ctx context.Context, id uuid.UUID, status string) error
}

// RouletteRoundRepository - репозиторий раундов рулетки
type RouletteRoundRepository interface {
    Create(ctx context.Context, round *domain.RouletteRound) (*domain.RouletteRound, error)
    GetByID(ctx context.Context, id uuid.UUID) (*domain.RouletteRound, error)
    GetCurrent(ctx context.Context, tableID uuid.UUID) (*domain.RouletteRound, error)
    SetResult(ctx context.Context, id uuid.UUID, result int, color string) error
    GetHistory(ctx context.Context, tableID uuid.UUID, limit int) ([]*domain.RouletteRound, error)
}

// RouletteBetRepository - репозиторий ставок на рулетку
type RouletteBetRepository interface {
    Create(ctx context.Context, bet *domain.RouletteBet) (*domain.RouletteBet, error)
    ListByRound(ctx context.Context, roundID uuid.UUID) ([]*domain.RouletteBet, error)
    UpdatePayout(ctx context.Context, id uuid.UUID, payout decimal.Decimal, status string) error
    BatchUpdatePayouts(ctx context.Context, bets []domain.BetPayout) error
}
