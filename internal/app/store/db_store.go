package store

import (
	"context"
	"embed"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type DBStore struct {
	logger *zap.Logger
	pool   *pgxpool.Pool
}

func newDBStore(ctx context.Context, dsn string, logger *zap.Logger) (*DBStore, error) {
	if err := runMigrations(dsn); err != nil {
		return nil, fmt.Errorf("failed to run DB migrations: %w", err)
	}

	pool, err := initPool(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}

	return &DBStore{pool: pool, logger: logger}, nil
}

//go:embed migrations/*.sql
var migrationsDir embed.FS

func runMigrations(dsn string) error {
	d, err := iofs.New(migrationsDir, "migrations")
	if err != nil {
		return fmt.Errorf("failed to return an iofs driver: %w", err)
	}

	m, err := migrate.NewWithSourceInstance("iofs", d, dsn)
	if err != nil {
		return fmt.Errorf("failed to get a new migrate instance: %w", err)
	}
	if err := m.Up(); err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			return fmt.Errorf("failed to apply migrations to the DB: %w", err)
		}
	}
	return nil
}

func initPool(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	poolCfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse DSN: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initalize a connection pool: %w", err)
	}

	return pool, nil
}

func (s *DBStore) SaveURL(ctx context.Context, fullURL string, shortURL string) error {
	_, err := s.pool.Exec(ctx, "INSERT INTO short_links (original_url, short_url) VALUES ($1, $2)", fullURL, shortURL)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) &&
			pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) && pgErr.ConstraintName == "short_links_original_url_key" {
			var existingShortURL string
			selectQuery := `
				SELECT short_url FROM short_links WHERE original_url = $1
			`
			err = s.pool.QueryRow(ctx, selectQuery, fullURL).Scan(&existingShortURL)
			if err != nil {
				return fmt.Errorf("failed to get existing short_url: %w", err)
			}

			return NewURLConflictError(existingShortURL, ErrConflict)
		}

		return fmt.Errorf("failed to save URL: %w", err)
	}

	return nil
}

func (s *DBStore) SaveURLsBatch(ctx context.Context, urls map[string]string) (map[string]string, error) {
	query := `INSERT INTO short_links (original_url, short_url) VALUES ($1, $2)`
	batch := &pgx.Batch{}
	for fullURL, shortURL := range urls {
		batch.Queue(query, fullURL, shortURL)
	}

	results := s.pool.SendBatch(ctx, batch)
	defer func() {
		err := results.Close()
		if err != nil {
			s.logger.Error("Failed to close batch", zap.Error(err))
		}
	}()

	res := make(map[string]string)
	for fullURL, shortURL := range urls {
		_, err := results.Exec()
		if err != nil {
			return nil, fmt.Errorf("unable to insert row: %w", err)
		}

		res[fullURL] = shortURL
	}

	return res, nil
}

func (s *DBStore) GetURL(ctx context.Context, shortURL string) (string, error) {
	var originalURL string

	err := s.pool.QueryRow(ctx, "SELECT original_url FROM short_links WHERE short_url = $1", shortURL).Scan(&originalURL)
	if err != nil {
		return "", fmt.Errorf("failed to get full URL: %w", err)
	}

	return originalURL, nil
}

func (s *DBStore) Ping(ctx context.Context) error {
	err := s.pool.Ping(ctx)
	if err != nil {
		return fmt.Errorf("failed to ping DB: %w", err)
	}

	return nil
}

func (s *DBStore) Close() {
	s.pool.Close()
}
