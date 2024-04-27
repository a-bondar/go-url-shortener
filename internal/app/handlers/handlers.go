package handlers

import (
	"encoding/base64"
	"io"
	"net/http"
	"strings"
)

var linksMap = map[string]string{}

func handlePost(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
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

	linkID := strings.Split(r.URL.Path[1:], "/")[0]

	if linkID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	link, ok := linksMap[linkID]

	if !ok {
		http.Error(w, `Link not found`, http.StatusNotFound)
		return
	}

	http.Redirect(w, r, link, http.StatusTemporaryRedirect)
}

func HandleRoot(w http.ResponseWriter, r *http.Request) {
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
