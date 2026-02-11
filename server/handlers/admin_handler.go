package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jokeoa/igaming/services"
)

type AdminHandler struct {
	adminService *services.AdminService
}

func NewAdminHandler(adminService *services.AdminService) *AdminHandler {
	return &AdminHandler{
		adminService: adminService,
	}
}

func (h *AdminHandler) GetAnalytics(w http.ResponseWriter, r *http.Request) {
	analytics, err := h.adminService.GetAnalytics(r.Context())
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	respondJSON(w, http.StatusOK, analytics)
}

func (h *AdminHandler) GetActiveSessions(w http.ResponseWriter, r *http.Request) {
	sessions, err := h.adminService.GetActiveSessions(r.Context())
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"sessions": sessions,
		"count":    len(sessions),
	})
}

func (h *AdminHandler) ForceCloseSession(w http.ResponseWriter, r *http.Request) {
	sessionIDStr := chi.URLParam(r, "id")
	sessionID, err := uuid.Parse(sessionIDStr)
	if err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid session_id"})
		return
	}

	err = h.adminService.ForceCloseSession(r.Context(), sessionID)
	if err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{
		"message": "session closed successfully",
		"id":      sessionID.String(),
	})
}

func (h *AdminHandler) GetUserStats(w http.ResponseWriter, r *http.Request) {
	userIDStr := chi.URLParam(r, "id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid user_id"})
		return
	}

	stats, err := h.adminService.GetUserStats(r.Context(), userID)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	respondJSON(w, http.StatusOK, stats)
}

func (h *AdminHandler) GetTopWinners(w http.ResponseWriter, r *http.Request) {
	winners, err := h.adminService.GetTopWinners(r.Context(), 10)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"winners": winners,
		"count":   len(winners),
	})
}

func (h *AdminHandler) RegisterRoutes(r chi.Router) {
	r.Route("/api/v1/admin", func(r chi.Router) {
		r.Get("/analytics", h.GetAnalytics)
		r.Get("/sessions", h.GetActiveSessions)
		r.Delete("/sessions/{id}", h.ForceCloseSession)
		r.Get("/users/{id}/stats", h.GetUserStats)
		r.Get("/top-winners", h.GetTopWinners)
	})
}
