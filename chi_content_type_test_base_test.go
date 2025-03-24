package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5/middleware"
)

func TestSetHeader(t *testing.T) {
	handler := middleware.SetHeader("X-Test-Header", "TestValue")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if got := w.Header().Get("X-Test-Header"); got != "TestValue" {
		t.Errorf("expected header X-Test-Header to be 'TestValue', got '%s'", got)
	}
}

func TestAllowContentType(t *testing.T) {
	handler := middleware.AllowContentType("application/json", "application/xml")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	tests := []struct {
		name          string
		contentType   string
		expectedCode  int
	}{
		{"Allowed JSON", "application/json", http.StatusOK},
		{"Allowed XML", "application/xml", http.StatusOK},
		{"Disallowed HTML", "text/html", http.StatusUnsupportedMediaType},
		{"Empty ContentType", "", http.StatusOK},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/", strings.NewReader("{}"))
			req.Header.Set("Content-Type", tc.contentType)
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)

			if w.Code != tc.expectedCode {
				t.Errorf("expected status %d, got %d", tc.expectedCode, w.Code)
			}
		})
	}
}