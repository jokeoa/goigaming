package domain

import (
    "time"
    "github.com/google/uuid"
)

type GameSession struct {
    ID        uuid.UUID `json:"id"`
    GameType  string    `json:"game_type"`
    TableName string    `json:"table_name"`
    Config    []byte    `json:"config"`
    State     []byte    `json:"state"`
    Status    string    `json:"status"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
    ClosedAt  *time.Time `json:"closed_at,omitempty"`
}

type SessionResult struct {
    Winners []uuid.UUID `json:"winners"`
    Losers  []uuid.UUID `json:"losers"`
    Payouts map[uuid.UUID]decimal.Decimal `json:"payouts"`
}
