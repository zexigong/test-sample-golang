package fasthttp

import (
	"bufio"
	"bytes"
	"errors"
	"net"
	"testing"
	"time"
)

func TestAppendHTMLEscape(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Hello & World", "Hello &amp; World"},
		{"<tag>", "&lt;tag&gt;"},
		{"\"quote\"", "&#34;quote&#34;"},
		{"'single quote'", "&#39;single quote&#39;"},
		{"No Escape", "No Escape"},
	}

	for _, tt := range tests {
		result := AppendHTMLEscape(nil, tt.input)
		if string(result) != tt.expected {
			t.Errorf("AppendHTMLEscape(%q) = %q, expected %q", tt.input, result, tt.expected)
		}
	}
}

func TestAppendHTMLEscapeBytes(t *testing.T) {
	input := []byte("Hello & World")
	expected := "Hello &amp; World"
	result := AppendHTMLEscapeBytes(nil, input)
	if string(result) != expected {
		t.Errorf("AppendHTMLEscapeBytes(%q) = %q, expected %q", input, result, expected)
	}
}

func TestAppendIPv4(t *testing.T) {
	tests := []struct {
		ip       net.IP
		expected string
	}{
		{net.IPv4(127, 0, 0, 1), "127.0.0.1"},
		{net.IPv4(192, 168, 1, 1), "192.168.1.1"},
		{net.IPv4(255, 255, 255, 255), "255.255.255.255"},
	}

	for _, tt := range tests {
		result := AppendIPv4(nil, tt.ip)
		if string(result) != tt.expected {
			t.Errorf("AppendIPv4(%v) = %q, expected %q", tt.ip, result, tt.expected)
		}
	}
}

func TestParseIPv4(t *testing.T) {
	tests := []struct {
		ipStr    []byte
		expected net.IP
		err      error
	}{
		{[]byte("127.0.0.1"), net.IPv4(127, 0, 0, 1), nil},
		{[]byte("192.168.1.1"), net.IPv4(192, 168, 1, 1), nil},
		{[]byte("255.255.255.255"), net.IPv4(255, 255, 255, 255), nil},
		{[]byte(""), nil, errEmptyIPStr},
	}

	for _, tt := range tests {
		ip, err := ParseIPv4(nil, tt.ipStr)
		if !ip.Equal(tt.expected) || !errors.Is(err, tt.err) {
			t.Errorf("ParseIPv4(%q) = %v, %v; expected %v, %v", tt.ipStr, ip, err, tt.expected, tt.err)
		}
	}
}

func TestAppendHTTPDate(t *testing.T) {
	date := time.Date(2025, 3, 23, 23, 28, 42, 0, time.UTC)
	expected := date.Format(time.RFC1123)
	result := AppendHTTPDate(nil, date)
	if string(result) != expected {
		t.Errorf("AppendHTTPDate(%v) = %q, expected %q", date, result, expected)
	}
}

func TestParseHTTPDate(t *testing.T) {
	dateStr := []byte("Sun, 23 Mar 2025 23:28:42 GMT")
	expected, _ := time.Parse(time.RFC1123, string(dateStr))
	result, err := ParseHTTPDate(dateStr)
	if !result.Equal(expected) || err != nil {
		t.Errorf("ParseHTTPDate(%q) = %v, %v; expected %v, nil", dateStr, result, err, expected)
	}
}

func TestAppendUint(t *testing.T) {
	tests := []struct {
		n        int
		expected string
	}{
		{0, "0"},
		{123, "123"},
		{999999, "999999"},
	}

	for _, tt := range tests {
		result := AppendUint(nil, tt.n)
		if string(result) != tt.expected {
			t.Errorf("AppendUint(%d) = %q, expected %q", tt.n, result, tt.expected)
		}
	}
}

func TestParseUint(t *testing.T) {
	tests := []struct {
		buf      []byte
		expected int
		err      error
	}{
		{[]byte("0"), 0, nil},
		{[]byte("123"), 123, nil},
		{[]byte("999999"), 999999, nil},
		{[]byte(""), -1, errEmptyInt},
		{[]byte("abc"), -1, errUnexpectedFirstChar},
	}

	for _, tt := range tests {
		result, err := ParseUint(tt.buf)
		if result != tt.expected || !errors.Is(err, tt.err) {
			t.Errorf("ParseUint(%q) = %d, %v; expected %d, %v", tt.buf, result, err, tt.expected, tt.err)
		}
	}
}

func TestParseUfloat(t *testing.T) {
	tests := []struct {
		buf      []byte
		expected float64
		err      error
	}{
		{[]byte("0.0"), 0.0, nil},
		{[]byte("123.456"), 123.456, nil},
		{[]byte("999999.999"), 999999.999, nil},
		{[]byte("-1.23"), -1, errors.New("negative input is invalid")},
		{[]byte("abc"), -1, strconv.ErrSyntax},
	}

	for _, tt := range tests {
		result, err := ParseUfloat(tt.buf)
		if result != tt.expected || (err != nil && !errors.Is(err, tt.err)) {
			t.Errorf("ParseUfloat(%q) = %f, %v; expected %f, %v", tt.buf, result, err, tt.expected, tt.err)
		}
	}
}

func TestReadHexInt(t *testing.T) {
	tests := []struct {
		input    string
		expected int
		err      error
	}{
		{"0", 0, nil},
		{"1A", 26, nil},
		{"FF", 255, nil},
		{"", -1, errEmptyHexNum},
		{"ZZ", -1, errEmptyHexNum},
	}

	for _, tt := range tests {
		r := bufio.NewReader(bytes.NewBufferString(tt.input))
		result, err := readHexInt(r)
		if result != tt.expected || !errors.Is(err, tt.err) {
			t.Errorf("readHexInt(%q) = %d, %v; expected %d, %v", tt.input, result, err, tt.expected, tt.err)
		}
	}
}

func TestWriteHexInt(t *testing.T) {
	tests := []struct {
		n        int
		expected string
	}{
		{0, "0"},
		{26, "1A"},
		{255, "FF"},
	}

	for _, tt := range tests {
		buf := &bytes.Buffer{}
		w := bufio.NewWriter(buf)
		err := writeHexInt(w, tt.n)
		w.Flush()
		if err != nil || buf.String() != tt.expected {
			t.Errorf("writeHexInt(%d) = %q, %v; expected %q, nil", tt.n, buf.String(), err, tt.expected)
		}
	}
}

func TestAppendUnquotedArg(t *testing.T) {
	src := []byte("Hello%20World")
	expected := "Hello World"
	result := AppendUnquotedArg(nil, src)
	if string(result) != expected {
		t.Errorf("AppendUnquotedArg(%q) = %q, expected %q", src, result, expected)
	}
}

func TestAppendQuotedArg(t *testing.T) {
	src := []byte("Hello World")
	expected := "Hello+World"
	result := AppendQuotedArg(nil, src)
	if string(result) != expected {
		t.Errorf("AppendQuotedArg(%q) = %q, expected %q", src, result, expected)
	}
}