package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAllowContentEncoding(t *testing.T) {
	tests := []struct {
		name              string
		allowedEncodings  []string
		requestEncoding   string
		expectedStatus    int
	}{
		{
			name:             "No Content-Encoding with content",
			allowedEncodings: []string{"gzip"},
			requestEncoding:  "",
			expectedStatus:   http.StatusOK,
		},
		{
			name:             "No Content-Encoding with no content",
			allowedEncodings: []string{"gzip"},
			requestEncoding:  "",
			expectedStatus:   http.StatusOK,
		},
		{
			name:             "Allowed encoding",
			allowedEncodings: []string{"gzip"},
			requestEncoding:  "gzip",
			expectedStatus:   http.StatusOK,
		},
		{
			name:             "Disallowed encoding",
			allowedEncodings: []string{"gzip"},
			requestEncoding:  "deflate",
			expectedStatus:   http.StatusUnsupportedMediaType,
		},
		{
			name:             "Multiple allowed encodings",
			allowedEncodings: []string{"gzip", "deflate"},
			requestEncoding:  "deflate",
			expectedStatus:   http.StatusOK,
		},
		{
			name:             "Encoding with spaces",
			allowedEncodings: []string{"gzip "},
			requestEncoding:  " gzip ",
			expectedStatus:   http.StatusOK,
		},
		{
			name:             "Encoding case insensitivity",
			allowedEncodings: []string{"GZIP"},
			requestEncoding:  "gzip",
			expectedStatus:   http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := AllowContentEncoding(tt.allowedEncodings...)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			}))

			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("request body"))
			req.Header.Set("Content-Encoding", tt.requestEncoding)

			if tt.requestEncoding == "" {
				req.ContentLength = 0
			} else {
				req.ContentLength = int64(len("request body"))
			}

			w := httptest.NewRecorder()
			handler.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %v, got %v", tt.expectedStatus, w.Code)
			}
		})
	}
}