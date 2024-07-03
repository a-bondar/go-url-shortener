package store

import (
	"context"
	"fmt"
)

type inMemoryStore struct {
	m map[string]map[string]string
}

func newInMemoryStore() *inMemoryStore {
	return &inMemoryStore{
		m: make(map[string]map[string]string),
	}
}

func (s *inMemoryStore) SaveURL(_ context.Context, fullURL string, shortURL string, userID string) (string, error) {
	userURLs, ok := s.m[userID]
	if !ok {
		// У пользователя еще нет данных, создаем пустой map
		userURLs = make(map[string]string)
		s.m[userID] = userURLs
	}

	// Проверка на конфликт
	for currentShortURL, currentFullURL := range userURLs {
		if currentFullURL == fullURL {
			return currentShortURL, nil
		}
	}

	userURLs[shortURL] = fullURL

	return shortURL, nil
}

func (s *inMemoryStore) GetURL(_ context.Context, shortURL string, userID string) (string, error) {
	userURLs, ok := s.m[userID]
	if !ok {
		return "", fmt.Errorf("%w", ErrUserHasNoURLs)
	}

	fullURL, ok := userURLs[shortURL]
	if !ok {
		return "", fmt.Errorf("%w", ErrURLNotFound)
	}

	return fullURL, nil
}

func (s *inMemoryStore) GetURLs(_ context.Context, userID string) (map[string]string, error) {
	userURLs, ok := s.m[userID]
	if !ok {
		return nil, fmt.Errorf("%w", ErrUserHasNoURLs)
	}

	return userURLs, nil
}

func (s *inMemoryStore) SaveURLsBatch(_ context.Context,
	urls map[string]string, userID string) (map[string]string, error) {
	userURLs, ok := s.m[userID]
	if !ok {
		return nil, fmt.Errorf("%w", ErrUserHasNoURLs)
	}

	res := make(map[string]string)

	for fullURL, shortURL := range urls {
		userURLs[shortURL] = fullURL
		res[fullURL] = shortURL
	}

	return res, nil
}

func (s *inMemoryStore) Ping(_ context.Context) error {
	return nil
}

func (s *inMemoryStore) Close() {}
