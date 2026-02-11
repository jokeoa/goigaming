package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	AdminRoleID = "admin"
)

func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get(ContextKeyUserID)
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "authentication required",
			})
			return
		}

		uid, ok := userID.(uuid.UUID)
		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"success": false,
				"error":   "access denied: admin only",
			})
			return
		}

		_ = uid
		c.Next()
	}
}
