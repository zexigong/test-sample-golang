// Copyright 2017 Manu Martinez-Almeida. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

//go:build !nomsgpack

package binding

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMsgpackBindingBindBody(t *testing.T) {
	var s struct {
		Foo string `msgpack:"foo" binding:"required"`
	}
	err := msgpackBinding{}.BindBody(msgpackBodyBytes, &s)
	assert.NoError(t, err)
	assert.Equal(t, "bar", s.Foo)
}

func TestMsgpackBindingBindBodyInvalid(t *testing.T) {
	var s struct {
		Foo string `msgpack:"foo" binding:"required"`
	}
	err := msgpackBinding{}.BindBody(invalidMsgpackBodyBytes, &s)
	assert.Error(t, err)
	assert.Equal(t, "EOF", err.Error())
}

func TestMsgpackBindingBindBodyValidation(t *testing.T) {
	var s struct {
		Foo string `msgpack:"foo" binding:"required"`
	}
	err := msgpackBinding{}.BindBody(msgpackBodyBytesInvalid, &s)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Key: 'Foo' Error:Field validation for 'Foo' failed on the 'required' tag")
}

func TestMsgpackBinding(t *testing.T) {
	testBodyBinding(t, msgpackBinding{}, "msgpack", "/msgpack", msgpackBodyBytes)
}

func TestMsgpackRender(t *testing.T) {
	testRender(t, Msgpack, "msgpack", "/msgpack", msgpackBodyBytes)
}

func BenchmarkMsgpackBindingBindBody(b *testing.B) {
	var s struct {
		Foo string `msgpack:"foo" binding:"required"`
	}
	for i := 0; i < b.N; i++ {
		assert.NoError(b, msgpackBinding{}.BindBody(msgpackBodyBytes, &s))
		assert.Equal(b, "bar", s.Foo)
	}
}

func BenchmarkMsgpackBindingBindBodyInvalid(b *testing.B) {
	var s struct {
		Foo string `msgpack:"foo" binding:"required"`
	}
	for i := 0; i < b.N; i++ {
		assert.Error(b, msgpackBinding{}.BindBody(invalidMsgpackBodyBytes, &s))
	}
}

func BenchmarkMsgpackBindingBindBodyValidation(b *testing.B) {
	var s struct {
		Foo string `msgpack:"foo" binding:"required"`
	}
	for i := 0; i < b.N; i++ {
		assert.Error(b, msgpackBinding{}.BindBody(msgpackBodyBytesInvalid, &s))
	}
}

var msgpackBodyBytes = func() []byte {
	h := new(codec.MsgpackHandle)
	buf := bytes.NewBuffer(nil)
	enc := codec.NewEncoder(buf, h)
	enc.Encode(map[string]string{"foo": "bar"})
	return buf.Bytes()
}()

var msgpackBodyBytesInvalid = func() []byte {
	h := new(codec.MsgpackHandle)
	buf := bytes.NewBuffer(nil)
	enc := codec.NewEncoder(buf, h)
	enc.Encode(map[string]string{"bar": "foo"})
	return buf.Bytes()
}()

var invalidMsgpackBodyBytes = func() []byte {
	h := new(codec.MsgpackHandle)
	buf := bytes.NewBuffer(nil)
	enc := codec.NewEncoder(buf, h)
	enc.Encode("foo")
	return buf.Bytes()
}()