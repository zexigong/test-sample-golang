package fasthttp

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"compress/zlib"
	"io"
	"testing"

	"github.com/valyala/bytebufferpool"
)

func TestAppendGzipBytesLevel(t *testing.T) {
	src := []byte("Hello, World!")
	dst := []byte{}

	result := AppendGzipBytesLevel(dst, src, CompressDefaultCompression)
	if len(result) == 0 {
		t.Fatalf("Expected non-empty result after compression")
	}

	// Ensure that the result can be decompressed correctly
	decompressed, err := AppendGunzipBytes(nil, result)
	if err != nil {
		t.Fatalf("Failed to decompress: %v", err)
	}
	if !bytes.Equal(decompressed, src) {
		t.Fatalf("Expected decompressed data to be %v, got %v", src, decompressed)
	}
}

func TestAppendGunzipBytes(t *testing.T) {
	src := []byte("Hello, World!")
	compressed := AppendGzipBytesLevel(nil, src, CompressDefaultCompression)

	decompressed, err := AppendGunzipBytes(nil, compressed)
	if err != nil {
		t.Fatalf("Failed to decompress: %v", err)
	}
	if !bytes.Equal(decompressed, src) {
		t.Fatalf("Expected decompressed data to be %v, got %v", src, decompressed)
	}
}

func TestAppendDeflateBytesLevel(t *testing.T) {
	src := []byte("Hello, World!")
	dst := []byte{}

	result := AppendDeflateBytesLevel(dst, src, CompressDefaultCompression)
	if len(result) == 0 {
		t.Fatalf("Expected non-empty result after compression")
	}

	// Ensure that the result can be decompressed correctly
	decompressed, err := AppendInflateBytes(nil, result)
	if err != nil {
		t.Fatalf("Failed to decompress: %v", err)
	}
	if !bytes.Equal(decompressed, src) {
		t.Fatalf("Expected decompressed data to be %v, got %v", src, decompressed)
	}
}

func TestAppendInflateBytes(t *testing.T) {
	src := []byte("Hello, World!")
	compressed := AppendDeflateBytesLevel(nil, src, CompressDefaultCompression)

	decompressed, err := AppendInflateBytes(nil, compressed)
	if err != nil {
		t.Fatalf("Failed to decompress: %v", err)
	}
	if !bytes.Equal(decompressed, src) {
		t.Fatalf("Expected decompressed data to be %v, got %v", src, decompressed)
	}
}

func TestNormalizeCompressLevel(t *testing.T) {
	testCases := []struct {
		level    int
		expected int
	}{
		{CompressNoCompression, 2},
		{CompressBestSpeed, 3},
		{CompressBestCompression, 11},
		{CompressDefaultCompression, 8},
		{CompressHuffmanOnly, 0},
		{-10, 8},
		{20, 8},
	}

	for _, tc := range testCases {
		result := normalizeCompressLevel(tc.level)
		if result != tc.expected {
			t.Errorf("Expected normalized level %d for input %d, got %d", tc.expected, tc.level, result)
		}
	}
}

func TestIsFileCompressible(t *testing.T) {
	// Create a mock file using bytebufferpool
	data := []byte("Hello, World!")
	b := bytebufferpool.Get()
	b.Set(data)

	minCompressRatio := 0.8
	isCompressible := isFileCompressible(b, minCompressRatio)
	if !isCompressible {
		t.Errorf("Expected file to be compressible")
	}

	bytebufferpool.Put(b)
}

func TestAcquireReleaseGzipReader(t *testing.T) {
	data := []byte("Hello, World!")
	compressed := AppendGzipBytesLevel(nil, data, CompressDefaultCompression)
	reader := bytes.NewReader(compressed)

	zr, err := acquireGzipReader(reader)
	if err != nil {
		t.Fatalf("Failed to acquire gzip reader: %v", err)
	}

	releaseGzipReader(zr)
}

func TestAcquireReleaseFlateReader(t *testing.T) {
	data := []byte("Hello, World!")
	compressed := AppendDeflateBytesLevel(nil, data, CompressDefaultCompression)
	reader := bytes.NewReader(compressed)

	zr, err := acquireFlateReader(reader)
	if err != nil {
		t.Fatalf("Failed to acquire flate reader: %v", err)
	}

	releaseFlateReader(zr)
}

func TestAcquireReleaseStacklessGzipWriter(t *testing.T) {
	var buf bytes.Buffer
	writer := acquireStacklessGzipWriter(&buf, CompressDefaultCompression)

	_, err := writer.Write([]byte("Hello, World!"))
	if err != nil {
		t.Fatalf("Failed to write using stackless gzip writer: %v", err)
	}

	releaseStacklessGzipWriter(writer, CompressDefaultCompression)
}

func TestAcquireReleaseRealGzipWriter(t *testing.T) {
	var buf bytes.Buffer
	writer := acquireRealGzipWriter(&buf, CompressDefaultCompression)

	_, err := writer.Write([]byte("Hello, World!"))
	if err != nil {
		t.Fatalf("Failed to write using real gzip writer: %v", err)
	}

	releaseRealGzipWriter(writer, CompressDefaultCompression)
}

func TestAcquireReleaseStacklessDeflateWriter(t *testing.T) {
	var buf bytes.Buffer
	writer := acquireStacklessDeflateWriter(&buf, CompressDefaultCompression)

	_, err := writer.Write([]byte("Hello, World!"))
	if err != nil {
		t.Fatalf("Failed to write using stackless deflate writer: %v", err)
	}

	releaseStacklessDeflateWriter(writer, CompressDefaultCompression)
}

func TestAcquireReleaseRealDeflateWriter(t *testing.T) {
	var buf bytes.Buffer
	writer := acquireRealDeflateWriter(&buf, CompressDefaultCompression)

	_, err := writer.Write([]byte("Hello, World!"))
	if err != nil {
		t.Fatalf("Failed to write using real deflate writer: %v", err)
	}

	releaseRealDeflateWriter(writer, CompressDefaultCompression)
}