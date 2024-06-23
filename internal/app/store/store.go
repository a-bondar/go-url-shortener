package store

import (
	"context"
	"errors"
	"fmt"
)

type Config struct {
	DatabaseDSN     string
	FileStoragePath string
}

type Store interface {
	SaveURL(ctx context.Context, fullURL string, shortURL string) error
	GetURL(ctx context.Context, shortURL string) (string, error)
	SaveURLsBatch(ctx context.Context, urls map[string]string) (map[string]string, error)
	Ping(ctx context.Context) error
	Close()
}

var ErrConflict = errors.New("data conflict")

type URLConflictError struct {
	Err error
	URL string
}

func (uce *URLConflictError) Error() string {
	return fmt.Sprintf("[%s] %v", uce.URL, uce.Err)
}

func (uce *URLConflictError) Unwrap() error {
	return uce.Err
}

func NewURLConflictError(shortURL string, err error) error {
	return &URLConflictError{
		URL: shortURL,
		Err: err,
	}
}

func NewStore(ctx context.Context, cfg Config) (Store, error) {
	if cfg.DatabaseDSN != "" {
		return newDBStore(ctx, cfg.DatabaseDSN)
	}

	if cfg.FileStoragePath != "" {
		return newFileStore(ctx, cfg.FileStoragePath)
	}

	return newInMemoryStore(), nil
}
