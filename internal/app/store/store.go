package store

type Store interface {
	SaveURL(fullURL string, shortURL string) error
	GetURL(shortURL string) (string, error)
}

func NewStore(fName string) (Store, error) {
	if fName == "" {
		return newInMemoryStore(), nil
	}

	return newFileStore(fName)
}
