package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type GameStage string

const (
	StageWaiting  GameStage = "waiting"
	StagePreflop  GameStage = "preflop"
	StageFlop     GameStage = "flop"
	StageTurn     GameStage = "turn"
	StageRiver    GameStage = "river"
	StageShowdown GameStage = "showdown"
	StageComplete GameStage = "complete"
)

type PokerHand struct {
	ID             uuid.UUID       `json:"id"`
	TableID        uuid.UUID       `json:"table_id"`
	HandNumber     int             `json:"hand_number"`
	Pot            decimal.Decimal `json:"pot"`
	CommunityCards string          `json:"community_cards"`
	Stage          GameStage       `json:"stage"`
	WinnerID       *uuid.UUID      `json:"winner_id,omitempty"`
	StartedAt      time.Time       `json:"started_at"`
	EndedAt        *time.Time      `json:"ended_at,omitempty"`
}

type PokerHandPlayer struct {
	ID         uuid.UUID       `json:"id"`
	HandID     uuid.UUID       `json:"hand_id"`
	PlayerID   uuid.UUID       `json:"player_id"`
	HoleCards  string          `json:"hole_cards"`
	BetAmount  decimal.Decimal `json:"bet_amount"`
	LastAction string          `json:"last_action"`
	IsActive   bool            `json:"is_active"`
}
