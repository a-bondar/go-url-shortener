package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/a-bondar/go-url-shortener/internal/app/config"
	"github.com/a-bondar/go-url-shortener/internal/app/handlers"
	"github.com/a-bondar/go-url-shortener/internal/app/router"
	"github.com/a-bondar/go-url-shortener/internal/app/service"
	"github.com/a-bondar/go-url-shortener/internal/app/store"
)

func main() {
	if err := Run(); err != nil {
		log.Fatal(err)
	}
}

func Run() error {
	cfg := config.GetConfig()
	fmt.Println("Running server on:", cfg.RunAddr)

	s := store.NewStore()
	svc := service.NewService(s)
	h := handlers.NewHanlder(cfg, svc)

	if err := http.ListenAndServe(cfg.RunAddr, router.Router(h)); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("HTTP server has encountered an error: %w", err)
		}
	}

	return nil
}
