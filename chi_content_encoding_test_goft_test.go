package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAllowContentEncoding(t *testing.T) {
	tests := []struct {
		name             string
		allowedEncodings []string
		requestEncoding  string
		statusCode       int
	}{
		{
			name:             "Allowed content encoding",
			allowedEncodings: []string{"gzip", "deflate"},
			requestEncoding:  "gzip",
			statusCode:       http.StatusOK,
		},
		{
			name:             "Allowed content encoding with multiple encodings",
			allowedEncodings: []string{"gzip", "deflate"},
			requestEncoding:  "gzip, deflate",
			statusCode:       http.StatusOK,
		},
		{
			name:             "Disallowed content encoding",
			allowedEncodings: []string{"gzip", "deflate"},
			requestEncoding:  "br",
			statusCode:       http.StatusUnsupportedMediaType,
		},
		{
			name:             "Empty content encoding",
			allowedEncodings: []string{"gzip", "deflate"},
			requestEncoding:  "",
			statusCode:       http.StatusOK,
		},
		{
			name:             "No content length",
			allowedEncodings: []string{"gzip", "deflate"},
			requestEncoding:  "gzip",
			statusCode:       http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := AllowContentEncoding(tt.allowedEncodings...)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			}))

			req := httptest.NewRequest(http.MethodPost, "/test", nil)
			if tt.requestEncoding != "" {
				req.Header.Set("Content-Encoding", tt.requestEncoding)
				req.ContentLength = 1
			}

			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)

			if rec.Code != tt.statusCode {
				t.Errorf("unexpected status code: got %v, want %v", rec.Code, tt.statusCode)
			}
		})
	}
}

func TestAllowContentEncodingMultiple(t *testing.T) {
	handler := AllowContentEncoding("gzip", "deflate")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodPost, "/test", nil)
	req.Header.Set("Content-Encoding", "gzip, deflate")
	req.ContentLength = 1

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("unexpected status code: got %v, want %v", rec.Code, http.StatusOK)
	}
}

func TestAllowContentEncodingMixedCase(t *testing.T) {
	handler := AllowContentEncoding("gzip")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodPost, "/test", nil)
	req.Header.Set("Content-Encoding", "GZIP")
	req.ContentLength = 1

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("unexpected status code: got %v, want %v", rec.Code, http.StatusOK)
	}
}

func TestAllowContentEncodingTrimSpaces(t *testing.T) {
	handler := AllowContentEncoding("gzip")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodPost, "/test", nil)
	req.Header.Set("Content-Encoding", " gzip ")
	req.ContentLength = 1

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("unexpected status code: got %v, want %v", rec.Code, http.StatusOK)
	}
}

func TestAllowContentEncodingNoContentLength(t *testing.T) {
	handler := AllowContentEncoding("gzip")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodPost, "/test", nil)
	req.Header.Set("Content-Encoding", "gzip")

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("unexpected status code: got %v, want %v", rec.Code, http.StatusOK)
	}
}

func TestAllowContentEncodingEmptyHeader(t *testing.T) {
	handler := AllowContentEncoding("gzip")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodPost, "/test", strings.NewReader("test"))
	req.Header.Set("Content-Encoding", "")

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("unexpected status code: got %v, want %v", rec.Code, http.StatusOK)
	}
}