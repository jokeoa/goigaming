package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jokeoa/goigaming/internal/core/domain"
	"github.com/jokeoa/goigaming/internal/core/ports"
)

type AuthHandler struct {
	authService ports.AuthService
}

func NewAuthHandler(authService ports.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

type registerRequest struct {
	Username string `json:"username" binding:"required,min=3,max=100"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,max=72"`
}

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "invalid request: check required fields and format",
		})
		return
	}

	user, err := h.authService.Register(c.Request.Context(), req.Username, req.Email, req.Password)
	if err != nil {
		respondError(c, err)
		return
	}

	respondSuccess(c, http.StatusCreated, domain.NewUserProfile(user))
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "invalid request: check required fields and format",
		})
		return
	}

	tokenPair, err := h.authService.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		respondError(c, err)
		return
	}

	respondSuccess(c, http.StatusOK, tokenPair)
}
