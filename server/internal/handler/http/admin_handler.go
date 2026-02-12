package http

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jokeoa/goigaming/internal/core/domain"
	"github.com/jokeoa/goigaming/internal/core/ports"
	"github.com/jokeoa/goigaming/models"
	"github.com/jokeoa/goigaming/repository"
	"github.com/shopspring/decimal"
)

type AdminHandler struct {
	pokerRepo    ports.PokerTableRepository
	rouletteRepo *repository.RouletteTableRepository
}

func NewAdminHandler(
	pokerRepo ports.PokerTableRepository,
	rouletteRepo *repository.RouletteTableRepository,
) *AdminHandler {
	return &AdminHandler{
		pokerRepo:    pokerRepo,
		rouletteRepo: rouletteRepo,
	}
}

// --- Poker Tables ---

type createPokerTableRequest struct {
	Name       string `json:"name" binding:"required"`
	SmallBlind string `json:"small_blind" binding:"required"`
	BigBlind   string `json:"big_blind" binding:"required"`
	MinBuyIn   string `json:"min_buy_in" binding:"required"`
	MaxBuyIn   string `json:"max_buy_in" binding:"required"`
	MaxPlayers int    `json:"max_players" binding:"required,min=2,max=10"`
}

type updatePokerTableRequest struct {
	Name       string `json:"name" binding:"required"`
	SmallBlind string `json:"small_blind" binding:"required"`
	BigBlind   string `json:"big_blind" binding:"required"`
	MinBuyIn   string `json:"min_buy_in" binding:"required"`
	MaxBuyIn   string `json:"max_buy_in" binding:"required"`
	MaxPlayers int    `json:"max_players" binding:"required,min=2,max=10"`
	Status     string `json:"status" binding:"required,oneof=waiting active closed"`
}

type pokerDecimals struct {
	smallBlind decimal.Decimal
	bigBlind   decimal.Decimal
	minBuyIn   decimal.Decimal
	maxBuyIn   decimal.Decimal
}

func parsePokerDecimals(sb, bb, minBI, maxBI string) (pokerDecimals, error) {
	smallBlind, err := decimal.NewFromString(sb)
	if err != nil {
		return pokerDecimals{}, fmt.Errorf("invalid small_blind")
	}
	bigBlind, err := decimal.NewFromString(bb)
	if err != nil {
		return pokerDecimals{}, fmt.Errorf("invalid big_blind")
	}
	minBuyIn, err := decimal.NewFromString(minBI)
	if err != nil {
		return pokerDecimals{}, fmt.Errorf("invalid min_buy_in")
	}
	maxBuyIn, err := decimal.NewFromString(maxBI)
	if err != nil {
		return pokerDecimals{}, fmt.Errorf("invalid max_buy_in")
	}

	if smallBlind.LessThanOrEqual(decimal.Zero) || bigBlind.LessThanOrEqual(decimal.Zero) {
		return pokerDecimals{}, fmt.Errorf("blinds must be positive")
	}
	if smallBlind.GreaterThanOrEqual(bigBlind) {
		return pokerDecimals{}, fmt.Errorf("small_blind must be less than big_blind")
	}
	if minBuyIn.LessThanOrEqual(decimal.Zero) || maxBuyIn.LessThanOrEqual(decimal.Zero) {
		return pokerDecimals{}, fmt.Errorf("buy-in amounts must be positive")
	}
	if minBuyIn.GreaterThanOrEqual(maxBuyIn) {
		return pokerDecimals{}, fmt.Errorf("min_buy_in must be less than max_buy_in")
	}

	return pokerDecimals{
		smallBlind: smallBlind,
		bigBlind:   bigBlind,
		minBuyIn:   minBuyIn,
		maxBuyIn:   maxBuyIn,
	}, nil
}

func (h *AdminHandler) CreatePokerTable(c *gin.Context) {
	var req createPokerTableRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "invalid request: " + err.Error(),
		})
		return
	}

	d, err := parsePokerDecimals(req.SmallBlind, req.BigBlind, req.MinBuyIn, req.MaxBuyIn)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{Success: false, Error: err.Error()})
		return
	}

	table, err := h.pokerRepo.Create(c.Request.Context(), domain.PokerTable{
		Name:       req.Name,
		SmallBlind: d.smallBlind,
		BigBlind:   d.bigBlind,
		MinBuyIn:   d.minBuyIn,
		MaxBuyIn:   d.maxBuyIn,
		MaxPlayers: req.MaxPlayers,
		Status:     domain.TableStatusWaiting,
	})
	if err != nil {
		respondError(c, err)
		return
	}

	respondSuccess(c, http.StatusCreated, table)
}

func (h *AdminHandler) ListPokerTables(c *gin.Context) {
	tables, err := h.pokerRepo.FindAll(c.Request.Context())
	if err != nil {
		respondError(c, err)
		return
	}

	respondSuccess(c, http.StatusOK, tables)
}

func (h *AdminHandler) GetPokerTable(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{Success: false, Error: "invalid table id"})
		return
	}

	table, err := h.pokerRepo.FindByID(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}

	respondSuccess(c, http.StatusOK, table)
}

func (h *AdminHandler) UpdatePokerTable(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{Success: false, Error: "invalid table id"})
		return
	}

	var req updatePokerTableRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "invalid request: " + err.Error(),
		})
		return
	}

	d, err := parsePokerDecimals(req.SmallBlind, req.BigBlind, req.MinBuyIn, req.MaxBuyIn)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{Success: false, Error: err.Error()})
		return
	}

	table, err := h.pokerRepo.Update(c.Request.Context(), domain.PokerTable{
		ID:         id,
		Name:       req.Name,
		SmallBlind: d.smallBlind,
		BigBlind:   d.bigBlind,
		MinBuyIn:   d.minBuyIn,
		MaxBuyIn:   d.maxBuyIn,
		MaxPlayers: req.MaxPlayers,
		Status:     domain.TableStatus(req.Status),
	})
	if err != nil {
		respondError(c, err)
		return
	}

	respondSuccess(c, http.StatusOK, table)
}

func (h *AdminHandler) DeletePokerTable(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{Success: false, Error: "invalid table id"})
		return
	}

	if err := h.pokerRepo.UpdateStatus(c.Request.Context(), id, domain.TableStatusClosed); err != nil {
		respondError(c, err)
		return
	}

	respondSuccess(c, http.StatusOK, gin.H{"message": "poker table closed"})
}

// --- Roulette Tables ---

type createRouletteTableRequest struct {
	Name   string  `json:"name" binding:"required"`
	MinBet float64 `json:"min_bet" binding:"required,gt=0"`
	MaxBet float64 `json:"max_bet" binding:"required,gtfield=MinBet"`
}

type updateRouletteTableRequest struct {
	Name   string  `json:"name" binding:"required"`
	MinBet float64 `json:"min_bet" binding:"required,gt=0"`
	MaxBet float64 `json:"max_bet" binding:"required,gtfield=MinBet"`
	Status string  `json:"status" binding:"required,oneof=active inactive maintenance"`
}

func (h *AdminHandler) CreateRouletteTable(c *gin.Context) {
	var req createRouletteTableRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "invalid request: " + err.Error(),
		})
		return
	}

	table, err := h.rouletteRepo.Create(c.Request.Context(), models.RouletteTable{
		Name:   req.Name,
		MinBet: req.MinBet,
		MaxBet: req.MaxBet,
		Status: "active",
	})
	if err != nil {
		respondError(c, err)
		return
	}

	respondSuccess(c, http.StatusCreated, table)
}

func (h *AdminHandler) ListRouletteTables(c *gin.Context) {
	tables, err := h.rouletteRepo.FindAll(c.Request.Context())
	if err != nil {
		respondError(c, err)
		return
	}

	respondSuccess(c, http.StatusOK, tables)
}

func (h *AdminHandler) GetRouletteTable(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{Success: false, Error: "invalid table id"})
		return
	}

	table, err := h.rouletteRepo.FindByID(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}

	respondSuccess(c, http.StatusOK, table)
}

func (h *AdminHandler) UpdateRouletteTable(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{Success: false, Error: "invalid table id"})
		return
	}

	var req updateRouletteTableRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "invalid request: " + err.Error(),
		})
		return
	}

	table, err := h.rouletteRepo.Update(c.Request.Context(), models.RouletteTable{
		ID:     id,
		Name:   req.Name,
		MinBet: req.MinBet,
		MaxBet: req.MaxBet,
		Status: req.Status,
	})
	if err != nil {
		respondError(c, err)
		return
	}

	respondSuccess(c, http.StatusOK, table)
}

func (h *AdminHandler) DeleteRouletteTable(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{Success: false, Error: "invalid table id"})
		return
	}

	if err := h.rouletteRepo.UpdateStatus(c.Request.Context(), id, "inactive"); err != nil {
		respondError(c, err)
		return
	}

	respondSuccess(c, http.StatusOK, gin.H{"message": "roulette table deactivated"})
}
