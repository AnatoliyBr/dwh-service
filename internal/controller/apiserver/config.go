package apiserver

import "time"

type Config struct {
	ReadTimeout  time.Duration `toml:"read_timeout"`
	WriteTimeout time.Duration `toml:"write_timeout"`
	BindAddr     string        `toml:"bind_addr"`

	ShutdownTimeout time.Duration `toml:"shutdown_time"`
	LogLevel        string        `toml:"log_level"`
}

func NewConfig() *Config {
	return &Config{
		ReadTimeout:     5 * time.Second,
		WriteTimeout:    5 * time.Second,
		BindAddr:        ":8080",
		ShutdownTimeout: 3 * time.Second,
		LogLevel:        "debug",
	}
}
