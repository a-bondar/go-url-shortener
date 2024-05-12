package handlers

import (
	"encoding/base64"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/a-bondar/go-url-shortener/internal/app/config"
	"github.com/a-bondar/go-url-shortener/internal/app/service"
	"github.com/go-chi/chi/v5"
)

//var linksMap = map[string]string{}

type Service interface {
	SaveURL(fullURL string) (string, error)
	GetURL(shortURL string) (string, error)
}

type Handler struct {
	cfg *config.Config
	s   Service
}

func NewHanlder(cfg *config.Config, s Service) *Handler {
	return &Handler{
		cfg: cfg,
		s:   s,
	}
}

func (h *Handler) HandlePost(w http.ResponseWriter, r *http.Request) {
	fullURL, err := io.ReadAll(r.Body)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	shortURL, err := service.SaveURL(string(fullURL))

	// check err and send short URL to user

	sEnc := base64.StdEncoding.EncodeToString(body)
	var fullURL, shortURL string
	h.s.SaveURL(fullURL, shortURL)
	//linksMap[sEnc] = string(body)
	w.WriteHeader(http.StatusCreated)

	resURL, err := url.JoinPath(h.cfg.ShortLinkBaseURL, sEnc)
	if err != nil {
		log.Printf("failed to build url: %v", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	if _, err := w.Write([]byte(resURL)); err != nil {
		log.Printf("failed to write result: %v", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) HandleGet(w http.ResponseWriter, r *http.Request) {
	linkID := chi.URLParam(r, "linkID")

	link, ok := linksMap[linkID]

	if !ok {
		http.Error(w, `Link not found`, http.StatusNotFound)
		return
	}

	http.Redirect(w, r, link, http.StatusTemporaryRedirect)
}
