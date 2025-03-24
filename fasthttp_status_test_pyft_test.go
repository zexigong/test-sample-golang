package fasthttp

import (
	"bytes"
	"net/http"
	"testing"
)

func TestStatusMessage(t *testing.T) {
	t.Parallel()

	for i := statusMessageMin; i <= statusMessageMax; i++ {
		hMsg := http.StatusText(i)
		fMsg := StatusMessage(i)
		if hMsg != fMsg {
			t.Fatalf("unexpected message for status code %d: %q. Expecting %q", i, fMsg, hMsg)
		}
	}

	hMsg := http.StatusText(1)
	fMsg := StatusMessage(1)
	if hMsg == fMsg {
		t.Fatalf("unexpected message for status code %d: %q. Expecting %q", 1, fMsg, hMsg)
	}
}

func TestFormatStatusLine(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		statusCode int
		statusText string
		proto      string
		expected   string
	}{
		{statusCode: 200, statusText: "OK", proto: "HTTP/1.1", expected: "HTTP/1.1 200 OK\r\n"},
		{statusCode: 404, statusText: "Not Found", proto: "HTTP/1.1", expected: "HTTP/1.1 404 Not Found\r\n"},
		{statusCode: 500, statusText: "Internal Server Error", proto: "HTTP/1.1", expected: "HTTP/1.1 500 Internal Server Error\r\n"},
		{statusCode: 201, statusText: "Created", proto: "HTTP/1.1", expected: "HTTP/1.1 201 Created\r\n"},
		{statusCode: 418, statusText: "I'm a teapot", proto: "HTTP/1.1", expected: "HTTP/1.1 418 I'm a teapot\r\n"},
	}

	for _, tc := range testCases {
		t.Run(tc.expected, func(t *testing.T) {
			t.Parallel()

			var b []byte
			b = formatStatusLine(b, s2b(tc.proto), tc.statusCode, s2b(tc.statusText))
			if !bytes.Equal(b, s2b(tc.expected)) {
				t.Errorf("unexpected status line: got %q, want %q", b, tc.expected)
			}
		})
	}
}