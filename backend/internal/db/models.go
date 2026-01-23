package db

import "time"

type URL struct {
	ID           int        `json:"id" db:"id"`
	ShortCode    string     `json:"short_code" db:"short_code"`
	OriginalURL  string     `json:"original_url" db:"original_url"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	ExpiresAt    *time.Time `json:"expires_at,omitempty" db:"expires_at"`
	ClickCount   int        `json:"click_count" db:"click_count"`
	LastClicked  *time.Time `json:"last_clicked,omitempty" db:"last_clicked"`
}

type ClickEvent struct {
	ID        int       `json:"id" db:"id"`
	URLId     int       `json:"url_id" db:"url_id"`
	IPAddress string    `json:"ip_address" db:"ip_address"`
	UserAgent string    `json:"user_agent" db:"user_agent"`
	Referer   string    `json:"referer" db:"referer"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// SetDefaultExpiration sets the expiration time to 60 minutes from now
func (s *URL) SetDefaultExpiration() {
	expiration := time.Now().Add(60 * time.Minute)
	s.ExpiresAt = &expiration
}

// IsExpired checks if the short URL has expired
func (s *URL) IsExpired() bool {
	if s.ExpiresAt == nil {
		return false
	}

	return s.ExpiresAt.Before(time.Now())
}
func (s *URL) GetStatus() string {
	if s.IsExpired() {
		return "expired"
	}
	return "active"
}

// IncrementClickCount increments the click count and updates last clicked time
func (s *URL) IncrementClickCount() {
	s.ClickCount++
	now := time.Now()
	s.LastClicked = &now
}