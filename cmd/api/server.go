package server

import (
	"library-management/backend/internal/api"
	"library-management/backend/internal/config"
	"library-management/backend/internal/database"
	"log"

	"github.com/joho/godotenv"
)

func Start() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	cfg := config.NewConfig()
	err = cfg.ParseFlag()
	if err != nil {
		log.Fatal("Failed to parse env variables")
	}
	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatal("failed to connect to database")
		panic(err)
	}
	api := api.NewAPI(cfg, db)

	err = api.Run()
	if err != nil {
		log.Fatal("Failed to start the server")
	}
}
