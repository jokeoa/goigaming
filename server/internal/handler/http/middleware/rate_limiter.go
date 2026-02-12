package middleware

import (
	"github.com/gin-gonic/gin"
)

// RateLimiter - заглушка (Redis не настроен)
func RateLimiter() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: добавить Redis интеграцию
		c.Next()
	}
}
