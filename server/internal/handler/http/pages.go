package http

import (
    "html/template"
    "net/http"
    
    "github.com/google/uuid"
    "github.com/jokeoa/goigaming/internal/service/roulette"
    tmpl "github.com/jokeoa/goigaming/internal/templates"
)

type PageHandler struct {
    templates       *template.Template
    rouletteService *roulette.Service
}

func NewPageHandler(rouletteService *roulette.Service) *PageHandler {
    // Парсим все templates
    templates := template.Must(template.New("").Funcs(tmpl.GetFuncMap()).ParseGlob("internal/templates/**/*.html"))
    
    return &PageHandler{
        templates:       templates,
        rouletteService: rouletteService,
    }
}

// RoulettePage отображает страницу рулетки
func (h *PageHandler) RoulettePage(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    userID := GetUserID(ctx)
    
    // Получаем данные для страницы
    tableID := uuid.MustParse("default-table-id") // или из параметров
    
    currentRound, _ := h.rouletteService.GetCurrentRound(ctx, tableID)
    history, _ := h.rouletteService.GetHistory(ctx, tableID, 10)
    
    var yourBets []*domain.RouletteBet
    if currentRound != nil {
        yourBets, _ = h.rouletteService.GetUserBets(ctx, userID, currentRound.ID)
    }
    
    data := map[string]interface{}{
        "CurrentRound": currentRound,
        "History":      history,
        "YourBets":     yourBets,
        "User":         GetUser(ctx),
    }
    
    if err := h.templates.ExecuteTemplate(w, "roulette.html", data); err != nil {
        http.Error(w, "Failed to render template", http.StatusInternalServerError)
        return
    }
}
