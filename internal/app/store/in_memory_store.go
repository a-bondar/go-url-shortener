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

func (s *inMemoryStore) Ping() error {
	return nil
}

func (s *inMemoryStore) Close() error {
	return nil
}
