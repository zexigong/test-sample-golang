package render

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

func TestReader_Render(t *testing.T) {
	tests := []struct {
		name           string
		contentType    string
		contentLength  int64
		readerContent  string
		headers        map[string]string
		expectedOutput string
		expectedHeader map[string]string
	}{
		{
			name:           "Basic rendering",
			contentType:    "text/plain",
			contentLength:  13,
			readerContent:  "Hello, World!",
			headers:        nil,
			expectedOutput: "Hello, World!",
			expectedHeader: map[string]string{
				"Content-Type":   "text/plain",
				"Content-Length": "13",
			},
		},
		{
			name:           "With custom header",
			contentType:    "text/html",
			contentLength:  5,
			readerContent:  "Hello",
			headers:        map[string]string{"X-Custom-Header": "CustomValue"},
			expectedOutput: "Hello",
			expectedHeader: map[string]string{
				"Content-Type":    "text/html",
				"Content-Length":  "5",
				"X-Custom-Header": "CustomValue",
			},
		},
		{
			name:           "With negative content length",
			contentType:    "application/json",
			contentLength:  -1,
			readerContent:  "{\"key\":\"value\"}",
			headers:        nil,
			expectedOutput: "{\"key\":\"value\"}",
			expectedHeader: map[string]string{
				"Content-Type": "application/json",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			reader := Reader{
				ContentType:   tt.contentType,
				ContentLength: tt.contentLength,
				Reader:        strings.NewReader(tt.readerContent),
				Headers:       tt.headers,
			}

			err := reader.Render(w)
			if err != nil {
				t.Errorf("Render() error = %v", err)
				return
			}

			result := w.Body.String()
			if result != tt.expectedOutput {
				t.Errorf("Render() = %v, want %v", result, tt.expectedOutput)
			}

			for key, expectedValue := range tt.expectedHeader {
				if value := w.Header().Get(key); value != expectedValue {
					t.Errorf("Header %v = %v, want %v", key, value, expectedValue)
				}
			}
		})
	}
}

func TestReader_WriteContentType(t *testing.T) {
	w := httptest.NewRecorder()
	reader := Reader{
		ContentType: "text/plain",
	}

	reader.WriteContentType(w)

	if contentType := w.Header().Get("Content-Type"); contentType != "text/plain" {
		t.Errorf("WriteContentType() = %v, want %v", contentType, "text/plain")
	}
}

func TestReader_writeHeaders(t *testing.T) {
	w := httptest.NewRecorder()
	headers := map[string]string{
		"Content-Type": "application/json",
		"X-Test":       "test-value",
	}

	reader := Reader{}
	reader.writeHeaders(w, headers)

	for key, expectedValue := range headers {
		if value := w.Header().Get(key); value != expectedValue {
			t.Errorf("writeHeaders() %v = %v, want %v", key, value, expectedValue)
		}
	}
}