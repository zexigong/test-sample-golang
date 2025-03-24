package fasthttp

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
	"testing"
)

func TestAppendBrotliBytesLevel(t *testing.T) {
	t.Parallel()

	testAppendBrotliBytesLevel(t, CompressBrotliNoCompression)
	testAppendBrotliBytesLevel(t, CompressBrotliBestSpeed)
	testAppendBrotliBytesLevel(t, CompressBrotliBestCompression)
	testAppendBrotliBytesLevel(t, CompressBrotliDefaultCompression)
	testAppendBrotliBytesLevel(t, 1234)
}

func testAppendBrotliBytesLevel(t *testing.T, level int) {
	dst := []byte("foobarbaz")
	for i := 0; i < 5; i++ {
		dst = AppendBrotliBytesLevel(dst, []byte("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"), level)
	}
	for i := 0; i < 5; i++ {
		dst = AppendBrotliBytesLevel(dst, []byte("bbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"), level)
	}
	for i := 0; i < 5; i++ {
		dst = AppendBrotliBytesLevel(dst, []byte("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"), level)
	}
	for i := 0; i < 5; i++ {
		dst = AppendBrotliBytesLevel(dst, []byte("bbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"), level)
	}
	_, err := AppendUnbrotliBytes(nil, dst)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
}

func TestAppendBrotliBytes(t *testing.T) {
	t.Parallel()

	s := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaa" + "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbb" + "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaa" + "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"
	src := []byte(s)

	dst := AppendBrotliBytes(nil, src)
	dst, err := AppendUnbrotliBytes(nil, dst)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if string(dst) != s {
		t.Fatalf("unexpected uncompressed data: %q. Expecting %q", dst, s)
	}
}

func TestWriteBrotliLevel(t *testing.T) {
	t.Parallel()

	testWriteBrotliLevel(t, CompressBrotliNoCompression)
	testWriteBrotliLevel(t, CompressBrotliBestSpeed)
	testWriteBrotliLevel(t, CompressBrotliBestCompression)
	testWriteBrotliLevel(t, CompressBrotliDefaultCompression)
	testWriteBrotliLevel(t, 1234)
}

func testWriteBrotliLevel(t *testing.T, level int) {
	bb := &bytebufferpool.ByteBuffer{}
	for i := 0; i < 5; i++ {
		WriteBrotliLevel(bb, []byte("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"), level) //nolint:errcheck
	}
	for i := 0; i < 5; i++ {
		WriteBrotliLevel(bb, []byte("bbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"), level) //nolint:errcheck
	}
	for i := 0; i < 5; i++ {
		WriteBrotliLevel(bb, []byte("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"), level) //nolint:errcheck
	}
	for i := 0; i < 5; i++ {
		WriteBrotliLevel(bb, []byte("bbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"), level) //nolint:errcheck
	}
	_, err := AppendUnbrotliBytes(nil, bb.B)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
}

func TestWriteBrotli(t *testing.T) {
	t.Parallel()

	s := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaa" + "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbb" + "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaa" + "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"
	src := []byte(s)

	var bb bytebufferpool.ByteBuffer
	if _, err := WriteBrotli(&bb, src); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	dst, err := AppendUnbrotliBytes(nil, bb.B)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if string(dst) != s {
		t.Fatalf("unexpected uncompressed data: %q. Expecting %q", dst, s)
	}
}

func TestAppendUnbrotliBytes(t *testing.T) {
	t.Parallel()

	s := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaa" + "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbb" + "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaa" + "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"
	src := []byte(s)

	dst := AppendBrotliBytes(nil, src)
	dst, err := AppendUnbrotliBytes(nil, dst)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if string(dst) != s {
		t.Fatalf("unexpected uncompressed data: %q. Expecting %q", dst, s)
	}

	// verify that unbrotliBytes returns all the uncompressed bytes
	dst = AppendBrotliBytes(nil, src)
	dst = append(dst, "foobar"...)
	dst, err = AppendUnbrotliBytes(nil, dst)
	if err == nil {
		t.Fatalf("expecting error")
	}
	if !strings.Contains(err.Error(), "too much data") {
		t.Fatalf("unexpected error: %s. Expecting 'too much data'", err)
	}
	if string(dst) != s {
		t.Fatalf("unexpected uncompressed data: %q. Expecting %q", dst, s)
	}
}

func TestWriteUnbrotli(t *testing.T) {
	t.Parallel()

	s := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaa" + "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbb" + "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaa" + "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"
	src := []byte(s)

	var bb bytebufferpool.ByteBuffer
	if _, err := WriteBrotli(&bb, src); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	var bb2 bytebufferpool.ByteBuffer
	if _, err := WriteUnbrotli(&bb2, bb.B); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if string(bb2.B) != s {
		t.Fatalf("unexpected uncompressed data: %q. Expecting %q", bb2.B, s)
	}
}

func TestBrotliCompressLevel(t *testing.T) {
	t.Parallel()

	testBrotliCompressLevel(t, CompressBrotliNoCompression)
	testBrotliCompressLevel(t, CompressBrotliBestSpeed)
	testBrotliCompressLevel(t, CompressBrotliBestCompression)
	testBrotliCompressLevel(t, CompressBrotliDefaultCompression)
	testBrotliCompressLevel(t, 1234)
}

func testBrotliCompressLevel(t *testing.T, level int) {
	s := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaa" + "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbb" + "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaa" + "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"
	src := []byte(s)

	dst := AppendBrotliBytesLevel(nil, src, level)
	dst, err := AppendUnbrotliBytes(nil, dst)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if string(dst) != s {
		t.Fatalf("unexpected uncompressed data: %q. Expecting %q", dst, s)
	}
}

func TestWriteUnbrotliBigChunks(t *testing.T) {
	t.Parallel()

	for _, src := range [][]byte{
		{},
		bytes.Repeat([]byte("a"), 3),
		bytes.Repeat([]byte("a"), 4),
		bytes.Repeat([]byte("a"), 5),
		bytes.Repeat([]byte("a"), 10),
		bytes.Repeat([]byte("a"), 100),
		bytes.Repeat([]byte("a"), 1000),
		bytes.Repeat([]byte("a"), 10000),
		bytes.Repeat([]byte("a"), 100000),
		bytes.Repeat([]byte("a"), 1000000),
	} {
		src := src
		t.Run(fmt.Sprintf("%d", len(src)), func(t *testing.T) {
			t.Parallel()
			testWriteUnbrotli(t, src)
		})
	}
}

func testWriteUnbrotli(t *testing.T, src []byte) {
	var bb bytebufferpool.ByteBuffer
	if _, err := WriteBrotli(&bb, src); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	var bb2 bytebufferpool.ByteBuffer
	if _, err := WriteUnbrotli(&bb2, bb.B); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if !bytes.Equal(bb2.B, src) {
		t.Fatalf("unexpected uncompressed data: got %q; want %q", bb2.B, src)
	}
}

func TestBrotliReaderError(t *testing.T) {
	t.Parallel()

	// Verify reading from a broken brotli stream returns an error
	var bb bytebufferpool.ByteBuffer
	if _, err := WriteBrotli(&bb, []byte("foobar")); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	bb.B[len(bb.B)-1]++
	_, err := AppendUnbrotliBytes(nil, bb.B)
	if err == nil {
		t.Fatalf("expecting error")
	}
}

func TestBrotliWriterError(t *testing.T) {
	t.Parallel()

	// Verify writing to a broken writer returns an error
	w := &errorWriter{}
	n, err := WriteBrotli(w, []byte("foobar"))
	if err == nil {
		t.Fatalf("expecting error")
	}
	if n != 0 {
		t.Fatalf("unexpected number of bytes written: %d. Expecting 0", n)
	}
}

func TestBrotliCompressLevelNone(t *testing.T) {
	t.Parallel()

	s := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaa" + "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbb" + "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaa" + "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"
	src := []byte(s)

	dst := AppendBrotliBytesLevel(nil, src, CompressBrotliNoCompression)
	if !bytes.Equal(dst, src) {
		t.Fatalf("unexpected compressed data %q. Expecting %q", dst, src)
	}

	dst, err := AppendUnbrotliBytes(nil, dst)
	if err == nil {
		t.Fatalf("expecting error")
	}
	if string(dst) != s {
		t.Fatalf("unexpected uncompressed data: %q. Expecting %q", dst, s)
	}
}

func TestBrotliCompressLevelNegative(t *testing.T) {
	t.Parallel()

	s := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaa" + "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbb" + "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaa" + "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"
	src := []byte(s)

	dst := AppendBrotliBytesLevel(nil, src, -1234)
	dst, err := AppendUnbrotliBytes(nil, dst)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if string(dst) != s {
		t.Fatalf("unexpected uncompressed data: %q. Expecting %q", dst, s)
	}
}

func TestBrotliCompressLevelTooBig(t *testing.T) {
	t.Parallel()

	s := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaa" + "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbb" + "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaa" + "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"
	src := []byte(s)

	dst := AppendBrotliBytesLevel(nil, src, 1234)
	dst, err := AppendUnbrotliBytes(nil, dst)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if string(dst) != s {
		t.Fatalf("unexpected uncompressed data: %q. Expecting %q", dst, s)
	}
}

func TestBrotliCompressLevelInvalid(t *testing.T) {
	t.Parallel()

	s := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaa" + "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbb" + "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaa" + "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"
	src := []byte(s)

	dst := AppendBrotliBytesLevel(nil, src, 1234)
	dst, err := AppendUnbrotliBytes(nil, dst)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if string(dst) != s {
		t.Fatalf("unexpected uncompressed data: %q. Expecting %q", dst, s)
	}
}

type errorWriter struct{}

func (w *errorWriter) Write(p []byte) (int, error) {
	return 0, errors.New("error")
}