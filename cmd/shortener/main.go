package main

import (
	"github.com/a-bondar/go-url-shortener/internal/app/handlers"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc(`/`, handlers.HandleRoot)

	err := http.ListenAndServe(`localhost:8080`, mux)
	if err != nil {
		panic(err)
	}
}
