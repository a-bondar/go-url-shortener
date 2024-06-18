package store

import (
	"database/sql"
	"fmt"
)

type DBStore struct {
	db *sql.DB
}

func newDBStore(dsn string) (*DBStore, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open DB: %w", err)
	}

	store := &DBStore{db: db}
	err = store.createTable()

	if err != nil {
		return nil, err
	}

	return store, nil
}

func (s *DBStore) SaveURL(fullURL string, shortURL string) error {
	_, err := s.db.Exec("INSERT INTO short_links (original_url, short_url) VALUES ($1, $2)", fullURL, shortURL)

	if err != nil {
		return fmt.Errorf("failed to save URL: %w", err)
	}

	return nil
}

func (s *DBStore) SaveURLsBatch(urls map[string]string) (map[string]string, error) {
	tx, err := s.db.Begin()

	if err != nil {
		return nil, fmt.Errorf("failed to start transaction: %w", err)
	}

	stmt, err := tx.Prepare("INSERT INTO short_links (original_url, short_url) VALUES ($1, $2)")
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return nil, fmt.Errorf("failed to rollback transaction: %w", rollbackErr)
		}

		return nil, fmt.Errorf("failed to prepare statement: %w", err)
	}

	defer func() {
		if err = stmt.Close(); err != nil {
			fmt.Printf("failed to close statement: %v\n", err)
		}
	}()

	res := make(map[string]string)

	for fullURL, shortURL := range urls {
		_, err = stmt.Exec(fullURL, shortURL)

		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				return nil, fmt.Errorf("failed to rollback transaction: %w", rollbackErr)
			}

			return nil, fmt.Errorf("failed to execute statement: %w", err)
		}

		res[fullURL] = shortURL
	}

	err = tx.Commit()

	if err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return res, nil
}

func (s *DBStore) GetURL(shortURL string) (string, error) {
	var originalURL string

	err := s.db.QueryRow("SELECT original_url FROM short_links WHERE short_url = $1", shortURL).Scan(&originalURL)

	if err != nil {
		return "", fmt.Errorf("failed to get full URL: %w", err)
	}

	return originalURL, nil
}

func (s *DBStore) createTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS short_links (
		id SERIAL PRIMARY KEY,
		short_url TEXT NOT NULL UNIQUE,
		original_url TEXT NOT NULL
	);`
	_, err := s.db.Exec(query)

	if err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}

	return nil
}

func (s *DBStore) Ping() error {
	err := s.db.Ping()

	if err != nil {
		return fmt.Errorf("failed to ping DB: %w", err)
	}

	return nil
}

func (s *DBStore) Close() error {
	err := s.db.Close()

	if err != nil {
		return fmt.Errorf("failed to close DB: %w", err)
	}

	return nil
}
