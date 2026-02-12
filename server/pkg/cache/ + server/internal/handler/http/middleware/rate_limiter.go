package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
)

// RateLimiterConfig holds rate limiter configuration
type RateLimiterConfig struct {
	RequestsPerWindow int           // Number of requests allowed
	Window            time.Duration // Time window
	KeyPrefix         string        // Redis key prefix
}

// RateLimiter implements sliding window rate limiting using Redis
type RateLimiter struct {
	client *redis.Client
	config RateLimiterConfig
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(client *redis.Client, config RateLimiterConfig) *RateLimiter {
	if config.KeyPrefix == "" {
		config.KeyPrefix = "ratelimit"
	}
	return &RateLimiter{
		client: client,
		config: config,
	}
}

// Middleware returns HTTP middleware that enforces rate limiting
func (rl *RateLimiter) Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Use IP address as identifier
			identifier := r.RemoteAddr

			// Check rate limit
			allowed, err := rl.Allow(r.Context(), identifier)
			if err != nil {
				http.Error(w, "Rate limiter error", http.StatusInternalServerError)
				return
			}

			if !allowed {
				w.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d", rl.config.RequestsPerWindow))
				w.Header().Set("X-RateLimit-Window", rl.config.Window.String())
				http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// Allow checks if request is allowed based on sliding window algorithm
func (rl *RateLimiter) Allow(ctx context.Context, identifier string) (bool, error) {
	key := fmt.Sprintf("%s:%s", rl.config.KeyPrefix, identifier)
	now := time.Now()
	windowStart := now.Add(-rl.config.Window)

	pipe := rl.client.Pipeline()

	// Remove old entries outside the window
	pipe.ZRemRangeByScore(ctx, key, "0", fmt.Sprintf("%d", windowStart.UnixNano()))

	// Count current requests in window
	pipe.ZCard(ctx, key)

	// Add current request
	pipe.ZAdd(ctx, key, redis.Z{
		Score:  float64(now.UnixNano()),
		Member: fmt.Sprintf("%d", now.UnixNano()),
	})

	// Set expiration on key
	pipe.Expire(ctx, key, rl.config.Window+time.Second)

	cmds, err := pipe.Exec(ctx)
	if err != nil {
		return false, fmt.Errorf("pipeline exec failed: %w", err)
	}

	// Get count from second command (ZCard)
	count := cmds[1].(*redis.IntCmd).Val()

	// Allow if count is less than limit
	return count < int64(rl.config.RequestsPerWindow), nil
}

// Reset clears rate limit for identifier
func (rl *RateLimiter) Reset(ctx context.Context, identifier string) error {
	key := fmt.Sprintf("%s:%s", rl.config.KeyPrefix, identifier)
	return rl.client.Del(ctx, key).Err()
}

// GetUsage returns current usage for identifier
func (rl *RateLimiter) GetUsage(ctx context.Context, identifier string) (int64, error) {
	key := fmt.Sprintf("%s:%s", rl.config.KeyPrefix, identifier)
	now := time.Now()
	windowStart := now.Add(-rl.config.Window)

	// Remove old entries
	err := rl.client.ZRemRangeByScore(ctx, key, "0", fmt.Sprintf("%d", windowStart.UnixNano())).Err()
	if err != nil {
		return 0, err
	}

	// Count current requests
	count, err := rl.client.ZCard(ctx, key).Result()
	if err != nil {
		return 0, err
	}

	return count, nil
}
