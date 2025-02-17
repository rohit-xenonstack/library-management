package database

import (
	"library-management/backend/internal/config"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConnect_WrongCredentials(t *testing.T) {
	db := config.DbConfig{
		DSN: "host=localhost user=wrong_user password=wrong_pass dbname=library port=5433 sslmode=disable",
	}
	cfg := &config.Config{
		DB: db,
	}
	_, err := Connect(cfg)
	assert.Error(t, err)
}

// func TestConnect_Timeout(t *testing.T) {
// 	dsn := "host=10.0.255.255 user=postgres password=postgres dbname=library port=5433 sslmode=disable"
// 	_, err := Connect(dsn)
// 	assert.Error(t, err)
// }

func TestConnect(t *testing.T) {
	dbcfg := config.DbConfig{
		DSN: "host=localhost user=postgres password=postgres dbname=library port=5433 sslmode=disable",
	}
	cfg := &config.Config{
		DB: dbcfg,
	}
	db, err := Connect(cfg)
	assert.NoError(t, err)

	sqlDB, err := db.DB()
	err = sqlDB.Ping()
	assert.NoError(t, err)
}
