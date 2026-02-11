package domain

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Pot struct {
	Amount      decimal.Decimal `json:"amount"`
	EligibleIDs []uuid.UUID    `json:"eligible_ids"`
}

type HandResult struct {
	HandID        uuid.UUID    `json:"hand_id"`
	Winners       []WinnerInfo `json:"winners"`
	Pots          []Pot        `json:"pots"`
	ShowdownCards map[uuid.UUID][]Card `json:"showdown_cards,omitempty"`
}

type WinnerInfo struct {
	PlayerID uuid.UUID       `json:"player_id"`
	Amount   decimal.Decimal `json:"amount"`
	HandRank string          `json:"hand_rank"`
}
