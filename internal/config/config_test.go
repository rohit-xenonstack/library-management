package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseFlag(t *testing.T) {
	sampleEnv := SampleEnv

	err := os.Setenv("PORT", sampleEnv.Server.Port)
	assert.NoError(t, err)
	err = os.Setenv("ENVIRONMENT", sampleEnv.Env)
	assert.NoError(t, err)
	err = os.Setenv("DATA_SOURCE_NAME", sampleEnv.DB.DSN)
	assert.NoError(t, err)
	err = os.Setenv("JWT_SECRET_KEY", sampleEnv.JWT.SecretKey)
	assert.NoError(t, err)
	err = os.Setenv("ACCESS_TOKEN_DURATION", sampleEnv.JWT.AccessTokenDuration.String())
	assert.NoError(t, err)

	cfg := *NewConfig()
	err = cfg.ParseFlag()
	assert.NoError(t, err)

	assert.Equal(t, sampleEnv.Server.Port, cfg.Server.Port)
	assert.Equal(t, sampleEnv.Env, cfg.Env)
	assert.Equal(t, sampleEnv.DB.DSN, cfg.DB.DSN)
	assert.Equal(t, sampleEnv.JWT.SecretKey, cfg.JWT.SecretKey)
	assert.Equal(t, sampleEnv.JWT.AccessTokenDuration, cfg.JWT.AccessTokenDuration)
}
