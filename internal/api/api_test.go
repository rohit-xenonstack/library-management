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

func TestNewApi(t *testing.T) {
	var cfg = &config.SampleEnv
	var db *gorm.DB

	api, err := NewAPI(cfg, db)
	assert.NoError(t, err)

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
