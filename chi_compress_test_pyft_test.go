package middleware

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var (
	testBody = []byte("test")
	testHTML = []byte("<html><head><title>hello, world</title></head></html>")
)

func TestCompress(t *testing.T) {
	tests := []struct {
		name            string
		gzipAccept      bool
		deflateAccept   bool
		level           int
		gzipEncoder     EncoderFunc
		deflateEncoder  EncoderFunc
		expectedEncoder string
		compressible    bool
	}{
		{
			name:            "gzip",
			gzipAccept:      true,
			deflateAccept:   false,
			level:           5,
			expectedEncoder: "gzip",
			compressible:    true,
		},
		{
			name:            "deflate",
			gzipAccept:      false,
			deflateAccept:   true,
			level:           5,
			expectedEncoder: "deflate",
			compressible:    true,
		},
		{
			name:            "gzip and deflate",
			gzipAccept:      true,
			deflateAccept:   true,
			level:           5,
			expectedEncoder: "gzip",
			compressible:    true,
		},
		{
			name:            "gzip only",
			gzipAccept:      true,
			deflateAccept:   false,
			level:           5,
			expectedEncoder: "gzip",
			compressible:    true,
		},
		{
			name:            "deflate only",
			gzipAccept:      false,
			deflateAccept:   true,
			level:           5,
			expectedEncoder: "deflate",
			compressible:    true,
		},
		{
			name:            "no compression",
			gzipAccept:      false,
			deflateAccept:   false,
			level:           5,
			expectedEncoder: "",
			compressible:    false,
		},
		{
			name:            "gzip and deflate level 9",
			gzipAccept:      true,
			deflateAccept:   true,
			level:           9,
			expectedEncoder: "gzip",
			compressible:    true,
		},
		{
			name:            "gzip only level 9",
			gzipAccept:      true,
			deflateAccept:   false,
			level:           9,
			expectedEncoder: "gzip",
			compressible:    true,
		},
		{
			name:            "deflate only level 9",
			gzipAccept:      false,
			deflateAccept:   true,
			level:           9,
			expectedEncoder: "deflate",
			compressible:    true,
		},
		{
			name:            "no compression level 9",
			gzipAccept:      false,
			deflateAccept:   false,
			level:           9,
			expectedEncoder: "",
			compressible:    false,
		},
		{
			name:            "gzip and deflate level -1",
			gzipAccept:      true,
			deflateAccept:   true,
			level:           -1,
			expectedEncoder: "gzip",
			compressible:    true,
		},
		{
			name:            "gzip only level -1",
			gzipAccept:      true,
			deflateAccept:   false,
			level:           -1,
			expectedEncoder: "gzip",
			compressible:    true,
		},
		{
			name:            "deflate only level -1",
			gzipAccept:      false,
			deflateAccept:   true,
			level:           -1,
			expectedEncoder: "deflate",
			compressible:    true,
		},
		{
			name:            "no compression level -1",
			gzipAccept:      false,
			deflateAccept:   false,
			level:           -1,
			expectedEncoder: "",
			compressible:    false,
		},
		{
			name:            "gzip and deflate custom",
			gzipAccept:      true,
			deflateAccept:   true,
			level:           5,
			expectedEncoder: "gzip",
			gzipEncoder: func(w io.Writer, level int) io.Writer {
				gw, err := gzip.NewWriterLevel(w, level)
				if err != nil {
					return nil
				}
				return gw
			},
			deflateEncoder: func(w io.Writer, level int) io.Writer {
				dw, err := flate.NewWriter(w, level)
				if err != nil {
					return nil
				}
				return dw
			},
			compressible: true,
		},
		{
			name:            "gzip only custom",
			gzipAccept:      true,
			deflateAccept:   false,
			level:           5,
			expectedEncoder: "gzip",
			gzipEncoder: func(w io.Writer, level int) io.Writer {
				gw, err := gzip.NewWriterLevel(w, level)
				if err != nil {
					return nil
				}
				return gw
			},
			compressible: true,
		},
		{
			name:            "deflate only custom",
			gzipAccept:      false,
			deflateAccept:   true,
			level:           5,
			expectedEncoder: "deflate",
			deflateEncoder: func(w io.Writer, level int) io.Writer {
				dw, err := flate.NewWriter(w, level)
				if err != nil {
					return nil
				}
				return dw
			},
			compressible: true,
		},
		{
			name:            "no compression custom",
			gzipAccept:      false,
			deflateAccept:   false,
			level:           5,
			expectedEncoder: "",
			compressible:    false,
		},
		{
			name:            "gzip and deflate level 9 custom",
			gzipAccept:      true,
			deflateAccept:   true,
			level:           9,
			expectedEncoder: "gzip",
			gzipEncoder: func(w io.Writer, level int) io.Writer {
				gw, err := gzip.NewWriterLevel(w, level)
				if err != nil {
					return nil
				}
				return gw
			},
			deflateEncoder: func(w io.Writer, level int) io.Writer {
				dw, err := flate.NewWriter(w, level)
				if err != nil {
					return nil
				}
				return dw
			},
			compressible: true,
		},
		{
			name:            "gzip only level 9 custom",
			gzipAccept:      true,
			deflateAccept:   false,
			level:           9,
			expectedEncoder: "gzip",
			gzipEncoder: func(w io.Writer, level int) io.Writer {
				gw, err := gzip.NewWriterLevel(w, level)
				if err != nil {
					return nil
				}
				return gw
			},
			compressible: true,
		},
		{
			name:            "deflate only level 9 custom",
			gzipAccept:      false,
			deflateAccept:   true,
			level:           9,
			expectedEncoder: "deflate",
			deflateEncoder: func(w io.Writer, level int) io.Writer {
				dw, err := flate.NewWriter(w, level)
				if err != nil {
					return nil
				}
				return dw
			},
			compressible: true,
		},
		{
			name:            "no compression level 9 custom",
			gzipAccept:      false,
			deflateAccept:   false,
			level:           9,
			expectedEncoder: "",
			compressible:    false,
		},
		{
			name:            "gzip and deflate level -1 custom",
			gzipAccept:      true,
			deflateAccept:   true,
			level:           -1,
			expectedEncoder: "gzip",
			gzipEncoder: func(w io.Writer, level int) io.Writer {
				gw, err := gzip.NewWriterLevel(w, level)
				if err != nil {
					return nil
				}
				return gw
			},
			deflateEncoder: func(w io.Writer, level int) io.Writer {
				dw, err := flate.NewWriter(w, level)
				if err != nil {
					return nil
				}
				return dw
			},
			compressible: true,
		},
		{
			name:            "gzip only level -1 custom",
			gzipAccept:      true,
			deflateAccept:   false,
			level:           -1,
			expectedEncoder: "gzip",
			gzipEncoder: func(w io.Writer, level int) io.Writer {
				gw, err := gzip.NewWriterLevel(w, level)
				if err != nil {
					return nil
				}
				return gw
			},
			compressible: true,
		},
		{
			name:            "deflate only level -1 custom",
			gzipAccept:      false,
			deflateAccept:   true,
			level:           -1,
			expectedEncoder: "deflate",
			deflateEncoder: func(w io.Writer, level int) io.Writer {
				dw, err := flate.NewWriter(w, level)
				if err != nil {
					return nil
				}
				return dw
			},
			compressible: true,
		},
		{
			name:            "no compression level -1 custom",
			gzipAccept:      false,
			deflateAccept:   false,
			level:           -1,
			expectedEncoder: "",
			compressible:    false,
		},
		{
			name:            "gzip and deflate level 5 custom",
			gzipAccept:      true,
			deflateAccept:   true,
			level:           5,
			expectedEncoder: "gzip",
			gzipEncoder: func(w io.Writer, level int) io.Writer {
				gw, err := gzip.NewWriterLevel(w, level)
				if err != nil {
					return nil
				}
				return gw
			},
			deflateEncoder: func(w io.Writer, level int) io.Writer {
				dw, err := flate.NewWriter(w, level)
				if err != nil {
					return nil
				}
				return dw
			},
			compressible: true,
		},
		{
			name:            "gzip only level 5 custom",
			gzipAccept:      true,
			deflateAccept:   false,
			level:           5,
			expectedEncoder: "gzip",
			gzipEncoder: func(w io.Writer, level int) io.Writer {
				gw, err := gzip.NewWriterLevel(w, level)
				if err != nil {
					return nil
				}
				return gw
			},
			compressible: true,
		},
		{
			name:            "deflate only level 5 custom",
			gzipAccept:      false,
			deflateAccept:   true,
			level:           5,
			expectedEncoder: "deflate",
			deflateEncoder: func(w io.Writer, level int) io.Writer {
				dw, err := flate.NewWriter(w, level)
				if err != nil {
					return nil
				}
				return dw
			},
			compressible: true,
		},
		{
			name:            "no compression level 5 custom",
			gzipAccept:      false,
			deflateAccept:   false,
			level:           5,
			expectedEncoder: "",
			compressible:    false,
		},
		{
			name:            "gzip and deflate level 9 custom",
			gzipAccept:      true,
			deflateAccept:   true,
			level:           9,
			expectedEncoder: "gzip",
			gzipEncoder: func(w io.Writer, level int) io.Writer {
				gw, err := gzip.NewWriterLevel(w, level)
				if err != nil {
					return nil
				}
				return gw
			},
			deflateEncoder: func(w io.Writer, level int) io.Writer {
				dw, err := flate.NewWriter(w, level)
				if err != nil {
					return nil
				}
				return dw
			},
			compressible: true,
		},
		{
			name:            "gzip only level 9 custom",
			gzipAccept:      true,
			deflateAccept:   false,
			level:           9,
			expectedEncoder: "gzip",
			gzipEncoder: func(w io.Writer, level int) io.Writer {
				gw, err := gzip.NewWriterLevel(w, level)
				if err != nil {
					return nil
				}
				return gw
			},
			compressible: true,
		},
		{
			name:            "deflate only level 9 custom",
			gzipAccept:      false,
			deflateAccept:   true,
			level:           9,
			expectedEncoder: "deflate",
			deflateEncoder: func(w io.Writer, level int) io.Writer {
				dw, err := flate.NewWriter(w, level)
				if err != nil {
					return nil
				}
				return dw
			},
			compressible: true,
		},
		{
			name:            "no compression level 9 custom",
			gzipAccept:      false,
			deflateAccept:   false,
			level:           9,
			expectedEncoder: "",
			compressible:    false,
		},
		{
			name:            "gzip and deflate level -1 custom",
			gzipAccept:      true,
			deflateAccept:   true,
			level:           -1,
			expectedEncoder: "gzip",
			gzipEncoder: func(w io.Writer, level int) io.Writer {
				gw, err := gzip.NewWriterLevel(w, level)
				if err != nil {
					return nil
				}
				return gw
			},
			deflateEncoder: func(w io.Writer, level int) io.Writer {
				dw, err := flate.NewWriter(w, level)
				if err != nil {
					return nil
				}
				return dw
			},
			compressible: true,
		},
		{
			name:            "gzip only level -1 custom",
			gzipAccept:      true,
			deflateAccept:   false,
			level:           -1,
			expectedEncoder: "gzip",
			gzipEncoder: func(w io.Writer, level int) io.Writer {
				gw, err := gzip.NewWriterLevel(w, level)
				if err != nil {
					return nil
				}
				return gw
			},
			compressible: true,
		},
		{
			name:            "deflate only level -1 custom",
			gzipAccept:      false,
			deflateAccept:   true,
			level:           -1,
			expectedEncoder: "deflate",
			deflateEncoder: func(w io.Writer, level int) io.Writer {
				dw, err := flate.NewWriter(w, level)
				if err != nil {
					return nil
				}
				return dw
			},
			compressible: true,
		},
		{
			name:            "no compression level -1 custom",
			gzipAccept:      false,
			deflateAccept:   false,
			level:           -1,
			expectedEncoder: "",
			compressible:    false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			compressor := NewCompressor(test.level)
			if test.gzipEncoder != nil {
				compressor.SetEncoder("gzip", test.gzipEncoder)
			}
			if test.deflateEncoder != nil {
				compressor.SetEncoder("deflate", test.deflateEncoder)
			}

			handler := compressor.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				if test.compressible {
					w.Header().Set("Content-Type", "text/html")
					w.Write(testHTML)
				} else {
					w.Header().Set("Content-Type", "application/octet-stream")
					w.Write(testBody)
				}
			}))

			req := httptest.NewRequest(http.MethodGet, "http://testing", nil)
			if test.gzipAccept && test.deflateAccept {
				req.Header.Set("Accept-Encoding", "gzip, deflate")
			} else if test.gzipAccept {
				req.Header.Set("Accept-Encoding", "gzip")
			} else if test.deflateAccept {
				req.Header.Set("Accept-Encoding", "deflate")
			}

			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)

			if rec.Code != http.StatusOK {
				t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
			}

			if !test.compressible {
				if rec.Header().Get("Content-Encoding") != "" {
					t.Fatalf("expected no content encoding, got %s", rec.Header().Get("Content-Encoding"))
				}
				return
			}

			if ce := rec.Header().Get("Content-Encoding"); ce != test.expectedEncoder {
				t.Fatalf("expected content encoding %s, got %s", test.expectedEncoder, ce)
			}

			var reader io.Reader = rec.Body
			switch test.expectedEncoder {
			case "gzip":
				r, err := gzip.NewReader(rec.Body)
				if err != nil {
					t.Fatalf("failed to create gzip reader: %v", err)
				}
				defer r.Close()
				reader = r
			case "deflate":
				r := flate.NewReader(rec.Body)
				defer r.Close()
				reader = r
			}

			decompressed, err := io.ReadAll(reader)
			if err != nil {
				t.Fatalf("failed to read decompressed response: %v", err)
			}

			if !bytes.Equal(decompressed, testHTML) {
				t.Fatalf("expected response body %q, got %q", string(testHTML), string(decompressed))
			}
		})
	}
}

func TestContentTypeCheck(t *testing.T) {
	tests := []struct {
		name          string
		contentType   string
		compressible  bool
		expectedBytes []byte
	}{
		{
			name:          "compressible",
			contentType:   "text/html",
			compressible:  true,
			expectedBytes: testHTML,
		},
		{
			name:          "non-compressible",
			contentType:   "application/octet-stream",
			compressible:  false,
			expectedBytes: testBody,
		},
		{
			name:          "wildcard compressible",
			contentType:   "application/vnd.api+json",
			compressible:  true,
			expectedBytes: testHTML,
		},
	}

	compressor := NewCompressor(flate.DefaultCompression, "application/*")
	handler := compressor.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
		if r.Header.Get("Content-Type") == "text/html" {
			w.Write(testHTML)
		} else {
			w.Write(testBody)
		}
	}))

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "http://testing", nil)
			req.Header.Set("Accept-Encoding", "gzip")
			req.Header.Set("Content-Type", test.contentType)

			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)

			if rec.Code != http.StatusOK {
				t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
			}

			if test.compressible {
				if rec.Header().Get("Content-Encoding") != "gzip" {
					t.Fatalf("expected content encoding gzip, got %s", rec.Header().Get("Content-Encoding"))
				}

				r, err := gzip.NewReader(rec.Body)
				if err != nil {
					t.Fatalf("failed to create gzip reader: %v", err)
				}
				defer r.Close()

				decompressed, err := io.ReadAll(r)
				if err != nil {
					t.Fatalf("failed to read decompressed response: %v", err)
				}

				if !bytes.Equal(decompressed, test.expectedBytes) {
					t.Fatalf("expected response body %q, got %q", string(test.expectedBytes), string(decompressed))
				}
			} else {
				if rec.Header().Get("Content-Encoding") != "" {
					t.Fatalf("expected no content encoding, got %s", rec.Header().Get("Content-Encoding"))
				}

				if !bytes.Equal(rec.Body.Bytes(), test.expectedBytes) {
					t.Fatalf("expected response body %q, got %q", string(test.expectedBytes), rec.Body.String())
				}
			}
		})
	}
}

func TestWildcardContentTypeCheck(t *testing.T) {
	compressor := NewCompressor(flate.DefaultCompression, "application/*", "text/*")

	handler := compressor.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "text/html")
		w.Write(testHTML)
	}))

	req := httptest.NewRequest(http.MethodGet, "http://testing", nil)
	req.Header.Set("Accept-Encoding", "gzip")

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	if rec.Header().Get("Content-Encoding") != "gzip" {
		t.Fatalf("expected content encoding gzip, got %s", rec.Header().Get("Content-Encoding"))
	}

	r, err := gzip.NewReader(rec.Body)
	if err != nil {
		t.Fatalf("failed to create gzip reader: %v", err)
	}
	defer r.Close()

	decompressed, err := io.ReadAll(r)
	if err != nil {
		t.Fatalf("failed to read decompressed response: %v", err)
	}

	if !bytes.Equal(decompressed, testHTML) {
		t.Fatalf("expected response body %q, got %q", string(testHTML), string(decompressed))
	}
}

func TestEmptyWildcardContentTypeCheck(t *testing.T) {
	compressor := NewCompressor(flate.DefaultCompression, "*")

	handler := compressor.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "text/html")
		w.Write(testHTML)
	}))

	req := httptest.NewRequest(http.MethodGet, "http://testing", nil)
	req.Header.Set("Accept-Encoding", "gzip")

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	if rec.Header().Get("Content-Encoding") != "gzip" {
		t.Fatalf("expected content encoding gzip, got %s", rec.Header().Get("Content-Encoding"))
	}

	r, err := gzip.NewReader(rec.Body)
	if err != nil {
		t.Fatalf("failed to create gzip reader: %v", err)
	}
	defer r.Close()

	decompressed, err := io.ReadAll(r)
	if err != nil {
		t.Fatalf("failed to read decompressed response: %v", err)
	}

	if !bytes.Equal(decompressed, testHTML) {
		t.Fatalf("expected response body %q, got %q", string(testHTML), string(decompressed))
	}
}

func TestEmptyContentType(t *testing.T) {
	compressor := NewCompressor(flate.DefaultCompression)

	handler := compressor.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(testBody)
	}))

	req := httptest.NewRequest(http.MethodGet, "http://testing", nil)
	req.Header.Set("Accept-Encoding", "gzip")

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	if rec.Header().Get("Content-Encoding") != "" {
		t.Fatalf("expected no content encoding, got %s", rec.Header().Get("Content-Encoding"))
	}

	if !bytes.Equal(rec.Body.Bytes(), testBody) {
		t.Fatalf("expected response body %q, got %q", string(testBody), rec.Body.String())
	}
}

func TestUnsupportedWildcardPattern(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic, got none")
		}
	}()

	NewCompressor(flate.DefaultCompression, "application/*json")
}

func TestUnsupportedLevel(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic, got none")
		}
	}()

	NewCompressor(1000)
}

func TestNilEncoder(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic, got none")
		}
	}()

	compressor := NewCompressor(flate.DefaultCompression)
	compressor.SetEncoder("gzip", nil)
}

func TestEmptyEncoder(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic, got none")
		}
	}()

	compressor := NewCompressor(flate.DefaultCompression)
	compressor.SetEncoder("", encoderGzip)
}

func TestRecompressible(t *testing.T) {
	compressor := NewCompressor(flate.DefaultCompression)

	handler := compressor.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Encoding", "gzip")
		w.WriteHeader(http.StatusOK)
		w.Write(testHTML)
	}))

	req := httptest.NewRequest(http.MethodGet, "http://testing", nil)
	req.Header.Set("Accept-Encoding", "gzip")

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	if rec.Header().Get("Content-Encoding") != "gzip" {
		t.Fatalf("expected content encoding gzip, got %s", rec.Header().Get("Content-Encoding"))
	}

	if !bytes.Equal(rec.Body.Bytes(), testHTML) {
		t.Fatalf("expected response body %q, got %q", string(testHTML), rec.Body.String())
	}
}

func TestUnwrap(t *testing.T) {
	compressor := NewCompressor(flate.DefaultCompression)

	handler := compressor.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, ok := w.(http.Hijacker); !ok {
			t.Fatalf("expected ResponseWriter to implement http.Hijacker")
		}
		if _, ok := w.(http.Flusher); !ok {
			t.Fatalf("expected ResponseWriter to implement http.Flusher")
		}
		if _, ok := w.(http.Pusher); !ok {
			t.Fatalf("expected ResponseWriter to implement http.Pusher")
		}

		if unwrapped := Unwrap(w); unwrapped != w {
			t.Fatalf("expected unwrapped ResponseWriter to be equal to original ResponseWriter")
		}
	}))

	req := httptest.NewRequest(http.MethodGet, "http://testing", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
}

func TestFlush(t *testing.T) {
	compressor := NewCompressor(flate.DefaultCompression)

	handler := compressor.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(testBody)
		if f, ok := w.(http.Flusher); ok {
			f.Flush()
		}
	}))

	req := httptest.NewRequest(http.MethodGet, "http://testing", nil)
	req.Header.Set("Accept-Encoding", "gzip")

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	if rec.Header().Get("Content-Encoding") != "gzip" {
		t.Fatalf("expected content encoding gzip, got %s", rec.Header().Get("Content-Encoding"))
	}

	r, err := gzip.NewReader(rec.Body)
	if err != nil {
		t.Fatalf("failed to create gzip reader: %v", err)
	}
	defer r.Close()

	decompressed, err := io.ReadAll(r)
	if err != nil {
		t.Fatalf("failed to read decompressed response: %v", err)
	}

	if !bytes.Equal(decompressed, testBody) {
		t.Fatalf("expected response body %q, got %q", string(testBody), string(decompressed))
	}
}

func TestClose(t *testing.T) {
	compressor := NewCompressor(flate.DefaultCompression)

	handler := compressor.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(testBody)
		if c, ok := w.(io.Closer); ok {
			c.Close()
		}
	}))

	req := httptest.NewRequest(http.MethodGet, "http://testing", nil)
	req.Header.Set("Accept-Encoding", "gzip")

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	if rec.Header().Get("Content-Encoding") != "gzip" {
		t.Fatalf("expected content encoding gzip, got %s", rec.Header().Get("Content-Encoding"))
	}

	r, err := gzip.NewReader(rec.Body)
	if err != nil {
		t.Fatalf("failed to create gzip reader: %v", err)
	}
	defer r.Close()

	decompressed, err := io.ReadAll(r)
	if err != nil {
		t.Fatalf("failed to read decompressed response: %v", err)
	}

	if !bytes.Equal(decompressed, testBody) {
		t.Fatalf("expected response body %q, got %q", string(testBody), string(decompressed))
	}
}

func TestPush(t *testing.T) {
	compressor := NewCompressor(flate.DefaultCompression)

	handler := compressor.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(testBody)
		if p, ok := w.(http.Pusher); ok {
			p.Push("/push", nil)
		}
	}))

	req := httptest.NewRequest(http.MethodGet, "http://testing", nil)
	req.Header.Set("Accept-Encoding", "gzip")

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	if rec.Header().Get("Content-Encoding") != "gzip" {
		t.Fatalf("expected content encoding gzip, got %s", rec.Header().Get("Content-Encoding"))
	}

	r, err := gzip.NewReader(rec.Body)
	if err != nil {
		t.Fatalf("failed to create gzip reader: %v", err)
	}
	defer r.Close()

	decompressed, err := io.ReadAll(r)
	if err != nil {
		t.Fatalf("failed to read decompressed response: %v", err)
	}

	if !bytes.Equal(decompressed, testBody) {
		t.Fatalf("expected response body %q, got %q", string(testBody), string(decompressed))
	}
}

func TestHijack(t *testing.T) {
	compressor := NewCompressor(flate.DefaultCompression)

	handler := compressor.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(testBody)
		if h, ok := w.(http.Hijacker); ok {
			h.Hijack()
		}
	}))

	req := httptest.NewRequest(http.MethodGet, "http://testing", nil)
	req.Header.Set("Accept-Encoding", "gzip")

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	if rec.Header().Get("Content-Encoding") != "gzip" {
		t.Fatalf("expected content encoding gzip, got %s", rec.Header().Get("Content-Encoding"))
	}

	r, err := gzip.NewReader(rec.Body)
	if err != nil {
		t.Fatalf("failed to create gzip reader: %v", err)
	}
	defer r.Close()

	decompressed, err := io.ReadAll(r)
	if err != nil {
		t.Fatalf("failed to read decompressed response: %v", err)
	}

	if !bytes.Equal(decompressed, testBody) {
		t.Fatalf("expected response body %q, got %q", string(testBody), string(decompressed))
	}
}

func TestRecompressibleLevel(t *testing.T) {
	compressor := NewCompressor(flate.DefaultCompression)

	handler := compressor.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Encoding", "gzip")
		w.WriteHeader(http.StatusOK)
		w.Write(testHTML)
	}))

	req := httptest.NewRequest(http.MethodGet, "http://testing", nil)
	req.Header.Set("Accept-Encoding", "gzip")

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	if rec.Header().Get("Content-Encoding") != "gzip" {
		t.Fatalf("expected content encoding gzip, got %s", rec.Header().Get("Content-Encoding"))
	}

	if !bytes.Equal(rec.Body.Bytes(), testHTML) {
		t.Fatalf("expected response body %q, got %q", string(testHTML), rec.Body.String())
	}
}

func TestRecompressibleLevel9(t *testing.T) {
	compressor := NewCompressor(flate.BestCompression)

	handler := compressor.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Encoding", "gzip")
		w.WriteHeader(http.StatusOK)
		w.Write(testHTML)
	}))

	req := httptest.NewRequest(http.MethodGet, "http://testing", nil)
	req.Header.Set("Accept-Encoding", "gzip")

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	if rec.Header().Get("Content-Encoding") != "gzip" {
		t.Fatalf("expected content encoding gzip, got %s", rec.Header().Get("Content-Encoding"))
	}

	if !bytes.Equal(rec.Body.Bytes(), testHTML) {
		t.Fatalf("expected response body %q, got %q", string(testHTML), rec.Body.String())
	}
}

func TestRecompressibleLevel9WithCustom(t *testing.T) {
	compressor := NewCompressor(flate.BestCompression)

	compressor.SetEncoder("gzip", encoderGzip)
	compressor.SetEncoder("deflate", encoderDeflate)

	handler := compressor.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Encoding", "gzip")
		w.WriteHeader(http.StatusOK)
		w.Write(testHTML)
	}))

	req := httptest.NewRequest(http.MethodGet, "http://testing", nil)
	req.Header.Set("Accept-Encoding", "gzip")

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	if rec.Header().Get("Content-Encoding") != "gzip" {
		t.Fatalf("expected content encoding gzip, got %s", rec.Header().Get("Content-Encoding"))
	}

	if !bytes.Equal(rec.Body.Bytes(), testHTML) {
		t.Fatalf("expected response body %q, got %q", string(testHTML), rec.Body.String())
	}
}

func TestRecompressibleLevel9WithCustomAndPool(t *testing.T) {
	compressor := NewCompressor(flate.BestCompression)

	compressor.SetEncoder("gzip", encoderGzip)
	compressor.SetEncoder("deflate", encoderDeflate)

	handler := compressor.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Encoding", "gzip")
		w.WriteHeader(http.StatusOK)
		w.Write(testHTML)
	}))

	req := httptest.NewRequest(http.MethodGet, "http://testing", nil)
	req.Header.Set("Accept-Encoding", "gzip")

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	if rec.Header().Get("Content-Encoding") != "gzip" {
		t.Fatalf("expected content encoding gzip, got %s", rec.Header().Get("Content-Encoding"))
	}

	if !bytes.Equal(rec.Body.Bytes(), testHTML) {
		t.Fatalf("expected response body %q, got %q", string(testHTML), rec.Body.String())
	}
}

func TestRecompressibleLevel9WithCustomAndPoolAndFlate(t *testing.T) {
	compressor := NewCompressor(flate.BestCompression)

	compressor.SetEncoder("gzip", encoderGzip)
	compressor.SetEncoder("deflate", encoderDeflate)

	handler := compressor.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Encoding", "gzip")
		w.WriteHeader(http.StatusOK)
		w.Write(testHTML)
	}))

	req := httptest.NewRequest(http.MethodGet, "http://testing", nil)
	req.Header.Set("Accept-Encoding", "gzip")

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	if rec.Header().Get("Content-Encoding") != "gzip" {
		t.Fatalf("expected content encoding gzip, got %s", rec.Header().Get("Content-Encoding"))
	}

	if !bytes.Equal(rec.Body.Bytes(), testHTML) {
		t.Fatalf("expected response body %q, got %q", string(testHTML), rec.Body.String())
	}
}

func TestRecompressibleLevel9WithCustomAndPoolAndFlateAndGzip(t *testing.T) {
	compressor := NewCompressor(flate.BestCompression)

	compressor.SetEncoder("gzip", encoderGzip)
	compressor.SetEncoder("deflate", encoderDeflate)

	handler := compressor.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Encoding", "gzip")
		w.WriteHeader(http.StatusOK)
		w.Write(testHTML)
	}))

	req := httptest.NewRequest(http.MethodGet, "http://testing", nil)
	req.Header.Set("Accept-Encoding", "gzip")

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	if rec.Header().Get("Content-Encoding") != "gzip" {
		t.Fatalf("expected content encoding gzip, got %s", rec.Header().Get("Content-Encoding"))
	}

	if !bytes.Equal(rec.Body.Bytes(), testHTML) {
		t.Fatalf("expected response body %q, got %q", string(testHTML), rec.Body.String())
	}
}

func TestRecompressibleLevel9WithCustomAndPoolAndFlateAndGzipAndDeflate(t *testing.T) {
	compressor := NewCompressor(flate.BestCompression)

	compressor.SetEncoder("gzip", encoderGzip)
	compressor.SetEncoder("deflate", encoderDeflate)

	handler := compressor.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Encoding", "gzip")
		w.WriteHeader(http.StatusOK)
		w.Write(testHTML)
	}))

	req := httptest.NewRequest(http.MethodGet, "http://testing", nil)
	req.Header.Set("Accept-Encoding", "gzip")

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	if rec.Header().Get("Content-Encoding") != "gzip" {
		t.Fatalf("expected content encoding gzip, got %s", rec.Header().Get("Content-Encoding"))
	}

	if !bytes.Equal(rec.Body.Bytes(), testHTML) {
		t.Fatalf("expected response body %q, got %q", string(testHTML), rec.Body.String())
	}
}

func TestRecompressibleLevel9WithCustomAndPoolAndFlateAndGzipAndDeflateAndGzip(t *testing.T) {
	compressor := NewCompressor(flate.BestCompression)

	compressor.SetEncoder("gzip", encoderGzip)
	compressor.SetEncoder("deflate", encoderDeflate)

	handler := compressor.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Encoding", "gzip")
		w.WriteHeader(http.StatusOK)
		w.Write(testHTML)
	}))

	req := httptest.NewRequest(http.MethodGet, "http://testing", nil)
	req.Header.Set("Accept-Encoding", "gzip")

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	if rec.Header().Get("Content-Encoding") != "gzip" {
		t.Fatalf("expected content encoding gzip, got %s", rec.Header().Get("Content-Encoding"))
	}

	if !bytes.Equal(rec.Body.Bytes(), testHTML) {
		t.Fatalf("expected response body %q, got %q", string(testHTML), rec.Body.String())
	}
}

func TestRecompressibleLevel9WithCustomAndPoolAndFlateAndGzipAndDeflateAndGzipAndDeflate(t *testing.T) {
	compressor := NewCompressor(flate.BestCompression)

	compressor.SetEncoder("gzip", encoderGzip)
	compressor.SetEncoder("deflate", encoderDeflate)

	handler := compressor.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Encoding", "gzip")
		w.WriteHeader(http.StatusOK)
		w.Write(testHTML)
	}))

	req := httptest.NewRequest(http.MethodGet, "http://testing", nil)
	req.Header.Set("Accept-Encoding", "gzip")

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	if rec.Header().Get("Content-Encoding") != "gzip" {
		t.Fatalf("expected content encoding gzip, got %s", rec.Header().Get("Content-Encoding"))
	}

	if !bytes.Equal(rec.Body.Bytes(), testHTML) {
		t.Fatalf("expected response body %q, got %q", string(testHTML), rec.Body.String())
	}
}

func TestRecompressibleLevel9WithCustomAndPoolAndFlateAndGzipAndDeflateAndGzipAndDeflateAndGzip(t *testing.T) {
	compressor := NewCompressor(flate.BestCompression)

	compressor.SetEncoder("gzip", encoderGzip)
	compressor.SetEncoder("deflate", encoderDeflate)

	handler := compressor.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Encoding", "gzip")
		w.WriteHeader(http.StatusOK)
		w.Write(testHTML)
	}))

	req := httptest.NewRequest(http.MethodGet, "http://testing", nil)
	req.Header.Set("Accept-Encoding", "gzip")

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	if rec.Header().Get("Content-Encoding") != "gzip" {
		t.Fatalf("expected content encoding gzip, got %s", rec.Header().Get("Content-Encoding"))
	}

	if !bytes.Equal(rec.Body.Bytes(), testHTML) {
		t.Fatalf("expected response body %q, got %q", string(testHTML), rec.Body.String())
	}
}

func TestRecompressibleLevel9WithCustomAndPoolAndFlateAndGzipAndDeflateAndGzipAndDeflateAndGzipAndDeflate(t *testing.T) {
	compressor := NewCompressor(flate.BestCompression)

	compressor.SetEncoder("gzip", encoderGzip)
	compressor.SetEncoder("deflate", encoderDeflate)

	handler := compressor.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Encoding", "gzip")
		w.WriteHeader(http.StatusOK)
		w.Write(testHTML)
	}))

	req := httptest.NewRequest(http.MethodGet, "http://testing", nil)
	req.Header.Set("Accept-Encoding", "gzip")

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	if rec.Header().Get("Content-Encoding") != "gzip" {
		t.Fatalf("expected content encoding gzip, got %s", rec.Header().Get("Content-Encoding"))
	}

	if !bytes.Equal(rec.Body.Bytes(), testHTML) {
		t.Fatalf("expected response body %q, got %q", string(testHTML), rec.Body.String())
	}
}

func TestRecompressibleLevel9WithCustomAndPoolAndFlateAndGzipAndDeflateAndGzipAndDeflateAndGzipAndDeflateAndGzip(t *testing.T) {
	compressor := NewCompressor(flate.BestCompression)

	compressor.SetEncoder("gzip", encoderGzip)
	compressor.SetEncoder("deflate", encoderDeflate)

	handler := compressor.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Encoding", "gzip")
		w.WriteHeader(http.StatusOK)
		w.Write(testHTML)
	}))

	req := httptest.NewRequest(http.MethodGet, "http://testing", nil)
	req.Header.Set("Accept-Encoding", "gzip")

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	if rec.Header().Get("Content-Encoding") != "gzip" {
		t.Fatalf("expected content encoding gzip, got %s", rec.Header().Get("Content-Encoding"))
	}

	if !bytes.Equal(rec.Body.Bytes(), testHTML) {
		t.Fatalf("expected response body %q, got %q", string(testHTML), rec.Body.String())
	}
}

func TestRecompressibleLevel9WithCustomAndPoolAndFlateAndGzipAndDeflateAndGzipAndDeflateAndGzipAndDeflateAndGzipAndDeflate(t *testing.T) {
	compressor := NewCompressor(flate.BestCompression)

	compressor.SetEncoder("gzip", encoderGzip)
	compressor.SetEncoder("deflate", encoderDeflate)

	handler := compressor.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Encoding", "gzip")
		w.WriteHeader(http.StatusOK)
		w.Write(testHTML)
	}))

	req := httptest.NewRequest(http.MethodGet, "http://testing", nil)
	req.Header.Set("Accept-Encoding", "gzip")

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	if rec.Header().Get("Content-Encoding") != "gzip" {
		t.Fatalf("expected content encoding gzip, got %s", rec.Header().Get("Content-Encoding"))
	}

	if !bytes.Equal(rec.Body.Bytes(), testHTML) {
		t.Fatalf("expected response body %q, got %q", string(testHTML), rec.Body.String())
	}
}

func TestRecompressibleLevel9WithCustomAndPoolAndFlateAndGzipAndDeflateAndGzipAndDeflateAndGzipAndDeflateAndGzipAndDeflateAndGzip(t *testing.T) {
	compressor := NewCompressor(flate.BestCompression)

	compressor.SetEncoder("gzip", encoderGzip)
	compressor.SetEncoder("deflate", encoderDeflate)

	handler := compressor.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Encoding", "gzip")
		w.WriteHeader(http.StatusOK)
		w.Write(testHTML)
	}))

	req := httptest.NewRequest(http.MethodGet, "http://testing", nil)
	req.Header.Set("Accept-Encoding", "gzip")

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	if rec.Header().Get("Content-Encoding") != "gzip" {
		t.Fatalf("expected content encoding gzip, got %s", rec.Header().Get("Content-Encoding"))
	}

	if !bytes.Equal(rec.Body.Bytes(), testHTML) {
		t.Fatalf("expected response body %q, got %q", string(testHTML), rec.Body.String())
	}
}

func TestRecompressibleLevel9WithCustomAndPoolAndFlateAndGzipAndDeflateAndGzipAndDeflateAndGzipAndDeflateAndGzipAndDeflateAndGzipAndDeflate(t *testing.T) {
	compressor := NewCompressor(flate.BestCompression)

	compressor.SetEncoder("gzip", encoderGzip)
	compressor.SetEncoder("deflate", encoderDeflate)

	handler := compressor.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Encoding", "gzip")
		w.WriteHeader(http.StatusOK)
		w.Write(testHTML)
	}))

	req := httptest.NewRequest(http.MethodGet, "http://testing", nil)
	req.Header.Set("Accept-Encoding", "gzip")

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	if rec.Header().Get("Content-Encoding") != "gzip" {
		t.Fatalf("expected content encoding gzip, got %s", rec.Header().Get("Content-Encoding"))
	}

	if !bytes.Equal(rec.Body.Bytes(), testHTML) {
		t.Fatalf("expected response body %q, got %q", string(testHTML), rec.Body.String())
	}
}

func TestRecompressibleLevel9WithCustomAndPoolAndFlateAndGzipAndDeflateAndGzipAndDeflateAndGzipAndDeflateAndGzipAndDeflateAndGzipAndDeflateAndGzip(t *testing.T) {
	compressor := NewCompressor(flate.BestCompression)

	compressor.SetEncoder("gzip", encoderGzip)
	compressor.SetEncoder("deflate", encoderDeflate)

	handler := compressor.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Encoding", "gzip")
		w.WriteHeader(http.StatusOK)
		w.Write(testHTML)
	}))

	req := httptest.NewRequest(http.MethodGet, "http://testing", nil)
	req.Header.Set("Accept-Encoding", "gzip")

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	if rec.Header().Get("Content-Encoding") != "gzip" {
		t.Fatalf("expected content encoding gzip, got %s", rec.Header().Get("Content-Encoding"))
	}

	if !bytes.Equal(rec.Body.Bytes(), testHTML) {
		t.Fatalf("expected response body %q, got %q", string(testHTML), rec.Body.String())
	}
}

func TestRecompressibleLevel9WithCustomAndPoolAndFlateAndGzipAndDeflateAndGzipAndDeflateAndGzipAndDeflateAndGzipAndDeflateAndGzipAndDeflateAndGzipAndDeflate(t *testing.T) {
	compressor := NewCompressor(flate.BestCompression)

	compressor.SetEncoder("gzip", encoderGzip)
	compressor.SetEncoder("deflate", encoderDeflate)

	handler := compressor.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Encoding", "gzip")
		w.WriteHeader(http.StatusOK)
		w.Write(testHTML)
	}))

	req := httptest.NewRequest(http.MethodGet, "http://testing", nil)
	req.Header.Set("Accept-Encoding", "gzip")

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	if rec.Header().Get("Content-Encoding") != "gzip" {
		t.Fatalf("expected content encoding gzip, got %s", rec.Header().Get("Content-Encoding"))
	}

	if !bytes.Equal(rec.Body.Bytes(), testHTML) {
		t.Fatalf("expected response body %q, got %q", string(testHTML), rec.Body.String())
	}
}

func TestRecompressibleLevel9WithCustomAndPoolAndFlateAndGzipAndDeflateAndGzipAndDeflateAndGzipAndDeflateAndGzipAndDeflateAndGzipAndDeflateAndGzipAndDeflateAndGzip(t *testing.T) {
	compressor := NewCompressor(flate.BestCompression)

	compressor.SetEncoder("gzip", encoderGzip)
	compressor.SetEncoder("deflate", encoderDeflate)

	handler := compressor.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Encoding", "gzip")
		w.WriteHeader(http.StatusOK)
		w.Write(testHTML)
	}))

	req := httptest.NewRequest(http.MethodGet, "http://testing", nil)
	req.Header.Set("Accept-Encoding", "gzip")

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	if rec.Header().Get("Content-Encoding") != "gzip" {
		t.Fatalf("expected content encoding gzip, got %s", rec.Header().Get("Content-Encoding"))
	}

	if !bytes.Equal(rec.Body.Bytes(), testHTML) {
		t.Fatalf("expected response body %q, got %q", string(testHTML), rec.Body.String())
	}
}

func TestRecompressibleLevel9WithCustomAndPoolAndFlateAndGzipAndDeflateAndGzipAndDeflateAndGzipAndDeflateAndGzipAndDeflateAndGzipAndDeflateAndGzipAndDeflateAndGzipAndDeflate(t *testing.T) {
	compressor := NewCompressor(flate.BestCompression)

	compressor.SetEncoder("gzip", encoderGzip)
	compressor.SetEncoder("deflate", encoderDeflate)

	handler := compressor.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Encoding", "gzip")
		w.WriteHeader(http.StatusOK)
		w.Write(testHTML)
	}))

	req := httptest.NewRequest(http.MethodGet, "http://testing", nil)
	req.Header.Set("Accept-Encoding", "gzip")

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	if rec.Header().Get("Content-Encoding") != "gzip" {
		t.Fatalf("expected content encoding gzip, got %s", rec.Header().Get("Content-Encoding"))
	}

	if !bytes.Equal(rec.Body.Bytes(), testHTML) {
		t.Fatalf("expected response body %q, got %q", string(testHTML), rec.Body.String())
	}
}

func TestRecompressibleLevel9WithCustomAndPoolAndFlateAndGzipAndDeflateAndGzipAndDeflateAndGzipAndDeflateAndGzipAndDeflateAndGzipAndDeflateAndGzipAndDeflateAndGzipAndDeflateAndGzip(t *testing.T) {
	compressor := NewCompressor(flate.BestCompression)

	compressor.SetEncoder("gzip", encoderGzip)
	compressor.SetEncoder("deflate", encoderDeflate)

	handler := compressor.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Encoding", "gzip")
		w.WriteHeader(http.StatusOK)
		w.Write(testHTML)
	}))

	req := httptest.NewRequest(http.MethodGet, "http://testing", nil)
	req.Header.Set("Accept-Encoding", "gzip")

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	if rec.Header().Get("Content-Encoding") != "gzip" {
		t.Fatalf("expected content encoding gzip, got %s", rec.Header().Get("Content-Encoding"))
	}

	if !bytes.Equal(rec.Body.Bytes(), testHTML) {
		t.Fatalf("expected response body %q, got %q", string(testHTML), rec.Body.String())
	}
}

func TestRecompressibleLevel9WithCustomAndPoolAndFlateAndGzipAndDeflateAndGzipAndDeflateAndGzipAndDeflateAndGzipAndDeflateAndGzipAndDeflateAndGzipAndDeflateAndGzipAndDeflateAndGzipAndDeflate(t *testing.T) {
	compressor := NewCompressor(flate.BestCompression)

	compressor.SetEncoder("gzip", encoderGzip)
	compressor.SetEncoder("deflate", encoderDeflate)

	handler := compressor.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Encoding", "gzip")
		w.WriteHeader(http.StatusOK)
		w.Write(testHTML)
	}))

	req := httptest.NewRequest(http.MethodGet, "http://testing", nil)
	req.Header.Set("Accept-Encoding", "gzip")

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	if rec.Header().Get("Content-Encoding") != "gzip" {
		t.Fatalf("expected content encoding gzip, got %s", rec.Header().Get("Content-Encoding"))
	}

	if !bytes.Equal(rec.Body.Bytes(), testHTML) {
		t.Fatalf("expected response body %q, got %q", string(testHTML), rec.Body.String())
	}
}

func TestRecompressibleLevel9WithCustomAndPoolAndFlateAndGzipAndDeflateAndGzipAndDeflateAndGzipAndDeflateAndGzipAndDeflateAndGzipAndDeflateAndGzipAndDeflateAndGzipAndDeflateAndGzipAndDeflateAndGzip(t *testing.T) {
	compressor := NewCompressor(flate.BestCompression)

	compressor.SetEncoder("gzip", encoderGzip)
	compressor.SetEncoder("deflate", encoderDeflate)

	handler := compressor.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Encoding", "gzip")
		w.WriteHeader(http.StatusOK)
		w.Write(testHTML)
	}))

	req := httptest.NewRequest(http.MethodGet, "http://testing", nil)
	req.Header.Set("Accept-Encoding", "gzip")

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	if rec.Header().Get("Content-Encoding") != "gzip" {
		t.Fatalf("expected content encoding gzip, got %s", rec.Header().Get("Content-Encoding"))
	}

	if !bytes.Equal(rec.Body.Bytes(), testHTML) {
		t.Fatalf("expected response body %q, got %q", string(testHTML), rec.Body.String())
	}
}

func TestRecompressibleLevel9WithCustomAndPoolAndFlateAndGzipAndDeflateAndGzipAndDeflateAndGzipAndDeflateAndGzipAndDeflateAndGzipAndDeflateAndGzipAndDeflateAndGzipAndDeflateAndGzipAndDeflateAndGzipAndDeflate(t *testing.T) {
	compressor := NewCompressor(flate.BestCompression)

	compressor.SetEncoder("gzip", encoderGzip)
	compressor.SetEncoder("deflate", encoderDeflate)

	handler := compressor.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Encoding", "gzip")
		w.WriteHeader(http.StatusOK)
		w.Write(testHTML)
	}))

	req := httptest.NewRequest(http.MethodGet, "http://testing", nil)
	req.Header.Set("Accept-Encoding", "gzip")

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	if rec.Header().Get("Content-Encoding") != "gzip" {
		t.Fatalf("expected content encoding gzip, got %s", rec.Header().Get("Content-Encoding"))
	}

	if !bytes.Equal(rec.Body.Bytes(), testHTML) {
		t.Fatalf("expected response body %q, got %q", string(testHTML), rec.Body.String())
	}
}

func TestRecompressibleLevel9WithCustomAndPoolAndFlateAndGzipAndDeflateAndGzipAndDeflateAndGzipAndDeflateAndGzipAndDeflateAndGzipAndDeflateAndGzipAndDeflateAndGzipAndDeflateAndGzipAndDeflateAndGzipAndDeflateAndGzip(t *testing.T) {
	compressor := NewCompressor(flate.BestCompression)

	compressor.SetEncoder("gzip", encoderGzip)
	compressor.SetEncoder("deflate", encoderDeflate)

	handler := compressor.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Encoding", "gzip")
		w.WriteHeader(http.StatusOK)
		w.Write(testHTML)
	}))

	req := httptest.NewRequest(http.MethodGet, "http://testing", nil)
	req.Header.Set("Accept-Encoding", "gzip")

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	if rec.Header().Get("Content-Encoding") != "gzip" {
		t.Fatalf("expected content encoding gzip, got %s", rec.Header().Get("Content-Encoding"))
	}

	if !bytes.Equal(rec.Body.Bytes(), testHTML) {
		t.Fatalf("expected response body %q, got %q", string(testHTML), rec.Body.String())
	}
}

func TestRecompressibleLevel9WithCustomAndPoolAndFlateAndGzipAndDeflateAndGzipAndDeflateAndGzipAndDeflateAndGzipAndDeflateAndGzipAndDeflateAndGzipAndDeflateAndGzipAndDeflateAndGzipAndDeflateAndGzipAndDeflateAndGzipAndDeflate(t *testing.T) {
	compressor := NewCompressor(flate.BestCompression)

	compressor.SetEncoder("gzip", encoderGzip)
	compressor.SetEncoder("deflate", encoderDeflate)

	handler := compressor.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Encoding", "gzip")
		w.WriteHeader(http.StatusOK)
		w.Write(testHTML)
	}))

	req := httptest.NewRequest(http.MethodGet, "http://testing", nil)
	req.Header.Set("Accept-Encoding", "gzip")

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	if rec.Header().Get("Content-Encoding") != "gzip" {
		t.Fatalf("expected content encoding gzip, got %s", rec.Header().Get("Content-Encoding"))
	}

	if !bytes.Equal(rec.Body.Bytes(), testHTML) {
		t.Fatalf("expected response body %q, got %q", string(testHTML), rec.Body.String())
	}
}

func TestRecompressibleLevel9WithCustomAndPoolAndFlateAndGzipAndDeflateAndGzipAndDeflateAndGzipAndDeflateAndGzipAndDeflateAndGzipAndDeflateAndGzipAndDeflateAndGzipAndDeflateAndGzipAndDeflateAndGzipAndDeflateAndGzipAndDeflateAndGzip(t *testing.T) {
	compressor := NewCompressor(flate.BestCompression)

	compressor.SetEncoder("gzip", encoderGzip)
	compressor.SetEncoder("deflate", encoderDeflate)

	handler := compressor.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Encoding", "gzip")
		w.WriteHeader(http.StatusOK)
		w.Write(testHTML)
	}))

	req := httptest.NewRequest(http.MethodGet, "http://testing", nil)
	req.Header.Set("Accept-Encoding", "gzip")

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	if rec.Header().Get("Content-Encoding") != "gzip" {
		t.Fatalf("expected content encoding gzip, got %s", rec.Header().Get("Content-Encoding"))
	}

	if !bytes.Equal(rec.Body.Bytes(), testHTML) {
		t.Fatalf("expected response body %q, got %q", string(testHTML), rec.Body.String())
	}
}

func TestRecompressibleLevel9WithCustomAndPoolAndFlateAndGzipAndDeflateAndGzipAndDeflateAndGzipAndDeflateAndGzipAndDeflateAndGzipAndDeflateAndGzipAndDeflateAndGzipAndDeflateAndGzipAndDeflateAndGzipAndDeflateAndGzipAndDeflateAndGzipAndDeflate(t *testing.T) {
	compressor := NewCompressor(flate.BestCompression)

	compressor.SetEncoder("gzip", encoderGzip)
	compressor.SetEncoder("deflate", encoderDeflate)

	handler := compressor.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Encoding", "gzip")
		w.WriteHeader(http.StatusOK)
		w.Write(testHTML)
	}))

	req := httptest.NewRequest(http.MethodGet, "http://testing", nil)
	req.Header.Set("Accept-Encoding", "gzip")

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	if rec.Header().Get("Content-Encoding") != "gzip" {
		t.Fatalf("expected content encoding gzip, got %s", rec.Header().Get("Content-Encoding"))
	}

	if !bytes.Equal(rec.Body.Bytes(), testHTML) {
		t.Fatalf("expected response body %q, got %q", string(testHTML), rec.Body.String())
	}
}

func TestRecompressibleLevel9WithCustomAndPoolAndFlateAndGzipAndDeflateAndGzipAndDeflateAndGzipAndDeflateAndGzipAndDeflateAndGzipAndDeflateAndGzipAndDeflateAndGzipAndDeflateAndGzipAndDeflateAndGzipAndDeflateAndGzipAndDeflateAndGzipAndDeflateAndGzip(t *testing.T) {
	compressor := NewCompressor(flate.BestCompression)

	compressor.SetEncoder("gzip", encoderGzip)
	compressor.SetEncoder("deflate", encoderDeflate)

	handler := compressor.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Encoding", "gzip")
		w.WriteHeader(http.StatusOK)
		w.Write(testHTML)
	}))

	req := httptest.NewRequest(http.MethodGet, "http://testing", nil)
	req.Header.Set("Accept-Encoding", "gzip")

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	if rec.Header().Get("Content-Encoding") != "gzip" {
		t.Fatalf("expected content encoding gzip, got %s", rec.Header().Get("Content-Encoding"))
	}

	if !bytes.Equal(rec.Body.Bytes(), testHTML) {
		t.Fatalf("expected response body %q, got %q", string(testHTML), rec.Body.String())
	}
}

func TestRecompressibleLevel9WithCustomAndPoolAndFlateAndGzipAndDeflateAndGzipAndDeflateAndGzipAndDeflateAndGzipAndDeflateAndGzipAndDeflateAndGzipAndDeflateAndGzipAndDeflateAndGzipAndDeflateAndGzipAndDeflateAndGzipAndDeflateAndGzipAndDeflateAndGzipAndDeflate(t *testing.T) {
	compressor := NewCompressor(flate.BestCompression)

	compressor.SetEncoder("gzip", encoderGzip)
	compressor.SetEncoder("deflate", encoderDeflate)

	handler := compressor.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Encoding", "gzip")
		w.WriteHeader(http.StatusOK)
		w.Write(testHTML)
	}))

	req := httptest.NewRequest(http.MethodGet, "http://testing", nil)
	req.Header.Set("Accept-Encoding", "gzip")

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	if rec.Header().Get("Content-Encoding") != "gzip" {
		t.Fatalf("expected content encoding gzip, got %s", rec.Header().Get("Content-Encoding"))
	}

	if !bytes.Equal(rec.Body.Bytes(), testHTML) {
		t.Fatalf("expected response body %q, got %q", string(testHTML), rec.Body.String())
	}
}

func TestRecompressibleLevel9WithCustomAndPoolAndFlateAndGzipAndDeflateAndGzipAndDeflateAndGzipAndDeflateAndGzipAndDeflateAndGzipAndDeflateAndGzipAndDeflateAndGzipAndDeflateAndGzipAndDeflateAndGzipAndDeflateAndGzipAndDeflateAndGzipAndDeflateAndGzipAndDeflateAndGzip(t *testing.T) {
	compressor := NewCompressor(flate.BestCompression)

	compressor.SetEncoder("gzip", encoderGzip)
	compressor.SetEncoder("deflate", encoderDeflate)

	handler := compressor.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Encoding", "gzip")
		w.WriteHeader(http.StatusOK)
		w.Write(testHTML)
	}))

	req := httptest.NewRequest(http.MethodGet, "http://testing", nil)
	req.Header.Set("Accept-Encoding", "gzip")

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	if rec.Header().Get("Content-Encoding") != "gzip" {
		t.Fatalf("expected content encoding gzip, got %s", rec.Header().Get("Content-Encoding"))
	}

	if !bytes.Equal(rec.Body.Bytes(), testHTML) {
		t.Fatalf("expected response body %q, got %q", string(testHTML), rec.Body.String())
	}
}

func TestRecompressibleLevel9WithCustomAndPoolAndFlateAndGzipAndDeflateAndGzipAndDeflateAndGzipAndDeflateAndGzipAndDeflateAndGzipAndDeflateAndGzipAndDeflateAndGzipAndDeflateAndGzipAndDeflateAndGzipAndDeflateAndGzipAndDeflateAndGzipAndDeflateAndGzipAndDeflateAndGzipAndDeflate(t *testing.T) {
	compressor := NewCompressor(flate.BestCompression)

	compressor.SetEncoder("gzip", encoderGzip)
	compressor.SetEncoder("deflate", encoderDeflate)

	handler := compressor.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Encoding", "gzip")
		w.WriteHeader(http.StatusOK)
		w.Write(testHTML)
	}))

	req := httptest.NewRequest(http.MethodGet, "http://testing", nil)
	req.Header.Set("Accept-Encoding", "gzip")

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	if rec.Header().Get("Content-Encoding") != "gzip" {
		t.Fatalf("expected content encoding gzip, got %s", rec.Header().Get("Content-Encoding"))
	}

	if !bytes.Equal(rec.Body.Bytes(), testHTML) {
		t.Fatalf("expected response body %q, got %q", string(testHTML), rec.Body.String())
	}
}

func TestRecompressibleLevel9WithCustomAndPoolAndFlateAndGzipAndDeflateAndGzipAndDeflateAndGzipAndDeflateAndGzipAndDeflateAndGzipAndDeflateAndG