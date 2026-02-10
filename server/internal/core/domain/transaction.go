package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Transaction struct {
	ID            uuid.UUID       `json:"id"`
	WalletID      uuid.UUID       `json:"wallet_id"`
	Amount        decimal.Decimal `json:"amount"`
	BalanceAfter  decimal.Decimal `json:"balance_after"`
	ReferenceType string          `json:"reference_type"`
	ReferenceID   *uuid.UUID      `json:"reference_id,omitempty"`
	CreatedAt     time.Time       `json:"created_at"`
}

type TransactionFilter struct {
	WalletID      uuid.UUID
	ReferenceType string
	Limit         int
	Offset        int
}
