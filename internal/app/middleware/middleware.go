package middleware

import (
	"compress/gzip"
	"net/http"
	"strings"
	"time"

	"go.uber.org/zap"
)

type (
	responseData struct {
		status int
		size   int
	}

	loggingResponseWriter struct {
		http.ResponseWriter
		responseData *responseData
	}
)

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size

	return size, err //nolint:wrapcheck // reimplement the interface and do not want to wrap the error
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}

func WithLogging(logger *zap.Logger) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			rd := &responseData{
				status: 0,
				size:   0,
			}
			lrw := &loggingResponseWriter{
				ResponseWriter: w,
				responseData:   rd,
			}

			h.ServeHTTP(lrw, r)

			duration := time.Since(start)

			logger.Info("Request",
				zap.String("method", r.Method),
				zap.String("uri", r.RequestURI),
				zap.Int("status", lrw.responseData.status),
				zap.Int("size", lrw.responseData.size),
				zap.Duration("duration", duration),
			)
		})
	}
}

type gzipResponseWriter struct {
	http.ResponseWriter
	Writer *gzip.Writer
}

func (w *gzipResponseWriter) Write(data []byte) (int, error) {
	return w.Writer.Write(data) //nolint:wrapcheck // reimplement the interface and do not want to wrap the error
}

func WithGzip(logger *zap.Logger) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
				gzReader, err := gzip.NewReader(r.Body)

				if err != nil {
					logger.Error("Cannot create gzip reader", zap.Error(err))
					http.Error(w, "Cannot create gzip reader", http.StatusInternalServerError)
					return
				}

				defer func() {
					if err := gzReader.Close(); err != nil {
						logger.Error("Cannot close gzip reader")
					}
				}()

				r.Body = gzReader
			}

			if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
				gzWriter := gzip.NewWriter(w)

				defer func() {
					if err := gzWriter.Close(); err != nil {
						logger.Error("Cannot close gzip writer")
					}
				}()

				w.Header().Set("Content-Encoding", "gzip")
				w = &gzipResponseWriter{ResponseWriter: w, Writer: gzWriter}
			}

			h.ServeHTTP(w, r)
		})
	}
}
