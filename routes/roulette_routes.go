package routes

import (
    "net/http"

    "goigaming/handler"
)

func RegisterRouletteRoutes(mux *http.ServeMux, h *handler.RouletteHandler) {
    mux.HandleFunc("/api/roulette/tables", h.GetTables)
    mux.HandleFunc("/api/roulette/bet", h.PlaceBet)
    mux.HandleFunc("/api/roulette/history/", h.GetHistory)
}
