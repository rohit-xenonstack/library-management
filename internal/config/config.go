package config

import (
	"flag"
	"library-management/backend/internal/api/handler"
	"library-management/backend/internal/database/repository"
	"os"
	"time"

	"gorm.io/gorm"
)

type Config struct {
	Env    string
	Server ServerConfig
	DB     DbConfig
	JWT    JWTConfig
}
type ServerConfig struct {
	Port string
}
type DbConfig struct {
	DSN string
}
type JWTConfig struct {
	SecretKey            string
	AccessTokenDuration  time.Duration
	RefreshTokenDuration time.Duration
}

func NewConfig() *Config {
	return &Config{}
}

func (cfg *Config) ParseFlag() error {
	flag.StringVar(&cfg.Server.Port, "port", os.Getenv("PORT"), "API server port")
	flag.StringVar(&cfg.Env, "env", os.Getenv("ENVIRONMENT"), "Environment(dev|prod)")
	flag.StringVar(&cfg.DB.DSN, "db-dsn", os.Getenv("DATA_SOURCE_NAME"), "PostgreSQL DSN")
	flag.StringVar(&cfg.JWT.SecretKey, "jwt-secret", os.Getenv("JWT_SECRET_KEY"), "JWT Secret Key")
	duration, err := time.ParseDuration(os.Getenv("ACCESS_TOKEN_DURATION"))
	if err != nil {
		return err
	}
	flag.DurationVar(&cfg.JWT.AccessTokenDuration, "jwt-access-token-duration", duration, "Access Token Lifetime")
	duration, err = time.ParseDuration(os.Getenv("REFRESH_TOKEN_DURATION"))
	if err != nil {
		return err
	}
	flag.DurationVar(&cfg.JWT.RefreshTokenDuration, "jwt-refresh-token-duration", duration, "Refresh Token Lifetime")
	return nil
}

func (cfg *Config) InitHandler(r *repository.Repository) *handler.Handler {
	return handler.NewHandler(r.AuthRepository, r.OwnerRepository, r.AdminRepository, r.ReaderRepository)
}

func (cfg *Config) InitRepository(db *gorm.DB) *repository.Repository {
	return repository.NewRepository(db)
}
