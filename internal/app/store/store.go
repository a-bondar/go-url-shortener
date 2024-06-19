package store

import (
	"errors"
	"fmt"
)

type Store interface {
	SaveURL(fullURL string, shortURL string) error
	GetURL(shortURL string) (string, error)
	SaveURLsBatch(urls map[string]string) (map[string]string, error)
	Ping() error
	Close() error
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

func NewStore(dsn string, fName string) (Store, error) {
	if dsn != "" {
		return newDBStore(dsn)
	}

	if fName != "" {
		return newFileStore(fName)
	}

	return newInMemoryStore(), nil
}
