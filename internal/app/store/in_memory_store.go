package store

import "errors"

type InMemoryStore struct {
	m map[string]string
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		m: make(map[string]string),
	}
}

func (s *InMemoryStore) SaveURL(fullURL string, shortURL string) error {
	s.m[shortURL] = fullURL

	return nil
}

func (s *InMemoryStore) GetURL(shortURL string) (string, error) {
	if URL, ok := s.m[shortURL]; ok {
		return URL, nil
	}

	return "", errors.New("URL not found")
}
