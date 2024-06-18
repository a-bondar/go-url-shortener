package store

import "errors"

type inMemoryStore struct {
	m map[string]string
}

func newInMemoryStore() *inMemoryStore {
	return &inMemoryStore{
		m: make(map[string]string),
	}
}

func (s *inMemoryStore) SaveURL(fullURL string, shortURL string) error {
	s.m[shortURL] = fullURL

	return nil
}

func (s *inMemoryStore) GetURL(shortURL string) (string, error) {
	if URL, ok := s.m[shortURL]; ok {
		return URL, nil
	}

	return "", errors.New("URL not found")
}

func (s *inMemoryStore) SaveURLsBatch(urls map[string]string) (map[string]string, error) {
	res := make(map[string]string)

	for fullURL, shortURL := range urls {
		s.m[shortURL] = fullURL
		res[fullURL] = shortURL
	}

	return res, nil
}

func (s *inMemoryStore) Ping() error {
	return nil
}

func (s *inMemoryStore) Close() error {
	return nil
}
