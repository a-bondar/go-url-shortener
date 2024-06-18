package service

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/a-bondar/go-url-shortener/internal/app/models"
)

type Store interface {
	SaveURL(fullURL string, shortURL string) error
	GetURL(shortURL string) (string, error)
	SaveBatchURLs(urls map[string]string) (map[string]string, error)
	Ping() error
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

func (s *Service) shortenURL() (string, error) {
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

	return shortenURL, nil
}

func (s *Service) SaveURL(fullURL string) (string, error) {
	shortenURL, err := s.shortenURL()
	if err != nil {
		return "", fmt.Errorf("failed to generate unique short URL: %w", err)
	}

	if err := s.s.SaveURL(fullURL, shortenURL); err != nil {
		return "", fmt.Errorf("failed to save URL: %w", err)
	}

	return shortenURL, nil
}

func (s *Service) SaveBatchURLs(urls []models.OriginalURLCorrelation) ([]models.ShortURLCorrelation, error) {
	// Мапа для связи корреляционных идентификаторов и полных URL
	fullURLbyCorrID := make(map[string]string)
	urlsMap := make(map[string]string)

	for _, url := range urls {
		shortURL, err := s.shortenURL()

		if err != nil {
			return nil, fmt.Errorf("failed to generate unique short URL: %w", err)
		}

		urlsMap[url.OriginalURL] = shortURL
		fullURLbyCorrID[url.OriginalURL] = url.CorrelationID
	}

	batchRes, err := s.s.SaveBatchURLs(urlsMap)
	if err != nil {
		return nil, fmt.Errorf("failed to save batch URLs: %w", err)
	}

	resp := make([]models.ShortURLCorrelation, 0, len(batchRes))
	for fullURL, shortURL := range batchRes {
		resp = append(resp, models.ShortURLCorrelation{
			CorrelationID: fullURLbyCorrID[fullURL],
			ShortURL:      shortURL,
		})
	}

	return resp, nil
}

func (s *Service) GetURL(shortURL string) (string, error) {
	fullURL, err := s.s.GetURL(shortURL)

	if err != nil {
		return "", fmt.Errorf("failed to get full URL: %w", err)
	}

	return fullURL, nil
}

func (s *Service) Ping() error {
	err := s.s.Ping()

	if err != nil {
		return fmt.Errorf("failed to reach store: %w", err)
	}

	return nil
}
