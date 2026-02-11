package models

import (
	"time"

	"github.com/google/uuid"
)

type SportEvent struct {
	ID          uuid.UUID  `json:"id" db:"id"`
	Sport       string     `json:"sport" db:"sport"`
	League      string     `json:"league" db:"league"`
	HomeTeam    string     `json:"home_team" db:"home_team"`
	AwayTeam    string     `json:"away_team" db:"away_team"`
	StartTime   time.Time  `json:"start_time" db:"start_time"`
	HomeOdds    float64    `json:"home_odds" db:"home_odds"`
	DrawOdds    *float64   `json:"draw_odds,omitempty" db:"draw_odds"`
	AwayOdds    float64    `json:"away_odds" db:"away_odds"`
	Status      string     `json:"status" db:"status"`
	Result      *string    `json:"result,omitempty" db:"result"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	SettledAt   *time.Time `json:"settled_at,omitempty" db:"settled_at"`
}

type SportBet struct {
	ID        uuid.UUID `json:"id" db:"id"`
	EventID   uuid.UUID `json:"event_id" db:"event_id"`
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	BetType   string    `json:"bet_type" db:"bet_type"`
	Amount    float64   `json:"amount" db:"amount"`
	Odds      float64   `json:"odds" db:"odds"`
	Payout    float64   `json:"payout" db:"payout"`
	Status    string    `json:"status" db:"status"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
