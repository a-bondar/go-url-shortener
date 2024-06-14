package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	"github.com/a-bondar/go-url-shortener/internal/app/config"
	"github.com/a-bondar/go-url-shortener/internal/app/models"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type Service interface {
	SaveURL(fullURL string) (string, error)
	GetURL(shortURL string) (string, error)
	Ping() error
}

type Handler struct {
	cfg    *config.Config
	s      Service
	logger *zap.Logger
}

func NewHandler(cfg *config.Config, s Service, logger *zap.Logger) *Handler {
	return &Handler{
		cfg:    cfg,
		s:      s,
		logger: logger,
	}
}

func (h *Handler) HandlePost(w http.ResponseWriter, r *http.Request) {
	fullURL, err := io.ReadAll(r.Body)

	if err != nil {
		h.logger.Error("Failed to read body", zap.Error(err))
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	shortURL, err := h.s.SaveURL(string(fullURL))

	if err != nil {
		h.logger.Error("Failed to shorten URL", zap.Error(err))
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)

	resURL, err := url.JoinPath(h.cfg.ShortLinkBaseURL, shortURL)

	if err != nil {
		h.logger.Error("Failed to build URL", zap.Error(err))
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	if _, err := w.Write([]byte(resURL)); err != nil {
		h.logger.Error("Failed to write result", zap.Error(err))
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) HandleGet(w http.ResponseWriter, r *http.Request) {
	linkID := chi.URLParam(r, "linkID")

	URL, err := h.s.GetURL(linkID)

	if err != nil {
		h.logger.Error("Failed to get URL", zap.Error(err))
		http.Error(w, `Link not found`, http.StatusNotFound)
		return
	}

	http.Redirect(w, r, URL, http.StatusTemporaryRedirect)
}

func (h *Handler) HandleShorten(w http.ResponseWriter, r *http.Request) {
	var request models.Request
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)

	if err != nil {
		h.logger.Error("Failed to read body", zap.Error(err))
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	if err := json.Unmarshal(buf.Bytes(), &request); err != nil {
		h.logger.Error("Failed to unmarshal request", zap.Error(err))
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	shortURL, err := h.s.SaveURL(request.URL)

	if err != nil {
		h.logger.Error("Failed to shorten URL", zap.Error(err))
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	resURL, err := url.JoinPath(h.cfg.ShortLinkBaseURL, shortURL)

	if err != nil {
		h.logger.Error("Failed to build URL", zap.Error(err))
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	resp := models.Response{
		Result: resURL,
	}

	w.Header().Set("Content-Type", "application/json")

	enc := json.NewEncoder(w)
	w.WriteHeader(http.StatusCreated)

	if err := enc.Encode(resp); err != nil {
		h.logger.Error("Failed to encode response", zap.Error(err))
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) HandleDatabasePing(w http.ResponseWriter, r *http.Request) {
	err := h.s.Ping()

	if err != nil {
		h.logger.Error("Unable to reach DB", zap.Error(err))
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if _, err := w.Write([]byte(`{"status": "ok"}`)); err != nil {
		h.logger.Error("Failed to write response", zap.Error(err))
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}
