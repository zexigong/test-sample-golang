// Copyright 2018 Gin Core Team. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package render

import (
	"io"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReaderRenderNoHeaders(t *testing.T) {
	content := "test"
	contentLength := int64(len(content))
	headers := map[string]string{}

	w := httptest.NewRecorder()
	reader := Reader{
		ContentType:   "text/plain",
		ContentLength: contentLength,
		Reader:        strings.NewReader(content),
		Headers:       headers,
	}

	assert.NoError(t, reader.Render(w))

	assert.Equal(t, "text/plain", w.Header().Get("Content-Type"))
	assert.Equal(t, strconv.FormatInt(contentLength, 10), w.Header().Get("Content-Length"))
	assert.Equal(t, content, w.Body.String())
}

func TestReaderRenderCustomHeaders(t *testing.T) {
	content := "test"
	contentLength := int64(len(content))
	headers := map[string]string{
		"X-Custom-Header": "value",
	}

	w := httptest.NewRecorder()
	reader := Reader{
		ContentType:   "text/plain",
		ContentLength: contentLength,
		Reader:        strings.NewReader(content),
		Headers:       headers,
	}

	assert.NoError(t, reader.Render(w))

	assert.Equal(t, "text/plain", w.Header().Get("Content-Type"))
	assert.Equal(t, strconv.FormatInt(contentLength, 10), w.Header().Get("Content-Length"))
	assert.Equal(t, "value", w.Header().Get("X-Custom-Header"))
	assert.Equal(t, content, w.Body.String())
}

func TestReaderRenderWithNegativeContentLength(t *testing.T) {
	content := "test"
	contentLength := int64(-1)
	headers := map[string]string{
		"X-Custom-Header": "value",
	}

	w := httptest.NewRecorder()
	reader := Reader{
		ContentType:   "text/plain",
		ContentLength: contentLength,
		Reader:        strings.NewReader(content),
		Headers:       headers,
	}

	assert.NoError(t, reader.Render(w))

	assert.Equal(t, "text/plain", w.Header().Get("Content-Type"))
	assert.Equal(t, "", w.Header().Get("Content-Length"))
	assert.Equal(t, "value", w.Header().Get("X-Custom-Header"))
	assert.Equal(t, content, w.Body.String())
}

type testReader struct {
	io.Reader
	closeFunc func() error
}

func (r *testReader) Close() error {
	return r.closeFunc()
}

func TestReaderRenderWithClose(t *testing.T) {
	content := "test"
	contentLength := int64(len(content))
	headers := map[string]string{
		"X-Custom-Header": "value",
	}

	w := httptest.NewRecorder()
	reader := Reader{
		ContentType:   "text/plain",
		ContentLength: contentLength,
		Reader: &testReader{
			Reader: strings.NewReader(content),
			closeFunc: func() error {
				return nil
			},
		},
		Headers: headers,
	}

	assert.NoError(t, reader.Render(w))

	assert.Equal(t, "text/plain", w.Header().Get("Content-Type"))
	assert.Equal(t, strconv.FormatInt(contentLength, 10), w.Header().Get("Content-Length"))
	assert.Equal(t, "value", w.Header().Get("X-Custom-Header"))
	assert.Equal(t, content, w.Body.String())
}