package ports

import (
	"context"

	"github.com/google/uuid"
	"github.com/jokeoa/goigaming/internal/core/domain"
)

type GameStateRepository interface {
	SaveTableState(ctx context.Context, tableID uuid.UUID, state map[string]string) error
	GetTableState(ctx context.Context, tableID uuid.UUID) (map[string]string, error)
	DeleteTableState(ctx context.Context, tableID uuid.UUID) error
	SavePlayerState(ctx context.Context, tableID, userID uuid.UUID, state map[string]string) error
	GetPlayerState(ctx context.Context, tableID, userID uuid.UUID) (map[string]string, error)
	PublishEvent(ctx context.Context, tableID uuid.UUID, event []byte) error
	SubscribeEvents(ctx context.Context, tableID uuid.UUID) (<-chan []byte, error)
}

type RNGService interface {
	GenerateServerSeed() (string, error)
	HashSeed(serverSeed string) string
	ShuffleDeck(serverSeed, clientSeed string, nonce int) []domain.Card
	VerifySeed(serverSeed, hash string) bool
}

type Broadcaster interface {
	BroadcastToTable(tableID uuid.UUID, msg domain.WSMessage)
	SendToPlayer(tableID, userID uuid.UUID, msg domain.WSMessage)
}
