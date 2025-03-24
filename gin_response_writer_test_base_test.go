package gin

import (
	"bufio"
	"bytes"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockHijacker struct {
	http.ResponseWriter
}

func (m *mockHijacker) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return nil, nil, nil
}

type mockFlusher struct {
	http.ResponseWriter
}

func (m *mockFlusher) Flush() {}

type mockCloseNotifier struct {
	http.ResponseWriter
	closeChan chan bool
}

func (m *mockCloseNotifier) CloseNotify() <-chan bool {
	return m.closeChan
}

type mockPusher struct {
	http.ResponseWriter
}

func (m *mockPusher) Push(target string, opts *http.PushOptions) error {
	return nil
}

func TestResponseWriterWriteHeader(t *testing.T) {
	rec := httptest.NewRecorder()
	writer := &responseWriter{ResponseWriter: rec}

	writer.WriteHeader(http.StatusNotFound)
	if writer.status != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, writer.status)
	}

	writer.WriteHeader(http.StatusInternalServerError)
	if writer.status != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, writer.status)
	}
}

func TestResponseWriterWriteHeaderNow(t *testing.T) {
	rec := httptest.NewRecorder()
	writer := &responseWriter{ResponseWriter: rec}

	writer.WriteHeaderNow()
	if rec.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rec.Code)
	}
}

func TestResponseWriterWrite(t *testing.T) {
	rec := httptest.NewRecorder()
	writer := &responseWriter{ResponseWriter: rec}

	data := []byte("Hello")
	n, err := writer.Write(data)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if n != len(data) {
		t.Errorf("Expected written length %d, got %d", len(data), n)
	}
	if rec.Body.String() != "Hello" {
		t.Errorf("Expected body %s, got %s", "Hello", rec.Body.String())
	}
}

func TestResponseWriterWriteString(t *testing.T) {
	rec := httptest.NewRecorder()
	writer := &responseWriter{ResponseWriter: rec}

	s := "Hello"
	n, err := writer.WriteString(s)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if n != len(s) {
		t.Errorf("Expected written length %d, got %d", len(s), n)
	}
	if rec.Body.String() != "Hello" {
		t.Errorf("Expected body %s, got %s", "Hello", rec.Body.String())
	}
}

func TestResponseWriterStatus(t *testing.T) {
	rec := httptest.NewRecorder()
	writer := &responseWriter{ResponseWriter: rec}

	if writer.Status() != http.StatusOK {
		t.Errorf("Expected default status %d, got %d", http.StatusOK, writer.Status())
	}

	writer.WriteHeader(http.StatusNotFound)
	if writer.Status() != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, writer.Status())
	}
}

func TestResponseWriterSize(t *testing.T) {
	rec := httptest.NewRecorder()
	writer := &responseWriter{ResponseWriter: rec}

	if writer.Size() != noWritten {
		t.Errorf("Expected size %d, got %d", noWritten, writer.Size())
	}

	writer.Write([]byte("Hello"))
	if writer.Size() != 5 {
		t.Errorf("Expected size %d, got %d", 5, writer.Size())
	}
}

func TestResponseWriterWritten(t *testing.T) {
	rec := httptest.NewRecorder()
	writer := &responseWriter{ResponseWriter: rec}

	if writer.Written() {
		t.Errorf("Expected not written, got written")
	}

	writer.Write([]byte("Hello"))
	if !writer.Written() {
		t.Errorf("Expected written, got not written")
	}
}

func TestResponseWriterHijack(t *testing.T) {
	rec := httptest.NewRecorder()
	hijacker := &mockHijacker{ResponseWriter: rec}
	writer := &responseWriter{ResponseWriter: hijacker}

	_, _, err := writer.Hijack()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
}

func TestResponseWriterCloseNotify(t *testing.T) {
	rec := httptest.NewRecorder()
	closeNotifier := &mockCloseNotifier{ResponseWriter: rec, closeChan: make(chan bool)}
	writer := &responseWriter{ResponseWriter: closeNotifier}

	if writer.CloseNotify() == nil {
		t.Errorf("Expected non-nil close notification channel")
	}
}

func TestResponseWriterFlush(t *testing.T) {
	rec := httptest.NewRecorder()
	flusher := &mockFlusher{ResponseWriter: rec}
	writer := &responseWriter{ResponseWriter: flusher}

	writer.Flush() // Should not panic
}

func TestResponseWriterPusher(t *testing.T) {
	rec := httptest.NewRecorder()
	pusher := &mockPusher{ResponseWriter: rec}
	writer := &responseWriter{ResponseWriter: pusher}

	if writer.Pusher() == nil {
		t.Errorf("Expected non-nil pusher")
	}
}