package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/htojiddinov77-png/Articles/internal/app"
)


func SetupRoutes(app *app.Application) *chi.Mux{
	r := chi.NewRouter()

	r.Get("/articles/{id}", app.ArticleHandler.HandlerGetArticleById)
	r.Get("/health", app.HealthCheck)

	r.Post("/articles", app.ArticleHandler.HandlerCreateArticle)
	r.Put("/articles/{id}", app.ArticleHandler.HandleUpdateArticleById)
	r.Delete("/articles/{id}", app.ArticleHandler.HandleDeleteWorkoutbyId)
	
	return r
}