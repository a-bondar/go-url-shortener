package store

import (
	"context"
	"fmt"
)

type ShortURLData struct {
	fullURL string
	deleted bool
}

type inMemoryStore struct {
	m map[string]map[string]ShortURLData
}

func newInMemoryStore() *inMemoryStore {
	return &inMemoryStore{
		m: make(map[string]map[string]ShortURLData),
	}
}

func (s *inMemoryStore) SaveURL(_ context.Context, fullURL string, shortURL string, userID string) (string, error) {
	userURLs, ok := s.m[userID]
	if !ok {
		// У пользователя еще нет данных, создаем пустой map
		userURLs = make(map[string]ShortURLData)
		s.m[userID] = userURLs
	}

	// Проверка на конфликт
	for currentShortURL, currentShortURLData := range userURLs {
		if currentShortURLData.fullURL == fullURL {
			return currentShortURL, nil
		}
	}

	userURLs[shortURL] = ShortURLData{fullURL: fullURL, deleted: false}

	return shortURL, nil
}

func (s *inMemoryStore) GetURL(_ context.Context, shortURL string) (string, bool, error) {
	for _, userURLs := range s.m {
		if shortURLData, ok := userURLs[shortURL]; ok {
			return shortURLData.fullURL, shortURLData.deleted, nil
		}
	}

	return "", false, fmt.Errorf("%w", ErrURLNotFound)
}

func (s *inMemoryStore) GetURLs(_ context.Context, userID string) (map[string]string, error) {
	userURLs, ok := s.m[userID]
	if !ok {
		return nil, fmt.Errorf("%w", ErrUserHasNoURLs)
	}

	res := make(map[string]string)
	for shortURL, shortURLData := range userURLs {
		res[shortURL] = shortURLData.fullURL
	}

	return res, nil
}

func (s *inMemoryStore) SaveURLsBatch(_ context.Context,
	urls map[string]string, userID string) (map[string]string, error) {
	userURLs, ok := s.m[userID]
	if !ok {
		// У пользователя еще нет данных, создаем пустой map
		userURLs = make(map[string]ShortURLData)
		s.m[userID] = userURLs
	}

	res := make(map[string]string)
	for fullURL, shortURL := range urls {
		userURLs[shortURL] = ShortURLData{fullURL: fullURL, deleted: false}
		res[fullURL] = shortURL
	}

	return res, nil
}

func (s *inMemoryStore) Ping(_ context.Context) error {
	return nil
}

func (s *inMemoryStore) Close() {}
