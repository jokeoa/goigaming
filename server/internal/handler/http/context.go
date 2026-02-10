package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jokeoa/goigaming/internal/handler/http/middleware"
)

func getUserID(c *gin.Context) (uuid.UUID, bool) {
	val, exists := c.Get(middleware.ContextKeyUserID)
	if !exists {
		c.AbortWithStatusJSON(http.StatusUnauthorized, Response{
			Success: false,
			Error:   "unauthorized",
		})
		return uuid.UUID{}, false
	}

	userID, ok := val.(uuid.UUID)
	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "internal server error",
		})
		return uuid.UUID{}, false
	}

	return userID, true
}
