package main

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/jokeoa/goigaming/internal/config"
	"github.com/jokeoa/goigaming/internal/core/ports"
	handler "github.com/jokeoa/goigaming/internal/handler/http"
	wsHandler "github.com/jokeoa/goigaming/internal/handler/ws"
	"github.com/jokeoa/goigaming/internal/repository/postgres"
	"github.com/jokeoa/goigaming/internal/service/game"
	rouletteService "github.com/jokeoa/goigaming/internal/service/roulette"
	"github.com/jokeoa/goigaming/repository"
	authService "github.com/jokeoa/goigaming/internal/service/auth"
	userService "github.com/jokeoa/goigaming/internal/service/user"
	walletService "github.com/jokeoa/goigaming/internal/service/wallet"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	pool, err := postgres.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	userRepo := postgres.NewUserRepository(pool)
	walletRepo := postgres.NewWalletRepository(pool)
	txRepo := postgres.NewTransactionRepository(pool)
	pokerTableRepo := postgres.NewPokerTableRepository(pool)
	pokerPlayerRepo := postgres.NewPokerPlayerRepository(pool)
	pokerHandRepo := postgres.NewPokerHandRepository(pool)
	rouletteTableRepo := repository.NewRouletteTableRepository(pool)
	rouletteTableRepoNew := postgres.NewRouletteTableRepo(pool)
	rouletteRoundRepo := postgres.NewRouletteRoundRepo(pool)
	rouletteBetRepo := postgres.NewRouletteBetRepo(pool)

	authSvc := authService.NewService(
		pool,
		userRepo,
		func(db postgres.DBTX) ports.UserRepository {
			return postgres.NewUserRepository(db)
		},
		func(db postgres.DBTX) ports.WalletRepository {
			return postgres.NewWalletRepository(db)
		},
		cfg.JWTSecret,
		cfg.JWTTokenTTL,
	)
	userSvc := userService.NewService(userRepo)
	walletSvc := walletService.NewService(
		pool,
		walletRepo,
		txRepo,
		func(db postgres.DBTX) ports.WalletRepository {
			return postgres.NewWalletRepository(db)
		},
		func(db postgres.DBTX) ports.TransactionRepository {
			return postgres.NewTransactionRepository(db)
		},
	)

	wsHub := wsHandler.NewHub(slog.Default())
	rngSvc := &game.SimpleRNGService{}
	hubManager := game.NewHubManager(
		ctx,
		30*time.Second,
		wsHub,
		walletSvc,
		rngSvc,
		pokerHandRepo,
		pokerPlayerRepo,
		slog.Default(),
	)
	pokerSvc := game.NewService(
		pool,
		pokerTableRepo,
		pokerPlayerRepo,
		pokerHandRepo,
		walletSvc,
		userSvc,
		hubManager,
		func(db postgres.DBTX) ports.PokerPlayerRepository {
			return postgres.NewPokerPlayerRepository(db)
		},
	)

	rouletteSvc := rouletteService.NewService(
		pool,
		walletSvc,
		rouletteTableRepoNew,
		rouletteRoundRepo,
		rouletteBetRepo,
		func(db postgres.DBTX) ports.RouletteBetRepository {
			return postgres.NewRouletteBetRepo(db)
		},
	)

	authHandler := handler.NewAuthHandler(authSvc)
	userHandler := handler.NewUserHandler(userSvc)
	walletHandler := handler.NewWalletHandler(walletSvc)
	adminHandler := handler.NewAdminHandler(pokerTableRepo, rouletteTableRepo)
	pokerHandler := handler.NewPokerHandler(pokerSvc)
	rouletteHandler := handler.NewRouletteHandler(rouletteSvc)
	ws := wsHandler.NewHandler(wsHub, authSvc, slog.Default())

	router := handler.NewRouter(authSvc, authHandler, userHandler, walletHandler, adminHandler, pokerHandler, rouletteHandler, ws)

	srv := &http.Server{
		Addr:              ":" + cfg.ServerPort,
		Handler:           router,
		ReadHeaderTimeout: 10 * time.Second,
	}

	go func() {
		log.Printf("server starting on :%s", cfg.ServerPort)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server error: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("server shutdown error: %v", err)
	}

	log.Println("server stopped")
}
