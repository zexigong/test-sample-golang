// Copyright 2017 Manu Martinez-Almeida. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package binding

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMsgpackBindingBindBody(t *testing.T) {
	var s struct {
		Foo string `msgpack:"foo"`
	}
	err := msgpackBinding{}.BindBody(msgpackBody, &s)
	assert.NoError(t, err)
	assert.Equal(t, "bar", s.Foo)
}

func TestMsgpackBindingBindBodyMap(t *testing.T) {
	s := make(map[string]any)
	err := msgpackBinding{}.BindBody(msgpackBody, &s)
	assert.NoError(t, err)
	assert.Equal(t, "bar", s["foo"])
}

func TestMsgpackBindingBindBodySlice(t *testing.T) {
	s := make([]any, 0)
	err := msgpackBinding{}.BindBody(msgpackBodySlice, &s)
	assert.NoError(t, err)
	assert.Len(t, s, 1)
	assert.Equal(t, map[any]any{"foo": "bar"}, s[0])
}

func TestMsgpackBindingBindBodyMapString(t *testing.T) {
	s := make(map[string]string)
	err := msgpackBinding{}.BindBody(msgpackBody, &s)
	assert.Error(t, err)
}

func TestMsgpackBindingBindBodySliceString(t *testing.T) {
	s := make([]string, 0)
	err := msgpackBinding{}.BindBody(msgpackBodySlice, &s)
	assert.Error(t, err)
}

func TestMsgpackBindingBindBodyInvalid(t *testing.T) {
	var s struct {
		Foo string `msgpack:"foo"`
	}
	err := msgpackBinding{}.BindBody(invalidMsgpackBody, &s)
	assert.Error(t, err)
}

func TestMsgpackBindingBindBodyNil(t *testing.T) {
	var s struct {
		Foo string `msgpack:"foo"`
	}
	err := msgpackBinding{}.BindBody(nil, &s)
	assert.Error(t, err)
}

func TestMsgpackBindingBindBodyEmpty(t *testing.T) {
	var s struct {
		Foo string `msgpack:"foo"`
	}
	err := msgpackBinding{}.BindBody([]byte{}, &s)
	assert.Error(t, err)
}

func TestDecodeMsgpack(t *testing.T) {
	var s struct {
		Foo string `msgpack:"foo" binding:"required"`
	}
	err := decodeMsgPack(bytes.NewReader(msgpackBody), &s)
	assert.NoError(t, err)
	assert.Equal(t, "bar", s.Foo)
}

func TestDecodeMsgpackFail(t *testing.T) {
	var s struct {
		Foo string `msgpack:"foo" binding:"required"`
	}
	err := decodeMsgPack(bytes.NewReader(msgpackBodyFail), &s)
	assert.Error(t, err)
}

var msgpackBody = []byte{
	0x81, 0xa3, 0x66, 0x6f, 0x6f, 0xa3, 0x62, 0x61, 0x72,
}

var msgpackBodySlice = []byte{
	0x91, 0x81, 0xa3, 0x66, 0x6f, 0x6f, 0xa3, 0x62, 0x61, 0x72,
}

var msgpackBodyFail = []byte{
	0x81, 0xa3, 0x62, 0x61, 0x72, 0xa3, 0x66, 0x6f, 0x6f,
}

var invalidMsgpackBody = []byte{
	0x81, 0xa3, 0x66, 0x6f, 0x6f,
}