package binding

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/ugorji/go/codec"
	"github.com/stretchr/testify/assert"
)

type testStruct struct {
	Name string `msgpack:"name"`
	Age  int    `msgpack:"age"`
}

func TestMsgpackBinding_Name(t *testing.T) {
	b := msgpackBinding{}
	assert.Equal(t, "msgpack", b.Name())
}

func TestMsgpackBinding_Bind(t *testing.T) {
	b := msgpackBinding{}

	obj := testStruct{}
	data := []byte{0x82, 0xa4, 0x6e, 0x61, 0x6d, 0x65, 0xa4, 0x4a, 0x6f, 0x68, 0x6e, 0xa3, 0x61, 0x67, 0x65, 0x1e}

	req, err := http.NewRequest("POST", "http://example.com", bytes.NewBuffer(data))
	assert.NoError(t, err)

	err = b.Bind(req, &obj)
	assert.NoError(t, err)
	assert.Equal(t, "John", obj.Name)
	assert.Equal(t, 30, obj.Age)
}

func TestMsgpackBinding_BindBody(t *testing.T) {
	b := msgpackBinding{}

	obj := testStruct{}
	data := []byte{0x82, 0xa4, 0x6e, 0x61, 0x6d, 0x65, 0xa4, 0x41, 0x6e, 0x6e, 0x65, 0xa3, 0x61, 0x67, 0x65, 0x25}

	err := b.BindBody(data, &obj)
	assert.NoError(t, err)
	assert.Equal(t, "Anne", obj.Name)
	assert.Equal(t, 37, obj.Age)
}

func TestDecodeMsgPack_Error(t *testing.T) {
	var obj testStruct
	data := []byte{0xc1} // Invalid Msgpack data

	err := decodeMsgPack(bytes.NewReader(data), &obj)
	assert.Error(t, err)
}