package middleware

import (
	"bytes"
	"context"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type testLogger struct {
	lastLine string
}

func (l *testLogger) Print(v ...interface{}) {
	l.lastLine = v[0].(string)
}

func TestLogger_Panic(t *testing.T) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, ReqIDCtxKey, "testID")
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, "/test", nil)

	logger := &testLogger{}
	formatter := &DefaultLogFormatter{
		Logger:  logger,
		NoColor: true,
	}
	entry := formatter.NewLogEntry(req)

	t.Run("panic and stacktrace", func(t *testing.T) {
		entry.Panic("test", []byte("stacktrace"))
		assert.Equal(t, "stacktrace", logger.lastLine)
	})
}

func TestLogger(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	req.RemoteAddr = "addr:123"
	req.Header.Add("X-Request-Id", "request-id")

	t.Run("default logger", func(t *testing.T) {
		rec := httptest.NewRecorder()
		logger := &testLogger{}
		handler := RequestLogger(&DefaultLogFormatter{
			Logger:  logger,
			NoColor: true,
		})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(10 * time.Millisecond)
		}))
		handler.ServeHTTP(rec, req)
		assert.Contains(t, logger.lastLine, "[request-id] \"GET http://example.com/test HTTP/1.1\" from addr:123 - 200 0B in")
	})

	t.Run("logger without request-id", func(t *testing.T) {
		rec := httptest.NewRecorder()
		logger := &testLogger{}
		handler := RequestLogger(&DefaultLogFormatter{
			Logger:  logger,
			NoColor: true,
		})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(10 * time.Millisecond)
		}))
		handler.ServeHTTP(rec, req)
		assert.Contains(t, logger.lastLine, "200 0B in")
	})

	t.Run("logger with header", func(t *testing.T) {
		rec := httptest.NewRecorder()
		logger := &testLogger{}
		handler := RequestLogger(&DefaultLogFormatter{
			Logger:  logger,
			NoColor: true,
		})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(10 * time.Millisecond)
			w.Header().Add("X-Test", "value")
		}))
		handler.ServeHTTP(rec, req)
		assert.Contains(t, logger.lastLine, "200 0B in")
	})

	t.Run("logger with body", func(t *testing.T) {
		rec := httptest.NewRecorder()
		logger := &testLogger{}
		handler := RequestLogger(&DefaultLogFormatter{
			Logger:  logger,
			NoColor: true,
		})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(10 * time.Millisecond)
			w.Write([]byte("test"))
		}))
		handler.ServeHTTP(rec, req)
		assert.Contains(t, logger.lastLine, "200 4B in")
	})
}

func TestGetLogEntry(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	req.RemoteAddr = "addr:123"

	t.Run("default logger", func(t *testing.T) {
		rec := httptest.NewRecorder()
		var entry LogEntry
		logger := &testLogger{}
		handler := RequestLogger(&DefaultLogFormatter{
			Logger:  logger,
			NoColor: true,
		})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			entry = GetLogEntry(r)
		}))
		handler.ServeHTTP(rec, req)
		assert.NotNil(t, entry)
	})

	t.Run("no logger", func(t *testing.T) {
		entry := GetLogEntry(req)
		assert.Nil(t, entry)
	})
}

func TestWithLogEntry(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	req.RemoteAddr = "addr:123"

	t.Run("default logger", func(t *testing.T) {
		rec := httptest.NewRecorder()
		var entry LogEntry
		logger := &testLogger{}
		handler := RequestLogger(&DefaultLogFormatter{
			Logger:  logger,
			NoColor: true,
		})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			entry = GetLogEntry(r)
			r = WithLogEntry(r, entry)
			entry = GetLogEntry(r)
		}))
		handler.ServeHTTP(rec, req)
		assert.NotNil(t, entry)
	})
}

func TestDefaultLogFormatter_NewLogEntry(t *testing.T) {
	formatter := &DefaultLogFormatter{
		Logger:  log.New(bytes.NewBuffer(nil), "", 0),
		NoColor: true,
	}
	req := httptest.NewRequest("GET", "http://example.com/foo", nil)
	entry := formatter.NewLogEntry(req)

	assert.NotNil(t, entry)
}