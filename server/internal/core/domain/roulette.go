package domain

import (
    "time"
    "github.com/google/uuid"
    "github.com/shopspring/decimal"
)

type RouletteRound struct {
    ID            uuid.UUID  `json:"id"`
    TableID       uuid.UUID  `json:"table_id"`
    RoundNumber   int        `json:"round_number"`
    Result        *int       `json:"result,omitempty"`
    ResultColor   *string    `json:"result_color,omitempty"`
    SeedHash      string     `json:"seed_hash"`
    SeedRevealed  *string    `json:"seed_revealed,omitempty"`
    BettingEndsAt *time.Time `json:"betting_ends_at,omitempty"`
    CreatedAt     time.Time  `json:"created_at"`
    SettledAt     *time.Time `json:"settled_at,omitempty"`
}

type RouletteBet struct {
    ID        uuid.UUID       `json:"id"`
    RoundID   uuid.UUID       `json:"round_id"`
    UserID    uuid.UUID       `json:"user_id"`
    BetType   string          `json:"bet_type"`
    BetValue  string          `json:"bet_value"`
    Amount    decimal.Decimal `json:"amount"`
    Payout    decimal.Decimal `json:"payout"`
    Status    string          `json:"status"`
    CreatedAt time.Time       `json:"created_at"`
}

type BetPayout struct {
    ID     uuid.UUID
    Payout decimal.Decimal
    Status string
}
