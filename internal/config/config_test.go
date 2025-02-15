package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseFlag(t *testing.T) {
	sampleEnv := Config{
		Env:    "prod",
		Server: ServerConfig{Port: ":8081"},
		DB:     DbConfig{DSN: "host=localhost user=postgres password=postgres dbname=library port=5433 sslmode=disable"},
	}

	err := os.Setenv("PORT", sampleEnv.Server.Port)
	assert.NoError(t, err)
	err = os.Setenv("ENV", sampleEnv.Env)
	assert.NoError(t, err)
	err = os.Setenv("DSN", sampleEnv.DB.DSN)
	assert.NoError(t, err)

	cfg := *NewConfig()
	err = cfg.ParseFlag()
	assert.NoError(t, err)

	assert.Equal(t, sampleEnv.Server.Port, cfg.Server.Port)
	assert.Equal(t, sampleEnv.Env, cfg.Env)
	assert.Equal(t, sampleEnv.DB.DSN, cfg.DB.DSN)
}
