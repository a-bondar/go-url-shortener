package store

import (
	"context"

	"go.uber.org/zap"
)

type Config struct {
	DatabaseDSN     string
	FileStoragePath string
}

type Store interface {
	SaveURL(ctx context.Context, fullURL string, shortURL string) (string, error)
	GetURL(ctx context.Context, shortURL string) (string, error)
	SaveURLsBatch(ctx context.Context, urls map[string]string) (map[string]string, error)
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
