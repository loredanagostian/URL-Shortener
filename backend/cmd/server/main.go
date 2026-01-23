package main

import (
	"go-url-shortener/internal/api"
	"go-url-shortener/internal/core"
	"go-url-shortener/internal/db"
	"log"
	"net/http"
	"os"
)

func main() {
    // Set environment variable to force PostgreSQL usage
    os.Setenv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/urlshortener?sslmode=disable")
    
    // Get database URL from environment
    databaseURL := db.GetDatabaseURL()
    log.Printf("Attempting to connect to database: %s", databaseURL)
    
    // Initialize repository with database support - force PostgreSQL, no fallback
    repo, err := db.NewPostgresRepository(databaseURL)
    if err != nil {
        log.Fatalf("Failed to initialize PostgreSQL database: %v\nPlease ensure PostgreSQL is running on localhost:5432", err)
    }
    
    log.Printf("Successfully connected to PostgreSQL database")
    
    // Initialize shortener
    shortener := core.NewShortener(repo)
    
    // Setup router
    router := api.NewRouter(shortener, repo)
    
    // Start server
    log.Println("Server starting on :8080")
    log.Fatal(http.ListenAndServe(":8080", router))
}