package middleware

import (
    "context"
    "fmt"
    "net/http"
    "time"
    
    "github.com/redis/go-redis/v9"
)

type RateLimiter struct {
    redis  *redis.Client
    limit  int           // максимум запросов
    window time.Duration // временное окно
}

func NewRateLimiter(redis *redis.Client, limit int, window time.Duration) *RateLimiter {
    return &RateLimiter{
        redis:  redis,
        limit:  limit,
        window: window,
    }
}

// Middleware создает middleware для rate limiting
func (rl *RateLimiter) Middleware() func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // Используем IP адрес как ключ
            key := fmt.Sprintf("rate_limit:%s:%s", r.URL.Path, r.RemoteAddr)
            
            ctx := context.Background()
            
            // Проверяем текущее количество запросов
            count, err := rl.redis.Incr(ctx, key).Result()
            if err != nil {
                // В случае ошибки Redis пропускаем запрос
                next.ServeHTTP(w, r)
                return
            }
            
            // Устанавливаем TTL только для первого запроса
            if count == 1 {
                rl.redis.Expire(ctx, key, rl.window)
            }
            
            // Проверяем лимит
            if count > int64(rl.limit) {
                // Получаем оставшееся время до сброса
                ttl, _ := rl.redis.TTL(ctx, key).Result()
                
                w.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d", rl.limit))
                w.Header().Set("X-RateLimit-Remaining", "0")
                w.Header().Set("X-RateLimit-Reset", fmt.Sprintf("%d", time.Now().Add(ttl).Unix()))
                w.Header().Set("Retry-After", fmt.Sprintf("%d", int(ttl.Seconds())))
                
                http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
                return
            }
            
            // Добавляем заголовки с информацией о лимите
            remaining := rl.limit - int(count)
            w.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d", rl.limit))
            w.Header().Set("X-RateLimit-Remaining", fmt.Sprintf("%d", remaining))
            
            next.ServeHTTP(w, r)
        })
    }
}

// PerEndpoint создает rate limiter для конкретного endpoint
func (rl *RateLimiter) PerEndpoint(endpoint string, limit int, window time.Duration) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            key := fmt.Sprintf("rate_limit:%s:%s", endpoint, r.RemoteAddr)
            
            ctx := context.Background()
            count, err := rl.redis.Incr(ctx, key).Result()
            if err != nil {
                next.ServeHTTP(w, r)
                return
            }
            
            if count == 1 {
                rl.redis.Expire(ctx, key, window)
            }
            
            if count > int64(limit) {
                ttl, _ := rl.redis.TTL(ctx, key).Result()
                w.Header().Set("Retry-After", fmt.Sprintf("%d", int(ttl.Seconds())))
                http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
                return
            }
            
            next.ServeHTTP(w, r)
        })
    }
}
