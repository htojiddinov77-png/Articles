package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/htojiddinov77-png/GolangArticles/internal/app"
)

func SetupRoutes(app *app.Application) *chi.Mux {
	r := chi.NewRouter()

	r.Get("/health", app.HealtheCheck)
	r.Get("/articles/{id}", app.ArticleHandler.HandleGetArticleByID)

	r.Post("/articles", app.ArticleHandler.HandleCreateArticle)
	return r
}