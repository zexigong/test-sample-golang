package middleware

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAllowContentType(t *testing.T) {
	noop := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	tests := []struct {
		name         string
		contentTypes []string
		req          *http.Request
		wantCode     int
	}{
		{
			name:         "allow no body",
			contentTypes: []string{"application/json"},
			req:          httptest.NewRequest("POST", "/", nil),
			wantCode:     http.StatusOK,
		},
		{
			name:         "allow empty body",
			contentTypes: []string{"application/json"},
			req:          httptest.NewRequest("POST", "/", strings.NewReader("")),
			wantCode:     http.StatusOK,
		},
		{
			name:         "allow json",
			contentTypes: []string{"application/json"},
			req: func() *http.Request {
				r := httptest.NewRequest("POST", "/", strings.NewReader("{}"))
				r.Header.Set("Content-Type", "application/json")
				return r
			}(),
			wantCode: http.StatusOK,
		},
		{
			name:         "allow json with charset",
			contentTypes: []string{"application/json"},
			req: func() *http.Request {
				r := httptest.NewRequest("POST", "/", strings.NewReader("{}"))
				r.Header.Set("Content-Type", "application/json; charset=UTF-8")
				return r
			}(),
			wantCode: http.StatusOK,
		},
		{
			name:         "allow case insensitive",
			contentTypes: []string{"application/JSON"},
			req: func() *http.Request {
				r := httptest.NewRequest("POST", "/", strings.NewReader("{}"))
				r.Header.Set("Content-Type", "APPLICATION/JSON")
				return r
			}(),
			wantCode: http.StatusOK,
		},
		{
			name:         "disallow xml",
			contentTypes: []string{"application/json"},
			req: func() *http.Request {
				r := httptest.NewRequest("POST", "/", strings.NewReader("<xml></xml>"))
				r.Header.Set("Content-Type", "application/xml")
				return r
			}(),
			wantCode: http.StatusUnsupportedMediaType,
		},
		{
			name:         "disallow xml no content type",
			contentTypes: []string{"application/json"},
			req: func() *http.Request {
				r := httptest.NewRequest("POST", "/", strings.NewReader("<xml></xml>"))
				return r
			}(),
			wantCode: http.StatusUnsupportedMediaType,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			handler := AllowContentType(tt.contentTypes...)(noop)
			handler.ServeHTTP(rr, tt.req)
			if rr.Code != tt.wantCode {
				t.Errorf("AllowContentType() statusCode = %d, want %d", rr.Code, tt.wantCode)
			}
		})
	}
}

func TestSetHeader(t *testing.T) {
	noop := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	tests := []struct {
		name     string
		key      string
		value    string
		req      *http.Request
		wantCode int
	}{
		{
			name:     "set header",
			key:      "Content-Type",
			value:    "application/json",
			req:      httptest.NewRequest("GET", "/", nil),
			wantCode: http.StatusOK,
		},
		{
			name:     "set header with body",
			key:      "Content-Type",
			value:    "application/json",
			req:      httptest.NewRequest("POST", "/", strings.NewReader("{}")),
			wantCode: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			handler := SetHeader(tt.key, tt.value)(noop)
			handler.ServeHTTP(rr, tt.req)

			if got := rr.Header().Get(tt.key); got != tt.value {
				t.Errorf("SetHeader() = %s, want %s", got, tt.value)
			}
			if rr.Code != tt.wantCode {
				t.Errorf("SetHeader() statusCode = %d, want %d", rr.Code, tt.wantCode)
			}
		})
	}
}