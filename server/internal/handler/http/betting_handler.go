package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type BettingHandler struct {
	// TODO: добавить betting service когда будет реализован
}

func NewBettingHandler() *BettingHandler {
	return &BettingHandler{}
}

type placeBetRequest struct {
	EventID uuid.UUID `json:"event_id" binding:"required"`
	BetType string    `json:"bet_type" binding:"required,oneof=home draw away"`
	Amount  float64   `json:"amount" binding:"required,gt=0"`
}

// POST /api/v1/betting/bets - Place a bet
func (h *BettingHandler) PlaceBet(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	var req placeBetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "invalid request: check event_id, bet_type (home/draw/away), and amount",
		})
		return
	}

	c.JSON(http.StatusCreated, Response{
		Success: true,
		Data: gin.H{
			"message": "Bet placed successfully (demo)",
			"user_id": userID,
			"bet":     req,
		},
	})
}

// GET /api/v1/betting/bets - Get user's bets
func (h *BettingHandler) GetBets(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	status := c.Query("status")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	if limit <= 0 || limit > 100 {
		limit = 20
	}

	// TODO: Fetch bets from database
	_ = status
	_ = offset

	c.JSON(http.StatusOK, Response{
		Success: true,
		Data: gin.H{
			"message": "Bets retrieved (demo)",
			"user_id": userID,
			"limit":   limit,
			"bets":    []interface{}{},
		},
	})
}

// DELETE /api/v1/betting/bets/:id - Cancel a bet
func (h *BettingHandler) CancelBet(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	betID := c.Param("id")
	betUUID, err := uuid.Parse(betID)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "invalid bet ID",
		})
		return
	}

	// TODO: Implement cancel logic
	// 1. Get bet from database
	// 2. Check bet belongs to user
	// 3. Check bet status is 'pending'
	// 4. Check event hasn't started
	// 5. Refund amount to wallet
	// 6. Update bet status to 'cancelled'

	c.JSON(http.StatusOK, Response{
		Success: true,
		Data: gin.H{
			"message": "Bet cancelled (demo)",
			"user_id": userID,
			"bet_id":  betUUID,
		},
	})
}

// GET /api/v1/betting/events - Get available sport events
func (h *BettingHandler) GetEvents(c *gin.Context) {
	sport := c.Query("sport")
	status := c.DefaultQuery("status", "upcoming")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))

	if limit <= 0 || limit > 100 {
		limit = 50
	}

	// TODO: Fetch events from database

	c.JSON(http.StatusOK, Response{
		Success: true,
		Data: gin.H{
			"message": "Events retrieved (demo)",
			"sport":   sport,
			"status":  status,
			"limit":   limit,
			"events":  []interface{}{},
		},
	})
}
