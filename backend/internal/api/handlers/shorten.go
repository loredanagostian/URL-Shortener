package handlers

import (
	"encoding/json"
	"go-url-shortener/internal/core"
	"go-url-shortener/internal/db"
	"net/http"
	"time"
)

type ShortenRequest struct {
    URL        string `json:"url" validate:"required,url"`
    CustomCode string `json:"custom_code,omitempty"`
    ExpiresAt  string `json:"expires_at,omitempty"`
}

type ShortenResponse struct {
    ShortURL    string `json:"short_url"`
    OriginalURL string `json:"original_url"`
    Code        string `json:"code"`
    QRCode      string `json:"qr_code,omitempty"`
}

type ShortenHandler struct {
    shortener *core.Shortener
    repo      db.RepositoryInterface
}

func NewShortenHandler(shortener *core.Shortener, repo db.RepositoryInterface) *ShortenHandler {
    return &ShortenHandler{
        shortener: shortener,
        repo:      repo,
    }
}

func (h *ShortenHandler) CreateShortURL(w http.ResponseWriter, r *http.Request) {
    var req ShortenRequest
    
    // Parse JSON request
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }
    
    // Parse expiration if provided
    var expiresAt *time.Time
    if req.ExpiresAt != "" {
        if parsed, err := time.Parse(time.RFC3339, req.ExpiresAt); err != nil {
            http.Error(w, "Invalid expiration date format", http.StatusBadRequest)
            return
        } else {
            expiresAt = &parsed
        }
    }
    
    // Create URL through business logic
    url, err := h.shortener.CreateShortURL(req.URL, req.CustomCode, expiresAt)
    if err != nil {
        if err.Error() == "short code already exists" {
            http.Error(w, err.Error(), http.StatusConflict)
            return
        }
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }
    
    // Prepare response
    response := ShortenResponse{
        ShortURL:    "http://localhost:8080/" + url.ShortCode, // TODO: Use config for base URL
        OriginalURL: url.OriginalURL,
        Code:        url.ShortCode,
    }
    
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(response)
}