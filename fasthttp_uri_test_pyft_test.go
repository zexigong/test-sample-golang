package fasthttp

import (
	"bytes"
	"fmt"
	"testing"
)

func TestURIPathNormalize(t *testing.T) {
	testURIPathNormalize(t, "", "/")
	testURIPathNormalize(t, "a/b", "/a/b")
	testURIPathNormalize(t, "/a/b", "/a/b")
	testURIPathNormalize(t, "/a/b/", "/a/b/")
	testURIPathNormalize(t, "/a/b//", "/a/b/")
	testURIPathNormalize(t, "/a/b//c", "/a/b/c")
	testURIPathNormalize(t, "/a/b/./c", "/a/b/c")
	testURIPathNormalize(t, "/a/b/../c", "/a/c")
	testURIPathNormalize(t, "/a/b/../../c", "/c")
	testURIPathNormalize(t, "/a/b/../../c/", "/c/")
	testURIPathNormalize(t, "/a/b/../../c/d", "/c/d")
	testURIPathNormalize(t, "/a/b/../../c/d/", "/c/d/")
	testURIPathNormalize(t, "/a/b/../../../c", "/c")
	testURIPathNormalize(t, "/a/b/../../../", "/")
	testURIPathNormalize(t, "/a/b/.././../c", "/c")
	testURIPathNormalize(t, "/a/b/.././../c/", "/c/")
	testURIPathNormalize(t, "/a/b/.././../c/d", "/c/d")
	testURIPathNormalize(t, "/a/b/.././../c/d/", "/c/d/")
	testURIPathNormalize(t, "/a/b/.././../c/d/.", "/c/d/")
	testURIPathNormalize(t, "/a/b/.././../c/d/..", "/c/")
	testURIPathNormalize(t, "/a/b/.././../c/d/./", "/c/d/")
	testURIPathNormalize(t, "/a/b/.././../c/d/./..", "/c/")
	testURIPathNormalize(t, "/a/b/.././../c/d/../", "/c/")
	testURIPathNormalize(t, "/a/b/.././../c/d/../.", "/c/")
	testURIPathNormalize(t, "/a/b/.././../c/d/../../..", "/")
	testURIPathNormalize(t, "/a/b/.././../c/d/../../../", "/")
	testURIPathNormalize(t, "/a/b/.././../c/d/../../..e", "/..e")

	testURIPathNormalize(t, "/a//b", "/a/b")
	testURIPathNormalize(t, "/a//b/c", "/a/b/c")
	testURIPathNormalize(t, "/a//b//c", "/a/b/c")
	testURIPathNormalize(t, "/a//b/./c", "/a/b/c")
	testURIPathNormalize(t, "/a//b/../c", "/a/c")
	testURIPathNormalize(t, "/a//b/../../c", "/c")
	testURIPathNormalize(t, "/a//b/../../c/", "/c/")
	testURIPathNormalize(t, "/a//b/../../c/d", "/c/d")
	testURIPathNormalize(t, "/a//b/../../c/d/", "/c/d/")
	testURIPathNormalize(t, "/a//b/../../../c", "/c")
	testURIPathNormalize(t, "/a//b/../../../", "/")
	testURIPathNormalize(t, "/a//b/.././../c", "/c")
	testURIPathNormalize(t, "/a//b/.././../c/", "/c/")
	testURIPathNormalize(t, "/a//b/.././../c/d", "/c/d")
	testURIPathNormalize(t, "/a//b/.././../c/d/", "/c/d/")
	testURIPathNormalize(t, "/a//b/.././../c/d/.", "/c/d/")
	testURIPathNormalize(t, "/a//b/.././../c/d/..", "/c/")
	testURIPathNormalize(t, "/a//b/.././../c/d/./", "/c/d/")
	testURIPathNormalize(t, "/a//b/.././../c/d/./..", "/c/")
	testURIPathNormalize(t, "/a//b/.././../c/d/../", "/c/")
	testURIPathNormalize(t, "/a//b/.././../c/d/../.", "/c/")
	testURIPathNormalize(t, "/a//b/.././../c/d/../../..", "/")
	testURIPathNormalize(t, "/a//b/.././../c/d/../../../", "/")
	testURIPathNormalize(t, "/a//b/.././../c/d/../../..e", "/..e")

	testURIPathNormalize(t, "/a//b", "/a/b")
	testURIPathNormalize(t, "/a//b/c", "/a/b/c")
	testURIPathNormalize(t, "/a//b//c", "/a/b/c")
	testURIPathNormalize(t, "/a//b/./c", "/a/b/c")
	testURIPathNormalize(t, "/a//b/../c", "/a/c")
	testURIPathNormalize(t, "/a//b/../../c", "/c")
	testURIPathNormalize(t, "/a//b/../../c/", "/c/")
	testURIPathNormalize(t, "/a//b/../../c/d", "/c/d")
	testURIPathNormalize(t, "/a//b/../../c/d/", "/c/d/")
	testURIPathNormalize(t, "/a//b/../../../c", "/c")
	testURIPathNormalize(t, "/a//b/../../../", "/")
	testURIPathNormalize(t, "/a//b/.././../c", "/c")
	testURIPathNormalize(t, "/a//b/.././../c/", "/c/")
	testURIPathNormalize(t, "/a//b/.././../c/d", "/c/d")
	testURIPathNormalize(t, "/a//b/.././../c/d/", "/c/d/")
	testURIPathNormalize(t, "/a//b/.././../c/d/.", "/c/d/")
	testURIPathNormalize(t, "/a//b/.././../c/d/..", "/c/")
	testURIPathNormalize(t, "/a//b/.././../c/d/./", "/c/d/")
	testURIPathNormalize(t, "/a//b/.././../c/d/./..", "/c/")
	testURIPathNormalize(t, "/a//b/.././../c/d/../", "/c/")
	testURIPathNormalize(t, "/a//b/.././../c/d/../.", "/c/")
	testURIPathNormalize(t, "/a//b/.././../c/d/../../..", "/")
	testURIPathNormalize(t, "/a//b/.././../c/d/../../../", "/")
	testURIPathNormalize(t, "/a//b/.././../c/d/../../..e", "/..e")

	testURIPathNormalize(t, "/a/./b", "/a/b")
	testURIPathNormalize(t, "/a/./b/c", "/a/b/c")
	testURIPathNormalize(t, "/a/./b//c", "/a/b/c")
	testURIPathNormalize(t, "/a/./b/./c", "/a/b/c")
	testURIPathNormalize(t, "/a/./b/../c", "/a/c")
	testURIPathNormalize(t, "/a/./b/../../c", "/c")
	testURIPathNormalize(t, "/a/./b/../../c/", "/c/")
	testURIPathNormalize(t, "/a/./b/../../c/d", "/c/d")
	testURIPathNormalize(t, "/a/./b/../../c/d/", "/c/d/")
	testURIPathNormalize(t, "/a/./b/../../../c", "/c")
	testURIPathNormalize(t, "/a/./b/../../../", "/")
	testURIPathNormalize(t, "/a/./b/.././../c", "/c")
	testURIPathNormalize(t, "/a/./b/.././../c/", "/c/")
	testURIPathNormalize(t, "/a/./b/.././../c/d", "/c/d")
	testURIPathNormalize(t, "/a/./b/.././../c/d/", "/c/d/")
	testURIPathNormalize(t, "/a/./b/.././../c/d/.", "/c/d/")
	testURIPathNormalize(t, "/a/./b/.././../c/d/..", "/c/")
	testURIPathNormalize(t, "/a/./b/.././../c/d/./", "/c/d/")
	testURIPathNormalize(t, "/a/./b/.././../c/d/./..", "/c/")
	testURIPathNormalize(t, "/a/./b/.././../c/d/../", "/c/")
	testURIPathNormalize(t, "/a/./b/.././../c/d/../.", "/c/")
	testURIPathNormalize(t, "/a/./b/.././../c/d/../../..", "/")
	testURIPathNormalize(t, "/a/./b/.././../c/d/../../../", "/")
	testURIPathNormalize(t, "/a/./b/.././../c/d/../../..e", "/..e")

	testURIPathNormalize(t, "/a/../b", "/b")
	testURIPathNormalize(t, "/a/../b/c", "/b/c")
	testURIPathNormalize(t, "/a/../b//c", "/b/c")
	testURIPathNormalize(t, "/a/../b/./c", "/b/c")
	testURIPathNormalize(t, "/a/../b/../c", "/c")
	testURIPathNormalize(t, "/a/../b/../../c", "/c")
	testURIPathNormalize(t, "/a/../b/../../c/", "/c/")
	testURIPathNormalize(t, "/a/../b/../../c/d", "/c/d")
	testURIPathNormalize(t, "/a/../b/../../c/d/", "/c/d/")
	testURIPathNormalize(t, "/a/../b/../../../c", "/c")
	testURIPathNormalize(t, "/a/../b/../../../", "/")
	testURIPathNormalize(t, "/a/../b/.././../c", "/c")
	testURIPathNormalize(t, "/a/../b/.././../c/", "/c/")
	testURIPathNormalize(t, "/a/../b/.././../c/d", "/c/d")
	testURIPathNormalize(t, "/a/../b/.././../c/d/", "/c/d/")
	testURIPathNormalize(t, "/a/../b/.././../c/d/.", "/c/d/")
	testURIPathNormalize(t, "/a/../b/.././../c/d/..", "/c/")
	testURIPathNormalize(t, "/a/../b/.././../c/d/./", "/c/d/")
	testURIPathNormalize(t, "/a/../b/.././../c/d/./..", "/c/")
	testURIPathNormalize(t, "/a/../b/.././../c/d/../", "/c/")
	testURIPathNormalize(t, "/a/../b/.././../c/d/../.", "/c/")
	testURIPathNormalize(t, "/a/../b/.././../c/d/../../..", "/")
	testURIPathNormalize(t, "/a/../b/.././../c/d/../../../", "/")
	testURIPathNormalize(t, "/a/../b/.././../c/d/../../..e", "/..e")

	testURIPathNormalize(t, "/a/../../b", "/b")
	testURIPathNormalize(t, "/a/../../b/c", "/b/c")
	testURIPathNormalize(t, "/a/../../b//c", "/b/c")
	testURIPathNormalize(t, "/a/../../b/./c", "/b/c")
	testURIPathNormalize(t, "/a/../../b/../c", "/c")
	testURIPathNormalize(t, "/a/../../b/../../c", "/c")
	testURIPathNormalize(t, "/a/../../b/../../c/", "/c/")
	testURIPathNormalize(t, "/a/../../b/../../c/d", "/c/d")
	testURIPathNormalize(t, "/a/../../b/../../c/d/", "/c/d/")
	testURIPathNormalize(t, "/a/../../b/../../../c", "/c")
	testURIPathNormalize(t, "/a/../../b/../../../", "/")
	testURIPathNormalize(t, "/a/../../b/.././../c", "/c")
	testURIPathNormalize(t, "/a/../../b/.././../c/", "/c/")
	testURIPathNormalize(t, "/a/../../b/.././../c/d", "/c/d")
	testURIPathNormalize(t, "/a/../../b/.././../c/d/", "/c/d/")
	testURIPathNormalize(t, "/a/../../b/.././../c/d/.", "/c/d/")
	testURIPathNormalize(t, "/a/../../b/.././../c/d/..", "/c/")
	testURIPathNormalize(t, "/a/../../b/.././../c/d/./", "/c/d/")
	testURIPathNormalize(t, "/a/../../b/.././../c/d/./..", "/c/")
	testURIPathNormalize(t, "/a/../../b/.././../c/d/../", "/c/")
	testURIPathNormalize(t, "/a/../../b/.././../c/d/../.", "/c/")
	testURIPathNormalize(t, "/a/../../b/.././../c/d/../../..", "/")
	testURIPathNormalize(t, "/a/../../b/.././../c/d/../../../", "/")
	testURIPathNormalize(t, "/a/../../b/.././../c/d/../../..e", "/..e")

	testURIPathNormalize(t, "/a/b/c/d/././..", "/a/b/c/")
	testURIPathNormalize(t, "/a/b/c/d/././../", "/a/b/c/")
	testURIPathNormalize(t, "/a/b/c/d/././..e", "/a/b/c/d/..e")
	testURIPathNormalize(t, "/a/b/c/d/././..e/", "/a/b/c/d/..e/")
	testURIPathNormalize(t, "/a/b/c/d/././..e/f", "/a/b/c/d/..e/f")
	testURIPathNormalize(t, "/a/b/c/d/././..e/f/", "/a/b/c/d/..e/f/")

	// last path parts
	testURIPathNormalize(t, "/a/..", "/")
	testURIPathNormalize(t, "/a/../", "/")
	testURIPathNormalize(t, "/a/..b", "/a/..b")
	testURIPathNormalize(t, "/a/..b/", "/a/..b/")
	testURIPathNormalize(t, "/a/..b/c", "/a/..b/c")
	testURIPathNormalize(t, "/a/..b/c/", "/a/..b/c/")
	testURIPathNormalize(t, "/a/..b/c/.", "/a/..b/c/")
	testURIPathNormalize(t, "/a/..b/c/..", "/a/")

	// '.' in the last path part
	testURIPathNormalize(t, "/a/.", "/a/")
	testURIPathNormalize(t, "/a/./", "/a/")
	testURIPathNormalize(t, "/a/././", "/a/")
	testURIPathNormalize(t, "/a/././.", "/a/")
	testURIPathNormalize(t, "/a/./././", "/a/")
	testURIPathNormalize(t, "/a/././.b", "/a/.b")
	testURIPathNormalize(t, "/a/././.b/", "/a/.b/")

	// '..' in the last path part
	testURIPathNormalize(t, "/a/..", "/")
	testURIPathNormalize(t, "/a/../", "/")
	testURIPathNormalize(t, "/a/../..", "/")
	testURIPathNormalize(t, "/a/../../", "/")
	testURIPathNormalize(t, "/a/../../.", "/")
	testURIPathNormalize(t, "/a/../.././", "/")
	testURIPathNormalize(t, "/a/../../.b", "/.b")
	testURIPathNormalize(t, "/a/../../.b/", "/.b/")
	testURIPathNormalize(t, "/a/../../..b", "/..b")
	testURIPathNormalize(t, "/a/../../..b/", "/..b/")
	testURIPathNormalize(t, "/a/../../..b/c", "/..b/c")
	testURIPathNormalize(t, "/a/../../..b/c/", "/..b/c/")
	testURIPathNormalize(t, "/a/../../..b/c/.", "/..b/c/")
	testURIPathNormalize(t, "/a/../../..b/c/..", "/")

	testURIPathNormalize(t, "/a/b/c/d/././..", "/a/b/c/")
	testURIPathNormalize(t, "/a/b/c/d/././../", "/a/b/c/")
	testURIPathNormalize(t, "/a/b/c/d/././..e", "/a/b/c/d/..e")
	testURIPathNormalize(t, "/a/b/c/d/././..e/", "/a/b/c/d/..e/")
	testURIPathNormalize(t, "/a/b/c/d/././..e/f", "/a/b/c/d/..e/f")
	testURIPathNormalize(t, "/a/b/c/d/././..e/f/", "/a/b/c/d/..e/f/")

	// test path normalization
	testURIPathNormalize(t, "/a%2fb", "/a%2Fb")
	testURIPathNormalize(t, "/a/b%c0%afc", "/a/b%C0%Af/") // invalid utf8
}

func testURIPathNormalize(t *testing.T, path, expectedNormalizedPath string) {
	t.Helper()

	var u URI

	u.SetPath(path)
	if string(u.Path()) != expectedNormalizedPath {
		t.Fatalf("unexpected path %q. Expecting %q. path=%q", u.Path(), expectedNormalizedPath, path)
	}

	u.DisablePathNormalizing = true
	u.SetPath(path)
	if string(u.Path()) != path {
		t.Fatalf("unexpected path %q. Expecting %q. path=%q", u.Path(), path, path)
	}
}

func TestURIFullURI(t *testing.T) {
	testURIFullURI(t, "", "http://foobar.com/")
	testURIFullURI(t, "/", "http://foobar.com/")
	testURIFullURI(t, "/aaa/bbb?ccc=ddd&qqq#zzz", "http://foobar.com/aaa/bbb?ccc=ddd&qqq#zzz")
	testURIFullURI(t, "aaa", "http://foobar.com/aaa")
	testURIFullURI(t, "aaa/aa?bb#qwerqwer", "http://foobar.com/aaa/aa?bb#qwerqwer")
	testURIFullURI(t, "http://aaa.com/aa/bbb?dd#aaa", "http://aaa.com/aa/bbb?dd#aaa")
	testURIFullURI(t, "https://aaa.com/aa/bbb?dd#aaa", "https://aaa.com/aa/bbb?dd#aaa")
	testURIFullURI(t, "https://foo.com", "https://foo.com/")
	testURIFullURI(t, "https://foo.com?bar=baz", "https://foo.com/?bar=baz")
	testURIFullURI(t, "https://foo.com#bar", "https://foo.com/#bar")
	testURIFullURI(t, "https://foo.com/#bar", "https://foo.com/#bar")
	testURIFullURI(t, "https://foo.com/?bar=baz", "https://foo.com/?bar=baz")
	testURIFullURI(t, "https://foo.com/bar", "https://foo.com/bar")
	testURIFullURI(t, "https://foo.com/bar#baz", "https://foo.com/bar#baz")
	testURIFullURI(t, "https://foo.com/bar?baz=quux", "https://foo.com/bar?baz=quux")
	testURIFullURI(t, "https://foo.com/bar?baz=quux#quuux", "https://foo.com/bar?baz=quux#quuux")
	testURIFullURI(t, "/%2F", "http://foobar.com/%2F")
}

func testURIFullURI(t *testing.T, uri, expectedFullURI string) {
	t.Helper()

	var u URI

	u.Parse([]byte("foobar.com"), []byte(uri))
	if string(u.FullURI()) != expectedFullURI {
		t.Fatalf("unexpected full uri %q. Expecting %q. uri=%q", u.FullURI(), expectedFullURI, uri)
	}

	u.DisablePathNormalizing = true
	u.Parse([]byte("foobar.com"), []byte(uri))
	if string(u.FullURI()) != expectedFullURI {
		t.Fatalf("unexpected full uri %q. Expecting %q. uri=%q", u.FullURI(), expectedFullURI, uri)
	}
}

func TestURIUpdate(t *testing.T) {
	testURIUpdate(t, "", "a", "http://foobar.com/a")
	testURIUpdate(t, "/", "a", "http://foobar.com/a")
	testURIUpdate(t, "/aa/bb?cc=dd&e", "a", "http://foobar.com/aa/a")
	testURIUpdate(t, "/aa/bb?cc=dd&e", "a/b", "http://foobar.com/aa/a/b")
	testURIUpdate(t, "/aa/bb?cc=dd&e", "/a/b", "http://foobar.com/a/b")
	testURIUpdate(t, "/aa/bb?cc=dd&e", "/a/b?zz", "http://foobar.com/a/b?zz")
	testURIUpdate(t, "/aa/bb?cc=dd&e", "/a/b?zz#ttt", "http://foobar.com/a/b?zz#ttt")
	testURIUpdate(t, "/aa/bb?cc=dd&e", "/a/b#ttt", "http://foobar.com/a/b#ttt")
	testURIUpdate(t, "/aa/bb?cc=dd&e", "http://aaa.com/xyz", "http://aaa.com/xyz")
	testURIUpdate(t, "/aa/bb?cc=dd&e", "https://aaa.com/xyz", "https://aaa.com/xyz")
	testURIUpdate(t, "/aa/bb?cc=dd&e", "ftp://aaa.com/xyz", "ftp://aaa.com/xyz")
	testURIUpdate(t, "/aa/bb?cc=dd&e", "//aaa.com/xyz", "http://aaa.com/xyz")
	testURIUpdate(t, "/aa/bb?cc=dd&e", "//aaa.com/xyz?op=90", "http://aaa.com/xyz?op=90")
	testURIUpdate(t, "/aa/bb?cc=dd&e", "//aaa.com/xyz#qwe", "http://aaa.com/xyz#qwe")
	testURIUpdate(t, "/aa/bb?cc=dd&e", "#qwe", "http://foobar.com/aa/bb?cc=dd&e#qwe")
	testURIUpdate(t, "/aa/bb?cc=dd&e", "?qwe", "http://foobar.com/aa/bb?qwe")
	testURIUpdate(t, "", "/%2F", "http://foobar.com/%2F")
}

func testURIUpdate(t *testing.T, uri, update, expectedFullURI string) {
	t.Helper()

	var u URI

	u.Parse([]byte("foobar.com"), []byte(uri))
	u.Update(update)
	if string(u.FullURI()) != expectedFullURI {
		t.Fatalf("unexpected full uri %q. Expecting %q. uri=%q, update=%q", u.FullURI(), expectedFullURI, uri, update)
	}

	u.DisablePathNormalizing = true
	u.Parse([]byte("foobar.com"), []byte(uri))
	u.Update(update)
	if string(u.FullURI()) != expectedFullURI {
		t.Fatalf("unexpected full uri %q. Expecting %q. uri=%q, update=%q", u.FullURI(), expectedFullURI, uri, update)
	}
}

func TestURICopyTo(t *testing.T) {
	var u URI
	u.Parse([]byte("foobar.com"), []byte("/aaa/bb?cc=dd&e"))

	var u1 URI
	u.CopyTo(&u1)

	if string(u1.FullURI()) != string(u.FullURI()) {
		t.Fatalf("unexpected full uri %q. Expecting %q", u1.FullURI(), u.FullURI())
	}

	u.DisablePathNormalizing = true
	u1.Reset()
	u.CopyTo(&u1)
	if string(u1.FullURI()) != string(u.FullURI()) {
		t.Fatalf("unexpected full uri %q. Expecting %q", u1.FullURI(), u.FullURI())
	}
}

func TestURICopyToQueryArgs(t *testing.T) {
	var u URI
	u.Parse([]byte("foobar.com"), []byte("/aaa/bb?cc=dd&e"))

	var u1 URI
	u.CopyTo(&u1)

	u.QueryArgs().Add("foo", "bar")
	if u1.QueryArgs().Peek("foo") != nil {
		t.Fatalf("unexpected foo value %q. Expecting %q", u1.QueryArgs().Peek("foo"), nil)
	}

	u.DisablePathNormalizing = true
	u1.Reset()
	u.CopyTo(&u1)
	u.QueryArgs().Add("foo", "bar")
	if u1.QueryArgs().Peek("foo") != nil {
		t.Fatalf("unexpected foo value %q. Expecting %q", u1.QueryArgs().Peek("foo"), nil)
	}
}

func TestURIPathOriginal(t *testing.T) {
	var u URI

	u.Parse(nil, []byte("/aa/bb/../c/./d/e/ff/g/../h"))
	pathOriginal := u.PathOriginal()
	path := u.Path()
	if string(pathOriginal) != "/aa/bb/../c/./d/e/ff/g/../h" {
		t.Fatalf("unexpected PathOriginal: %q. Expecting %q", pathOriginal, "/aa/bb/../c/./d/e/ff/g/../h")
	}
	if string(path) != "/aa/c/d/e/ff/h" {
		t.Fatalf("unexpected Path: %q. Expecting %q", path, "/aa/c/d/e/ff/h")
	}

	u.DisablePathNormalizing = true
	u.Parse(nil, []byte("/aa/bb/../c/./d/e/ff/g/../h"))
	pathOriginal = u.PathOriginal()
	path = u.Path()
	if string(pathOriginal) != "/aa/bb/../c/./d/e/ff/g/../h" {
		t.Fatalf("unexpected PathOriginal: %q. Expecting %q", pathOriginal, "/aa/bb/../c/./d/e/ff/g/../h")
	}
	if string(path) != "/aa/bb/../c/./d/e/ff/g/../h" {
		t.Fatalf("unexpected Path: %q. Expecting %q", path, "/aa/bb/../c/./d/e/ff/g/../h")
	}
}

func TestURIPathTrailingSlash(t *testing.T) {
	var u URI

	u.SetPath("/foo")
	if string(u.Path()) != "/foo" {
		t.Fatalf("unexpected Path: %q. Expecting %q", u.Path(), "/foo")
	}

	u.SetPath("/foo/")
	if string(u.Path()) != "/foo/" {
		t.Fatalf("unexpected Path: %q. Expecting %q", u.Path(), "/foo/")
	}

	u.SetPath("/foo/bar")
	if string(u.Path()) != "/foo/bar" {
		t.Fatalf("unexpected Path: %q. Expecting %q", u.Path(), "/foo/bar")
	}

	u.SetPath("/foo/bar/")
	if string(u.Path()) != "/foo/bar/" {
		t.Fatalf("unexpected Path: %q. Expecting %q", u.Path(), "/foo/bar/")
	}

	u.SetPath("foo/bar/")
	if string(u.Path()) != "/foo/bar/" {
		t.Fatalf("unexpected Path: %q. Expecting %q", u.Path(), "/foo/bar/")
	}

	u.SetPath("foo/bar")
	if string(u.Path()) != "/foo/bar" {
		t.Fatalf("unexpected Path: %q. Expecting %q", u.Path(), "/foo/bar")
	}

	u.SetPath("foo/")
	if string(u.Path()) != "/foo/" {
		t.Fatalf("unexpected Path: %q. Expecting %q", u.Path(), "/foo/")
	}

	u.SetPath("foo")
	if string(u.Path()) != "/foo" {
		t.Fatalf("unexpected Path: %q. Expecting %q", u.Path(), "/foo")
	}
}

func TestURIHostPort(t *testing.T) {
	var u URI

	u.SetHost("foobar.com:8080")
	if string(u.Host()) != "foobar.com:8080" {
		t.Fatalf("unexpected host: %q. Expecting %q", u.Host(), "foobar.com:8080")
	}
	if string(u.FullURI()) != "http://foobar.com:8080/" {
		t.Fatalf("unexpected full uri: %q. Expecting %q", u.FullURI(), "http://foobar.com:8080/")
	}

	u.SetHost("foobar.com")
	if string(u.Host()) != "foobar.com" {
		t.Fatalf("unexpected host: %q. Expecting %q", u.Host(), "foobar.com")
	}
	if string(u.FullURI()) != "http://foobar.com/" {
		t.Fatalf("unexpected full uri: %q. Expecting %q", u.FullURI(), "http://foobar.com/")
	}

	u.SetScheme("https")
	u.SetHost("foobar.com:443")
	if string(u.Host()) != "foobar.com:443" {
		t.Fatalf("unexpected host: %q. Expecting %q", u.Host(), "foobar.com:443")
	}
	if string(u.FullURI()) != "https://foobar.com:443/" {
		t.Fatalf("unexpected full uri: %q. Expecting %q", u.FullURI(), "https://foobar.com:443/")
	}

	u.SetScheme("http")
	if string(u.Host()) != "foobar.com:443" {
		t.Fatalf("unexpected host: %q. Expecting %q", u.Host(), "foobar.com:443")
	}
	if string(u.FullURI()) != "http://foobar.com:443/" {
		t.Fatalf("unexpected full uri: %q. Expecting %q", u.FullURI(), "http://foobar.com:443/")
	}
}

func TestURISetQueryString(t *testing.T) {
	var u URI

	u.SetQueryString("foo=bar&baz=123")
	if string(u.QueryString()) != "foo=bar&baz=123" {
		t.Fatalf("unexpected query string: %q. Expecting %q", u.QueryString(), "foo=bar&baz=123")
	}
	if string(u.QueryArgs().Peek("foo")) != "bar" {
		t.Fatalf("unexpected query args value: %q. Expecting %q", u.QueryArgs().Peek("foo"), "bar")
	}
	if string(u.QueryArgs().Peek("baz")) != "123" {
		t.Fatalf("unexpected query args value: %q. Expecting %q", u.QueryArgs().Peek("baz"), "123")
	}
	if string(u.QueryArgs().Peek("missing")) != "" {
		t.Fatalf("unexpected query args value: %q. Expecting %q", u.QueryArgs().Peek("missing"), "")
	}

	u.QueryArgs().Set("foo", "barbar")
	if string(u.QueryString()) != "foo=barbar&baz=123" {
		t.Fatalf("unexpected query string: %q. Expecting %q", u.QueryString(), "foo=barbar&baz=123")
	}
	if string(u.QueryArgs().Peek("foo")) != "barbar" {
		t.Fatalf("unexpected query args value: %q. Expecting %q", u.QueryArgs().Peek("foo"), "barbar")
	}
	if string(u.QueryArgs().Peek("baz")) != "123" {
		t.Fatalf("unexpected query args value: %q. Expecting %q", u.QueryArgs().Peek("baz"), "123")
	}
	if string(u.QueryArgs().Peek("missing")) != "" {
		t.Fatalf("unexpected query args value: %q. Expecting %q", u.QueryArgs().Peek("missing"), "")
	}

	u.SetQueryString("foo=bar&baz=123")
	if string(u.QueryArgs().Peek("foo")) != "bar" {
		t.Fatalf("unexpected query args value: %q. Expecting %q", u.QueryArgs().Peek("foo"), "bar")
	}
	if string(u.QueryArgs().Peek("baz")) != "123" {
		t.Fatalf("unexpected query args value: %q. Expecting %q", u.QueryArgs().Peek("baz"), "123")
	}
	if string(u.QueryArgs().Peek("missing")) != "" {
		t.Fatalf("unexpected query args value: %q. Expecting %q", u.QueryArgs().Peek("missing"), "")
	}

	u.QueryArgs().SetBytesV("baz", []byte("432"))
	if string(u.QueryString()) != "foo=bar&baz=432" {
		t.Fatalf("unexpected query string: %q. Expecting %q", u.QueryString(), "foo=bar&baz=432")
	}
	if string(u.QueryArgs().Peek("foo")) != "bar" {
		t.Fatalf("unexpected query args value: %q. Expecting %q", u.QueryArgs().Peek("foo"), "bar")
	}
	if string(u.QueryArgs().Peek("baz")) != "432" {
		t.Fatalf("unexpected query args value: %q. Expecting %q", u.QueryArgs().Peek("baz"), "432")
	}
	if string(u.QueryArgs().Peek("missing")) != "" {
		t.Fatalf("unexpected query args value: %q. Expecting %q", u.QueryArgs().Peek("missing"), "")
	}
}

func TestURIMinimal(t *testing.T) {
	testURIMinimal(t, "", "http://foobar.com/")
	testURIMinimal(t, "/", "http://foobar.com/")
	testURIMinimal(t, "aaa", "http://foobar.com/aaa")
	testURIMinimal(t, "foo/bar", "http://foobar.com/foo/bar")
	testURIMinimal(t, "https://aaa.com", "https://aaa.com/")
	testURIMinimal(t, "http://a/", "http://a/")
	testURIMinimal(t, "http://a:123/", "http://a:123/")
	testURIMinimal(t, "http://foo.com", "http://foo.com/")
	testURIMinimal(t, "http://foo.com/bar", "http://foo.com/bar")
	testURIMinimal(t, "http://foo.com/bar?baz", "http://foo.com/bar?baz")
	testURIMinimal(t, "http://foo.com/bar#baz", "http://foo.com/bar#baz")
	testURIMinimal(t, "http://foo.com/bar?baz#quux", "http://foo.com/bar?baz#quux")
	testURIMinimal(t, "http://foo.com/bar?baz=quux", "http://foo.com/bar?baz=quux")
	testURIMinimal(t, "http://foo.com/bar?baz=quux#quuux", "http://foo.com/bar?baz=quux#quuux")
	testURIMinimal(t, "/%2F", "http://foobar.com/%2F")
	testURIMinimal(t, "http://user@foo.com/bar", "http://user@foo.com/bar")
	testURIMinimal(t, "http://user:password@foo.com/bar", "http://user:password@foo.com/bar")
}

func testURIMinimal(t *testing.T, uri, expectedMinimalURI string) {
	t.Helper()

	var u URI

	u.Parse([]byte("foobar.com"), []byte(uri))
	if string(u.FullURI()) != expectedMinimalURI {
		t.Fatalf("unexpected full uri %q. Expecting %q. uri=%q", u.FullURI(), expectedMinimalURI, uri)
	}

	u.DisablePathNormalizing = true
	u.Parse([]byte("foobar.com"), []byte(uri))
	if string(u.FullURI()) != expectedMinimalURI {
		t.Fatalf("unexpected full uri %q. Expecting %q. uri=%q", u.FullURI(), expectedMinimalURI, uri)
	}
}

func TestURIParse(t *testing.T) {
	testURIParseSuccess(t, "http://foobar.com", nil, nil, "foobar.com")
	testURIParseSuccess(t, "http://user@foobar.com", []byte("user"), nil, "foobar.com")
	testURIParseSuccess(t, "http://user:password@foobar.com", []byte("user"), []byte("password"), "foobar.com")
	testURIParseSuccess(t, "http://user:password@foobar.com:8080", []byte("user"), []byte("password"), "foobar.com:8080")
	testURIParseSuccess(t, "http://user:password@[::1]", []byte("user"), []byte("password"), "[::1]")
	testURIParseSuccess(t, "http://user:password@[::1]:8080", []byte("user"), []byte("password"), "[::1]:8080")

	testURIParseFailure(t, "http://user@")
	testURIParseFailure(t, "http://user@:")
	testURIParseFailure(t, "http://user@foobar.com@foobar.com")
	testURIParseFailure(t, "http://user:password@foobar.com@foobar.com")
	testURIParseFailure(t, "http://user:password@foobar.com:8080@foobar.com")
	testURIParseFailure(t, "http://user:password@[::1]@foobar.com")
	testURIParseFailure(t, "http://user:password@[::1]:8080@foobar.com")
}

func testURIParseSuccess(t *testing.T, uri string, expectedUsername, expectedPassword []byte, expectedHost string) {
	t.Helper()

	var u URI
	if err := u.Parse(nil, []byte(uri)); err != nil {
		t.Fatalf("unexpected error when parsing uri %q: %s", uri, err)
	}
	if !bytes.Equal(u.Username(), expectedUsername) {
		t.Fatalf("unexpected username %q for uri %q. Expecting %q", u.Username(), uri, expectedUsername)
	}
	if !bytes.Equal(u.Password(), expectedPassword) {
		t.Fatalf("unexpected password %q for uri %q. Expecting %q", u.Password(), uri, expectedPassword)
	}
	if string(u.Host()) != expectedHost {
		t.Fatalf("unexpected host %q for uri %q. Expecting %q", u.Host(), uri, expectedHost)
	}
}

func testURIParseFailure(t *testing.T, uri string) {
	t.Helper()

	var u URI
	if err := u.Parse(nil, []byte(uri)); err == nil {
		t.Fatalf("expecting error when parsing uri %q", uri)
	}
}

func TestURIHostParse(t *testing.T) {
	testURIHostParseSuccess(t, "http://[::1]", "[::1]")
	testURIHostParseSuccess(t, "http://[::1]:8080", "[::1]:8080")
	testURIHostParseSuccess(t, "http://[a:b:c::d:e:f]", "[a:b:c::d:e:f]")
	testURIHostParseSuccess(t, "http://[a:b:c::d:e:f]:8080", "[a:b:c::d:e:f]:8080")
	testURIHostParseSuccess(t, "http://[a:b:c::d:e:f%25en0]", "[a:b:c::d:e:f%25en0]")
	testURIHostParseSuccess(t, "http://[a:b:c::d:e:f%25en0]:8080", "[a:b:c::d:e:f%25en0]:8080")
	testURIHostParseSuccess(t, "http://[fe80::1]", "[fe80::1]")
	testURIHostParseSuccess(t, "http://[fe80::1%25en0]", "[fe80::1%25en0]")
	testURIHostParseSuccess(t, "http://[fe80::1]:8080", "[fe80::1]:8080")
	testURIHostParseSuccess(t, "http://[fe80::1%25en0]:8080", "[fe80::1%25en0]:8080")

	testURIHostParseFailure(t, "http://[::1%25]")
	testURIHostParseFailure(t, "http://[::1%25]8080")
	testURIHostParseFailure(t, "http://[::1%25foo]")
	testURIHostParseFailure(t, "http://[::1%25foo]8080")
}

func testURIHostParseSuccess(t *testing.T, uri string, expectedHost string) {
	t.Helper()

	var u URI
	if err := u.Parse(nil, []byte(uri)); err != nil {
		t.Fatalf("unexpected error when parsing uri %q: %s", uri, err)
	}
	if string(u.Host()) != expectedHost {
		t.Fatalf("unexpected host %q for uri %q. Expecting %q", u.Host(), uri, expectedHost)
	}
}

func testURIHostParseFailure(t *testing.T, uri string) {
	t.Helper()

	var u URI
	if err := u.Parse(nil, []byte(uri)); err == nil {
		t.Fatalf("expecting error when parsing uri %q", uri)
	}
}

func TestURIString(t *testing.T) {
	var u URI
	s := u.String()
	if len(s) == 0 {
		t.Fatalf("unexpected empty uri string")
	}
}

func TestURIAcquireRelease(t *testing.T) {
	u := AcquireURI()
	ReleaseURI(u)
}

func TestURIAppendBytes(t *testing.T) {
	var u URI
	u.Parse([]byte("foobar.com"), []byte("/aaa/bb?cc=dd&e"))
	dst := u.AppendBytes(nil)
	if string(dst) != string(u.FullURI()) {
		t.Fatalf("unexpected bytes %q. Expecting %q", dst, u.FullURI())
	}

	u.DisablePathNormalizing = true
	u.Parse([]byte("foobar.com"), []byte("/aaa/bb?cc=dd&e"))
	dst = u.AppendBytes(nil)
	if string(dst) != string(u.FullURI()) {
		t.Fatalf("unexpected bytes %q. Expecting %q", dst, u.FullURI())
	}
}

func TestURIWriteTo(t *testing.T) {
	var u URI
	u.Parse([]byte("foobar.com"), []byte("/aaa/bb?cc=dd&e"))
	var buf bytes.Buffer
	n, err := u.WriteTo(&buf)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if n != int64(buf.Len()) {
		t.Fatalf("unexpected n %d. Expecting %d", n, buf.Len())
	}
	if string(buf.Bytes()) != string(u.FullURI()) {
		t.Fatalf("unexpected bytes %q. Expecting %q", buf.Bytes(), u.FullURI())
	}

	u.DisablePathNormalizing = true
	u.Parse([]byte("foobar.com"), []byte("/aaa/bb?cc=dd&e"))
	buf.Reset()
	n, err = u.WriteTo(&buf)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if n != int64(buf.Len()) {
		t.Fatalf("unexpected n %d. Expecting %d", n, buf.Len())
	}
	if string(buf.Bytes()) != string(u.FullURI()) {
		t.Fatalf("unexpected bytes %q. Expecting %q", buf.Bytes(), u.FullURI())
	}
}

func TestURIQueryArgsCopyTo(t *testing.T) {
	var u URI
	u.Parse([]byte("foobar.com"), []byte("/aaa/bb?cc=dd&e"))

	dst := AcquireArgs()
	u.QueryArgs().CopyTo(dst)
	s := dst.String()
	if s != "cc=dd&e=" {
		t.Fatalf("unexpected query args %q. Expecting %q", s, "cc=dd&e=")
	}

	u.QueryArgs().Set("foo", "bar")
	dst = AcquireArgs()
	u.QueryArgs().CopyTo(dst)
	s = dst.String()
	if s != "cc=dd&e=&foo=bar" {
		t.Fatalf("unexpected query args %q. Expecting %q", s, "cc=dd&e=&foo=bar")
	}

	u.DisablePathNormalizing = true
	u.Parse([]byte("foobar.com"), []byte("/aaa/bb?cc=dd&e"))

	dst = AcquireArgs()
	u.QueryArgs().CopyTo(dst)
	s = dst.String()
	if s != "cc=dd&e=" {
		t.Fatalf("unexpected query args %q. Expecting %q", s, "cc=dd&e=")
	}

	u.QueryArgs().Set("foo", "bar")
	dst = AcquireArgs()
	u.QueryArgs().CopyTo(dst)
	s = dst.String()
	if s != "cc=dd&e=&foo=bar" {
		t.Fatalf("unexpected query args %q. Expecting %q", s, "cc=dd&e=&foo=bar")
	}
}

func TestURIRemoveDuplicateSlashes(t *testing.T) {
	uri := []byte("https://example.com//foo/bar")
	parsedURI := []byte("/foo/bar")
	var u URI
	u.Parse(nil, uri)
	if !bytes.Equal(u.Path(), parsedURI) {
		t.Fatalf("unexpected path: %q. Expecting %q", u.Path(), parsedURI)
	}
}

func TestURINormalizePathBackslash(t *testing.T) {
	testURINormalizePathBackslash(t, "/foo/bar")
	testURINormalizePathBackslash(t, "/foo/bar/..")
	testURINormalizePathBackslash(t, "/foo/bar/./..")
	testURINormalizePathBackslash(t, "/foo/bar/.././")
	testURINormalizePathBackslash(t, "/foo/bar/.././.")
	testURINormalizePathBackslash(t, "/foo/bar/../././")
	testURINormalizePathBackslash(t, "/foo/bar/../././.")

	testURINormalizePathBackslash(t, "/foo/bar")
	testURINormalizePathBackslash(t, "/foo/bar/")
	testURINormalizePathBackslash(t, "/foo/bar/..")
	testURINormalizePathBackslash(t, "/foo/bar/./..")
	testURINormalizePathBackslash(t, "/foo/bar/.././")
	testURINormalizePathBackslash(t, "/foo/bar/.././.")
	testURINormalizePathBackslash(t, "/foo/bar/../././")
	testURINormalizePathBackslash(t, "/foo/bar/../././.")

	testURINormalizePathBackslash(t, "/foo/bar/")
	testURINormalizePathBackslash(t, "/foo/bar//")
	testURINormalizePathBackslash(t, "/foo/bar///")
	testURINormalizePathBackslash(t, "/foo/bar////")
	testURINormalizePathBackslash(t, "/foo/bar/////")
	testURINormalizePathBackslash(t, "/foo/bar//////")
	testURINormalizePathBackslash(t, "/foo/bar///////")
	testURINormalizePathBackslash(t, "/foo/bar////////")
	testURINormalizePathBackslash(t, "/foo/bar/////////")

	testURINormalizePathBackslash(t, "/foo/./bar")
	testURINormalizePathBackslash(t, "/foo/../bar")
	testURINormalizePathBackslash(t, "/foo/././bar")
	testURINormalizePathBackslash(t, "/foo/./bar/.")
	testURINormalizePathBackslash(t, "/foo/./bar/./")
	testURINormalizePathBackslash(t, "/foo/./bar/./.")
	testURINormalizePathBackslash(t, "/foo/./bar/././")
	testURINormalizePathBackslash(t, "/foo/./bar/././.")
}

func testURINormalizePathBackslash(t *testing.T, path string) {
	t.Helper()

	var u URI

	u.SetPath(path)
	if string(u.Path()) != path {
		t.Fatalf("unexpected path %q. Expecting %q", u.Path(), path)
	}
}

func TestURIParseHost(t *testing.T) {
	testURIParseHostSuccess(t, "http://foo.com", "foo.com")
	testURIParseHostSuccess(t, "http://foo.com:80", "foo.com:80")
	testURIParseHostSuccess(t, "http://foo.com:8080", "foo.com:8080")
	testURIParseHostSuccess(t, "http://foo.com:8080:8080", "foo.com:8080:8080")
	testURIParseHostSuccess(t, "http://foo.com:8080:8080:8080", "foo.com:8080:8080:8080")
	testURIParseHostSuccess(t, "http://foo.com:8080:8080:8080:8080", "foo.com:8080:8080:8080:8080")
	testURIParseHostSuccess(t, "http://foo.com:8080:8080:8080:8080:8080", "foo.com:8080:8080:8080:8080:8080")
	testURIParseHostSuccess(t, "http://foo.com:8080:8080:8080:8080:8080:8080", "foo.com:8080:8080:8080:8080:8080:8080")

	testURIParseHostFailure(t, "http://foo.com:")
	testURIParseHostFailure(t, "http://foo.com:abc")
	testURIParseHostFailure(t, "http://foo.com:8080:")
	testURIParseHostFailure(t, "http://foo.com:8080:abc")
}

func testURIParseHostSuccess(t *testing.T, uri, expectedHost string) {
	t.Helper()

	var u URI
	if err := u.Parse(nil, []byte(uri)); err != nil {
		t.Fatalf("unexpected error when parsing uri %q: %s", uri, err)
	}
	if string(u.Host()) != expectedHost {
		t.Fatalf("unexpected host %q for uri %q. Expecting %q", u.Host(), uri, expectedHost)
	}
}

func testURIParseHostFailure(t *testing.T, uri string) {
	t.Helper()

	var u URI
	if err := u.Parse(nil, []byte(uri)); err == nil {
		t.Fatalf("expecting error when parsing uri %q", uri)
	}
}

func TestURILastPathSegment(t *testing.T) {
	testURILastPathSegment(t, "", "")
	testURILastPathSegment(t, "/", "")
	testURILastPathSegment(t, "foo", "foo")
	testURILastPathSegment(t, "foo/", "")
	testURILastPathSegment(t, "foo/bar", "bar")
	testURILastPathSegment(t, "foo/bar/", "")
	testURILastPathSegment(t, "/foobar.js", "foobar.js")
	testURILastPathSegment(t, "/foo/bar/baz.html", "baz.html")
	testURILastPathSegment(t, "/foo/bar/baz.html?q=1", "baz.html")
	testURILastPathSegment(t, "/foo/bar/baz.html?q=1&q2=2", "baz.html")
	testURILastPathSegment(t, "/foo/bar/baz.html?q=1&q2=2&", "baz.html")
	testURILastPathSegment(t, "/foo/bar/baz.html?q=1&q2=2&.", "baz.html")
}

func testURILastPathSegment(t *testing.T, uri, expectedLastSegment string) {
	t.Helper()

	var u URI
	u.SetPath(uri)
	if string(u.LastPathSegment()) != expectedLastSegment {
		t.Fatalf("unexpected last path segment %q. Expecting %q. path=%q", u.LastPathSegment(), expectedLastSegment, uri)
	}
}

func TestURIParseFullPath(t *testing.T) {
	testURIParseFullPathSuccess(t, "/foo/bar", "/foo/bar")
	testURIParseFullPathSuccess(t, "/foo/bar?q=1", "/foo/bar?q=1")
	testURIParseFullPathSuccess(t, "/foo/bar?q=1&q2=2", "/foo/bar?q=1&q2=2")
	testURIParseFullPathSuccess(t, "/foo/bar?q=1&q2=2&", "/foo/bar?q=1&q2=2&")
	testURIParseFullPathSuccess(t, "/foo/bar?q=1&q2=2&.", "/foo/bar?q=1&q2=2&.")

	testURIParseFullPathFailure(t, "foo/bar")
	testURIParseFullPathFailure(t, "foo/bar?q=1")
	testURIParseFullPathFailure(t, "foo/bar?q=1&q2=2")
	testURIParseFullPathFailure(t, "foo/bar?q=1&q2=2&")
	testURIParseFullPathFailure(t, "foo/bar?q=1&q2=2&.")
}

func testURIParseFullPathSuccess(t *testing.T, fullPath, expectedFullPath string) {
	t.Helper()

	var u URI
	if err := u.Parse(nil, []byte(fullPath)); err != nil {
		t.Fatalf("unexpected error when parsing fullPath %q: %s", fullPath, err)
	}
	if string(u.FullURI()) != expectedFullPath {
		t.Fatalf("unexpected full uri %q. Expecting %q. fullPath=%q", u.FullURI(), expectedFullPath, fullPath)
	}
}

func testURIParseFullPathFailure(t *testing.T, fullPath string) {
	t.Helper()

	var u URI
	if err := u.Parse(nil, []byte(fullPath)); err == nil {
		t.Fatalf("expecting error when parsing fullPath %q", fullPath)
	}
}

func TestURIParsePort(t *testing.T) {
	testURIParsePortSuccess(t, "http://foo.com", "")
	testURIParsePortSuccess(t, "http://foo.com:80", ":80")
	testURIParsePortSuccess(t, "http://foo.com:8080", ":8080")
	testURIParsePortSuccess(t, "http://foo.com:8080:8080", ":8080:8080")
	testURIParsePortSuccess(t, "http://foo.com:8080:8080:8080", ":8080:8080:8080")
	testURIParsePortSuccess(t, "http://foo.com:8080:8080:8080:8080", ":8080:8080:8080:8080")
	testURIParsePortSuccess(t, "http://foo.com:8080:8080:8080:8080:8080", ":8080:8080:8080:8080:8080")
	testURIParsePortSuccess(t, "http://foo.com:8080:8080:8080:8080:8080:8080", ":8080:8080:8080:8080:8080:8080")

	testURIParsePortFailure(t, "http://foo.com:")
	testURIParsePortFailure(t, "http://foo.com:abc")
	testURIParsePortFailure(t, "http://foo.com:8080:")
	testURIParsePortFailure(t, "http://foo.com:8080:abc")
}

func testURIParsePortSuccess(t *testing.T, uri, expectedPort string) {
	t.Helper()

	var u URI
	if err := u.Parse(nil, []byte(uri)); err != nil {
		t.Fatalf("unexpected error when parsing uri %q: %s", uri, err)
	}
	if string(u.QueryArgs().Peek("port")) != expectedPort {
		t.Fatalf("unexpected port %q for uri %q. Expecting %q", u.QueryArgs().Peek("port"), uri, expectedPort)
	}
}

func testURIParsePortFailure(t *testing.T, uri string) {
	t.Helper()

	var u URI
	if err := u.Parse(nil, []byte(uri)); err == nil {
		t.Fatalf("expecting error when parsing uri %q", uri)
	}
}

func TestURIParseHostIPv6(t *testing.T) {
	testURIParseHostIPv6Success(t, "http://[::1]", "[::1]")
	testURIParseHostIPv6Success(t, "http://[::1]:8080", "[::1]:8080")
	testURIParseHostIPv6Success(t, "http://[a:b:c::d:e:f]", "[a:b:c::d:e:f]")
	testURIParseHostIPv6Success(t, "http://[a:b:c::d:e:f]:8080", "[a:b:c::d:e:f]:8080")
	testURIParseHostIPv6Success(t, "http://[a:b:c::d:e:f%25en0]", "[a:b:c::d:e:f%25en0]")
	testURIParseHostIPv6Success(t, "http://[a:b:c::d:e:f%25en0]:8080", "[a:b:c::d:e:f%25en0]:8080")
	testURIParseHostIPv6Success(t, "http://[fe80::1]", "[fe80::1]")
	testURIParseHostIPv6Success(t, "http://[fe80::1%25en0]", "[fe80::1%25en0]")
	testURIParseHostIPv6Success(t, "http://[fe80::1]:8080", "[fe80::1]:8080")
	testURIParseHostIPv6Success(t, "http://[fe80::1%25en0]:8080", "[fe80::1%25en0]:8080")

	testURIParseHostIPv6Failure(t, "http://[::1%25]")
	testURIParseHostIPv6Failure(t, "http://[::1%25]8080")
	testURIParseHostIPv6Failure(t, "http://[::1%25foo]")
	testURIParseHostIPv6Failure(t, "http://[::1%25foo]8080")
}

func testURIParseHostIPv6Success(t *testing.T, uri, expectedHost string) {
	t.Helper()

	var u URI
	if err := u.Parse(nil, []byte(uri)); err != nil {
		t.Fatalf("unexpected error when parsing uri %q: %s", uri, err)
	}
	if string(u.Host()) != expectedHost {
		t.Fatalf("unexpected host %q for uri %q. Expecting %q", u.Host(), uri, expectedHost)
	}
}

func testURIParseHostIPv6Failure(t *testing.T, uri string) {
	t.Helper()

	var u URI
	if err := u.Parse(nil, []byte(uri)); err == nil {
		t.Fatalf("expecting error when parsing uri %q", uri)
	}
}

func TestURIAppendBytesIPv6(t *testing.T) {
	var u URI
	u.Parse([]byte("foobar.com"), []byte("/aaa/bb?cc=dd&e"))
	dst := u.AppendBytes(nil)
	if string(dst) != string(u.FullURI()) {
		t.Fatalf("unexpected bytes %q. Expecting %q", dst, u.FullURI())
	}

	u.DisablePathNormalizing = true
	u.Parse([]byte("foobar.com"), []byte("/aaa/bb?cc=dd&e"))
	dst = u.AppendBytes(nil)
	if string(dst) != string(u.FullURI()) {
		t.Fatalf("unexpected bytes %q. Expecting %q", dst, u.FullURI())
	}
}

func TestURIRemoveDuplicateSlashesIPv6(t *testing.T) {
	uri := []byte("https://example.com//foo/bar")
	parsedURI := []byte("/foo/bar")
	var u URI
	u.Parse(nil, uri)
	if !bytes.Equal(u.Path(), parsedURI) {
		t.Fatalf("unexpected path: %q. Expecting %q", u.Path(), parsedURI)
	}
}

func TestURIParseHostIPv6Multicast(t *testing.T) {
	testURIParseHostIPv6MulticastSuccess(t, "http://[ff02::1]", "[ff02::1]")
	testURIParseHostIPv6MulticastSuccess(t, "http://[ff02::1]:8080", "[ff02::1]:8080")
	testURIParseHostIPv6MulticastSuccess(t, "http://[ff02::2]", "[ff02::2]")
	testURIParseHostIPv6MulticastSuccess(t, "http://[ff02::2]:8080", "[ff02::2]:8080")

	testURIParseHostIPv6MulticastFailure(t, "http://[ff02::]")
	testURIParseHostIPv6MulticastFailure(t, "http://[ff02::]8080")
	testURIParseHostIPv6MulticastFailure(t, "http://[ff02::1%25]")
	testURIParseHostIPv6MulticastFailure(t, "http://[ff02::1%25]8080")
}

func testURIParseHostIPv6MulticastSuccess(t *testing.T, uri, expectedHost string) {
	t.Helper()

	var u URI
	if err := u.Parse(nil, []byte(uri)); err != nil {
		t.Fatalf("unexpected error when parsing uri %q: %s", uri, err)
	}
	if string(u.Host()) != expectedHost {
		t.Fatalf("unexpected host %q for uri %q. Expecting %q", u.Host(), uri, expectedHost)
	}
}

func testURIParseHostIPv6MulticastFailure(t *testing.T, uri string) {
	t.Helper()

	var u URI
	if err := u.Parse(nil, []byte(uri)); err == nil {
		t.Fatalf("expecting error when parsing uri %q", uri)
	}
}

func TestURIParseHostIPv6LinkLocal(t *testing.T) {
	testURIParseHostIPv6LinkLocalSuccess(t, "http://[fe80::1%25en0]", "[fe80::1%25en0]")
	testURIParseHostIPv6LinkLocalSuccess(t, "http://[fe80::1%25en0]:8080", "[fe80::1%25en0]:8080")
	testURIParseHostIPv6LinkLocalSuccess(t, "http://[fe80::1%25en1]", "[fe80::1%25en1]")
	testURIParseHostIPv6LinkLocalSuccess(t, "http://[fe80::1%25en1]:8080", "[fe80::1%25en1]:8080")

	testURIParseHostIPv6LinkLocalFailure(t, "http://[fe80::1]")
	testURIParseHostIPv6LinkLocalFailure(t, "http://[fe80::1]8080")
	testURIParseHostIPv6LinkLocalFailure(t, "http://[fe80::1%25]")
	testURIParseHostIPv6LinkLocalFailure(t, "http://[fe80::1%25]8080")
}

func testURIParseHostIPv6LinkLocalSuccess(t *testing.T, uri, expectedHost string) {
	t.Helper()

	var u URI
	if err := u.Parse(nil, []byte(uri)); err != nil {
		t.Fatalf("unexpected error when parsing uri %q: %s", uri, err)
	}
	if string(u.Host()) != expectedHost {
		t.Fatalf("unexpected host %q for uri %q. Expecting %q", u.Host(), uri, expectedHost)
	}
}

func testURIParseHostIPv6LinkLocalFailure(t *testing.T, uri string) {
	t.Helper()

	var u URI
	if err := u.Parse(nil, []byte(uri)); err == nil {
		t.Fatalf("expecting error when parsing uri %q", uri)
	}
}

func TestURIHostPortIPv6(t *testing.T) {
	var u URI

	u.SetHost("[::1]:8080")
	if string(u.Host()) != "[::1]:8080" {
		t.Fatalf("unexpected host: %q. Expecting %q", u.Host(), "[::1]:8080")
	}
	if string(u.FullURI()) != "http://[::1]:8080/" {
		t.Fatalf("unexpected full uri: %q. Expecting %q", u.FullURI(), "http://[::1]:8080/")
	}

	u.SetHost("[::1]")
	if string(u.Host()) != "[::1]" {
		t.Fatalf("unexpected host: %q. Expecting %q", u.Host(), "[::1]")
	}
	if string(u.FullURI()) != "http://[::1]/" {
		t.Fatalf("unexpected full uri: %q. Expecting %q", u.FullURI(), "http://[::1]/")
	}

	u.SetScheme("https")
	u.SetHost("[::1]:443")
	if string(u.Host()) != "[::1]:443" {
		t.Fatalf("unexpected host: %q. Expecting %q", u.Host(), "[::1]:443")
	}
	if string(u.FullURI()) != "https://[::1]:443/" {
		t.Fatalf("unexpected full uri: %q. Expecting %q", u.FullURI(), "https://[::1]:443/")
	}

	u.SetScheme("http")
	if string(u.Host()) != "[::1]:443" {
		t.Fatalf("unexpected host: %q. Expecting %q", u.Host(), "[::1]:443")
	}
	if string(u.FullURI()) != "http://[::1]:443/" {
		t.Fatalf("unexpected full uri: %q. Expecting %q", u.FullURI(), "http://[::1]:443/")
	}
}

func TestURIHostPortIPv6Loopback(t *testing.T) {
	var u URI

	u.SetHost("[::1]:8080")
	if string(u.Host()) != "[::1]:8080" {
		t.Fatalf("unexpected host: %q. Expecting %q", u.Host(), "[::1]:8080")
	}
	if string(u.FullURI()) != "http://[::1]:8080/" {
		t.Fatalf("unexpected full uri: %q. Expecting %q", u.FullURI(), "http://[::1]:8080/")
	}

	u.SetHost("[::1]")
	if string(u.Host()) != "[::1]" {
		t.Fatalf("unexpected host: %q. Expecting %q", u.Host(), "[::1]")
	}
	if string(u.FullURI()) != "http://[::1]/" {
		t.Fatalf("unexpected full uri: %q. Expecting %q", u.FullURI(), "http://[::1]/")
	}

	u.SetScheme("https")
	u.SetHost("[::1]:443")
	if string(u.Host()) != "[::1]:443" {
		t.Fatalf("unexpected host: %q. Expecting %q", u.Host(), "[::1]:443")
	}
	if string(u.FullURI()) != "https://[::1]:443/" {
		t.Fatalf("unexpected full uri: %q. Expecting %q", u.FullURI(), "https://[::1]:443/")
	}

	u.SetScheme("http")
	if string(u.Host()) != "[::1]:443" {
		t.Fatalf("unexpected host: %q. Expecting %q", u.Host(), "[::1]:443")
	}
	if string(u.FullURI()) != "http://[::1]:443/" {
		t.Fatalf("unexpected full uri: %q. Expecting %q", u.FullURI(), "http://[::1]:443/")
	}
}

func TestURIHostPortIPv6LinkLocal(t *testing.T) {
	var u URI

	u.SetHost("[fe80::1]:8080")
	if string(u.Host()) != "[fe80::1]:8080" {
		t.Fatalf("unexpected host: %q. Expecting %q", u.Host(), "[fe80::1]:8080")
	}
	if string(u.FullURI()) != "http://[fe80::1]:8080/" {
		t.Fatalf("unexpected full uri: %q. Expecting %q", u.FullURI(), "http://[fe80::1]:8080/")
	}

	u.SetHost("[fe80::1]")
	if string(u.Host()) != "[fe80::1]" {
		t.Fatalf("unexpected host: %q. Expecting %q", u.Host(), "[fe80::1]")
	}
	if string(u.FullURI()) != "http://[fe80::1]/" {
		t.Fatalf("unexpected full uri: %q. Expecting %q", u.FullURI(), "http://[fe80::1]/")
	}

	u.SetScheme("https")
	u.SetHost("[fe80::1]:443")
	if string(u.Host()) != "[fe80::1]:443" {
		t.Fatalf("unexpected host: %q. Expecting %q", u.Host(), "[fe80::1]:443")
	}
	if string(u.FullURI()) != "https://[fe80::1]:443/" {
		t.Fatalf("unexpected full uri: %q. Expecting %q", u.FullURI(), "https://[fe80::1]:443/")
	}

	u.SetScheme("http")
	if string(u.Host()) != "[fe80::1]:443" {
		t.Fatalf("unexpected host: %q. Expecting %q", u.Host(), "[fe80::1]:443")
	}
	if string(u.FullURI()) != "http://[fe80::1]:443/" {
		t.Fatalf("unexpected full uri: %q. Expecting %q", u.FullURI(), "http://[fe80::1]:443/")
	}
}

func TestURIHostPortIPv6Multicast(t *testing.T) {
	var u URI

	u.SetHost("[ff02::1]:8080")
	if string(u.Host()) != "[ff02::1]:8080" {
		t.Fatalf("unexpected host: %q. Expecting %q", u.Host(), "[ff02::1]:8080")
	}
	if string(u.FullURI()) != "http://[ff02::1]:8080/" {
		t.Fatalf("unexpected full uri: %q. Expecting %q", u.FullURI(), "http://[ff02::1]:8080/")
	}

	u.SetHost("[ff02::1]")
	if string(u.Host()) != "[ff02::1]" {
		t.Fatalf("unexpected host: %q. Expecting %q", u.Host(), "[ff02::1]")
	}
	if string(u.FullURI()) != "http://[ff02::1]/" {
		t.Fatalf("unexpected full uri: %q. Expecting %q", u.FullURI(), "http://[ff02::1]/")
	}

	u.SetScheme("https")
	u.SetHost("[ff02::1]:443")
	if string(u.Host()) != "[ff02::1]:443" {
		t.Fatalf("unexpected host: %q. Expecting %q", u.Host(), "[ff02::1]:443")
	}
	if string(u.FullURI()) != "https://[ff02::1]:443/" {
		t.Fatalf("unexpected full uri: %q. Expecting %q", u.FullURI(), "https://[ff02::1]:443/")
	}

	u.SetScheme("http")
	if string(u.Host()) != "[ff02::1]:443" {
		t.Fatalf("unexpected host: %q. Expecting %q", u.Host(), "[ff02::1]:443")
	}
	if string(u.FullURI()) != "http://[ff02::1]:443/" {
		t.Fatalf("unexpected full uri: %q. Expecting %q", u.FullURI(), "http://[ff02::1]:443/")
	}
}

func TestURIParseHostIPv6LinkLocalMulticast(t *testing.T) {
	testURIParseHostIPv6LinkLocalMulticastSuccess(t, "http://[fe80::1%25en0]", "[fe80::1%25en0]")
	testURIParseHostIPv6LinkLocalMulticastSuccess(t, "http://[fe80::1%25en0]:8080", "[fe80::1%25en0]:8080")
	testURIParseHostIPv6LinkLocalMulticastSuccess(t, "http://[fe80::1%25en1]", "[fe80::1%25en1]")
	testURIParseHostIPv6LinkLocalMulticastSuccess(t, "http://[fe80::1%25en1]:8080", "[fe80::1%25en1]:8080")

	testURIParseHostIPv6LinkLocalMulticastFailure(t, "http://[fe80::1]")
	testURIParseHostIPv6LinkLocalMulticastFailure(t, "http://[fe80::1]8080")
	testURIParseHostIPv6LinkLocalMulticastFailure(t, "http://[fe80::1%25]")
	testURIParseHostIPv6LinkLocalMulticastFailure(t, "http://[fe80::1%25]8080")
}

func testURIParseHostIPv6LinkLocalMulticastSuccess(t *testing.T, uri, expectedHost string) {
	t.Helper()

	var u URI
	if err := u.Parse(nil, []byte(uri)); err != nil {
		t.Fatalf("unexpected error when parsing uri %q: %s", uri, err)
	}
	if string(u.Host()) != expectedHost {
		t.Fatalf("unexpected host %q for uri %q. Expecting %q", u.Host(), uri, expectedHost)
	}
}

func testURIParseHostIPv6LinkLocalMulticastFailure(t *testing.T, uri string) {
	t.Helper()

	var u URI
	if err := u.Parse(nil, []byte(uri)); err == nil {
		t.Fatalf("expecting error when parsing uri %q", uri)
	}
}

func TestURIParseHostIPv6Multicast(t *testing.T) {
	testURIParseHostIPv6MulticastSuccess(t, "http://[ff02::1]", "[ff02::1]")
	testURIParseHostIPv6MulticastSuccess(t, "http://[ff02::1]:8080", "[ff02::1]:8080")
	testURIParseHostIPv6MulticastSuccess(t, "http://[ff02::2]", "[ff02::2]")
	testURIParseHostIPv6MulticastSuccess(t, "http://[ff02::2]:8080", "[ff02::2]:8080")

	testURIParseHostIPv6MulticastFailure(t, "http://[ff02::]")
	testURIParseHostIPv6MulticastFailure(t, "http://[ff02::]8080