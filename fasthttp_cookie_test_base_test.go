package fasthttp

import (
	"bytes"
	"errors"
	"io"
	"testing"
	"time"
)

func TestCookieAcquireRelease(t *testing.T) {
	c := AcquireCookie()
	if c == nil {
		t.Fatalf("AcquireCookie() returned nil")
	}
	ReleaseCookie(c)
}

func TestCookieCopyTo(t *testing.T) {
	c1 := &Cookie{}
	c1.SetKey("testKey")
	c1.SetValue("testValue")
	c1.SetExpire(time.Now().Add(10 * time.Hour))
	c1.SetDomain("example.com")
	c1.SetPath("/path")
	c1.SetHTTPOnly(true)
	c1.SetSecure(true)
	c1.SetSameSite(CookieSameSiteStrictMode)
	c1.SetPartitioned(true)

	c2 := &Cookie{}
	c1.CopyTo(c2)

	if !bytes.Equal(c1.Key(), c2.Key()) || !bytes.Equal(c1.Value(), c2.Value()) ||
		!c1.Expire().Equal(c2.Expire()) || !bytes.Equal(c1.Domain(), c2.Domain()) ||
		!bytes.Equal(c1.Path(), c2.Path()) || c1.HTTPOnly() != c2.HTTPOnly() ||
		c1.Secure() != c2.Secure() || c1.SameSite() != c2.SameSite() ||
		c1.Partitioned() != c2.Partitioned() {
		t.Fatalf("CopyTo failed: %v vs %v", c1, c2)
	}
}

func TestCookieSetGet(t *testing.T) {
	c := &Cookie{}

	c.SetKey("testKey")
	if !bytes.Equal(c.Key(), []byte("testKey")) {
		t.Fatalf("unexpected key: %q", c.Key())
	}

	c.SetValue("testValue")
	if !bytes.Equal(c.Value(), []byte("testValue")) {
		t.Fatalf("unexpected value: %q", c.Value())
	}

	expire := time.Now().Add(10 * time.Hour)
	c.SetExpire(expire)
	if !c.Expire().Equal(expire) {
		t.Fatalf("unexpected expire: %v", c.Expire())
	}

	c.SetMaxAge(3600)
	if c.MaxAge() != 3600 {
		t.Fatalf("unexpected maxAge: %d", c.MaxAge())
	}

	c.SetDomain("example.com")
	if !bytes.Equal(c.Domain(), []byte("example.com")) {
		t.Fatalf("unexpected domain: %q", c.Domain())
	}

	c.SetPath("/path")
	if !bytes.Equal(c.Path(), []byte("/path")) {
		t.Fatalf("unexpected path: %q", c.Path())
	}

	c.SetHTTPOnly(true)
	if !c.HTTPOnly() {
		t.Fatalf("expected httpOnly to be true")
	}

	c.SetSecure(true)
	if !c.Secure() {
		t.Fatalf("expected secure to be true")
	}

	c.SetSameSite(CookieSameSiteStrictMode)
	if c.SameSite() != CookieSameSiteStrictMode {
		t.Fatalf("unexpected SameSite: %d", c.SameSite())
	}

	c.SetPartitioned(true)
	if !c.Partitioned() {
		t.Fatalf("expected partitioned to be true")
	}
}

func TestCookieParse(t *testing.T) {
	c := &Cookie{}
	err := c.Parse("testKey=testValue; Domain=example.com; Path=/path; Max-Age=3600; HttpOnly; Secure; SameSite=Strict; Partitioned")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !bytes.Equal(c.Key(), []byte("testKey")) {
		t.Fatalf("unexpected key: %q", c.Key())
	}
	if !bytes.Equal(c.Value(), []byte("testValue")) {
		t.Fatalf("unexpected value: %q", c.Value())
	}
	if !bytes.Equal(c.Domain(), []byte("example.com")) {
		t.Fatalf("unexpected domain: %q", c.Domain())
	}
	if !bytes.Equal(c.Path(), []byte("/path")) {
		t.Fatalf("unexpected path: %q", c.Path())
	}
	if c.MaxAge() != 3600 {
		t.Fatalf("unexpected maxAge: %d", c.MaxAge())
	}
	if !c.HTTPOnly() {
		t.Fatalf("expected httpOnly to be true")
	}
	if !c.Secure() {
		t.Fatalf("expected secure to be true")
	}
	if c.SameSite() != CookieSameSiteStrictMode {
		t.Fatalf("unexpected SameSite: %d", c.SameSite())
	}
	if !c.Partitioned() {
		t.Fatalf("expected partitioned to be true")
	}
}

func TestCookieParseErrors(t *testing.T) {
	c := &Cookie{}

	err := c.Parse("invalid_cookie")
	if !errors.Is(err, errNoCookies) {
		t.Fatalf("unexpected error: %v", err)
	}

	err = c.Parse("Max-Age=invalid")
	if err == nil {
		t.Fatalf("expected error but got none")
	}
}

func TestCookieAppendBytes(t *testing.T) {
	c := &Cookie{}
	c.SetKey("testKey")
	c.SetValue("testValue")
	c.SetDomain("example.com")
	c.SetPath("/path")
	c.SetHTTPOnly(true)
	c.SetSecure(true)
	c.SetSameSite(CookieSameSiteStrictMode)
	c.SetPartitioned(true)

	expected := "testKey=testValue; Domain=example.com; Path=/path; HttpOnly; Secure; SameSite=Strict; Partitioned"
	result := string(c.AppendBytes(nil))
	if result != expected {
		t.Fatalf("unexpected cookie representation: %s", result)
	}
}

func TestCookieWriteTo(t *testing.T) {
	c := &Cookie{}
	c.SetKey("testKey")
	c.SetValue("testValue")

	var buf bytes.Buffer
	n, err := c.WriteTo(&buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != int64(buf.Len()) {
		t.Fatalf("unexpected written length: %d", n)
	}
	expected := "testKey=testValue"
	if buf.String() != expected {
		t.Fatalf("unexpected buffer content: %s", buf.String())
	}
}