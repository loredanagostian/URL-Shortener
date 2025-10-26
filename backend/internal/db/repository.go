package db

// URL Operations
type Repository interface {
	CreateShortURL(shortURL *URL) error
	GetShortURL(code string) (*URL, error)
	DeleteShortURL(code string) error
	GetAllShortURLs() ([]*URL, error)
}

func InitRepository(databaseURL string) (Repository, error) {
	return NewMemoryRepository(), nil
}