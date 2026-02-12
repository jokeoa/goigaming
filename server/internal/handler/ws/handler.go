package ws

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/jokeoa/goigaming/internal/core/ports"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Handler struct {
	hub     *Hub
	authSvc ports.AuthService
	logger  *slog.Logger
}

func NewHandler(hub *Hub, authSvc ports.AuthService, logger *slog.Logger) *Handler {
	return &Handler{
		hub:     hub,
		authSvc: authSvc,
		logger:  logger,
	}
}

func (h *Handler) HandleConnection(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
		return
	}

	gameIDStr := c.Query("game_id")
	if gameIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing game_id"})
		return
	}

	claims, err := h.authSvc.ValidateToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
		return
	}

	tableID, err := uuid.Parse(gameIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid game_id"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		h.logger.Error("websocket upgrade failed", "error", err)
		return
	}

	h.hub.register(tableID, claims.UserID, conn)
	h.logger.Info("ws client connected",
		"user_id", claims.UserID,
		"username", claims.Username,
		"table_id", tableID,
	)

	defer func() {
		h.hub.unregister(tableID, claims.UserID)
		conn.Close()
		h.logger.Info("ws client disconnected",
			"user_id", claims.UserID,
			"table_id", tableID,
		)
	}()

	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			break
		}
	}
}
