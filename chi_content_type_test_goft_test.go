package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSetHeader(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()

	h := SetHeader("X-Test-Header", "TestValue")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	h.ServeHTTP(rr, req)

	if rr.Header().Get("X-Test-Header") != "TestValue" {
		t.Errorf("expected header X-Test-Header to be set to 'TestValue', got '%s'", rr.Header().Get("X-Test-Header"))
	}
}

func TestAllowContentType(t *testing.T) {
	tests := []struct {
		contentType  string
		allowedTypes []string
		expectedCode int
	}{
		{"application/json", []string{"application/json"}, http.StatusOK},
		{"application/xml", []string{"application/json", "application/xml"}, http.StatusOK},
		{"text/plain", []string{"application/json", "application/xml"}, http.StatusUnsupportedMediaType},
	}

	for _, tt := range tests {
		t.Run(tt.contentType, func(t *testing.T) {
			req, _ := http.NewRequest("POST", "/", strings.NewReader("{}"))
			req.Header.Set("Content-Type", tt.contentType)

			rr := httptest.NewRecorder()

			h := AllowContentType(tt.allowedTypes...)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			}))
			h.ServeHTTP(rr, req)

			if rr.Code != tt.expectedCode {
				t.Errorf("expected status code %d, got %d", tt.expectedCode, rr.Code)
			}
		})
	}
}

func TestAllowContentType_EmptyBody(t *testing.T) {
	req, _ := http.NewRequest("POST", "/", nil)
	req.Header.Set("Content-Type", "text/plain")
	rr := httptest.NewRecorder()

	h := AllowContentType("application/json")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, rr.Code)
	}
}