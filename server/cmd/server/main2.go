package main

import (
    "log"
    "log/slog"
    "net/http"
    "time"
    
    // ... другие импорты ...
    "github.com/jokeoa/goigaming/internal/repository/redis"
    "github.com/jokeoa/goigaming/internal/handler/http/middleware"
)

func main() {
    // ... существующий код ...
    
    // Инициализация Redis
    redisClient, err := redis.NewClient(redis.Config{
        Host:     cfg.Redis.Host,
        Port:     cfg.Redis.Port,
        Password: cfg.Redis.Password,
        DB:       cfg.Redis.DB,
    })
    if err != nil {
        log.Fatal("Failed to connect to Redis:", err)
    }
    defer redisClient.Close()
    
    // Создаем cache
    cache := redis.NewCache(redisClient)
    
    // Создаем rate limiter
    rateLimiter := middleware.NewRateLimiter(
        redisClient,
        100,              // 100 запросов
        time.Minute,      // в минуту
    )
    
    // Добавляем middleware в router
    router := chi.NewRouter()
    
    // Global middleware
    router.Use(middleware.RequestLogger(logger))
    router.Use(rateLimiter.Middleware())
    router.Use(middleware.CORS)
    
    // Специальный rate limit для auth endpoints
    router.Route("/api/v1/auth", func(r chi.Router) {
        r.Use(rateLimiter.PerEndpoint("login", 10, time.Minute))
        // ... auth routes ...
    })
    
    // ... остальные routes ...
}
