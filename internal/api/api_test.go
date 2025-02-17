package api

import (
	"encoding/json"
	"library-management/backend/internal/config"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

type Response struct {
	Message string `json:"message"`
}

var cfg = &config.Config{
	Env:    "prod",
	Server: config.ServerConfig{Port: "8081"},
	DB:     config.DbConfig{DSN: "host=localhost user=postgres password=postgres dbname=library port=5433 sslmode=disable"},
}
var db *gorm.DB

// var db, _ = database.Connect(cfg.DB.DSN)

// func LoadTestCreds() {
// 	sampleEnv = &config.Config{
// 		Env:    "prod",
// 		Server: config.ServerConfig{Port: "8081"},
// 		DB:     config.DbConfig{DSN: "host=localhost user=postgres password=postgres dbname=library port=5433 sslmode=disable"},
// 	}
// 	db, _ = database.Connect(sampleEnv)
// }

func TestNewApi(t *testing.T) {
	api := NewAPI(cfg, db)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	api.Router.ServeHTTP(w, req)

	response := Response{
		Message: "pong",
	}
	responseJson, _ := json.Marshal(response)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, string(responseJson), w.Body.String())
}
