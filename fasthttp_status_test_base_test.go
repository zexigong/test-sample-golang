package fasthttp

import (
	"testing"
)

func TestStatusMessage(t *testing.T) {
	tests := []struct {
		statusCode int
		expected   string
	}{
		{StatusContinue, "Continue"},
		{StatusSwitchingProtocols, "Switching Protocols"},
		{StatusProcessing, "Processing"},
		{StatusEarlyHints, "Early Hints"},
		{StatusOK, "OK"},
		{StatusCreated, "Created"},
		{StatusAccepted, "Accepted"},
		{StatusNonAuthoritativeInfo, "Non-Authoritative Information"},
		{StatusNoContent, "No Content"},
		{StatusResetContent, "Reset Content"},
		{StatusPartialContent, "Partial Content"},
		{StatusMultiStatus, "Multi-Status"},
		{StatusAlreadyReported, "Already Reported"},
		{StatusIMUsed, "IM Used"},
		{StatusMultipleChoices, "Multiple Choices"},
		{StatusMovedPermanently, "Moved Permanently"},
		{StatusFound, "Found"},
		{StatusSeeOther, "See Other"},
		{StatusNotModified, "Not Modified"},
		{StatusUseProxy, "Use Proxy"},
		{StatusTemporaryRedirect, "Temporary Redirect"},
		{StatusPermanentRedirect, "Permanent Redirect"},
		{StatusBadRequest, "Bad Request"},
		{StatusUnauthorized, "Unauthorized"},
		{StatusPaymentRequired, "Payment Required"},
		{StatusForbidden, "Forbidden"},
		{StatusNotFound, "Not Found"},
		{StatusMethodNotAllowed, "Method Not Allowed"},
		{StatusNotAcceptable, "Not Acceptable"},
		{StatusProxyAuthRequired, "Proxy Authentication Required"},
		{StatusRequestTimeout, "Request Timeout"},
		{StatusConflict, "Conflict"},
		{StatusGone, "Gone"},
		{StatusLengthRequired, "Length Required"},
		{StatusPreconditionFailed, "Precondition Failed"},
		{StatusRequestEntityTooLarge, "Request Entity Too Large"},
		{StatusRequestURITooLong, "Request URI Too Long"},
		{StatusUnsupportedMediaType, "Unsupported Media Type"},
		{StatusRequestedRangeNotSatisfiable, "Requested Range Not Satisfiable"},
		{StatusExpectationFailed, "Expectation Failed"},
		{StatusTeapot, "I'm a teapot"},
		{StatusMisdirectedRequest, "Misdirected Request"},
		{StatusUnprocessableEntity, "Unprocessable Entity"},
		{StatusLocked, "Locked"},
		{StatusFailedDependency, "Failed Dependency"},
		{StatusUpgradeRequired, "Upgrade Required"},
		{StatusPreconditionRequired, "Precondition Required"},
		{StatusTooManyRequests, "Too Many Requests"},
		{StatusRequestHeaderFieldsTooLarge, "Request Header Fields Too Large"},
		{StatusUnavailableForLegalReasons, "Unavailable For Legal Reasons"},
		{StatusInternalServerError, "Internal Server Error"},
		{StatusNotImplemented, "Not Implemented"},
		{StatusBadGateway, "Bad Gateway"},
		{StatusServiceUnavailable, "Service Unavailable"},
		{StatusGatewayTimeout, "Gateway Timeout"},
		{StatusHTTPVersionNotSupported, "HTTP Version Not Supported"},
		{StatusVariantAlsoNegotiates, "Variant Also Negotiates"},
		{StatusInsufficientStorage, "Insufficient Storage"},
		{StatusLoopDetected, "Loop Detected"},
		{StatusNotExtended, "Not Extended"},
		{StatusNetworkAuthenticationRequired, "Network Authentication Required"},
		{99, "Unknown Status Code"},
		{512, "Unknown Status Code"},
	}

	for _, test := range tests {
		t.Run(strconv.Itoa(test.statusCode), func(t *testing.T) {
			result := StatusMessage(test.statusCode)
			if result != test.expected {
				t.Errorf("Expected %s, got %s", test.expected, result)
			}
		})
	}
}

func TestFormatStatusLine(t *testing.T) {
	tests := []struct {
		protocol   []byte
		statusCode int
		statusText []byte
		expected   string
	}{
		{[]byte("HTTP/1.1"), StatusOK, nil, "HTTP/1.1 200 OK\r\n"},
		{[]byte("HTTP/1.1"), StatusNotFound, nil, "HTTP/1.1 404 Not Found\r\n"},
		{[]byte("HTTP/2"), StatusInternalServerError, []byte("Custom Error"), "HTTP/2 500 Custom Error\r\n"},
		{[]byte("HTTP/1.1"), 99, nil, "HTTP/1.1 99 Unknown Status Code\r\n"},
	}

	for _, test := range tests {
		t.Run(test.expected, func(t *testing.T) {
			result := formatStatusLine(nil, test.protocol, test.statusCode, test.statusText)
			if string(result) != test.expected {
				t.Errorf("Expected %s, got %s", test.expected, string(result))
			}
		})
	}
}