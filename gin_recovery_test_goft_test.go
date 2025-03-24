// Copyright 2014 Manu Martinez-Almeida. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package gin

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRecoveryWithWriter(t *testing.T) {
	buffer := new(bytes.Buffer)
	router := New()
	router.Use(RecoveryWithWriter(buffer))
	router.GET("/recovery", func(_ *Context) {
		panic("Oupps")
	})
	router.GET("/recovery-err", func(_ *Context) {
		panic(errors.New("my error"))
	})

	// RUN
	w := performRequest(router, "GET", "/recovery")

	// TEST
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, buffer.String(), "GET /recovery HTTP/1.1")
	assert.Contains(t, buffer.String(), "Oupps")
	assert.Contains(t, buffer.String(), "TestRecoveryWithWriter")
	assert.Contains(t, buffer.String(), "gin-recovery_test.go:25")

	buffer.Reset()

	// RUN
	w = performRequest(router, "GET", "/recovery-err")

	// TEST
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, buffer.String(), "GET /recovery-err HTTP/1.1")
	assert.Contains(t, buffer.String(), "my error")
	assert.Contains(t, buffer.String(), "TestRecoveryWithWriter")
	assert.Contains(t, buffer.String(), "gin-recovery_test.go:29")
}

func TestCustomRecoveryWithWriter(t *testing.T) {
	buffer := new(bytes.Buffer)
	router := New()
	router.Use(CustomRecoveryWithWriter(buffer, func(c *Context, recovered any) {
		c.String(http.StatusBadRequest, "error")
	}))
	router.GET("/recovery", func(_ *Context) {
		panic("Oupps")
	})

	// RUN
	w := performRequest(router, "GET", "/recovery")

	// TEST
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "error", w.Body.String())
	assert.Contains(t, buffer.String(), "GET /recovery HTTP/1.1")
	assert.Contains(t, buffer.String(), "Oupps")
	assert.Contains(t, buffer.String(), "TestCustomRecoveryWithWriter")
	assert.Contains(t, buffer.String(), "gin-recovery_test.go:54")
}

func TestPanicWithBrokenPipe(t *testing.T) {
	const expectCode = 204

	var buf bytes.Buffer
	router := New()
	router.Use(RecoveryWithWriter(&buf))
	router.GET("/recovery", func(c *Context) {
		// Start writing response
		c.Header("X-Test", "Value")
		c.Status(expectCode)

		// Oops. Client connection closed
		_, err := c.Writer.Write([]byte("data"))
		assert.Error(t, err)
		assert.True(t, errors.Is(err, syscall.EPIPE))

		// This will fail with a broken pipe error and call the recovery
		panic(err)
	})

	w := httptest.NewRecorder()
	_, w.Flusher = w.Body, w.Body // emulate `http.response` (min. Go 1.14)
	w.Code = 200

	// RUN
	router.ServeHTTP(w, httptest.NewRequest("GET", "/recovery", nil))

	// TEST
	assert.Equal(t, expectCode, w.Code)
	assert.Equal(t, "Value", w.HeaderMap.Get("x-test"))
	assert.Equal(t, "data", w.Body.String())
	assert.Empty(t, buf.String()) // no panic log
}

func TestSource(t *testing.T) {
	// Test for unknown file
	assert.Equal(t, []byte("???"), source(nil, 0))

	// Test for known file
	data := [][]byte{[]byte("Hello world.")}
	assert.Equal(t, []byte("Hello world."), source(data, 1))
}

func TestFunction(t *testing.T) {
	assert.Equal(t, []byte("???"), function(0))

	pc := testFunction()
	assert.Equal(t, []byte("github.com/gin-gonic/gin.testFunction"), function(pc))
}

func testFunction() uintptr {
	pc, _, _, _ := runtime.Caller(1)
	return pc
}

func TestTimeFormat(t *testing.T) {
	assert.Equal(t, "2015/10/31 - 16:39:46", timeFormat(time.Unix(1446310786, 0)))
}

func TestRecoveryCustomType(t *testing.T) {
	buffer := new(bytes.Buffer)
	router := New()
	router.Use(RecoveryWithWriter(buffer))
	router.GET("/recovery", func(_ *Context) {
		panic(&customErrorType{
			customMessage: "user defined",
		})
	})

	// RUN
	w := performRequest(router, "GET", "/recovery")

	// TEST
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, buffer.String(), "GET /recovery HTTP/1.1")
	assert.Contains(t, buffer.String(), "user defined")
	assert.Contains(t, buffer.String(), "TestRecoveryCustomType")
	assert.Contains(t, buffer.String(), "gin-recovery_test.go:128")
}

func TestRecoveryWithError(t *testing.T) {
	buffer := new(bytes.Buffer)
	router := New()
	router.Use(RecoveryWithWriter(buffer))
	router.GET("/recovery", func(_ *Context) {
		panic(errors.New("error"))
	})

	// RUN
	w := performRequest(router, "GET", "/recovery")

	// TEST
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, buffer.String(), "GET /recovery HTTP/1.1")
	assert.Contains(t, buffer.String(), "error")
	assert.Contains(t, buffer.String(), "TestRecoveryWithError")
	assert.Contains(t, buffer.String(), "gin-recovery_test.go:146")
}

func TestRecoveryWithString(t *testing.T) {
	buffer := new(bytes.Buffer)
	router := New()
	router.Use(RecoveryWithWriter(buffer))
	router.GET("/recovery", func(_ *Context) {
		panic("error")
	})

	// RUN
	w := performRequest(router, "GET", "/recovery")

	// TEST
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, buffer.String(), "GET /recovery HTTP/1.1")
	assert.Contains(t, buffer.String(), "error")
	assert.Contains(t, buffer.String(), "TestRecoveryWithString")
	assert.Contains(t, buffer.String(), "gin-recovery_test.go:164")
}

func TestRecoveryWithEmpty(t *testing.T) {
	buffer := new(bytes.Buffer)
	router := New()
	router.Use(RecoveryWithWriter(buffer))
	router.GET("/recovery", func(_ *Context) {
		panic(nil)
	})

	// RUN
	w := performRequest(router, "GET", "/recovery")

	// TEST
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, buffer.String(), "GET /recovery HTTP/1.1")
	assert.Contains(t, buffer.String(), "TestRecoveryWithEmpty")
	assert.Contains(t, buffer.String(), "gin-recovery_test.go:182")
}

type customErrorType struct {
	customMessage string
}

func (e *customErrorType) Error() string {
	return e.customMessage
}

func (e *customErrorType) String() string {
	return e.customMessage
}