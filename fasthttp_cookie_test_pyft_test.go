package fasthttp

import (
	"testing"
	"time"
)

func TestAcquireReleaseCookie(t *testing.T) {
	c := AcquireCookie()
	defer ReleaseCookie(c)
}

func TestCookie(t *testing.T) {
	var c Cookie

	// empty cookie must have empty key, value and expire unlimited
	testCookie(t, &c, "", "", CookieExpireUnlimited)

	c.SetKey("foo")
	testCookie(t, &c, "foo", "", CookieExpireUnlimited)
	c.SetValue("bar")
	testCookie(t, &c, "foo", "bar", CookieExpireUnlimited)
	c.SetExpire(CookieExpireDelete)
	testCookie(t, &c, "foo", "bar", CookieExpireDelete)
	c.SetKey("aaa")
	testCookie(t, &c, "aaa", "bar", CookieExpireDelete)
	c.SetValue("bbb")
	testCookie(t, &c, "aaa", "bbb", CookieExpireDelete)
	c.SetExpire(CookieExpireUnlimited)
	testCookie(t, &c, "aaa", "bbb", CookieExpireUnlimited)

	c.Reset()
	testCookie(t, &c, "", "", CookieExpireUnlimited)
}

func TestCookieSecureHttpOnly(t *testing.T) {
	var c Cookie

	if c.Secure() {
		t.Fatalf("cookie.Secure must be false")
	}
	if c.HTTPOnly() {
		t.Fatalf("cookie.HTTPOnly must be false")
	}

	c.SetSecure(true)
	c.SetHTTPOnly(true)

	if !c.Secure() {
		t.Fatalf("cookie.Secure must be true")
	}
	if !c.HTTPOnly() {
		t.Fatalf("cookie.HTTPOnly must be true")
	}

	c.SetSecure(false)
	c.SetHTTPOnly(false)

	if c.Secure() {
		t.Fatalf("cookie.Secure must be false")
	}
	if c.HTTPOnly() {
		t.Fatalf("cookie.HTTPOnly must be false")
	}
}

func TestCookieSameSite(t *testing.T) {
	var c Cookie

	if c.SameSite() != CookieSameSiteDisabled {
		t.Fatalf("cookie.SameSite must be disabled by default")
	}

	c.SetSameSite(CookieSameSiteLaxMode)
	if c.SameSite() != CookieSameSiteLaxMode {
		t.Fatalf("cookie.SameSite must be lax")
	}

	c.SetSameSite(CookieSameSiteStrictMode)
	if c.SameSite() != CookieSameSiteStrictMode {
		t.Fatalf("cookie.SameSite must be strict")
	}

	c.SetSameSite(CookieSameSiteNoneMode)
	if c.SameSite() != CookieSameSiteNoneMode {
		t.Fatalf("cookie.SameSite must be none")
	}
	if !c.Secure() {
		t.Fatalf("cookie.Secure must be true when cookie.SameSite is none")
	}

	c.SetSecure(false)
	c.SetSameSite(CookieSameSiteDefaultMode)
	if c.SameSite() != CookieSameSiteDefaultMode {
		t.Fatalf("cookie.SameSite must be default")
	}
}

func TestCookiePartitioned(t *testing.T) {
	var c Cookie

	if c.Partitioned() {
		t.Fatalf("cookie.Partitioned must be false by default")
	}

	c.SetPartitioned(true)
	if !c.Partitioned() {
		t.Fatalf("cookie.Partitioned must be true")
	}
	if !c.Secure() {
		t.Fatalf("cookie.Secure must be true when cookie.Partitioned is true")
	}
	if string(c.Path()) != "/" {
		t.Fatalf("cookie.Path must be '/' when cookie.Partitioned is true")
	}

	c.SetSecure(false)
	c.SetPartitioned(false)
	if c.Partitioned() {
		t.Fatalf("cookie.Partitioned must be false")
	}
}

func TestCookiePath(t *testing.T) {
	var c Cookie

	if len(c.Path()) > 0 {
		t.Fatalf("cookie.Path must be empty")
	}

	c.SetPath("/foo/bar")
	if string(c.Path()) != "/foo/bar" {
		t.Fatalf("unexpected cookie.Path %q. Expecting %q", c.Path(), "/foo/bar")
	}

	c.SetPath("")
	if len(c.Path()) > 0 {
		t.Fatalf("cookie.Path must be empty")
	}

	c.SetPath("/aaabbb")
	if string(c.Path()) != "/aaabbb" {
		t.Fatalf("unexpected cookie.Path %q. Expecting %q", c.Path(), "/aaabbb")
	}
}

func TestCookieDomain(t *testing.T) {
	var c Cookie

	if len(c.Domain()) > 0 {
		t.Fatalf("cookie.Domain must be empty")
	}

	c.SetDomain("foobar.com")
	if string(c.Domain()) != "foobar.com" {
		t.Fatalf("unexpected cookie.Domain %q. Expecting %q", c.Domain(), "foobar.com")
	}

	c.SetDomain("")
	if len(c.Domain()) > 0 {
		t.Fatalf("cookie.Domain must be empty")
	}

	c.SetDomain("aaabbb.com")
	if string(c.Domain()) != "aaabbb.com" {
		t.Fatalf("unexpected cookie.Domain %q. Expecting %q", c.Domain(), "aaabbb.com")
	}
}

func TestResponseCookieMaxAge(t *testing.T) {
	var c Cookie

	if c.MaxAge() != 0 {
		t.Fatalf("cookie.MaxAge must be empty")
	}

	c.SetMaxAge(3600)
	if c.MaxAge() != 3600 {
		t.Fatalf("unexpected cookie.MaxAge %q. Expecting %q", c.MaxAge(), 3600)
	}

	c.SetMaxAge(0)
	if c.MaxAge() != 0 {
		t.Fatalf("cookie.MaxAge must be empty")
	}

	c.SetMaxAge(1200)
	if c.MaxAge() != 1200 {
		t.Fatalf("unexpected cookie.MaxAge %q. Expecting %q", c.MaxAge(), 1200)
	}
}

func TestParseResponseCookie(t *testing.T) {
	t.Parallel()

	var c Cookie
	testParseResponseCookie(t, &c, "foo=bar; Path=/; Domain=aaa.com; Max-Age=150; Secure; HttpOnly",
		"foo", "bar", "/",
		"aaa.com", 150, true, true, CookieExpireUnlimited, CookieSameSiteDisabled, false)
	testParseResponseCookie(t, &c, "httpOnly=1; path=/foo/bar; domain=bar.com; max-age=100; httponly; expires=Mon, 02-Jan-2006 15:04:05 MST",
		"httpOnly", "1", "/foo/bar",
		"bar.com", 100, false, true, time.Date(2006, 1, 2, 15, 4, 5, 0, time.FixedZone("MST", -7*60*60)), CookieSameSiteDisabled, false)
	testParseResponseCookie(t, &c, "foo=bar; Path=/; HttpOnly; SameSite",
		"foo", "bar", "/", "", 0, false, true, CookieExpireUnlimited, CookieSameSiteDefaultMode, false)
	testParseResponseCookie(t, &c, "foo=bar; Path=/; HttpOnly; SameSite=Lax",
		"foo", "bar", "/", "", 0, false, true, CookieExpireUnlimited, CookieSameSiteLaxMode, false)
	testParseResponseCookie(t, &c, "foo=bar; Path=/; HttpOnly; SameSite=Strict",
		"foo", "bar", "/", "", 0, false, true, CookieExpireUnlimited, CookieSameSiteStrictMode, false)
	testParseResponseCookie(t, &c, "foo=bar; Path=/; HttpOnly; SameSite=None",
		"foo", "bar", "/", "", 0, true, true, CookieExpireUnlimited, CookieSameSiteNoneMode, false)
	testParseResponseCookie(t, &c, "foo=bar; Path=/foo/bar; Secure; Partitioned",
		"foo", "bar", "/foo/bar", "", 0, true, false, CookieExpireUnlimited, CookieSameSiteDisabled, true)

	// Invalid cookie expires format.
	testParseResponseCookie(t, &c, "foo=bar; Path=/; expires=invalid",
		"foo", "bar", "/", "", 0, false, false, CookieExpireUnlimited, CookieSameSiteDisabled, false)

	// Invalid cookie expires format.
	testParseResponseCookie(t, &c, "foo=bar; Path=/; Expires=Thu, 01-Jan-1970 00:00:00 GMT",
		"foo", "bar", "/", "", 0, false, false, time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC), CookieSameSiteDisabled, false)
}

func testParseResponseCookie(t *testing.T, c *Cookie, s, expectedKey, expectedValue, expectedPath, expectedDomain string, expectedMaxAge int, expectedSecure, expectedHTTPOnly bool, expectedExpire time.Time, expectedSameSite CookieSameSite, expectedPartitioned bool) {
	c.Parse(s)
	if string(c.Key()) != expectedKey {
		t.Fatalf("unexpected cookie key %q. Expecting %q. cookie %q", c.Key(), expectedKey, s)
	}
	if string(c.Value()) != expectedValue {
		t.Fatalf("unexpected cookie value %q. Expecting %q. cookie %q", c.Value(), expectedValue, s)
	}
	if string(c.Path()) != expectedPath {
		t.Fatalf("unexpected cookie path %q. Expecting %q. cookie %q", c.Path(), expectedPath, s)
	}
	if string(c.Domain()) != expectedDomain {
		t.Fatalf("unexpected cookie domain %q. Expecting %q. cookie %q", c.Domain(), expectedDomain, s)
	}
	if c.MaxAge() != expectedMaxAge {
		t.Fatalf("unexpected cookie max-age %d. Expecting %d. cookie %q", c.MaxAge(), expectedMaxAge, s)
	}
	if c.Secure() != expectedSecure {
		t.Fatalf("unexpected cookie secure flag %v. Expecting %v. cookie %q", c.Secure(), expectedSecure, s)
	}
	if c.HTTPOnly() != expectedHTTPOnly {
		t.Fatalf("unexpected cookie httpOnly flag %v. Expecting %v. cookie %q", c.HTTPOnly(), expectedHTTPOnly, s)
	}
	if c.Expire() != expectedExpire {
		t.Fatalf("unexpected cookie expire %v. Expecting %v. cookie %q", c.Expire(), expectedExpire, s)
	}
	if c.SameSite() != expectedSameSite {
		t.Fatalf("unexpected cookie SameSite %v. Expecting %v. cookie %q", c.SameSite(), expectedSameSite, s)
	}
	if c.Partitioned() != expectedPartitioned {
		t.Fatalf("unexpected cookie Partitioned %v. Expecting %v. cookie %q", c.Partitioned(), expectedPartitioned, s)
	}
}

func TestParseRequestCookie(t *testing.T) {
	t.Parallel()

	var c Cookie
	testParseRequestCookie(t, &c, "foo=bar", "foo", "bar")
	testParseRequestCookie(t, &c, "key=value", "key", "value")
	testParseRequestCookie(t, &c, "key=1", "key", "1")
	testParseRequestCookie(t, &c, "key=", "key", "")
	testParseRequestCookie(t, &c, "=value", "", "value")
	testParseRequestCookie(t, &c, "key", "key", "")
	testParseRequestCookie(t, &c, "key=", "key", "")
}

func testParseRequestCookie(t *testing.T, c *Cookie, s, expectedKey, expectedValue string) {
	c.Parse(s)
	if string(c.Key()) != expectedKey {
		t.Fatalf("unexpected cookie key %q. Expecting %q. cookie %q", c.Key(), expectedKey, s)
	}
	if string(c.Value()) != expectedValue {
		t.Fatalf("unexpected cookie value %q. Expecting %q. cookie %q", c.Value(), expectedValue, s)
	}
}

func TestRequestCookieKeyValue(t *testing.T) {
	t.Parallel()

	testRequestCookieKeyValue(t, "foo=bar", "foo", "bar")
	testRequestCookieKeyValue(t, "key=value", "key", "value")
	testRequestCookieKeyValue(t, "key=1", "key", "1")
	testRequestCookieKeyValue(t, "key=", "key", "")
	testRequestCookieKeyValue(t, "=value", "", "value")
	testRequestCookieKeyValue(t, "key", "key", "")
	testRequestCookieKeyValue(t, "key=", "key", "")
	testRequestCookieKeyValue(t, "key=; key2=abc", "key2", "abc")
}

func testRequestCookieKeyValue(t *testing.T, s, expectedKey, expectedValue string) {
	k, v := ParseRequestCookie([]byte(s))
	if string(k) != expectedKey {
		t.Fatalf("unexpected cookie key %q. Expecting %q. cookie %q", k, expectedKey, s)
	}
	if string(v) != expectedValue {
		t.Fatalf("unexpected cookie value %q. Expecting %q. cookie %q", v, expectedValue, s)
	}
}

func TestResponseCookieKeyValue(t *testing.T) {
	t.Parallel()

	testResponseCookieKeyValue(t, "foo=bar", "foo", "bar")
	testResponseCookieKeyValue(t, "key=value", "key", "value")
	testResponseCookieKeyValue(t, "key=1", "key", "1")
	testResponseCookieKeyValue(t, "key=", "key", "")
	testResponseCookieKeyValue(t, "=value", "", "value")
	testResponseCookieKeyValue(t, "key", "key", "")
	testResponseCookieKeyValue(t, "key=", "key", "")
	testResponseCookieKeyValue(t, "key=; key2=abc", "key2", "abc")
}

func testResponseCookieKeyValue(t *testing.T, s, expectedKey, expectedValue string) {
	k, v := ParseResponseCookie([]byte(s))
	if string(k) != expectedKey {
		t.Fatalf("unexpected cookie key %q. Expecting %q. cookie %q", k, expectedKey, s)
	}
	if string(v) != expectedValue {
		t.Fatalf("unexpected cookie value %q. Expecting %q. cookie %q", v, expectedValue, s)
	}
}

func TestParseRequestCookies(t *testing.T) {
	t.Parallel()

	testParseRequestCookies(t, "foo=bar", map[string]string{"foo": "bar"})
	testParseRequestCookies(t, "foo=bar; aaa=bbb", map[string]string{"foo": "bar", "aaa": "bbb"})
	testParseRequestCookies(t, "  foo=bar;  aaa=bbb  ", map[string]string{"foo": "bar", "aaa": "bbb"})
	testParseRequestCookies(t, "  foo=bar;  aaa=", map[string]string{"foo": "bar", "aaa": ""})
	testParseRequestCookies(t, "foobar", map[string]string{"foobar": ""})
	testParseRequestCookies(t, "foo=bar=aaa", map[string]string{"foo": "bar=aaa"})
}

func testParseRequestCookies(t *testing.T, s string, expectedKV map[string]string) {
	c := &RequestHeader{}
	c.Set("Cookie", s)

	c.VisitAllCookie(func(key, value []byte) {
		expectedValue, ok := expectedKV[string(key)]
		if !ok {
			t.Fatalf("unexpected key %q. expecting one of %v", key, expectedKV)
		}
		if string(value) != expectedValue {
			t.Fatalf("unexpected value %q for key %q. expecting %q. cookies %q", value, key, expectedValue, s)
		}
		delete(expectedKV, string(key))
	})

	if len(expectedKV) > 0 {
		t.Fatalf("missing cookies %v. cookies %q", expectedKV, s)
	}
}

func TestAppendResponseCookieBytes(t *testing.T) {
	t.Parallel()

	testAppendResponseCookieBytes(t, "", "")
	testAppendResponseCookieBytes(t, "foo=bar; aaa=b", "Set-Cookie: foo=bar\r\nSet-Cookie: aaa=b\r\n")
	testAppendResponseCookieBytes(t, "foo=bar, aaa=b", "Set-Cookie: foo=bar\r\nSet-Cookie: aaa=b\r\n")
	testAppendResponseCookieBytes(t, "foo=bar; aaa=b; bbb=c", "Set-Cookie: foo=bar\r\nSet-Cookie: aaa=b\r\nSet-Cookie: bbb=c\r\n")
	testAppendResponseCookieBytes(t, "foo=bar, aaa=b, bbb=c", "Set-Cookie: foo=bar\r\nSet-Cookie: aaa=b\r\nSet-Cookie: bbb=c\r\n")
}

func testAppendResponseCookieBytes(t *testing.T, s, expectedResult string) {
	c := &ResponseHeader{}
	c.Set("Set-Cookie", s)
	c.VisitAllCookie(func(key, value []byte) {
		t.Fatalf("unexpected cookie %q=%q", key, value)
	})

	result := c.String()
	if result != expectedResult {
		t.Fatalf("unexpected result %q. Expecting %q", result, expectedResult)
	}
}

func TestAppendRequestCookieBytes(t *testing.T) {
	t.Parallel()

	testAppendRequestCookieBytes(t, "", "")
	testAppendRequestCookieBytes(t, "foo=bar; aaa=b", "Cookie: foo=bar; aaa=b\r\n")
	testAppendRequestCookieBytes(t, "foo=bar, aaa=b", "Cookie: foo=bar; aaa=b\r\n")
	testAppendRequestCookieBytes(t, "foo=bar; aaa=b; bbb=c", "Cookie: foo=bar; aaa=b; bbb=c\r\n")
	testAppendRequestCookieBytes(t, "foo=bar, aaa=b, bbb=c", "Cookie: foo=bar; aaa=b; bbb=c\r\n")
}

func testAppendRequestCookieBytes(t *testing.T, s, expectedResult string) {
	c := &RequestHeader{}
	c.Set("Cookie", s)
	c.VisitAllCookie(func(key, value []byte) {
		t.Fatalf("unexpected cookie %q=%q", key, value)
	})

	result := c.String()
	if result != expectedResult {
		t.Fatalf("unexpected result %q. Expecting %q", result, expectedResult)
	}
}

func TestAppendBytesCookie(t *testing.T) {
	t.Parallel()

	testAppendBytesCookie(t, "foo=bar; Path=/path; Domain=aaa.com; Max-Age=150; Secure; HttpOnly",
		"foo", "bar", "/path", "aaa.com", 150, true, true, CookieExpireUnlimited, CookieSameSiteDisabled, false)
	testAppendBytesCookie(t, "foo=bar; Path=/path; Domain=aaa.com; Max-Age=150; Secure; HttpOnly; SameSite=Lax",
		"foo", "bar", "/path", "aaa.com", 150, true, true, CookieExpireUnlimited, CookieSameSiteLaxMode, false)
	testAppendBytesCookie(t, "foo=bar; Path=/path; Domain=aaa.com; Max-Age=150; Secure; HttpOnly; SameSite=Strict",
		"foo", "bar", "/path", "aaa.com", 150, true, true, CookieExpireUnlimited, CookieSameSiteStrictMode, false)
	testAppendBytesCookie(t, "foo=bar; Path=/; Domain=aaa.com; Max-Age=150; Secure; HttpOnly; SameSite=None",
		"foo", "bar", "/", "aaa.com", 150, true, true, CookieExpireUnlimited, CookieSameSiteNoneMode, false)
	testAppendBytesCookie(t, "foo=bar; Path=/foo/bar; Domain=aaa.com; Max-Age=150; Secure; Partitioned",
		"foo", "bar", "/foo/bar", "aaa.com", 150, true, false, CookieExpireUnlimited, CookieSameSiteDisabled, true)
	testAppendBytesCookie(t, "foo=bar; Path=/foo/bar; Domain=aaa.com; Max-Age=150; Secure; HttpOnly; SameSite; Partitioned",
		"foo", "bar", "/foo/bar", "aaa.com", 150, true, true, CookieExpireUnlimited, CookieSameSiteDefaultMode, true)
	testAppendBytesCookie(t, "foo=bar; Path=/foo/bar; Domain=aaa.com; Max-Age=150; Secure; HttpOnly; SameSite=Lax; Partitioned",
		"foo", "bar", "/foo/bar", "aaa.com", 150, true, true, CookieExpireUnlimited, CookieSameSiteLaxMode, true)
	testAppendBytesCookie(t, "foo=bar; Path=/foo/bar; Domain=aaa.com; Max-Age=150; Secure; HttpOnly; SameSite=Strict; Partitioned",
		"foo", "bar", "/foo/bar", "aaa.com", 150, true, true, CookieExpireUnlimited, CookieSameSiteStrictMode, true)
	testAppendBytesCookie(t, "foo=bar; Path=/foo/bar; Domain=aaa.com; Max-Age=150; Secure; HttpOnly; SameSite=None; Partitioned",
		"foo", "bar", "/foo/bar", "aaa.com", 150, true, true, CookieExpireUnlimited, CookieSameSiteNoneMode, true)
}

func testAppendBytesCookie(t *testing.T, expectedResult, key, value, path, domain string, maxAge int, secure, httpOnly bool, expire time.Time, sameSite CookieSameSite, partitioned bool) {
	c := &Cookie{}
	c.SetKey(key)
	c.SetValue(value)
	c.SetPath(path)
	c.SetDomain(domain)
	c.SetMaxAge(maxAge)
	c.SetSecure(secure)
	c.SetHTTPOnly(httpOnly)
	c.SetExpire(expire)
	c.SetSameSite(sameSite)
	c.SetPartitioned(partitioned)

	result := string(c.Cookie())
	if result != expectedResult {
		t.Fatalf("unexpected result %q. Expecting %q", result, expectedResult)
	}
}

func TestParseSetCookie(t *testing.T) {
	t.Parallel()

	testParseSetCookie(t, "foo=bar; Path=/path; Domain=aaa.com; Max-Age=150; Secure; HttpOnly",
		"foo", "bar", "/path", "aaa.com", 150, true, true, CookieExpireUnlimited, CookieSameSiteDisabled, false)
	testParseSetCookie(t, "foo=bar; Path=/path; Domain=aaa.com; Max-Age=150; Secure; HttpOnly; SameSite=Lax",
		"foo", "bar", "/path", "aaa.com", 150, true, true, CookieExpireUnlimited, CookieSameSiteLaxMode, false)
	testParseSetCookie(t, "foo=bar; Path=/path; Domain=aaa.com; Max-Age=150; Secure; HttpOnly; SameSite=Strict",
		"foo", "bar", "/path", "aaa.com", 150, true, true, CookieExpireUnlimited, CookieSameSiteStrictMode, false)
	testParseSetCookie(t, "foo=bar; Path=/; Domain=aaa.com; Max-Age=150; Secure; HttpOnly; SameSite=None",
		"foo", "bar", "/", "aaa.com", 150, true, true, CookieExpireUnlimited, CookieSameSiteNoneMode, false)
	testParseSetCookie(t, "foo=bar; Path=/foo/bar; Domain=aaa.com; Max-Age=150; Secure; Partitioned",
		"foo", "bar", "/foo/bar", "aaa.com", 150, true, false, CookieExpireUnlimited, CookieSameSiteDisabled, true)
	testParseSetCookie(t, "foo=bar; Path=/foo/bar; Domain=aaa.com; Max-Age=150; Secure; HttpOnly; SameSite; Partitioned",
		"foo", "bar", "/foo/bar", "aaa.com", 150, true, true, CookieExpireUnlimited, CookieSameSiteDefaultMode, true)
	testParseSetCookie(t, "foo=bar; Path=/foo/bar; Domain=aaa.com; Max-Age=150; Secure; HttpOnly; SameSite=Lax; Partitioned",
		"foo", "bar", "/foo/bar", "aaa.com", 150, true, true, CookieExpireUnlimited, CookieSameSiteLaxMode, true)
	testParseSetCookie(t, "foo=bar; Path=/foo/bar; Domain=aaa.com; Max-Age=150; Secure; HttpOnly; SameSite=Strict; Partitioned",
		"foo", "bar", "/foo/bar", "aaa.com", 150, true, true, CookieExpireUnlimited, CookieSameSiteStrictMode, true)
	testParseSetCookie(t, "foo=bar; Path=/foo/bar; Domain=aaa.com; Max-Age=150; Secure; HttpOnly; SameSite=None; Partitioned",
		"foo", "bar", "/foo/bar", "aaa.com", 150, true, true, CookieExpireUnlimited, CookieSameSiteNoneMode, true)
}

func testParseSetCookie(t *testing.T, s, expectedKey, expectedValue, expectedPath, expectedDomain string, expectedMaxAge int, expectedSecure, expectedHTTPOnly bool, expectedExpire time.Time, expectedSameSite CookieSameSite, expectedPartitioned bool) {
	c := &Cookie{}
	if err := c.Parse(s); err != nil {
		t.Fatalf("unexpected error: %s. cookie %q", err, s)
	}
	if string(c.Key()) != expectedKey {
		t.Fatalf("unexpected cookie key %q. Expecting %q. cookie %q", c.Key(), expectedKey, s)
	}
	if string(c.Value()) != expectedValue {
		t.Fatalf("unexpected cookie value %q. Expecting %q. cookie %q", c.Value(), expectedValue, s)
	}
	if string(c.Path()) != expectedPath {
		t.Fatalf("unexpected cookie path %q. Expecting %q. cookie %q", c.Path(), expectedPath, s)
	}
	if string(c.Domain()) != expectedDomain {
		t.Fatalf("unexpected cookie domain %q. Expecting %q. cookie %q", c.Domain(), expectedDomain, s)
	}
	if c.MaxAge() != expectedMaxAge {
		t.Fatalf("unexpected cookie max-age %d. Expecting %d. cookie %q", c.MaxAge(), expectedMaxAge, s)
	}
	if c.Secure() != expectedSecure {
		t.Fatalf("unexpected cookie secure flag %v. Expecting %v. cookie %q", c.Secure(), expectedSecure, s)
	}
	if c.HTTPOnly() != expectedHTTPOnly {
		t.Fatalf("unexpected cookie httpOnly flag %v. Expecting %v. cookie %q", c.HTTPOnly(), expectedHTTPOnly, s)
	}
	if c.Expire() != expectedExpire {
		t.Fatalf("unexpected cookie expire %v. Expecting %v. cookie %q", c.Expire(), expectedExpire, s)
	}
	if c.SameSite() != expectedSameSite {
		t.Fatalf("unexpected cookie SameSite %v. Expecting %v. cookie %q", c.SameSite(), expectedSameSite, s)
	}
	if c.Partitioned() != expectedPartitioned {
		t.Fatalf("unexpected cookie Partitioned %v. Expecting %v. cookie %q", c.Partitioned(), expectedPartitioned, s)
	}

	c.Reset()
	c.SetKey(expectedKey)
	c.SetValue(expectedValue)
	c.SetPath(expectedPath)
	c.SetDomain(expectedDomain)
	c.SetMaxAge(expectedMaxAge)
	c.SetSecure(expectedSecure)
	c.SetHTTPOnly(expectedHTTPOnly)
	c.SetExpire(expectedExpire)
	c.SetSameSite(expectedSameSite)
	c.SetPartitioned(expectedPartitioned)

	result := string(c.Cookie())
	if result != s {
		t.Fatalf("unexpected result %q. Expecting %q", result, s)
	}
}

func TestParseCookieError(t *testing.T) {
	t.Parallel()

	var c Cookie
	if err := c.Parse(";foo=bar"); err != errNoCookies {
		t.Fatalf("unexpected error: %s. Expecting %s", err, errNoCookies)
	}
	if err := c.Parse("="); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if err := c.Parse("=bar"); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if err := c.Parse("foo="); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if err := c.Parse("="); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if err := c.Parse("foo=bar; max-age=foobar"); err == nil {
		t.Fatalf("expecting non-nil error")
	}
	if err := c.Parse("foo=bar; expires=foobar"); err == nil {
		t.Fatalf("expecting non-nil error")
	}
}

func TestCookieCopyTo(t *testing.T) {
	t.Parallel()

	var c Cookie
	c.SetKey("foo")
	c.SetValue("bar")
	c.SetPath("/path")
	c.SetDomain("aaa.com")
	c.SetMaxAge(150)
	c.SetSecure(true)
	c.SetHTTPOnly(true)
	c.SetExpire(CookieExpireDelete)
	c.SetSameSite(CookieSameSiteLaxMode)
	c.SetPartitioned(true)

	var c1 Cookie
	c.CopyTo(&c1)

	if string(c1.Key()) != string(c.Key()) {
		t.Fatalf("unexpected cookie key %q. Expecting %q", c1.Key(), c.Key())
	}
	if string(c1.Value()) != string(c.Value()) {
		t.Fatalf("unexpected cookie value %q. Expecting %q", c1.Value(), c.Value())
	}
	if string(c1.Path()) != string(c.Path()) {
		t.Fatalf("unexpected cookie path %q. Expecting %q", c1.Path(), c.Path())
	}
	if string(c1.Domain()) != string(c.Domain()) {
		t.Fatalf("unexpected cookie domain %q. Expecting %q", c1.Domain(), c.Domain())
	}
	if c1.MaxAge() != c.MaxAge() {
		t.Fatalf("unexpected cookie max-age %d. Expecting %d", c1.MaxAge(), c.MaxAge())
	}
	if c1.Secure() != c.Secure() {
		t.Fatalf("unexpected cookie secure flag %v. Expecting %v", c1.Secure(), c.Secure())
	}
	if c1.HTTPOnly() != c.HTTPOnly() {
		t.Fatalf("unexpected cookie httpOnly flag %v. Expecting %v", c1.HTTPOnly(), c.HTTPOnly())
	}
	if c1.Expire() != c.Expire() {
		t.Fatalf("unexpected cookie expire %v. Expecting %v", c1.Expire(), c.Expire())
	}
	if c1.SameSite() != c.SameSite() {
		t.Fatalf("unexpected cookie SameSite %v. Expecting %v", c1.SameSite(), c.SameSite())
	}
	if c1.Partitioned() != c.Partitioned() {
		t.Fatalf("unexpected cookie Partitioned %v. Expecting %v", c1.Partitioned(), c.Partitioned())
	}
}

func TestCookieCopyToOverwrite(t *testing.T) {
	t.Parallel()

	var c Cookie
	c.SetKey("foo")
	c.SetValue("bar")
	c.SetPath("/path")
	c.SetDomain("aaa.com")
	c.SetMaxAge(150)
	c.SetSecure(true)
	c.SetHTTPOnly(true)
	c.SetExpire(CookieExpireDelete)
	c.SetSameSite(CookieSameSiteLaxMode)
	c.SetPartitioned(true)

	var c1 Cookie
	c1.SetKey("bar")
	c1.SetValue("foo")
	c1.SetPath("/path/")
	c1.SetDomain("bbb.com")
	c1.SetMaxAge(100)
	c1.SetSecure(false)
	c1.SetHTTPOnly(false)
	c1.SetExpire(CookieExpireUnlimited)
	c1.SetSameSite(CookieSameSiteStrictMode)
	c1.SetPartitioned(false)

	c.CopyTo(&c1)

	if string(c1.Key()) != string(c.Key()) {
		t.Fatalf("unexpected cookie key %q. Expecting %q", c1.Key(), c.Key())
	}
	if string(c1.Value()) != string(c.Value()) {
		t.Fatalf("unexpected cookie value %q. Expecting %q", c1.Value(), c.Value())
	}
	if string(c1.Path()) != string(c.Path()) {
		t.Fatalf("unexpected cookie path %q. Expecting %q", c1.Path(), c.Path())
	}
	if string(c1.Domain()) != string(c.Domain()) {
		t.Fatalf("unexpected cookie domain %q. Expecting %q", c1.Domain(), c.Domain())
	}
	if c1.MaxAge() != c.MaxAge() {
		t.Fatalf("unexpected cookie max-age %d. Expecting %d", c1.MaxAge(), c.MaxAge())
	}
	if c1.Secure() != c.Secure() {
		t.Fatalf("unexpected cookie secure flag %v. Expecting %v", c1.Secure(), c.Secure())
	}
	if c1.HTTPOnly() != c.HTTPOnly() {
		t.Fatalf("unexpected cookie httpOnly flag %v. Expecting %v", c1.HTTPOnly(), c.HTTPOnly())
	}
	if c1.Expire() != c.Expire() {
		t.Fatalf("unexpected cookie expire %v. Expecting %v", c1.Expire(), c.Expire())
	}
	if c1.SameSite() != c.SameSite() {
		t.Fatalf("unexpected cookie SameSite %v. Expecting %v", c1.SameSite(), c.SameSite())
	}
	if c1.Partitioned() != c.Partitioned() {
		t.Fatalf("unexpected cookie Partitioned %v. Expecting %v", c1.Partitioned(), c.Partitioned())
	}
}

func TestCookieAppendBytes(t *testing.T) {
	t.Parallel()

	var c Cookie
	c.SetKey("foo")
	c.SetValue("bar")
	c.SetPath("/path")
	c.SetDomain("aaa.com")
	c.SetMaxAge(150)
	c.SetSecure(true)
	c.SetHTTPOnly(true)
	c.SetExpire(CookieExpireDelete)
	c.SetSameSite(CookieSameSiteLaxMode)
	c.SetPartitioned(true)

	expectedResult := "foo=bar; Max-Age=150; Domain=aaa.com; Path=/path; HttpOnly; Secure; SameSite=Lax; Partitioned"

	result := string(c.AppendBytes(nil))
	if result != expectedResult {
		t.Fatalf("unexpected result %q. Expecting %q", result, expectedResult)
	}
}

func TestCookieSetKeyBytes(t *testing.T) {
	t.Parallel()

	var c Cookie
	c.SetKeyBytes([]byte("foo"))
	c.SetValue("bar")

	expectedResult := "foo=bar"

	result := string(c.AppendBytes(nil))
	if result != expectedResult {
		t.Fatalf("unexpected result %q. Expecting %q", result, expectedResult)
	}
}

func TestCookieSetValueBytes(t *testing.T) {
	t.Parallel()

	var c Cookie
	c.SetKey("foo")
	c.SetValueBytes([]byte("bar"))

	expectedResult := "foo=bar"

	result := string(c.AppendBytes(nil))
	if result != expectedResult {
		t.Fatalf("unexpected result %q. Expecting %q", result, expectedResult)
	}
}

func TestCookieSetPathBytes(t *testing.T) {
	t.Parallel()

	var c Cookie
	c.SetKey("foo")
	c.SetValue("bar")
	c.SetPathBytes([]byte("/path"))

	expectedResult := "foo=bar; Path=/path"

	result := string(c.AppendBytes(nil))
	if result != expectedResult {
		t.Fatalf("unexpected result %q. Expecting %q", result, expectedResult)
	}
}

func TestCookieSetDomainBytes(t *testing.T) {
	t.Parallel()

	var c Cookie
	c.SetKey("foo")
	c.SetValue("bar")
	c.SetDomainBytes([]byte("aaa.com"))

	expectedResult := "foo=bar; Domain=aaa.com"

	result := string(c.AppendBytes(nil))
	if result != expectedResult {
		t.Fatalf("unexpected result %q. Expecting %q", result, expectedResult)
	}
}

func testCookie(t *testing.T, c *Cookie, expectedKey, expectedValue string, expectedExpire time.Time) {
	if string(c.Key()) != expectedKey {
		t.Fatalf("unexpected key %q. Expecting %q", c.Key(), expectedKey)
	}
	if string(c.Value()) != expectedValue {
		t.Fatalf("unexpected value %q. Expecting %q", c.Value(), expectedValue)
	}
	if c.Expire() != expectedExpire {
		t.Fatalf("unexpected expire %v. Expecting %v", c.Expire(), expectedExpire)
	}
}