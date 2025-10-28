package handlers

import (
	"net/http"

	"go-url-shortener/internal/db"
)

type RedirectHandler struct {
	repo      db.RepositoryInterface
}

func NewRedirectHandler(repo db.RepositoryInterface) *RedirectHandler {
	return &RedirectHandler{
		repo: repo,
	}
}

func (h *RedirectHandler) RedirectToOriginal(w http.ResponseWriter, r *http.Request) {
	// Extract short code from URL path
	shortCode := r.URL.Path[1:]
	
	if shortCode == "" {
		http.Error(w, "Short code is required", http.StatusBadRequest)
		return
	}

	// Get original URL from database and increment click count
	shortURL, err := h.repo.GetShortURLForRedirect(shortCode)
	if err != nil {
		http.Error(w, "Short URL not found", http.StatusNotFound)
		return
	}

	// Check if URL is expired (already handled in GetShortURLForRedirect, but double-check)
	if shortURL.IsExpired() {
		http.Error(w, "Short URL has expired", http.StatusGone)
		return
	}

	// Redirect to original URL
	http.Redirect(w, r, shortURL.OriginalURL, http.StatusMovedPermanently)
}