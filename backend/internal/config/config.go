// Package config loads runtime configuration from the environment.
package config

import "os"

// Config holds every runtime setting the server needs.
type Config struct {
	DatabaseURL string
	Port        string
}

// Load reads configuration from the environment, applying local-friendly defaults.
func Load() Config {
	return Config{
		DatabaseURL: env("DATABASE_URL", "postgres://tutor:tutor@localhost:5434/english_tutor?sslmode=disable"),
		Port:        env("PORT", "8096"),
	}
}

func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
