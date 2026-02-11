package handlers

import (
	"html/template"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jokeoa/igaming/services"
)

type PageHandler struct {
	rouletteService *services.RouletteService
	adminService    *services.AdminService
	templates       *template.Template
}

func NewPageHandler(rouletteService *services.RouletteService, adminService *services.AdminService) *PageHandler {
	templates := template.Must(template.ParseGlob("templates/*.html"))
	
	return &PageHandler{
		rouletteService: rouletteService,
		adminService:    adminService,
		templates:       templates,
	}
}

func (h *PageHandler) RoulettePage(w http.ResponseWriter, r *http.Request) {
	tableID := uuid.MustParse("00000000-0000-0000-0000-000000000001")

	currentRound, err := h.rouletteService.GetHistory(r.Context(), tableID, 1)
	var round interface{}
	if err == nil && len(currentRound) > 0 {
		round = currentRound[0]
	}

	history, _ := h.rouletteService.GetHistory(r.Context(), tableID, 20)

	stats := struct {
		TotalBets    int
		TotalWagered float64
		NetProfit    float64
	}{
		TotalBets:    0,
		TotalWagered: 0.0,
		NetProfit:    0.0,
	}

	data := map[string]interface{}{
		"TableID":      tableID.String(),
		"CurrentRound": round,
		"LastSpin":     round,
		"History":      history,
		"Stats":        stats,
	}

	if err := h.templates.ExecuteTemplate(w, "roulette.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *PageHandler) AdminPage(w http.ResponseWriter, r *http.Request) {
	analytics, err := h.adminService.GetAnalytics(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sessions, err := h.adminService.GetActiveSessions(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	topWinners, err := h.adminService.GetTopWinners(r.Context(), 10)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Analytics":  analytics,
		"Sessions":   sessions,
		"TopWinners": topWinners,
	}

	if err := h.templates.ExecuteTemplate(w, "admin.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *PageHandler) RegisterRoutes(r chi.Router) {
	r.Get("/roulette", h.RoulettePage)
	r.Get("/admin", h.AdminPage)
}
