package service

import (
	"errors"
	"fmt"
)

type Store interface {
	SaveURL(fullURL string, shortURL string) error
	GetURL(shortURL string) (string, error)
}

type Utils interface {
	GenerateRandomString(size int) string
}

type Service struct {
	s Store
	u Utils
}

func NewService(s Store, u Utils) *Service {
	return &Service{s: s, u: u}
}

const maxRetries = 3
const maxShortURLLength = 8

func (s *Service) SaveURL(fullURL string) (string, error) {
	var shortenURL string

	for range maxRetries {
		shortenURL = s.u.GenerateRandomString(maxShortURLLength)

		if _, err := s.s.GetURL(shortenURL); err != nil {
			break
		} else {
			shortenURL = ""
		}
	}

	if shortenURL == "" {
		return "", errors.New("failed to generate unique short URL")
	}

	if err := s.s.SaveURL(fullURL, shortenURL); err != nil {
		return "", fmt.Errorf("failed to save URL: %w", err)
	}

	return shortenURL, nil
}

func (s *Service) GetURL(shortURL string) (string, error) {
	fullURL, err := s.s.GetURL(shortURL)

	if err != nil {
		return "", fmt.Errorf("failed to get full URL: %w", err)
	}

	return fullURL, nil
}
