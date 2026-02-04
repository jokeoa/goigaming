package services

import (
	"github.com/jokeoa/igaming/models"
	"github.com/olahol/melody"
)

type BroadcastService struct {
	melody *melody.Melody
}

func NewBroadcastService(m *melody.Melody) *BroadcastService {
	return &BroadcastService{melody: m}
}

func (s *BroadcastService) BroadcastEvent(event models.GameEvent) error {
	return s.melody.Broadcast([]byte(event.Message))
}

func (s *BroadcastService) GetConnectedClients() int {
	return len(s.melody.Sessions())
}