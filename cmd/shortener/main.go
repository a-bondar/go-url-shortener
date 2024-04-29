package main

import (
	"github.com/a-bondar/go-url-shortener/internal/app/router"
	"net/http"
)

func main() {
	err := http.ListenAndServe(`localhost:8080`, router.Router())
	if err != nil {
		panic(err)
	}
}
