// Copyright 2014 Manu Mtz-Almeida. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package gin

import (
	"bytes"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
	SetMode(TestMode)
}

func TestLoggerWithConfig(t *testing.T) {
	buffer := new(bytes.Buffer)
	router := New()
	router.Use(LoggerWithConfig(LoggerConfig{Output: buffer}))

	router.GET("/example", func(c *Context) {
		c.Request.URL.Path = "/new_path"
		c.Request.URL.RawPath = "/new_path"
	})

	// RUN
	w := performRequest(router, "GET", "/example")

	// TEST
	assert.Equal(t, 404, w.Code)
	assert.Contains(t, buffer.String(), "GET /new_path")
	assert.NotContains(t, buffer.String(), "GET /example")
}

func TestLoggerWithWriter(t *testing.T) {
	buffer := new(bytes.Buffer)
	router := New()
	router.Use(LoggerWithWriter(buffer))

	router.GET("/example", func(c *Context) {
		c.String(200, "this is an example")
	})

	// RUN
	w := performRequest(router, "GET", "/example")

	// TEST
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "this is an example", w.Body.String())
	assert.Contains(t, buffer.String(), "200")
	assert.Contains(t, buffer.String(), "GET /example")
	assert.NotContains(t, buffer.String(), "GET /notfound")
}

func TestLoggerWithWriterAndWithConfig(t *testing.T) {
	buffer1 := new(bytes.Buffer)
	buffer2 := new(bytes.Buffer)

	router := New()
	router.Use(LoggerWithWriter(buffer1))
	router.Use(LoggerWithConfig(LoggerConfig{Output: buffer2}))

	router.GET("/example", func(c *Context) {
		c.String(200, "this is an example")
	})

	// RUN
	w := performRequest(router, "GET", "/example")

	// TEST
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "this is an example", w.Body.String())

	assert.Contains(t, buffer1.String(), "200")
	assert.Contains(t, buffer1.String(), "GET /example")
	assert.NotContains(t, buffer1.String(), "GET /notfound")

	assert.Contains(t, buffer2.String(), "200")
	assert.Contains(t, buffer2.String(), "GET /example")
	assert.NotContains(t, buffer2.String(), "GET /notfound")
}

func TestLoggerWithConfigFormatting(t *testing.T) {
	buffer := new(bytes.Buffer)
	router := New()
	router.Use(LoggerWithFormatter(func(param LogFormatterParams) string {
		return param.Method + " " + param.Path + " " + param.ErrorMessage + "\n"
	}))

	router.GET("/example", func(c *Context) {
		c.String(200, "this is an example")
	})

	// RUN
	w := performRequest(router, "GET", "/example")

	// TEST
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "this is an example", w.Body.String())
	assert.Equal(t, "GET /example \n", buffer.String())
}

func TestLoggerWithConfigOutput(t *testing.T) {
	buffer := new(bytes.Buffer)
	router := New()
	router.Use(LoggerWithConfig(LoggerConfig{Output: buffer}))

	router.GET("/example", func(c *Context) {
		c.String(200, "this is an example")
	})

	// RUN
	w := performRequest(router, "GET", "/example")

	// TEST
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "this is an example", w.Body.String())
	assert.Contains(t, buffer.String(), "200")
	assert.Contains(t, buffer.String(), "GET /example")
	assert.NotContains(t, buffer.String(), "GET /notfound")
}

func TestLoggerWithConfigSkipPaths(t *testing.T) {
	buffer := new(bytes.Buffer)
	router := New()
	router.Use(LoggerWithConfig(LoggerConfig{Output: buffer, SkipPaths: []string{"/skip"}}))

	router.GET("/example", func(c *Context) {
		c.String(200, "this is an example")
	})

	router.GET("/skip", func(c *Context) {
		c.String(200, "this is a skip")
	})

	// RUN
	w := performRequest(router, "GET", "/example")

	// TEST
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "this is an example", w.Body.String())
	assert.Contains(t, buffer.String(), "200")
	assert.Contains(t, buffer.String(), "GET /example")
	assert.NotContains(t, buffer.String(), "GET /notfound")

	buffer.Reset()
	w = performRequest(router, "GET", "/skip")
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "this is a skip", w.Body.String())
	assert.NotContains(t, buffer.String(), "200")
	assert.NotContains(t, buffer.String(), "GET /skip")
}

func TestLoggerWithConfigSkipPathsFormatter(t *testing.T) {
	buffer := new(bytes.Buffer)
	router := New()
	router.Use(LoggerWithConfig(LoggerConfig{Output: buffer, SkipPaths: []string{"/skip"}, Formatter: func(param LogFormatterParams) string {
		return param.Method + " " + param.Path + "\n"
	}}))

	router.GET("/example", func(c *Context) {
		c.String(200, "this is an example")
	})

	router.GET("/skip", func(c *Context) {
		c.String(200, "this is a skip")
	})

	// RUN
	w := performRequest(router, "GET", "/example")

	// TEST
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "this is an example", w.Body.String())
	assert.Contains(t, buffer.String(), "GET /example")
	assert.NotContains(t, buffer.String(), "GET /skip")

	buffer.Reset()
	w = performRequest(router, "GET", "/skip")
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "this is a skip", w.Body.String())
	assert.NotContains(t, buffer.String(), "GET /example")
	assert.NotContains(t, buffer.String(), "GET /skip")
}

func TestLoggerWithConfigSkipPathsEmpty(t *testing.T) {
	buffer := new(bytes.Buffer)
	router := New()
	router.Use(LoggerWithConfig(LoggerConfig{Output: buffer, SkipPaths: nil}))

	router.GET("/example", func(c *Context) {
		c.String(200, "this is an example")
	})

	router.GET("/skip", func(c *Context) {
		c.String(200, "this is a skip")
	})

	// RUN
	w := performRequest(router, "GET", "/example")

	// TEST
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "this is an example", w.Body.String())
	assert.Contains(t, buffer.String(), "200")
	assert.Contains(t, buffer.String(), "GET /example")
	assert.NotContains(t, buffer.String(), "GET /notfound")

	buffer.Reset()
	w = performRequest(router, "GET", "/skip")
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "this is a skip", w.Body.String())
	assert.Contains(t, buffer.String(), "200")
	assert.Contains(t, buffer.String(), "GET /skip")
	assert.NotContains(t, buffer.String(), "GET /notfound")
}

func TestLoggerWithConfigEnsureBody(t *testing.T) {
	router := New()
	router.Use(LoggerWithConfig(LoggerConfig{Output: ioutil.Discard}))

	router.POST("/example", func(c *Context) {
		body, _ := ioutil.ReadAll(c.Request.Body)
		c.String(200, string(body))
	})

	// RUN
	w := performRequest(router, "POST", "/example", "foo=bar&bar=baz")

	// TEST
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "foo=bar&bar=baz", w.Body.String())
}

func TestLoggerWithConfigEnsureBodySingleMiddleware(t *testing.T) {
	router := New()
	router.Use(LoggerWithConfig(LoggerConfig{Output: ioutil.Discard}))

	router.POST("/example", func(c *Context) {
		body, _ := ioutil.ReadAll(c.Request.Body)
		c.String(200, string(body))
	})

	// RUN
	w := performRequest(router, "POST", "/example", "foo=bar&bar=baz")

	// TEST
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "foo=bar&bar=baz", w.Body.String())
}

func TestLoggerWithConfigEnsureBodyNoUse(t *testing.T) {
	router := New()

	router.POST("/example", func(c *Context) {
		body, _ := ioutil.ReadAll(c.Request.Body)
		c.String(200, string(body))
	})

	// RUN
	w := performRequest(router, "POST", "/example", "foo=bar&bar=baz")

	// TEST
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "foo=bar&bar=baz", w.Body.String())
}

func TestLoggerWithConfigEnsureBodyMultipleUse(t *testing.T) {
	router := New()
	router.Use(LoggerWithConfig(LoggerConfig{Output: ioutil.Discard}))
	router.Use(LoggerWithConfig(LoggerConfig{Output: ioutil.Discard}))

	router.POST("/example", func(c *Context) {
		body, _ := ioutil.ReadAll(c.Request.Body)
		c.String(200, string(body))
	})

	// RUN
	w := performRequest(router, "POST", "/example", "foo=bar&bar=baz")

	// TEST
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "foo=bar&bar=baz", w.Body.String())
}

func TestLoggerWithConfigEnsureBodyMultipleUseFirst(t *testing.T) {
	router := New()
	router.Use(LoggerWithConfig(LoggerConfig{Output: ioutil.Discard}))

	router.POST("/example", func(c *Context) {
		body, _ := ioutil.ReadAll(c.Request.Body)
		c.String(200, string(body))
	})

	router.Use(LoggerWithConfig(LoggerConfig{Output: ioutil.Discard}))

	// RUN
	w := performRequest(router, "POST", "/example", "foo=bar&bar=baz")

	// TEST
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "foo=bar&bar=baz", w.Body.String())
}

func TestLoggerWithConfigEnsureBodyMultipleUseSecond(t *testing.T) {
	router := New()

	router.POST("/example", func(c *Context) {
		body, _ := ioutil.ReadAll(c.Request.Body)
		c.String(200, string(body))
	})

	router.Use(LoggerWithConfig(LoggerConfig{Output: ioutil.Discard}))
	router.Use(LoggerWithConfig(LoggerConfig{Output: ioutil.Discard}))

	// RUN
	w := performRequest(router, "POST", "/example", "foo=bar&bar=baz")

	// TEST
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "foo=bar&bar=baz", w.Body.String())
}

func TestLoggerWithConfigEnsureBodyMultipleUseFirstAndSecond(t *testing.T) {
	router := New()
	router.Use(LoggerWithConfig(LoggerConfig{Output: ioutil.Discard}))

	router.POST("/example", func(c *Context) {
		body, _ := ioutil.ReadAll(c.Request.Body)
		c.String(200, string(body))
	})

	router.Use(LoggerWithConfig(LoggerConfig{Output: ioutil.Discard}))

	// RUN
	w := performRequest(router, "POST", "/example", "foo=bar&bar=baz")

	// TEST
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "foo=bar&bar=baz", w.Body.String())
}

func TestLoggerWithConfigEnsureBodyMultipleUseWithSkipPaths(t *testing.T) {
	router := New()
	router.Use(LoggerWithConfig(LoggerConfig{Output: ioutil.Discard, SkipPaths: []string{"/skip"}}))
	router.Use(LoggerWithConfig(LoggerConfig{Output: ioutil.Discard}))

	router.POST("/example", func(c *Context) {
		body, _ := ioutil.ReadAll(c.Request.Body)
		c.String(200, string(body))
	})

	router.POST("/skip", func(c *Context) {
		body, _ := ioutil.ReadAll(c.Request.Body)
		c.String(200, string(body))
	})

	// RUN
	w := performRequest(router, "POST", "/example", "foo=bar&bar=baz")

	// TEST
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "foo=bar&bar=baz", w.Body.String())

	// RUN
	w = performRequest(router, "POST", "/skip", "foo=bar&bar=baz")

	// TEST
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "foo=bar&bar=baz", w.Body.String())
}

func TestLoggerWithConfigEnsureBodyMultipleUseWithSkipPathsFirst(t *testing.T) {
	router := New()
	router.Use(LoggerWithConfig(LoggerConfig{Output: ioutil.Discard}))

	router.POST("/example", func(c *Context) {
		body, _ := ioutil.ReadAll(c.Request.Body)
		c.String(200, string(body))
	})

	router.POST("/skip", func(c *Context) {
		body, _ := ioutil.ReadAll(c.Request.Body)
		c.String(200, string(body))
	})

	router.Use(LoggerWithConfig(LoggerConfig{Output: ioutil.Discard, SkipPaths: []string{"/skip"}}))

	// RUN
	w := performRequest(router, "POST", "/example", "foo=bar&bar=baz")

	// TEST
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "foo=bar&bar=baz", w.Body.String())

	// RUN
	w = performRequest(router, "POST", "/skip", "foo=bar&bar=baz")

	// TEST
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "foo=bar&bar=baz", w.Body.String())
}

func TestLoggerWithConfigEnsureBodyMultipleUseWithSkipPathsSecond(t *testing.T) {
	router := New()

	router.POST("/example", func(c *Context) {
		body, _ := ioutil.ReadAll(c.Request.Body)
		c.String(200, string(body))
	})

	router.POST("/skip", func(c *Context) {
		body, _ := ioutil.ReadAll(c.Request.Body)
		c.String(200, string(body))
	})

	router.Use(LoggerWithConfig(LoggerConfig{Output: ioutil.Discard, SkipPaths: []string{"/skip"}}))
	router.Use(LoggerWithConfig(LoggerConfig{Output: ioutil.Discard}))

	// RUN
	w := performRequest(router, "POST", "/example", "foo=bar&bar=baz")

	// TEST
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "foo=bar&bar=baz", w.Body.String())

	// RUN
	w = performRequest(router, "POST", "/skip", "foo=bar&bar=baz")

	// TEST
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "foo=bar&bar=baz", w.Body.String())
}

func TestLoggerWithConfigEnsureBodyMultipleUseWithSkipPathsFirstAndSecond(t *testing.T) {
	router := New()
	router.Use(LoggerWithConfig(LoggerConfig{Output: ioutil.Discard}))

	router.POST("/example", func(c *Context) {
		body, _ := ioutil.ReadAll(c.Request.Body)
		c.String(200, string(body))
	})

	router.POST("/skip", func(c *Context) {
		body, _ := ioutil.ReadAll(c.Request.Body)
		c.String(200, string(body))
	})

	router.Use(LoggerWithConfig(LoggerConfig{Output: ioutil.Discard, SkipPaths: []string{"/skip"}}))

	// RUN
	w := performRequest(router, "POST", "/example", "foo=bar&bar=baz")

	// TEST
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "foo=bar&bar=baz", w.Body.String())

	// RUN
	w = performRequest(router, "POST", "/skip", "foo=bar&bar=baz")

	// TEST
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "foo=bar&bar=baz", w.Body.String())
}

func TestLoggerWithConfigEnsureBodyMultipleUseWithSkipPathsEmpty(t *testing.T) {
	router := New()
	router.Use(LoggerWithConfig(LoggerConfig{Output: ioutil.Discard, SkipPaths: nil}))
	router.Use(LoggerWithConfig(LoggerConfig{Output: ioutil.Discard}))

	router.POST("/example", func(c *Context) {
		body, _ := ioutil.ReadAll(c.Request.Body)
		c.String(200, string(body))
	})

	router.POST("/skip", func(c *Context) {
		body, _ := ioutil.ReadAll(c.Request.Body)
		c.String(200, string(body))
	})

	// RUN
	w := performRequest(router, "POST", "/example", "foo=bar&bar=baz")

	// TEST
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "foo=bar&bar=baz", w.Body.String())

	// RUN
	w = performRequest(router, "POST", "/skip", "foo=bar&bar=baz")

	// TEST
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "foo=bar&bar=baz", w.Body.String())
}

func TestLoggerWithConfigEnsureBodyMultipleUseWithSkipPathsEmptyFirst(t *testing.T) {
	router := New()
	router.Use(LoggerWithConfig(LoggerConfig{Output: ioutil.Discard}))

	router.POST("/example", func(c *Context) {
		body, _ := ioutil.ReadAll(c.Request.Body)
		c.String(200, string(body))
	})

	router.POST("/skip", func(c *Context) {
		body, _ := ioutil.ReadAll(c.Request.Body)
		c.String(200, string(body))
	})

	router.Use(LoggerWithConfig(LoggerConfig{Output: ioutil.Discard, SkipPaths: nil}))

	// RUN
	w := performRequest(router, "POST", "/example", "foo=bar&bar=baz")

	// TEST
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "foo=bar&bar=baz", w.Body.String())

	// RUN
	w = performRequest(router, "POST", "/skip", "foo=bar&bar=baz")

	// TEST
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "foo=bar&bar=baz", w.Body.String())
}

func TestLoggerWithConfigEnsureBodyMultipleUseWithSkipPathsEmptySecond(t *testing.T) {
	router := New()

	router.POST("/example", func(c *Context) {
		body, _ := ioutil.ReadAll(c.Request.Body)
		c.String(200, string(body))
	})

	router.POST("/skip", func(c *Context) {
		body, _ := ioutil.ReadAll(c.Request.Body)
		c.String(200, string(body))
	})

	router.Use(LoggerWithConfig(LoggerConfig{Output: ioutil.Discard, SkipPaths: nil}))
	router.Use(LoggerWithConfig(LoggerConfig{Output: ioutil.Discard}))

	// RUN
	w := performRequest(router, "POST", "/example", "foo=bar&bar=baz")

	// TEST
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "foo=bar&bar=baz", w.Body.String())

	// RUN
	w = performRequest(router, "POST", "/skip", "foo=bar&bar=baz")

	// TEST
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "foo=bar&bar=baz", w.Body.String())
}

func TestLoggerWithConfigEnsureBodyMultipleUseWithSkipPathsEmptyFirstAndSecond(t *testing.T) {
	router := New()
	router.Use(LoggerWithConfig(LoggerConfig{Output: ioutil.Discard}))

	router.POST("/example", func(c *Context) {
		body, _ := ioutil.ReadAll(c.Request.Body)
		c.String(200, string(body))
	})

	router.POST("/skip", func(c *Context) {
		body, _ := ioutil.ReadAll(c.Request.Body)
		c.String(200, string(body))
	})

	router.Use(LoggerWithConfig(LoggerConfig{Output: ioutil.Discard, SkipPaths: nil}))

	// RUN
	w := performRequest(router, "POST", "/example", "foo=bar&bar=baz")

	// TEST
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "foo=bar&bar=baz", w.Body.String())

	// RUN
	w = performRequest(router, "POST", "/skip", "foo=bar&bar=baz")

	// TEST
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "foo=bar&bar=baz", w.Body.String())
}

func TestLoggerWithConfigBodyEmpty(t *testing.T) {
	router := New()
	router.Use(LoggerWithConfig(LoggerConfig{Output: ioutil.Discard}))

	router.POST("/example", func(c *Context) {
		body, _ := ioutil.ReadAll(c.Request.Body)
		c.String(200, string(body))
	})

	// RUN
	w := performRequest(router, "POST", "/example", "")

	// TEST
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "", w.Body.String())
}

func TestLoggerWithConfigBodyEmptyFirst(t *testing.T) {
	router := New()
	router.Use(LoggerWithConfig(LoggerConfig{Output: ioutil.Discard}))

	router.POST("/example", func(c *Context) {
		body, _ := ioutil.ReadAll(c.Request.Body)
		c.String(200, string(body))
	})

	router.Use(LoggerWithConfig(LoggerConfig{Output: ioutil.Discard}))

	// RUN
	w := performRequest(router, "POST", "/example", "")

	// TEST
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "", w.Body.String())
}

func TestLoggerWithConfigBodyEmptySecond(t *testing.T) {
	router := New()

	router.POST("/example", func(c *Context) {
		body, _ := ioutil.ReadAll(c.Request.Body)
		c.String(200, string(body))
	})

	router.Use(LoggerWithConfig(LoggerConfig{Output: ioutil.Discard}))
	router.Use(LoggerWithConfig(LoggerConfig{Output: ioutil.Discard}))

	// RUN
	w := performRequest(router, "POST", "/example", "")

	// TEST
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "", w.Body.String())
}

func TestLoggerWithConfigBodyEmptyFirstAndSecond(t *testing.T) {
	router := New()
	router.Use(LoggerWithConfig(LoggerConfig{Output: ioutil.Discard}))

	router.POST("/example", func(c *Context) {
		body, _ := ioutil.ReadAll(c.Request.Body)
		c.String(200, string(body))
	})

	router.Use(LoggerWithConfig(LoggerConfig{Output: ioutil.Discard}))

	// RUN
	w := performRequest(router, "POST", "/example", "")

	// TEST
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "", w.Body.String())
}

func TestLoggerWithConfigBodyEmptyWithSkipPaths(t *testing.T) {
	router := New()
	router.Use(LoggerWithConfig(LoggerConfig{Output: ioutil.Discard, SkipPaths: []string{"/skip"}}))

	router.POST("/example", func(c *Context) {
		body, _ := ioutil.ReadAll(c.Request.Body)
		c.String(200, string(body))
	})

	router.POST("/skip", func(c *Context) {
		body, _ := ioutil.ReadAll(c.Request.Body)
		c.String(200, string(body))
	})

	// RUN
	w := performRequest(router, "POST", "/example", "")

	// TEST
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "", w.Body.String())

	// RUN
	w = performRequest(router, "POST", "/skip", "")

	// TEST
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "", w.Body.String())
}

func TestLoggerWithConfigBodyEmptyWithSkipPathsFirst(t *testing.T) {
	router := New()
	router.Use(LoggerWithConfig(LoggerConfig{Output: ioutil.Discard}))

	router.POST("/example", func(c *Context) {
		body, _ := ioutil.ReadAll(c.Request.Body)
		c.String(200, string(body))
	})

	router.POST("/skip", func(c *Context) {
		body, _ := ioutil.ReadAll(c.Request.Body)
		c.String(200, string(body))
	})

	router.Use(LoggerWithConfig(LoggerConfig{Output: ioutil.Discard, SkipPaths: []string{"/skip"}}))

	// RUN
	w := performRequest(router, "POST", "/example", "")

	// TEST
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "", w.Body.String())

	// RUN
	w = performRequest(router, "POST", "/skip", "")

	// TEST
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "", w.Body.String())
}

func TestLoggerWithConfigBodyEmptyWithSkipPathsSecond(t *testing.T) {
	router := New()

	router.POST("/example", func(c *Context) {
		body, _ := ioutil.ReadAll(c.Request.Body)
		c.String(200, string(body))
	})

	router.POST("/skip", func(c *Context) {
		body, _ := ioutil.ReadAll(c.Request.Body)
		c.String(200, string(body))
	})

	router.Use(LoggerWithConfig(LoggerConfig{Output: ioutil.Discard, SkipPaths: []string{"/skip"}}))
	router.Use(LoggerWithConfig(LoggerConfig{Output: ioutil.Discard}))

	// RUN
	w := performRequest(router, "POST", "/example", "")

	// TEST
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "", w.Body.String())

	// RUN
	w = performRequest(router, "POST", "/skip", "")

	// TEST
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "", w.Body.String())
}

func TestLoggerWithConfigBodyEmptyWithSkipPathsFirstAndSecond(t *testing.T) {
	router := New()
	router.Use(LoggerWithConfig(LoggerConfig{Output: ioutil.Discard}))

	router.POST("/example", func(c *Context) {
		body, _ := ioutil.ReadAll(c.Request.Body)
		c.String(200, string(body))
	})

	router.POST("/skip", func(c *Context) {
		body, _ := ioutil.ReadAll(c.Request.Body)
		c.String(200, string(body))
	})

	router.Use(LoggerWithConfig(LoggerConfig{Output: ioutil.Discard, SkipPaths: []string{"/skip"}}))

	// RUN
	w := performRequest(router, "POST", "/example", "")

	// TEST
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "", w.Body.String())

	// RUN
	w = performRequest(router, "POST", "/skip", "")

	// TEST
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "", w.Body.String())
}

func TestLoggerWithConfigBodyEmptyWithSkipPathsEmpty(t *testing.T) {
	router := New()
	router.Use(LoggerWithConfig(LoggerConfig{Output: ioutil.Discard, SkipPaths: nil}))

	router.POST("/example", func(c *Context) {
		body, _ := ioutil.ReadAll(c.Request.Body)
		c.String(200, string(body))
	})

	router.POST("/skip", func(c *Context) {
		body, _ := ioutil.ReadAll(c.Request.Body)
		c.String(200, string(body))
	})

	// RUN
	w := performRequest(router, "POST", "/example", "")

	// TEST
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "", w.Body.String())

	// RUN
	w = performRequest(router, "POST", "/skip", "")

	// TEST
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "", w.Body.String())
}

func TestLoggerWithConfigBodyEmptyWithSkipPathsEmptyFirst(t *testing.T) {
	router := New()
	router.Use(LoggerWithConfig(LoggerConfig{Output: ioutil.Discard}))

	router.POST("/example", func(c *Context) {
		body, _ := ioutil.ReadAll(c.Request.Body)
		c.String(200, string(body))
	})

	router.POST("/skip", func(c *Context) {
		body, _ := ioutil.ReadAll(c.Request.Body)
		c.String(200, string(body))
	})

	router.Use(LoggerWithConfig(LoggerConfig{Output: ioutil.Discard, SkipPaths: nil}))

	// RUN
	w := performRequest(router, "POST", "/example", "")

	// TEST
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "", w.Body.String())

	// RUN
	w = performRequest(router, "POST", "/skip", "")

	// TEST
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "", w.Body.String())
}

func TestLoggerWithConfigBodyEmptyWithSkipPathsEmptySecond(t *testing.T) {
	router := New()

	router.POST("/example", func(c *Context) {
		body, _ := ioutil.ReadAll(c.Request.Body)
		c.String(200, string(body))
	})

	router.POST("/skip", func(c *Context) {
		body, _ := ioutil.ReadAll(c.Request.Body)
		c.String(200, string(body))
	})

	router.Use(LoggerWithConfig(LoggerConfig{Output: ioutil.Discard, SkipPaths: nil}))
	router.Use(LoggerWithConfig(LoggerConfig{Output: ioutil.Discard}))

	// RUN
	w := performRequest(router, "POST", "/example", "")

	// TEST
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "", w.Body.String())

	// RUN
	w = performRequest(router, "POST", "/skip", "")

	// TEST
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "", w.Body.String())
}

func TestLoggerWithConfigBodyEmptyWithSkipPathsEmptyFirstAndSecond(t *testing.T) {
	router := New()
	router.Use(LoggerWithConfig(LoggerConfig{Output: ioutil.Discard}))

	router.POST("/example", func(c *Context) {
		body, _ := ioutil.ReadAll(c.Request.Body)
		c.String(200, string(body))
	})

	router.POST("/skip", func(c *Context) {
		body, _ := ioutil.ReadAll(c.Request.Body)
		c.String(200, string(body))
	})

	router.Use(LoggerWithConfig(LoggerConfig{Output: ioutil.Discard, SkipPaths: nil}))

	// RUN
	w := performRequest(router, "POST", "/example", "")

	// TEST
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "", w.Body.String())

	// RUN
	w = performRequest(router, "POST", "/skip", "")

	// TEST
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "", w.Body.String())
}

func TestLoggerWithConfigEnsureBodyWithSkipPaths(t *testing.T) {
	router := New()
	router.Use(LoggerWithConfig(LoggerConfig{Output: ioutil.Discard, SkipPaths: []string{"/skip"}}))

	router.POST("/example", func(c *Context) {
		body, _ := ioutil.ReadAll(c.Request.Body)
		c.String(200, string(body))
	})

	router.POST("/skip", func(c *Context) {
		body, _ := ioutil.ReadAll(c.Request.Body)
		c.String(200, string(body))
	})

	// RUN
	w := performRequest(router, "POST", "/example", "foo=bar&bar=baz")

	// TEST
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "foo=bar&bar=baz", w.Body.String())

	// RUN
	w = performRequest(router, "POST", "/skip", "foo=bar&bar=baz")

	// TEST
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "foo=bar&bar=baz", w.Body.String())
}

func TestLoggerWithConfigEnsureBodyWithSkipPathsFirst(t *testing.T) {
	router := New()
	router.Use(LoggerWithConfig(LoggerConfig{Output: ioutil.Discard}))

	router.POST("/example", func(c *Context) {
		body, _ := ioutil.ReadAll(c.Request.Body)
		c.String(200, string(body))
	})

	router.POST("/skip", func(c *Context) {
		body, _ := ioutil.ReadAll(c.Request.Body)
		c.String(200, string(body))
	})

	router.Use(LoggerWithConfig(LoggerConfig{Output: ioutil.Discard, SkipPaths: []string{"/skip"}}))

	// RUN
	w := performRequest(router, "POST", "/example", "foo=bar&bar=baz")

	// TEST
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "foo=bar&bar=baz", w.Body.String())

	// RUN
	w = performRequest(router, "POST", "/skip", "foo=bar&bar=baz")

	// TEST
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "foo=bar&bar=baz", w.Body.String())
}

func TestLoggerWithConfigEnsureBodyWithSkipPathsSecond(t *testing.T) {
	router := New()

	router.POST("/example", func(c *Context) {
		body, _ := ioutil.ReadAll(c.Request.Body)
		c.String(200, string(body))
	})

	router.POST("/skip", func(c *Context) {
		body, _ := ioutil.ReadAll(c.Request.Body)
		c.String(200, string(body))
	})

	router.Use(LoggerWithConfig(LoggerConfig{Output: ioutil.Discard, SkipPaths: []string{"/skip"}}))
	router.Use(LoggerWithConfig(LoggerConfig{Output: ioutil.Discard}))

	// RUN
	w := performRequest(router, "POST", "/example", "foo=bar&bar=baz")

	// TEST
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "foo=bar&bar=baz", w.Body.String())

	// RUN
	w = performRequest(router, "POST", "/skip", "foo=bar&bar=baz")

	// TEST
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "foo=bar&bar=baz", w.Body.String())
}

func TestLoggerWithConfigEnsureBodyWithSkipPathsFirstAndSecond(t *testing.T) {
	router := New()
	router.Use(LoggerWithConfig(LoggerConfig{Output: ioutil.Discard}))

	router.POST("/example", func(c *Context) {
		body, _ := ioutil.ReadAll(c.Request.Body)
		c.String(200, string(body))
	})

	router.POST("/skip", func(c *Context) {
		body, _ := ioutil.ReadAll(c.Request.Body)
		c.String(200, string(body))
	})

	router.Use(LoggerWithConfig(LoggerConfig{Output: ioutil.Discard, SkipPaths: []string{"/skip"}}))

	// RUN
	w := performRequest(router, "POST", "/example", "foo=bar&bar=baz")

	// TEST
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "foo=bar&bar=baz", w.Body.String())

	// RUN
	w = performRequest(router, "POST", "/skip", "foo=bar&bar=baz")

	// TEST
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "foo=bar&bar=baz", w.Body.String())
}

func TestLoggerWithConfigEnsureBodyWithSkipPathsEmpty(t *testing.T) {
	router := New()
	router.Use(LoggerWithConfig(LoggerConfig{Output: ioutil.Discard, SkipPaths: nil}))

	router.POST("/example", func(c *Context) {
		body, _ := ioutil.ReadAll(c.Request.Body)
		c.String(200, string(body))
	})

	router.POST("/skip", func(c *Context) {
		body, _ := ioutil.ReadAll(c.Request.Body)
		c.String(200, string(body))
	})

	// RUN
	w := performRequest(router, "POST", "/example", "foo=bar&bar=baz")

	// TEST
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "foo=bar&bar=baz", w.Body.String())

	// RUN
	w = performRequest(router, "POST", "/skip", "foo=bar&bar=baz")

	// TEST
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "foo=bar&bar=baz", w.Body.String())
}

func TestLoggerWithConfigEnsureBodyWithSkipPathsEmptyFirst(t *testing.T) {
	router := New()
	router.Use(LoggerWithConfig(LoggerConfig{Output: ioutil.Discard}))

	router.POST("/example", func(c *Context) {
		body, _ := ioutil.ReadAll(c.Request.Body)
		c.String(200, string(body))
	})

	router.POST("/skip", func(c *Context) {
		body, _ := ioutil.ReadAll(c.Request.Body)
		c.String(200, string(body))
	})

	router.Use(LoggerWithConfig(LoggerConfig{Output: ioutil.Discard, SkipPaths: nil}))

	// RUN
	w := performRequest(router, "POST", "/example", "foo=bar&bar=baz")

	// TEST
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "foo=bar&bar=baz", w.Body.String())

	// RUN
	w = performRequest(router, "POST", "/skip", "foo=bar&bar=baz")

	// TEST
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "foo=bar&bar=baz", w.Body.String())
}

func TestLoggerWithConfigEnsureBodyWithSkipPathsEmptySecond(t *testing.T) {
	router := New()

	router.POST("/example", func(c *Context) {
		body, _ := ioutil.ReadAll(c.Request.Body)
		c.String(200, string(body))
	})

	router.POST("/skip", func(c *Context) {
		body, _ := ioutil.ReadAll(c.Request.Body)
		c.String(200, string(body))
	})

	router.Use(LoggerWithConfig(LoggerConfig{Output: ioutil.Discard, SkipPaths: nil}))
	router.Use(LoggerWithConfig(LoggerConfig{Output: ioutil.Discard}))

	// RUN
	w := performRequest(router, "POST", "/example", "foo=bar&bar=baz")

	// TEST
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "foo=bar&bar=baz", w.Body.String())

	// RUN
	w = performRequest(router, "POST", "/skip", "foo=bar&bar=baz")

	// TEST
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "foo=bar&bar=baz", w.Body.String())
}

func TestLoggerWithConfigEnsureBodyWithSkipPathsEmptyFirstAndSecond(t *testing.T) {
	router := New()
	router.Use(LoggerWithConfig(LoggerConfig{Output: ioutil.Discard}))

	router.POST("/example", func(c *Context) {
		body, _ := ioutil.ReadAll(c.Request.Body)
		c.String(200, string(body))
	})

	router.POST("/skip", func(c *Context) {
		body, _ := ioutil.ReadAll(c.Request.Body)
		c.String(200, string(body))
	})

	router.Use(LoggerWithConfig(LoggerConfig{Output: ioutil.Discard, SkipPaths: nil}))

	// RUN
	w := performRequest(router, "POST", "/example", "foo=bar&bar=baz")

	// TEST
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "foo=bar&bar=baz", w.Body.String())

	// RUN
	w = performRequest(router, "POST", "/skip", "foo=bar&bar=baz")

	// TEST
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "foo=bar&bar=baz", w.Body.String())
}

func TestLoggerWithConfigEnsureBodyWithSkipPathsEmptyMultipleUse(t *testing.T) {
	router := New()
	router.Use(LoggerWithConfig(LoggerConfig{Output: ioutil.Discard, SkipPaths: nil}))
	router.Use(LoggerWithConfig(LoggerConfig{Output: ioutil.Discard}))

	router.POST("/example", func(c *Context) {
		body, _ := ioutil.ReadAll(c.Request.Body)
		c.String(200, string(body))
	})

	router.POST("/skip", func(c *Context) {
		body, _ := ioutil.ReadAll(c.Request.Body)
		c.String(200, string(body))
	})

	// RUN
	w := performRequest(router, "POST", "/example", "foo=bar&bar=baz")

	// TEST
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "foo=bar&bar=baz", w.Body.String())

	// RUN
	w = performRequest(router, "POST", "/skip", "foo=bar&bar=baz")

	// TEST
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "foo=bar&bar=baz", w.Body.String())
}

func TestLoggerWithConfigEnsureBodyWithSkipPathsEmptyMultipleUseFirst(t *testing.T) {
	router := New()
	router.Use(LoggerWithConfig(LoggerConfig{Output: ioutil.Discard}))

	router.POST("/example", func(c *Context) {
		body, _ := ioutil.ReadAll(c.Request.Body)
		c.String(200, string(body))
	})

	router.POST("/skip", func(c *Context) {
		body, _ := ioutil.ReadAll(c.Request.Body)
		c.String(200, string(body))
	})

	router.Use(LoggerWithConfig(LoggerConfig{Output: ioutil.Discard, SkipPaths: nil}))

	// RUN
	w := performRequest(router, "POST", "/example", "foo=bar&bar=baz")

	// TEST
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "foo=bar&bar=baz", w.Body.String())

	// RUN
	w = performRequest(router, "POST", "/skip", "foo=bar&bar=baz")

	// TEST
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "foo=bar&bar=baz", w.Body.String())
}

func TestLoggerWithConfigEnsureBodyWithSkipPathsEmptyMultipleUseSecond(t *testing.T) {
	router := New()

	router.POST("/example", func(c *Context) {
		body, _ := ioutil.ReadAll(c.Request.Body)
		c.String(200, string(body))
	})

	router.POST("/skip", func(c *Context) {
		body, _ := ioutil.ReadAll(c.Request.Body)
		c.String(200, string(body))
	})

	router.Use(LoggerWithConfig(LoggerConfig{Output: ioutil.Discard, SkipPaths: nil}))
	router.Use(LoggerWithConfig(LoggerConfig{Output: ioutil.Discard}))

	// RUN
	w := performRequest(router, "POST", "/example", "foo=bar&bar=baz")

	// TEST
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "foo=bar&bar=baz", w.Body.String())

	// RUN
	w = performRequest(router, "POST", "/skip", "foo=bar&bar=baz")

	// TEST
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "foo=bar&bar=baz", w.Body.String())
}

func TestLoggerWithConfigEnsureBodyWithSkipPathsEmptyMultipleUseFirstAndSecond(t *testing.T) {
	router := New()
	router.Use(LoggerWithConfig(LoggerConfig{Output: ioutil.Discard}))

	router.POST("/example", func(c *Context) {
		body, _ := ioutil.ReadAll(c.Request.Body)
		c.String(200, string(body))
	})

	router.POST("/skip", func(c *Context) {
		body, _ := ioutil.ReadAll(c.Request.Body)
		c.String(200, string(body))
	})

	router.Use(LoggerWithConfig(LoggerConfig{Output: ioutil.Discard, SkipPaths: nil}))

	// RUN
	w := performRequest(router, "POST", "/example", "foo=bar&bar=baz")

	// TEST
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "foo=bar&bar=baz", w.Body.String())

	// RUN
	w = performRequest(router, "POST", "/skip", "foo=bar&bar=baz")

	// TEST
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "foo=bar&bar=baz", w.Body.String())
}

func TestLoggerWithConfigEnsureBodyWithSkipPathsEmptyMultipleUseWithSkipPaths(t *testing.T) {
	router := New()
	router.Use(LoggerWithConfig(LoggerConfig{Output: ioutil.Discard, SkipPaths: nil}))
	router.Use(LoggerWithConfig(LoggerConfig{Output: ioutil.Discard, SkipPaths: []string{"/skip"}}))

	router.POST("/example", func(c *Context) {
		body, _ := ioutil.ReadAll(c.Request.Body)
		c.String(200, string(body))
	})

	router.POST("/skip", func(c *Context) {
		body, _ := ioutil.ReadAll(c.Request.Body)
		c.String(200, string(body))
	})

	// RUN
	w := performRequest(router, "POST", "/example", "foo=bar&bar=baz")

	// TEST
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "foo=bar&bar=baz", w.Body.String())

	// RUN
	w = performRequest(router, "POST", "/skip", "foo=bar&bar=baz")

	// TEST
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "foo=bar&bar=baz", w.Body.String())
}

func TestLoggerWithConfigEnsureBodyWithSkipPathsEmptyMultipleUseWithSkipPathsFirst(t *testing.T) {
	router := New()
	router.Use(LoggerWithConfig(LoggerConfig{Output: ioutil.Discard}))

	router.POST("/example", func(c *Context) {
		body, _ := ioutil.ReadAll(c.Request.Body)
		c.String(200, string(body))
	})

	router.POST("/skip", func(c *Context) {
		body, _ := ioutil.ReadAll(c.Request.Body)
		c.String(200, string(body))
	})

	router.Use(LoggerWithConfig(LoggerConfig{Output: ioutil.Discard, SkipPaths: nil}))
	router.Use(LoggerWithConfig(LoggerConfig{Output: ioutil.Discard, SkipPaths: []string{"/skip"}}))

	// RUN
	w := performRequest(router, "POST", "/example", "foo=bar&bar=baz")

	// TEST
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "foo=bar&bar=baz", w.Body.String())

	// RUN
	w = performRequest(router, "POST", "/skip", "foo=bar&bar=baz")

	// TEST
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "foo=bar&bar=baz", w.Body.String())
}

func TestLoggerWithConfigEnsureBodyWithSkipPathsEmptyMultipleUseWithSkipPathsSecond(t *testing.T) {
	router := New()

	router.POST("/example", func(c *Context) {
		body, _ := ioutil.ReadAll(c.Request.Body)
		c.String(200, string(body))
	})

	router.POST("/skip", func(c *Context) {
		body, _ := ioutil.ReadAll(c.Request.Body)
		c.String(200, string(body))
	})

	router.Use(LoggerWithConfig(LoggerConfig{Output: ioutil.Discard, SkipPaths: nil}))
	router.Use(LoggerWithConfig(LoggerConfig{Output: ioutil.Discard, SkipPaths: []string{"/skip"}}))

	// RUN
	w := performRequest(router, "POST", "/example", "foo=bar&bar=baz")

	// TEST
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "foo=bar&bar=baz", w.Body.String())

	// RUN
	w = performRequest(router, "POST", "/skip", "foo=bar&bar=baz")

	// TEST
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "foo=bar&bar=baz", w.Body.String())
}

func TestLoggerWithConfigEnsureBodyWithSkipPathsEmptyMultipleUseWithSkipPathsFirstAndSecond(t *testing.T) {
	router := New()
	router.Use(LoggerWithConfig(LoggerConfig{Output: ioutil.Discard}))

	router.POST("/example", func(c *Context) {
		body, _ := ioutil.ReadAll(c.Request.Body)
		c.String(200, string(body))
	})

	router.POST("/skip", func(c *Context) {
		body, _ := ioutil.ReadAll(c.Request.Body)
		c.String(200, string(body))
	})

	router.Use(LoggerWithConfig(LoggerConfig{Output: ioutil.Discard, SkipPaths: nil}))
	router.Use(LoggerWithConfig(LoggerConfig{Output: ioutil.Discard, SkipPaths: []string{"/skip"}}))

	// RUN
	w := performRequest(router, "POST", "/example", "foo=bar&bar=baz")

	// TEST
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "foo=bar&bar=baz", w.Body.String())

	// RUN
	w = performRequest(router, "POST", "/skip", "foo=bar&bar=baz")

	// TEST
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "foo=bar&bar=baz", w.Body.String())
}

func TestLoggerWithConfigEnsureBodyWithSkipPathsEmptyMultipleUseWithSkipPathsEmpty(t *testing.T) {
	router := New()
	router.Use(LoggerWithConfig(LoggerConfig{Output: ioutil.Discard, SkipPaths: nil}))
	router.Use(LoggerWithConfig(LoggerConfig{Output: ioutil.Discard, SkipPaths: nil}))

	router.POST("/example", func(c *Context) {
		body, _ := ioutil.ReadAll(c.Request.Body)
		c.String(200, string(body))
	})

	router.POST("/skip", func(c *Context) {
		body, _ := ioutil.ReadAll(c.Request.Body)
		c.String(200, string(body))
	})

	// RUN
	w := performRequest(router, "POST", "/example", "foo=bar&bar=baz")

	// TEST
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "foo=bar&bar=baz", w.Body.String())

	// RUN
	w = performRequest(router, "POST", "/skip", "foo=bar&bar=baz")

	// TEST
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "foo=bar&bar=baz", w.Body.String())
}

func TestLoggerWithConfigEnsureBodyWithSkipPathsEmptyMultipleUseWithSkipPathsEmptyFirst(t *testing.T) {
	router := New()
	router.Use(LoggerWithConfig(LoggerConfig{Output: ioutil.Discard}))

	router.POST("/example", func(c *Context) {
		body, _ := ioutil.ReadAll(c.Request.Body)
		c.String(200, string(body))
	})

	router.POST("/skip", func(c *Context) {
		body, _ := ioutil.ReadAll(c.Request.Body)
		c.String(200, string(body))
	})

	router.Use(LoggerWithConfig(LoggerConfig{Output: ioutil.Discard, SkipPaths: nil}))
	router.Use(LoggerWithConfig(LoggerConfig{Output: ioutil.Discard, SkipPaths: nil}))

	// RUN
	w := performRequest(router, "POST", "/example", "foo=bar&bar=baz")

	// TEST
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "foo=bar&bar=baz", w.Body.String())

	// RUN
	w = performRequest(router, "POST", "/skip", "foo=bar&bar=baz")

	// TEST
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "foo=bar&bar=baz", w.Body.String())
}

func TestLoggerWithConfigEnsureBodyWithSkipPathsEmptyMultipleUseWithSkipPathsEmptySecond(t *testing.T) {
	router := New()

	router.POST("/example", func(c *Context) {
		body, _ := ioutil.ReadAll(c.Request.Body)
		c.String(200, string(body))
	})

	router.POST("/skip", func(c *Context) {
		body, _ := ioutil.ReadAll(c.Request.Body)
		c.String(200, string(body))
	})

	router.Use(LoggerWithConfig(LoggerConfig{Output: ioutil.Discard, SkipPaths: nil}))
	router.Use(LoggerWithConfig(LoggerConfig{Output: ioutil.Discard, SkipPaths: nil}))

	// RUN
	w := performRequest(router, "POST", "/example", "foo=bar&bar=baz")

	// TEST
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "foo=bar&bar=baz", w.Body.String())

	// RUN
	w = performRequest(router, "POST", "/skip", "foo=bar&bar=baz")

	// TEST
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "foo=bar&bar=baz", w.Body.String())
}

func TestLoggerWithConfigEnsureBodyWithSkipPathsEmptyMultipleUseWithSkipPathsEmptyFirstAndSecond(t *testing.T) {
	router := New()
	router.Use(LoggerWithConfig(LoggerConfig{Output: ioutil.Discard}))

	router.POST("/example", func(c *Context) {
		body, _ := ioutil.ReadAll(c.Request.Body)
		c.String(200, string(body))
	})

	router.POST("/skip", func(c *Context) {
		body, _ := ioutil.ReadAll(c.Request.Body)
		c.String(200, string(body))
	})

	router.Use(LoggerWithConfig(LoggerConfig{Output: ioutil.Discard, SkipPaths: nil}))
	router.Use(LoggerWithConfig(LoggerConfig{Output: ioutil.Discard, SkipPaths: nil}))

	// RUN
	w := performRequest(router, "POST", "/example", "foo=bar&bar=baz")

	// TEST
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "foo=bar&bar=baz", w.Body.String())

	// RUN
	w = performRequest(router, "POST", "/skip", "foo=bar&bar=baz")

	// TEST
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "foo=bar&bar=baz", w.Body.String())
}

func TestLoggerWithConfigEnsureBodyWithSkipPathsMultipleUse(t *testing.T) {
	router := New()
	router.Use(LoggerWithConfig(LoggerConfig{Output: ioutil.Discard, SkipPaths: []string{"/skip"}}))
	router.Use(LoggerWithConfig(LoggerConfig{Output: ioutil.Discard, SkipPaths: []string{"/skip"}}))

	router.POST("/example", func(c *Context) {
		body, _ := ioutil.ReadAll(c.Request.Body)
		c.String(200, string(body))
	})

	router.POST("/skip", func(c *Context) {
		body, _ := ioutil.ReadAll(c.Request.Body)
		c.String(200, string(body))
	})

	// RUN
	w := performRequest(router, "POST", "/example", "foo=bar&bar=baz")

	// TEST
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "foo=bar&bar=baz", w.Body.String())

	// RUN
	w = performRequest(router, "POST", "/skip", "foo=bar&bar=baz")

	// TEST
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "foo=bar&bar=baz", w.Body.String())
}

func TestLoggerWithConfigEnsureBodyWithSkipPathsMultipleUseFirst(t *testing.T) {
	router := New()
	router.Use(LoggerWithConfig(LoggerConfig{Output: ioutil.Discard}))

	router.POST("/example", func(c *Context) {
		body, _ := ioutil.ReadAll(c.Request.Body)
		c.String(200, string(body))
	})

	router.POST("/skip", func(c *Context) {
		body, _ := ioutil.ReadAll(c.Request.Body)
		c.String(200, string(body))
	})

	router.Use(LoggerWithConfig(LoggerConfig{Output: ioutil.Discard, SkipPaths: []string{"/skip"}}))
	router.Use(LoggerWithConfig(LoggerConfig{Output: ioutil.Discard, SkipPaths: []string{"/skip"}}))

	// RUN
	w := performRequest(router, "POST", "/example", "foo=bar&bar=baz")

	// TEST
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "foo=bar&bar=baz", w.Body.String())

	// RUN
	w = performRequest(router, "POST", "/skip", "foo=bar&bar=baz")

	// TEST
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "foo=bar&bar=baz", w.Body.String())
}

func TestLoggerWithConfigEnsureBodyWithSkipPathsMultipleUseSecond(t *testing.T) {
	router := New()

	router.POST("/example", func(c *Context) {
		body, _ := ioutil.ReadAll(c.Request.Body)
		c.String(200, string(body))
	})

	router.POST("/skip", func(c *Context) {
		body, _ := ioutil.ReadAll(c.Request.Body)
		c.String(200, string(body))
	})

	router.Use(LoggerWithConfig(LoggerConfig{Output: ioutil.Discard, SkipPaths: []string{"/skip"}}))
	router.Use(LoggerWithConfig(LoggerConfig{Output: ioutil.Discard, SkipPaths: []string{"/skip"}}))

	// RUN
	w := performRequest(router, "POST", "/example", "foo=bar&bar=baz")

	// TEST
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "foo=bar&bar=baz", w.Body.String())

	// RUN
	w = performRequest(router, "POST", "/skip", "foo=bar&bar=baz")

	// TEST
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "foo=bar&bar=baz", w.Body.String())
}

func TestLoggerWithConfigEnsureBodyWithSkipPathsMultipleUseFirstAndSecond(t *testing.T) {
	router := New()
	router.Use(LoggerWithConfig(LoggerConfig{Output: ioutil.Discard}))

	router.POST("/example", func(c *Context) {
		body, _ := ioutil.ReadAll(c.Request.Body)
		c.String(200, string(body))
	})

	router.POST("/skip", func(c *Context) {
		body, _ := ioutil.ReadAll(c.Request.Body)
		c.String(200, string(body))
	})

	router.Use(LoggerWithConfig(LoggerConfig{Output: ioutil.Discard, SkipPaths: []string{"/skip"}}))
	router.Use(LoggerWithConfig(LoggerConfig{Output: ioutil.Discard, SkipPaths: []string{"/skip"}}))

	// RUN
	w := performRequest(router, "POST", "/example", "foo=bar&bar=baz")

	// TEST
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "foo=bar&bar=baz", w.Body.String())

	// RUN
	w = performRequest(router, "POST", "/skip", "foo=bar&bar=baz")

	// TEST
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "foo=bar&bar=baz", w.Body.String())
}

func TestLoggerWithConfigEnsureBodyWithSkipPathsMultipleUseWithSkipPathsEmpty(t *testing.T) {
	router := New()
	router.Use(LoggerWithConfig(LoggerConfig{Output: ioutil.Discard, SkipPaths: []string{"/skip"}}))
	router.Use(LoggerWithConfig(LoggerConfig{Output: ioutil.Discard, SkipPaths: nil}))

	router.POST("/example", func(c *Context) {
		body, _ := ioutil.ReadAll(c.Request.Body)
		c.String(200, string(body))
	})

	router.POST("/skip", func(c *Context) {
		body, _ := ioutil.ReadAll(c.Request.Body)
		c.String(200, string(body))
	})

	// RUN
	w := performRequest(router, "POST", "/example", "foo=bar&bar=baz")

	// TEST
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "foo=bar&bar=baz", w.Body.String())

	// RUN
	w = performRequest(router, "POST", "/skip", "foo=bar&bar=baz")

	// TEST
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "foo=bar&bar=baz", w.Body.String())
}

func TestLoggerWithConfigEnsureBodyWithSkipPathsMultipleUseWithSkipPathsEmptyFirst(t *testing.T) {
	router := New()
	router.Use(LoggerWithConfig(LoggerConfig{Output: ioutil.Discard}))

	router.POST("/example", func(c *Context) {
		body, _ := ioutil.ReadAll(c.Request.Body)
		c.String(200, string(body))
	})

	router.POST("/skip", func(c *Context) {
		body, _ := ioutil.ReadAll(c.Request.Body)
		c.String(200, string(body))
	})

	router.Use(LoggerWithConfig(LoggerConfig{Output: ioutil.Discard, SkipPaths: []string{"/skip"}}))
	router.Use(LoggerWithConfig(LoggerConfig{Output: ioutil.Discard, SkipPaths: nil}))

	// RUN
	w := performRequest(router, "POST", "/example", "foo=bar&bar=baz")

	// TEST
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "foo=bar&bar=baz", w.Body.String())

	// RUN
	w = performRequest(router, "POST", "/skip", "foo=bar&bar=baz")

	// TEST
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "foo=bar&bar=baz", w.Body.String())
}

func TestLoggerWithConfigEnsureBodyWithSkipPathsMultipleUseWithSkipPathsEmptySecond(t *testing.T) {
	router := New()

	router.POST("/example", func(c *Context) {
		body, _ := ioutil.ReadAll(c.Request.Body)
		c.String(200, string(body))
	})

	router.POST("/skip", func(c *Context) {
		body, _ := ioutil.ReadAll(c.Request.Body)
		c.String(200, string(body))
	})

	router.Use(LoggerWithConfig(LoggerConfig{Output: ioutil.Discard, SkipPaths: []string{"/skip"}}))
	router.Use(LoggerWithConfig(LoggerConfig{Output: ioutil.Discard, SkipPaths: nil}))

	// RUN
	w := performRequest(router, "POST", "/example", "foo=bar&bar=baz")

	// TEST
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "foo=bar&bar=baz", w.Body.String())

	// RUN
	w = performRequest(router, "POST", "/skip", "foo=bar&bar=baz")

	// TEST
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "foo=bar&bar=baz", w.Body.String())
}

func TestLoggerWithConfigEnsureBodyWithSkipPathsMultipleUseWithSkipPathsEmptyFirstAndSecond(t *testing.T) {
	router := New()
	router.Use(LoggerWithConfig(LoggerConfig{Output: ioutil.Discard}))

	router.POST("/example", func(c *Context) {
		body, _ := ioutil.ReadAll(c.Request.Body)
		c.String(200, string(body))
	})

	router.POST("/skip", func(c *Context) {
		body, _ := ioutil.ReadAll(c.Request.Body)
		c.String(200, string(body))
	})

	router.Use(LoggerWithConfig(LoggerConfig{Output: ioutil.Discard, SkipPaths: []string{"/skip"}}))
	router.Use(LoggerWithConfig(LoggerConfig{Output: ioutil.Discard, SkipPaths: nil}))

	// RUN
	w := performRequest(router, "POST", "/example", "foo=bar&bar=baz")

	// TEST
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "foo=bar&bar=baz", w.Body.String())

	// RUN
	w = performRequest(router, "POST", "/skip", "foo=bar&bar=baz")

	// TEST
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "foo=bar&bar=baz", w.Body.String())
}

func TestLoggerWithConfigEnsureBodyWithSkipPathsMultipleUseWithSkipPathsEmptyMultipleUse(t *testing.T) {
	router := New()
	router.Use(LoggerWithConfig(LoggerConfig{Output: ioutil.Discard, SkipPaths: []string{"/skip"}}))
	router.Use(LoggerWithConfig(LoggerConfig{Output: ioutil.Discard, SkipPaths: nil}))
	router.Use(LoggerWithConfig(LoggerConfig{Output: ioutil.Discard, SkipPaths: nil}))

	router.POST("/example", func(c *Context) {
		body, _ := ioutil.ReadAll(c.Request.Body)
		c.String(200, string(body))
	})

	router.POST("/skip", func(c *Context) {
		body, _ := ioutil.ReadAll(c.Request.Body)
		c.String(200, string(body))
	})

	// RUN
	w := performRequest(router, "POST", "/example", "foo=bar&bar=baz")

	// TEST
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "foo=bar&bar=baz", w.Body.String())

	// RUN
	w = performRequest(router, "POST", "/skip", "foo=bar&bar=baz")

	// TEST
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "foo=bar&bar=baz", w.Body.String())
}

func TestLoggerWithConfigEnsureBodyWithSkipPathsMultipleUseWithSkipPathsEmptyMultipleUseFirst(t *testing.T) {
	router := New()
	router.Use(LoggerWithConfig(LoggerConfig{Output: ioutil.Discard}))

	router.POST("/example", func(c *Context) {
		body, _ := ioutil.ReadAll(c.Request.Body)
		c.String(200, string(body))
	})

	router.POST("/skip", func(c *Context) {
		body, _ := ioutil.ReadAll(c.Request.Body)
		c.String(200, string(body))
	})

	router.Use(LoggerWithConfig(LoggerConfig{Output: ioutil.Discard, SkipPaths: []string{"/skip"}}))
	router.Use(LoggerWithConfig(LoggerConfig{Output: ioutil.Discard, SkipPaths: nil}))
	router.Use(LoggerWithConfig(LoggerConfig{Output: ioutil.Discard, SkipPaths: nil}))

	// RUN
	w := performRequest(router, "POST", "/example", "foo=bar&bar=baz")

	// TEST
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "foo=bar&bar=baz", w.Body.String())

	// RUN
	w = performRequest(router, "POST", "/skip", "foo=bar&bar=baz")

	// TEST
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "foo=bar&bar=baz", w.Body.String())
}

func TestLoggerWithConfigEnsureBodyWithSkipPathsMultipleUseWithSkipPathsEmptyMultipleUseSecond(t *testing.T) {
	router := New()

	router.POST("/example", func(c *Context) {
		body, _ := ioutil.ReadAll(c.Request.Body)
		c.String(200, string(body))
	})

	router.POST("/skip", func(c *Context) {
		body, _ := ioutil.ReadAll(c.Request.Body)
		c.String(200, string(body))
	})

	router.Use(LoggerWithConfig(LoggerConfig{Output: ioutil.Discard, SkipPaths: []string{"/skip"}}))
	router.Use(LoggerWithConfig(LoggerConfig{Output: ioutil.Discard, SkipPaths: nil}))
	router.Use(LoggerWithConfig(LoggerConfig{Output: ioutil.Discard, SkipPaths: nil}))

	// RUN
	w := performRequest(router, "POST", "/example", "foo=bar&bar=baz")

	// TEST
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "foo=bar&bar=baz", w.Body.String())

	// RUN
	w = performRequest(router, "POST", "/skip", "foo=bar&bar=baz")

	// TEST
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "foo=bar&bar=baz", w.Body.String())
}

func TestLoggerWithConfigEnsureBodyWithSkipPathsMultipleUseWithSkipPathsEmptyMultipleUseFirstAndSecond(t *testing.T) {
	router := New()
	router.Use(LoggerWithConfig(LoggerConfig{Output: ioutil.Discard}))

	router.POST("/example", func(c *Context) {
		body, _ := ioutil.ReadAll(c.Request.Body)
		c.String(200, string(body))
	})

	router.POST("/skip", func(c *Context) {
		body, _ := ioutil.ReadAll(c.Request.Body)
		c.String(200, string(body))
	})

	router.Use(LoggerWithConfig(LoggerConfig{Output: ioutil.Discard, SkipPaths: []string{"/skip"}}))
	router.Use(LoggerWithConfig(LoggerConfig{Output: ioutil.Discard, SkipPaths: nil}))
	router.Use(LoggerWithConfig(LoggerConfig{Output: ioutil.Discard, SkipPaths: nil}))

	// RUN
	w := performRequest(router, "POST", "/example", "foo=bar&bar=baz")

	// TEST
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "foo=bar&bar=baz", w.Body.String())

	// RUN
	w = performRequest(router, "POST", "/skip", "foo=bar&bar=baz")

	// TEST
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "foo=bar&bar=baz", w.Body.String())
}

func TestLoggerWithConfigEnsureBodyWithSkipPathsMultipleUseWithSkipPathsEmptyMultipleUseWithSkipPaths(t *testing.T) {
	router := New()
	router.Use(LoggerWithConfig(LoggerConfig{Output: ioutil.Discard, SkipPaths: []string{"/skip"}}))
	router.Use(LoggerWithConfig(LoggerConfig{Output: ioutil.Discard, SkipPaths: nil}))
	router.Use(LoggerWithConfig(LoggerConfig{Output: ioutil.Discard, SkipPaths: []string{"/skip"}}))

	router.POST("/example", func(c *Context) {
		body, _ := ioutil.ReadAll(c.Request.Body)
		c.String(200, string(body))
	})

	router.POST("/skip", func(c *Context) {
		body, _ := ioutil.ReadAll(c.Request.Body)
		c.String(200, string(body))
	})

	// RUN
	w := performRequest(router, "POST", "/example", "foo=bar&bar=baz")

	// TEST
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "foo=bar&bar=baz", w.Body.String())

	// RUN
	w = performRequest(router, "POST", "/skip", "foo=bar&bar=baz")

	// TEST
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "foo=bar&bar=baz", w.Body.String())
}

func TestLoggerWithConfigEnsureBodyWithSkipPathsMultipleUseWithSkipPathsEmptyMultipleUseWithSkipPathsFirst(t *testing.T) {
	router := New()
	router.Use(LoggerWithConfig(LoggerConfig{Output: ioutil.Discard}))

	router.POST("/example", func(c *Context) {
		body, _ := ioutil.ReadAll(c.Request.Body)
		c.String(200, string(body))
	})

	router.POST("/skip", func(c *Context) {
		body, _ := ioutil.ReadAll(c.Request.Body)
		c.String(200, string(body))
	})

	router.Use(LoggerWithConfig(LoggerConfig{Output: ioutil.Discard, SkipPaths: []string{"/skip"}}))
	router.Use(LoggerWithConfig(LoggerConfig{Output: ioutil.Discard, SkipPaths: nil}))
	router.Use(LoggerWithConfig(LoggerConfig{Output: ioutil.Discard, SkipPaths: []string{"/skip"}}))

	// RUN
	w := performRequest(router, "POST", "/example", "foo=bar&bar=baz")

	// TEST
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "foo=bar&bar=baz", w.Body.String())

	// RUN
	w = performRequest(router, "POST", "/skip", "foo=bar&bar=baz")

	// TEST
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "foo=bar&bar=baz", w.Body.String())
}

func TestLoggerWithConfigEnsureBodyWithSkipPathsMultipleUseWithSkipPathsEmptyMultipleUseWithSkipPathsSecond(t *testing.T) {
	router := New()

	router.POST("/example", func(c *Context) {
		body, _ := ioutil.ReadAll(c.Request.Body)
		c.String(200, string(body))
	})

	router.POST("/skip", func(c *Context) {
		body, _ := ioutil.ReadAll(c.Request.Body)
		c.String(200, string(body))
	})

	router.Use(LoggerWithConfig(LoggerConfig{Output: ioutil.Discard, SkipPaths: []string{"/skip"}}))
	router.Use(LoggerWithConfig(LoggerConfig{Output: ioutil.Discard, SkipPaths: nil}))
	router.Use(LoggerWithConfig(LoggerConfig{Output: ioutil.Discard, SkipPaths: []string{"/skip"}}))

	// RUN
	w := performRequest(router, "POST", "/example", "foo=bar&bar=baz")

	// TEST
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "foo=bar&bar=baz", w.Body.String())

	// RUN
	w = performRequest(router, "POST", "/skip", "foo=bar&bar=baz")

	// TEST
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "foo=bar&bar=baz", w.Body.String())
}

func TestLoggerWithConfigEnsureBodyWithSkipPathsMultipleUseWithSkipPathsEmptyMultipleUseWithSkipPathsFirstAndSecond(t *testing.T) {
	router := New()
	router.Use(LoggerWithConfig(LoggerConfig{Output: ioutil.Discard}))

	router.POST("/example", func(c *Context) {
		body, _ := ioutil.ReadAll(c.Request.Body)
		c.String(200, string(body))
	})

	router.POST("/skip", func(c *Context) {
		body, _ := ioutil.ReadAll(c.Request.Body)
		c.String(200, string(body))
	})

	router.Use(LoggerWithConfig(LoggerConfig{Output: ioutil.Discard, SkipPaths: []string{"/skip"}}))
	router.Use(LoggerWithConfig(LoggerConfig{Output: ioutil.Discard, SkipPaths: nil}))
	router.Use(LoggerWithConfig(LoggerConfig{Output: ioutil.Discard, SkipPaths: []string{"/skip"}}))

	// RUN
	w := performRequest(router, "POST", "/example", "foo=bar&bar=baz")

	// TEST
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "foo=bar&bar=baz", w.Body.String())

	// RUN
	w = performRequest(router, "POST", "/skip", "foo=bar&bar=baz")

	// TEST
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "foo=bar&bar=baz", w.Body.String())
}

func TestLoggerWithConfigEnsureBodyWithSkipPathsMultipleUseWithSkipPathsEmptyMultipleUseWithSkipPathsEmpty(t *testing.T) {
	router := New()
	router.Use(LoggerWithConfig(LoggerConfig{Output: ioutil.Discard, SkipPaths: []string{"/skip"}}))
	router.Use(LoggerWithConfig(LoggerConfig{Output: ioutil.Discard, SkipPaths: nil}))
	router.Use(LoggerWithConfig(LoggerConfig{Output: ioutil.Discard, SkipPaths: nil}))

	router.POST("/example", func(c *Context) {
		body, _ := ioutil.ReadAll(c.Request.Body)
		c.String(200, string(body))
	})

	router.POST("/skip", func(c *Context) {
		body, _ := ioutil.ReadAll(c.Request.Body)
		c.String(200, string(body))
	})

	// RUN
	w := performRequest(router, "POST", "/example", "foo=bar&bar=baz")

	// TEST
	assert.Equal(t,