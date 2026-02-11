package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type ActionType string

const (
	ActionFold  ActionType = "fold"
	ActionCheck ActionType = "check"
	ActionCall  ActionType = "call"
	ActionRaise ActionType = "raise"
	ActionAllIn ActionType = "all_in"
	ActionBet   ActionType = "bet"
)

type PokerAction struct {
	ID          uuid.UUID       `json:"id"`
	HandID      uuid.UUID       `json:"hand_id"`
	PlayerID    uuid.UUID       `json:"player_id"`
	Action      ActionType      `json:"action"`
	Amount      decimal.Decimal `json:"amount"`
	Stage       GameStage       `json:"stage"`
	ActionOrder int             `json:"action_order"`
	CreatedAt   time.Time       `json:"created_at"`
}
