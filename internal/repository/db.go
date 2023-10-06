package repository

import (
	"database/sql"

	_ "github.com/lib/pq"
)

func NewDB(config *Config) (*sql.DB, error) {
	db, err := sql.Open("postgres", config.DatabaseURL)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
