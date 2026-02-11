package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type TableStatus string

const (
	TableStatusWaiting TableStatus = "waiting"
	TableStatusActive  TableStatus = "active"
	TableStatusClosed  TableStatus = "closed"
)

type PokerTable struct {
	ID         uuid.UUID       `json:"id"`
	Name       string          `json:"name"`
	SmallBlind decimal.Decimal `json:"small_blind"`
	BigBlind   decimal.Decimal `json:"big_blind"`
	MinBuyIn   decimal.Decimal `json:"min_buy_in"`
	MaxBuyIn   decimal.Decimal `json:"max_buy_in"`
	MaxPlayers int             `json:"max_players"`
	Status     TableStatus     `json:"status"`
	CreatedAt  time.Time       `json:"created_at"`
}
