package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAllowContentEncoding(t *testing.T) {
	tests := []struct {
		name              string
		allowedEncodings  []string
		requestEncodings  []string
		requestContentLen int64
		wantStatusCode    int
	}{
		{"accepts empty body", []string{"gzip"}, []string{"gzip"}, 0, http.StatusOK},
		{"accepts valid encoding", []string{"gzip"}, []string{"gzip"}, 1, http.StatusOK},
		{"accepts any encoding", []string{"gzip"}, []string{"identity"}, 1, http.StatusOK},
		{"accepts multiple encodings", []string{"gzip", "deflate"}, []string{"gzip"}, 1, http.StatusOK},
		{"accepts chained encodings", []string{"gzip", "deflate"}, []string{"gzip", "deflate"}, 1, http.StatusOK},
		{"rejects invalid encoding", []string{"gzip"}, []string{"deflate"}, 1, http.StatusUnsupportedMediaType},
		{"rejects chained invalid encoding", []string{"gzip"}, []string{"gzip", "deflate"}, 1, http.StatusUnsupportedMediaType},
		{"rejects empty allowed encodings", []string{}, []string{"gzip"}, 1, http.StatusUnsupportedMediaType},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			h := AllowContentEncoding(tt.allowedEncodings...)(next)

			req := httptest.NewRequest("POST", "/", nil)
			req.ContentLength = tt.requestContentLen
			req.Header["Content-Encoding"] = tt.requestEncodings
			rec := httptest.NewRecorder()
			h.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatusCode {
				t.Errorf("AllowContentEncoding() = %v, want %v", rec.Code, tt.wantStatusCode)
			}
		})
	}
}