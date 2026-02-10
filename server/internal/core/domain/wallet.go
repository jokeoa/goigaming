package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Wallet struct {
	UserID    uuid.UUID       `json:"user_id"`
	Balance   decimal.Decimal `json:"balance"`
	Version   int             `json:"version"`
	UpdatedAt time.Time       `json:"updated_at"`
}
