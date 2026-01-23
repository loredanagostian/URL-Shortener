package db

import (
	"fmt"
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
	databaseURL = strings.TrimSpace(databaseURL)
	if databaseURL == "" {
		return nil, fmt.Errorf("database URL is required (set DATABASE_URL or POSTGRES_URL)")
	}

	// PostgreSQL only
	if strings.HasPrefix(databaseURL, "postgres://") || strings.HasPrefix(databaseURL, "postgresql://") {
		return NewPostgresRepository(databaseURL)
	}

	return nil, fmt.Errorf("unsupported database URL scheme (expected postgres:// or postgresql://)")
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
