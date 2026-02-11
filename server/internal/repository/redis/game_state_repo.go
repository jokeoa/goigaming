package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	goredis "github.com/redis/go-redis/v9"
)

const stateTTL = 24 * time.Hour

type GameStateRepository struct {
	client *goredis.Client
}

func NewGameStateRepository(client *goredis.Client) *GameStateRepository {
	return &GameStateRepository{client: client}
}

func tableKey(tableID uuid.UUID) string {
	return fmt.Sprintf("game:table:%s", tableID.String())
}

func playerKey(tableID, userID uuid.UUID) string {
	return fmt.Sprintf("game:table:%s:player:%s", tableID.String(), userID.String())
}

func channelKey(tableID uuid.UUID) string {
	return fmt.Sprintf("channel:table:%s", tableID.String())
}

func (r *GameStateRepository) SaveTableState(ctx context.Context, tableID uuid.UUID, state map[string]string) error {
	if len(state) == 0 {
		return nil
	}

	pairs := make([]any, 0, len(state)*2)
	for k, v := range state {
		pairs = append(pairs, k, v)
	}

	key := tableKey(tableID)
	if err := r.client.HSet(ctx, key, pairs...).Err(); err != nil {
		return fmt.Errorf("GameStateRepository.SaveTableState: %w", err)
	}

	r.client.Expire(ctx, key, stateTTL)

	return nil
}

func (r *GameStateRepository) GetTableState(ctx context.Context, tableID uuid.UUID) (map[string]string, error) {
	result, err := r.client.HGetAll(ctx, tableKey(tableID)).Result()
	if err != nil {
		return nil, fmt.Errorf("GameStateRepository.GetTableState: %w", err)
	}

	return result, nil
}

func (r *GameStateRepository) DeleteTableState(ctx context.Context, tableID uuid.UUID) error {
	if err := r.client.Del(ctx, tableKey(tableID)).Err(); err != nil {
		return fmt.Errorf("GameStateRepository.DeleteTableState: %w", err)
	}

	return nil
}

func (r *GameStateRepository) SavePlayerState(ctx context.Context, tableID, userID uuid.UUID, state map[string]string) error {
	if len(state) == 0 {
		return nil
	}

	pairs := make([]any, 0, len(state)*2)
	for k, v := range state {
		pairs = append(pairs, k, v)
	}

	key := playerKey(tableID, userID)
	if err := r.client.HSet(ctx, key, pairs...).Err(); err != nil {
		return fmt.Errorf("GameStateRepository.SavePlayerState: %w", err)
	}

	r.client.Expire(ctx, key, stateTTL)

	return nil
}

func (r *GameStateRepository) GetPlayerState(ctx context.Context, tableID, userID uuid.UUID) (map[string]string, error) {
	result, err := r.client.HGetAll(ctx, playerKey(tableID, userID)).Result()
	if err != nil {
		return nil, fmt.Errorf("GameStateRepository.GetPlayerState: %w", err)
	}

	return result, nil
}

func (r *GameStateRepository) PublishEvent(ctx context.Context, tableID uuid.UUID, event []byte) error {
	if err := r.client.Publish(ctx, channelKey(tableID), event).Err(); err != nil {
		return fmt.Errorf("GameStateRepository.PublishEvent: %w", err)
	}

	return nil
}

func (r *GameStateRepository) SubscribeEvents(ctx context.Context, tableID uuid.UUID) (<-chan []byte, error) {
	sub := r.client.Subscribe(ctx, channelKey(tableID))

	ch := make(chan []byte, 64)
	go func() {
		defer close(ch)
		defer sub.Close()

		msgCh := sub.Channel()
		for {
			select {
			case <-ctx.Done():
				return
			case msg, ok := <-msgCh:
				if !ok {
					return
				}
				ch <- []byte(msg.Payload)
			}
		}
	}()

	return ch, nil
}
