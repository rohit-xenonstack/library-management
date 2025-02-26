package config

import "time"

var SampleEnv = Config{
	Env:    "prod",
	Server: ServerConfig{Port: ":8081"},
	DB:     DbConfig{DSN: "host=localhost user=postgres password=postgres dbname=library port=5433 sslmode=disable"},
	JWT: JWTConfig{
		SecretKey:           "bgeab3wbna3gh3p83hw8hgf83hg8hp8ghp8g38w8h3",
		AccessTokenDuration: 30 * time.Minute,
	},
}
