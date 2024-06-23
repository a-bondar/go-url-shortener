package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/a-bondar/go-url-shortener/internal/app/config"
	"github.com/a-bondar/go-url-shortener/internal/app/handlers"
	"github.com/a-bondar/go-url-shortener/internal/app/logger"
	"github.com/a-bondar/go-url-shortener/internal/app/router"
	"github.com/a-bondar/go-url-shortener/internal/app/service"
	"github.com/a-bondar/go-url-shortener/internal/app/store"
	"go.uber.org/zap"
)

func main() {
	if err := Run(); err != nil {
		log.Fatal(err)
	}
}

func Run() error {
	l, err := logger.NewLogger()

	if err != nil {
		return fmt.Errorf("failed to create logger: %w", err)
	}

	defer func(l *zap.Logger) {
		err := l.Sync()
		if err != nil {
			l.Error("Failed to sync logger", zap.Error(err))
		}
	}(l)

	cfg := config.NewConfig()

	s, err := store.NewStore(context.Background(), store.Config{
		DatabaseDSN:     cfg.DatabaseDSN,
		FileStoragePath: cfg.FileStoragePath,
	})
	if err != nil {
		return fmt.Errorf("failed to initialize store: %w", err)
	}

	defer s.Close()

	svc := service.NewService(s, cfg)
	h := handlers.NewHandler(svc, l)

	l.Info("Running server", zap.String("address", cfg.RunAddr))

	if err := http.ListenAndServe(cfg.RunAddr, router.Router(h, l)); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			l.Error("HTTP server has encountered an error", zap.Error(err))

			return fmt.Errorf("HTTP server has encountered an error: %w", err)
		}
	}

	return nil
}
