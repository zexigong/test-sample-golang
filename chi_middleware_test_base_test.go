package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewMiddleware(t *testing.T) {
	// Create a sample handler to be wrapped by the middleware
	sampleHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, World!"))
	})

	// Create the middleware
	middleware := New(sampleHandler)

	// Create a test server using the middleware
	ts := httptest.NewServer(middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Middleware Test"))
	})))
	defer ts.Close()

	resp, err := http.Get(ts.URL)
	if err != nil {
		t.Fatalf("Failed to send GET request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	// Check if the middleware works by checking the response body
	expectedBody := "Hello, World!"
	body := make([]byte, len(expectedBody))
	_, err = resp.Body.Read(body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	if string(body) != expectedBody {
		t.Errorf("Expected body %s, got %s", expectedBody, string(body))
	}
}

func TestContextKeyString(t *testing.T) {
	keyName := "testKey"
	key := &contextKey{name: keyName}

	expectedString := "chi/middleware context value " + keyName
	if key.String() != expectedString {
		t.Errorf("Expected context key string %s, got %s", expectedString, key.String())
	}
}