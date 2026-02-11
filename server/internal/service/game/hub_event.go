package game

import (
	"github.com/google/uuid"
	"github.com/jokeoa/goigaming/internal/core/domain"
	"github.com/shopspring/decimal"
)

type EventType int

const (
	EventPlayerAction EventType = iota
	EventPlayerJoin
	EventPlayerLeave
	EventStartHand
	EventTimerExpired
	EventShutdown
)

type HubEvent struct {
	Type     EventType
	PlayerID uuid.UUID
	Action   domain.ActionType
	Amount   decimal.Decimal
	SeatNum  int
	BuyIn    decimal.Decimal
	Username string
	ResultCh chan HubResult
}

type HubResult struct {
	Err  error
	Data any
}
