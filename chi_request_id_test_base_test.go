package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRequestIDMiddleware(t *testing.T) {
	handler := RequestID(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := GetReqID(r.Context())
		if reqID == "" {
			t.Error("expected a request ID, got empty string")
		}
		w.WriteHeader(http.StatusOK)
	}))

	req, err := http.NewRequest("GET", "http://example.com/foo", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status OK; got %v", rr.Code)
	}

	reqID := rr.Header().Get(RequestIDHeader)
	if reqID == "" {
		t.Error("expected a request ID in the response header, got empty string")
	}
}

func TestRequestIDWithExistingHeader(t *testing.T) {
	handler := RequestID(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := GetReqID(r.Context())
		expectedID := "test-existing-id"
		if reqID != expectedID {
			t.Errorf("expected request ID %s, got %s", expectedID, reqID)
		}
		w.WriteHeader(http.StatusOK)
	}))

	req, err := http.NewRequest("GET", "http://example.com/foo", nil)
	if err != nil {
		t.Fatal(err)
	}

	existingID := "test-existing-id"
	req.Header.Set(RequestIDHeader, existingID)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status OK; got %v", rr.Code)
	}

	reqID := rr.Header().Get(RequestIDHeader)
	if reqID != existingID {
		t.Errorf("expected request ID in the response header to be %s, got %s", existingID, reqID)
	}
}

func TestGetReqID(t *testing.T) {
	ctx := context.Background()
	if id := GetReqID(ctx); id != "" {
		t.Errorf("expected empty request ID, got %s", id)
	}

	ctx = context.WithValue(ctx, RequestIDKey, "12345")
	if id := GetReqID(ctx); id != "12345" {
		t.Errorf("expected request ID '12345', got %s", id)
	}
}

func TestNextRequestID(t *testing.T) {
	id1 := NextRequestID()
	id2 := NextRequestID()

	if id2 != id1+1 {
		t.Errorf("expected NextRequestID to increment by 1, got %d, %d", id1, id2)
	}
}

func TestInitPrefix(t *testing.T) {
	if prefix == "" {
		t.Error("expected non-empty prefix after initialization")
	}

	if strings.Contains(prefix, "+") || strings.Contains(prefix, "/") {
		t.Error("prefix should not contain '+' or '/' characters")
	}
}