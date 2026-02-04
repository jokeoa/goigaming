package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/jokeoa/igaming/models"
	"github.com/jokeoa/igaming/services"
	"github.com/olahol/melody"
)

type WebSocketHandler struct {
	melody           *melody.Melody
	broadcastService *services.BroadcastService
}

func NewWebSocketHandler(m *melody.Melody, bs *services.BroadcastService) *WebSocketHandler {
	return &WebSocketHandler{
		melody:           m,
		broadcastService: bs,
	}
}

func (h *WebSocketHandler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	h.melody.HandleRequest(w, r)
}

func (h *WebSocketHandler) SendEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var event models.GameEvent
	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	event.Timestamp = time.Now()
	eventJSON, _ := json.Marshal(event)
	h.melody.Broadcast(eventJSON)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "event sent"})
}

func (h *WebSocketHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	stats := map[string]int{
		"connected_clients": h.broadcastService.GetConnectedClients(),
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

func (h *WebSocketHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status": "ok"}`))
}

func (h *WebSocketHandler) SetupHandlers(m *melody.Melody) {
	m.HandleConnect(func(s *melody.Session) {
		log.Println("Client connected")
	})

	m.HandleDisconnect(func(s *melody.Session) {
		log.Println("Client disconnected")
	})
}