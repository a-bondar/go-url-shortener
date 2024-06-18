package store

type Store interface {
	SaveURL(fullURL string, shortURL string) error
	GetURL(shortURL string) (string, error)
	SaveURLsBatch(urls map[string]string) (map[string]string, error)
	Ping() error
	Close() error
}

func NewStore(dsn string, fName string) (Store, error) {
	if dsn != "" {
		return newDBStore(dsn)
	}

	if fName != "" {
		return newFileStore(fName)
	}

	return newInMemoryStore(), nil
}
