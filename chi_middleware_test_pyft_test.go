package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNew(t *testing.T) {
	// Create a mock handler to test the middleware
	mockHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Use the middleware to wrap the mock handler
	middlewareHandler := New(mockHandler)(mockHandler)

	// Create a request to pass to the handler
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Call the middleware handler with the request and ResponseRecorder
	middlewareHandler.ServeHTTP(rr, req)

	// Check the status code is what we expect
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestContextKey_String(t *testing.T) {
	// Create a new context key with a name
	key := &contextKey{name: "test"}

	// Get the string representation of the context key
	got := key.String()
	want := "chi/middleware context value test"

	// Check if the string representation matches what we expect
	if got != want {
		t.Errorf("contextKey.String() = %v, want %v", got, want)
	}
}