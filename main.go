package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/htojiddinov77-png/GolangArticles/internal/app"
	"github.com/htojiddinov77-png/GolangArticles/internal/routes"
)

func main() {
	var port int
	flag.IntVar(&port, "port", 8080, "go backend server port")
	flag.Parse()

	app, err := app.NewApplication()
	if err != nil{
		panic(err)
	}
	r := routes.SetupRoutes(app)

	server := &http.Server{
		Addr: fmt.Sprintf(":%d", port),
		Handler: r,
		IdleTimeout: time.Minute,
		ReadTimeout: 10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	app.Logger.Printf("we are running on port %d\n", port)

	err = server.ListenAndServe()
	if err != nil{
		app.Logger.Fatal(err)
	}

}