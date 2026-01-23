package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"go-url-shortener/internal/db"
)

type AnalyticsHandler struct {
	repo db.RepositoryInterface
}

func NewAnalyticsHandler(repo db.RepositoryInterface) *AnalyticsHandler {
	return &AnalyticsHandler{
		repo: repo,
	}
}

type AnalyticsResponse struct {
	URL         *db.URL           `json:"url"`
	ClickEvents []*db.ClickEvent  `json:"click_events"`
	Summary     AnalyticsSummary  `json:"summary"`
}

type AnalyticsSummary struct {
	TotalClicks     int                    `json:"total_clicks"`
	UniqueIPs       int                    `json:"unique_ips"`
	TopReferers     map[string]int         `json:"top_referers"`
	TopUserAgents   map[string]int         `json:"top_user_agents"`
	ClicksByHour    map[string]int         `json:"clicks_by_hour"`
}

// GetURLAnalytics returns analytics data for a specific short URL
func (h *AnalyticsHandler) GetURLAnalytics(w http.ResponseWriter, r *http.Request) {
	// Extract short code from URL path
	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(pathParts) < 3 { // /api/analytics/{shortCode}
		http.Error(w, "Short code is required", http.StatusBadRequest)
		return
	}
	
	shortCode := pathParts[2]

	// Get URL information
	url, err := h.repo.GetShortURL(shortCode)
	if err != nil {
		http.Error(w, "Short URL not found", http.StatusNotFound)
		return
	}

	// Get click events
	clickEvents, err := h.repo.GetClickEvents(url.ID)
	if err != nil {
		http.Error(w, "Failed to retrieve click events", http.StatusInternalServerError)
		return
	}

	// Generate summary analytics
	summary := generateSummary(clickEvents)

	response := AnalyticsResponse{
		URL:         url,
		ClickEvents: clickEvents,
		Summary:     summary,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// generateSummary creates analytics summary from click events
func generateSummary(events []*db.ClickEvent) AnalyticsSummary {
	summary := AnalyticsSummary{
		TotalClicks:   len(events),
		TopReferers:   make(map[string]int),
		TopUserAgents: make(map[string]int),
		ClicksByHour:  make(map[string]int),
	}

	uniqueIPs := make(map[string]bool)

	for _, event := range events {
		// Count unique IPs
		uniqueIPs[event.IPAddress] = true

		// Count referers
		referer := event.Referer
		if referer == "" {
			referer = "Direct"
		}
		summary.TopReferers[referer]++

		// Count user agents (simplified - browser name only)
		browser := extractBrowser(event.UserAgent)
		summary.TopUserAgents[browser]++

		// Count clicks by hour
		hour := event.CreatedAt.Format("2006-01-02 15:00")
		summary.ClicksByHour[hour]++
	}

	summary.UniqueIPs = len(uniqueIPs)

	return summary
}

// extractBrowser extracts browser name from user agent string (simplified)
func extractBrowser(userAgent string) string {
	userAgent = strings.ToLower(userAgent)
	
	if strings.Contains(userAgent, "chrome") {
		return "Chrome"
	} else if strings.Contains(userAgent, "firefox") {
		return "Firefox"
	} else if strings.Contains(userAgent, "safari") {
		return "Safari"
	} else if strings.Contains(userAgent, "edge") {
		return "Edge"
	} else if strings.Contains(userAgent, "opera") {
		return "Opera"
	}
	
	return "Other"
}