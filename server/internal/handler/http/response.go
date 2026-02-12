package http

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jokeoa/goigaming/internal/core/domain"
)

type Response struct {
	Success bool   `json:"success"`
	Data    any    `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
}

func respondSuccess(c *gin.Context, status int, data any) {
	c.JSON(status, Response{
		Success: true,
		Data:    data,
	})
}

func respondError(c *gin.Context, err error) {
	status, message := mapDomainError(err)
	c.JSON(status, Response{
		Success: false,
		Error:   message,
	})
}

func mapDomainError(err error) (int, string) {
	switch {
	case errors.Is(err, domain.ErrUserNotFound):
		return http.StatusNotFound, "user not found"
	case errors.Is(err, domain.ErrUserAlreadyExists):
		return http.StatusConflict, "user already exists"
	case errors.Is(err, domain.ErrInvalidCredentials):
		return http.StatusUnauthorized, "invalid credentials"
	case errors.Is(err, domain.ErrInvalidToken):
		return http.StatusUnauthorized, "invalid or expired token"
	case errors.Is(err, domain.ErrWalletNotFound):
		return http.StatusNotFound, "wallet not found"
	case errors.Is(err, domain.ErrInsufficientFunds):
		return http.StatusUnprocessableEntity, "insufficient funds"
	case errors.Is(err, domain.ErrInvalidAmount):
		return http.StatusBadRequest, "amount must be greater than zero"
	case errors.Is(err, domain.ErrTableNotFound):
		return http.StatusNotFound, "table not found"
	case errors.Is(err, domain.ErrTableFull):
		return http.StatusConflict, "table is full"
	case errors.Is(err, domain.ErrPlayerNotFound):
		return http.StatusNotFound, "player not found"
	case errors.Is(err, domain.ErrPlayerAlreadySeated):
		return http.StatusConflict, "player already seated at this table"
	case errors.Is(err, domain.ErrHandNotFound):
		return http.StatusNotFound, "hand not found"
	case errors.Is(err, domain.ErrNotPlayerTurn):
		return http.StatusBadRequest, "not your turn"
	case errors.Is(err, domain.ErrInvalidAction):
		return http.StatusBadRequest, "invalid action"
	case errors.Is(err, domain.ErrInvalidBetAmount):
		return http.StatusBadRequest, "invalid bet amount"
	case errors.Is(err, domain.ErrInsufficientStack):
		return http.StatusBadRequest, "insufficient stack"
	case errors.Is(err, domain.ErrGameNotStarted):
		return http.StatusBadRequest, "game not started"
	case errors.Is(err, domain.ErrInvalidBuyIn):
		return http.StatusBadRequest, "invalid buy-in amount"
	case errors.Is(err, domain.ErrSeatTaken):
		return http.StatusConflict, "seat is taken"
	case errors.Is(err, domain.ErrMinPlayersRequired):
		return http.StatusBadRequest, "minimum 2 players required"
	case errors.Is(err, domain.ErrForbidden):
		return http.StatusForbidden, "forbidden"
	case errors.Is(err, domain.ErrRoundNotFound):
		return http.StatusNotFound, "round not found"
	case errors.Is(err, domain.ErrBettingClosed):
		return http.StatusUnprocessableEntity, "betting is closed"
	case errors.Is(err, domain.ErrInvalidBetType):
		return http.StatusBadRequest, "invalid bet type"
	case errors.Is(err, domain.ErrBetAmountOutOfRange):
		return http.StatusBadRequest, "bet amount out of range"
	case errors.Is(err, domain.ErrTableNotActive):
		return http.StatusBadRequest, "table is not active"
	default:
		return http.StatusInternalServerError, "internal server error"
	}
}
