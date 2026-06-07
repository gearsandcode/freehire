package config

import "os"

// Config holds application configuration read from environment variables.
type Config struct {
	Port        string
	DatabaseURL string
}

// Load reads configuration from the environment, falling back to sensible defaults.
func Load() Config {
	return Config{
		Port:        env("PORT", "8080"),
		DatabaseURL: env("DATABASE_URL", "postgres://hire:hire@localhost:5432/hire?sslmode=disable"),
	}
}

func env(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}
