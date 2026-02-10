package config

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	ServerPort  string        `env:"SERVER_PORT" envDefault:"8080"`
	DatabaseURL string        `env:"DATABASE_URL,required"`
	JWTSecret   string        `env:"JWT_SECRET,required"`
	JWTTokenTTL time.Duration `env:"JWT_TOKEN_TTL" envDefault:"24h"`
}

func Load() (Config, error) {
	cfg, err := env.ParseAs[Config]()
	if err != nil {
		return Config{}, fmt.Errorf("config.Load: %w", err)
	}
	if len(cfg.JWTSecret) < 32 {
		return Config{}, fmt.Errorf("config.Load: JWT_SECRET must be at least 32 characters")
	}
	return cfg, nil
}
