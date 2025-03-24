package fasthttp

import (
	"bytes"
	"io"
	"testing"

	"github.com/andybalholm/brotli"
	"github.com/valyala/bytebufferpool"
)

func TestAppendBrotliBytesLevel(t *testing.T) {
	src := []byte("Hello, World!")
	dst := make([]byte, 0)
	result := AppendBrotliBytesLevel(dst, src, CompressBrotliDefaultCompression)

	// Verify that the result is not empty
	if len(result) == 0 {
		t.Errorf("Expected non-empty result")
	}

	// Decompress and verify contents
	decompressed, err := AppendUnbrotliBytes(nil, result)
	if err != nil {
		t.Fatalf("Failed to decompress: %v", err)
	}

	if !bytes.Equal(decompressed, src) {
		t.Errorf("Expected %q, got %q", src, decompressed)
	}
}

func TestWriteBrotliLevel(t *testing.T) {
	src := []byte("Hello, World!")
	var buf bytes.Buffer
	n, err := WriteBrotliLevel(&buf, src, CompressBrotliDefaultCompression)

	if err != nil {
		t.Fatalf("WriteBrotliLevel failed: %v", err)
	}

	if n != len(src) {
		t.Errorf("Expected %d bytes written, got %d", len(src), n)
	}

	// Decompress and verify contents
	decompressed, err := AppendUnbrotliBytes(nil, buf.Bytes())
	if err != nil {
		t.Fatalf("Failed to decompress: %v", err)
	}

	if !bytes.Equal(decompressed, src) {
		t.Errorf("Expected %q, got %q", src, decompressed)
	}
}

func TestWriteUnbrotli(t *testing.T) {
	src := []byte("Hello, World!")
	compressed := AppendBrotliBytes(nil, src)

	var buf bytes.Buffer
	n, err := WriteUnbrotli(&buf, compressed)
	if err != nil {
		t.Fatalf("WriteUnbrotli failed: %v", err)
	}

	if n != len(src) {
		t.Errorf("Expected %d bytes written, got %d", len(src), n)
	}

	if !bytes.Equal(buf.Bytes(), src) {
		t.Errorf("Expected %q, got %q", src, buf.Bytes())
	}
}

func TestAppendUnbrotliBytes(t *testing.T) {
	src := []byte("Hello, World!")
	compressed := AppendBrotliBytes(nil, src)
	dst := make([]byte, 0)

	result, err := AppendUnbrotliBytes(dst, compressed)
	if err != nil {
		t.Fatalf("AppendUnbrotliBytes failed: %v", err)
	}

	if !bytes.Equal(result, src) {
		t.Errorf("Expected %q, got %q", src, result)
	}
}

func TestBrotliReaderPool(t *testing.T) {
	src := []byte("Hello, World!")
	compressed := AppendBrotliBytes(nil, src)
	r := bytes.NewReader(compressed)

	zr, err := acquireBrotliReader(r)
	if err != nil {
		t.Fatalf("Failed to acquire brotli reader: %v", err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, zr)
	if err != nil {
		t.Fatalf("Failed to read: %v", err)
	}

	releaseBrotliReader(zr)

	if !bytes.Equal(buf.Bytes(), src) {
		t.Errorf("Expected %q, got %q", src, buf.Bytes())
	}
}

func TestBrotliWriterPool(t *testing.T) {
	src := []byte("Hello, World!")
	var buf bytes.Buffer

	zw := acquireStacklessBrotliWriter(&buf, CompressBrotliDefaultCompression)
	_, err := zw.Write(src)
	if err != nil {
		t.Fatalf("Failed to write: %v", err)
	}
	zw.Close()

	releaseStacklessBrotliWriter(zw, CompressBrotliDefaultCompression)

	decompressed, err := AppendUnbrotliBytes(nil, buf.Bytes())
	if err != nil {
		t.Fatalf("Failed to decompress: %v", err)
	}

	if !bytes.Equal(decompressed, src) {
		t.Errorf("Expected %q, got %q", src, decompressed)
	}
}