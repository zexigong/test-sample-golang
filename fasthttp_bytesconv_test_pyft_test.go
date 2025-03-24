package fasthttp

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
	"strings"
	"testing"
	"time"
)

func TestAppendHTMLEscape(t *testing.T) {
	t.Parallel()

	testAppendHTMLEscape(t, "", "")
	testAppendHTMLEscape(t, "foo", "foo")
	testAppendHTMLEscape(t, "a&b", "a&amp;b")
	testAppendHTMLEscape(t, "<&>\"'", "&lt;&amp;&gt;&#34;&#39;")
	testAppendHTMLEscape(t, "1<2>3&4", "1&lt;2&gt;3&amp;4")
}

func TestAppendHTMLEscapeBytes(t *testing.T) {
	t.Parallel()

	testAppendHTMLEscapeBytes(t, "", "")
	testAppendHTMLEscapeBytes(t, "foo", "foo")
	testAppendHTMLEscapeBytes(t, "a&b", "a&amp;b")
	testAppendHTMLEscapeBytes(t, "<&>\"'", "&lt;&amp;&gt;&#34;&#39;")
	testAppendHTMLEscapeBytes(t, "1<2>3&4", "1&lt;2&gt;3&amp;4")
}

func TestAppendIPv4(t *testing.T) {
	t.Parallel()

	testAppendIPv4(t, "0.0.0.0")
	testAppendIPv4(t, "1.2.3.4")
	testAppendIPv4(t, "192.168.0.1")
	testAppendIPv4(t, "255.255.255.255")
	testAppendIPv4(t, "10.0.0.1")
}

func TestAppendHTTPDate(t *testing.T) {
	t.Parallel()

	date := time.Now().Truncate(time.Second).UTC()

	b := AppendHTTPDate(nil, date)
	date1, err := ParseHTTPDate(b)
	if err != nil {
		t.Fatalf("unexpected error when parsing HTTP date %q: %s", b, err)
	}
	if !date1.Equal(date) {
		t.Fatalf("unexpected date: %s. Expecting %s", date1, date)
	}

	b = date.AppendFormat(nil, time.RFC1123)
	date1, err = ParseHTTPDate(b)
	if err != nil {
		t.Fatalf("unexpected error when parsing HTTP date %q: %s", b, err)
	}
	if !date1.Equal(date) {
		t.Fatalf("unexpected date: %s. Expecting %s", date1, date)
	}
}

func TestParseHTTPDate(t *testing.T) {
	t.Parallel()

	testParseHTTPDate(t, "Sun, 06 Nov 1994 08:49:37 GMT")
	testParseHTTPDate(t, "Monday, 01-Jan-01 08:49:37 GMT")
}

func TestAppendUint(t *testing.T) {
	t.Parallel()

	testAppendUint(t, 0)
	testAppendUint(t, 1)
	testAppendUint(t, 123456)
	testAppendUint(t, 10)
	testAppendUint(t, 2345)
	testAppendUint(t, 6789)
	testAppendUint(t, 90909)
	testAppendUint(t, 90909)
	testAppendUint(t, 999999)
}

func TestParseUintSuccess(t *testing.T) {
	t.Parallel()

	testParseUintSuccess(t, "0", 0)
	testParseUintSuccess(t, "1", 1)
	testParseUintSuccess(t, "123456", 123456)
	testParseUintSuccess(t, "10", 10)
	testParseUintSuccess(t, "2345", 2345)
	testParseUintSuccess(t, "6789", 6789)
	testParseUintSuccess(t, "90909", 90909)
	testParseUintSuccess(t, "90909", 90909)
	testParseUintSuccess(t, "999999", 999999)

	testParseUintSuccess(t, "12345ab", 12345)
	testParseUintSuccess(t, "0x5f", 0)
}

func TestParseUintError(t *testing.T) {
	t.Parallel()

	testParseUintError(t, "")
	testParseUintError(t, "foobar")
	testParseUintError(t, "0x123456")
	testParseUintError(t, "123456x")
	testParseUintError(t, "x123456")
}

func TestParseUfloatSuccess(t *testing.T) {
	t.Parallel()

	testParseUfloatSuccess(t, "0", 0)
	testParseUfloatSuccess(t, "1", 1)
	testParseUfloatSuccess(t, "123456", 123456)
	testParseUfloatSuccess(t, "10", 10)
	testParseUfloatSuccess(t, "2345", 2345)
	testParseUfloatSuccess(t, "6789", 6789)
	testParseUfloatSuccess(t, "90909", 90909)
	testParseUfloatSuccess(t, "90909", 90909)
	testParseUfloatSuccess(t, "999999", 999999)
	testParseUfloatSuccess(t, "1.2345", 1.2345)
	testParseUfloatSuccess(t, "1.2345e2", 1.2345e2)
	testParseUfloatSuccess(t, "1.2345e-2", 1.2345e-2)
	testParseUfloatSuccess(t, "0.000001", 0.000001)

	testParseUfloatSuccess(t, "12345ab", 12345)
	testParseUfloatSuccess(t, "0x5f", 0)
}

func TestParseUfloatError(t *testing.T) {
	t.Parallel()

	testParseUfloatError(t, "")
	testParseUfloatError(t, "foobar")
	testParseUfloatError(t, "0x123456")
	testParseUfloatError(t, "123456x")
	testParseUfloatError(t, "x123456")
	testParseUfloatError(t, "1.234.56")
	testParseUfloatError(t, "-1.234")
}

func TestReadHexInt(t *testing.T) {
	t.Parallel()

	testReadHexInt(t, "0", 0)
	testReadHexInt(t, "1", 1)
	testReadHexInt(t, "a", 0xa)
	testReadHexInt(t, "f", 0xf)
	testReadHexInt(t, "10", 0x10)
	testReadHexInt(t, "1234567890abcdef", 0x1234567890abcdef)
	testReadHexInt(t, "1234567890ABCDEF", 0x1234567890ABCDEF)

	testReadHexInt(t, "1234567890abcdef-", 0x1234567890abcdef)
	testReadHexInt(t, "1234567890abcdef-foobar", 0x1234567890abcdef)
}

func TestWriteHexInt(t *testing.T) {
	t.Parallel()

	testWriteHexInt(t, 0)
	testWriteHexInt(t, 1)
	testWriteHexInt(t, 0xa)
	testWriteHexInt(t, 0xf)
	testWriteHexInt(t, 0x10)
	testWriteHexInt(t, 0x1234567890abcdef)
}

func TestLowercaseBytes(t *testing.T) {
	t.Parallel()

	testLowercaseBytes(t, "", "")
	testLowercaseBytes(t, "a", "a")
	testLowercaseBytes(t, "A", "a")
	testLowercaseBytes(t, "foo", "foo")
	testLowercaseBytes(t, "fOo", "foo")
	testLowercaseBytes(t, "Foo", "foo")
	testLowercaseBytes(t, "FOO", "foo")
	testLowercaseBytes(t, "fOO", "foo")
	testLowercaseBytes(t, "123", "123")
	testLowercaseBytes(t, "123aBC", "123abc")
	testLowercaseBytes(t, "123ABC", "123abc")
	testLowercaseBytes(t, "123abc", "123abc")
}

func TestAppendUnquotedArg(t *testing.T) {
	t.Parallel()

	testAppendUnquotedArg(t, "", "")
	testAppendUnquotedArg(t, "foo", "foo")
	testAppendUnquotedArg(t, "f+o+o", "f o o")
	testAppendUnquotedArg(t, "%20", " ")
	testAppendUnquotedArg(t, "f%20o%20o", "f o o")
	testAppendUnquotedArg(t, "%25", "%")
	testAppendUnquotedArg(t, "f%25o%25o", "f%o%o")
	testAppendUnquotedArg(t, "+%25%20+", " %+ ")
}

func TestAppendQuotedArg(t *testing.T) {
	t.Parallel()

	testAppendQuotedArg(t, "", "")
	testAppendQuotedArg(t, "foo", "foo")
	testAppendQuotedArg(t, "f o o", "f+o+o")
	testAppendQuotedArg(t, " ", "+")
	testAppendQuotedArg(t, "f o o", "f+o+o")
	testAppendQuotedArg(t, "%", "%25")
	testAppendQuotedArg(t, "f%o%o", "f%25o%25o")
	testAppendQuotedArg(t, " %+ ", "+%25%20+")
}

func TestAppendQuotedPath(t *testing.T) {
	t.Parallel()

	testAppendQuotedPath(t, "", "")
	testAppendQuotedPath(t, "foo", "foo")
	testAppendQuotedPath(t, "foo/bar", "foo/bar")
	testAppendQuotedPath(t, "foo/bar baz", "foo/bar%20baz")
	testAppendQuotedPath(t, "foo/bar%20baz", "foo/bar%2520baz")
	testAppendQuotedPath(t, "foo/bar%baz", "foo/bar%25baz")
	testAppendQuotedPath(t, "*", "*")
	testAppendQuotedPath(t, "/foo/bar/*", "/foo/bar/*")
}

func testAppendHTMLEscape(t *testing.T, s, expectedS string) {
	t.Helper()

	b := AppendHTMLEscape(nil, s)
	if string(b) != expectedS {
		t.Fatalf("unexpected escaped string: %q. Expecting %q. s=%q", b, expectedS, s)
	}

	b = AppendHTMLEscape([]byte("foobarbaz"), s)
	if string(b) != "foobarbaz"+expectedS {
		t.Fatalf("unexpected escaped string: %q. Expecting %q. s=%q", b, "foobarbaz"+expectedS, s)
	}
}

func testAppendHTMLEscapeBytes(t *testing.T, s, expectedS string) {
	t.Helper()

	b := AppendHTMLEscapeBytes(nil, []byte(s))
	if string(b) != expectedS {
		t.Fatalf("unexpected escaped string: %q. Expecting %q. s=%q", b, expectedS, s)
	}

	b = AppendHTMLEscapeBytes([]byte("foobarbaz"), []byte(s))
	if string(b) != "foobarbaz"+expectedS {
		t.Fatalf("unexpected escaped string: %q. Expecting %q. s=%q", b, "foobarbaz"+expectedS, s)
	}
}

func testAppendIPv4(t *testing.T, ipStr string) {
	t.Helper()

	ip := net.ParseIP(ipStr)

	b := AppendIPv4(nil, ip)
	if string(b) != ipStr {
		t.Fatalf("unexpected ip %q. Expecting %q", b, ipStr)
	}

	b = AppendIPv4([]byte("foobar"), ip)
	if string(b) != "foobar"+ipStr {
		t.Fatalf("unexpected ip %q. Expecting %q", b, "foobar"+ipStr)
	}
}

func testParseHTTPDate(t *testing.T, date string) {
	t.Helper()

	d, err := ParseHTTPDate([]byte(date))
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	b := AppendHTTPDate(nil, d)
	if string(b) != date {
		t.Fatalf("unexpected date: %q. Expecting %q", b, date)
	}
}

func testAppendUint(t *testing.T, n int) {
	t.Helper()

	s := fmt.Sprintf("%d", n)
	b := AppendUint(nil, n)
	if string(b) != s {
		t.Fatalf("unexpected uint %q. Expecting %q", b, s)
	}

	b = AppendUint([]byte("foobar"), n)
	if string(b) != "foobar"+s {
		t.Fatalf("unexpected uint %q. Expecting %q", b, "foobar"+s)
	}
}

func testParseUintSuccess(t *testing.T, buf string, expectedN int) {
	t.Helper()

	n, err := ParseUint([]byte(buf))
	if err != nil {
		t.Fatalf("unexpected error: %s. buf=%q", err, buf)
	}
	if n != expectedN {
		t.Fatalf("unexpected number %d. Expecting %d. buf=%q", n, expectedN, buf)
	}
}

func testParseUintError(t *testing.T, buf string) {
	t.Helper()

	n, err := ParseUint([]byte(buf))
	if err == nil {
		t.Fatalf("expecting error")
	}
	if n >= 0 {
		t.Fatalf("unexpected number %d. Expecting negative number. buf=%q", n, buf)
	}
}

func testParseUfloatSuccess(t *testing.T, buf string, expectedN float64) {
	t.Helper()

	n, err := ParseUfloat([]byte(buf))
	if err != nil {
		t.Fatalf("unexpected error: %s. buf=%q", err, buf)
	}
	if n != expectedN {
		t.Fatalf("unexpected number %f. Expecting %f. buf=%q", n, expectedN, buf)
	}
}

func testParseUfloatError(t *testing.T, buf string) {
	t.Helper()

	n, err := ParseUfloat([]byte(buf))
	if err == nil {
		t.Fatalf("expecting error")
	}
	if n >= 0 {
		t.Fatalf("unexpected number %f. Expecting negative number. buf=%q", n, buf)
	}
}

func testReadHexInt(t *testing.T, buf string, expectedN int) {
	t.Helper()

	r := bufio.NewReader(strings.NewReader(buf))
	n, err := readHexInt(r)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if n != expectedN {
		t.Fatalf("unexpected hex number %x. Expecting %x. buf=%q", n, expectedN, buf)
	}
}

func testWriteHexInt(t *testing.T, n int) {
	t.Helper()

	s := fmt.Sprintf("%x", n)
	w := &bytes.Buffer{}
	bw := bufio.NewWriter(w)
	if err := writeHexInt(bw, n); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if err := bw.Flush(); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if w.String() != s {
		t.Fatalf("unexpected hex number %q. Expecting %q", w.String(), s)
	}
}

func TestReadHexIntError(t *testing.T) {
	t.Parallel()

	testReadHexIntError(t, "")
	testReadHexIntError(t, "foobar")
	testReadHexIntError(t, "1234567890abcdef1234567890abcdef")
	testReadHexIntError(t, "1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef")
}

func TestWriteHexIntError(t *testing.T) {
	t.Parallel()

	testWriteHexIntError(t, -1)
	testWriteHexIntError(t, -123456)
	testWriteHexIntError(t, -1234567890)
}

func testReadHexIntError(t *testing.T, buf string) {
	t.Helper()

	r := bufio.NewReader(strings.NewReader(buf))
	n, err := readHexInt(r)
	if err == nil {
		t.Fatalf("expecting error")
	}
	if n >= 0 {
		t.Fatalf("unexpected hex number %x. Expecting negative number. buf=%q", n, buf)
	}
}

func testWriteHexIntError(t *testing.T, n int) {
	t.Helper()

	w := &bytes.Buffer{}
	bw := bufio.NewWriter(w)
	defer func() {
		if x := recover(); x == nil {
			t.Fatalf("expecting panic")
		}
	}()
	if err := writeHexInt(bw, n); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if err := bw.Flush(); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	t.Fatalf("unreachable")
}

func testLowercaseBytes(t *testing.T, s, expectedS string) {
	t.Helper()

	b := []byte(s)
	lowercaseBytes(b)
	if string(b) != expectedS {
		t.Fatalf("unexpected lowercased bytes %q. Expecting %q. s=%q", b, expectedS, s)
	}
}

func testAppendUnquotedArg(t *testing.T, s, expectedS string) {
	t.Helper()

	b := AppendUnquotedArg(nil, []byte(s))
	if string(b) != expectedS {
		t.Fatalf("unexpected unquoted arg: %q. Expecting %q. s=%q", b, expectedS, s)
	}

	b = AppendUnquotedArg([]byte("foobar"), []byte(s))
	if string(b) != "foobar"+expectedS {
		t.Fatalf("unexpected unquoted arg: %q. Expecting %q. s=%q", b, "foobar"+expectedS, s)
	}
}

func testAppendQuotedArg(t *testing.T, s, expectedS string) {
	t.Helper()

	b := AppendQuotedArg(nil, []byte(s))
	if string(b) != expectedS {
		t.Fatalf("unexpected quoted arg: %q. Expecting %q. s=%q", b, expectedS, s)
	}

	b = AppendQuotedArg([]byte("foobar"), []byte(s))
	if string(b) != "foobar"+expectedS {
		t.Fatalf("unexpected quoted arg: %q. Expecting %q. s=%q", b, "foobar"+expectedS, s)
	}
}

func testAppendQuotedPath(t *testing.T, s, expectedS string) {
	t.Helper()

	b := appendQuotedPath(nil, []byte(s))
	if string(b) != expectedS {
		t.Fatalf("unexpected quoted path: %q. Expecting %q. s=%q", b, expectedS, s)
	}

	b = appendQuotedPath([]byte("foobar"), []byte(s))
	if string(b) != "foobar"+expectedS {
		t.Fatalf("unexpected quoted path: %q. Expecting %q. s=%q", b, "foobar"+expectedS, s)
	}
}