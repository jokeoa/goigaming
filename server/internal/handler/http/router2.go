package http
func setupRoutes(r chi.Router, handlers *Handlers) {
    // ... API routes ...
    
    // Page routes
    r.Get("/roulette", handlers.Pages.RoulettePage)
    r.Get("/lobby", handlers.Pages.LobbyPage)
    r.Get("/", handlers.Pages.HomePage)
}
