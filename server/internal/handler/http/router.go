package http

import (
	"github.com/gin-gonic/gin"
	"github.com/jokeoa/goigaming/internal/core/ports"
	"github.com/jokeoa/goigaming/internal/handler/http/middleware"
	wsHandler "github.com/jokeoa/goigaming/internal/handler/ws"
)

func NewRouter(
	authService ports.AuthService,
	authHandler *AuthHandler,
	userHandler *UserHandler,
	walletHandler *WalletHandler,
	adminHandler *AdminHandler,
	pokerHandler *PokerHandler,
	rouletteHandler *RouletteHandler,
	ws *wsHandler.Handler,
) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())

	r.GET("/ws", ws.HandleConnection)

	api := r.Group("/api/v1")

	auth := api.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
	}

	protected := api.Group("")
	protected.Use(middleware.Auth(authService))
	{
		users := protected.Group("/users")
		{
			users.GET("/me", userHandler.GetMe)
		}

		wallet := protected.Group("/wallet")
		{
			wallet.GET("/balance", walletHandler.GetBalance)
			wallet.POST("/deposit", walletHandler.Deposit)
			wallet.POST("/withdraw", walletHandler.Withdraw)
			wallet.GET("/transactions", walletHandler.GetTransactions)
		}

		poker := protected.Group("/poker/tables")
		{
			poker.GET("", pokerHandler.ListTables)
			poker.GET("/:id", pokerHandler.GetTable)
			poker.POST("/:id/join", pokerHandler.JoinTable)
			poker.POST("/:id/leave", pokerHandler.LeaveTable)
			poker.GET("/:id/state", pokerHandler.GetTableState)
		}

		roulette := protected.Group("/roulette")
		{
			tables := roulette.Group("/tables")
			tables.GET("", rouletteHandler.ListTables)
			tables.GET("/:id", rouletteHandler.GetTable)
			tables.GET("/:id/current-round", rouletteHandler.GetCurrentRound)
			tables.POST("/:id/bets", rouletteHandler.PlaceBet)
			tables.GET("/:id/history", rouletteHandler.GetRoundHistory)

			roulette.GET("/bets/me", rouletteHandler.GetMyBets)
			roulette.GET("/rounds/:id", rouletteHandler.GetRound)
		}
	}

	admin := api.Group("/admin")
	admin.Use(middleware.Auth(authService), middleware.Admin())
	{
		pokerTables := admin.Group("/poker-tables")
		{
			pokerTables.POST("", adminHandler.CreatePokerTable)
			pokerTables.GET("", adminHandler.ListPokerTables)
			pokerTables.GET("/:id", adminHandler.GetPokerTable)
			pokerTables.PUT("/:id", adminHandler.UpdatePokerTable)
			pokerTables.DELETE("/:id", adminHandler.DeletePokerTable)
		}

		rouletteTables := admin.Group("/roulette-tables")
		{
			rouletteTables.POST("", adminHandler.CreateRouletteTable)
			rouletteTables.GET("", adminHandler.ListRouletteTables)
			rouletteTables.GET("/:id", adminHandler.GetRouletteTable)
			rouletteTables.PUT("/:id", adminHandler.UpdateRouletteTable)
			rouletteTables.DELETE("/:id", adminHandler.DeleteRouletteTable)
		}
	}

	return r
}
