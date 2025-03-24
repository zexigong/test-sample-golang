package fasthttp

import (
	"bytes"
	"io"
	"testing"
)

type mockCloser struct {
	closed bool
}

func (m *mockCloser) Close() error {
	m.closed = true
	return nil
}

func TestUserData_SetGet(t *testing.T) {
	var data userData

	data.Set("key1", "value1")
	if v := data.Get("key1"); v != "value1" {
		t.Fatalf("unexpected value: %v. Expecting %v", v, "value1")
	}

	data.Set("key1", "value2")
	if v := data.Get("key1"); v != "value2" {
		t.Fatalf("unexpected value: %v. Expecting %v", v, "value2")
	}

	data.Set("key2", 123)
	if v := data.Get("key2"); v != 123 {
		t.Fatalf("unexpected value: %v. Expecting %v", v, 123)
	}

	data.Set("key3", nil)
	if v := data.Get("key3"); v != nil {
		t.Fatalf("unexpected value: %v. Expecting nil", v)
	}
}

func TestUserData_SetGetBytes(t *testing.T) {
	var data userData

	key := []byte("key1")
	data.SetBytes(key, "value1")
	if v := data.GetBytes(key); v != "value1" {
		t.Fatalf("unexpected value: %v. Expecting %v", v, "value1")
	}

	data.SetBytes(key, "value2")
	if v := data.GetBytes(key); v != "value2" {
		t.Fatalf("unexpected value: %v. Expecting %v", v, "value2")
	}

	data.SetBytes([]byte("key2"), 123)
	if v := data.GetBytes([]byte("key2")); v != 123 {
		t.Fatalf("unexpected value: %v. Expecting %v", v, 123)
	}

	data.SetBytes([]byte("key3"), nil)
	if v := data.GetBytes([]byte("key3")); v != nil {
		t.Fatalf("unexpected value: %v. Expecting nil", v)
	}
}

func TestUserData_Reset(t *testing.T) {
	var data userData

	data.Set("key1", "value1")
	data.Set("key2", 123)

	data.Reset()

	if v := data.Get("key1"); v != nil {
		t.Fatalf("expected nil, got %v", v)
	}

	if v := data.Get("key2"); v != nil {
		t.Fatalf("expected nil, got %v", v)
	}
}

func TestUserData_Remove(t *testing.T) {
	var data userData

	data.Set("key1", "value1")
	data.Set("key2", 123)

	data.Remove("key1")

	if v := data.Get("key1"); v != nil {
		t.Fatalf("expected nil, got %v", v)
	}

	if v := data.Get("key2"); v != 123 {
		t.Fatalf("unexpected value: %v. Expecting %v", v, 123)
	}
}

func TestUserData_RemoveBytes(t *testing.T) {
	var data userData

	data.SetBytes([]byte("key1"), "value1")
	data.SetBytes([]byte("key2"), 123)

	data.RemoveBytes([]byte("key1"))

	if v := data.GetBytes([]byte("key1")); v != nil {
		t.Fatalf("expected nil, got %v", v)
	}

	if v := data.GetBytes([]byte("key2")); v != 123 {
		t.Fatalf("unexpected value: %v. Expecting %v", v, 123)
	}
}

func TestUserData_ResetWithCloser(t *testing.T) {
	var data userData

	mc := &mockCloser{}
	data.Set("key1", mc)

	data.Reset()

	if !mc.closed {
		t.Fatalf("expected mockCloser to be closed")
	}

	if v := data.Get("key1"); v != nil {
		t.Fatalf("expected nil, got %v", v)
	}
}

func b2s(b []byte) string {
	return string(b)
}