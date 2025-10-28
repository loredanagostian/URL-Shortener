package handlers

import (
	"encoding/json"
	"go-url-shortener/internal/db"
	"net/http"

	"github.com/gorilla/mux"
)

type URLHandler struct {
    repo db.RepositoryInterface
}

func NewURLHandler(repo db.RepositoryInterface) *URLHandler {
    return &URLHandler{repo: repo}
}

func (h *URLHandler) GetShortURL(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    shortCode := vars["shortCode"]
    
    url, err := h.repo.GetShortURL(shortCode)
    if err != nil {
        http.Error(w, "URL not found", http.StatusNotFound)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(url)
}

func (h *URLHandler) DeleteShortURL(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    shortCode := vars["shortCode"]
    
    err := h.repo.DeleteShortURL(shortCode)
    if err != nil {
        http.Error(w, "URL not found", http.StatusNotFound)
        return
    }
    
    w.WriteHeader(http.StatusNoContent)
}

func (h *URLHandler) GetAllURLs(w http.ResponseWriter, r *http.Request) {
    urls, err := h.repo.GetAllShortURLs()
    if err != nil {
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(urls)
}

func (h *URLHandler) GetURLHistory(w http.ResponseWriter, r *http.Request) {
    urls, err := h.repo.GetAllURLsHistory()
    if err != nil {
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }
    
    // Add status information to each URL
    type URLHistoryItem struct {
        *db.URL
        Status string `json:"status"`
    }
    
    var historyItems []URLHistoryItem
    for _, url := range urls {
        historyItems = append(historyItems, URLHistoryItem{
            URL:    url,
            Status: url.GetStatus(),
        })
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(historyItems)
}