package redis

import (
    "context"
    "encoding/json"
    "fmt"
    "time"
    
    "github.com/redis/go-redis/v9"
)

type Cache struct {
    client *redis.Client
}

func NewCache(client *redis.Client) *Cache {
    return &Cache{client: client}
}

// Set сохраняет значение в кэш с TTL
func (c *Cache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
    data, err := json.Marshal(value)
    if err != nil {
        return fmt.Errorf("failed to marshal value: %w", err)
    }
    
    return c.client.Set(ctx, key, data, ttl).Err()
}

// Get получает значение из кэша
func (c *Cache) Get(ctx context.Context, key string, dest interface{}) error {
    data, err := c.client.Get(ctx, key).Bytes()
    if err == redis.Nil {
        return fmt.Errorf("key not found")
    }
    if err != nil {
        return fmt.Errorf("failed to get value: %w", err)
    }
    
    if err := json.Unmarshal(data, dest); err != nil {
        return fmt.Errorf("failed to unmarshal value: %w", err)
    }
    
    return nil
}

// Delete удаляет значение из кэша
func (c *Cache) Delete(ctx context.Context, keys ...string) error {
    return c.client.Del(ctx, keys...).Err()
}

// Exists проверяет существование ключа
func (c *Cache) Exists(ctx context.Context, key string) (bool, error) {
    count, err := c.client.Exists(ctx, key).Result()
    if err != nil {
        return false, err
    }
    return count > 0, nil
}

// Increment инкрементирует счетчик
func (c *Cache) Increment(ctx context.Context, key string) (int64, error) {
    return c.client.Incr(ctx, key).Result()
}

// Expire устанавливает TTL для ключа
func (c *Cache) Expire(ctx context.Context, key string, ttl time.Duration) error {
    return c.client.Expire(ctx, key, ttl).Err()
}
