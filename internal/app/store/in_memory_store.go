package store

import (
	"context"
	"fmt"
)

type ShortURLData struct {
	FullURL string
	Deleted bool
}

type inMemoryStore struct {
	m map[string]map[string]*ShortURLData
}

func newInMemoryStore() *inMemoryStore {
	return &inMemoryStore{
		m: make(map[string]map[string]*ShortURLData),
	}
}

func (s *inMemoryStore) SaveURL(_ context.Context, fullURL string, shortURL string, userID string) (string, error) {
	userURLs, ok := s.m[userID]
	if !ok {
		// У пользователя еще нет данных, создаем пустой map
		userURLs = make(map[string]*ShortURLData)
		s.m[userID] = userURLs
	}

	// Проверка на конфликт с учетом флага deleted
	for currentShortURL, currentShortURLData := range userURLs {
		if currentShortURLData.FullURL == fullURL {
			if currentShortURLData.Deleted {
				userURLs[shortURL] = &ShortURLData{FullURL: fullURL, Deleted: false}
				return shortURL, nil
			} else {
				return currentShortURL, nil
			}
		}
	}

	userURLs[shortURL] = &ShortURLData{FullURL: fullURL, Deleted: false}

	return shortURL, nil
}

func (s *inMemoryStore) GetURL(_ context.Context, shortURL string) (string, bool, error) {
	for _, userURLs := range s.m {
		if shortURLData, ok := userURLs[shortURL]; ok {
			return shortURLData.FullURL, shortURLData.Deleted, nil
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
		res[shortURL] = shortURLData.FullURL
	}

	return res, nil
}

func (s *inMemoryStore) DeleteURLs(_ context.Context, urls []string, userID string) error {
	userURLs, ok := s.m[userID]
	if !ok {
		return fmt.Errorf("%w", ErrUserHasNoURLs)
	}

	for _, url := range urls {
		if userURL, ok := userURLs[url]; ok {
			userURL.Deleted = true
		}
	}

	return nil
}

func (s *inMemoryStore) SaveURLsBatch(_ context.Context,
	urls map[string]string, userID string) (map[string]string, error) {
	userURLs, ok := s.m[userID]
	if !ok {
		// У пользователя еще нет данных, создаем пустой map
		userURLs = make(map[string]*ShortURLData)
		s.m[userID] = userURLs
	}

	res := make(map[string]string)
	for fullURL, shortURL := range urls {
		userURLs[shortURL] = &ShortURLData{FullURL: fullURL, Deleted: false}
		res[fullURL] = shortURL
	}

	return res, nil
}

func (s *inMemoryStore) Ping(_ context.Context) error {
	return nil
}

func (s *inMemoryStore) Close() {}
