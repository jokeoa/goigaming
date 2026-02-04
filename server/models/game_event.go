package models

import "time"

type GameEvent struct {
	Type      string    `json:"type"`
	GameID    string    `json:"game_id"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`