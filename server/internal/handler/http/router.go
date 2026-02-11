package http

import (
	"github.com/gin-gonic/gin"
	"github.com/jokeoa/goigaming/internal/core/ports"
	"github.com/jokeoa/goigaming/internal/handler/http/middleware"
)

func NewRouter(
	authService ports.AuthService,
	authHandler *AuthHandler,
	userHandler *UserHandler,
	walletHandler *WalletHandler,
) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())

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

		betting := protected.Group("/betting")  
	{
		    betting.GET("/events", bettingHandler.GetEvents)
		    betting.POST("/bets", bettingHandler.PlaceBet)
		    betting.GET("/bets", bettingHandler.GetBets)
		    betting.DELETE("/bets/:id", bettingHandler.CancelBet)
	    }
		
	}

	return r
}
