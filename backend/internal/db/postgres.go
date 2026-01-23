package db

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(databaseURL string) (*PostgresRepository, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	repo := &PostgresRepository{db: db}
	
	// Initialize database schema
	if err := repo.initSchema(); err != nil {
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	return repo, nil
}

func (r *PostgresRepository) initSchema() error {
	// Create urls table
	urlsSchema := `
	CREATE TABLE IF NOT EXISTS urls (
		id SERIAL PRIMARY KEY,
		short_code VARCHAR(50) UNIQUE NOT NULL,
		original_url TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		expires_at TIMESTAMP,
		click_count INTEGER DEFAULT 0,
		last_clicked TIMESTAMP
	);`

	if _, err := r.db.Exec(urlsSchema); err != nil {
		return fmt.Errorf("failed to create urls table: %w", err)
	}

	// Create click_events table
	clickEventsSchema := `
	CREATE TABLE IF NOT EXISTS click_events (
		id SERIAL PRIMARY KEY,
		url_id INTEGER REFERENCES urls(id) ON DELETE CASCADE,
		ip_address INET,
		user_agent TEXT,
		referer TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	if _, err := r.db.Exec(clickEventsSchema); err != nil {
		return fmt.Errorf("failed to create click_events table: %w", err)
	}

	// Create indexes for better performance
	indexes := []string{
		"CREATE INDEX IF NOT EXISTS idx_urls_short_code ON urls(short_code);",
		"CREATE INDEX IF NOT EXISTS idx_click_events_url_id ON click_events(url_id);",
		"CREATE INDEX IF NOT EXISTS idx_urls_expires_at ON urls(expires_at);",
	}

	for _, idx := range indexes {
		if _, err := r.db.Exec(idx); err != nil {
			return fmt.Errorf("failed to create index: %w", err)
		}
	}

	return nil
}

func (r *PostgresRepository) CreateShortURL(shortURL *URL) error {
	query := `
		INSERT INTO urls (short_code, original_url, created_at, expires_at, click_count)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id`

	now := time.Now()
	shortURL.CreatedAt = now

	// Set default expiration if not set
	if shortURL.ExpiresAt == nil {
		shortURL.SetDefaultExpiration()
	}

	err := r.db.QueryRow(
		query,
		shortURL.ShortCode,
		shortURL.OriginalURL,
		shortURL.CreatedAt,
		shortURL.ExpiresAt,
		shortURL.ClickCount,
	).Scan(&shortURL.ID)

	if err != nil {
		return fmt.Errorf("failed to create short URL: %w", err)
	}

	return nil
}

func (r *PostgresRepository) GetShortURL(code string) (*URL, error) {
	query := `
		SELECT id, short_code, original_url, created_at, expires_at, click_count, last_clicked
		FROM urls 
		WHERE short_code = $1`

	url := &URL{}
	err := r.db.QueryRow(query, code).Scan(
		&url.ID,
		&url.ShortCode,
		&url.OriginalURL,
		&url.CreatedAt,
		&url.ExpiresAt,
		&url.ClickCount,
		&url.LastClicked,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("short URL not found")
		}
		return nil, fmt.Errorf("failed to get short URL: %w", err)
	}

	// Check if expired
	if url.IsExpired() {
		return nil, fmt.Errorf("short URL has expired")
	}

	return url, nil
}

func (r *PostgresRepository) GetShortURLForRedirect(code string) (*URL, error) {
	// Start transaction for atomic read and update
	tx, err := r.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()

	// Get URL
	url, err := r.GetShortURL(code)
	if err != nil {
		return nil, err
	}

	// Update click count and last clicked
	updateQuery := `
		UPDATE urls 
		SET click_count = click_count + 1, last_clicked = $1 
		WHERE short_code = $2`

	now := time.Now()
	if _, err := tx.Exec(updateQuery, now, code); err != nil {
		return nil, fmt.Errorf("failed to update click count: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	url.ClickCount++
	url.LastClicked = &now

	return url, nil
}

func (r *PostgresRepository) DeleteShortURL(code string) error {
	query := `DELETE FROM urls WHERE short_code = $1`
	result, err := r.db.Exec(query, code)
	if err != nil {
		return fmt.Errorf("failed to delete short URL: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("short URL not found")
	}

	return nil
}

func (r *PostgresRepository) GetAllShortURLs() ([]*URL, error) {
	query := `
		SELECT id, short_code, original_url, created_at, expires_at, click_count, last_clicked
		FROM urls 
		WHERE expires_at IS NULL OR expires_at > CURRENT_TIMESTAMP
		ORDER BY created_at DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all URLs: %w", err)
	}
	defer rows.Close()

	var urls []*URL
	for rows.Next() {
		url := &URL{}
		err := rows.Scan(
			&url.ID,
			&url.ShortCode,
			&url.OriginalURL,
			&url.CreatedAt,
			&url.ExpiresAt,
			&url.ClickCount,
			&url.LastClicked,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan URL: %w", err)
		}
		urls = append(urls, url)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate rows: %w", err)
	}

	return urls, nil
}

func (r *PostgresRepository) GetAllURLsHistory() ([]*URL, error) {
	query := `
		SELECT id, short_code, original_url, created_at, expires_at, click_count, last_clicked
		FROM urls 
		ORDER BY created_at DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get URL history: %w", err)
	}
	defer rows.Close()

	var urls []*URL
	for rows.Next() {
		url := &URL{}
		err := rows.Scan(
			&url.ID,
			&url.ShortCode,
			&url.OriginalURL,
			&url.CreatedAt,
			&url.ExpiresAt,
			&url.ClickCount,
			&url.LastClicked,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan URL: %w", err)
		}
		urls = append(urls, url)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate rows: %w", err)
	}

	return urls, nil
}

// AddClickEvent records a click event
func (r *PostgresRepository) AddClickEvent(urlId int, ipAddress, userAgent, referer string) error {
	query := `
		INSERT INTO click_events (url_id, ip_address, user_agent, referer, created_at)
		VALUES ($1, $2, $3, $4, $5)`

	_, err := r.db.Exec(query, urlId, ipAddress, userAgent, referer, time.Now())
	if err != nil {
		return fmt.Errorf("failed to add click event: %w", err)
	}

	return nil
}

// GetClickEvents returns click events for a URL
func (r *PostgresRepository) GetClickEvents(urlId int) ([]*ClickEvent, error) {
	query := `
		SELECT id, url_id, ip_address, user_agent, referer, created_at
		FROM click_events 
		WHERE url_id = $1 
		ORDER BY created_at DESC`

	rows, err := r.db.Query(query, urlId)
	if err != nil {
		return nil, fmt.Errorf("failed to get click events: %w", err)
	}
	defer rows.Close()

	var events []*ClickEvent
	for rows.Next() {
		event := &ClickEvent{}
		err := rows.Scan(
			&event.ID,
			&event.URLId,
			&event.IPAddress,
			&event.UserAgent,
			&event.Referer,
			&event.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan click event: %w", err)
		}
		events = append(events, event)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate rows: %w", err)
	}

	return events, nil
}

func (r *PostgresRepository) Close() error {
	return r.db.Close()
}