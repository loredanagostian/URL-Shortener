package tests

import (
	"go-url-shortener/internal/db"
	"testing"
)

// setupTestRepo initializes a test repository connection
func setupTestRepo(t *testing.T) db.RepositoryInterface {
	databaseURL := db.GetDatabaseURL()
	repo, err := db.InitRepository(databaseURL)
	if err != nil {
		t.Fatalf("Failed to initialize database repository: %v", err)
	}
	return repo
}

// cleanupTestURL removes a test URL by short code
func cleanupTestURL(t *testing.T, repo db.RepositoryInterface, shortCode string) {
	_ = repo.DeleteShortURL(shortCode)
}

func TestCreateShortURL(t *testing.T) {
	repo := setupTestRepo(t)
	defer cleanupTestURL(t, repo, "tstcrt")

	url := &db.URL{
		OriginalURL: "https://example.com/create-test",
		ShortCode:   "tstcrt",
	}

	err := repo.CreateShortURL(url)
	if err != nil {
		t.Errorf("Error creating short URL: %v", err)
		return
	}

	if url.ID == 0 {
		t.Error("Expected URL ID to be set after creation")
	}

	if url.ShortCode != "tstcrt" {
		t.Errorf("Expected short code 'tstcrt', got '%s'", url.ShortCode)
	}
}

func TestCreateDuplicateShortCode(t *testing.T) {
	repo := setupTestRepo(t)
	defer cleanupTestURL(t, repo, "tstdup")

	url := &db.URL{
		OriginalURL: "https://example.com/original",
		ShortCode:   "tstdup",
	}

	err := repo.CreateShortURL(url)
	if err != nil {
		t.Fatalf("Failed to create initial URL: %v", err)
	}

	urlDup := &db.URL{
		OriginalURL: "https://duplicate.com",
		ShortCode:   "tstdup",
	}

	err = repo.CreateShortURL(urlDup)
	if err == nil {
		t.Error("Expected error when creating duplicate short code, got nil")
	}
}

func TestGetShortURL(t *testing.T) {
	repo := setupTestRepo(t)
	defer cleanupTestURL(t, repo, "tstget")

	url := &db.URL{
		OriginalURL: "https://example.com/get-test",
		ShortCode:   "tstget",
	}

	err := repo.CreateShortURL(url)
	if err != nil {
		t.Fatalf("Failed to create URL for retrieval test: %v", err)
	}

	retrievedURL, err := repo.GetShortURL("tstget")
	if err != nil {
		t.Errorf("Error retrieving short URL: %v", err)
		return
	}

	if retrievedURL.OriginalURL != url.OriginalURL {
		t.Errorf("Expected original URL '%s', got '%s'", url.OriginalURL, retrievedURL.OriginalURL)
	}

	if retrievedURL.ShortCode != url.ShortCode {
		t.Errorf("Expected short code '%s', got '%s'", url.ShortCode, retrievedURL.ShortCode)
	}
}

func TestGetNonExistentShortURL(t *testing.T) {
	repo := setupTestRepo(t)

	_, err := repo.GetShortURL("nonexist")
	if err == nil {
		t.Error("Expected error when retrieving non-existent short URL, got nil")
	}
}

func TestCreateExpiringURL(t *testing.T) {
	repo := setupTestRepo(t)
	defer cleanupTestURL(t, repo, "tstexp")

	expiringURL := &db.URL{
		OriginalURL: "https://tempurl.com",
		ShortCode:   "tstexp",
	}

	expiringURL.SetDefaultExpiration() // Set to expire in 60 minutes

	err := repo.CreateShortURL(expiringURL)
	if err != nil {
		t.Errorf("Error creating expiring short URL: %v", err)
		return
	}

	if expiringURL.ExpiresAt.IsZero() {
		t.Error("Expected ExpiresAt to be set")
	}

	// URL should not be expired immediately after creation
	if expiringURL.IsExpired() {
		t.Error("URL should not be expired immediately after creation")
	}
}

func TestDeleteShortURL(t *testing.T) {
	repo := setupTestRepo(t)

	url := &db.URL{
		OriginalURL: "https://example.com/delete-test",
		ShortCode:   "tstdel",
	}

	err := repo.CreateShortURL(url)
	if err != nil {
		t.Fatalf("Failed to create URL for deletion test: %v", err)
	}

	err = repo.DeleteShortURL("tstdel")
	if err != nil {
		t.Errorf("Error deleting short URL: %v", err)
		return
	}

	// Verify deletion
	_, err = repo.GetShortURL("tstdel")
	if err == nil {
		t.Error("Expected error when retrieving deleted short URL, got nil")
	}
}

func TestGetAllShortURLs(t *testing.T) {
	repo := setupTestRepo(t)

	allURLs, err := repo.GetAllShortURLs()
	if err != nil {
		t.Errorf("Error retrieving all short URLs: %v", err)
		return
	}

	// Just verify the function works and returns a slice
	if allURLs == nil {
		t.Error("Expected non-nil slice from GetAllShortURLs")
	}
}
