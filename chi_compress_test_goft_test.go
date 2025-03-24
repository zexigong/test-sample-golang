package middleware

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func Test_Compress_deflate(t *testing.T) {
	testCompress(t, "deflate")
}

func Test_Compress_gzip(t *testing.T) {
	testCompress(t, "gzip")
}

func testCompress(t *testing.T, encoding string) {
	// A message larger than 1500 bytes will trigger writes > the default io.Copy
	// buffer size of 32 bytes. This message should be larger than the io.Copy
	// buffer size but still compressible.
	message := strings.Repeat("0", 1500)
	// Store the compressed message for comparison later.
	msgBuf := new(bytes.Buffer)
	var err error
	switch encoding {
	case "deflate":
		gzw, _ := flate.NewWriter(msgBuf, flate.DefaultCompression)
		_, err = gzw.Write([]byte(message))
		if err != nil {
			t.Fatal(err)
		}
		err = gzw.Close()
		if err != nil {
			t.Fatal(err)
		}
	case "gzip":
		gzw, _ := gzip.NewWriterLevel(msgBuf, gzip.DefaultCompression)
		_, err = gzw.Write([]byte(message))
		if err != nil {
			t.Fatal(err)
		}
		err = gzw.Close()
		if err != nil {
			t.Fatal(err)
		}
	}

	// This message will be too small to be compressed.
	smallMsg := strings.Repeat("0", 10)

	// This message is already compressed.
	gzw, _ := gzip.NewWriterLevel(io.Discard, gzip.DefaultCompression)
	gzw.Write([]byte(message))
	gzw.Close()

	tests := []struct {
		encoding  string
		accept    string
		skipMsg   bool
		msg       string
		compressed bool
	}{
		{"deflate", "", false, message, false},
		{"gzip", "", false, message, false},
		{"deflate", "gzip", false, message, false},
		{"gzip", "deflate", false, message, false},
		{"deflate", "deflate", false, message, true},
		{"gzip", "gzip", false, message, true},
		{"deflate", "deflate, gzip", false, message, true},
		{"gzip", "deflate, gzip", false, message, true},
		{"deflate", "gzip, deflate", false, message, true},
		{"gzip", "gzip, deflate", false, message, true},
		{"deflate", "*", false, message, true},
		{"gzip", "*", false, message, true},
		{"deflate", "deflate", false, smallMsg, false},
		{"gzip", "gzip", false, smallMsg, false},
		{"deflate", "deflate", true, message, false},
		{"gzip", "gzip", true, message, false},
	}

	for _, tc := range tests {
		t.Run(tc.encoding+":"+tc.accept, func(t *testing.T) {
			compressor := NewCompressor(-1)
			compressor.SetEncoder(tc.encoding, func(w io.Writer, level int) io.Writer {
				switch tc.encoding {
				case "deflate":
					gzw, _ := flate.NewWriter(w, level)
					return gzw
				case "gzip":
					gzw, _ := gzip.NewWriterLevel(w, level)
					return gzw
				}
				panic("unexpected encoding: " + tc.encoding)
			})

			handler := compressor.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tc.skipMsg {
					w.Header().Set("Content-Encoding", tc.encoding)
				}
				w.Write([]byte(tc.msg))
			}))

			r, _ := http.NewRequest("GET", "/", nil)
			r.Header.Set("Accept-Encoding", tc.accept)

			w := httptest.NewRecorder()
			handler.ServeHTTP(w, r)

			if tc.compressed {
				if w.Header().Get("Content-Encoding") != tc.encoding {
					t.Errorf("expected Content-Encoding: %q, got: %q", tc.encoding, w.Header().Get("Content-Encoding"))
				}
				if !bytes.Equal(w.Body.Bytes(), msgBuf.Bytes()) {
					t.Errorf("expected compressed message to be %q, got %q", msgBuf.Bytes(), w.Body.Bytes())
				}
			} else {
				if w.Header().Get("Content-Encoding") == tc.encoding {
					t.Errorf("unexpected Content-Encoding: %q", tc.encoding)
				}
				if !bytes.Equal(w.Body.Bytes(), []byte(tc.msg)) {
					t.Errorf("expected uncompressed message to be %q, got %q", tc.msg, w.Body.Bytes())
				}
			}
		})
	}
}

func Test_UnsupportedWildcard(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected panic for unsupported wildcard")
		}
	}()
	NewCompressor(-1, "*")
	t.Fatal("should panic")
}

func Test_SupportsContentType(t *testing.T) {
	c := NewCompressor(5, "text/html")

	tests := []struct {
		contentType string
		shouldBeValid bool
	}{
		{"text/html", true},
		{"text/html; charset=utf-8", true},
		{"text/css", false},
		{"", false},
	}

	for _, test := range tests {
		t.Run(test.contentType, func(t *testing.T) {
			cw := &compressResponseWriter{
				ResponseWriter: httptest.NewRecorder(),
				contentTypes:   c.allowedTypes,
				encoding:       "gzip",
			}
			cw.Header().Set("Content-Type", test.contentType)

			if actual := cw.isCompressible(); actual != test.shouldBeValid {
				t.Errorf("expected %t, got %t", test.shouldBeValid, actual)
			}
		})
	}
}