package store

type Store struct {
	m map[string]string
}

func NewStore() *Store {
	return &Store{
		m: make(map[string]string),
	}
}

func (s *Store) SaveURL(fullURL string, shortURL string) error {
	return nil
}

func (s *Store) GetURL(shortURL string) (string, error) {
	return "", nil
}
