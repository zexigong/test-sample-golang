package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRealIP(t *testing.T) {
	tests := []struct {
		name         string
		headers      map[string]string
		expectedAddr string
	}{
		{
			name: "Using True-Client-IP",
			headers: map[string]string{
				"True-Client-IP": "192.0.2.1",
			},
			expectedAddr: "192.0.2.1",
		},
		{
			name: "Using X-Real-IP",
			headers: map[string]string{
				"X-Real-IP": "192.0.2.2",
			},
			expectedAddr: "192.0.2.2",
		},
		{
			name: "Using X-Forwarded-For",
			headers: map[string]string{
				"X-Forwarded-For": "192.0.2.3, 198.51.100.1",
			},
			expectedAddr: "192.0.2.3",
		},
		{
			name: "Using all headers",
			headers: map[string]string{
				"True-Client-IP":  "203.0.113.1",
				"X-Real-IP":       "192.0.2.2",
				"X-Forwarded-For": "192.0.2.3, 198.51.100.1",
			},
			expectedAddr: "203.0.113.1",
		},
		{
			name: "No headers",
			headers: map[string]string{
				"Authorization": "foo",
			},
			expectedAddr: "",
		},
		{
			name: "Using invalid IP",
			headers: map[string]string{
				"True-Client-IP": "Invalid IP",
			},
			expectedAddr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.RemoteAddr != tt.expectedAddr {
					t.Errorf("realIP = %v, want %v", r.RemoteAddr, tt.expectedAddr)
				}
			})

			handler := RealIP(next)

			req := httptest.NewRequest("GET", "/", nil)
			for k, v := range tt.headers {
				req.Header.Set(k, v)
			}

			w := httptest.NewRecorder()
			handler.ServeHTTP(w, req)
		})
	}
}