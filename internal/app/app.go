package app

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/htojiddinov77-png/Articles/internal/api"
	"github.com/htojiddinov77-png/Articles/internal/migrations"
	"github.com/htojiddinov77-png/Articles/internal/store"
)

type Application struct {
	Logger         *log.Logger
	ArticleHandler *api.ArticleHandler
	UserHandler    *api.UserHandler
	ReviewHandler  *api.ReviewHandler
	TokenHandler   *api.TokenHandler
	DB             *sql.DB
}

func NewApplication() (*Application, error) {
	pgDB, err := store.Open()
	if err != nil {
		return nil, err
	}

	err = store.MigrateFs(pgDB, migrations.FS, ".")
	if err != nil {
		panic(err)
	}

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	articleStore := store.NewPostgresArticleStore(pgDB)
	tokenStore := store.NewPostgresTokenStore(pgDB)
	userStore := store.NewPostgresUserStore(pgDB)
	reviewStore := store.NewPostgresReviewStore(pgDB)

	articleHandler := api.NewArticleHandler(articleStore, logger)
	userHandler := api.NewUserHandler(userStore, logger)
	reviewHandler := api.NewReviewHandler(reviewStore, logger)
	tokenHandler := api.NewTokenHandler(tokenStore, userStore, logger)

	app := &Application{
		Logger:         logger,
		ArticleHandler: articleHandler,
		UserHandler:    userHandler,
		ReviewHandler:  reviewHandler,
		TokenHandler:   tokenHandler,
		DB:             pgDB,
	}
	return app, nil
}

func (a *Application) HealthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Status is available\n")
}
