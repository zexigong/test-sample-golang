package gin

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestRecoveryMiddleware(t *testing.T) {
	router := New()
	router.Use(Recovery())
	router.GET("/panic", func(c *Context) {
		panic("test panic")
	})

	req, _ := http.NewRequest("GET", "/panic", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("Expected status code %d, got %d", http.StatusInternalServerError, w.Code)
	}
}

func TestCustomRecoveryMiddleware(t *testing.T) {
	var recoveredErr any
	recoveryFunc := func(c *Context, err any) {
		recoveredErr = err
		c.AbortWithStatus(http.StatusTeapot)
	}

	router := New()
	router.Use(CustomRecovery(recoveryFunc))
	router.GET("/panic", func(c *Context) {
		panic("custom panic")
	})

	req, _ := http.NewRequest("GET", "/panic", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusTeapot {
		t.Fatalf("Expected status code %d, got %d", http.StatusTeapot, w.Code)
	}

	if recoveredErr != "custom panic" {
		t.Fatalf("Expected recovered error 'custom panic', got '%v'", recoveredErr)
	}
}

func TestRecoveryWithWriter(t *testing.T) {
	var buf bytes.Buffer
	router := New()
	router.Use(RecoveryWithWriter(&buf))
	router.GET("/panic", func(c *Context) {
		panic("writer panic")
	})

	req, _ := http.NewRequest("GET", "/panic", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("Expected status code %d, got %d", http.StatusInternalServerError, w.Code)
	}

	if !bytes.Contains(buf.Bytes(), []byte("writer panic")) {
		t.Fatalf("Expected log to contain 'writer panic', got '%s'", buf.String())
	}
}

func TestCustomRecoveryWithWriter(t *testing.T) {
	var buf bytes.Buffer
	var recoveredErr any
	recoveryFunc := func(c *Context, err any) {
		recoveredErr = err
		c.AbortWithStatus(http.StatusServiceUnavailable)
	}

	router := New()
	router.Use(CustomRecoveryWithWriter(&buf, recoveryFunc))
	router.GET("/panic", func(c *Context) {
		panic("custom writer panic")
	})

	req, _ := http.NewRequest("GET", "/panic", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusServiceUnavailable {
		t.Fatalf("Expected status code %d, got %d", http.StatusServiceUnavailable, w.Code)
	}

	if recoveredErr != "custom writer panic" {
		t.Fatalf("Expected recovered error 'custom writer panic', got '%v'", recoveredErr)
	}

	if !bytes.Contains(buf.Bytes(), []byte("custom writer panic")) {
		t.Fatalf("Expected log to contain 'custom writer panic', got '%s'", buf.String())
	}
}

func TestBrokenPipeRecovery(t *testing.T) {
	var buf bytes.Buffer
	router := New()
	router.Use(RecoveryWithWriter(&buf))
	router.GET("/panic", func(c *Context) {
		panic(&net.OpError{Err: os.NewSyscallError("write", errors.New("broken pipe"))})
	})

	req, _ := http.NewRequest("GET", "/panic", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("Expected status code %d, got %d", http.StatusInternalServerError, w.Code)
	}

	if bytes.Contains(buf.Bytes(), []byte("broken pipe")) {
		t.Fatalf("Expected log not to contain 'broken pipe', got '%s'", buf.String())
	}
}