package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/htojiddinov77-png/Articles/internal/app"
)

func SetupRoutes(app *app.Application) *chi.Mux {
	r := chi.NewRouter()
	// Authenticate runs for every request.
	// If token present -> sets real user.
	// If no token -> sets AnonymousUser.
	r.Use(app.Middleware.Authenticate)

	// PUBLIC ROUTES (no login required) 
	r.Get("/health", app.HealthCheck)

	r.Get("/articles/{id}", app.ArticleHandler.HandlerGetArticleById)
	r.Get("/reviews/{id}", app.ReviewHandler.HandleGetReviewByid)

	r.Post("/users/register/", app.UserHandler.HandleRegisterUser)
	r.Post("/tokens/authentication", app.TokenHandler.HandleCreateToken)
	// // user password change
	r.Post("/users/{id}/password-change/", app.UserHandler.HandleChangePassword)        // password change
	r.Post("/users/password-reset-request", app.UserHandler.HandlePasswordResetRequest) // password reset requst
	r.Post("/users/password-reset/{token}", app.UserHandler.HandlePasswordReset)        // password reset


	// PROTECTED ROUTES (login required)
	r.Group(func(r chi.Router){

		r.Use(app.Middleware.RequireUser)
		r.Post("/articles", app.ArticleHandler.HandlerCreateArticle)
		r.Put("/articles/{id}", app.ArticleHandler.HandleUpdateArticleById)
		r.Delete("/articles/{id}", app.ArticleHandler.HandleDeleteArticlebyId)

		
		r.Get("/users/{id}", app.UserHandler.HandleGetUserById)
		r.Put("/users/{id}", app.UserHandler.HandleUpdateUser)
		r.Delete("/users/{id}", app.UserHandler.HandleDeleteUser)

		r.Post("/reviews", app.ReviewHandler.HandleCreateReview)
		r.Put("/reviews/{id}", app.ReviewHandler.HandleUpdateReviewById)
		r.Delete("/reviews/{id}", app.ReviewHandler.HandleDeleteReview)
	})

	return r
}
