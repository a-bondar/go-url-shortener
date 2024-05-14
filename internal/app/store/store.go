package store

import "errors"

type Store struct {
	m map[string]string
}

func NewStore() *Store {
	return &Store{
		m: make(map[string]string),
	}
}

func (s *Store) SaveURL(fullURL string, shortURL string) error {
	s.m[shortURL] = fullURL

	return nil
}

func (s *Store) GetURL(shortURL string) (string, error) {
	if URL, ok := s.m[shortURL]; ok {
		return URL, nil
	}

	return "", errors.New("URL not found")
}
