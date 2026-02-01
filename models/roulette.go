package models

import (
    "time"

    "github.com/google/uuid"
)

type RouletteTable struct {
    ID        uuid.UUID `json:"id" db:"id"`
    Name      string    `json:"name" db:"name"`
    MinBet    float64   `json:"min_bet" db:"min_bet"`
    MaxBet    float64   `json:"max_bet" db:"max_bet"`
    Status    string    `json:"status" db:"status"`
    CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type RouletteRound struct {
    ID            uuid.UUID  `json:"id" db:"id"`
    TableID       uuid.UUID  `json:"table_id" db:"table_id"`
    RoundNumber   int        `json:"round_number" db:"round_number"`
    Result        *int       `json:"result" db:"result"`
    ResultColor   *string    `json:"result_color" db:"result_color"`
    SeedHash      *string    `json:"seed_hash" db:"seed_hash"`
    SeedRevealed  *string    `json:"seed_revealed" db:"seed_revealed"`
    BettingEndsAt *time.Time `json:"betting_ends_at" db:"betting_ends_at"`
    CreatedAt     time.Time  `json:"created_at" db:"created_at"`
    SettledAt     *time.Time `json:"settled_at" db:"settled_at"`
}

type RouletteBet struct {
    ID        uuid.UUID `json:"id" db:"id"`
    RoundID   uuid.UUID `json:"round_id" db:"round_id"`
    UserID    uuid.UUID `json:"user_id" db:"user_id"`
    BetType   string    `json:"bet_type" db:"bet_type"`
    BetValue  string    `json:"bet_value" db:"bet_value"`
    Amount    float64   `json:"amount" db:"amount"`
    Payout    float64   `json:"payout" db:"payout"`
    Status    string    `json:"status" db:"status"`
    CreatedAt time.Time `json:"created_at" db:"created_at"`
}
