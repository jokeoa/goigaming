package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
)

type RateLimiterConfig struct {
	RequestsPerMinute int
	BurstSize         int
}

func RateLimiter(redisClient *redis.Client, config RateLimiterConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			clientIP := r.RemoteAddr
			key := fmt.Sprintf("rate_limit:%s:%s", r.URL.Path, clientIP)

			allowed, err := checkRateLimit(ctx, redisClient, key, config)
			if err != nil {
				http.Error(w, "Rate limiter error", http.StatusInternalServerError)
				return
			}

			if !allowed {
				w.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d", config.RequestsPerMinute))
				w.Header().Set("Retry-After", "60")
				http.Error(w, "Too many requests", http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func checkRateLimit(ctx context.Context, client *redis.Client, key string, config RateLimiterConfig) (bool, error) {
	now := time.Now().Unix()
	windowStart := now - 60

	pipe := client.Pipeline()

	pipe.ZRemRangeByScore(ctx, key, "0", fmt.Sprintf("%d", windowStart))

	pipe.ZAdd(ctx, key, redis.Z{
		Score:  float64(now),
		Member: fmt.Sprintf("%d", now),
	})

	pipe.ZCard(ctx, key)

	pipe.Expire(ctx, key, 61*time.Second)

	cmds, err := pipe.Exec(ctx)
	if err != nil {
		return false, err
	}

	count := cmds[2].(*redis.IntCmd).Val()

	return count <= int64(config.RequestsPerMinute), nil
}

func PerEndpointRateLimiter(redisClient *redis.Client, limits map[string]RateLimiterConfig, defaultConfig RateLimiterConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			config, exists := limits[r.URL.Path]
			if !exists {
				config = defaultConfig
			}

			limiter := RateLimiter(redisClient, config)
			limiter(next).ServeHTTP(w, r)
		})
	}
}
