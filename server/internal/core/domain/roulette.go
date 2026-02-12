package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type RouletteTableStatus string

const (
	RouletteTableStatusActive      RouletteTableStatus = "active"
	RouletteTableStatusInactive    RouletteTableStatus = "inactive"
	RouletteTableStatusMaintenance RouletteTableStatus = "maintenance"
)

type RouletteBetStatus string

const (
	RouletteBetStatusPending RouletteBetStatus = "pending"
	RouletteBetStatusWon     RouletteBetStatus = "won"
	RouletteBetStatusLost    RouletteBetStatus = "lost"
)

var validBetTypes = map[string]bool{
	"straight": true,
	"split":    true,
	"street":   true,
	"corner":   true,
	"line":     true,
	"dozen":    true,
	"column":   true,
	"red":      true,
	"black":    true,
	"odd":      true,
	"even":     true,
	"high":     true,
	"low":      true,
}

func IsValidBetType(betType string) bool {
	return validBetTypes[betType]
}

type RouletteTable struct {
	ID        uuid.UUID           `json:"id"`
	Name      string              `json:"name"`
	MinBet    decimal.Decimal     `json:"min_bet"`
	MaxBet    decimal.Decimal     `json:"max_bet"`
	Status    RouletteTableStatus `json:"status"`
	CreatedAt time.Time           `json:"created_at"`
}

type RouletteRound struct {
	ID            uuid.UUID  `json:"id"`
	TableID       uuid.UUID  `json:"table_id"`
	RoundNumber   int        `json:"round_number"`
	Result        *int       `json:"result"`
	ResultColor   *string    `json:"result_color"`
	SeedHash      *string    `json:"seed_hash"`
	SeedRevealed  *string    `json:"seed_revealed"`
	BettingEndsAt *time.Time `json:"betting_ends_at"`
	CreatedAt     time.Time  `json:"created_at"`
	SettledAt     *time.Time `json:"settled_at"`
}

type RouletteBet struct {
	ID        uuid.UUID         `json:"id"`
	RoundID   uuid.UUID         `json:"round_id"`
	UserID    uuid.UUID         `json:"user_id"`
	BetType   string            `json:"bet_type"`
	BetValue  string            `json:"bet_value"`
	Amount    decimal.Decimal   `json:"amount"`
	Payout    decimal.Decimal   `json:"payout"`
	Status    RouletteBetStatus `json:"status"`
	CreatedAt time.Time         `json:"created_at"`
}
