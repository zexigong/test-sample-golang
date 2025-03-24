package fasthttp

import (
	"testing"
)

func TestUserData(t *testing.T) {
	var d userData

	d.Set("foo", "bar")
	d.SetBytes([]byte("foo"), "baz")
	d.Set("bar", "baz")

	v := d.Get("foo")
	if v != "baz" {
		t.Fatalf("unexpected value for key %q: %q. Expecting %q", "foo", v, "baz")
	}

	v = d.GetBytes([]byte("foo"))
	if v != "baz" {
		t.Fatalf("unexpected value for key %q: %q. Expecting %q", "foo", v, "baz")
	}

	d.SetBytes([]byte("foo"), nil)

	v = d.Get("foo")
	if v != nil {
		t.Fatalf("unexpected value for key %q: %q. Expecting %q", "foo", v, nil)
	}

	v = d.GetBytes([]byte("foo"))
	if v != nil {
		t.Fatalf("unexpected value for key %q: %q. Expecting %q", "foo", v, nil)
	}

	d.Set("foo", "bar")
	d.SetBytes([]byte("foo"), "baz")
	d.Set("bar", "baz")

	v = d.Get("foo")
	if v != "baz" {
		t.Fatalf("unexpected value for key %q: %q. Expecting %q", "foo", v, "baz")
	}

	v = d.GetBytes([]byte("foo"))
	if v != "baz" {
		t.Fatalf("unexpected value for key %q: %q. Expecting %q", "foo", v, "baz")
	}

	d.Remove("foo")

	v = d.Get("foo")
	if v != nil {
		t.Fatalf("unexpected value for key %q: %q. Expecting %q", "foo", v, nil)
	}

	v = d.GetBytes([]byte("foo"))
	if v != nil {
		t.Fatalf("unexpected value for key %q: %q. Expecting %q", "foo", v, nil)
	}

	v = d.Get("bar")
	if v != "baz" {
		t.Fatalf("unexpected value for key %q: %q. Expecting %q", "bar", v, "baz")
	}

	v = d.GetBytes([]byte("bar"))
	if v != "baz" {
		t.Fatalf("unexpected value for key %q: %q. Expecting %q", "bar", v, "baz")
	}

	d.RemoveBytes([]byte("bar"))

	v = d.Get("bar")
	if v != nil {
		t.Fatalf("unexpected value for key %q: %q. Expecting %q", "bar", v, nil)
	}

	v = d.GetBytes([]byte("bar"))
	if v != nil {
		t.Fatalf("unexpected value for key %q: %q. Expecting %q", "bar", v, nil)
	}
}