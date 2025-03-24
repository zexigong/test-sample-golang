package fasthttp

import (
	"bytes"
	"io"
	"reflect"
	"testing"
)

func TestAcquireReleaseArgs(t *testing.T) {
	a := AcquireArgs()
	if a == nil {
		t.Fatalf("Expected non-nil Args")
	}

	ReleaseArgs(a)
}

func TestArgsReset(t *testing.T) {
	a := &Args{}
	a.Add("key1", "value1")
	a.Reset()
	if a.Len() != 0 {
		t.Fatalf("Expected args to be empty after reset, got %d", a.Len())
	}
}

func TestArgsCopyTo(t *testing.T) {
	src := &Args{}
	src.Add("key1", "value1")

	dst := &Args{}
	src.CopyTo(dst)

	if dst.Len() != 1 {
		t.Fatalf("Expected dst to have 1 arg, got %d", dst.Len())
	}
	if !bytes.Equal(dst.Peek("key1"), []byte("value1")) {
		t.Fatalf("Expected dst to have value 'value1', got %s", dst.Peek("key1"))
	}
}

func TestArgsVisitAll(t *testing.T) {
	a := &Args{}
	a.Add("key1", "value1")
	a.Add("key2", "value2")

	result := make(map[string]string)
	a.VisitAll(func(k, v []byte) {
		result[string(k)] = string(v)
	})

	expected := map[string]string{"key1": "value1", "key2": "value2"}
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("Expected %v, got %v", expected, result)
	}
}

func TestArgsParse(t *testing.T) {
	a := &Args{}
	a.Parse("key1=value1&key2=value2")

	if a.Len() != 2 {
		t.Fatalf("Expected 2 args, got %d", a.Len())
	}
	if !bytes.Equal(a.Peek("key1"), []byte("value1")) {
		t.Fatalf("Expected value 'value1', got %s", a.Peek("key1"))
	}
}

func TestArgsString(t *testing.T) {
	a := &Args{}
	a.Add("key1", "value1")
	a.Add("key2", "value2")

	expected := "key1=value1&key2=value2"
	if a.String() != expected {
		t.Fatalf("Expected %s, got %s", expected, a.String())
	}
}

func TestArgsSort(t *testing.T) {
	a := &Args{}
	a.Add("key2", "value2")
	a.Add("key1", "value1")

	a.Sort(bytes.Compare)

	expected := "key1=value1&key2=value2"
	if a.String() != expected {
		t.Fatalf("Expected %s, got %s", expected, a.String())
	}
}

func TestArgsAppendBytes(t *testing.T) {
	a := &Args{}
	a.Add("key1", "value1")
	a.Add("key2", "value2")

	dst := a.AppendBytes(nil)
	expected := []byte("key1=value1&key2=value2")
	if !bytes.Equal(dst, expected) {
		t.Fatalf("Expected %s, got %s", expected, dst)
	}
}

func TestArgsWriteTo(t *testing.T) {
	a := &Args{}
	a.Add("key1", "value1")
	a.Add("key2", "value2")

	var buf bytes.Buffer
	n, err := a.WriteTo(&buf)
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}

	expected := "key1=value1&key2=value2"
	if buf.String() != expected {
		t.Fatalf("Expected %s, got %s", expected, buf.String())
	}

	if n != int64(len(expected)) {
		t.Fatalf("Expected %d bytes written, got %d", len(expected), n)
	}
}

func TestArgsDel(t *testing.T) {
	a := &Args{}
	a.Add("key1", "value1")
	a.Del("key1")

	if a.Len() != 0 {
		t.Fatalf("Expected 0 args, got %d", a.Len())
	}
}

func TestArgsAddAndSet(t *testing.T) {
	a := &Args{}
	a.Add("key1", "value1")
	a.Set("key1", "value2")

	if !bytes.Equal(a.Peek("key1"), []byte("value2")) {
		t.Fatalf("Expected value 'value2', got %s", a.Peek("key1"))
	}
}

func TestArgsPeek(t *testing.T) {
	a := &Args{}
	a.Add("key1", "value1")

	if !bytes.Equal(a.Peek("key1"), []byte("value1")) {
		t.Fatalf("Expected value 'value1', got %s", a.Peek("key1"))
	}

	if a.Peek("key2") != nil {
		t.Fatalf("Expected nil, got %s", a.Peek("key2"))
	}
}

func TestArgsHas(t *testing.T) {
	a := &Args{}
	a.Add("key1", "value1")

	if !a.Has("key1") {
		t.Fatalf("Expected key 'key1' to exist")
	}

	if a.Has("key2") {
		t.Fatalf("Expected key 'key2' to not exist")
	}
}

func TestArgsGetUint(t *testing.T) {
	a := &Args{}
	a.Add("key1", "123")

	value, err := a.GetUint("key1")
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}
	if value != 123 {
		t.Fatalf("Expected value 123, got %d", value)
	}

	value, err = a.GetUint("key2")
	if err == nil {
		t.Fatalf("Expected error, got nil")
	}
	if value != -1 {
		t.Fatalf("Expected value -1, got %d", value)
	}
}

func TestArgsSetUint(t *testing.T) {
	a := &Args{}
	a.SetUint("key1", 123)

	if !bytes.Equal(a.Peek("key1"), []byte("123")) {
		t.Fatalf("Expected value '123', got %s", a.Peek("key1"))
	}
}

func TestArgsGetUfloat(t *testing.T) {
	a := &Args{}
	a.Add("key1", "123.45")

	value, err := a.GetUfloat("key1")
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}
	if value != 123.45 {
		t.Fatalf("Expected value 123.45, got %f", value)
	}

	value, err = a.GetUfloat("key2")
	if err == nil {
		t.Fatalf("Expected error, got nil")
	}
	if value != -1 {
		t.Fatalf("Expected value -1, got %f", value)
	}
}

func TestArgsGetBool(t *testing.T) {
	a := &Args{}
	a.Add("key1", "true")
	a.Add("key2", "0")

	if !a.GetBool("key1") {
		t.Fatalf("Expected true, got false")
	}

	if a.GetBool("key2") {
		t.Fatalf("Expected false, got true")
	}
}