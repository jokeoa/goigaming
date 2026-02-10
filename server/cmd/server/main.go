package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/jokeoa/goigaming/internal/config"
	"github.com/jokeoa/goigaming/internal/core/ports"
	handler "github.com/jokeoa/goigaming/internal/handler/http"
	"github.com/jokeoa/goigaming/internal/repository/postgres"
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

	authHandler := handler.NewAuthHandler(authSvc)
	userHandler := handler.NewUserHandler(userSvc)
	walletHandler := handler.NewWalletHandler(walletSvc)

	router := handler.NewRouter(authSvc, authHandler, userHandler, walletHandler)

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
