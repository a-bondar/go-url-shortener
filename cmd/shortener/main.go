package main

import (
	"encoding/base64"
	"io"
	"net/http"
	"strings"
)

var linksMap = map[string]string{}

func handlePost(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" || r.Header.Get("Content-Type") != "text/plain" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sEnc := base64.StdEncoding.EncodeToString(body)
	linksMap[sEnc] = string(body)
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("http" + "://" + r.Host + r.URL.String() + sEnc))
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

	err := http.ListenAndServe(`localhost:8080`, mux)
	if err != nil {
		panic(err)
	}
}
