package main

import (
	"go-url-shortener/internal/api"
	"go-url-shortener/internal/core"
	"go-url-shortener/internal/db"
	"log"
	"net/http"
)

func main() {
    // Initialize dependencies
    repo := db.NewMemoryRepository()
    shortener := core.NewShortener(repo)
    
    // Setup router
    router := api.NewRouter(shortener, repo)
    
    // Start server
    log.Println("Server starting on :8080")
    log.Fatal(http.ListenAndServe(":8080", router))
}