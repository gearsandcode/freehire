package config

import "os"

// Settings holds application configuration read from environment variables.
type Settings struct {
	Port        string
	DatabaseURL string
}

// Load reads configuration from the environment, falling back to sensible defaults.
func Load() Settings {
	return Settings{
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
