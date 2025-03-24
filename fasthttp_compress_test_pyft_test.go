package fasthttp

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"io"
	"testing"

	"github.com/klauspost/compress/zlib"
	"github.com/valyala/bytebufferpool"
)

func TestAppendGzipBytesLevel(t *testing.T) {
	t.Parallel()

	dst := []byte("foobar")
	src := []byte("12345")
	dstNew := AppendGzipBytesLevel(dst, src, CompressDefaultCompression)
	if !bytes.Equal(dst, []byte("foobar")) {
		t.Fatalf("unexpected dst: %q. Expecting %q", dst, "foobar")
	}
	if string(dstNew[:6]) != "foobar" {
		t.Fatalf("unexpected dstNew prefix: %q. Expecting %q", dstNew[:6], "foobar")
	}
	if len(dstNew) == len(dst) {
		t.Fatalf("dstNew must be different from dst")
	}

	dstNew = AppendGunzipBytes(dst[:0], dstNew[6:])
	if !bytes.Equal(dstNew, src) {
		t.Fatalf("unexpected dstNew: %q. Expecting %q", dstNew, src)
	}
}

func TestAppendGzipBytes(t *testing.T) {
	t.Parallel()

	dst := []byte("foobar")
	src := []byte("12345")
	dstNew := AppendGzipBytes(dst, src)
	if !bytes.Equal(dst, []byte("foobar")) {
		t.Fatalf("unexpected dst: %q. Expecting %q", dst, "foobar")
	}
	if string(dstNew[:6]) != "foobar" {
		t.Fatalf("unexpected dstNew prefix: %q. Expecting %q", dstNew[:6], "foobar")
	}
	if len(dstNew) == len(dst) {
		t.Fatalf("dstNew must be different from dst")
	}

	dstNew, err := AppendGunzipBytes(dst[:0], dstNew[6:])
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if !bytes.Equal(dstNew, src) {
		t.Fatalf("unexpected dstNew: %q. Expecting %q", dstNew, src)
	}
}

func TestWriteGzipLevel(t *testing.T) {
	t.Parallel()

	dst := &bytebufferpool.ByteBuffer{}
	src := []byte("12345")
	n, err := WriteGzipLevel(dst, src, CompressDefaultCompression)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if n != len(src) {
		t.Fatalf("unexpected n: %d. Expecting %d", n, len(src))
	}

	dstNew := &bytebufferpool.ByteBuffer{}
	n, err = WriteGunzip(dstNew, dst.B)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if n != len(src) {
		t.Fatalf("unexpected n: %d. Expecting %d", n, len(src))
	}
	if !bytes.Equal(dstNew.B, src) {
		t.Fatalf("unexpected dstNew: %q. Expecting %q", dstNew.B, src)
	}
}

func TestWriteGzipLevelError(t *testing.T) {
	t.Parallel()

	t.Run("error on close", func(t *testing.T) {
		t.Parallel()

		dst := &errCloseWriter{}
		src := []byte("12345")
		n, err := WriteGzipLevel(dst, src, CompressDefaultCompression)
		if err != errClose {
			t.Fatalf("unexpected error: %s. Expecting %s", err, errClose)
		}
		if n != len(src) {
			t.Fatalf("unexpected n: %d. Expecting %d", n, len(src))
		}
	})

	t.Run("error on flush", func(t *testing.T) {
		t.Parallel()

		dst := &errFlushWriter{}
		src := []byte("12345")
		n, err := WriteGzipLevel(dst, src, CompressDefaultCompression)
		if err != errFlush {
			t.Fatalf("unexpected error: %s. Expecting %s", err, errFlush)
		}
		if n != len(src) {
			t.Fatalf("unexpected n: %d. Expecting %d", n, len(src))
		}
	})

	t.Run("error on write", func(t *testing.T) {
		t.Parallel()

		dst := &errWriteWriter{}
		src := []byte("12345")
		n, err := WriteGzipLevel(dst, src, CompressDefaultCompression)
		if err != errWrite {
			t.Fatalf("unexpected error: %s. Expecting %s", err, errWrite)
		}
		if n != 0 {
			t.Fatalf("unexpected n: %d. Expecting %d", n, 0)
		}
	})
}

var (
	errClose = fmt.Errorf("error on close")
	errFlush = fmt.Errorf("error on flush")
	errWrite = fmt.Errorf("error on write")
)

type errCloseWriter struct{}

func (w *errCloseWriter) Write(p []byte) (int, error) {
	return len(p), nil
}

func (w *errCloseWriter) Close() error {
	return errClose
}

type errFlushWriter struct{}

func (w *errFlushWriter) Write(p []byte) (int, error) {
	return len(p), nil
}

func (w *errFlushWriter) Flush() error {
	return errFlush
}

type errWriteWriter struct{}

func (w *errWriteWriter) Write(p []byte) (int, error) {
	return 0, errWrite
}

func TestAppendGunzipBytes(t *testing.T) {
	t.Parallel()

	dst := []byte("foobar")
	src := []byte("12345")
	dstNew, err := AppendGunzipBytes(dst, src)
	if err == nil {
		t.Fatalf("expecting error")
	}
	if !bytes.Equal(dstNew, dst) {
		t.Fatalf("unexpected dstNew: %q. Expecting %q", dstNew, dst)
	}

	dstNew, err = AppendGunzipBytes(dst[:0], nil)
	if err == nil {
		t.Fatalf("expecting error")
	}
	if len(dstNew) != 0 {
		t.Fatalf("unexpected dstNew: %q. Expecting %q", dstNew, dst[:0])
	}

	dstNew, err = AppendGunzipBytes(dst[:0], gunzipAppendableData)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if !bytes.Equal(dstNew, gunzipAppendableDataExpected) {
		t.Fatalf("unexpected dstNew: %q. Expecting %q", dstNew, gunzipAppendableDataExpected)
	}
}

var (
	gunzipAppendableData = []byte{
		0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xff, 0x2a, 0xca,
		0x49, 0x51, 0x40, 0x04, 0x00, 0x00, 0xff, 0xff, 0xc4, 0x20, 0x5f, 0x2c,
		0x05, 0x00, 0x00, 0x00,
	}
	gunzipAppendableDataExpected = []byte("foobar")
)

func TestAppendDeflateBytesLevel(t *testing.T) {
	t.Parallel()

	dst := []byte("foobar")
	src := []byte("12345")
	dstNew := AppendDeflateBytesLevel(dst, src, CompressDefaultCompression)
	if !bytes.Equal(dst, []byte("foobar")) {
		t.Fatalf("unexpected dst: %q. Expecting %q", dst, "foobar")
	}
	if string(dstNew[:6]) != "foobar" {
		t.Fatalf("unexpected dstNew prefix: %q. Expecting %q", dstNew[:6], "foobar")
	}
	if len(dstNew) == len(dst) {
		t.Fatalf("dstNew must be different from dst")
	}

	dstNew, err := AppendInflateBytes(dst[:0], dstNew[6:])
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if !bytes.Equal(dstNew, src) {
		t.Fatalf("unexpected dstNew: %q. Expecting %q", dstNew, src)
	}
}

func TestAppendDeflateBytes(t *testing.T) {
	t.Parallel()

	dst := []byte("foobar")
	src := []byte("12345")
	dstNew := AppendDeflateBytes(dst, src)
	if !bytes.Equal(dst, []byte("foobar")) {
		t.Fatalf("unexpected dst: %q. Expecting %q", dst, "foobar")
	}
	if string(dstNew[:6]) != "foobar" {
		t.Fatalf("unexpected dstNew prefix: %q. Expecting %q", dstNew[:6], "foobar")
	}
	if len(dstNew) == len(dst) {
		t.Fatalf("dstNew must be different from dst")
	}

	dstNew, err := AppendInflateBytes(dst[:0], dstNew[6:])
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if !bytes.Equal(dstNew, src) {
		t.Fatalf("unexpected dstNew: %q. Expecting %q", dstNew, src)
	}
}

func TestWriteDeflateLevel(t *testing.T) {
	t.Parallel()

	dst := &bytebufferpool.ByteBuffer{}
	src := []byte("12345")
	n, err := WriteDeflateLevel(dst, src, CompressDefaultCompression)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if n != len(src) {
		t.Fatalf("unexpected n: %d. Expecting %d", n, len(src))
	}

	dstNew := &bytebufferpool.ByteBuffer{}
	n, err = WriteInflate(dstNew, dst.B)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if n != len(src) {
		t.Fatalf("unexpected n: %d. Expecting %d", n, len(src))
	}
	if !bytes.Equal(dstNew.B, src) {
		t.Fatalf("unexpected dstNew: %q. Expecting %q", dstNew.B, src)
	}
}

func TestWriteDeflateLevelError(t *testing.T) {
	t.Parallel()

	t.Run("error on close", func(t *testing.T) {
		t.Parallel()

		dst := &errCloseWriter{}
		src := []byte("12345")
		n, err := WriteDeflateLevel(dst, src, CompressDefaultCompression)
		if err != errClose {
			t.Fatalf("unexpected error: %s. Expecting %s", err, errClose)
		}
		if n != len(src) {
			t.Fatalf("unexpected n: %d. Expecting %d", n, len(src))
		}
	})

	t.Run("error on flush", func(t *testing.T) {
		t.Parallel()

		dst := &errFlushWriter{}
		src := []byte("12345")
		n, err := WriteDeflateLevel(dst, src, CompressDefaultCompression)
		if err != errFlush {
			t.Fatalf("unexpected error: %s. Expecting %s", err, errFlush)
		}
		if n != len(src) {
			t.Fatalf("unexpected n: %d. Expecting %d", n, len(src))
		}
	})

	t.Run("error on write", func(t *testing.T) {
		t.Parallel()

		dst := &errWriteWriter{}
		src := []byte("12345")
		n, err := WriteDeflateLevel(dst, src, CompressDefaultCompression)
		if err != errWrite {
			t.Fatalf("unexpected error: %s. Expecting %s", err, errWrite)
		}
		if n != 0 {
			t.Fatalf("unexpected n: %d. Expecting %d", n, 0)
		}
	})
}

func TestAppendInflateBytes(t *testing.T) {
	t.Parallel()

	dst := []byte("foobar")
	src := []byte("12345")
	dstNew, err := AppendInflateBytes(dst, src)
	if err == nil {
		t.Fatalf("expecting error")
	}
	if !bytes.Equal(dstNew, dst) {
		t.Fatalf("unexpected dstNew: %q. Expecting %q", dstNew, dst)
	}

	dstNew, err = AppendInflateBytes(dst[:0], nil)
	if err == nil {
		t.Fatalf("expecting error")
	}
	if len(dstNew) != 0 {
		t.Fatalf("unexpected dstNew: %q. Expecting %q", dstNew, dst[:0])
	}

	dstNew, err = AppendInflateBytes(dst[:0], inflateAppendableData)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if !bytes.Equal(dstNew, inflateAppendableDataExpected) {
		t.Fatalf("unexpected dstNew: %q. Expecting %q", dstNew, inflateAppendableDataExpected)
	}
}

var (
	inflateAppendableData = []byte{
		0x78, 0x9c, 0x2a, 0xca, 0x49, 0x51, 0x40, 0x04, 0x00, 0x00, 0xff, 0xff,
		0xc4, 0x20, 0x5f, 0x2c,
	}
	inflateAppendableDataExpected = []byte("foobar")
)

func TestCompressLevel(t *testing.T) {
	t.Parallel()

	testCompressLevel(t, CompressNoCompression)
	testCompressLevel(t, CompressBestSpeed)
	testCompressLevel(t, CompressBestCompression)
	testCompressLevel(t, CompressDefaultCompression)
	testCompressLevel(t, CompressHuffmanOnly)
	testCompressLevel(t, -1234)
	testCompressLevel(t, 1234)
}

func testCompressLevel(t *testing.T, level int) {
	for i := 0; i < 10; i++ {
		testCompressLevelSerial(t, level)
	}
}

func testCompressLevelSerial(t *testing.T, level int) {
	for i := 0; i < 100; i++ {
		testCompressLevelConcurrent(t, level)
	}
}

func testCompressLevelConcurrent(t *testing.T, level int) {
	t.Helper()

	c := make(chan struct{}, 10)
	for i := 0; i < cap(c); i++ {
		go func() {
			testCompressLevelConcurrentGoroutine(t, level)
			c <- struct{}{}
		}()
	}
	for i := 0; i < cap(c); i++ {
		<-c
	}
}

func testCompressLevelConcurrentGoroutine(t *testing.T, level int) {
	t.Helper()

	dst := &bytebufferpool.ByteBuffer{}
	dstNew := &bytebufferpool.ByteBuffer{}
	for i := 0; i < 10; i++ {
		src := []byte(fmt.Sprintf("foobar baz compress level testing %d", i))
		WriteGzipLevel(dst, src, level) //nolint:errcheck
		WriteGunzip(dstNew, dst.B)      //nolint:errcheck
		if !bytes.Equal(dstNew.B, src) {
			t.Fatalf("unexpected dstNew: %q. Expecting %q", dstNew.B, src)
		}
		dst.Reset()
		dstNew.Reset()

		WriteDeflateLevel(dst, src, level) //nolint:errcheck
		WriteInflate(dstNew, dst.B)        //nolint:errcheck
		if !bytes.Equal(dstNew.B, src) {
			t.Fatalf("unexpected dstNew: %q. Expecting %q", dstNew.B, src)
		}
		dst.Reset()
		dstNew.Reset()
	}
}

func TestAppendInflateBytesError(t *testing.T) {
	t.Parallel()

	// Verify that invalid compressed input results in an error
	src := make([]byte, 100)
	for i := range src {
		src[i] = byte(i)
	}
	dst, err := AppendInflateBytes(nil, src)
	if err == nil {
		t.Fatal("expecting non-nil error")
	}
	if dst != nil {
		t.Fatalf("unexpected dst: %q. Expecting nil", dst)
	}
}

func TestAppendGunzipBytesError(t *testing.T) {
	t.Parallel()

	// Verify that invalid compressed input results in an error
	src := make([]byte, 100)
	for i := range src {
		src[i] = byte(i)
	}
	dst, err := AppendGunzipBytes(nil, src)
	if err == nil {
		t.Fatal("expecting non-nil error")
	}
	if dst != nil {
		t.Fatalf("unexpected dst: %q. Expecting nil", dst)
	}
}

func TestInflateTooMuchData(t *testing.T) {
	t.Parallel()

	// Verify that too much data is rejected
	data := base64.StdEncoding.EncodeToString([]byte("Hello, World!"))
	var buf bytes.Buffer
	w := zlib.NewWriter(&buf)
	w.Write([]byte(data)) //nolint:errcheck
	w.Close()             //nolint:errcheck

	data = ""
	for i := 0; i < 2048; i++ {
		data += "A"
	}
	data = base64.StdEncoding.EncodeToString([]byte(data))
	w = zlib.NewWriter(&buf)
	w.Write([]byte(data)) //nolint:errcheck
	w.Close()             //nolint:errcheck

	_, err := AppendInflateBytes(nil, buf.Bytes())
	if err == nil {
		t.Fatal("expecting non-nil error")
	}
}

func TestGunzipTooMuchData(t *testing.T) {
	t.Parallel()

	// Verify that too much data is rejected
	data := base64.StdEncoding.EncodeToString([]byte("Hello, World!"))
	var buf bytes.Buffer
	w := gzip.NewWriter(&buf)
	w.Write([]byte(data)) //nolint:errcheck
	w.Close()             //nolint:errcheck

	data = ""
	for i := 0; i < 2048; i++ {
		data += "A"
	}
	data = base64.StdEncoding.EncodeToString([]byte(data))
	w = gzip.NewWriter(&buf)
	w.Write([]byte(data)) //nolint:errcheck
	w.Close()             //nolint:errcheck

	_, err := AppendGunzipBytes(nil, buf.Bytes())
	if err == nil {
		t.Fatal("expecting non-nil error")
	}
}