package http

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jokeoa/igaming/internal/core/domain"
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
	default:
		return http.StatusInternalServerError, "internal server error"
	}
}
