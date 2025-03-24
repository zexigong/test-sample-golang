// Copyright 2014 Manu Martinez-Almeida. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package gin

import (
	"bufio"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

var _ ResponseWriter = &responseWriter{}

func TestResponseWriterUnwrap(t *testing.T) {
	w := httptest.NewRecorder()
	rw := &responseWriter{ResponseWriter: w}

	assert.Equal(t, w, rw.Unwrap())
}

func TestResponseWriterReset(t *testing.T) {
	w := httptest.NewRecorder()
	rw := &responseWriter{}

	rw.reset(w)
	assert.Equal(t, w, rw.ResponseWriter)
	assert.Equal(t, noWritten, rw.size)
	assert.Equal(t, defaultStatus, rw.status)
}

func TestResponseWriterWriteHeader(t *testing.T) {
	w := httptest.NewRecorder()
	rw := &responseWriter{ResponseWriter: w}

	rw.WriteHeader(300)
	assert.Equal(t, 300, rw.status)
	assert.False(t, rw.Written())

	rw.WriteHeader(200)
	assert.Equal(t, 300, rw.status)
}

func TestResponseWriterWriteHeadersNow(t *testing.T) {
	w := httptest.NewRecorder()
	rw := &responseWriter{ResponseWriter: w}

	rw.WriteHeaderNow()

	assert.Equal(t, defaultStatus, rw.status)
	assert.Equal(t, 0, rw.size)
	assert.True(t, rw.Written())
}

func TestResponseWriterWrite(t *testing.T) {
	w := httptest.NewRecorder()
	rw := &responseWriter{ResponseWriter: w}

	n, err := rw.Write([]byte("hola"))
	assert.Equal(t, 4, n)
	assert.Equal(t, 4, rw.size)
	assert.NoError(t, err)
	assert.Equal(t, defaultStatus, rw.status)
	assert.True(t, rw.Written())

	n, err = rw.Write([]byte(" adios"))
	assert.Equal(t, 6, n)
	assert.Equal(t, 10, rw.size)
	assert.NoError(t, err)
	assert.Equal(t, defaultStatus, rw.status)
	assert.True(t, rw.Written())
}

func TestResponseWriterWriteString(t *testing.T) {
	w := httptest.NewRecorder()
	rw := &responseWriter{ResponseWriter: w}

	n, err := rw.WriteString("hola")
	assert.Equal(t, 4, n)
	assert.Equal(t, 4, rw.size)
	assert.NoError(t, err)
	assert.Equal(t, defaultStatus, rw.status)
	assert.True(t, rw.Written())

	n, err = rw.WriteString(" adios")
	assert.Equal(t, 6, n)
	assert.Equal(t, 10, rw.size)
	assert.NoError(t, err)
	assert.Equal(t, defaultStatus, rw.status)
	assert.True(t, rw.Written())
}

type hijackableResponseWriter struct {
	http.ResponseWriter
}

func (w *hijackableResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return nil, nil, nil
}

func TestResponseWriterHijack(t *testing.T) {
	w := &hijackableResponseWriter{httptest.NewRecorder()}
	rw := &responseWriter{ResponseWriter: w}

	assert.Equal(t, noWritten, rw.size)
	assert.Equal(t, defaultStatus, rw.status)

	conn, buf, err := rw.Hijack()
	assert.Nil(t, conn)
	assert.Nil(t, buf)
	assert.NoError(t, err)
	assert.Equal(t, 0, rw.size)
	assert.Equal(t, defaultStatus, rw.status)
}

func TestResponseWriterFlush(t *testing.T) {
	w := httptest.NewRecorder()
	rw := &responseWriter{ResponseWriter: w}

	rw.Flush()
	assert.Equal(t, 0, rw.size)
	assert.Equal(t, defaultStatus, rw.status)
}

func TestResponseWriterStatus(t *testing.T) {
	w := httptest.NewRecorder()
	rw := &responseWriter{ResponseWriter: w}

	assert.Equal(t, defaultStatus, rw.Status())

	rw.WriteHeader(300)
	assert.Equal(t, 300, rw.Status())
}

func TestResponseWriterSize(t *testing.T) {
	w := httptest.NewRecorder()
	rw := &responseWriter{ResponseWriter: w}

	assert.Equal(t, noWritten, rw.Size())

	rw.size = 10
	assert.Equal(t, 10, rw.Size())
}

func TestResponseWriterWritten(t *testing.T) {
	w := httptest.NewRecorder()
	rw := &responseWriter{ResponseWriter: w}

	assert.False(t, rw.Written())

	rw.size = 10
	assert.True(t, rw.Written())
}

func TestResponseWriterPusher(t *testing.T) {
	w := httptest.NewRecorder()
	rw := &responseWriter{ResponseWriter: w}

	assert.Nil(t, rw.Pusher())
}

func TestResponseWriterCloseNotify(t *testing.T) {
	w := &closeNotifyingResponseWriter{httptest.NewRecorder()}
	rw := &responseWriter{ResponseWriter: w}

	assert.NotNil(t, rw.CloseNotify())
}

type closeNotifyingResponseWriter struct {
	http.ResponseWriter
}

func (c *closeNotifyingResponseWriter) CloseNotify() <-chan bool {
	return nil
}

type flushableResponseWriter struct {
	http.ResponseWriter
}

func (w *flushableResponseWriter) Flush() {
}

func TestResponseWriterFlushWithFlushableResponseWriter(t *testing.T) {
	w := &flushableResponseWriter{httptest.NewRecorder()}
	rw := &responseWriter{ResponseWriter: w}

	rw.Flush()
	assert.Equal(t, 0, rw.size)
	assert.Equal(t, defaultStatus, rw.status)
}

type pushableResponseWriter struct {
	http.ResponseWriter
}

func (w *pushableResponseWriter) Push(target string, opts *http.PushOptions) error {
	return nil
}

func TestResponseWriterPusherWithPushableResponseWriter(t *testing.T) {
	w := &pushableResponseWriter{httptest.NewRecorder()}
	rw := &responseWriter{ResponseWriter: w}

	assert.NotNil(t, rw.Pusher())
}

func TestResponseWriterCloseNotifyWithCloseNotifyingResponseWriter(t *testing.T) {
	w := &closeNotifyingResponseWriter{httptest.NewRecorder()}
	rw := &responseWriter{ResponseWriter: w}

	assert.NotNil(t, rw.CloseNotify())
}

type notHijackableResponseWriter struct {
	http.ResponseWriter
}

func TestResponseWriterHijackWithoutHijackableResponseWriter(t *testing.T) {
	w := &notHijackableResponseWriter{httptest.NewRecorder()}
	rw := &responseWriter{ResponseWriter: w}

	conn, buf, err := rw.Hijack()
	assert.Nil(t, conn)
	assert.Nil(t, buf)
	assert.Error(t, err)
}

type notFlushableResponseWriter struct {
	http.ResponseWriter
}

func TestResponseWriterFlushWithoutFlushableResponseWriter(t *testing.T) {
	w := &notFlushableResponseWriter{httptest.NewRecorder()}
	rw := &responseWriter{ResponseWriter: w}

	rw.Flush()
	assert.Equal(t, 0, rw.size)
	assert.Equal(t, defaultStatus, rw.status)
}

type notPushableResponseWriter struct {
	http.ResponseWriter
}

func TestResponseWriterPusherWithoutPushableResponseWriter(t *testing.T) {
	w := &notPushableResponseWriter{httptest.NewRecorder()}
	rw := &responseWriter{ResponseWriter: w}

	assert.Nil(t, rw.Pusher())
}

type notCloseNotifyingResponseWriter struct {
	http.ResponseWriter
}

func TestResponseWriterCloseNotifyWithoutCloseNotifyingResponseWriter(t *testing.T) {
	w := &notCloseNotifyingResponseWriter{httptest.NewRecorder()}
	rw := &responseWriter{ResponseWriter: w}

	assert.Nil(t, rw.CloseNotify())
}

type writerWithoutWriteHeader struct {
	io.Writer
}

func (w *writerWithoutWriteHeader) Header() http.Header {
	return http.Header{}
}

func (w *writerWithoutWriteHeader) WriteHeader(code int) {
}

func TestResponseWriterWriteHeaderWithWriterWithoutWriteHeader(t *testing.T) {
	w := &writerWithoutWriteHeader{}
	rw := &responseWriter{ResponseWriter: w}

	assert.NotPanics(t, func() {
		rw.WriteHeader(200)
	})
}

func TestResponseWriterWriteHeaderNowWithWriterWithoutWriteHeader(t *testing.T) {
	w := &writerWithoutWriteHeader{}
	rw := &responseWriter{ResponseWriter: w}

	assert.NotPanics(t, func() {
		rw.WriteHeaderNow()
	})
}

func TestResponseWriterWriteWithWriterWithoutWriteHeader(t *testing.T) {
	w := &writerWithoutWriteHeader{}
	rw := &responseWriter{ResponseWriter: w}

	assert.NotPanics(t, func() {
		rw.Write(nil)
	})
}

func TestResponseWriterWriteStringWithWriterWithoutWriteHeader(t *testing.T) {
	w := &writerWithoutWriteHeader{}
	rw := &responseWriter{ResponseWriter: w}

	assert.NotPanics(t, func() {
		rw.WriteString("")
	})
}