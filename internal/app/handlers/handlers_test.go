package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandleRoot(t *testing.T) {
	testCases := []struct {
		name         string
		method       string
		body         io.Reader
		input        string
		expectedCode int
		expectedBody string
	}{
		{
			name:         "PUT method is not allowed",
			method:       http.MethodPut,
			expectedCode: http.StatusMethodNotAllowed,
		},
		{
			name:         "HEAD method is not allowed",
			method:       http.MethodHead,
			expectedCode: http.StatusMethodNotAllowed,
		},
		{
			name:         "PATCH method is not allowed",
			method:       http.MethodPatch,
			expectedCode: http.StatusMethodNotAllowed,
		},
		{
			name:         "DELETE method is not allowed",
			method:       http.MethodDelete,
			expectedCode: http.StatusMethodNotAllowed,
		},
		{
			name:         "Status 201 if link was shortened successfully",
			method:       http.MethodPost,
			body:         io.NopCloser(strings.NewReader("http://hello.world")),
			expectedCode: http.StatusCreated,
			expectedBody: "http://example.com/aHR0cDovL2hlbGxvLndvcmxk",
		},
		{
			name:         "Status 400 if linkID was not given",
			method:       http.MethodGet,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "Status 404 if link doesn't exist",
			method:       http.MethodGet,
			input:        "12131kjhjhjk",
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "Status 307 if link was found successfully",
			method:       http.MethodGet,
			input:        "aHR0cDovL2hlbGxvLndvcmxk",
			expectedCode: http.StatusTemporaryRedirect,
			expectedBody: "<a href=\"http://hello.world\">Temporary Redirect</a>.\n\n",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := httptest.NewRequest(tc.method, "/"+tc.input, tc.body)
			w := httptest.NewRecorder()

			HandleRoot(w, r)

			assert.Equal(t, tc.expectedCode, w.Code, "Response code is not correct")
			if tc.expectedBody != "" {
				assert.Equal(t, tc.expectedBody, w.Body.String(), "Response body is not correct")
			}
		})
	}
}
