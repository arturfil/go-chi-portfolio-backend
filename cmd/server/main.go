package main

import (
	"log"
	"os"
	"portfolio-api/application"
	"portfolio-api/db"
	"portfolio-api/models"
)

func main() {
	var cfg application.Config
	port := os.Getenv("PORT")
	cfg.Port = port

	dsn := os.Getenv("DSN")
	dbConn, err := db.ConnectPostgres(dsn)
	if err != nil {
		log.Fatal("Cannot connect to database", err)
	}

	defer dbConn.DB.Close()

	App := &application.Application{
		Config: cfg,
		Models: models.New(dbConn.DB),
	}

	err = App.Serve()
	if err != nil {
		log.Fatal(err)
	}
}
