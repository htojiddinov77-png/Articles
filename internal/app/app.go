package app

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/htojiddinov77-png/Articles/internal/api"
	"github.com/htojiddinov77-png/Articles/internal/store"
	"github.com/htojiddinov77-png/Articles/migrations"
)

type Application struct {
	Logger *log.Logger
	ArticleHandler *api.ArticleHandler
	DB             *sql.DB
}

func NewApplication() (*Application, error) {
	pgDB, err := store.Open()
	if err != nil {
		return nil, err
	}
	err = store.MigrateFS(pgDB, migrations.FS, ".")
	if err != nil {
		panic(err)
	}



	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	articleHandler := api.NewArticleHandler()
	app := &Application{
		Logger: logger,
		ArticleHandler: articleHandler,
		DB:             pgDB,
	}
	return app,nil
}

func (a *Application) HealtheCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w , "Status is available\n")
}
