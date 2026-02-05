package main

import (
	"log"
	"net/http"

	"github.com/jokeoa/igaming/handlers"
	"github.com/jokeoa/igaming/services"
	"github.com/jokeoa/igaming/models
	"github.com/olahol/melody"
)

func main() {
	m := melody.New()

	broadcastService := services.NewBroadcastService(m)
	wsHandler := handlers.NewWebSocketHandler(m, broadcastService)
	
	wsHandler.SetupHandlers(m)

	http.HandleFunc("/", wsHandler.HealthCheck)
	http.HandleFunc("/ws", wsHandler.HandleWebSocket)
	http.HandleFunc("/api/game/event", wsHandler.SendEvent)
	http.HandleFunc("/api/game/stats", wsHandler.GetStats)

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}