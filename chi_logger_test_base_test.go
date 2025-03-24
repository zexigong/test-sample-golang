package middleware_test

import (
	"bytes"
	"context"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

type mockLogFormatter struct{}

func (f *mockLogFormatter) NewLogEntry(r *http.Request) middleware.LogEntry {
	return &mockLogEntry{}
}

type mockLogEntry struct {
	buf bytes.Buffer
}

func (e *mockLogEntry) Write(status, bytes int, header http.Header, elapsed time.Duration, extra interface{}) {
	e.buf.WriteString("mock write")
}

func (e *mockLogEntry) Panic(v interface{}, stack []byte) {
	e.buf.WriteString("mock panic")
}

type mockLogger struct {
	output *bytes.Buffer
}

func (ml *mockLogger) Print(v ...interface{}) {
	ml.output.WriteString(v[0].(string))
}

func TestLoggerMiddleware(t *testing.T) {
	mockLogBuffer := &bytes.Buffer{}
	mockLogger := &mockLogger{output: mockLogBuffer}
	middleware.DefaultLogger = middleware.RequestLogger(&middleware.DefaultLogFormatter{Logger: mockLogger, NoColor: true})

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, World!"))
	})

	req := httptest.NewRequest("GET", "http://example.com/foo", nil)
	w := httptest.NewRecorder()

	middleware.Logger(handler).ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	logOutput := mockLogBuffer.String()
	if logOutput == "" {
		t.Fatalf("expected log output, got empty string")
	}
}

func TestRequestLoggerWithCustomFormatter(t *testing.T) {
	mockLogBuffer := &bytes.Buffer{}
	mockLogger := &mockLogger{output: mockLogBuffer}
	formatter := &mockLogFormatter{}
	logger := middleware.RequestLogger(formatter)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, World!"))
	})

	req := httptest.NewRequest("GET", "http://example.com/foo", nil)
	w := httptest.NewRecorder()

	logger(handler).ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	// Check if mock formatter's Write method was called
	expected := "mock write"
	if !bytes.Contains(mockLogBuffer.Bytes(), []byte(expected)) {
		t.Fatalf("expected log output to contain %q, got %q", expected, mockLogBuffer.String())
	}
}

func TestGetLogEntry(t *testing.T) {
	req := httptest.NewRequest("GET", "http://example.com/foo", nil)
	entry := &mockLogEntry{}
	req = middleware.WithLogEntry(req, entry)

	logEntry := middleware.GetLogEntry(req)
	if logEntry != entry {
		t.Fatalf("expected log entry to be %v, got %v", entry, logEntry)
	}
}