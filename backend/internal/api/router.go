package api

import (
	"go-url-shortener/internal/api/handlers"
	"go-url-shortener/internal/core"
	"go-url-shortener/internal/db"
	"go-url-shortener/internal/middleware"

	"github.com/gorilla/mux"
)

func NewRouter(shortener *core.Shortener, repo db.RepositoryInterface) *mux.Router {
    r := mux.NewRouter()
    
    // Add CORS middleware
    r.Use(middleware.CORS)
    
    // Initialize handlers
    shortenHandler := handlers.NewShortenHandler(shortener, repo)
    redirectHandler := handlers.NewRedirectHandler(repo)
    urlHandler := handlers.NewURLHandler(repo)
    
    // API routes
    api := r.PathPrefix("/api").Subrouter()
    api.HandleFunc("/shorten", shortenHandler.CreateShortURL).Methods("POST")
    api.HandleFunc("/urls/{shortCode}", urlHandler.GetShortURL).Methods("GET")
    api.HandleFunc("/urls/{shortCode}", urlHandler.DeleteShortURL).Methods("DELETE")
    api.HandleFunc("/urls", urlHandler.GetAllURLs).Methods("GET")
    api.HandleFunc("/history", urlHandler.GetURLHistory).Methods("GET")
    
    // Redirect route
    r.HandleFunc("/{shortCode}", redirectHandler.RedirectToOriginal).Methods("GET")
    
    return r
}
