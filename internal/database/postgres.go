package database

import (
	"fmt"
	"library-management/backend/internal/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(cfg *config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(cfg.DB.DSN), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	fmt.Println("Database connected successfuly")
	return db, nil
}
