package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewMiddleware(t *testing.T) {
	// Create a simple handler that responds with "Hello, world!"
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, world!"))
	})

	// Create the middleware using the New function
	middleware := New(handler)

	// Create a test handler that uses the middleware
	testHandler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

	// Create a request and response recorder
	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()

	// Serve the request using the test handler
	testHandler.ServeHTTP(rec, req)

	// Check the response
	if rec.Body.String() != "Hello, world!" {
		t.Errorf("Expected response body to be 'Hello, world!', got '%s'", rec.Body.String())
	}
}

func TestContextKeyString(t *testing.T) {
	key := &contextKey{name: "test"}
	expected := "chi/middleware context value test"
	actual := key.String()
	if actual != expected {
		t.Errorf("Expected %s, got %s", expected, actual)
	}
}