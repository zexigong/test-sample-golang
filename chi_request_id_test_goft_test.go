package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRequestID(t *testing.T) {
	t.Parallel()

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := GetReqID(r.Context())
		if reqID == "" {
			t.Errorf("expected a request ID, got empty string")
		}
		if !strings.HasPrefix(reqID, prefix) {
			t.Errorf("request ID should have prefix '%s', got '%s'", prefix, reqID)
		}
		w.WriteHeader(http.StatusOK)
	})

	handler := RequestID(next)

	t.Run("request without RequestID", func(t *testing.T) {
		t.Parallel()

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("expected status code %d, got %d", http.StatusOK, rr.Code)
		}
	})

	t.Run("request with RequestID", func(t *testing.T) {
		t.Parallel()

		const requestID = "test-request-id"

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set(RequestIDHeader, requestID)
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("expected status code %d, got %d", http.StatusOK, rr.Code)
		}
		reqID := GetReqID(req.Context())
		if reqID != requestID {
			t.Errorf("expected request ID %q, got %q", requestID, reqID)
		}
	})
}

func TestGetReqID(t *testing.T) {
	t.Parallel()

	t.Run("without RequestID", func(t *testing.T) {
		t.Parallel()

		if id := GetReqID(context.Background()); id != "" {
			t.Errorf("request ID should be empty, got %s", id)
		}
	})

	t.Run("with RequestID", func(t *testing.T) {
		t.Parallel()

		const id = "test-request-id"

		ctx := context.WithValue(context.Background(), RequestIDKey, id)

		if reqID := GetReqID(ctx); reqID != id {
			t.Errorf("expected request ID %q, got %q", id, reqID)
		}
	})
}