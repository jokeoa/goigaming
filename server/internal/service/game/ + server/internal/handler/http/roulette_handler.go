package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// PlaceBetRequest represents a request to place a roulette bet
type PlaceBetRequest struct {
	BetType string   `json:"bet_type"`
	Numbers []int    `json:"numbers,omitempty"`
	Amount  float64  `json:"amount"`
}

// PlaceBetsRequest for multiple bets at once
type PlaceBetsRequest struct {
	Bets []PlaceBetRequest `json:"bets"`
}

// SpinResponse represents the spin result response
type SpinResponse struct {
	SpinID       string             `json:"spin_id"`
	Number       int                `json:"number"`
	Color        string             `json:"color"`
	Winnings     float64            `json:"winnings"`
	TotalWagered float64            `json:"total_wagered"`
	TotalPaidOut float64            `json:"total_paid_out"`
	SpinHash     string             `json:"spin_hash"`
	BetResults   []BetResult        `json:"bet_results"`
}

// BetResult shows individual bet outcome
type BetResult struct {
	BetType  string  `json:"bet_type"`
	Amount   float64 `json:"amount"`
	Won      bool    `json:"won"`
	Payout   float64 `json:"payout"`
}

// RouletteHandler handles roulette game endpoints
type RouletteHandler struct {
	// engine *service.RouletteEngine
	// walletService ports.WalletService
}

// NewRouletteHandler creates a new handler
func NewRouletteHandler() *RouletteHandler {
	return &RouletteHandler{}
}

// PlaceBet handles single bet placement
// POST /api/v1/roulette/bets
func (h *RouletteHandler) PlaceBet(w http.ResponseWriter, r *http.Request) {
	var req PlaceBetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate request
	if req.Amount <= 0 {
		http.Error(w, "Bet amount must be positive", http.StatusBadRequest)
		return
	}

	// Get user ID from context (set by auth middleware)
	userID := getUserIDFromContext(r.Context())
	if userID == uuid.Nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// TODO: Place bet via engine
	// bet, err := h.engine.PlaceBet(r.Context(), tableID, userID, betType, numbers, amount)

	respondJSON(w, http.StatusCreated, map[string]interface{}{
		"message": "Bet placed successfully",
		"bet_id":  uuid.New().String(),
	})
}

// PlaceBets handles multiple bet placement
// POST /api/v1/roulette/bets/multi
func (h *RouletteHandler) PlaceBets(w http.ResponseWriter, r *http.Request) {
	var req PlaceBetsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if len(req.Bets) == 0 {
		http.Error(w, "No bets provided", http.StatusBadRequest)
		return
	}

	userID := getUserIDFromContext(r.Context())
	if userID == uuid.Nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Calculate total bet amount
	totalAmount := 0.0
	for _, bet := range req.Bets {
		totalAmount += bet.Amount
	}

	// TODO: Deduct from wallet and place bets
	// err := h.walletService.PlaceBet(r.Context(), userID, totalAmount, "roulette")

	respondJSON(w, http.StatusCreated, map[string]interface{}{
		"message":      "Bets placed successfully",
		"bet_count":    len(req.Bets),
		"total_amount": totalAmount,
	})
}

// Spin triggers the roulette spin
// POST /api/v1/roulette/spin
func (h *RouletteHandler) Spin(w http.ResponseWriter, r *http.Request) {
	// In production, only authenticated users or dealers can trigger spins
	userID := getUserIDFromContext(r.Context())
	if userID == uuid.Nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Default table ID (in production, get from request or session)
	tableID := uuid.MustParse("00000000-0000-0000-0000-000000000001")

	// TODO: Execute spin via engine
	// result, err := h.engine.Spin(r.Context(), tableID)
	// if err != nil {
	//     http.Error(w, err.Error(), http.StatusBadRequest)
	//     return
	// }

	// Mock response
	response := SpinResponse{
		SpinID:       uuid.New().String(),
		Number:       17,
		Color:        "black",
		Winnings:     175.00,
		TotalWagered: 50.00,
		TotalPaidOut: 175.00,
		SpinHash:     "abc123def456...",
		BetResults: []BetResult{
			{
				BetType: "straight",
				Amount:  5.00,
				Won:     true,
				Payout:  175.00,
			},
		},
	}

	respondJSON(w, http.StatusOK, response)
}

// GetHistory returns spin history
// GET /api/v1/roulette/history
func (h *RouletteHandler) GetHistory(w http.ResponseWriter, r *http.Request) {
	limit := 10
	if l := r.URL.Query().Get("limit"); l != "" {
		// Parse limit
	}

	tableID := uuid.MustParse("00000000-0000-0000-0000-000000000001")

	// TODO: Get history from engine
	// history, err := h.engine.GetHistory(r.Context(), tableID, limit)

	// Mock history
	history := []map[string]interface{}{
		{
			"spin_id":   uuid.New().String(),
			"number":    17,
			"color":     "black",
			"timestamp": "2026-02-12T10:30:00Z",
		},
		{
			"spin_id":   uuid.New().String(),
			"number":    32,
			"color":     "red",
			"timestamp": "2026-02-12T10:28:00Z",
		},
		{
			"spin_id":   uuid.New().String(),
			"number":    0,
			"color":     "green",
			"timestamp": "2026-02-12T10:26:00Z",
		},
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"history": history,
		"count":   len(history),
	})
}

// GetStatistics returns table statistics
// GET /api/v1/roulette/stats
func (h *RouletteHandler) GetStatistics(w http.ResponseWriter, r *http.Request) {
	tableID := uuid.MustParse("00000000-0000-0000-0000-000000000001")

	// TODO: Calculate stats from database
	_ = tableID

	stats := map[string]interface{}{
		"total_spins": 1523,
		"hot_numbers": []int{17, 32, 24, 7, 11},
		"cold_numbers": []int{15, 28, 3, 22, 35},
		"color_distribution": map[string]int{
			"red":   720,
			"black": 762,
			"green": 41,
		},
		"last_10_numbers": []int{17, 32, 0, 24, 7, 15, 28, 11, 3, 22},
	}

	respondJSON(w, http.StatusOK, stats)
}

// VerifyResult verifies a provably fair result
// GET /api/v1/roulette/verify/:hash
func (h *RouletteHandler) VerifyResult(w http.ResponseWriter, r *http.Request) {
	hash := chi.URLParam(r, "hash")

	if hash == "" {
		http.Error(w, "Hash required", http.StatusBadRequest)
		return
	}

	// TODO: Verify hash against stored server seed
	// valid, result := h.engine.VerifyHash(r.Context(), hash)

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"valid":       true,
		"number":      17,
		"server_seed": "abc123...",
		"hash":        hash,
	})
}

// Helper functions

func getUserIDFromContext(ctx interface{}) uuid.UUID {
	// TODO: Extract from actual context
	// This would be set by auth middleware
	return uuid.New()
}

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// RegisterRoutes registers all roulette routes
func (h *RouletteHandler) RegisterRoutes(r chi.Router) {
	r.Post("/roulette/bets", h.PlaceBet)
	r.Post("/roulette/bets/multi", h.PlaceBets)
	r.Post("/roulette/spin", h.Spin)
	r.Get("/roulette/history", h.GetHistory)
	r.Get("/roulette/stats", h.GetStatistics)
	r.Get("/roulette/verify/{hash}", h.VerifyResult)
}
