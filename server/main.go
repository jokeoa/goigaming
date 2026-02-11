package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/olahol/melody"
	"github.com/redis/go-redis/v9"

	"github.com/jokeoa/igaming/handlers"
	appMiddleware "github.com/jokeoa/igaming/middleware"
	"github.com/jokeoa/igaming/repository"
	"github.com/jokeoa/igaming/services"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5432/igaming?sslmode=disable"
	}

	dbPool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer dbPool.Close()

	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		redisURL = "localhost:6379"
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr: redisURL,
	})
	defer redisClient.Close()

	if err := dbPool.Ping(context.Background()); err != nil {
		log.Fatalf("Database ping failed: %v", err)
	}
	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("Redis ping failed: %v", err)
	}

	logger.Info("Connected to database and redis")

	tableRepo := repository.NewRouletteTableRepository(dbPool)
	roundRepo := repository.NewRouletteRoundRepository(dbPool)
	betRepo := repository.NewRouletteBetRepository(dbPool)

	rouletteService := services.NewRouletteService(tableRepo, roundRepo, betRepo)
	adminService := services.NewAdminService(dbPool)

	m := melody.New()
	broadcastService := services.NewBroadcastService(m)
	wsHandler := handlers.NewWebSocketHandler(m, broadcastService)
	wsHandler.SetupHandlers(m)

	rouletteHandler := handlers.NewRouletteHandler(rouletteService)
	adminHandler := handlers.NewAdminHandler(adminService)
	pageHandler := handlers.NewPageHandler(rouletteService, adminService)

	r := chi.NewRouter()

	r.Use(middleware.Recoverer)
	r.Use(appMiddleware.RequestLogger(logger))

	rateLimits := map[string]appMiddleware.RateLimiterConfig{
		"/api/v1/roulette/bets": {RequestsPerMinute: 30, BurstSize: 5},
		"/api/v1/roulette/spin": {RequestsPerMinute: 10, BurstSize: 2},
	}
	defaultRateLimit := appMiddleware.RateLimiterConfig{RequestsPerMinute: 60, BurstSize: 10}
	r.Use(appMiddleware.PerEndpointRateLimiter(redisClient, rateLimits, defaultRateLimit))

	rouletteHandler.RegisterRoutes(r)
	adminHandler.RegisterRoutes(r)
	pageHandler.RegisterRoutes(r)

	r.Get("/", wsHandler.HealthCheck)
	r.Get("/ws", wsHandler.HandleWebSocket)
	r.Post("/api/game/event", wsHandler.SendEvent)
	r.Get("/api/game/stats", wsHandler.GetStats)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	logger.Info("Server starting", slog.String("port", port))
	log.Fatal(http.ListenAndServe(":"+port, r))
}
