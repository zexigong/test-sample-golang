package gin

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Define a helper function to simulate a HTTP request and response
func performRequest(r http.Handler, method, path string, body string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, strings.NewReader(body))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

// TestLoggerWithFormatter tests the LoggerWithFormatter middleware
func TestLoggerWithFormatter(t *testing.T) {
	buffer := new(bytes.Buffer)
	router := New()
	router.Use(LoggerWithFormatter(func(param LogFormatterParams) string {
		return fmt.Sprintf("[%s] - %s %s %d %s\n",
			param.TimeStamp.Format(time.RFC3339),
			param.Method,
			param.Path,
			param.StatusCode,
			param.Latency,
		)
	}))

	router.GET("/ping", func(c *Context) {
		c.String(http.StatusOK, "pong")
	})

	// Simulate a GET request to /ping
	performRequest(router, "GET", "/ping", "")

	logOutput := buffer.String()

	// Check the log output contains the expected information
	assert.Contains(t, logOutput, "GET /ping")
	assert.Contains(t, logOutput, "200")
	assert.Contains(t, logOutput, "pong")
}

// TestLoggerWithWriter tests the LoggerWithWriter middleware
func TestLoggerWithWriter(t *testing.T) {
	buffer := new(bytes.Buffer)
	router := New()
	router.Use(LoggerWithWriter(buffer))

	router.GET("/ping", func(c *Context) {
		c.String(http.StatusOK, "pong")
	})

	// Simulate a GET request to /ping
	performRequest(router, "GET", "/ping", "")

	logOutput := buffer.String()

	// Check the log output contains the expected information
	assert.Contains(t, logOutput, "GET /ping")
	assert.Contains(t, logOutput, "200")
	assert.Contains(t, logOutput, "pong")
}

// TestLoggerWithConfig tests the LoggerWithConfig middleware
func TestLoggerWithConfig(t *testing.T) {
	buffer := new(bytes.Buffer)
	router := New()
	router.Use(LoggerWithConfig(LoggerConfig{
		Output: buffer,
	}))

	router.GET("/ping", func(c *Context) {
		c.String(http.StatusOK, "pong")
	})

	// Simulate a GET request to /ping
	performRequest(router, "GET", "/ping", "")

	logOutput := buffer.String()

	// Check the log output contains the expected information
	assert.Contains(t, logOutput, "GET /ping")
	assert.Contains(t, logOutput, "200")
	assert.Contains(t, logOutput, "pong")
}

// TestLoggerWithConfigDisableColors tests the LoggerWithConfig middleware with DisableColors option
func TestLoggerWithConfigDisableColors(t *testing.T) {
	buffer := new(bytes.Buffer)
	router := New()
	router.Use(LoggerWithConfig(LoggerConfig{
		Output:       buffer,
		DisableColors: true,
	}))

	router.GET("/ping", func(c *Context) {
		c.String(http.StatusOK, "pong")
	})

	// Simulate a GET request to /ping
	performRequest(router, "GET", "/ping", "")

	logOutput := buffer.String()

	// Check the log output contains the expected information
	assert.Contains(t, logOutput, "GET /ping")
	assert.Contains(t, logOutput, "200")
	assert.Contains(t, logOutput, "pong")
}

// TestLoggerWithConfigSkipPaths tests the LoggerWithConfig middleware with SkipPaths option
func TestLoggerWithConfigSkipPaths(t *testing.T) {
	buffer := new(bytes.Buffer)
	router := New()
	router.Use(LoggerWithConfig(LoggerConfig{
		Output:       buffer,
		SkipPaths:    []string{"/skip"},
	}))

	router.GET("/skip", func(c *Context) {
		c.String(http.StatusOK, "skip")
	})

	router.GET("/ping", func(c *Context) {
		c.String(http.StatusOK, "pong")
	})

	// Simulate a GET request to /skip
	performRequest(router, "GET", "/skip", "")

	// Simulate a GET request to /ping
	performRequest(router, "GET", "/ping", "")

	logOutput := buffer.String()

	// Check the log output contains the expected information
	assert.NotContains(t, logOutput, "GET /skip")
	assert.Contains(t, logOutput, "GET /ping")
	assert.Contains(t, logOutput, "200")
	assert.Contains(t, logOutput, "pong")
}

// TestLoggerWithConfigEnableFunc tests the LoggerWithConfig middleware with EnableFunc option
func TestLoggerWithConfigEnableFunc(t *testing.T) {
	buffer := new(bytes.Buffer)
	router := New()
	router.Use(LoggerWithConfig(LoggerConfig{
		Output: buffer,
		EnableFunc: func(c *Context) bool {
			return c.Request.Method == "POST"
		},
	}))

	router.POST("/ping", func(c *Context) {
		c.String(http.StatusOK, "pong")
	})

	router.GET("/ping", func(c *Context) {
		c.String(http.StatusOK, "pong")
	})

	// Simulate a POST request to /ping
	performRequest(router, "POST", "/ping", "")

	// Simulate a GET request to /ping
	performRequest(router, "GET", "/ping", "")

	logOutput := buffer.String()

	// Check the log output contains the expected information
	assert.Contains(t, logOutput, "POST /ping")
	assert.Contains(t, logOutput, "200")
	assert.Contains(t, logOutput, "pong")
	assert.NotContains(t, logOutput, "GET /ping")
}