package redis

import (
    "context"
    "fmt"
    
    "github.com/redis/go-redis/v9"
)

type Config struct {
    Host     string
    Port     string
    Password string
    DB       int
}

func NewClient(cfg Config) (*redis.Client, error) {
    client := redis.NewClient(&redis.Options{
        Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
        Password: cfg.Password,
        DB:       cfg.DB,
    })
    
    // Проверка подключения
    ctx := context.Background()
    if err := client.Ping(ctx).Err(); err != nil {
        return nil, fmt.Errorf("failed to connect to redis: %w", err)
    }
    
    return client, nil
}
