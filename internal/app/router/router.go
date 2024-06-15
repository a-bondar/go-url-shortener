package router

import (
	"github.com/a-bondar/go-url-shortener/internal/app/handlers"
	"github.com/a-bondar/go-url-shortener/internal/app/middleware"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func Router(h *handlers.Handler, logger *zap.Logger) chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.WithLogging(logger))
	r.Use(middleware.WithGzip(logger))
	r.Post("/", h.HandlePost)
	r.Get("/{linkID}", h.HandleGet)
	r.Get("/ping", h.HandleDatabasePing)
	r.Post("/api/shorten", h.HandleShorten)
	r.Post("/api/shorten/batch", h.HandleShortenBatch)

	return r
}
