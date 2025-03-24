package gin

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLogger(t *testing.T) {
	router := New()
	router.Use(Logger())
	router.GET("/ping", func(c *Context) {
		c.String(http.StatusOK, "pong")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "pong", w.Body.String())
}

func TestLoggerWithFormatter(t *testing.T) {
	var buffer bytes.Buffer
	logger := log.New(&buffer, "", 0)

	router := New()
	router.Use(LoggerWithFormatter(func(param LogFormatterParams) string {
		return param.Path + " " + param.Method + " " + param.ClientIP + "\n"
	}))
	router.GET("/ping", func(c *Context) {
		c.String(http.StatusOK, "pong")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "pong", w.Body.String())

	// Check if the custom log formatter is used
	expectedLog := "/ping GET 127.0.0.1\n"
	assert.Contains(t, buffer.String(), expectedLog)
}

func TestLoggerWithConfig(t *testing.T) {
	router := New()
	router.Use(LoggerWithConfig(LoggerConfig{
		Formatter: func(param LogFormatterParams) string {
			return param.Path + " " + param.Method + "\n"
		},
		Output: log.Writer(),
	}))
	router.GET("/ping", func(c *Context) {
		c.String(http.StatusOK, "pong")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "pong", w.Body.String())
}

func TestLoggerWithoutOutput(t *testing.T) {
	router := New()
	router.Use(LoggerWithConfig(LoggerConfig{
		SkipPaths: []string{"/skip"},
	}))
	router.GET("/ping", func(c *Context) {
		c.String(http.StatusOK, "pong")
	})
	router.GET("/skip", func(c *Context) {
		c.String(http.StatusOK, "skip")
	})

	w1 := httptest.NewRecorder()
	req1, _ := http.NewRequest("GET", "/ping", nil)
	router.ServeHTTP(w1, req1)

	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("GET", "/skip", nil)
	router.ServeHTTP(w2, req2)

	assert.Equal(t, http.StatusOK, w1.Code)
	assert.Equal(t, "pong", w1.Body.String())

	assert.Equal(t, http.StatusOK, w2.Code)
	assert.Equal(t, "skip", w2.Body.String())
}