package store

import (
	"context"
	"errors"
)

type inMemoryStore struct {
	m map[string]string
}

func newInMemoryStore() *inMemoryStore {
	return &inMemoryStore{
		m: make(map[string]string),
	}
}

func (s *inMemoryStore) SaveURL(_ context.Context, fullURL string, shortURL string) error {
	// Проверка на конфликт
	for currentShortURL, currentFullURL := range s.m {
		if currentFullURL == fullURL {
			return NewURLConflictError(currentShortURL, ErrConflict)
		}
	}

	s.m[shortURL] = fullURL

	return nil
}

func (s *inMemoryStore) GetURL(_ context.Context, shortURL string) (string, error) {
	if URL, ok := s.m[shortURL]; ok {
		return URL, nil
	}

	return "", errors.New("URL not found")
}

func (s *inMemoryStore) SaveURLsBatch(_ context.Context, urls map[string]string) (map[string]string, error) {
	res := make(map[string]string)

	for fullURL, shortURL := range urls {
		s.m[shortURL] = fullURL
		res[fullURL] = shortURL
	}

	return res, nil
}

func (s *inMemoryStore) Ping(_ context.Context) error {
	return nil
}

func (s *inMemoryStore) Close() {}
