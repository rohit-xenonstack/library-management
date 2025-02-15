package config

import (
	"flag"
	"os"
)

type Config struct {
	Env    string
	Server ServerConfig
	DB     DbConfig
}
type ServerConfig struct {
	Port string
}
type DbConfig struct {
	DSN string
}

func NewConfig() *Config {
	return &Config{}
}

func (cfg *Config) ParseFlag() error {
	flag.StringVar(&cfg.Server.Port, "port", os.Getenv("PORT"), "API server port")
	flag.StringVar(&cfg.Env, "env", os.Getenv("ENV"), "Environment(dev|prod)")
	flag.StringVar(&cfg.DB.DSN, "db-dsn", os.Getenv("DSN"), "PostgreSQL DSN")
	return nil
}
