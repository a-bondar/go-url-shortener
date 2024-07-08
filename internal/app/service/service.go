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
)

type Store interface {
	SaveURL(ctx context.Context, fullURL string, shortURL string, userID string) (string, error)
	GetURL(ctx context.Context, shortURL string) (string, bool, error)
	GetURLs(ctx context.Context, userID string) (map[string]string, error)
	DeleteURLs(ctx context.Context, urls []string, userID string) error
	SaveURLsBatch(ctx context.Context, urls map[string]string, userID string) (map[string]string, error)
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

var ErrConflict = errors.New("data conflict")

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

		if _, _, err := s.s.GetURL(ctx, shortenURL); err != nil {
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

func (s *Service) SaveURL(ctx context.Context, fullURL string, userID string) (string, error) {
	shortenURL, err := s.shortenURL(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to generate unique short URL: %w", err)
	}

	resultedShortURL, err := s.s.SaveURL(ctx, fullURL, shortenURL, userID)
	if err != nil {
		return "", fmt.Errorf("failed to save URL: %w", err)
	}

	resURL, err := s.buildURL(resultedShortURL)
	if err != nil {
		return "", fmt.Errorf(failedToBuildURLError, err)
	}

	if shortenURL != resultedShortURL {
		err = ErrConflict
	}

	return resURL, err
}

func (s *Service) SaveBatchURLs(
	ctx context.Context,
	urls []models.OriginalURLCorrelation,
	userID string,
) ([]models.ShortURLCorrelation, error) {
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

	batchRes, err := s.s.SaveURLsBatch(ctx, urlsMap, userID)
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

func (s *Service) GetURL(ctx context.Context, shortURL string) (string, bool, error) {
	fullURL, deleted, err := s.s.GetURL(ctx, shortURL)
	if err != nil {
		return "", false, fmt.Errorf("failed to get full URL: %w", err)
	}

	return fullURL, deleted, nil
}

func (s *Service) GetURLs(ctx context.Context, userID string) ([]models.URLsPair, error) {
	userURLs, err := s.s.GetURLs(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	resp := make([]models.URLsPair, 0, len(userURLs))
	for shortURL, fullURL := range userURLs {
		resURL, buildErr := s.buildURL(shortURL)
		if buildErr != nil {
			return nil, fmt.Errorf(failedToBuildURLError, buildErr)
		}

		resp = append(resp, models.URLsPair{
			ShortURL:    resURL,
			OriginalURL: fullURL,
		})
	}

	return resp, nil
}

func (s *Service) DeleteURLs(ctx context.Context, urls []string, userID string) error {
	err := s.s.DeleteURLs(ctx, urls, userID)
	if err != nil {
		return fmt.Errorf("failed to delete urls: %w", err)
	}

	return nil
}

func (s *Service) Ping(ctx context.Context) error {
	err := s.s.Ping(ctx)
	if err != nil {
		return fmt.Errorf("failed to reach store: %w", err)
	}

	return nil
}
