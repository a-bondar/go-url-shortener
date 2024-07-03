package router

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/a-bondar/go-url-shortener/internal/app/handlers"
	"github.com/a-bondar/go-url-shortener/internal/app/middleware"
	"github.com/a-bondar/go-url-shortener/internal/app/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

const userID = "12345"

func testRequest(t *testing.T, ts *httptest.Server, method, path string, body io.Reader) (*http.Response, string) {
	t.Helper()

	req, err := http.NewRequest(method, ts.URL+path, body)
	require.NoError(t, err)

	token, _ := middleware.CreateAccessToken(userID)
	req.AddCookie(&http.Cookie{
		Name:  "auth_token",
		Value: token,
		Path:  "/",
	})

	ts.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	resp, err := ts.Client().Do(req)
	require.NoError(t, err)

	defer func() {
		if err := resp.Body.Close(); err != nil {
			fmt.Println("Error closing response body:", err)
		}
	}()

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return resp, string(respBody)
}

type serviceMock struct{}

func (s *serviceMock) SaveURL(_ context.Context, _ string, _ string) (string, error) {
	return "http://localhost:8080/qw12qw", nil
}

func (s *serviceMock) GetURL(_ context.Context, shortURL string, _ string) (string, error) {
	if shortURL != "qw12qw" {
		return "", errors.New("link not found")
	}

	return "https://hello.world", nil
}

func (s *serviceMock) GetURLs(_ context.Context, _ string) ([]models.URLsPair, error) {
	return nil, nil
}

func (s *serviceMock) SaveBatchURLs(
	_ context.Context,
	urls []models.OriginalURLCorrelation, _ string) ([]models.ShortURLCorrelation, error) {
	res := make([]models.ShortURLCorrelation, 0, len(urls))

	for _, url := range urls {
		res = append(res, models.ShortURLCorrelation{
			CorrelationID: url.CorrelationID,
			ShortURL:      "qw12qw",
		})
	}

	return res, nil
}

func (s *serviceMock) Ping(_ context.Context) error {
	return nil
}

func TestRouter(t *testing.T) {
	logger := zap.NewNop()
	svc := &serviceMock{}
	h := handlers.NewHandler(svc, logger)

	ts := httptest.NewServer(Router(h, logger))
	defer ts.Close()

	testCases := []struct {
		name             string
		method           string
		body             string
		path             string
		expectedCode     int
		expectedBody     string
		expectedLocation string
	}{
		{
			name:         "PUT method is not allowed",
			method:       http.MethodPut,
			path:         "/",
			expectedCode: http.StatusMethodNotAllowed,
		},
		{
			name:         "HEAD method is not allowed",
			method:       http.MethodHead,
			path:         "/",
			expectedCode: http.StatusMethodNotAllowed,
		},
		{
			name:         "PATCH method is not allowed",
			method:       http.MethodPatch,
			path:         "/",
			expectedCode: http.StatusMethodNotAllowed,
		},
		{
			name:         "DELETE method is not allowed",
			method:       http.MethodDelete,
			path:         "/",
			expectedCode: http.StatusMethodNotAllowed,
		},
		{
			name:         "GET method on \"/\" is not allowed",
			method:       http.MethodGet,
			path:         "/",
			expectedCode: http.StatusMethodNotAllowed,
		},
		{
			name:         "POST method on \"/{linkID}\" is not allowed",
			method:       http.MethodPost,
			path:         "/121212",
			expectedCode: http.StatusMethodNotAllowed,
		},
		{
			name:         "Status 201 if link was shortened successfully",
			method:       http.MethodPost,
			path:         "/api/shorten",
			body:         `{"url":  "https://hello.world"}`,
			expectedCode: http.StatusCreated,
			expectedBody: `{"result": "http://localhost:8080/qw12qw"}`,
		},
		{
			name:   "Status 201 if links was shortened successfully",
			method: http.MethodPost,
			path:   "/api/shorten/batch",
			body: `[{"correlation_id": "1", "original_url": "https://example.com/1"},
					{"correlation_id": "2", "original_url": "https://example.com/2"}]`,
			expectedCode: http.StatusCreated,
			expectedBody: `[{"correlation_id":"1","short_url":"qw12qw"},{"correlation_id":"2","short_url":"qw12qw"}]`,
		},
		{
			name:         "Status 404 if link doesn't exist",
			method:       http.MethodGet,
			path:         "/12131kjhjhjk",
			expectedCode: http.StatusNotFound,
		},
		{
			name:             "Status 307 if link was found successfully",
			method:           http.MethodGet,
			path:             "/qw12qw",
			expectedCode:     http.StatusTemporaryRedirect,
			expectedLocation: "https://hello.world",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resp, body := testRequest(t, ts, tc.method, tc.path, bytes.NewBufferString(tc.body))

			defer func() {
				if err := resp.Body.Close(); err != nil {
					fmt.Println("Error closing response body:", err)
				}
			}()

			assert.Equal(t, tc.expectedCode, resp.StatusCode, "Response code is not correct")

			if tc.expectedBody != "" {
				assert.JSONEq(t, tc.expectedBody, body, "Response body is not correct")
			}

			if tc.expectedLocation != "" {
				assert.Equal(t, tc.expectedLocation, resp.Header.Get("Location"))
			}
		})
	}
}
