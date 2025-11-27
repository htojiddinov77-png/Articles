package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/htojiddinov77-png/Articles/internal/app"
)

func SetupRoutes(app *app.Application) *chi.Mux {
	r := chi.NewRouter()
	// article routes
	r.Get("/articles/{id}", app.ArticleHandler.HandlerGetArticleById)
	r.Get("/health", app.HealthCheck)
	r.Post("/articles", app.ArticleHandler.HandlerCreateArticle)
	r.Put("/articles/{id}", app.ArticleHandler.HandleUpdateArticleById)
	r.Delete("/articles/{id}", app.ArticleHandler.HandleDeleteArticlebyId)

	// user routes
	r.Post("/users/register/", app.UserHandler.HandleRegisterUser)
	r.Get("/users/{id}", app.UserHandler.HandleGetUserById)
	r.Put("/users/{id}", app.UserHandler.HandleUpdateUser)
	r.Delete("/users/{id}", app.UserHandler.HandleDeleteUser)

	
	// // user password change
	r.Post("/users/{id}/password-change/", app.UserHandler.HandleChangePassword) // password change
	r.Post("/users/password-reset-request", app.UserHandler.HandlePasswordResetRequest) // password reset requst
	r.Post("/users/password-reset/{token}", app.UserHandler.HandlePasswordRequst) 	// password reset



	r.Post("/reviews", app.ReviewHandler.HandleCreateReview)
	r.Get("/reviews/{id}", app.ReviewHandler.HandleGetReviewByid)
	r.Put("/reviews/{id}", app.ReviewHandler.HandleUpdateReviewById)
	r.Delete("/reviews/{id}", app.ReviewHandler.HandleDeleteReview)

	r.Post("/tokens/authentication", app.TokenHandler.HandleCreateToken)

	return r
}
