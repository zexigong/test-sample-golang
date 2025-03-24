package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestContentCharset(t *testing.T) {
	tests := []struct {
		name           string
		charsets       []string
		contentType    string
		expectedStatus int
	}{
		{
			name:           "Supported charset",
			charsets:       []string{"utf-8"},
			contentType:    "text/plain; charset=utf-8",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Unsupported charset",
			charsets:       []string{"iso-8859-1"},
			contentType:    "text/plain; charset=utf-8",
			expectedStatus: http.StatusUnsupportedMediaType,
		},
		{
			name:           "No charset in Content-Type",
			charsets:       []string{"utf-8"},
			contentType:    "text/plain",
			expectedStatus: http.StatusUnsupportedMediaType,
		},
		{
			name:           "Empty charset allows all",
			charsets:       []string{""},
			contentType:    "text/plain; charset=utf-8",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Empty charset with no Content-Type",
			charsets:       []string{""},
			contentType:    "",
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := ContentCharset(tt.charsets...)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			}))

			req, err := http.NewRequest("GET", "/", nil)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}
			req.Header.Set("Content-Type", tt.contentType)

			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("Expected status %v, got %v", tt.expectedStatus, rr.Code)
			}
		})
	}
}

func TestContentEncoding(t *testing.T) {
	tests := []struct {
		name     string
		ce       string
		charsets []string
		expected bool
	}{
		{
			name:     "Match charset",
			ce:       "text/html; charset=utf-8",
			charsets: []string{"utf-8"},
			expected: true,
		},
		{
			name:     "No charset in content type",
			ce:       "text/html",
			charsets: []string{"utf-8"},
			expected: false,
		},
		{
			name:     "Empty charset list",
			ce:       "text/html; charset=utf-8",
			charsets: []string{},
			expected: false,
		},
		{
			name:     "Empty charset allows all",
			ce:       "text/html; charset=iso-8859-1",
			charsets: []string{""},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := contentEncoding(tt.ce, tt.charsets...)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestSplit(t *testing.T) {
	tests := []struct {
		str      string
		sep      string
		expectedA string
		expectedB string
	}{
		{
			str:      "text/html; charset=utf-8",
			sep:      ";",
			expectedA: "text/html",
			expectedB: "charset=utf-8",
		},
		{
			str:      "charset=utf-8",
			sep:      "=",
			expectedA: "charset",
			expectedB: "utf-8",
		},
		{
			str:      "text/html",
			sep:      ";",
			expectedA: "text/html",
			expectedB: "",
		},
		{
			str:      "  key=value  ",
			sep:      "=",
			expectedA: "key",
			expectedB: "value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.str, func(t *testing.T) {
			a, b := split(tt.str, tt.sep)
			if a != tt.expectedA || b != tt.expectedB {
				t.Errorf("Expected (%v, %v), got (%v, %v)", tt.expectedA, tt.expectedB, a, b)
			}
		})
	}
}