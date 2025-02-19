package database

import (
	"library-management/backend/internal/config"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	// Setup test database
	cfg := &config.Config{
		DB: config.DbConfig{
			DSN: "host=localhost user=postgres password=postgres dbname=library port=5433 sslmode=disable TimeZone=UTC",
		},
	}

	db, err := Connect(cfg)
	if err != nil {
		panic("failed to connect database: " + err.Error())
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic("failed to get database instance: " + err.Error())
	}

	// Run tests
	code := m.Run()

	// Cleanup
	sqlDB.Close()

	os.Exit(code)
}

func TestConnect(t *testing.T) {
	tests := []struct {
		name    string
		cfg     *config.Config
		wantErr bool
	}{
		{
			name: "successful connection",
			cfg: &config.Config{
				DB: config.DbConfig{
					DSN: "host=localhost user=postgres password=postgres dbname=library port=5433 sslmode=disable TimeZone=UTC",
				},
			},
			wantErr: false,
		},
		{
			name: "invalid connection string",
			cfg: &config.Config{
				DB: config.DbConfig{
					DSN: "invalid-connection-string",
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, err := Connect(tt.cfg)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, db)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, db)

			// Test connection pool settings
			sqlDB, err := db.DB()
			require.NoError(t, err)

			// Test actual connection
			err = sqlDB.Ping()
			assert.NoError(t, err, "Database should be reachable")

			// Check GORM configuration
			assert.True(t, db.Config.SkipDefaultTransaction)
			assert.True(t, db.Config.PrepareStmt)
		})
	}
}

func TestConnect_NilConfig(t *testing.T) {
	db, err := Connect(nil)
	assert.Error(t, err)
	assert.Nil(t, db)
}
