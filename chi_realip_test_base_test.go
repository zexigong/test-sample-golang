package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRealIP(t *testing.T) {
	tests := []struct {
		name           string
		headers        map[string]string
		expectedRemote string
	}{
		{
			name:           "True-Client-IP header",
			headers:        map[string]string{"True-Client-IP": "192.168.1.1"},
			expectedRemote: "192.168.1.1",
		},
		{
			name:           "X-Real-IP header",
			headers:        map[string]string{"X-Real-IP": "192.168.1.2"},
			expectedRemote: "192.168.1.2",
		},
		{
			name:           "X-Forwarded-For header",
			headers:        map[string]string{"X-Forwarded-For": "192.168.1.3, 192.168.1.4"},
			expectedRemote: "192.168.1.3",
		},
		{
			name:           "No headers",
			headers:        map[string]string{},
			expectedRemote: "",
		},
		{
			name:           "Invalid IP in headers",
			headers:        map[string]string{"X-Real-IP": "invalid-ip"},
			expectedRemote: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "http://example.com", nil)
			for k, v := range tt.headers {
				req.Header.Set(k, v)
			}
			rr := httptest.NewRecorder()
			handler := RealIP(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.RemoteAddr != tt.expectedRemote {
					t.Errorf("expected RemoteAddr to be %v, got %v", tt.expectedRemote, r.RemoteAddr)
				}
			}))
			handler.ServeHTTP(rr, req)
		})
	}
}