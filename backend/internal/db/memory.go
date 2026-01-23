package db

import (
	"fmt"
	"sync"
	"time"
)

type MemoryRepository struct {
	urls   map[string]*URL
	nextID int
	mutex  sync.RWMutex
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		urls:   make(map[string]*URL),
		nextID: 1,
	}
}

// CreateShortURL creates a new short URL record with default 60-minute expiration
func (r *MemoryRepository) CreateShortURL(shortURL *URL) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	// Check if code already exists
	if _, exists := r.urls[shortURL.ShortCode]; exists {
		return fmt.Errorf("short code already exists")
	}

	now := time.Now()
	shortURL.ID = r.nextID
	shortURL.CreatedAt = now

	// Set default expiration (60 minutes) if not already set
	if shortURL.ExpiresAt == nil {
		shortURL.SetDefaultExpiration()
	}

	r.urls[shortURL.ShortCode] = shortURL
	r.nextID++

	return nil
}

// GetShortURL retrieves a short URL by code and checks expiration
func (r *MemoryRepository) GetShortURL(code string) (*URL, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	shortURL, exists := r.urls[code]
	if !exists {
		return nil, fmt.Errorf("short URL not found")
	}

	// Check if URL has expired
	if shortURL.IsExpired() {
		return nil, fmt.Errorf("short URL has expired")
	}

	// Return a copy to avoid external modifications
	result := *shortURL
	return &result, nil
}

// GetShortURLForRedirect retrieves a URL for redirect and increments click count
func (r *MemoryRepository) GetShortURLForRedirect(code string) (*URL, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	shortURL, exists := r.urls[code]
	if !exists {
		return nil, fmt.Errorf("short URL not found")
	}

	// Check if URL has expired
	if shortURL.IsExpired() {
		return nil, fmt.Errorf("short URL has expired")
	}

	// Increment click count
	shortURL.IncrementClickCount()

	// Return a copy
	result := *shortURL
	return &result, nil
}

// GetAllURLsHistory retrieves ALL URLs including expired ones for history
func (r *MemoryRepository) GetAllURLsHistory() ([]*URL, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	// Convert map to slice including expired URLs
	var allURLs []*URL
	for _, url := range r.urls {
		urlCopy := *url
		allURLs = append(allURLs, &urlCopy)
	}

	// Simple sorting by creation time (newest first)
	for i := 0; i < len(allURLs)-1; i++ {
		for j := i + 1; j < len(allURLs); j++ {
			if allURLs[i].CreatedAt.Before(allURLs[j].CreatedAt) {
				allURLs[i], allURLs[j] = allURLs[j], allURLs[i]
			}
		}
	}

	return allURLs, nil
}

// GetAllShortURLs retrieves all non-expired short URLs with pagination
func (r *MemoryRepository) GetAllShortURLs() ([]*URL, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	// Convert map to slice and filter out expired URLs
	var allURLs []*URL
	for _, url := range r.urls {
		if !url.IsExpired() { // Only include non-expired URLs
			urlCopy := *url
			allURLs = append(allURLs, &urlCopy)
		}
	}

	// Simple sorting by creation time (newest first)
	for i := 0; i < len(allURLs)-1; i++ {
		for j := i + 1; j < len(allURLs); j++ {
			if allURLs[i].CreatedAt.Before(allURLs[j].CreatedAt) {
				allURLs[i], allURLs[j] = allURLs[j], allURLs[i]
			}
		}
	}

	return allURLs, nil
}

// DeleteShortURL deletes a short URL by code
func (r *MemoryRepository) DeleteShortURL(code string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.urls[code]; !exists {
		return fmt.Errorf("short URL not found")
	}

	delete(r.urls, code)

	return nil
}

// AddClickEvent - memory repository doesn't persist click events
func (r *MemoryRepository) AddClickEvent(urlId int, ipAddress, userAgent, referer string) error {
	// In memory repository, we don't store individual click events
	// Click counts are updated in GetShortURLForRedirect
	return nil
}

// GetClickEvents - memory repository doesn't store individual click events
func (r *MemoryRepository) GetClickEvents(urlId int) ([]*ClickEvent, error) {
	// In memory repository, we don't store individual click events
	return []*ClickEvent{}, nil
}