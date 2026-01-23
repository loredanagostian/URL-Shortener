package db

import (
	"os"
	"strings"
)

// URL Operations
type RepositoryInterface interface {
	CreateShortURL(shortURL *URL) error
	GetShortURL(code string) (*URL, error)
	GetShortURLForRedirect(code string) (*URL, error)
	DeleteShortURL(code string) error
	GetAllShortURLs() ([]*URL, error)
	GetAllURLsHistory() ([]*URL, error)
	AddClickEvent(urlId int, ipAddress, userAgent, referer string) error
	GetClickEvents(urlId int) ([]*ClickEvent, error)
}

func InitRepository(databaseURL string) (RepositoryInterface, error) {
	// If no database URL provided, use memory repository
	if databaseURL == "" {
		return NewMemoryRepository(), nil
	}

	// If PostgreSQL URL provided, use PostgreSQL
	if strings.Contains(databaseURL, "postgres://") || strings.Contains(databaseURL, "postgresql://") {
		return NewPostgresRepository(databaseURL)
	}

	// Default to memory repository
	return NewMemoryRepository(), nil
}

func GetDatabaseURL() string {
	// Try to get from environment variable
	if dbURL := os.Getenv("DATABASE_URL"); dbURL != "" {
		return dbURL
	}

	// Try to get from environment variable (alternative name)
	if dbURL := os.Getenv("POSTGRES_URL"); dbURL != "" {
		return dbURL
	}

	// Default to local PostgreSQL if available
	return "postgres://postgres:postgres@localhost:5432/urlshortener?sslmode=disable"
}