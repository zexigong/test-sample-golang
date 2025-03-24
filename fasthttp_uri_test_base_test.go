package fasthttp

import (
	"bytes"
	"testing"
)

func TestURIReset(t *testing.T) {
	uri := &URI{}
	uri.SetScheme("https")
	uri.SetHost("example.com")
	uri.SetPath("/test/path")
	uri.SetQueryString("param=value")
	uri.SetHash("section")
	uri.SetUsername("user")
	uri.SetPassword("pass")
	uri.DisablePathNormalizing = true

	uri.Reset()

	if len(uri.Scheme()) != 0 {
		t.Fatalf("expected empty scheme, got %q", uri.Scheme())
	}
	if len(uri.Host()) != 0 {
		t.Fatalf("expected empty host, got %q", uri.Host())
	}
	if len(uri.Path()) != 0 {
		t.Fatalf("expected empty path, got %q", uri.Path())
	}
	if len(uri.QueryString()) != 0 {
		t.Fatalf("expected empty query string, got %q", uri.QueryString())
	}
	if len(uri.Hash()) != 0 {
		t.Fatalf("expected empty hash, got %q", uri.Hash())
	}
	if len(uri.Username()) != 0 {
		t.Fatalf("expected empty username, got %q", uri.Username())
	}
	if len(uri.Password()) != 0 {
		t.Fatalf("expected empty password, got %q", uri.Password())
	}
	if uri.DisablePathNormalizing {
		t.Fatal("expected DisablePathNormalizing to be false")
	}
}

func TestURISetAndGet(t *testing.T) {
	uri := &URI{}

	uri.SetScheme("https")
	if !bytes.Equal(uri.Scheme(), []byte("https")) {
		t.Fatalf("expected scheme https, got %q", uri.Scheme())
	}

	uri.SetHost("example.com")
	if !bytes.Equal(uri.Host(), []byte("example.com")) {
		t.Fatalf("expected host example.com, got %q", uri.Host())
	}

	uri.SetPath("/test/path")
	if !bytes.Equal(uri.Path(), []byte("/test/path")) {
		t.Fatalf("expected path /test/path, got %q", uri.Path())
	}

	uri.SetQueryString("param=value")
	if !bytes.Equal(uri.QueryString(), []byte("param=value")) {
		t.Fatalf("expected query string param=value, got %q", uri.QueryString())
	}

	uri.SetHash("section")
	if !bytes.Equal(uri.Hash(), []byte("section")) {
		t.Fatalf("expected hash section, got %q", uri.Hash())
	}

	uri.SetUsername("user")
	if !bytes.Equal(uri.Username(), []byte("user")) {
		t.Fatalf("expected username user, got %q", uri.Username())
	}

	uri.SetPassword("pass")
	if !bytes.Equal(uri.Password(), []byte("pass")) {
		t.Fatalf("expected password pass, got %q", uri.Password())
	}
}

func TestURIFullURI(t *testing.T) {
	uri := &URI{}
	uri.SetScheme("http")
	uri.SetHost("example.com")
	uri.SetPath("/test/path")
	uri.SetQueryString("param=value")
	uri.SetHash("section")

	expected := "http://example.com/test/path?param=value#section"
	if string(uri.FullURI()) != expected {
		t.Fatalf("expected full URI %q, got %q", expected, uri.FullURI())
	}
}

func TestURIParse(t *testing.T) {
	uri := &URI{}
	err := uri.Parse([]byte("example.com"), []byte("/test/path?param=value#section"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !bytes.Equal(uri.Host(), []byte("example.com")) {
		t.Fatalf("expected host example.com, got %q", uri.Host())
	}

	if !bytes.Equal(uri.Path(), []byte("/test/path")) {
		t.Fatalf("expected path /test/path, got %q", uri.Path())
	}

	if !bytes.Equal(uri.QueryString(), []byte("param=value")) {
		t.Fatalf("expected query string param=value, got %q", uri.QueryString())
	}

	if !bytes.Equal(uri.Hash(), []byte("section")) {
		t.Fatalf("expected hash section, got %q", uri.Hash())
	}
}

func TestURIUpdate(t *testing.T) {
	uri := &URI{}
	uri.SetScheme("http")
	uri.SetHost("example.com")
	uri.SetPath("/test/path")

	uri.Update("https://newexample.com/newpath?newparam=newvalue#newsection")

	expected := "https://newexample.com/newpath?newparam=newvalue#newsection"
	if string(uri.FullURI()) != expected {
		t.Fatalf("expected updated URI %q, got %q", expected, uri.FullURI())
	}
}

func TestURIRequestURI(t *testing.T) {
	uri := &URI{}
	uri.SetPath("/test/path")
	uri.SetQueryString("param=value")

	expected := "/test/path?param=value"
	if string(uri.RequestURI()) != expected {
		t.Fatalf("expected request URI %q, got %q", expected, uri.RequestURI())
	}
}