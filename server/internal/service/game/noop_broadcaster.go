package game

import (
	"github.com/google/uuid"
	"github.com/jokeoa/goigaming/internal/core/domain"
)

type NoopBroadcaster struct{}

func (n *NoopBroadcaster) BroadcastToTable(_ uuid.UUID, _ domain.WSMessage) {}
func (n *NoopBroadcaster) SendToPlayer(_, _ uuid.UUID, _ domain.WSMessage)  {}
