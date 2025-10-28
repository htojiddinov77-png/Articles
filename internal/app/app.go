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

	articleStore := store.NewPosgresArticleStore(pgDB)

	articleHandler := api.NewArticleHandler(articleStore, logger)

	app := &Application{
		Logger:         logger,
		ArticleHandler: articleHandler,
		DB:             pgDB,
	}
	return app, nil
}

func (a *Application) HealthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Status is available\n")
}
