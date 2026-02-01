package handler

import (
    "encoding/json"
    "net/http"

    "github.com/google/uuid"
    "goigaming/models"
    "goigaming/repository"
)

type RouletteHandler struct {
    betRepo   *repository.RouletteBetRepository
    tableRepo *repository.RouletteTableRepository
    roundRepo *repository.RouletteRoundRepository
}

func NewRouletteHandler(
    betRepo *repository.RouletteBetRepository,
    tableRepo *repository.RouletteTableRepository,
    roundRepo *repository.RouletteRoundRepository,
) *RouletteHandler {
    return &RouletteHandler{
        betRepo:   betRepo,
        tableRepo: tableRepo,
        roundRepo: roundRepo,
    }
}

func (h *RouletteHandler) GetTables(w http.ResponseWriter, r *http.Request) {
    tables, err := h.tableRepo.FindActive(r.Context())
    if err != nil {
        http.Error(w, `{"error": "failed to fetch tables"}`, http.StatusInternalServerError)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(tables)
}

func (h *RouletteHandler) PlaceBet(w http.ResponseWriter, r *http.Request) {
    var req struct {
        RoundID  uuid.UUID `json:"round_id"`
        UserID   uuid.UUID `json:"user_id"`
        BetType  string    `json:"bet_type"`
        BetValue string    `json:"bet_value"`
        Amount   float64   `json:"amount"`
    }

    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, `{"error": "invalid request body"}`, http.StatusBadRequest)
        return
    }

    _, err := h.roundRepo.FindById(r.Context(), req.RoundID)
    if err != nil {
        http.Error(w, `{"error": "round not found"}`, http.StatusNotFound)
        return
    }

    bet, err := h.betRepo.Create(r.Context(), models.RouletteBet{
        RoundID:  req.RoundID,
        UserID:   req.UserID,
        BetType:  req.BetType,
        BetValue: req.BetValue,
        Amount:   req.Amount,
        Status:   "pending",
    })
    if err != nil {
        http.Error(w, `{"error": "failed to place bet"}`, http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(bet)
}

func (h *RouletteHandler) GetHistory(w http.ResponseWriter, r *http.Request) {
    path := r.URL.Path
    userIdStr := path[len("/api/roulette/history/"):]

    userID, err := uuid.Parse(userIdStr)
    if err != nil {
        http.Error(w, `{"error": "invalid user id"}`, http.StatusBadRequest)
        return
    }

    bets, err := h.betRepo.FindByUserId(r.Context(), userID)
    if err != nil {
        http.Error(w, `{"error": "failed to fetch history"}`, http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(bets)
}
