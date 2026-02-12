package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jokeoa/goigaming/internal/core/ports"
)

type RouletteHandler struct {
	rouletteService ports.RouletteService
}

func NewRouletteHandler(rouletteService ports.RouletteService) *RouletteHandler {
	return &RouletteHandler{rouletteService: rouletteService}
}

func (h *RouletteHandler) ListTables(c *gin.Context) {
	tables, err := h.rouletteService.ListActiveTables(c.Request.Context())
	if err != nil {
		respondError(c, err)
		return
	}

	respondSuccess(c, http.StatusOK, tables)
}

func (h *RouletteHandler) GetTable(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{Success: false, Error: "invalid table id"})
		return
	}

	table, err := h.rouletteService.GetTable(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}

	respondSuccess(c, http.StatusOK, table)
}

func (h *RouletteHandler) GetCurrentRound(c *gin.Context) {
	tableID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{Success: false, Error: "invalid table id"})
		return
	}

	round, err := h.rouletteService.GetCurrentRound(c.Request.Context(), tableID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondSuccess(c, http.StatusOK, round)
}

type placeBetRequest struct {
	RoundID  string `json:"round_id" binding:"required"`
	BetType  string `json:"bet_type" binding:"required"`
	BetValue string `json:"bet_value" binding:"required"`
	Amount   string `json:"amount" binding:"required"`
}

func (h *RouletteHandler) PlaceBet(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	tableID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{Success: false, Error: "invalid table id"})
		return
	}

	var req placeBetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "invalid request body",
		})
		return
	}

	roundID, err := uuid.Parse(req.RoundID)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{Success: false, Error: "invalid round_id"})
		return
	}

	bet, err := h.rouletteService.PlaceBet(
		c.Request.Context(), userID, tableID, roundID,
		req.BetType, req.BetValue, req.Amount,
	)
	if err != nil {
		respondError(c, err)
		return
	}

	respondSuccess(c, http.StatusCreated, bet)
}

func (h *RouletteHandler) GetRoundHistory(c *gin.Context) {
	tableID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{Success: false, Error: "invalid table id"})
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	rounds, err := h.rouletteService.GetRoundHistory(c.Request.Context(), tableID, limit, offset)
	if err != nil {
		respondError(c, err)
		return
	}

	respondSuccess(c, http.StatusOK, rounds)
}

func (h *RouletteHandler) GetMyBets(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	bets, err := h.rouletteService.GetUserBets(c.Request.Context(), userID, limit, offset)
	if err != nil {
		respondError(c, err)
		return
	}

	respondSuccess(c, http.StatusOK, bets)
}

func (h *RouletteHandler) GetRound(c *gin.Context) {
	roundID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{Success: false, Error: "invalid round id"})
		return
	}

	round, err := h.rouletteService.GetRound(c.Request.Context(), roundID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondSuccess(c, http.StatusOK, round)
}
