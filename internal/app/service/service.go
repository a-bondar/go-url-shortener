package service

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"net/url"
	"time"

	"github.com/a-bondar/go-url-shortener/internal/app/config"
	"github.com/a-bondar/go-url-shortener/internal/app/models"
	"github.com/a-bondar/go-url-shortener/internal/app/store"
)

type Store interface {
	SaveURL(ctx context.Context, fullURL string, shortURL string) error
	GetURL(ctx context.Context, shortURL string) (string, error)
	SaveURLsBatch(ctx context.Context, urls map[string]string) (map[string]string, error)
	Ping(ctx context.Context) error
}

type Service struct {
	s   Store
	cfg *config.Config
}

func NewService(s Store, cfg *config.Config) *Service {
	return &Service{s: s, cfg: cfg}
}

const maxRetries = 3
const maxShortURLLength = 8
const (
	chars                 = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	failedToBuildURLError = "failed to build URL: %w"
)

func generateRandomString(size int) string {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	b := make([]rune, size)
	for i := range b {
		b[i] = rune(chars[rnd.Intn(len(chars))])
	}

	return string(b)
}

func (s *Service) shortenURL(ctx context.Context) (string, error) {
	var shortenURL string

	for range maxRetries {
		shortenURL = generateRandomString(maxShortURLLength)

		if _, err := s.s.GetURL(ctx, shortenURL); err != nil {
			break
		}
	}

	if shortenURL == "" {
		return "", errors.New("failed to generate unique short URL")
	}

	return shortenURL, nil
}

func (s *Service) buildURL(shortenURL string) (string, error) {
	res, err := url.JoinPath(s.cfg.ShortLinkBaseURL, shortenURL)
	if err != nil {
		return "", fmt.Errorf(failedToBuildURLError, err)
	}

	return res, nil
}

func (s *Service) SaveURL(ctx context.Context, fullURL string) (string, error) {
	shortenURL, err := s.shortenURL(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to generate unique short URL: %w", err)
	}

	err = s.s.SaveURL(ctx, fullURL, shortenURL)
	if err != nil {
		if !errors.Is(err, store.ErrConflict) {
			return "", fmt.Errorf("failed to save URL: %w", err)
		}

		var conflictErr *store.URLConflictError
		if errors.As(err, &conflictErr) {
			resURL, buildErr := s.buildURL(conflictErr.URL)
			if buildErr != nil {
				return "", fmt.Errorf(failedToBuildURLError, buildErr)
			}

			return "", fmt.Errorf("%w", store.NewURLConflictError(resURL, store.ErrConflict))
		}
	}

	resURL, err := s.buildURL(shortenURL)
	if err != nil {
		return "", fmt.Errorf(failedToBuildURLError, err)
	}

	return resURL, nil
}

func (s *Service) SaveBatchURLs(
	ctx context.Context,
	urls []models.OriginalURLCorrelation) ([]models.ShortURLCorrelation, error) {
	// Мапа для связи корреляционных идентификаторов и полных URL
	fullURLbyCorrID := make(map[string]string)
	urlsMap := make(map[string]string)

	for _, URL := range urls {
		shortURL, err := s.shortenURL(ctx)

		if err != nil {
			return nil, fmt.Errorf("failed to generate unique short URL: %w", err)
		}

		urlsMap[URL.OriginalURL] = shortURL
		fullURLbyCorrID[URL.OriginalURL] = URL.CorrelationID
	}

	batchRes, err := s.s.SaveURLsBatch(ctx, urlsMap)
	if err != nil {
		return nil, fmt.Errorf("failed to save batch URLs: %w", err)
	}

	resp := make([]models.ShortURLCorrelation, 0, len(batchRes))
	for fullURL, shortURL := range batchRes {
		resURL, err := s.buildURL(shortURL)
		if err != nil {
			return nil, fmt.Errorf(failedToBuildURLError, err)
		}

		resp = append(resp, models.ShortURLCorrelation{
			CorrelationID: fullURLbyCorrID[fullURL],
			ShortURL:      resURL,
		})
	}

	return resp, nil
}

func (s *Service) GetURL(ctx context.Context, shortURL string) (string, error) {
	fullURL, err := s.s.GetURL(ctx, shortURL)

	if err != nil {
		return "", fmt.Errorf("failed to get full URL: %w", err)
	}

	return fullURL, nil
}

func (s *Service) Ping(ctx context.Context) error {
	err := s.s.Ping(ctx)
	if err != nil {
		return fmt.Errorf("failed to reach store: %w", err)
	}

	return nil
}
