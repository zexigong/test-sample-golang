package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestThrottle(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	throttledHandler := Throttle(1)(handler)

	req, _ := http.NewRequest("GET", "/", nil)
	resp := httptest.NewRecorder()

	throttledHandler.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("expected status %v, got %v", http.StatusOK, resp.Code)
	}
}

func TestThrottleBacklog(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	throttledHandler := ThrottleBacklog(1, 1, time.Second*2)(handler)

	req, _ := http.NewRequest("GET", "/", nil)
	resp := httptest.NewRecorder()

	throttledHandler.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("expected status %v, got %v", http.StatusOK, resp.Code)
	}

	// Sending another request to test backlog
	resp2 := httptest.NewRecorder()
	throttledHandler.ServeHTTP(resp2, req)

	if resp2.Code != http.StatusOK {
		t.Errorf("expected status %v, got %v", http.StatusOK, resp2.Code)
	}
}

func TestThrottleWithOpts(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	opts := ThrottleOpts{
		Limit:          1,
		BacklogLimit:   1,
		BacklogTimeout: time.Second * 2,
		StatusCode:     http.StatusTooManyRequests,
	}

	throttledHandler := ThrottleWithOpts(opts)(handler)

	req, _ := http.NewRequest("GET", "/", nil)
	resp := httptest.NewRecorder()

	throttledHandler.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("expected status %v, got %v", http.StatusOK, resp.Code)
	}

	// Sending another request to test backlog
	resp2 := httptest.NewRecorder()
	throttledHandler.ServeHTTP(resp2, req)

	if resp2.Code != http.StatusOK {
		t.Errorf("expected status %v, got %v", http.StatusOK, resp2.Code)
	}

	// Sending another request to exceed capacity
	resp3 := httptest.NewRecorder()
	throttledHandler.ServeHTTP(resp3, req)

	if resp3.Code != opts.StatusCode {
		t.Errorf("expected status %v, got %v", opts.StatusCode, resp3.Code)
	}
}