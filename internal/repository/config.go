package repository

import "os"

type Config struct {
	DatabaseURL string
}

func NewConfig() *Config {
	databaseURL := os.Getenv("DATABASE_URL_LOCALHOST")
	if databaseURL == "" {
		databaseURL = "postgres://dev:qwerty@localhost:5432/dwh_service_dev"
	}

	databaseURL += "?sslmode=disable"

	return &Config{
		DatabaseURL: databaseURL,
	}
}
