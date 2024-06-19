package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/a-bondar/go-url-shortener/internal/app/models"
	"github.com/a-bondar/go-url-shortener/internal/app/store"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type Service interface {
	SaveURL(fullURL string) (string, error)
	GetURL(shortURL string) (string, error)
	SaveBatchURLs(urls []models.OriginalURLCorrelation) ([]models.ShortURLCorrelation, error)
	Ping() error
}

type Handler struct {
	s      Service
	logger *zap.Logger
}

func NewHandler(s Service, logger *zap.Logger) *Handler {
	return &Handler{
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

	resURL, err := h.s.SaveURL(string(fullURL))
	if err != nil && !errors.Is(err, store.ErrConflict) {
		h.logger.Error("Failed to shorten URL", zap.Error(err))
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	statusCode := http.StatusCreated
	var conflictErr *store.URLConflictError
	if errors.As(err, &conflictErr) {
		resURL = conflictErr.URL
		statusCode = http.StatusConflict
	}
	w.WriteHeader(statusCode)

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
	var request models.HandleShortenRequest
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)

	if err != nil {
		h.logger.Error("Failed to read body", zap.Error(err))
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &request); err != nil {
		h.logger.Error("Failed to unmarshal request", zap.Error(err))
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	resURL, err := h.s.SaveURL(request.URL)
	if err != nil && !errors.Is(err, store.ErrConflict) {
		h.logger.Error("Failed to shorten URL", zap.Error(err))
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	statusCode := http.StatusCreated
	var conflictErr *store.URLConflictError
	if errors.As(err, &conflictErr) {
		resURL = conflictErr.URL
		statusCode = http.StatusConflict
	}

	resp := models.HandleShortenResponse{
		Result: resURL,
	}

	w.Header().Set("Content-Type", "application/json")

	enc := json.NewEncoder(w)
	w.WriteHeader(statusCode)

	if err := enc.Encode(resp); err != nil {
		h.logger.Error("Failed to encode response", zap.Error(err))
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) HandleShortenBatch(w http.ResponseWriter, r *http.Request) {
	var request models.HandleShortenBatchRequest
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		h.logger.Error("Failed to read body", zap.Error(err))
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &request); err != nil {
		h.logger.Error("Failed to unmarshal request", zap.Error(err))
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	var response models.HandleShortenBatchResponse
	response, err = h.s.SaveBatchURLs(request)
	if err != nil {
		h.logger.Error("Failed to shorten URLs", zap.Error(err))
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	w.WriteHeader(http.StatusCreated)

	if err = enc.Encode(response); err != nil {
		h.logger.Error("Failed to encode response", zap.Error(err))
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) HandleDatabasePing(w http.ResponseWriter, _ *http.Request) {
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
