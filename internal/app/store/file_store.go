package store

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	"github.com/a-bondar/go-url-shortener/internal/app/models"
	"github.com/google/uuid"
)

type FileStore struct {
	InMemoryStore *InMemoryStore
	fName         string
}

func NewFileStore(fName string) (*FileStore, error) {
	store := &FileStore{
		InMemoryStore: NewInMemoryStore(),
		fName:         fName,
	}

	err := store.loadFromFile()
	if err != nil {
		return nil, err
	}

	return store, nil
}

func (s *FileStore) SaveURL(fullURL string, shortURL string) error {
	err := s.InMemoryStore.SaveURL(fullURL, shortURL)

	if err != nil {
		return err
	}

	return s.writeToFile(fullURL, shortURL)
}

func (s *FileStore) GetURL(shortURL string) (string, error) {
	return s.InMemoryStore.GetURL(shortURL)
}

func (s *FileStore) writeToFile(fullURL string, shortURL string) error {
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

func (s *FileStore) loadFromFile() error {
	file, err := os.Open(s.fName)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}

		return fmt.Errorf("failed to open file: %w", err)
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println("Error closing file:", err)
		}
	}(file)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var data models.Data

		if err := json.Unmarshal(scanner.Bytes(), &data); err != nil {
			return fmt.Errorf("failed to unmarshal data: %w", err)
		}

		if err := s.InMemoryStore.SaveURL(data.OriginalURL, data.ShortURL); err != nil {
			return fmt.Errorf("failed to save URL: %w", err)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("failed to scan file: %w", err)
	}

	return nil
}
