// Copyright 2018 Gin Core Team. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package render

import (
	"bytes"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReaderRenderNoHeaders(t *testing.T) {
	w := httptest.NewRecorder()
	data := "test"
	headers := make(map[string]string)
	headers["Content-Length"] = "4"
	reader := Reader{
		ContentType:   "text/html; charset=utf-8",
		ContentLength: int64(len(data)),
		Reader:        bytes.NewBufferString(data),
		Headers:       nil,
	}
	err := reader.Render(w)
	assert.NoError(t, err)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, data, w.Body.String())
	assert.Equal(t, "text/html; charset=utf-8", w.Header().Get("Content-Type"))
	assert.Equal(t, headers["Content-Length"], w.Header().Get("Content-Length"))
}

func TestReaderRenderHeaders(t *testing.T) {
	w := httptest.NewRecorder()
	data := "test"
	headers := make(map[string]string)
	headers["Content-Length"] = "4"
	headers["X-Custom"] = "value"
	reader := Reader{
		ContentType:   "text/html; charset=utf-8",
		ContentLength: int64(len(data)),
		Reader:        bytes.NewBufferString(data),
		Headers:       headers,
	}
	err := reader.Render(w)
	assert.NoError(t, err)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, data, w.Body.String())
	assert.Equal(t, "text/html; charset=utf-8", w.Header().Get("Content-Type"))
	assert.Equal(t, headers["Content-Length"], w.Header().Get("Content-Length"))
	assert.Equal(t, headers["X-Custom"], w.Header().Get("X-Custom"))
}