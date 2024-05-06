package router

import (
	"net/http"

	"github.com/a-bondar/go-url-shortener/internal/app/config"
	"github.com/a-bondar/go-url-shortener/internal/app/handlers"
	"github.com/go-chi/chi/v5"
)

func Router(cfg *config.Config) chi.Router {
	r := chi.NewRouter()

	r.Post("/", func(writer http.ResponseWriter, request *http.Request) {
		handlers.HandlePost(cfg, writer, request)
	})
	r.Get("/{linkID}", handlers.HandleGet)

	return r
}
