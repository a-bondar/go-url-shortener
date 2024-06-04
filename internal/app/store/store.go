package store

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/a-bondar/go-url-shortener/internal/app/models"
	"github.com/google/uuid"
)

type Store struct {
	m map[string]string
}

func NewStore() *Store {
	return &Store{
		m: make(map[string]string),
	}
}

func (s *Store) SaveURL(fullURL string, shortURL string) error {
	s.m[shortURL] = fullURL

	return nil
}

func (s *Store) GetURL(shortURL string) (string, error) {
	if URL, ok := s.m[shortURL]; ok {
		return URL, nil
	}

	return "", errors.New("URL not found")
}

func (s *Store) WriteToFile(shortURL string, fullURL string, fName string) error {
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

	file, err := os.OpenFile(fName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o600)
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

func (s *Store) LoadFromFile(fName string) error {
	file, err := os.Open(fName)
	if err != nil {
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

		if err := s.SaveURL(data.OriginalURL, data.ShortURL); err != nil {
			return fmt.Errorf("failed to save URL: %w", err)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("failed to scan file: %w", err)
	}

	return nil
}
