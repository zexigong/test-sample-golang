// Copyright 2013 Julien Schmidt. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be found
// at https://github.com/julienschmidt/httprouter/blob/master/LICENSE

package gin

import (
	"net/url"
	"reflect"
	"strings"
	"testing"

	"github.com/gin-gonic/gin/internal/bytesconv"
)

func TestParams(t *testing.T) {
	ps := Params{
		Param{"param1", "value1"},
		Param{"param2", "value2"},
		Param{"param3", "value3"},
	}

	if ps.ByName("param1") != "value1" {
		t.Fail()
	}
	if ps.ByName("param2") != "value2" {
		t.Fail()
	}
	if ps.ByName("param3") != "value3" {
		t.Fail()
	}
	if ps.ByName("param4") != "" {
		t.Fail()
	}

	value, ok := ps.Get("param1")
	if value != "value1" {
		t.Errorf("Expected 'value1', got %s", value)
	}
	if !ok {
		t.Errorf("Expected 'true', got %t", ok)
	}

	value, ok = ps.Get("param2")
	if value != "value2" {
		t.Errorf("Expected 'value2', got %s", value)
	}
	if !ok {
		t.Errorf("Expected 'true', got %t", ok)
	}

	value, ok = ps.Get("param3")
	if value != "value3" {
		t.Errorf("Expected 'value3', got %s", value)
	}
	if !ok {
		t.Errorf("Expected 'true', got %t", ok)
	}

	value, ok = ps.Get("param4")
	if value != "" {
		t.Errorf("Expected '', got %s", value)
	}
	if ok {
		t.Errorf("Expected 'false', got %t", ok)
	}
}

type testRequests []struct {
	path       string
	nilHandler bool
	route      string
	ps         Params
}

func TestTreeAddAndGet(t *testing.T) {
	var trees methodTrees

	trees = addRoute(trees, "GET", "/", "root handler")
	trees = addRoute(trees, "GET", "/cmd/:tool/:sub", "cmdSub handler")
	trees = addRoute(trees, "GET", "/cmd/:tool/", "cmd handler")
	trees = addRoute(trees, "GET", "/cmd/whoami", "whoami handler")
	trees = addRoute(trees, "GET", "/cmd/whoami/root/", "root handler")
	trees = addRoute(trees, "GET", "/cmd/whoami/root/insert", "insert handler")
	trees = addRoute(trees, "GET", "/cmd/whoami/:name", "hello handler")
	trees = addRoute(trees, "GET", "/cmd/whoami/:name/hello", "hello2 handler")
	trees = addRoute(trees, "GET", "/cmd/whoami/root/:name", "rootName handler")
	trees = addRoute(trees, "GET", "/cmd/whoami/root1/:name", "rootName handler")
	trees = addRoute(trees, "GET", "/cmd/whoami/root12/:name", "rootName handler")
	trees = addRoute(trees, "GET", "/src/*filepath", "src handler")
	trees = addRoute(trees, "GET", "/search/:query", "search handler")
	trees = addRoute(trees, "GET", "/search/invalid", "searchInvalid handler")
	trees = addRoute(trees, "GET", "/user_:name", "user_ handler")
	trees = addRoute(trees, "GET", "/user_x", "userx handler")
	trees = addRoute(trees, "GET", "/user_:name/about", "userAbout handler")
	trees = addRoute(trees, "GET", "/files/:dir/*filepath", "files handler")
	trees = addRoute(trees, "GET", "/doc/", "doc handler")
	trees = addRoute(trees, "GET", "/doc/go_faq.html", "goFaq handler")
	trees = addRoute(trees, "GET", "/doc/go1.html", "go1 handler")
	trees = addRoute(trees, "GET", "/info/:user/public", "info handler")
	trees = addRoute(trees, "GET", "/info/:user/project/:project", "project handler")
	trees = addRoute(trees, "GET", "/aa", "aa handler")
	trees = addRoute(trees, "GET", "/a/aa", "a/aa handler")
	trees = addRoute(trees, "GET", "/a/ab", "a/ab handler")
	trees = addRoute(trees, "GET", "/a/\\:foo\\:bar", "a/:foo:bar handler")
	trees = addRoute(trees, "GET", "/a/:bar", "a/:bar handler")
	trees = addRoute(trees, "GET", "/a/\\:bar/:baz", "a/:bar/:baz handler")
	trees = addRoute(trees, "GET", "/a/\\:bar/baz", "a/:bar/baz handler")
	trees = addRoute(trees, "GET", "/a/b/:bar/:baz", "a/b/:bar/:baz handler")
	trees = addRoute(trees, "GET", "/a/b/:bar/baz", "a/b/:bar/baz handler")

	trees = addRoute(trees, "POST", "/c", "c handler")
	trees = addRoute(trees, "POST", "/c/:bar", "c/:bar handler")
	trees = addRoute(trees, "POST", "/c/:bar/baz", "c/:bar/baz handler")

	trees = addRoute(trees, "GET", "/d/foo\\:bar", "d/foo:bar handler")
	trees = addRoute(trees, "GET", "/d/foo\\:bar/baz", "d/foo:bar/baz handler")
	trees = addRoute(trees, "GET", "/d/foo:bar", "d/foo:bar handler")
	trees = addRoute(trees, "GET", "/d/foo:bar/baz", "d/foo:bar/baz handler")

	trees = addRoute(trees, "GET", "/e/:baz/foo\\:bar", "e/:baz/foo:bar handler")
	trees = addRoute(trees, "GET", "/e/:baz/foo\\:bar/baz", "e/:baz/foo:bar/baz handler")
	trees = addRoute(trees, "GET", "/e/:baz/foo:bar", "e/:baz/foo:bar handler")
	trees = addRoute(trees, "GET", "/e/:baz/foo:bar/baz", "e/:baz/foo:bar/baz handler")

	trees = addRoute(trees, "GET", "/f/:bar", "f/:bar handler")
	trees = addRoute(trees, "GET", "/f/:bar/bar", "f/:bar/bar handler")
	trees = addRoute(trees, "GET", "/f/:bar/foo\\:bar", "f/:bar/foo:bar handler")
	trees = addRoute(trees, "GET", "/f/:bar/foo\\:bar/baz", "f/:bar/foo:bar/baz handler")
	trees = addRoute(trees, "GET", "/f/:bar/foo:bar", "f/:bar/foo:bar handler")
	trees = addRoute(trees, "GET", "/f/:bar/foo:bar/baz", "f/:bar/foo:bar/baz handler")

	trees = addRoute(trees, "GET", "/g\\:foo", "g:foo handler")
	trees = addRoute(trees, "GET", "/g\\:foo/bar", "g:foo/bar handler")
	trees = addRoute(trees, "GET", "/g\\:foo/:bar", "g:foo/:bar handler")

	trees = addRoute(trees, "GET", "/h:foo", "h:foo handler")
	trees = addRoute(trees, "GET", "/h:foo/bar", "h:foo/bar handler")
	trees = addRoute(trees, "GET", "/h:foo/:bar", "h:foo/:bar handler")

	trees = addRoute(trees, "GET", "/i/:foo/\\:bar", "i/:foo/:bar handler")
	trees = addRoute(trees, "GET", "/i/:foo/\\:bar/baz", "i/:foo/:bar/baz handler")
	trees = addRoute(trees, "GET", "/i/:foo/:bar", "i/:foo/:bar handler")
	trees = addRoute(trees, "GET", "/i/:foo/:bar/baz", "i/:foo/:bar/baz handler")

	trees = addRoute(trees, "GET", "/j/:foo", "j/:foo handler")
	trees = addRoute(trees, "GET", "/j/:foo/bar", "j/:foo/bar handler")
	trees = addRoute(trees, "GET", "/j/:foo/:bar", "j/:foo/:bar handler")

	trees = addRoute(trees, "GET", "/k\\:foo", "k:foo handler")
	trees = addRoute(trees, "GET", "/k\\:foo/bar", "k:foo/bar handler")
	trees = addRoute(trees, "GET", "/k\\:foo/:bar", "k:foo/:bar handler")

	trees = addRoute(trees, "GET", "/l:foo", "l:foo handler")
	trees = addRoute(trees, "GET", "/l:foo/bar", "l:foo/bar handler")
	trees = addRoute(trees, "GET", "/l:foo/:bar", "l:foo/:bar handler")

	trees = addRoute(trees, "GET", "/m/:foo", "m/:foo handler")
	trees = addRoute(trees, "GET", "/m/:foo/bar", "m/:foo/bar handler")
	trees = addRoute(trees, "GET", "/m/:foo/:bar", "m/:foo/:bar handler")

	trees = addRoute(trees, "GET", "/n\\:foo", "n:foo handler")
	trees = addRoute(trees, "GET", "/n\\:foo/bar", "n:foo/bar handler")
	trees = addRoute(trees, "GET", "/n\\:foo/:bar", "n:foo/:bar handler")

	trees = addRoute(trees, "GET", "/o:foo", "o:foo handler")
	trees = addRoute(trees, "GET", "/o:foo/bar", "o:foo/bar handler")
	trees = addRoute(trees, "GET", "/o:foo/:bar", "o:foo/:bar handler")

	trees = addRoute(trees, "GET", "/p/:foo/\\:bar", "p/:foo/:bar handler")
	trees = addRoute(trees, "GET", "/p/:foo/\\:bar/baz", "p/:foo/:bar/baz handler")
	trees = addRoute(trees, "GET", "/p/:foo/:bar", "p/:foo/:bar handler")
	trees = addRoute(trees, "GET", "/p/:foo/:bar/baz", "p/:foo/:bar/baz handler")

	trees = addRoute(trees, "GET", "/q/:foo", "q/:foo handler")
	trees = addRoute(trees, "GET", "/q/:foo/bar", "q/:foo/bar handler")
	trees = addRoute(trees, "GET", "/q/:foo/:bar", "q/:foo/:bar handler")

	trees = addRoute(trees, "GET", "/r/:foo/:bar", "r/:foo/:bar handler")
	trees = addRoute(trees, "GET", "/r/:foo/:bar/baz", "r/:foo/:bar/baz handler")
	trees = addRoute(trees, "GET", "/r/:foo/\\:bar", "r/:foo/:bar handler")
	trees = addRoute(trees, "GET", "/r/:foo/\\:bar/baz", "r/:foo/:bar/baz handler")

	trees = addRoute(trees, "GET", "/s/:foo/bar", "s/:foo/bar handler")
	trees = addRoute(trees, "GET", "/s/:foo/:bar", "s/:foo/:bar handler")

	trees = addRoute(trees, "GET", "/t\\:foo", "t:foo handler")
	trees = addRoute(trees, "GET", "/t\\:foo/bar", "t:foo/bar handler")
	trees = addRoute(trees, "GET", "/t\\:foo/:bar", "t:foo/:bar handler")

	trees = addRoute(trees, "GET", "/u:foo", "u:foo handler")
	trees = addRoute(trees, "GET", "/u:foo/bar", "u:foo/bar handler")
	trees = addRoute(trees, "GET", "/u:foo/:bar", "u:foo/:bar handler")

	trees = addRoute(trees, "GET", "/v/:foo/\\:bar", "v/:foo/:bar handler")
	trees = addRoute(trees, "GET", "/v/:foo/\\:bar/baz", "v/:foo/:bar/baz handler")
	trees = addRoute(trees, "GET", "/v/:foo/:bar", "v/:foo/:bar handler")
	trees = addRoute(trees, "GET", "/v/:foo/:bar/baz", "v/:foo/:bar/baz handler")

	trees = addRoute(trees, "GET", "/w/:foo", "w/:foo handler")
	trees = addRoute(trees, "GET", "/w/:foo/bar", "w/:foo/bar handler")
	trees = addRoute(trees, "GET", "/w/:foo/:bar", "w/:foo/:bar handler")

	trees = addRoute(trees, "GET", "/x\\:foo", "x:foo handler")
	trees = addRoute(trees, "GET", "/x\\:foo/bar", "x:foo/bar handler")
	trees = addRoute(trees, "GET", "/x\\:foo/:bar", "x:foo/:bar handler")

	trees = addRoute(trees, "GET", "/y:foo", "y:foo handler")
	trees = addRoute(trees, "GET", "/y:foo/bar", "y:foo/bar handler")
	trees = addRoute(trees, "GET", "/y:foo/:bar", "y:foo/:bar handler")

	trees = addRoute(trees, "GET", "/z/:foo/\\:bar", "z/:foo/:bar handler")
	trees = addRoute(trees, "GET", "/z/:foo/\\:bar/baz", "z/:foo/:bar/baz handler")
	trees = addRoute(trees, "GET", "/z/:foo/:bar", "z/:foo/:bar handler")
	trees = addRoute(trees, "GET", "/z/:foo/:bar/baz", "z/:foo/:bar/baz handler")

	trees = addRoute(trees, "GET", "/a1:foo", "a1:foo handler")
	trees = addRoute(trees, "GET", "/a1:foo/bar", "a1:foo/bar handler")
	trees = addRoute(trees, "GET", "/a1:foo/:bar", "a1:foo/:bar handler")

	trees = addRoute(trees, "GET", "/b1/:foo", "b1/:foo handler")
	trees = addRoute(trees, "GET", "/b1/:foo/bar", "b1/:foo/bar handler")
	trees = addRoute(trees, "GET", "/b1/:foo/:bar", "b1/:foo/:bar handler")

	trees = addRoute(trees, "GET", "/c1\\:foo", "c1:foo handler")
	trees = addRoute(trees, "GET", "/c1\\:foo/bar", "c1:foo/bar handler")
	trees = addRoute(trees, "GET", "/c1\\:foo/:bar", "c1:foo/:bar handler")

	trees = addRoute(trees, "GET", "/d1:foo", "d1:foo handler")
	trees = addRoute(trees, "GET", "/d1:foo/bar", "d1:foo/bar handler")
	trees = addRoute(trees, "GET", "/d1:foo/:bar", "d1:foo/:bar handler")

	trees = addRoute(trees, "GET", "/e1/:foo/\\:bar", "e1/:foo/:bar handler")
	trees = addRoute(trees, "GET", "/e1/:foo/\\:bar/baz", "e1/:foo/:bar/baz handler")
	trees = addRoute(trees, "GET", "/e1/:foo/:bar", "e1/:foo/:bar handler")
	trees = addRoute(trees, "GET", "/e1/:foo/:bar/baz", "e1/:foo/:bar/baz handler")

	trees = addRoute(trees, "GET", "/f1/:foo", "f1/:foo handler")
	trees = addRoute(trees, "GET", "/f1/:foo/bar", "f1/:foo/bar handler")
	trees = addRoute(trees, "GET", "/f1/:foo/:bar", "f1/:foo/:bar handler")

	trees = addRoute(trees, "GET", "/g1\\:foo", "g1:foo handler")
	trees = addRoute(trees, "GET", "/g1\\:foo/bar", "g1:foo/bar handler")
	trees = addRoute(trees, "GET", "/g1\\:foo/:bar", "g1:foo/:bar handler")

	trees = addRoute(trees, "GET", "/h1:foo", "h1:foo handler")
	trees = addRoute(trees, "GET", "/h1:foo/bar", "h1:foo/bar handler")
	trees = addRoute(trees, "GET", "/h1:foo/:bar", "h1:foo/:bar handler")

	trees = addRoute(trees, "GET", "/i1/:foo/\\:bar", "i1/:foo/:bar handler")
	trees = addRoute(trees, "GET", "/i1/:foo/\\:bar/baz", "i1/:foo/:bar/baz handler")
	trees = addRoute(trees, "GET", "/i1/:foo/:bar", "i1/:foo/:bar handler")
	trees = addRoute(trees, "GET", "/i1/:foo/:bar/baz", "i1/:foo/:bar/baz handler")

	trees = addRoute(trees, "GET", "/j1/:foo", "j1/:foo handler")
	trees = addRoute(trees, "GET", "/j1/:foo/bar", "j1/:foo/bar handler")
	trees = addRoute(trees, "GET", "/j1/:foo/:bar", "j1/:foo/:bar handler")

	trees = addRoute(trees, "GET", "/k1\\:foo", "k1:foo handler")
	trees = addRoute(trees, "GET", "/k1\\:foo/bar", "k1:foo/bar handler")
	trees = addRoute(trees, "GET", "/k1\\:foo/:bar", "k1:foo/:bar handler")

	trees = addRoute(trees, "GET", "/l1:foo", "l1:foo handler")
	trees = addRoute(trees, "GET", "/l1:foo/bar", "l1:foo/bar handler")
	trees = addRoute(trees, "GET", "/l1:foo/:bar", "l1:foo/:bar handler")

	trees = addRoute(trees, "GET", "/m1/:foo", "m1/:foo handler")
	trees = addRoute(trees, "GET", "/m1/:foo/bar", "m1/:foo/bar handler")
	trees = addRoute(trees, "GET", "/m1/:foo/:bar", "m1/:foo/:bar handler")

	trees = addRoute(trees, "GET", "/n1\\:foo", "n1:foo handler")
	trees = addRoute(trees, "GET", "/n1\\:foo/bar", "n1:foo/bar handler")
	trees = addRoute(trees, "GET", "/n1\\:foo/:bar", "n1:foo/:bar handler")

	trees = addRoute(trees, "GET", "/o1:foo", "o1:foo handler")
	trees = addRoute(trees, "GET", "/o1:foo/bar", "o1:foo/bar handler")
	trees = addRoute(trees, "GET", "/o1:foo/:bar", "o1:foo/:bar handler")

	trees = addRoute(trees, "GET", "/p1/:foo/\\:bar", "p1/:foo/:bar handler")
	trees = addRoute(trees, "GET", "/p1/:foo/\\:bar/baz", "p1/:foo/:bar/baz handler")
	trees = addRoute(trees, "GET", "/p1/:foo/:bar", "p1/:foo/:bar handler")
	trees = addRoute(trees, "GET", "/p1/:foo/:bar/baz", "p1/:foo/:bar/baz handler")

	trees = addRoute(trees, "GET", "/q1/:foo", "q1/:foo handler")
	trees = addRoute(trees, "GET", "/q1/:foo/bar", "q1/:foo/bar handler")
	trees = addRoute(trees, "GET", "/q1/:foo/:bar", "q1/:foo/:bar handler")

	trees = addRoute(trees, "GET", "/r1/:foo/:bar", "r1/:foo/:bar handler")
	trees = addRoute(trees, "GET", "/r1/:foo/:bar/baz", "r1/:foo/:bar/baz handler")
	trees = addRoute(trees, "GET", "/r1/:foo/\\:bar", "r1/:foo/:bar handler")
	trees = addRoute(trees, "GET", "/r1/:foo/\\:bar/baz", "r1/:foo/:bar/baz handler")

	trees = addRoute(trees, "GET", "/s1/:foo/bar", "s1/:foo/bar handler")
	trees = addRoute(trees, "GET", "/s1/:foo/:bar", "s1/:foo/:bar handler")

	trees = addRoute(trees, "GET", "/t1\\:foo", "t1:foo handler")
	trees = addRoute(trees, "GET", "/t1\\:foo/bar", "t1:foo/bar handler")
	trees = addRoute(trees, "GET", "/t1\\:foo/:bar", "t1:foo/:bar handler")

	trees = addRoute(trees, "GET", "/u1:foo", "u1:foo handler")
	trees = addRoute(trees, "GET", "/u1:foo/bar", "u1:foo/bar handler")
	trees = addRoute(trees, "GET", "/u1:foo/:bar", "u1:foo/:bar handler")

	trees = addRoute(trees, "GET", "/v1/:foo/\\:bar", "v1/:foo/:bar handler")
	trees = addRoute(trees, "GET", "/v1/:foo/\\:bar/baz", "v1/:foo/:bar/baz handler")
	trees = addRoute(trees, "GET", "/v1/:foo/:bar", "v1/:foo/:bar handler")
	trees = addRoute(trees, "GET", "/v1/:foo/:bar/baz", "v1/:foo/:bar/baz handler")

	trees = addRoute(trees, "GET", "/w1/:foo", "w1/:foo handler")
	trees = addRoute(trees, "GET", "/w1/:foo/bar", "w1/:foo/bar handler")
	trees = addRoute(trees, "GET", "/w1/:foo/:bar", "w1/:foo/:bar handler")

	trees = addRoute(trees, "GET", "/x1\\:foo", "x1:foo handler")
	trees = addRoute(trees, "GET", "/x1\\:foo/bar", "x1:foo/bar handler")
	trees = addRoute(trees, "GET", "/x1\\:foo/:bar", "x1:foo/:bar handler")

	trees = addRoute(trees, "GET", "/y1:foo", "y1:foo handler")
	trees = addRoute(trees, "GET", "/y1:foo/bar", "y1:foo/bar handler")
	trees = addRoute(trees, "GET", "/y1:foo/:bar", "y1:foo/:bar handler")

	trees = addRoute(trees, "GET", "/z1/:foo/\\:bar", "z1/:foo/:bar handler")
	trees = addRoute(trees, "GET", "/z1/:foo/\\:bar/baz", "z1/:foo/:bar/baz handler")
	trees = addRoute(trees, "GET", "/z1/:foo/:bar", "z1/:foo/:bar handler")
	trees = addRoute(trees, "GET", "/z1/:foo/:bar/baz", "z1/:foo/:bar/baz handler")

	trees = addRoute(trees, "GET", "/a2:foo", "a2:foo handler")
	trees = addRoute(trees, "GET", "/a2:foo/bar", "a2:foo/bar handler")
	trees = addRoute(trees, "GET", "/a2:foo/:bar", "a2:foo/:bar handler")

	trees = addRoute(trees, "GET", "/b2/:foo", "b2/:foo handler")
	trees = addRoute(trees, "GET", "/b2/:foo/bar", "b2/:foo/bar handler")
	trees = addRoute(trees, "GET", "/b2/:foo/:bar", "b2/:foo/:bar handler")

	trees = addRoute(trees, "GET", "/c2\\:foo", "c2:foo handler")
	trees = addRoute(trees, "GET", "/c2\\:foo/bar", "c2:foo/bar handler")
	trees = addRoute(trees, "GET", "/c2\\:foo/:bar", "c2:foo/:bar handler")

	trees = addRoute(trees, "GET", "/d2:foo", "d2:foo handler")
	trees = addRoute(trees, "GET", "/d2:foo/bar", "d2:foo/bar handler")
	trees = addRoute(trees, "GET", "/d2:foo/:bar", "d2:foo/:bar handler")

	trees = addRoute(trees, "GET", "/e2/:foo/\\:bar", "e2/:foo/:bar handler")
	trees = addRoute(trees, "GET", "/e2/:foo/\\:bar/baz", "e2/:foo/:bar/baz handler")
	trees = addRoute(trees, "GET", "/e2/:foo/:bar", "e2/:foo/:bar handler")
	trees = addRoute(trees, "GET", "/e2/:foo/:bar/baz", "e2/:foo/:bar/baz handler")

	trees = addRoute(trees, "GET", "/f2/:foo", "f2/:foo handler")
	trees = addRoute(trees, "GET", "/f2/:foo/bar", "f2/:foo/bar handler")
	trees = addRoute(trees, "GET", "/f2/:foo/:bar", "f2/:foo/:bar handler")

	trees = addRoute(trees, "GET", "/g2\\:foo", "g2:foo handler")
	trees = addRoute(trees, "GET", "/g2\\:foo/bar", "g2:foo/bar handler")
	trees = addRoute(trees, "GET", "/g2\\:foo/:bar", "g2:foo/:bar handler")

	trees = addRoute(trees, "GET", "/h2:foo", "h2:foo handler")
	trees = addRoute(trees, "GET", "/h2:foo/bar", "h2:foo/bar handler")
	trees = addRoute(trees, "GET", "/h2:foo/:bar", "h2:foo/:bar handler")

	trees = addRoute(trees, "GET", "/i2/:foo/\\:bar", "i2/:foo/:bar handler")
	trees = addRoute(trees, "GET", "/i2/:foo/\\:bar/baz", "i2/:foo/:bar/baz handler")
	trees = addRoute(trees, "GET", "/i2/:foo/:bar", "i2/:foo/:bar handler")
	trees = addRoute(trees, "GET", "/i2/:foo/:bar/baz", "i2/:foo/:bar/baz handler")

	trees = addRoute(trees, "GET", "/j2/:foo", "j2/:foo handler")
	trees = addRoute(trees, "GET", "/j2/:foo/bar", "j2/:foo/bar handler")
	trees = addRoute(trees, "GET", "/j2/:foo/:bar", "j2/:foo/:bar handler")

	trees = addRoute(trees, "GET", "/k2\\:foo", "k2:foo handler")
	trees = addRoute(trees, "GET", "/k2\\:foo/bar", "k2:foo/bar handler")
	trees = addRoute(trees, "GET", "/k2\\:foo/:bar", "k2:foo/:bar handler")

	trees = addRoute(trees, "GET", "/l2:foo", "l2:foo handler")
	trees = addRoute(trees, "GET", "/l2:foo/bar", "l2:foo/bar handler")
	trees = addRoute(trees, "GET", "/l2:foo/:bar", "l2:foo/:bar handler")

	trees = addRoute(trees, "GET", "/m2/:foo", "m2/:foo handler")
	trees = addRoute(trees, "GET", "/m2/:foo/bar", "m2/:foo/bar handler")
	trees = addRoute(trees, "GET", "/m2/:foo/:bar", "m2/:foo/:bar handler")

	trees = addRoute(trees, "GET", "/n2\\:foo", "n2:foo handler")
	trees = addRoute(trees, "GET", "/n2\\:foo/bar", "n2:foo/bar handler")
	trees = addRoute(trees, "GET", "/n2\\:foo/:bar", "n2:foo/:bar handler")

	trees = addRoute(trees, "GET", "/o2:foo", "o2:foo handler")
	trees = addRoute(trees, "GET", "/o2:foo/bar", "o2:foo/bar handler")
	trees = addRoute(trees, "GET", "/o2:foo/:bar", "o2:foo/:bar handler")

	trees = addRoute(trees, "GET", "/p2/:foo/\\:bar", "p2/:foo/:bar handler")
	trees = addRoute(trees, "GET", "/p2/:foo/\\:bar/baz", "p2/:foo/:bar/baz handler")
	trees = addRoute(trees, "GET", "/p2/:foo/:bar", "p2/:foo/:bar handler")
	trees = addRoute(trees, "GET", "/p2/:foo/:bar/baz", "p2/:foo/:bar/baz handler")

	trees = addRoute(trees, "GET", "/q2/:foo", "q2/:foo handler")
	trees = addRoute(trees, "GET", "/q2/:foo/bar", "q2/:foo/bar handler")
	trees = addRoute(trees, "GET", "/q2/:foo/:bar", "q2/:foo/:bar handler")

	trees = addRoute(trees, "GET", "/r2/:foo/:bar", "r2/:foo/:bar handler")
	trees = addRoute(trees, "GET", "/r2/:foo/:bar/baz", "r2/:foo/:bar/baz handler")
	trees = addRoute(trees, "GET", "/r2/:foo/\\:bar", "r2/:foo/:bar handler")
	trees = addRoute(trees, "GET", "/r2/:foo/\\:bar/baz", "r2/:foo/:bar/baz handler")

	trees = addRoute(trees, "GET", "/s2/:foo/bar", "s2/:foo/bar handler")
	trees = addRoute(trees, "GET", "/s2/:foo/:bar", "s2/:foo/:bar handler")

	trees = addRoute(trees, "GET", "/t2\\:foo", "t2:foo handler")
	trees = addRoute(trees, "GET", "/t2\\:foo/bar", "t2:foo/bar handler")
	trees = addRoute(trees, "GET", "/t2\\:foo/:bar", "t2:foo/:bar handler")

	trees = addRoute(trees, "GET", "/u2:foo", "u2:foo handler")
	trees = addRoute(trees, "GET", "/u2:foo/bar", "u2:foo/bar handler")
	trees = addRoute(trees, "GET", "/u2:foo/:bar", "u2:foo/:bar handler")

	trees = addRoute(trees, "GET", "/v2/:foo/\\:bar", "v2/:foo/:bar handler")
	trees = addRoute(trees, "GET", "/v2/:foo/\\:bar/baz", "v2/:foo/:bar/baz handler")
	trees = addRoute(trees, "GET", "/v2/:foo/:bar", "v2/:foo/:bar handler")
	trees = addRoute(trees, "GET", "/v2/:foo/:bar/baz", "v2/:foo/:bar/baz handler")

	trees = addRoute(trees, "GET", "/w2/:foo", "w2/:foo handler")
	trees = addRoute(trees, "GET", "/w2/:foo/bar", "w2/:foo/bar handler")
	trees = addRoute(trees, "GET", "/w2/:foo/:bar", "w2/:foo/:bar handler")

	trees = addRoute(trees, "GET", "/x2\\:foo", "x2:foo handler")
	trees = addRoute(trees, "GET", "/x2\\:foo/bar", "x2:foo/bar handler")
	trees = addRoute(trees, "GET", "/x2\\:foo/:bar", "x2:foo/:bar handler")

	trees = addRoute(trees, "GET", "/y2:foo", "y2:foo handler")
	trees = addRoute(trees, "GET", "/y2:foo/bar", "y2:foo/bar handler")
	trees = addRoute(trees, "GET", "/y2:foo/:bar", "y2:foo/:bar handler")

	trees = addRoute(trees, "GET", "/z2/:foo/\\:bar", "z2/:foo/:bar handler")
	trees = addRoute(trees, "GET", "/z2/:foo/\\:bar/baz", "z2/:foo/:bar/baz handler")
	trees = addRoute(trees, "GET", "/z2/:foo/:bar", "z2/:foo/:bar handler")
	trees = addRoute(trees, "GET", "/z2/:foo/:bar/baz", "z2/:foo/:bar/baz handler")

	trees = addRoute(trees, "GET", "/a3:foo", "a3:foo handler")
	trees = addRoute(trees, "GET", "/a3:foo/bar", "a3:foo/bar handler")
	trees = addRoute(trees, "GET", "/a3:foo/:bar", "a3:foo/:bar handler")

	trees = addRoute(trees, "GET", "/b3/:foo", "b3/:foo handler")
	trees = addRoute(trees, "GET", "/b3/:foo/bar", "b3/:foo/bar handler")
	trees = addRoute(trees, "GET", "/b3/:foo/:bar", "b3/:foo/:bar handler")

	trees = addRoute(trees, "GET", "/c3\\:foo", "c3:foo handler")
	trees = addRoute(trees, "GET", "/c3\\:foo/bar", "c3:foo/bar handler")
	trees = addRoute(trees, "GET", "/c3\\:foo/:bar", "c3:foo/:bar handler")

	trees = addRoute(trees, "GET", "/d3:foo", "d3:foo handler")
	trees = addRoute(trees, "GET", "/d3:foo/bar", "d3:foo/bar handler")
	trees = addRoute(trees, "GET", "/d3:foo/:bar", "d3:foo/:bar handler")

	trees = addRoute(trees, "GET", "/e3/:foo/\\:bar", "e3/:foo/:bar handler")
	trees = addRoute(trees, "GET", "/e3/:foo/\\:bar/baz", "e3/:foo/:bar/baz handler")
	trees = addRoute(trees, "GET", "/e3/:foo/:bar", "e3/:foo/:bar handler")
	trees = addRoute(trees, "GET", "/e3/:foo/:bar/baz", "e3/:foo/:bar/baz handler")

	trees = addRoute(trees, "GET", "/f3/:foo", "f3/:foo handler")
	trees = addRoute(trees, "GET", "/f3/:foo/bar", "f3/:foo/bar handler")
	trees = addRoute(trees, "GET", "/f3/:foo/:bar", "f3/:foo/:bar handler")

	trees = addRoute(trees, "GET", "/g3\\:foo", "g3:foo handler")
	trees = addRoute(trees, "GET", "/g3\\:foo/bar", "g3:foo/bar handler")
	trees = addRoute(trees, "GET", "/g3\\:foo/:bar", "g3:foo/:bar handler")

	trees = addRoute(trees, "GET", "/h3:foo", "h3:foo handler")
	trees = addRoute(trees, "GET", "/h3:foo/bar", "h3:foo/bar handler")
	trees = addRoute(trees, "GET", "/h3:foo/:bar", "h3:foo/:bar handler")

	trees = addRoute(trees, "GET", "/i3/:foo/\\:bar", "i3/:foo/:bar handler")
	trees = addRoute(trees, "GET", "/i3/:foo/\\:bar/baz", "i3/:foo/:bar/baz handler")
	trees = addRoute(trees, "GET", "/i3/:foo/:bar", "i3/:foo/:bar handler")
	trees = addRoute(trees, "GET", "/i3/:foo/:bar/baz", "i3/:foo/:bar/baz handler")

	trees = addRoute(trees, "GET", "/j3/:foo", "j3/:foo handler")
	trees = addRoute(trees, "GET", "/j3/:foo/bar", "j3/:foo/bar handler")
	trees = addRoute(trees, "GET", "/j3/:foo/:bar", "j3/:foo/:bar handler")

	trees = addRoute(trees, "GET", "/k3\\:foo", "k3:foo handler")
	trees = addRoute(trees, "GET", "/k3\\:foo/bar", "k3:foo/bar handler")
	trees = addRoute(trees, "GET", "/k3\\:foo/:bar", "k3:foo/:bar handler")

	trees = addRoute(trees, "GET", "/l3:foo", "l3:foo handler")
	trees = addRoute(trees, "GET", "/l3:foo/bar", "l3:foo/bar handler")
	trees = addRoute(trees, "GET", "/l3:foo/:bar", "l3:foo/:bar handler")

	trees = addRoute(trees, "GET", "/m3/:foo", "m3/:foo handler")
	trees = addRoute(trees, "GET", "/m3/:foo/bar", "m3/:foo/bar handler")
	trees = addRoute(trees, "GET", "/m3/:foo/:bar", "m3/:foo/:bar handler")

	trees = addRoute(trees, "GET", "/n3\\:foo", "n3:foo handler")
	trees = addRoute(trees, "GET", "/n3\\:foo/bar", "n3:foo/bar handler")
	trees = addRoute(trees, "GET", "/n3\\:foo/:bar", "n3:foo/:bar handler")

	trees = addRoute(trees, "GET", "/o3:foo", "o3:foo handler")
	trees = addRoute(trees, "GET", "/o3:foo/bar", "o3:foo/bar handler")
	trees = addRoute(trees, "GET", "/o3:foo/:bar", "o3:foo/:bar handler")

	trees = addRoute(trees, "GET", "/p3/:foo/\\:bar", "p3/:foo/:bar handler")
	trees = addRoute(trees, "GET", "/p3/:foo/\\:bar/baz", "p3/:foo/:bar/baz handler")
	trees = addRoute(trees, "GET", "/p3/:foo/:bar", "p3/:foo/:bar handler")
	trees = addRoute(trees, "GET", "/p3/:foo/:bar/baz", "p3/:foo/:bar/baz handler")

	trees = addRoute(trees, "GET", "/q3/:foo", "q3/:foo handler")
	trees = addRoute(trees, "GET", "/q3/:foo/bar", "q3/:foo/bar handler")
	trees = addRoute(trees, "GET", "/q3/:foo/:bar", "q3/:foo/:bar handler")

	trees = addRoute(trees, "GET", "/r3/:foo/:bar", "r3/:foo/:bar handler")
	trees = addRoute(trees, "GET", "/r3/:foo/:bar/baz", "r3/:foo/:bar/baz handler")
	trees = addRoute(trees, "GET", "/r3/:foo/\\:bar", "r3/:foo/:bar handler")
	trees = addRoute(trees, "GET", "/r3/:foo/\\:bar/baz", "r3/:foo/:bar/baz handler")

	trees = addRoute(trees, "GET", "/s3/:foo/bar", "s3/:foo/bar handler")
	trees = addRoute(trees, "GET", "/s3/:foo/:bar", "s3/:foo/:bar handler")

	trees = addRoute(trees, "GET", "/t3\\:foo", "t3:foo handler")
	trees = addRoute(trees, "GET", "/t3\\:foo/bar", "t3:foo/bar handler")
	trees = addRoute(trees, "GET", "/t3\\:foo/:bar", "t3:foo/:bar handler")

	trees = addRoute(trees, "GET", "/u3:foo", "u3:foo handler")
	trees = addRoute(trees, "GET", "/u3:foo/bar", "u3:foo/bar handler")
	trees = addRoute(trees, "GET", "/u3:foo/:bar", "u3:foo/:bar handler")

	trees = addRoute(trees, "GET", "/v3/:foo/\\:bar", "v3/:foo/:bar handler")
	trees = addRoute(trees, "GET", "/v3/:foo/\\:bar/baz", "v3/:foo/:bar/baz handler")
	trees = addRoute(trees, "GET", "/v3/:foo/:bar", "v3/:foo/:bar handler")
	trees = addRoute(trees, "GET", "/v3/:foo/:bar/baz", "v3/:foo/:bar/baz handler")

	trees = addRoute(trees, "GET", "/w3/:foo", "w3/:foo handler")
	trees = addRoute(trees, "GET", "/w3/:foo/bar", "w3/:foo/bar handler")
	trees = addRoute(trees, "GET", "/w3/:foo/:bar", "w3/:foo/:bar handler")

	trees = addRoute(trees, "GET", "/x3\\:foo", "x3:foo handler")
	trees = addRoute(trees, "GET", "/x3\\:foo/bar", "x3:foo/bar handler")
	trees = addRoute(trees, "GET", "/x3\\:foo/:bar", "x3:foo/:bar handler")

	trees = addRoute(trees, "GET", "/y3:foo", "y3:foo handler")
	trees = addRoute(trees, "GET", "/y3:foo/bar", "y3:foo/bar handler")
	trees = addRoute(trees, "GET", "/y3:foo/:bar", "y3:foo/:bar handler")

	trees = addRoute(trees, "GET", "/z3/:foo/\\:bar", "z3/:foo/:bar handler")
	trees = addRoute(trees, "GET", "/z3/:foo/\\:bar/baz", "z3/:foo/:bar/baz handler")
	trees = addRoute(trees, "GET", "/z3/:foo/:bar", "z3/:foo/:bar handler")
	trees = addRoute(trees, "GET", "/z3/:foo/:bar/baz", "z3/:foo/:bar/baz handler")

	trees = addRoute(trees, "GET", "/:foo", "/:foo handler")
	trees = addRoute(trees, "GET", "/:foo/bar", "/:foo/bar handler")
	trees = addRoute(trees, "GET", "/:foo/:bar", "/:foo/:bar handler")

	trees = addRoute(trees, "GET", "/\\:foo", "/:foo handler")
	trees = addRoute(trees, "GET", "/\\:foo/bar", "/:foo/bar handler")
	trees = addRoute(trees, "GET", "/\\:foo/:bar", "/:foo/:bar handler")

	trees = addRoute(trees, "GET", "/:foo/\\:bar", "/:foo/:bar handler")
	trees = addRoute(trees, "GET", "/:foo/\\:bar/baz", "/:foo/:bar/baz handler")
	trees = addRoute(trees, "GET", "/:foo/:bar", "/:foo/:bar handler")
	trees = addRoute(trees, "GET", "/:foo/:bar/baz", "/:foo/:bar/baz handler")

	requests := testRequests{
		{"/", false, "/", nil},
		{"/cmd/test/", false, "/cmd/:tool/", Params{Param{"tool", "test"}}},
		{"/cmd/test", false, "/cmd/:tool/", Params{Param{"tool", "test"}}},
		{"/cmd/test//", false, "/cmd/:tool/", Params{Param{"tool", "test"}}},
		{"/cmd/whoami", false, "/cmd/whoami", nil},
		{"/cmd/who", true, "", nil},
		{"/cmd/whoami/r", false, "/cmd/whoami/:name", Params{Param{"name", "r"}}},
		{"/cmd/whoami/root/", false, "/cmd/whoami/root/", nil},
		{"/cmd/whoami/root/insert", false, "/cmd/whoami/root/insert", nil},
		{"/cmd/whoami/root1", true, "", nil},
		{"/cmd/whoami/root1/", false, "/cmd/whoami/root1/:name", Params{Param{"name", ""}}},
		{"/cmd/whoami/root1/r", false, "/cmd/whoami/root1/:name", Params{Param{"name", "r"}}},
		{"/cmd/whoami/root1/r1", false, "/cmd/whoami/root1/:name", Params{Param{"name", "r1"}}},
		{"/cmd/whoami/root12", true, "", nil},
		{"/cmd/whoami/root12/", false, "/cmd/whoami/root12/:name", Params{Param{"name", ""}}},
		{"/cmd/whoami/root12/r", false, "/cmd/whoami/root12/:name", Params{Param{"name", "r"}}},
		{"/cmd/whoami/root12/r1", false, "/cmd/whoami/root12/:name", Params{Param{"name", "r1"}}},
		{"/src/", false, "/src/*filepath", Params{Param{"filepath", ""}}},
		{"/src/some/file.png", false, "/src/*filepath", Params{Param{"filepath", "some/file.png"}}},
		{"/search/", true, "", nil},
		{"/search/someth!ng+in+ünìcodé", false, "/search/:query", Params{Param{"query", "someth!ng+in+ünìcodé"}}},
		{"/search/someth!ng+in+ünìcodé/", false, "/search/:query", Params{Param{"query", "someth!ng+in+ünìcodé"}}},
		{"/search/someth!ng+in+ünìcodé/other", false, "/search/:query", Params{Param{"query", "someth!ng+in+ünìcodé/other"}}},
		{"/search/someth!ng+in+ünìcodé/other/", false, "/search/:query", Params{Param{"query", "someth!ng+in+ünìcodé/other"}}},
		{"/search/invalid", false, "/search/invalid", nil},
		{"/user_gopher", false, "/user_:name", Params{Param{"name", "gopher"}}},
		{"/user_gopher/about", false, "/user_:name/about", Params{Param{"name", "gopher"}}},
		{"/user_x", false, "/user_x", nil},
		{"/id/123", true, "", nil},
		{"/id/123/\\:bar", true, "", nil},
		{"/id/123/\\:bar/baz", true, "", nil},
		{"/id/123/:bar", true, "", nil},
		{"/id/123/:bar/baz", true, "", nil},
		{"/files/js/inc/framework.js", false, "/files/:dir/*filepath", Params{Param{"dir", "js"}, Param{"filepath", "inc/framework.js"}}},
		{"/info/test/public", false, "/info/:user/public", Params{Param{"user", "test"}}},
		{"/info/test/project/go", false, "/info/:user/project/:project", Params{Param{"user", "test"}, Param{"project", "go"}}},
		{"/a", true, "", nil},
		{"/a/", true, "", nil},
		{"/aa", false, "/aa", nil},
		{"/a/aa", false, "/a/aa", nil},
		{"/a/ab", false, "/a/ab", nil},
		{"/a/\\:foo\\:bar", false, "/a/\\:foo\\:bar", nil},
		{"/a/\\:foo\\:bar/", false, "/a/\\:foo\\:bar", nil},
		{"/a/\\:foo\\:bar/baz", false, "/a/\\:bar/:baz", Params{Param{"bar", "foo:bar"}, Param{"baz", "baz"}}},
		{"/a/\\:foo\\:bar/baz/", false, "/a/\\:bar/:baz", Params{Param{"bar", "foo:bar"}, Param{"baz", "baz"}}},
		{"/a/:bar", false, "/a/:bar", Params{Param{"bar", ":bar"}}},
		{"/a/:bar/", false, "/a/:bar", Params{Param{"bar", ":bar"}}},
		{"/a/:bar/baz", false, "/a/:bar/:baz", Params{Param{"bar", ":bar"}, Param{"baz", "baz"}}},
		{"/a/:bar/baz/", false, "/a/:bar/:baz", Params{Param{"bar", ":bar"}, Param{"baz", "baz"}}},
		{"/a/b/:bar/:baz", false, "/a/b/:bar/:baz", Params{Param{"bar", ":bar"}, Param{"baz", ":baz"}}},
		{"/a/b/:bar/:baz/", false, "/a/b/:bar/:baz", Params{Param{"bar", ":bar"}, Param{"baz", ":baz"}}},
		{"/a/b/:bar/baz", false, "/a/b/:bar/baz", Params{Param{"bar", ":bar"}}},
		{"/a/b/:bar/baz/", false, "/a/b/:bar/baz", Params{Param{"bar", ":bar"}}},
		{"/c", true, "", nil},
		{"/c/", true, "", nil},
		{"/c/\\:bar", true, "", nil},
		{"/c/\\:bar/", true, "", nil},
		{"/c/\\:bar/baz", true, "", nil},
		{"/c/\\:bar/baz/", true, "", nil},
		{"/c/:bar", false, "/c/:bar", Params{Param{"bar", ":bar"}}},
		{"/c/:bar/", false, "/c/:bar", Params{Param{"bar", ":bar"}}},
		{"/c/:bar/baz", false, "/c/:bar/baz", Params{Param{"bar", ":bar"}}},
		{"/c/:bar/baz/", false, "/c/:bar/baz", Params{Param{"bar", ":bar"}}},
		{"/d/foo\\:bar", false, "/d/foo\\:bar", nil},
		{"/d/foo\\:bar/", false, "/d/foo\\:bar", nil},
		{"/d/foo\\:bar/baz", false, "/d/foo\\:bar/baz", nil},
		{"/d/foo\\:bar/baz/", false, "/d/foo\\:bar/baz", nil},
		{"/d/foo:bar", false, "/d/foo:bar", nil},
		{"/d/foo:bar/", false, "/d/foo:bar", nil},
		{"/d/foo:bar/baz", false, "/d/foo:bar/baz", nil},
		{"/d/foo:bar/baz/", false, "/d/foo:bar/baz", nil},
		{"/e/:baz/foo\\:bar", false, "/e/:baz/foo\\:bar", Params{Param{"baz", ":baz"}}},
		{"/e/:baz/foo\\:bar/", false, "/e/:baz/foo\\:bar", Params{Param{"baz", ":baz"}}},
		{"/e/:baz/foo\\:bar/baz", false, "/e/:baz/foo\\:bar/baz", Params{Param{"baz", ":baz"}}},
		{"/e/:baz/foo\\:bar/baz/", false, "/e/:baz/foo\\:bar/baz", Params{Param{"baz", ":baz"}}},
		{"/e/:baz/foo:bar", false, "/e/:baz/foo:bar", Params{Param{"baz", ":baz"}}},
		{"/e/:baz/foo:bar/", false, "/e/:baz/foo:bar", Params{Param{"baz", ":baz"}}},
		{"/e/:baz/foo:bar/baz", false, "/e/:baz/foo:bar/baz", Params{Param{"baz", ":baz"}}},
		{"/e/:baz/foo:bar/baz/", false, "/e/:baz/foo:bar/baz", Params{Param{"baz", ":baz"}}},
		{"/f/:bar", false, "/f/:bar", Params{Param{"bar", ":bar"}}},
		{"/f/:bar/", false, "/f/:bar", Params{Param{"bar", ":bar"}}},
		{"/f/:bar/bar", false, "/f/:bar/bar", Params{Param{"bar", ":bar"}}},
		{"/f/:bar/bar/", false, "/f/:bar/bar", Params{Param{"bar", ":bar"}}},
		{"/f/:bar/foo\\:bar", false, "/f/:bar/foo\\:bar", Params{Param{"bar", ":bar"}}},
		{"/f/:bar/foo\\:bar/", false, "/f/:bar/foo\\:bar", Params{Param{"bar", ":bar"}}},
		{"/f/:bar/foo\\:bar/baz", false, "/f/:bar/foo\\:bar/baz", Params{Param{"bar", ":bar"}}},
		{"/f/:bar/foo\\:bar/baz/", false, "/f/:bar/foo\\:bar/baz", Params{Param{"bar", ":bar"}}},
		{"/f/:bar/foo:bar", false, "/f/:bar/foo:bar", Params{Param{"bar", ":bar"}}},
		{"/f/:bar/foo:bar/", false, "/f/:bar/foo:bar", Params{Param{"bar", ":bar"}}},
		{"/f/:bar/foo:bar/baz", false, "/f/:bar/foo:bar/baz", Params{Param{"bar", ":bar"}}},
		{"/f/:bar/foo:bar/baz/", false, "/f/:bar/foo:bar/baz", Params{Param{"bar", ":bar"}}},
		{"/g\\:foo", false, "/g\\:foo", nil},
		{"/g\\:foo/", false, "/g\\:foo", nil},
		{"/g\\:foo/bar", false, "/g\\:foo/bar", nil},
		{"/g\\:foo/bar/", false, "/g\\:foo/bar", nil},
		{"/g\\:foo/:bar", false, "/g\\:foo/:bar", Params{Param{"bar", ":bar"}}},
		{"/g\\:foo/:bar/", false, "/g\\:foo/:bar", Params{Param{"bar", ":bar"}}},
		{"/h:foo", false, "/h:foo", nil},
		{"/h:foo/", false, "/h:foo", nil},
		{"/h:foo/bar", false, "/h:foo/bar", nil},
		{"/h:foo/bar/", false, "/h:foo/bar", nil},
		{"/h:foo/:bar", false, "/h:foo/:bar", Params{Param{"bar", ":bar"}}},
		{"/h:foo/:bar/", false, "/h:foo/:bar", Params{Param{"bar", ":bar"}}},
		{"/i/:foo/\\:bar", false, "/i/:foo/\\:bar", Params{Param{"foo", ":foo"}}},
		{"/i/:foo/\\:bar/", false, "/i/:foo/\\:bar", Params{Param{"foo", ":foo"}}},
		{"/i/:foo/\\:bar/baz", false, "/i/:foo/\\:bar/baz", Params{Param{"foo", ":foo"}}},
		{"/i/:foo/\\:bar/baz/", false, "/i/:foo/\\:bar/baz", Params{Param{"foo", ":foo"}}},
		{"/i/:foo/:bar", false, "/i/:foo/:bar", Params{Param{"foo", ":foo"}, Param{"bar", ":bar"}}},
		{"/i/:foo/:bar/", false, "/i/:foo/:bar", Params{Param{"foo", ":foo"}, Param{"bar", ":bar"}}},
		{"/i/:foo/:bar/baz", false, "/i/:foo/:bar/baz", Params{Param{"foo", ":foo"}, Param{"bar", ":bar"}}},
		{"/i/:foo/:bar/baz/", false, "/i/:foo/:bar/baz", Params{Param{"foo", ":foo"}, Param{"bar", ":bar"}}},
		{"/j/:foo", false, "/j/:foo", Params{Param{"foo", ":foo"}}},
		{"/j/:foo/", false, "/j/:foo", Params{Param{"foo", ":foo"}}},
		{"/j/:foo/bar", false, "/j/:foo/bar", Params{Param{"foo", ":foo"}}},
		{"/j/:foo/bar/", false, "/j/:foo/bar", Params{Param{"foo", ":foo"}}},
		{"/j/:foo/:bar", false, "/j/:foo/:bar", Params{Param{"foo", ":foo"}, Param{"bar", ":bar"}}},
		{"/j/:foo/:bar/", false, "/j/:foo/:bar", Params{Param{"foo", ":foo"}, Param{"bar", ":bar"}}},
		{"/k\\:foo", false, "/k\\:foo", nil},
		{"/k\\:foo/", false, "/k\\:foo", nil},
		{"/k\\:foo/bar", false, "/k\\:foo/bar", nil},
		{"/k\\:foo/bar/", false, "/k\\:foo/bar", nil},
		{"/k\\:foo/:bar", false, "/k\\:foo/:bar", Params{Param{"bar", ":bar"}}},
		{"/k\\:foo/:bar/", false, "/k\\:foo/:bar", Params{Param{"bar", ":bar"}}},
		{"/l:foo", false, "/l:foo", nil},
		{"/l:foo/", false, "/l:foo", nil},
		{"/l:foo/bar", false, "/l:foo/bar", nil},
		{"/l:foo/bar/", false, "/l:foo/bar", nil},
		{"/l:foo/:bar", false, "/l:foo/:bar", Params{Param{"bar", ":bar"}}},
		{"/l:foo/:bar/", false, "/l:foo/:bar", Params{Param{"bar", ":bar"}}},
		{"/m/:foo", false, "/m/:foo", Params{Param{"foo", ":foo"}}},
		{"/m/:foo/", false, "/m/:foo", Params{Param{"foo", ":foo"}}},
		{"/m/:foo/bar", false, "/m/:foo/bar", Params{Param{"foo", ":foo"}}},
		{"/m/:foo/bar/", false, "/m/:foo/bar", Params{Param{"foo", ":foo"}}},
		{"/m/:foo/:bar", false, "/m/:foo/:bar", Params{Param{"foo", ":foo"}, Param{"bar", ":bar"}}},
		{"/m/:foo/:bar/", false, "/m/:foo/:bar", Params{Param{"foo", ":foo"}, Param{"bar", ":bar"}}},
		{"/n\\:foo", false, "/n\\:foo", nil},
		{"/n\\:foo/", false, "/n\\:foo", nil},
		{"/n\\:foo/bar", false, "/n\\:foo/bar", nil},
		{"/n\\:foo/bar/", false, "/n\\:foo/bar", nil},
		{"/n\\:foo/:bar", false, "/n\\:foo/:bar", Params{Param{"bar", ":bar"}}},
		{"/n\\:foo/:bar/", false, "/n\\:foo/:bar", Params{Param{"bar", ":bar"}}},
		{"/o:foo", false, "/o:foo", nil},
		{"/o:foo/", false, "/o:foo", nil},
		{"/o:foo/bar", false, "/o:foo/bar", nil},
		{"/o:foo/bar/", false, "/o:foo/bar", nil},
		{"/o:foo/:bar", false, "/o:foo/:bar", Params{Param{"bar", ":bar"}}},
		{"/o:foo/:bar/", false, "/o:foo/:bar", Params{Param{"bar", ":bar"}}},
		{"/p/:foo/\\:bar", false, "/p/:foo/\\:bar", Params{Param{"foo", ":foo"}}},
		{"/p/:foo/\\:bar/", false, "/p/:foo/\\:bar", Params{Param{"foo", ":foo"}}},
		{"/p/:foo/\\:bar/baz", false, "/p/:foo/\\:bar/baz", Params{Param{"foo", ":foo"}}},
		{"/p/:foo/\\:bar/baz/", false, "/p/:foo/\\:bar/baz", Params{Param{"foo", ":foo"}}},
		{"/p/:foo/:bar", false, "/p/:foo/:bar", Params{Param{"foo", ":foo"}, Param{"bar", ":bar"}}},
		{"/p/:foo/:bar/", false, "/p/:foo/:bar", Params{Param{"foo", ":foo"}, Param{"bar", ":bar"}}},
		{"/p/:foo/:bar/baz", false, "/p/:foo/:bar/baz", Params{Param{"foo", ":foo"}, Param{"bar", ":bar"}}},
		{"/p/:foo/:bar/baz/", false, "/p/:foo/:bar/baz", Params{Param{"foo", ":foo"}, Param{"bar", ":bar"}}},
		{"/q/:foo", false, "/q/:foo", Params{Param{"foo", ":foo"}}},
		{"/q/:foo/", false, "/q/:foo", Params{Param{"foo", ":foo"}}},
		{"/q/:foo/bar", false, "/q/:foo/bar", Params{Param{"foo", ":foo"}}},
		{"/q/:foo/bar/", false, "/q/:foo/bar", Params{Param{"foo", ":foo"}}},
		{"/q/:foo/:bar", false, "/q/:foo/:bar", Params{Param{"foo", ":foo"}, Param{"bar", ":bar"}}},
		{"/q/:foo/:bar/", false, "/q/:foo/:bar", Params{Param{"foo", ":foo"}, Param{"bar", ":bar"}}},
		{"/r/:foo/:bar", false, "/r/:foo/:bar", Params{Param{"foo", ":foo"}, Param{"bar", ":bar"}}},
		{"/r/:foo/:bar/", false, "/r/:foo/:bar", Params{Param{"foo", ":foo"}, Param{"bar", ":bar"}}},
		{"/r/:foo/:bar/baz", false, "/r/:foo/:bar/baz", Params{Param{"foo", ":foo"}, Param{"bar", ":bar"}}},
		{"/r/:foo/:bar/baz/", false, "/r/:foo/:bar/baz", Params{Param{"foo", ":foo"}, Param{"bar", ":bar"}}},
		{"/r/:foo/\\:bar", false, "/r/:foo/\\:bar", Params{Param{"foo", ":foo"}}},
		{"/r/:foo/\\:bar/", false, "/r/:foo/\\:bar", Params{Param{"foo", ":foo"}}},
		{"/r/:foo/\\:bar/baz", false, "/r/:foo/\\:bar/baz", Params{Param{"foo", ":foo"}}},
		{"/r/:foo/\\:bar/baz/", false, "/r/:foo/\\:bar/baz", Params{Param{"foo", ":foo"}}},
		{"/s/:foo/bar", false, "/s/:foo/bar", Params{Param{"foo", ":foo"}}},
		{"/s/:foo/bar/", false, "/s/:foo/bar", Params{Param{"foo", ":foo"}}},
		{"/s/:foo/:bar", false, "/s/:foo/:bar", Params{Param{"foo", ":foo"}, Param{"bar", ":bar"}}},
		{"/s/:foo/:bar/", false, "/s/:foo/:bar", Params{Param{"foo", ":foo"}, Param{"bar", ":bar"}}},
		{"/t\\:foo", false, "/t\\:foo", nil},
		{"/t\\:foo/", false, "/t\\:foo", nil},
		{"/t\\:foo/bar", false, "/t\\:foo/bar", nil},
		{"/t\\:foo/bar/", false, "/t\\:foo/bar", nil},
		{"/t\\:foo/:bar", false, "/t\\:foo/:bar", Params{Param{"bar", ":bar"}}},
		{"/t\\:foo/:bar/", false, "/t\\:foo/:bar", Params{Param{"bar", ":bar"}}},
		{"/u:foo", false, "/u:foo", nil},
		{"/u:foo/", false, "/u:foo", nil},
		{"/u:foo/bar", false, "/u:foo/bar", nil},
		{"/u:foo/bar/", false, "/u:foo/bar", nil},
		{"/u:foo/:bar", false, "/u:foo/:bar", Params{Param{"bar", ":bar"}}},
		{"/u:foo/:bar/", false, "/u:foo/:bar", Params{Param{"bar", ":bar"}}},
		{"/v/:foo/\\:bar", false, "/v/:foo/\\:bar", Params{Param{"foo", ":foo"}}},
		{"/v/:foo/\\:bar/", false, "/v/:foo/\\:bar", Params{Param{"foo", ":foo"}}},
		{"/v/:foo/\\:bar/baz", false, "/v/:foo/\\:bar/baz", Params{Param{"foo", ":foo"}}},
		{"/v/:foo/\\:bar/baz/", false, "/v/:foo/\\:bar/baz", Params{Param{"foo", ":foo"}}},
		{"/v/:foo/:bar", false, "/v/:foo/:bar", Params{Param{"foo", ":foo"}, Param{"bar", ":bar"}}},
		{"/v/:foo/:bar/", false, "/v/:foo/:bar", Params{Param{"foo", ":foo"}, Param{"bar", ":bar"}}},
		{"/v/:foo/:bar/baz", false, "/v/:foo/:bar/baz", Params{Param{"foo", ":foo"}, Param{"bar", ":bar"}}},
		{"/v/:foo/:bar/baz/", false, "/v/:foo/:bar/baz", Params{Param{"foo", ":foo"}, Param{"bar", ":bar"}}},
		{"/w/:foo", false, "/w/:foo", Params{Param{"foo", ":foo"}}},
		{"/w/:foo/", false, "/w/:foo", Params{Param{"foo", ":foo"}}},
		{"/w/:foo/bar", false, "/w/:foo/bar", Params{Param{"foo", ":foo"}}},
		{"/w/:foo/bar/", false, "/w/:foo/bar", Params{Param{"foo", ":foo"}}},
		{"/w/:foo/:bar", false, "/w/:foo/:bar", Params{Param{"foo", ":foo"}, Param{"bar", ":bar"}}},
		{"/w/:foo/:bar/", false, "/w/:foo/:bar", Params{Param{"foo", ":foo"}, Param{"bar", ":bar"}}},
		{"/x\\:foo", false, "/x\\:foo", nil},
		{"/x\\:foo/", false, "/