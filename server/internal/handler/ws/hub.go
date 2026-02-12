package ws

import (
	"encoding/json"
	"log/slog"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/jokeoa/goigaming/internal/core/domain"
)

type client struct {
	ws     *websocket.Conn
	userID uuid.UUID
	mu     sync.Mutex
}

func (c *client) write(data []byte) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.ws.WriteMessage(websocket.TextMessage, data)
}

// Hub manages WebSocket connections grouped by table.
// It implements ports.Broadcaster.
type Hub struct {
	mu     sync.RWMutex
	tables map[uuid.UUID]map[uuid.UUID]*client
	logger *slog.Logger
}

func NewHub(logger *slog.Logger) *Hub {
	return &Hub{
		tables: make(map[uuid.UUID]map[uuid.UUID]*client),
		logger: logger,
	}
}

func (h *Hub) register(tableID, userID uuid.UUID, ws *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.tables[tableID] == nil {
		h.tables[tableID] = make(map[uuid.UUID]*client)
	}

	if existing, ok := h.tables[tableID][userID]; ok {
		existing.ws.Close()
	}

	h.tables[tableID][userID] = &client{ws: ws, userID: userID}
}

func (h *Hub) unregister(tableID, userID uuid.UUID) {
	h.mu.Lock()
	defer h.mu.Unlock()

	conns, ok := h.tables[tableID]
	if !ok {
		return
	}

	delete(conns, userID)
	if len(conns) == 0 {
		delete(h.tables, tableID)
	}
}

func (h *Hub) BroadcastToTable(tableID uuid.UUID, msg domain.WSMessage) {
	data, err := json.Marshal(msg)
	if err != nil {
		h.logger.Error("failed to marshal ws message", "error", err)
		return
	}

	h.mu.RLock()
	targets := make([]*client, 0, len(h.tables[tableID]))
	for _, c := range h.tables[tableID] {
		targets = append(targets, c)
	}
	h.mu.RUnlock()

	for _, c := range targets {
		if err := c.write(data); err != nil {
			h.logger.Debug("broadcast write failed", "user_id", c.userID, "error", err)
		}
	}
}

func (h *Hub) SendToPlayer(tableID, userID uuid.UUID, msg domain.WSMessage) {
	data, err := json.Marshal(msg)
	if err != nil {
		h.logger.Error("failed to marshal ws message", "error", err)
		return
	}

	h.mu.RLock()
	c, ok := h.tables[tableID][userID]
	h.mu.RUnlock()

	if !ok {
		return
	}

	if err := c.write(data); err != nil {
		h.logger.Debug("send to player failed", "user_id", userID, "error", err)
	}
}
