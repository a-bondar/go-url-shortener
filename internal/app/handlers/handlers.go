package handlers

import (
	"encoding/base64"
	"io"
	"log"
	"net/http"

	"github.com/a-bondar/go-url-shortener/internal/app/config"
	"github.com/go-chi/chi/v5"
)

var linksMap = map[string]string{}

func HandlePost(cfg *config.Config, w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sEnc := base64.StdEncoding.EncodeToString(body)
	linksMap[sEnc] = string(body)
	w.WriteHeader(http.StatusCreated)
	if _, err := w.Write([]byte(cfg.ShortLinkBaseURL + "/" + sEnc)); err != nil {
		log.Printf("failed to write result: %v", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}

func HandleGet(w http.ResponseWriter, r *http.Request) {
	linkID := chi.URLParam(r, "linkID")

	link, ok := linksMap[linkID]

	if !ok {
		http.Error(w, `Link not found`, http.StatusNotFound)
		return
	}

	http.Redirect(w, r, link, http.StatusTemporaryRedirect)
}
