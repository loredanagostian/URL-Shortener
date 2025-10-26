package db

import "time"

type URL struct {
	ID           int        `json:"id"`
	ShortCode    string     `json:"short_code"`
	OriginalURL  string     `json:"original_url"`
	CreatedAt    time.Time  `json:"created_at"`
	ExpiresAt    *time.Time `json:"expires_at,omitempty"`
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