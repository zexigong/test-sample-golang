// Copyright 2014 Manu Martinez-Almeida. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package gin

import (
	"bufio"
	"net"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

type DummyResponseWriter struct{}

func (_ DummyResponseWriter) Header() http.Header {
	return http.Header{}
}

func (_ DummyResponseWriter) Write([]byte) (int, error) {
	return 0, nil
}

func (_ DummyResponseWriter) WriteHeader(int) {}

func (_ DummyResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return nil, nil, nil
}

func (_ DummyResponseWriter) CloseNotify() <-chan bool {
	return nil
}

func (_ DummyResponseWriter) Flush() {}

func TestResponseWriterReset(t *testing.T) {
	w := &responseWriter{}
	w.reset(&DummyResponseWriter{})
	assert.Equal(t, -1, w.size)
	assert.Equal(t, 200, w.status)
	assert.Equal(t, w.ResponseWriter, &DummyResponseWriter{})
}

func TestResponseWriterUnwrap(t *testing.T) {
	w := &responseWriter{}
	w.reset(&DummyResponseWriter{})
	assert.Equal(t, w.Unwrap(), &DummyResponseWriter{})
}

func TestResponseWriterStatus(t *testing.T) {
	w := &responseWriter{}
	w.reset(&DummyResponseWriter{})
	assert.Equal(t, 200, w.Status())
}

func TestResponseWriterSize(t *testing.T) {
	w := &responseWriter{}
	w.reset(&DummyResponseWriter{})
	assert.Equal(t, -1, w.Size())
}

func TestResponseWriterWriteHeaderNow(t *testing.T) {
	w := &responseWriter{}
	w.reset(&DummyResponseWriter{})
	w.WriteHeaderNow()
	assert.Equal(t, 0, w.Size())
}

func TestResponseWriterWritten(t *testing.T) {
	w := &responseWriter{}
	w.reset(&DummyResponseWriter{})
	assert.False(t, w.Written())
	w.size = 1
	assert.True(t, w.Written())
}

func TestResponseWriterWriteString(t *testing.T) {
	w := &responseWriter{}
	w.reset(&DummyResponseWriter{})
	n, err := w.WriteString("")
	assert.Equal(t, 0, n)
	assert.NoError(t, err)
}

func TestResponseWriterHijack(t *testing.T) {
	w := &responseWriter{}
	w.reset(&DummyResponseWriter{})
	_, _, err := w.Hijack()
	assert.NoError(t, err)
}

func TestResponseWriterCloseNotify(t *testing.T) {
	w := &responseWriter{}
	w.reset(&DummyResponseWriter{})
	assert.Nil(t, w.CloseNotify())
}

func TestResponseWriterFlush(t *testing.T) {
	w := &responseWriter{}
	w.reset(&DummyResponseWriter{})
	assert.NotPanics(t, func() {
		w.Flush()
	})
}