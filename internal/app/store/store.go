package store

import (
	"context"
	"errors"

	"go.uber.org/zap"
)

var (
	ErrUserHasNoURLs = errors.New("user has no URLs")
	ErrURLNotFound   = errors.New("URL not found for the given short URL")
)

type Config struct {
	DatabaseDSN     string
	FileStoragePath string
}

type Store interface {
	SaveURL(ctx context.Context, fullURL string, shortURL string, userID string) (string, error)
	GetURL(ctx context.Context, shortURL string) (string, error)
	GetURLs(ctx context.Context, userID string) (map[string]string, error)
	SaveURLsBatch(ctx context.Context, urls map[string]string, userID string) (map[string]string, error)
	Ping(ctx context.Context) error
	Close()
}

func NewStore(ctx context.Context, cfg Config, logger *zap.Logger) (Store, error) {
	if cfg.DatabaseDSN != "" {
		return newDBStore(ctx, cfg.DatabaseDSN, logger)
	}

	if cfg.FileStoragePath != "" {
		return newFileStore(ctx, cfg.FileStoragePath)
	}

	return newInMemoryStore(), nil
}
