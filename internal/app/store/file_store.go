package store

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/a-bondar/go-url-shortener/internal/app/models"
	"github.com/google/uuid"
)

type fileStore struct {
	inMemoryStore *inMemoryStore
	fName         string
}

func newFileStore(ctx context.Context, fName string) (*fileStore, error) {
	store := &fileStore{
		inMemoryStore: newInMemoryStore(),
		fName:         fName,
	}

	err := store.loadFromFile(ctx)
	if err != nil {
		return nil, err
	}

	return store, nil
}

func (s *fileStore) SaveURL(ctx context.Context, fullURL string, shortURL string) error {
	err := s.inMemoryStore.SaveURL(ctx, fullURL, shortURL)

	if err != nil {
		return err
	}

	return s.writeToFile(fullURL, shortURL)
}

func (s *fileStore) SaveURLsBatch(ctx context.Context, urls map[string]string) (map[string]string, error) {
	res := make(map[string]string)

	for fullURL, shortURL := range urls {
		err := s.inMemoryStore.SaveURL(ctx, fullURL, shortURL)

		if err != nil {
			return nil, err
		}

		err = s.writeToFile(fullURL, shortURL)
		if err != nil {
			return nil, err
		}

		res[fullURL] = shortURL
	}

	return res, nil
}

func (s *fileStore) GetURL(ctx context.Context, shortURL string) (string, error) {
	return s.inMemoryStore.GetURL(ctx, shortURL)
}

func (s *fileStore) writeToFile(fullURL string, shortURL string) error {
	data := models.Data{
		UUID:        uuid.NewString(),
		ShortURL:    shortURL,
		OriginalURL: fullURL,
	}

	dataToJSON, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	dataToJSON = append(dataToJSON, '\n')

	const fileModeOwnerReadWrite = 0o600
	file, err := os.OpenFile(s.fName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, fileModeOwnerReadWrite)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println("Error closing file:", err)
		}
	}(file)

	if _, err := file.Write(dataToJSON); err != nil {
		return fmt.Errorf("failed to write data to file: %w", err)
	}

	return nil
}

func (s *fileStore) loadFromFile(ctx context.Context) error {
	file, err := os.Open(s.fName)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}

		return fmt.Errorf("failed to open file: %w", err)
	}

	defer func() {
		err = file.Close()
		if err != nil {
			fmt.Println("Error closing file:", err)
		}
	}()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var data models.Data

		if err := json.Unmarshal(scanner.Bytes(), &data); err != nil {
			return fmt.Errorf("failed to unmarshal data: %w", err)
		}

		if err := s.inMemoryStore.SaveURL(ctx, data.OriginalURL, data.ShortURL); err != nil {
			return fmt.Errorf("failed to save URL: %w", err)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("failed to scan file: %w", err)
	}

	return nil
}

func (s *fileStore) Ping(_ context.Context) error {
	return nil
}

func (s *fileStore) Close() {}
