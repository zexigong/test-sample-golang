package middleware

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCompressMiddleware(t *testing.T) {
	tests := []struct {
		name               string
		acceptEncoding     string
		expectedEncoding   string
		expectedStatusCode int
	}{
		{"gzip", "gzip", "gzip", http.StatusOK},
		{"deflate", "deflate", "deflate", http.StatusOK},
		{"unsupported", "br", "", http.StatusOK},
		{"no encoding", "", "", http.StatusOK},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "text/plain")
				w.Write([]byte("Hello, World!"))
			})

			compressor := Compress(flate.DefaultCompression)
			ts := httptest.NewServer(compressor(handler))
			defer ts.Close()

			req, _ := http.NewRequest("GET", ts.URL, nil)
			req.Header.Set("Accept-Encoding", tt.acceptEncoding)

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatalf("failed to perform request: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tt.expectedStatusCode {
				t.Errorf("expected status code %d, got %d", tt.expectedStatusCode, resp.StatusCode)
			}

			if tt.expectedEncoding != "" {
				if resp.Header.Get("Content-Encoding") != tt.expectedEncoding {
					t.Errorf("expected encoding %s, got %s", tt.expectedEncoding, resp.Header.Get("Content-Encoding"))
				}

				// Verify decompressed content
				var reader io.Reader
				if tt.expectedEncoding == "gzip" {
					reader, err = gzip.NewReader(resp.Body)
					if err != nil {
						t.Fatalf("failed to create gzip reader: %v", err)
					}
				} else if tt.expectedEncoding == "deflate" {
					reader = flate.NewReader(resp.Body)
				}

				body, err := ioutil.ReadAll(reader)
				if err != nil {
					t.Fatalf("failed to read response body: %v", err)
				}

				if string(body) != "Hello, World!" {
					t.Errorf("expected body %s, got %s", "Hello, World!", string(body))
				}
			} else {
				body, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					t.Fatalf("failed to read response body: %v", err)
				}
				if string(body) != "Hello, World!" {
					t.Errorf("expected body %s, got %s", "Hello, World!", string(body))
				}
			}
		})
	}
}

func TestNewCompressor(t *testing.T) {
	compressor := NewCompressor(flate.DefaultCompression, "text/plain", "application/json")

	if _, ok := compressor.allowedTypes["text/plain"]; !ok {
		t.Errorf("expected text/plain to be allowed")
	}

	if _, ok := compressor.allowedTypes["application/json"]; !ok {
		t.Errorf("expected application/json to be allowed")
	}

	if _, ok := compressor.allowedTypes["text/html"]; ok {
		t.Errorf("expected text/html not to be allowed")
	}
}

func TestSetEncoder(t *testing.T) {
	compressor := NewCompressor(flate.DefaultCompression)

	compressor.SetEncoder("custom", func(w io.Writer, level int) io.Writer {
		return &nopCloser{w}
	})

	if _, ok := compressor.encoders["custom"]; !ok {
		t.Errorf("expected custom encoder to be set")
	}
}

type nopCloser struct {
	io.Writer
}

func (nopCloser) Close() error { return nil }