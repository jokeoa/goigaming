package domain

import (
    "time"
    "github.com/google/uuid"
    "github.com/shopspring/decimal"
)

type SportEvent struct {
    ID          uuid.UUID       `json:"id"`
    SportType   string          `json:"sport_type"`
    League      string          `json:"league"`
    HomeTeam    string          `json:"home_team"`
    AwayTeam    string          `json:"away_team"`
    HomeOdds    decimal.Decimal `json:"home_odds"`
    DrawOdds    *decimal.Decimal `json:"draw_odds,omitempty"`
    AwayOdds    decimal.Decimal `json:"away_odds"`
    EventTime   time.Time       `json:"event_time"`
    Status      string          `json:"status"`
    HomeScore   int             `json:"home_score"`
    AwayScore   int             `json:"away_score"`
    CreatedBy   uuid.UUID       `json:"created_by"`
    CreatedAt   time.Time       `json:"created_at"`
    UpdatedAt   time.Time       `json:"updated_at"`
    SettledAt   *time.Time      `json:"settled_at,omitempty"`
}

type SportEventFilter struct {
    SportType *string
    Status    *string
    FromDate  *time.Time
    ToDate    *time.Time
    Limit     int
    Offset    int
}

type SportBet struct {
    ID           uuid.UUID       `json:"id"`
    EventID      uuid.UUID       `json:"event_id"`
    UserID       uuid.UUID       `json:"user_id"`
    BetType      string          `json:"bet_type"` // home, draw, away
    Odds         decimal.Decimal `json:"odds"`
    Amount       decimal.Decimal `json:"amount"`
    PotentialWin decimal.Decimal `json:"potential_win"`
    Status       string          `json:"status"`
    CreatedAt    time.Time       `json:"created_at"`
    SettledAt    *time.Time      `json:"settled_at,omitempty"`
}

type BetFilter struct {
    Status    *string
    FromDate  *time.Time
    ToDate    *time.Time
    Limit     int
    Offset    int
}
