package config

import "os"

type Config struct {
	Port        string
	DatabaseURL string
	BaseURL     string
}

func Load() *Config {
	return &Config{
		Port: getEnv("PORT", "8080"),
		DatabaseURL: getEnv(
			"DATABASE_URL",
			getEnv("POSTGRES_URL", "postgres://postgres:postgres@localhost:5432/urlshortener?sslmode=disable"),
		),
		BaseURL: getEnv("BASE_URL", "http://localhost:8080"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return defaultValue
}
