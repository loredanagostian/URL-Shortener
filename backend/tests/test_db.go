package main

import (
	"fmt"
	"go-url-shortener/internal/db"
)

func main() {
	databaseURL := db.GetDatabaseURL()

	repo, err := db.InitRepository(databaseURL)
	if err != nil {
		fmt.Println("Failed to initialize database repository:", err)
		return
	}

	// Test 1: Create a short URL
	url := &db.URL{
		OriginalURL: "https://example.com",
		ShortCode:   "exmpl",
	}

	err = repo.CreateShortURL(url)

	if err != nil {
		fmt.Println("[1] Error creating short URL:", err)
		return
	}

	fmt.Printf("[1] Created URL: ID=%d, Code=%s, Original=%s\n", url.ID, url.ShortCode, url.OriginalURL)
	fmt.Printf("[1] Created at: %s\n", url.CreatedAt.Format("2006-01-02 15:04:05"))

	// Test 2: Try to create a duplicate short URL
	urlDup := &db.URL{
		OriginalURL: "https://duplicate.com",
		ShortCode:   "exmpl",
	}

	err = repo.CreateShortURL(urlDup)

	if err != nil {
		fmt.Println("[2] Correctly caught duplicate short code:", err)
	} else {
		fmt.Println("[2] Error: Duplicate short code was allowed")
	}

	// Test 3: Retrieve the created short URL
	retrievedURL, err := repo.GetShortURL("exmpl")

	if err != nil {
		fmt.Println("[3] Error retrieving short URL:", err)
		return
	}

	fmt.Printf("[3] Retrieved URL: ID=%d, Code=%s, Original=%s\n", retrievedURL.ID, retrievedURL.ShortCode, retrievedURL.OriginalURL)

	// Test 4: Attempt to retrieve a non-existent short URL
	_, err = repo.GetShortURL("nonexist")
	if err != nil {
		fmt.Println("[4] Correctly handled non-existent short URL:", err)
	} else {
		fmt.Println("[4] Error: Non-existent short URL was found")
	}

	// Test 5: Create a short URL with expiration and check expiration
	expiringURL := &db.URL{
		OriginalURL: "https://tempurl.com",
		ShortCode:   "temp",
	}

	expiringURL.SetDefaultExpiration() // Set to expire in 60 minutes

	err = repo.CreateShortURL(expiringURL)

	if err != nil {
		fmt.Println("[5] Error creating expiring short URL:", err)
		return
	}

	fmt.Printf("[5] Created expiring URL: ID=%d, Code=%s, ExpiresAt=%s\n", expiringURL.ID, expiringURL.ShortCode, expiringURL.ExpiresAt.Format("2006-01-02 15:04:05"))

	// Simulate expiration check
	if expiringURL.IsExpired() {
		fmt.Println("[5] Error: URL should not be expired yet")
	} else {
		fmt.Println("[5] URL is not expired as expected")
	}

	// Test 6: Delete the created short URL
	err = repo.DeleteShortURL("exmpl")

	if err != nil {
		fmt.Println("[6] Error deleting short URL:", err)
		return
	}

	fmt.Println("[6] Deleted short URL: Code=exmpl")

	// Verify deletion
	_, err = repo.GetShortURL("exmpl")

	if err != nil {
		fmt.Println("[6] Correctly confirmed deletion of short URL:", err)
	} else {
		fmt.Println("[6] Error: Deleted short URL was still found")
	}

	// Test 7: Retrieve all short URLs
	allURLs, err := repo.GetAllShortURLs()

	if err != nil {
		fmt.Println("[7] Error retrieving all short URLs:", err)
		return
	}

	fmt.Printf("[7] Retrieved all short URLs, count=%d\n", len(allURLs))

	for _, u := range allURLs {
		fmt.Printf("[7] ID=%d, Code=%s, Original=%s\n", u.ID, u.ShortCode, u.OriginalURL)
	}

	fmt.Println("=== All database layer tests completed. ===")
}
