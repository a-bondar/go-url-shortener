package main

import (
	"fmt"
	"github.com/a-bondar/go-url-shortener/internal/app/config"
	"github.com/a-bondar/go-url-shortener/internal/app/router"
	"net/http"
)

func main() {
	config.ParseFlags()

	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	fmt.Println("Running server on:", config.FlagOptions.RunAddr)

	return http.ListenAndServe(config.FlagOptions.RunAddr, router.Router())
}
