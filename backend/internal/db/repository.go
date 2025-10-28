package db

// URL Operations
type RepositoryInterface interface {
	CreateShortURL(shortURL *URL) error
	GetShortURL(code string) (*URL, error)
	GetShortURLForRedirect(code string) (*URL, error)
	DeleteShortURL(code string) error
	GetAllShortURLs() ([]*URL, error)
	GetAllURLsHistory() ([]*URL, error)
}

func InitRepository(databaseURL string) (RepositoryInterface, error) {
	return NewMemoryRepository(), nil
}