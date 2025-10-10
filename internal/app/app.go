package app

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/htojiddinov77-png/GolangArticles/internal/api"
)

type Application struct {
	Logger *log.Logger
	ArticleHandler *api.ArticleHandler
}

func NewApplication() (*Application, error) {
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	articleHandler := api.NewArticleHandler()
	app := &Application{
		Logger: logger,
		ArticleHandler: articleHandler,
	}
	return app,nil
}

func (a *Application) HealtheCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w , "Status is available\n")
}
