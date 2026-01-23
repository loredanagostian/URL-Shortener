package main

import (
	"go-url-shortener/internal/api"
	"go-url-shortener/internal/config"
	"go-url-shortener/internal/core"
	"go-url-shortener/internal/db"
	"log"
	"net/http"
)

func main() {
	cfg := config.Load()

	repo, err := db.InitRepository(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to initialize repository: %v", err)
	}

	shortener := core.NewShortener(repo)
	router := api.NewRouter(shortener, repo)

	addr := ":" + cfg.Port
	log.Printf("Server starting on %s", addr)
	log.Fatal(http.ListenAndServe(addr, router))
}
