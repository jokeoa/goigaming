package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/olahol/melody"
)

// GameEvent represents a game event to broadcast
type GameEvent struct {
	Type      string    `json:"type"`
	GameID    string    `json:"game_id"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

func main() {
	m := melody.New()

	http.HandleFunc("/", healthHandler)

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		m.HandleRequest(w, r)
	})

	http.HandleFunc("/api/game/event", func(w http.ResponseWriter, r *http.Request) {
		sendEventHandler(m, w, r)
	})

	m.HandleConnect(func(s *melody.Session) {
		log.Println("Client connected")
	})

	m.HandleDisconnect(func(s *melody.Session) {
		log.Println("Client disconnected")
	})

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// sendEventHandler receives a game event via POST and broadcasts it
func sendEventHandler(m *melody.Melody, w http.ResponseWriter, r *http.Request) {
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

	eventJSON, err := json.Marshal(event)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	m.Broadcast(eventJSON)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "event sent"})
}

// healthHandler returns server status
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"status": "ok"}`)
}
