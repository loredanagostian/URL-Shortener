package handlers

import (
	"net/http"
	"strings"

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

	// Record click event with visitor information
	ipAddress := getClientIP(r)
	userAgent := r.UserAgent()
	referer := r.Referer()
	
	// Add click event (ignore errors as it's not critical for redirect functionality)
	_ = h.repo.AddClickEvent(shortURL.ID, ipAddress, userAgent, referer)

	// Redirect to original URL
	http.Redirect(w, r, shortURL.OriginalURL, http.StatusMovedPermanently)
}

// getClientIP extracts the real client IP address from the request
func getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header first (when behind a proxy)
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		// X-Forwarded-For can contain multiple IPs, get the first one
		if idx := strings.Index(xff, ","); idx != -1 {
			return strings.TrimSpace(xff[:idx])
		}
		return strings.TrimSpace(xff)
	}

	// Check X-Real-IP header
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return strings.TrimSpace(xri)
	}

	// Fall back to RemoteAddr
	ip := r.RemoteAddr
	if idx := strings.LastIndex(ip, ":"); idx != -1 {
		ip = ip[:idx]
	}
	return ip
}