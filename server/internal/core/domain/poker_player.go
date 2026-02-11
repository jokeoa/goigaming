package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type PlayerStatus string

const (
	PlayerStatusActive     PlayerStatus = "active"
	PlayerStatusSittingOut PlayerStatus = "sitting_out"
	PlayerStatusAllIn      PlayerStatus = "all_in"
	PlayerStatusFolded     PlayerStatus = "folded"
)

type PokerPlayer struct {
	ID         uuid.UUID       `json:"id"`
	TableID    uuid.UUID       `json:"table_id"`
	UserID     uuid.UUID       `json:"user_id"`
	Username   string          `json:"username"`
	Stack      decimal.Decimal `json:"stack"`
	SeatNumber int             `json:"seat_number"`
	Status     PlayerStatus    `json:"status"`
	JoinedAt   time.Time       `json:"joined_at"`
}
