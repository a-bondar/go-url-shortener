package router

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/a-bondar/go-url-shortener/internal/app/config"
	"github.com/a-bondar/go-url-shortener/internal/app/handlers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testRequest(t *testing.T, ts *httptest.Server, method, path string, body io.Reader) (*http.Response, string) {
	t.Helper()

	req, err := http.NewRequest(method, ts.URL+path, body)
	require.NoError(t, err)

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

func (s *serviceMock) SaveURL(_ string) (string, error) {
	return "qw12qw", nil
}

func (s *serviceMock) GetURL(shortURL string) (string, error) {
	if shortURL != "qw12qw" {
		return "", errors.New("link not found")
	}

	return "https://hello.world", nil
}

func TestRouter(t *testing.T) {
	cfg := config.NewConfig()
	svc := &serviceMock{}
	h := handlers.NewHandler(cfg, svc)

	ts := httptest.NewServer(Router(h))
	defer ts.Close()

	testCases := []struct {
		name             string
		method           string
		body             io.Reader
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
			path:         "/",
			body:         io.NopCloser(strings.NewReader("https://hello.world")),
			expectedCode: http.StatusCreated,
			expectedBody: "http://localhost:8080/qw12qw",
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
			resp, body := testRequest(t, ts, tc.method, tc.path, tc.body)

			defer func() {
				if err := resp.Body.Close(); err != nil {
					fmt.Println("Error closing response body:", err)
				}
			}()

			assert.Equal(t, tc.expectedCode, resp.StatusCode, "Response code is not correct")

			if tc.expectedBody != "" {
				assert.Equal(t, tc.expectedBody, body, "Response body is not correct")
			}

			if tc.expectedLocation != "" {
				assert.Equal(t, tc.expectedLocation, resp.Header.Get("Location"))
			}
		})
	}
}
