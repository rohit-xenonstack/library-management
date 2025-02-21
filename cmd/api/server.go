package server

import (
	"library-management/backend/internal/api"
	"library-management/backend/internal/api/model"
	"library-management/backend/internal/config"
	"library-management/backend/internal/database"
	"log"

	"github.com/joho/godotenv"
)

func Start() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
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

	err = db.AutoMigrate(&model.Library{}, &model.Users{}, &model.BookInventory{}, &model.RequestEvents{}, &model.IssueRegistry{})
	if err != nil {
		log.Fatal("failed to migrate DB")
	}
	// err = db.AutoMigrate()
	// if err != nil {
	// 	log.Fatal("failed to migrate DB")
	// }

	h := cfg.InitHandler(cfg.InitRepository(db))
	api := api.NewAPI(cfg, h)
	if err != nil {
		log.Fatal("cannot create api server")
	}

	err = api.Run()
	if err != nil {
		log.Fatal("failed to start the server")
	}
}
