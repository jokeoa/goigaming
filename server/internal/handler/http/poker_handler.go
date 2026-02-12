package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jokeoa/goigaming/internal/core/ports"
)

type PokerHandler struct {
	pokerService ports.PokerService
}

func NewPokerHandler(pokerService ports.PokerService) *PokerHandler {
	return &PokerHandler{pokerService: pokerService}
}

func (h *PokerHandler) ListTables(c *gin.Context) {
	tables, err := h.pokerService.ListTables(c.Request.Context())
	if err != nil {
		respondError(c, err)
		return
	}

	respondSuccess(c, http.StatusOK, tables)
}

func (h *PokerHandler) GetTable(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{Success: false, Error: "invalid table id"})
		return
	}

	table, err := h.pokerService.GetTable(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}

	respondSuccess(c, http.StatusOK, table)
}

type joinTableRequest struct {
	SeatNumber int    `json:"seat_number" binding:"required,min=1"`
	BuyIn      string `json:"buy_in" binding:"required"`
}

func (h *PokerHandler) JoinTable(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	tableID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{Success: false, Error: "invalid table id"})
		return
	}

	var req joinTableRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "invalid request body",
		})
		return
	}

	player, err := h.pokerService.JoinTable(c.Request.Context(), tableID, userID, req.SeatNumber, req.BuyIn)
	if err != nil {
		respondError(c, err)
		return
	}

	respondSuccess(c, http.StatusCreated, player)
}

func (h *PokerHandler) LeaveTable(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	tableID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{Success: false, Error: "invalid table id"})
		return
	}

	if err := h.pokerService.LeaveTable(c.Request.Context(), tableID, userID); err != nil {
		respondError(c, err)
		return
	}

	respondSuccess(c, http.StatusOK, gin.H{"message": "left table"})
}

func (h *PokerHandler) GetTableState(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{Success: false, Error: "invalid table id"})
		return
	}

	state, err := h.pokerService.GetTableState(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}

	respondSuccess(c, http.StatusOK, state)
}
