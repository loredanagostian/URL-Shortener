package core

import (
	"crypto/rand"
	"errors"
	"fmt"
	"go-url-shortener/internal/db"
	"net/url"
	"strings"
	"time"
)

type Shortener struct {
    baseURL string
    repo    db.RepositoryInterface
}

func NewShortener(repo db.RepositoryInterface) *Shortener {
    return &Shortener{
        baseURL: "http://localhost:8080", // TODO: Move to config
        repo:    repo,
    }
}

// CreateShortURL is the main business logic method
func (s *Shortener) CreateShortURL(originalURL, customCode string, expiresAt *time.Time) (*db.URL, error) {
    // Validate URL
    if !s.isValidURL(originalURL) {
        return nil, errors.New("invalid URL format")
    }

    // Generate or validate short code
    shortCode, err := s.GenerateShortCode(originalURL, customCode)
    if err != nil {
        return nil, err
    }

    // Check if short code already exists
    existingURL, err := s.repo.GetShortURL(shortCode)
    if err == nil && existingURL != nil {
        return nil, errors.New("short code already exists")
    }

    // Create URL object
    urlObj := &db.URL{
        OriginalURL: originalURL,
        ShortCode:   shortCode,
        CreatedAt:   time.Now(),
    }

    // Set expiration
    if expiresAt != nil {
        urlObj.ExpiresAt = expiresAt
    } else {
        urlObj.SetDefaultExpiration() // 60 minutes default
    }

    // Save to database
    err = s.repo.CreateShortURL(urlObj)
    if err != nil {
        return nil, fmt.Errorf("failed to create short URL: %w", err)
    }

    return urlObj, nil
}

// GenerateShortCode creates a short code for the given URL
func (s *Shortener) GenerateShortCode(originalURL, customCode string) (string, error) {
    // Use custom code if provided
    if customCode != "" {
        if !s.isValidCustomCode(customCode) {
            return "", errors.New("invalid custom code format")
        }
        return customCode, nil
    }

    // Generate random code with retry logic for uniqueness
    maxRetries := 5
    for i := 0; i < maxRetries; i++ {
        code, err := s.generateRandomCode(6)
        if err != nil {
            return "", err
        }

        // Check if code already exists
        _, err = s.repo.GetShortURL(code)
        if err != nil {
            // Code doesn't exist, we can use it
            return code, nil
        }
    }

    return "", errors.New("failed to generate unique short code after retries")
}

// GetOriginalURL retrieves and validates the original URL
func (s *Shortener) GetOriginalURL(shortCode string) (string, error) {
    urlObj, err := s.repo.GetShortURL(shortCode)
    if err != nil {
        return "", errors.New("short URL not found")
    }

    if urlObj.IsExpired() {
        return "", errors.New("short URL has expired")
    }

    return urlObj.OriginalURL, nil
}

// DeleteShortURL removes a short URL
func (s *Shortener) DeleteShortURL(shortCode string) error {
    return s.repo.DeleteShortURL(shortCode)
}

// GetAllShortURLs returns all URLs (for admin/management)
func (s *Shortener) GetAllShortURLs() ([]*db.URL, error) {
    return s.repo.GetAllShortURLs()
}

// ValidateAndSanitizeURL performs comprehensive URL validation
func (s *Shortener) ValidateAndSanitizeURL(urlStr string) (string, error) {
    // Basic format validation
    if !s.isValidURL(urlStr) {
        return "", errors.New("invalid URL format")
    }

    // Parse URL for further validation
    parsedURL, err := url.Parse(urlStr)
    if err != nil {
        return "", errors.New("failed to parse URL")
    }

    // Security checks
    if err := s.checkURLSecurity(parsedURL); err != nil {
        return "", err
    }

    // Return cleaned URL
    return parsedURL.String(), nil
}

// checkURLSecurity performs security validation
func (s *Shortener) checkURLSecurity(parsedURL *url.URL) error {
    // Block localhost and internal IPs
    if strings.Contains(parsedURL.Host, "localhost") || 
       strings.Contains(parsedURL.Host, "127.0.0.1") ||
       strings.Contains(parsedURL.Host, "0.0.0.0") {
        return errors.New("localhost URLs are not allowed")
    }

    // Block known malicious domains (example)
    blockedDomains := []string{
        "malware.com",
        "phishing.net",
        // Add more as needed
    }

    for _, domain := range blockedDomains {
        if strings.Contains(parsedURL.Host, domain) {
            return errors.New("URL domain is blocked")
        }
    }

    return nil
}

// isValidURL checks if the URL is valid
func (s *Shortener) isValidURL(urlStr string) bool {
    _, err := url.ParseRequestURI(urlStr)
    if err != nil {
        return false
    }

    u, err := url.Parse(urlStr)
    if err != nil || u.Scheme == "" || u.Host == "" {
        return false
    }

    return u.Scheme == "http" || u.Scheme == "https"
}

// isValidCustomCode checks if custom code is valid
func (s *Shortener) isValidCustomCode(code string) bool {
    if len(code) < 3 || len(code) > 20 {
        return false
    }

    // Only allow alphanumeric characters and hyphens
    for _, char := range code {
        if !((char >= 'a' && char <= 'z') || 
             (char >= 'A' && char <= 'Z') || 
             (char >= '0' && char <= '9') || 
             char == '-') {
            return false
        }
    }

    // Don't allow codes that might conflict with API routes
    reservedCodes := []string{"api", "admin", "health", "metrics"}
    for _, reserved := range reservedCodes {
        if strings.ToLower(code) == reserved {
            return false
        }
    }

    return true
}

// generateRandomCode creates a random alphanumeric code
func (s *Shortener) generateRandomCode(length int) (string, error) {
    const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
    bytes := make([]byte, length)
    
    for i := range bytes {
        randomIndex := make([]byte, 1)
        if _, err := rand.Read(randomIndex); err != nil {
            return "", err
        }
        bytes[i] = charset[randomIndex[0]%byte(len(charset))]
    }

    return string(bytes), nil
}

// GetShortURL builds the complete short URL
func (s *Shortener) GetShortURL(code string) string {
    return fmt.Sprintf("%s/%s", s.baseURL, code)
}

// SetBaseURL allows changing the base URL (useful for different environments)
func (s *Shortener) SetBaseURL(baseURL string) {
    s.baseURL = strings.TrimSuffix(baseURL, "/")
}