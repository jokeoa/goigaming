package http

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jokeoa/igaming/internal/core/domain"
	"github.com/jokeoa/igaming/internal/core/ports"
	"github.com/shopspring/decimal"
)

type WalletHandler struct {
	walletService ports.WalletService
}

func NewWalletHandler(walletService ports.WalletService) *WalletHandler {
	return &WalletHandler{walletService: walletService}
}

type amountRequest struct {
	Amount string `json:"amount" binding:"required"`
}

type walletResponse struct {
	UserID    uuid.UUID       `json:"user_id"`
	Balance   decimal.Decimal `json:"balance"`
	UpdatedAt time.Time       `json:"updated_at"`
}

func newWalletResponse(w domain.Wallet) walletResponse {
	return walletResponse{
		UserID:    w.UserID,
		Balance:   w.Balance,
		UpdatedAt: w.UpdatedAt,
	}
}

func (h *WalletHandler) GetBalance(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	wallet, err := h.walletService.GetBalance(c.Request.Context(), userID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondSuccess(c, http.StatusOK, newWalletResponse(wallet))
}

func (h *WalletHandler) Deposit(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	var req amountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "invalid request: amount is required",
		})
		return
	}

	wallet, err := h.walletService.Deposit(c.Request.Context(), userID, req.Amount)
	if err != nil {
		respondError(c, err)
		return
	}

	respondSuccess(c, http.StatusOK, newWalletResponse(wallet))
}

func (h *WalletHandler) Withdraw(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	var req amountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "invalid request: amount is required",
		})
		return
	}

	wallet, err := h.walletService.Withdraw(c.Request.Context(), userID, req.Amount)
	if err != nil {
		respondError(c, err)
		return
	}

	respondSuccess(c, http.StatusOK, newWalletResponse(wallet))
}

func (h *WalletHandler) GetTransactions(c *gin.Context) {
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

	txs, err := h.walletService.GetTransactions(c.Request.Context(), userID, limit, offset)
	if err != nil {
		respondError(c, err)
		return
	}

	respondSuccess(c, http.StatusOK, txs)
}
