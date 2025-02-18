package server

import (
	"library-management/backend/internal/api"
	"library-management/backend/internal/config"
	"library-management/backend/internal/database"
	"library-management/backend/internal/model"
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

	err = db.AutoMigrate(&model.Library{}, &model.Users{}, &model.IssueRegistry{}, &model.BookInventory{}, &model.RequestEvents{})
	if err != nil {
		log.Fatal("Failed to migrate DB")
	}

	api, err := api.NewAPI(cfg, db)
	if err != nil {
		log.Fatal("Error creating api service")
	}

	err = api.Run()
	if err != nil {
		log.Fatal("Failed to start the server")
	}
}
