package main

import (
	"io"
	"net/http"
	"net/url"
	"strings"
)

var linksMap map[string]string

func isValidURL(input string) bool {
	_, err := url.ParseRequestURI(input)
	if err != nil {
		return false
	}
	return true
}

func handlePost(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if r.Header.Get("Content-Type") != "text/plain" {
		http.Error(w, "Header is not correct", http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !isValidURL(string(body)) {
		http.Error(w, "Data is not valid", http.StatusBadRequest)
		return
	}

	// @TODO - add base64 implementation

	w.Write([]byte(body))
}

func handleGet(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	linkId := strings.Split(r.URL.Path[1:], "/")[0]

	if linkId == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	link, ok := linksMap[linkId]

	if !ok {
		http.Error(w, `Link not found`, http.StatusNotFound)
		return
	}

	http.Redirect(w, r, link, http.StatusTemporaryRedirect)
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		handlePost(w, r)
		return
	}

	if r.Method == http.MethodGet {
		handleGet(w, r)
		return
	}

	w.WriteHeader(http.StatusMethodNotAllowed)
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc(`/`, handleRoot)

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}
