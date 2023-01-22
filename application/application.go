package application

import (
	"fmt"
	"net/http"
	"os"
	"portfolio-api/helpers"
	"portfolio-api/models"
	"portfolio-api/routes"
)

type Config struct {
	Port string
}

type Application struct {
	Config Config
	Models models.Models
}

func (app *Application) Serve() error {
	port := os.Getenv("PORT")
	helpers.MessageLogs.Infolog.Println("API listening on port", port)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: routes.Routes(),
	}
	return srv.ListenAndServe()
}
