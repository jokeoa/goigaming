package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jokeoa/igaming/services"
)

type RouletteHandler struct {
	rouletteService *services.RouletteService
}

func NewRouletteHandler(rouletteService *services.RouletteService) *RouletteHandler {
	return &RouletteHandler{
		rouletteService: rouletteService,
	}
}

type PlaceBetRequest struct {
	TableID  string  `json:"table_id"`
	BetType  string  `json:"bet_type"`
	BetValue string  `json:"bet_value"`
	Amount   float64 `json:"amount"`
}

func (h *RouletteHandler) PlaceBet(w http.ResponseWriter, r *http.Request) {
	var req PlaceBetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request"})
		return
	}

	userID := uuid.MustParse("00000000-0000-0000-0000-000000000001")

	tableID, err := uuid.Parse(req.TableID)
	if err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid table_id"})
		return
	}

	bet, err := h.rouletteService.PlaceBet(r.Context(), userID, tableID, req.BetType, req.BetValue, req.Amount)
	if err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	respondJSON(w, http.StatusCreated, bet)
}

type SpinRequest struct {
	TableID string `json:"table_id"`
}

func (h *RouletteHandler) Spin(w http.ResponseWriter, r *http.Request) {
	var req SpinRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request"})
		return
	}

	tableID, err := uuid.Parse(req.TableID)
	if err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid table_id"})
		return
	}

	round, err := h.rouletteService.Spin(r.Context(), tableID)
	if err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	respondJSON(w, http.StatusOK, round)
}

func (h *RouletteHandler) GetHistory(w http.ResponseWriter, r *http.Request) {
	tableIDStr := r.URL.Query().Get("table_id")
	if tableIDStr == "" {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "table_id required"})
		return
	}

	tableID, err := uuid.Parse(tableIDStr)
	if err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid table_id"})
		return
	}

	limitStr := r.URL.Query().Get("limit")
	limit := 20
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	history, err := h.rouletteService.GetHistory(r.Context(), tableID, limit)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"history": history,
		"count":   len(history),
	})
}

func (h *RouletteHandler) RegisterRoutes(r chi.Router) {
	r.Route("/api/v1/roulette", func(r chi.Router) {
		r.Post("/bets", h.PlaceBet)
		r.Post("/spin", h.Spin)
		r.Get("/history", h.GetHistory)
	})
}

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
