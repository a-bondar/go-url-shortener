package router

import (
	"github.com/a-bondar/go-url-shortener/internal/app/handlers"
	"github.com/go-chi/chi/v5"
)

func Router() chi.Router {
	r := chi.NewRouter()

	r.Post("/", handlers.HandlePost)
	r.Get("/{linkID}", handlers.HandleGet)

	return r
}
