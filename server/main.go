package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// GameEvent represents a game event to broadcast
type GameEvent struct {
	Type      string    `json:"type"`
	GameID    string    `json:"game_id"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

// Hub manages WebSocket client connections
type Hub struct {
	clients   map[*websocket.Conn]bool
	broadcast chan GameEvent
	mu        sync.Mutex
}

func NewHub() *Hub {
	return &Hub{
		clients:   make(map[*websocket.Conn]bool),
		broadcast: make(chan GameEvent),
	}
}

// Run listens for events and broadcasts them to all connected clients
func (h *Hub) Run() {
	for event := range h.broadcast {
		h.mu.Lock()
		for client := range h.clients {
			err := client.WriteJSON(event)
			if err != nil {
				client.Close()
				delete(h.clients, client)
			}
		}
		h.mu.Unlock()
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// wsHandler upgrades HTTP connection to WebSocket
func (h *Hub) wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}

	h.mu.Lock()
	h.clients[conn] = true
	h.mu.Unlock()

	log.Println("Client connected")

	// Goroutine to keep connection alive and detect disconnects
	go func() {
		defer func() {
			h.mu.Lock()
			delete(h.clients, conn)
			h.mu.Unlock()
			conn.Close()
			log.Println("Client disconnected")
		}()

		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				break
			}
		}
	}()
}

// sendEventHandler receives a game event via POST and broadcasts it
func (h *Hub) sendEventHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var event GameEvent
	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	event.Timestamp = time.Now()
	h.broadcast <- event

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "event sent"})
}

// healthHandler returns server status
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"status": "ok"}`)
}

func main() {
	hub := NewHub()
	go hub.Run()

	http.HandleFunc("/", healthHandler)
	http.HandleFunc("/ws", hub.wsHandler)
	http.HandleFunc("/api/game/event", hub.sendEventHandler)

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
