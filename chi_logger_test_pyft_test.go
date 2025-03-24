package middleware

import (
	"bytes"
	"context"
	"log"
	"net/http"
	"os"
	"runtime"
	"testing"
	"time"
)

func TestLogger(t *testing.T) {
	// Test the Logger middleware function
	handler := Logger(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	req, _ := http.NewRequest("GET", "/", nil)
	rr := &responseRecorder{ResponseWriter: newNopResponseWriter(), status: http.StatusOK}
	handler.ServeHTTP(rr, req)

	if rr.status != http.StatusOK {
		t.Errorf("expected status %v, got %v", http.StatusOK, rr.status)
	}
}

func TestRequestLogger(t *testing.T) {
	// Test the RequestLogger function
	formatter := &DefaultLogFormatter{Logger: log.New(os.Stdout, "", log.LstdFlags), NoColor: true}
	handler := RequestLogger(formatter)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	req, _ := http.NewRequest("GET", "/", nil)
	rr := &responseRecorder{ResponseWriter: newNopResponseWriter(), status: http.StatusOK}
	handler.ServeHTTP(rr, req)

	if rr.status != http.StatusOK {
		t.Errorf("expected status %v, got %v", http.StatusOK, rr.status)
	}
}

func TestGetLogEntry(t *testing.T) {
	// Test the GetLogEntry function
	req, _ := http.NewRequest("GET", "/", nil)
	entry := &defaultLogEntry{
		request: req,
		buf:     &bytes.Buffer{},
	}
	ctx := context.WithValue(req.Context(), LogEntryCtxKey, entry)
	req = req.WithContext(ctx)

	gotEntry := GetLogEntry(req)
	if gotEntry != entry {
		t.Errorf("expected %v, got %v", entry, gotEntry)
	}
}

func TestWithLogEntry(t *testing.T) {
	// Test the WithLogEntry function
	req, _ := http.NewRequest("GET", "/", nil)
	entry := &defaultLogEntry{
		request: req,
		buf:     &bytes.Buffer{},
	}

	reqWithEntry := WithLogEntry(req, entry)
	gotEntry := GetLogEntry(reqWithEntry)
	if gotEntry != entry {
		t.Errorf("expected %v, got %v", entry, gotEntry)
	}
}

type responseRecorder struct {
	http.ResponseWriter
	status int
}

func (rr *responseRecorder) WriteHeader(status int) {
	rr.status = status
	rr.ResponseWriter.WriteHeader(status)
}

func newNopResponseWriter() http.ResponseWriter {
	return &nopResponseWriter{}
}

type nopResponseWriter struct{}

func (w *nopResponseWriter) Header() http.Header {
	return http.Header{}
}

func (w *nopResponseWriter) Write([]byte) (int, error) {
	return 0, nil
}

func (w *nopResponseWriter) WriteHeader(statusCode int) {}