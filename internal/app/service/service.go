package service

import (
	"errors"
	"fmt"
	"math/rand"
	"time"
)

type Store interface {
	SaveURL(fullURL string, shortURL string) error
	GetURL(shortURL string) (string, error)
	WriteToFile(shortURL string, fullURL string, fName string) error
}

type Service struct {
	s Store
}

func NewService(s Store) *Service {
	return &Service{s: s}
}

const maxRetries = 3
const maxShortURLLength = 8
const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

func generateRandomString(size int) string {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	b := make([]rune, size)
	for i := range b {
		b[i] = rune(chars[rnd.Intn(len(chars))])
	}

	return string(b)
}

func (s *Service) SaveURL(fullURL string, fName string) (string, error) {
	var shortenURL string

	for range maxRetries {
		shortenURL = generateRandomString(maxShortURLLength)

		if _, err := s.s.GetURL(shortenURL); err != nil {
			break
		}
	}

	if shortenURL == "" {
		return "", errors.New("failed to generate unique short URL")
	}

	if err := s.s.SaveURL(fullURL, shortenURL); err != nil {
		return "", fmt.Errorf("failed to save URL: %w", err)
	}

	if fName != "" {
		if err := s.s.WriteToFile(shortenURL, fullURL, fName); err != nil {
			return "", fmt.Errorf("failed to write to file: %w", err)
		}
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
