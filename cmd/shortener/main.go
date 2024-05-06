package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/a-bondar/go-url-shortener/internal/app/config"
	"github.com/a-bondar/go-url-shortener/internal/app/router"
)

func main() {
	if err := Run(); err != nil {
		log.Fatal(err)
	}
}

func Run() error {
	cfg := config.GetConfig()

	fmt.Println("Running server on:", cfg.RunAddr)

	if err := http.ListenAndServe(cfg.RunAddr, router.Router(cfg)); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("HTTP server has encountered an error: %w", err)
		}
	}

	return nil
}
