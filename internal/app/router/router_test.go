package router

import (
	"github.com/a-bondar/go-url-shortener/internal/app/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func testRequest(t *testing.T, ts *httptest.Server, method, path string, body io.Reader) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, body)
	require.NoError(t, err)

	ts.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	resp, err := ts.Client().Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return resp, string(respBody)
}

func TestRouter(t *testing.T) {
	config.ParseFlags()
	ts := httptest.NewServer(Router())
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
			body:         io.NopCloser(strings.NewReader("http://hello.world")),
			expectedCode: http.StatusCreated,
			expectedBody: ts.URL + "/aHR0cDovL2hlbGxvLndvcmxk",
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
			path:             "/aHR0cDovL2hlbGxvLndvcmxk",
			expectedCode:     http.StatusTemporaryRedirect,
			expectedLocation: "http://hello.world",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resp, body := testRequest(t, ts, tc.method, tc.path, tc.body)
			defer resp.Body.Close()

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
