package fasthttp

import (
	"bytes"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"
	"testing"
)

func TestArgsWriteTo(t *testing.T) {
	var a Args
	a.Parse("aaa=bbb&ccc=ddd")
	var buf bytes.Buffer
	n, err := a.WriteTo(&buf)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if n != int64(buf.Len()) {
		t.Fatalf("unexpected n: %d. Expecting %d", n, buf.Len())
	}
	if buf.String() != "aaa=bbb&ccc=ddd" {
		t.Fatalf("unexpected result: %q. Expecting %q", buf.String(), "aaa=bbb&ccc=ddd")
	}
}

func TestArgsWriteToSingleArg(t *testing.T) {
	var a Args
	a.Add("foo", "bar")
	var buf bytes.Buffer
	n, err := a.WriteTo(&buf)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if n != int64(buf.Len()) {
		t.Fatalf("unexpected n: %d. Expecting %d", n, buf.Len())
	}
	if buf.String() != "foo=bar" {
		t.Fatalf("unexpected result: %q. Expecting %q", buf.String(), "foo=bar")
	}
}

func TestArgsWriteToNoArgs(t *testing.T) {
	var a Args
	var buf bytes.Buffer
	n, err := a.WriteTo(&buf)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if n != int64(buf.Len()) {
		t.Fatalf("unexpected n: %d. Expecting %d", n, buf.Len())
	}
	if buf.String() != "" {
		t.Fatalf("unexpected result: %q. Expecting %q", buf.String(), "")
	}
}

func TestArgsAddNoValue(t *testing.T) {
	var a Args
	a.Parse("foo=bar&baz")
	a.AddNoValue("foobar")
	if string(a.QueryString()) != "foo=bar&baz&foobar" {
		t.Fatalf("unexpected query string %q. Expecting %q", a.QueryString(), "foo=bar&baz&foobar")
	}
}

func TestArgsSetNoValue(t *testing.T) {
	var a Args
	a.Parse("foo=bar&baz")
	a.SetNoValue("foobar")
	if string(a.QueryString()) != "foo=bar&baz&foobar" {
		t.Fatalf("unexpected query string %q. Expecting %q", a.QueryString(), "foo=bar&baz&foobar")
	}

	a.Set("foo", "barbaz")
	if string(a.QueryString()) != "foo=barbaz&baz&foobar" {
		t.Fatalf("unexpected query string %q. Expecting %q", a.QueryString(), "foo=barbaz&baz&foobar")
	}

	a.SetNoValue("foo")
	if string(a.QueryString()) != "foo&baz&foobar" {
		t.Fatalf("unexpected query string %q. Expecting %q", a.QueryString(), "foo&baz&foobar")
	}
}

func TestArgsHas(t *testing.T) {
	var a Args
	a.Parse("foo=bar&baz")
	if !a.Has("foo") {
		t.Error("Has: expected foo to exist")
	}
	if !a.Has("baz") {
		t.Error("Has: expected baz to exist")
	}
	if a.Has("notexist") {
		t.Error("Has: expected notexist to not exist")
	}
}

func TestArgsSort(t *testing.T) {
	var a Args
	a.Parse("d=c&b=a&a=f&b=b&f=z")
	a.Sort(bytes.Compare)
	if string(a.QueryString()) != "a=f&b=a&b=b&d=c&f=z" {
		t.Fatalf("unexpected sort result: %q. Expecting %q", a.QueryString(), "a=f&b=a&b=b&d=c&f=z")
	}
}

func TestArgsSetUint(t *testing.T) {
	var a Args
	a.SetUint("foo", 12345)
	if v, err := a.GetUint("foo"); err != nil {
		t.Fatalf("unexpected error: %s", err)
	} else if v != 12345 {
		t.Fatalf("unexpected value: %d. Expecting %d", v, 12345)
	}

	a.SetUint("foo", 0)
	if v, err := a.GetUint("foo"); err != nil {
		t.Fatalf("unexpected error: %s", err)
	} else if v != 0 {
		t.Fatalf("unexpected value: %d. Expecting %d", v, 0)
	}
}

func TestArgsSetUintBytes(t *testing.T) {
	var a Args
	a.SetUintBytes([]byte("foo"), 12345)
	if v, err := a.GetUint("foo"); err != nil {
		t.Fatalf("unexpected error: %s", err)
	} else if v != 12345 {
		t.Fatalf("unexpected value: %d. Expecting %d", v, 12345)
	}

	a.SetUintBytes([]byte("foo"), 0)
	if v, err := a.GetUint("foo"); err != nil {
		t.Fatalf("unexpected error: %s", err)
	} else if v != 0 {
		t.Fatalf("unexpected value: %d. Expecting %d", v, 0)
	}
}

func TestArgsGetUintOrZero(t *testing.T) {
	var a Args
	a.SetUint("foo", 12345)
	if v := a.GetUintOrZero("foo"); v != 12345 {
		t.Fatalf("unexpected value: %d. Expecting %d", v, 12345)
	}

	a.SetUint("foo", 0)
	if v := a.GetUintOrZero("foo"); v != 0 {
		t.Fatalf("unexpected value: %d. Expecting %d", v, 0)
	}

	if v := a.GetUintOrZero("bar"); v != 0 {
		t.Fatalf("unexpected value: %d. Expecting %d", v, 0)
	}

	a.Set("foo", "bar")
	if v := a.GetUintOrZero("foo"); v != 0 {
		t.Fatalf("unexpected value: %d. Expecting %d", v, 0)
	}
}

func TestArgsGetUfloatOrZero(t *testing.T) {
	var a Args
	a.Set("foo", "123.456")
	if v := a.GetUfloatOrZero("foo"); v != 123.456 {
		t.Fatalf("unexpected value: %f. Expecting %f", v, 123.456)
	}

	a.Set("foo", "0")
	if v := a.GetUfloatOrZero("foo"); v != 0 {
		t.Fatalf("unexpected value: %f. Expecting %f", v, 0.)
	}

	if v := a.GetUfloatOrZero("bar"); v != 0 {
		t.Fatalf("unexpected value: %f. Expecting %f", v, 0.)
	}

	a.Set("foo", "bar")
	if v := a.GetUfloatOrZero("foo"); v != 0 {
		t.Fatalf("unexpected value: %f. Expecting %f", v, 0.)
	}
}

func TestArgsUint(t *testing.T) {
	var a Args
	a.Parse("a=123&b=foobar&c=322")
	var err error
	var v int

	v, err = a.GetUint("a")
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if v != 123 {
		t.Fatalf("unexpected value: %d. Expecting %d", v, 123)
	}

	v, err = a.GetUint("c")
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if v != 322 {
		t.Fatalf("unexpected value: %d. Expecting %d", v, 322)
	}

	v, err = a.GetUint("foobar")
	if err == nil {
		t.Fatalf("expecting error")
	}
	if v != -1 {
		t.Fatalf("unexpected value: %d. Expecting %d", v, -1)
	}

	v, err = a.GetUint("b")
	if err == nil {
		t.Fatalf("expecting error")
	}
	if v != -1 {
		t.Fatalf("unexpected value: %d. Expecting %d", v, -1)
	}

	v, err = a.GetUint("d")
	if err == nil {
		t.Fatalf("expecting error")
	}
	if v != -1 {
		t.Fatalf("unexpected value: %d. Expecting %d", v, -1)
	}
}

func TestArgsUfloat(t *testing.T) {
	var a Args
	a.Parse("a=123.456&b=foobar&c=322.22")
	var err error
	var v float64

	v, err = a.GetUfloat("a")
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if v != 123.456 {
		t.Fatalf("unexpected value: %f. Expecting %f", v, 123.456)
	}

	v, err = a.GetUfloat("c")
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if v != 322.22 {
		t.Fatalf("unexpected value: %f. Expecting %f", v, 322.22)
	}

	v, err = a.GetUfloat("foobar")
	if err == nil {
		t.Fatalf("expecting error")
	}
	if v != -1 {
		t.Fatalf("unexpected value: %f. Expecting %f", v, -1.0)
	}

	v, err = a.GetUfloat("b")
	if err == nil {
		t.Fatalf("expecting error")
	}
	if v != -1 {
		t.Fatalf("unexpected value: %f. Expecting %f", v, -1.0)
	}

	v, err = a.GetUfloat("d")
	if err == nil {
		t.Fatalf("expecting error")
	}
	if v != -1 {
		t.Fatalf("unexpected value: %f. Expecting %f", v, -1.0)
	}
}

func TestArgsBool(t *testing.T) {
	var a Args
	a.Parse("a=1&b=true&c=yes&d&f=False&g=0")
	if !a.GetBool("a") {
		t.Fatalf("unexpected value for a")
	}
	if !a.GetBool("b") {
		t.Fatalf("unexpected value for b")
	}
	if !a.GetBool("c") {
		t.Fatalf("unexpected value for c")
	}
	if !a.GetBool("d") {
		t.Fatalf("unexpected value for d")
	}
	if a.GetBool("f") {
		t.Fatalf("unexpected value for f")
	}
	if a.GetBool("g") {
		t.Fatalf("unexpected value for g")
	}
	if a.GetBool("h") {
		t.Fatalf("unexpected value for h")
	}
}

func TestArgsString(t *testing.T) {
	var a Args
	a.Parse("foo=x&bar=y&baz=z&zzz")
	s := a.String()
	if s != "foo=x&bar=y&baz=z&zzz" {
		t.Fatalf("unexpected args string %q. Expecting %q", s, "foo=x&bar=y&baz=z&zzz")
	}
}

func TestArgsParse(t *testing.T) {
	var a Args
	a.Parse("foo=bar&baz=sss&aaa=bbb")
	if v := string(a.Peek("foo")); v != "bar" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "bar")
	}
	if v := string(a.Peek("baz")); v != "sss" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "sss")
	}
	if v := string(a.Peek("aaa")); v != "bbb" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "bbb")
	}

	// test empty args
	a.Parse("")
	if v := string(a.Peek("foo")); v != "" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "")
	}

	// test empty key
	a.Parse("=bar")
	if v := string(a.Peek("")); v != "bar" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "bar")
	}
	if v := string(a.Peek("foo")); v != "" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "")
	}

	// test empty value
	a.Parse("foo=")
	if v := string(a.Peek("foo")); v != "" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "")
	}

	// test missing value
	a.Parse("foo")
	if v := string(a.Peek("foo")); v != "" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "")
	}
	if v := string(a.Peek("")); v != "" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "")
	}

	// test missing value
	a.Parse("&&&foo&&&bar&&")
	if v := string(a.Peek("foo")); v != "" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "")
	}
	if v := string(a.Peek("bar")); v != "" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "")
	}
	if v := string(a.Peek("")); v != "" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "")
	}

	// test empty key and value
	a.Parse("=")
	if v := string(a.Peek("")); v != "" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "")
	}
	if v := string(a.Peek("foo")); v != "" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "")
	}

	// test empty key and missing value
	a.Parse("===")
	if v := string(a.Peek("")); v != "==" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "==")
	}
	if v := string(a.Peek("foo")); v != "" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "")
	}
}

func TestArgsParseBytes(t *testing.T) {
	var a Args
	a.ParseBytes([]byte("foo=bar&baz=sss&aaa=bbb"))
	if v := string(a.Peek("foo")); v != "bar" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "bar")
	}
	if v := string(a.Peek("baz")); v != "sss" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "sss")
	}
	if v := string(a.Peek("aaa")); v != "bbb" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "bbb")
	}
}

func TestArgsParseError(t *testing.T) {
	a := &Args{}
	a.Parse("=bar")
	a.Parse("foo=")
	a.Parse("&&&foo&&&bar&&")
	a.Parse("===")
	a.Parse("foo=bar&baz=sss&aaa=bbb")
}

func TestArgsCopyTo(t *testing.T) {
	var a Args
	a.Parse("foo=bar&baz=sss&aaa=bbb")
	var b Args
	a.CopyTo(&b)
	if v := string(b.Peek("foo")); v != "bar" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "bar")
	}
	if v := string(b.Peek("baz")); v != "sss" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "sss")
	}
	if v := string(b.Peek("aaa")); v != "bbb" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "bbb")
	}
}

func TestArgsPeek(t *testing.T) {
	var a Args
	a.Parse("foo=bar&baz=sss&aaa=bbb")
	if v := string(a.Peek("foo")); v != "bar" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "bar")
	}
	if v := string(a.Peek("baz")); v != "sss" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "sss")
	}
	if v := string(a.Peek("aaa")); v != "bbb" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "bbb")
	}
	if v := string(a.Peek("xxxx")); v != "" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "")
	}
}

func TestArgsPeekBytes(t *testing.T) {
	var a Args
	a.Parse("foo=bar&baz=sss&aaa=bbb")
	if v := string(a.PeekBytes([]byte("foo"))); v != "bar" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "bar")
	}
	if v := string(a.PeekBytes([]byte("baz"))); v != "sss" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "sss")
	}
	if v := string(a.PeekBytes([]byte("aaa"))); v != "bbb" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "bbb")
	}
	if v := string(a.PeekBytes([]byte("xxxx"))); v != "" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "")
	}
}

func TestArgsPeekMulti(t *testing.T) {
	var a Args
	a.Parse("foo=bar&baz=sss&aaa=bbb&foo=baz&foo=aaa&foo=zzz")
	if v := a.PeekMulti("foo"); !testEq(v, [][]byte{[]byte("bar"), []byte("baz"), []byte("aaa"), []byte("zzz")}) {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, [][]byte{[]byte("bar"), []byte("baz"), []byte("aaa"), []byte("zzz")})
	}
	if v := a.PeekMulti("baz"); !testEq(v, [][]byte{[]byte("sss")}) {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, [][]byte{[]byte("sss")})
	}
	if v := a.PeekMulti("aaa"); !testEq(v, [][]byte{[]byte("bbb")}) {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, [][]byte{[]byte("bbb")})
	}
	if v := a.PeekMulti("xxxx"); !testEq(v, [][]byte{}) {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, [][]byte{})
	}
}

func TestArgsPeekMultiBytes(t *testing.T) {
	var a Args
	a.Parse("foo=bar&baz=sss&aaa=bbb&foo=baz&foo=aaa&foo=zzz")
	if v := a.PeekMultiBytes([]byte("foo")); !testEq(v, [][]byte{[]byte("bar"), []byte("baz"), []byte("aaa"), []byte("zzz")}) {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, [][]byte{[]byte("bar"), []byte("baz"), []byte("aaa"), []byte("zzz")})
	}
	if v := a.PeekMultiBytes([]byte("baz")); !testEq(v, [][]byte{[]byte("sss")}) {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, [][]byte{[]byte("sss")})
	}
	if v := a.PeekMultiBytes([]byte("aaa")); !testEq(v, [][]byte{[]byte("bbb")}) {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, [][]byte{[]byte("bbb")})
	}
	if v := a.PeekMultiBytes([]byte("xxxx")); !testEq(v, [][]byte{}) {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, [][]byte{})
	}
}

func testEq(a, b [][]byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if !bytes.Equal(a[i], b[i]) {
			return false
		}
	}
	return true
}

func TestArgsDel(t *testing.T) {
	var a Args
	a.Parse("foo=bar&baz=sss&aaa=bbb")
	a.Del("baz")
	if v := string(a.Peek("foo")); v != "bar" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "bar")
	}
	if v := string(a.Peek("baz")); v != "" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "")
	}
	if v := string(a.Peek("aaa")); v != "bbb" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "bbb")
	}

	a.Del("aaa")
	if v := string(a.Peek("foo")); v != "bar" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "bar")
	}
	if v := string(a.Peek("baz")); v != "" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "")
	}
	if v := string(a.Peek("aaa")); v != "" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "")
	}

	a.Del("foo")
	if v := string(a.Peek("foo")); v != "" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "")
	}
	if v := string(a.Peek("baz")); v != "" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "")
	}
	if v := string(a.Peek("aaa")); v != "" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "")
	}

	a.Del("xxx")
	if v := string(a.Peek("foo")); v != "" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "")
	}
	if v := string(a.Peek("baz")); v != "" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "")
	}
	if v := string(a.Peek("aaa")); v != "" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "")
	}
}

func TestArgsSet(t *testing.T) {
	var a Args
	a.Parse("foo=bar&baz=sss&aaa=bbb")
	a.Set("baz", "111")
	if v := string(a.Peek("foo")); v != "bar" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "bar")
	}
	if v := string(a.Peek("baz")); v != "111" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "111")
	}
	if v := string(a.Peek("aaa")); v != "bbb" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "bbb")
	}

	a.Set("aaa", "xxx")
	if v := string(a.Peek("foo")); v != "bar" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "bar")
	}
	if v := string(a.Peek("baz")); v != "111" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "111")
	}
	if v := string(a.Peek("aaa")); v != "xxx" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "xxx")
	}

	a.Set("foo", "barbar")
	if v := string(a.Peek("foo")); v != "barbar" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "barbar")
	}
	if v := string(a.Peek("baz")); v != "111" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "111")
	}
	if v := string(a.Peek("aaa")); v != "xxx" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "xxx")
	}

	a.Set("xxx", "yyy")
	if v := string(a.Peek("foo")); v != "barbar" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "barbar")
	}
	if v := string(a.Peek("baz")); v != "111" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "111")
	}
	if v := string(a.Peek("aaa")); v != "xxx" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "xxx")
	}
	if v := string(a.Peek("xxx")); v != "yyy" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "yyy")
	}
}

func TestArgsSetBytesK(t *testing.T) {
	var a Args
	a.Parse("foo=bar&baz=sss&aaa=bbb")
	a.SetBytesK([]byte("baz"), "111")
	if v := string(a.Peek("foo")); v != "bar" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "bar")
	}
	if v := string(a.Peek("baz")); v != "111" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "111")
	}
	if v := string(a.Peek("aaa")); v != "bbb" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "bbb")
	}

	a.SetBytesK([]byte("aaa"), "xxx")
	if v := string(a.Peek("foo")); v != "bar" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "bar")
	}
	if v := string(a.Peek("baz")); v != "111" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "111")
	}
	if v := string(a.Peek("aaa")); v != "xxx" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "xxx")
	}

	a.SetBytesK([]byte("foo"), "barbar")
	if v := string(a.Peek("foo")); v != "barbar" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "barbar")
	}
	if v := string(a.Peek("baz")); v != "111" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "111")
	}
	if v := string(a.Peek("aaa")); v != "xxx" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "xxx")
	}

	a.SetBytesK([]byte("xxx"), "yyy")
	if v := string(a.Peek("foo")); v != "barbar" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "barbar")
	}
	if v := string(a.Peek("baz")); v != "111" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "111")
	}
	if v := string(a.Peek("aaa")); v != "xxx" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "xxx")
	}
	if v := string(a.Peek("xxx")); v != "yyy" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "yyy")
	}
}

func TestArgsSetBytesV(t *testing.T) {
	var a Args
	a.Parse("foo=bar&baz=sss&aaa=bbb")
	a.SetBytesV("baz", []byte("111"))
	if v := string(a.Peek("foo")); v != "bar" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "bar")
	}
	if v := string(a.Peek("baz")); v != "111" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "111")
	}
	if v := string(a.Peek("aaa")); v != "bbb" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "bbb")
	}

	a.SetBytesV("aaa", []byte("xxx"))
	if v := string(a.Peek("foo")); v != "bar" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "bar")
	}
	if v := string(a.Peek("baz")); v != "111" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "111")
	}
	if v := string(a.Peek("aaa")); v != "xxx" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "xxx")
	}

	a.SetBytesV("foo", []byte("barbar"))
	if v := string(a.Peek("foo")); v != "barbar" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "barbar")
	}
	if v := string(a.Peek("baz")); v != "111" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "111")
	}
	if v := string(a.Peek("aaa")); v != "xxx" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "xxx")
	}

	a.SetBytesV("xxx", []byte("yyy"))
	if v := string(a.Peek("foo")); v != "barbar" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "barbar")
	}
	if v := string(a.Peek("baz")); v != "111" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "111")
	}
	if v := string(a.Peek("aaa")); v != "xxx" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "xxx")
	}
	if v := string(a.Peek("xxx")); v != "yyy" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "yyy")
	}
}

func TestArgsSetBytesKV(t *testing.T) {
	var a Args
	a.Parse("foo=bar&baz=sss&aaa=bbb")
	a.SetBytesKV([]byte("baz"), []byte("111"))
	if v := string(a.Peek("foo")); v != "bar" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "bar")
	}
	if v := string(a.Peek("baz")); v != "111" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "111")
	}
	if v := string(a.Peek("aaa")); v != "bbb" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "bbb")
	}

	a.SetBytesKV([]byte("aaa"), []byte("xxx"))
	if v := string(a.Peek("foo")); v != "bar" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "bar")
	}
	if v := string(a.Peek("baz")); v != "111" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "111")
	}
	if v := string(a.Peek("aaa")); v != "xxx" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "xxx")
	}

	a.SetBytesKV([]byte("foo"), []byte("barbar"))
	if v := string(a.Peek("foo")); v != "barbar" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "barbar")
	}
	if v := string(a.Peek("baz")); v != "111" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "111")
	}
	if v := string(a.Peek("aaa")); v != "xxx" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "xxx")
	}

	a.SetBytesKV([]byte("xxx"), []byte("yyy"))
	if v := string(a.Peek("foo")); v != "barbar" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "barbar")
	}
	if v := string(a.Peek("baz")); v != "111" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "111")
	}
	if v := string(a.Peek("aaa")); v != "xxx" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "xxx")
	}
	if v := string(a.Peek("xxx")); v != "yyy" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "yyy")
	}
}

func TestArgsSetCanonical(t *testing.T) {
	var a Args
	a.Set("foo", "bar")
	a.Set("baz", "sss")
	a.Set("aaa", "bbb")
	a.Set("baz", "111")
	if v := string(a.Peek("foo")); v != "bar" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "bar")
	}
	if v := string(a.Peek("baz")); v != "111" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "111")
	}
	if v := string(a.Peek("aaa")); v != "bbb" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "bbb")
	}

	a.Set("aaa", "xxx")
	if v := string(a.Peek("foo")); v != "bar" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "bar")
	}
	if v := string(a.Peek("baz")); v != "111" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "111")
	}
	if v := string(a.Peek("aaa")); v != "xxx" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "xxx")
	}

	a.Set("foo", "barbar")
	if v := string(a.Peek("foo")); v != "barbar" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "barbar")
	}
	if v := string(a.Peek("baz")); v != "111" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "111")
	}
	if v := string(a.Peek("aaa")); v != "xxx" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "xxx")
	}

	a.Set("xxx", "yyy")
	if v := string(a.Peek("foo")); v != "barbar" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "barbar")
	}
	if v := string(a.Peek("baz")); v != "111" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "111")
	}
	if v := string(a.Peek("aaa")); v != "xxx" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "xxx")
	}
	if v := string(a.Peek("xxx")); v != "yyy" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "yyy")
	}
}

func TestArgsVisitAll(t *testing.T) {
	var a Args
	a.Parse("foo=bar&baz=sss&aaa=bbb")
	a.VisitAll(func(k, v []byte) {
		switch string(k) {
		case "foo":
			if string(v) != "bar" {
				t.Fatalf("unexpected value %q. Expecting %q", v, "bar")
			}
		case "baz":
			if string(v) != "sss" {
				t.Fatalf("unexpected value %q. Expecting %q", v, "sss")
			}
		case "aaa":
			if string(v) != "bbb" {
				t.Fatalf("unexpected value %q. Expecting %q", v, "bbb")
			}
		default:
			t.Fatalf("unexpected key %q. Expecting %q", k, "foo/baz/aaa")
		}
	})
}

func TestArgsVisitAllBytes(t *testing.T) {
	var a Args
	a.Parse("foo=bar&baz=sss&aaa=bbb")
	a.VisitAll(func(k, v []byte) {
		switch string(k) {
		case "foo":
			if string(v) != "bar" {
				t.Fatalf("unexpected value %q. Expecting %q", v, "bar")
			}
		case "baz":
			if string(v) != "sss" {
				t.Fatalf("unexpected value %q. Expecting %q", v, "sss")
			}
		case "aaa":
			if string(v) != "bbb" {
				t.Fatalf("unexpected value %q. Expecting %q", v, "bbb")
			}
		default:
			t.Fatalf("unexpected key %q. Expecting %q", k, "foo/baz/aaa")
		}
	})
}

func TestArgsAdd(t *testing.T) {
	var a Args
	a.Add("foo", "bar")
	a.Add("baz", "sss")
	a.Add("aaa", "bbb")
	if v := string(a.Peek("foo")); v != "bar" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "bar")
	}
	if v := string(a.Peek("baz")); v != "sss" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "sss")
	}
	if v := string(a.Peek("aaa")); v != "bbb" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "bbb")
	}
}

func TestArgsAddBytesK(t *testing.T) {
	var a Args
	a.AddBytesK([]byte("foo"), "bar")
	a.AddBytesK([]byte("baz"), "sss")
	a.AddBytesK([]byte("aaa"), "bbb")
	if v := string(a.Peek("foo")); v != "bar" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "bar")
	}
	if v := string(a.Peek("baz")); v != "sss" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "sss")
	}
	if v := string(a.Peek("aaa")); v != "bbb" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "bbb")
	}
}

func TestArgsAddBytesV(t *testing.T) {
	var a Args
	a.AddBytesV("foo", []byte("bar"))
	a.AddBytesV("baz", []byte("sss"))
	a.AddBytesV("aaa", []byte("bbb"))
	if v := string(a.Peek("foo")); v != "bar" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "bar")
	}
	if v := string(a.Peek("baz")); v != "sss" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "sss")
	}
	if v := string(a.Peek("aaa")); v != "bbb" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "bbb")
	}
}

func TestArgsAddBytesKV(t *testing.T) {
	var a Args
	a.AddBytesKV([]byte("foo"), []byte("bar"))
	a.AddBytesKV([]byte("baz"), []byte("sss"))
	a.AddBytesKV([]byte("aaa"), []byte("bbb"))
	if v := string(a.Peek("foo")); v != "bar" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "bar")
	}
	if v := string(a.Peek("baz")); v != "sss" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "sss")
	}
	if v := string(a.Peek("aaa")); v != "bbb" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "bbb")
	}
}

func TestArgsAddCanonical(t *testing.T) {
	var a Args
	a.Add("foo", "bar")
	a.Add("baz", "sss")
	a.Add("aaa", "bbb")
	if v := string(a.Peek("foo")); v != "bar" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "bar")
	}
	if v := string(a.Peek("baz")); v != "sss" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "sss")
	}
	if v := string(a.Peek("aaa")); v != "bbb" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "bbb")
	}
}

func TestArgsSetGet(t *testing.T) {
	var a Args

	a.Set("foo", "bar")
	if v := string(a.Peek("foo")); v != "bar" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "bar")
	}

	a.Set("foo", "barbar")
	if v := string(a.Peek("foo")); v != "barbar" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "barbar")
	}

	a.Set("bar", "baz")
	if v := string(a.Peek("bar")); v != "baz" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "baz")
	}
	if v := string(a.Peek("foo")); v != "barbar" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "barbar")
	}

	a.Set("foo", "xxxxx")
	if v := string(a.Peek("foo")); v != "xxxxx" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "xxxxx")
	}
	if v := string(a.Peek("bar")); v != "baz" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "baz")
	}
}

func TestArgsDelSetGet(t *testing.T) {
	var a Args

	a.Del("foo")
	if v := string(a.Peek("foo")); v != "" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "")
	}

	a.Del("foo")
	a.Set("foo", "bar")
	if v := string(a.Peek("foo")); v != "bar" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "bar")
	}

	a.Del("foo")
	if v := string(a.Peek("foo")); v != "" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "")
	}
}

func TestArgsSetDel(t *testing.T) {
	var a Args

	a.Set("foo", "bar")
	a.Set("bar", "baz")
	a.Set("baz", "aaaaaa")
	if v := string(a.Peek("foo")); v != "bar" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "bar")
	}
	if v := string(a.Peek("bar")); v != "baz" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "baz")
	}
	if v := string(a.Peek("baz")); v != "aaaaaa" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "aaaaaa")
	}

	a.Del("bar")
	if v := string(a.Peek("foo")); v != "bar" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "bar")
	}
	if v := string(a.Peek("bar")); v != "" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "")
	}
	if v := string(a.Peek("baz")); v != "aaaaaa" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "aaaaaa")
	}

	a.Del("baz")
	if v := string(a.Peek("foo")); v != "bar" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "bar")
	}
	if v := string(a.Peek("bar")); v != "" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "")
	}
	if v := string(a.Peek("baz")); v != "" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "")
	}

	a.Del("foo")
	if v := string(a.Peek("foo")); v != "" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "")
	}
	if v := string(a.Peek("bar")); v != "" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "")
	}
	if v := string(a.Peek("baz")); v != "" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "")
	}
}

func TestArgsDelSetGet2(t *testing.T) {
	var a Args

	a.Del("foo")
	if v := string(a.Peek("foo")); v != "" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "")
	}

	a.Del("foo")
	a.Set("foo", "bar")
	if v := string(a.Peek("foo")); v != "bar" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "bar")
	}

	a.Del("foo")
	if v := string(a.Peek("foo")); v != "" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "")
	}

	a.Set("foo", "bar")
	if v := string(a.Peek("foo")); v != "bar" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "bar")
	}

	a.Del("foo")
	if v := string(a.Peek("foo")); v != "" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "")
	}
}

func TestArgsReset(t *testing.T) {
	var a Args

	a.Set("foo", "bar")
	a.Set("bar", "baz")
	a.Reset()
	if v := string(a.Peek("foo")); v != "" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "")
	}
	if v := string(a.Peek("bar")); v != "" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "")
	}
	if v := string(a.Peek("xxx")); v != "" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "")
	}

	a.Reset()
	if v := string(a.Peek("foo")); v != "" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "")
	}
	if v := string(a.Peek("bar")); v != "" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "")
	}
	if v := string(a.Peek("xxx")); v != "" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "")
	}

	a.Set("aaa", "bbb")
	if v := string(a.Peek("aaa")); v != "bbb" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "bbb")
	}
	a.Reset()
	if v := string(a.Peek("aaa")); v != "" {
		t.Fatalf("unexpected arg value %q. Expecting %q", v, "")
	}
}

func TestArgsAppendBytes(t *testing.T) {
	var a Args

	a.Set("foo", "bar")
	a.Set("bar", "baz")
	a.Set("baz", "aaaaaa")

	s := string(a.AppendBytes(nil))
	if s != "foo=bar&bar=baz&baz=aaaaaa" {
		t.Fatalf("unexpected result %q. Expecting %q", s, "foo=bar&bar=baz&baz=aaaaaa")
	}

	s = string(a.AppendBytes([]byte("xx")))
	if s != "xxfoo=bar&bar=baz&baz=aaaaaa" {
		t.Fatalf("unexpected result %q. Expecting %q", s, "xxfoo=bar&bar=baz&baz=aaaaaa")
	}

	s = string(a.AppendBytes([]byte("xx&")))
	if s != "xx&foo=bar&bar=baz&baz=aaaaaa" {
		t.Fatalf("unexpected result %q. Expecting %q", s, "xx&foo=bar&bar=baz&baz=aaaaaa")
	}

	s = string(a.AppendBytes([]byte("xx&=&")))
	if s != "xx&=&foo=bar&bar=baz&baz=aaaaaa" {
		t.Fatalf("unexpected result %q. Expecting %q", s, "xx&=&foo=bar&bar=baz&baz=aaaaaa")
	}

	s = string(a.AppendBytes([]byte("xx&=&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&")))
	if s != "xx&=&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&foo=bar&bar=baz&baz=aaaaaa" {
		t.Fatalf("unexpected result %q. Expecting %q", s, "xx&=&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&foo=bar&bar=baz&baz=aaaaaa")
	}
}

func TestArgsAppendBytesNoValue(t *testing.T) {
	var a Args

	a.Set("foo", "bar")
	a.SetNoValue("bar")
	a.Set("baz", "aaaaaa")

	s := string(a.AppendBytes(nil))
	if s != "foo=bar&bar&baz=aaaaaa" {
		t.Fatalf("unexpected result %q. Expecting %q", s, "foo=bar&bar&baz=aaaaaa")
	}

	s = string(a.AppendBytes([]byte("xx")))
	if s != "xxfoo=bar&bar&baz=aaaaaa" {
		t.Fatalf("unexpected result %q. Expecting %q", s, "xxfoo=bar&bar&baz=aaaaaa")
	}

	s = string(a.AppendBytes([]byte("xx&")))
	if s != "xx&foo=bar&bar&baz=aaaaaa" {
		t.Fatalf("unexpected result %q. Expecting %q", s, "xx&foo=bar&bar&baz=aaaaaa")
	}

	s = string(a.AppendBytes([]byte("xx&=&")))
	if s != "xx&=&foo=bar&bar&baz=aaaaaa" {
		t.Fatalf("unexpected result %q. Expecting %q", s, "xx&=&foo=bar&bar&baz=aaaaaa")
	}

	s = string(a.AppendBytes([]byte("xx&=&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&")))
	if s != "xx&=&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&foo=bar&bar&baz=aaaaaa" {
		t.Fatalf("unexpected result %q. Expecting %q", s, "xx&=&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&foo=bar&bar&baz=aaaaaa")
	}
}

func TestArgsStringCompose(t *testing.T) {
	var a Args

	a.Set("foo", "bar")
	a.Set("bar", "baz")
	a.Set("baz", "aaaaaa")

	s := string(a.QueryString())
	if s != "foo=bar&bar=baz&baz=aaaaaa" {
		t.Fatalf("unexpected result %q. Expecting %q", s, "foo=bar&bar=baz&baz=aaaaaa")
	}
}

func TestArgsStringComposeNoValue(t *testing.T) {
	var a Args

	a.Set("foo", "bar")
	a.SetNoValue("bar")
	a.Set("baz", "aaaaaa")

	s := string(a.QueryString())
	if s != "foo=bar&bar&baz=aaaaaa" {
		t.Fatalf("unexpected result %q. Expecting %q", s, "foo=bar&bar&baz=aaaaaa")
	}
}

func TestArgsStringComposeDelete(t *testing.T) {
	var a Args

	a.Set("foo", "bar")
	a.Set("bar", "baz")
	a.Set("baz", "aaaaaa")

	a.Del("bar")
	s := string(a.QueryString())
	if s != "foo=bar&baz=aaaaaa" {
		t.Fatalf("unexpected result %q. Expecting %q", s, "foo=bar&baz=aaaaaa")
	}

	a.Del("foo")
	s = string(a.QueryString())
	if s != "baz=aaaaaa" {
		t.Fatalf("unexpected result %q. Expecting %q", s, "baz=aaaaaa")
	}

	a.Del("baz")
	s = string(a.QueryString())
	if s != "" {
		t.Fatalf("unexpected result %q. Expecting %q", s, "")
	}
}

func TestArgsStringComposeDeleteNoValue(t *testing.T) {
	var a Args

	a.Set("foo", "bar")
	a.SetNoValue("bar")
	a.Set("baz", "aaaaaa")

	a.Del("bar")
	s := string(a.QueryString())
	if s != "foo=bar&baz=aaaaaa" {
		t.Fatalf("unexpected result %q. Expecting %q", s, "foo=bar&baz=aaaaaa")
	}

	a.Del("foo")
	s = string(a.QueryString())
	if s != "baz=aaaaaa" {
		t.Fatalf("unexpected result %q. Expecting %q", s, "baz=aaaaaa")
	}

	a.Del("baz")
	s = string(a.QueryString())
	if s != "" {
		t.Fatalf("unexpected result %q. Expecting %q", s, "")
	}
}

func TestArgsStringComposeSet(t *testing.T) {
	var a Args

	a.Set("foo", "bar")
	a.Set("bar", "baz")
	a.Set("baz", "aaaaaa")

	a.Set("bar", "x")
	s := string(a.QueryString())
	if s != "foo=bar&bar=x&baz=aaaaaa" {
		t.Fatalf("unexpected result %q. Expecting %q", s, "foo=bar&bar=x&baz=aaaaaa")
	}

	a.Set("foo", "y")
	s = string(a.QueryString())
	if s != "foo=y&bar=x&baz=aaaaaa" {
		t.Fatalf("unexpected result %q. Expecting %q", s, "foo=y&bar=x&baz=aaaaaa")
	}

	a.Set("baz", "z")
	s = string(a.QueryString())
	if s != "foo=y&bar=x&baz=z" {
		t.Fatalf("unexpected result %q. Expecting %q", s, "foo=y&bar=x&baz=z")
	}

	a.Set("aaa", "bbb")
	s = string(a.QueryString())
	if s != "foo=y&bar=x&baz=z&aaa=bbb" {
		t.Fatalf("unexpected result %q. Expecting %q", s, "foo=y&bar=x&baz=z&aaa=bbb")
	}
}

func TestArgsStringComposeSetNoValue(t *testing.T) {
	var a Args

	a.Set("foo", "bar")
	a.SetNoValue("bar")
	a.Set("baz", "aaaaaa")

	a.SetNoValue("bar")
	s := string(a.QueryString())
	if s != "foo=bar&bar&baz=aaaaaa" {
		t.Fatalf("unexpected result %q. Expecting %q", s, "foo=bar&bar&baz=aaaaaa")
	}

	a.Set("foo", "y")
	s = string(a.QueryString())
	if s != "foo=y&bar&baz=aaaaaa" {
		t.Fatalf("unexpected result %q. Expecting %q", s, "foo=y&bar&baz=aaaaaa")
	}

	a.Set("baz", "z")
	s = string(a.QueryString())
	if s != "foo=y&bar&baz=z" {
		t.Fatalf("unexpected result %q. Expecting %q", s, "foo=y&bar&baz=z")
	}

	a.SetNoValue("aaa")
	s = string(a.QueryString())
	if s != "foo=y&bar&baz=z&aaa" {
		t.Fatalf("unexpected result %q. Expecting %q", s, "foo=y&bar&baz=z&aaa")
	}
}

func TestArgsStringComposeAdd(t *testing.T) {
	var a Args

	a.Add("foo", "bar")
	a.Add("bar", "baz")
	a.Add("baz", "aaaaaa")

	a.Add("bar", "x")
	s := string(a.QueryString())
	if s != "foo=bar&bar=baz&baz=aaaaaa&bar=x" {
		t.Fatalf("unexpected result %q. Expecting %q", s, "foo=bar&bar=baz&baz=aaaaaa&bar=x")
	}

	a.Add("foo", "y")
	s = string(a.QueryString())
	if s != "foo=bar&bar=baz&baz=aaaaaa&bar=x&foo=y" {
		t.Fatalf("unexpected result %q. Expecting %q", s, "foo=bar&bar=baz&baz=aaaaaa&bar=x&foo=y")
	}

	a.Add("baz", "z")
	s = string(a.QueryString())
	if s != "foo=bar&bar=baz&baz=aaaaaa&bar=x&foo=y&baz=z" {
		t.Fatalf("unexpected result %q. Expecting %q", s, "foo=bar&bar=baz&baz=aaaaaa&bar=x&foo=y&baz=z")
	}

	a.Add("aaa", "bbb")
	s = string(a.QueryString())
	if s != "foo=bar&bar=baz&baz=aaaaaa&bar=x&foo=y&baz=z&aaa=bbb" {
		t.Fatalf("unexpected result %q. Expecting %q", s, "foo=bar&bar=baz&baz=aaaaaa&bar=x&foo=y&baz=z&aaa=bbb")
	}
}

func TestArgsStringComposeAddNoValue(t *testing.T) {
	var a Args

	a.Add("foo", "bar")
	a.AddNoValue("bar")
	a.Add("baz", "aaaaaa")

	a.AddNoValue("bar")
	s := string(a.QueryString())
	if s != "foo=bar&bar&baz=aaaaaa&bar" {
		t.Fatalf("unexpected result %q. Expecting %q", s, "foo=bar&bar&baz=aaaaaa&bar")
	}

	a.Add("foo", "y")
	s = string(a.QueryString())
	if s != "foo=bar&bar&baz=aaaaaa&bar&foo=y" {
		t.Fatalf("unexpected result %q. Expecting %q", s, "foo=bar&bar&baz=aaaaaa&bar&foo=y")
	}

	a.Add("baz", "z")
	s = string(a.QueryString())
	if s != "foo=bar&bar&baz=aaaaaa&bar&foo=y&baz=z" {
		t.Fatalf("unexpected result %q. Expecting %q", s, "foo=bar&bar&baz=aaaaaa&bar&foo=y&baz=z")
	}

	a.AddNoValue("aaa")
	s = string(a.QueryString())
	if s != "foo=bar&bar&baz=aaaaaa&bar&foo=y&baz=z&aaa" {
		t.Fatalf("unexpected result %q. Expecting %q", s, "foo=bar&bar&baz=aaaaaa&bar&foo=y&baz=z&aaa")
	}
}

func TestArgsGetAll(t *testing.T) {
	var a Args

	a.Add("foo", "bar")
	a.Add("bar", "baz")
	a.Add("baz", "aaaaaa")

	if v := string(a.Peek("foo")); v != "bar" {
		t.Fatalf("unexpected result %q. Expecting %q", v, "bar")
	}
	if v := string(a.Peek("bar")); v != "baz" {
		t.Fatalf("unexpected result %q. Expecting %q", v, "baz")
	}
	if v := string(a.Peek("baz")); v != "aaaaaa" {
		t.Fatalf("unexpected result %q. Expecting %q", v, "aaaaaa")
	}
	if v := string(a.Peek("xxxx")); v != "" {
		t.Fatalf("unexpected result %q. Expecting %q", v, "")
	}
}

func TestArgsGetAllNoValue(t *testing.T) {
	var a Args

	a.Add("foo", "bar")
	a.AddNoValue("bar")
	a.Add("baz", "aaaaaa")

	if v := string(a.Peek("foo")); v != "bar" {
		t.Fatalf("unexpected result %q. Expecting %q", v, "bar")
	}
	if v := string(a.Peek("bar")); v != "" {
		t.Fatalf("unexpected result %q. Expecting %q", v, "")
	}
	if v := string(a.Peek("baz")); v != "aaaaaa" {
		t.Fatalf("unexpected result %q. Expecting %q", v, "aaaaaa")
	}
	if v := string(a.Peek("xxxx")); v != "" {
		t.Fatalf("unexpected result %q. Expecting %q", v, "")
	}
}

func TestArgsGetAll2(t *testing.T) {
	var a Args

	a.Add("foo", "bar")
	a.Add("bar", "baz")
	a.Add("baz", "aaaaaa")

	if v := string(a.Peek("foo")); v != "bar" {
		t.Fatalf("unexpected result %q. Expecting %q", v, "bar")
	}
	if v := string(a.Peek("bar")); v != "baz" {
		t.Fatalf("unexpected result %q. Expecting %q", v, "baz")
	}
	if v := string(a.Peek("baz")); v != "aaaaaa" {
		t.Fatalf("unexpected result %q. Expecting %q", v, "aaaaaa")
	}
	if v := string(a.Peek("xxxx")); v != "" {
		t.Fatalf("unexpected result %q. Expecting %q", v, "")
	}

	a.Del("foo")
	if v := string(a.Peek("foo")); v != "" {
		t.Fatalf("unexpected result %q. Expecting %q", v, "")
	}
	if v := string(a.Peek("bar")); v != "baz" {
		t.Fatalf("unexpected result %q. Expecting %q", v, "baz")
	}
	if v := string(a.Peek("baz")); v != "aaaaaa" {
		t.Fatalf("unexpected result %q. Expecting %q", v, "aaaaaa")
	}
	if v := string(a.Peek("xxxx")); v != "" {
		t.Fatalf("unexpected result %q. Expecting %q", v, "")
	}
}

func TestArgsGetAll2NoValue(t *testing.T) {
	var a Args

	a.Add("foo", "bar")
	a.AddNoValue("bar")
	a.Add("baz", "aaaaaa")

	if v := string(a.Peek("foo")); v != "bar" {
		t.Fatalf("unexpected result %q. Expecting %q", v, "bar")
	}
	if v := string(a.Peek("bar")); v != "" {
		t.Fatalf("unexpected result %q. Expecting %q", v, "")
	}
	if v := string(a.Peek("baz")); v != "aaaaaa" {
		t.Fatalf("unexpected result %q. Expecting %q", v, "aaaaaa")
	}
	if v := string(a.Peek("xxxx")); v != "" {
		t.Fatalf("unexpected result %q. Expecting %q", v, "")
	}

	a.Del("foo")
	if v := string(a.Peek("foo")); v != "" {
		t.Fatalf("unexpected result %q. Expecting %q", v, "")
	}
	if v := string(a.Peek("bar")); v != "" {
		t.Fatalf("unexpected result %q. Expecting %q", v, "")
	}
	if v := string(a.Peek("baz")); v != "aaaaaa" {
		t.Fatalf("unexpected result %q. Expecting %q", v, "aaaaaa")
	}
	if v := string(a.Peek("xxxx")); v != "" {
		t.Fatalf("unexpected result %q. Expecting %q", v, "")
	}
}

func TestArgsGetAll3(t *testing.T) {
	var a Args

	a.Add("foo", "bar")
	a.Add("bar", "baz")
	a.Add("baz", "aaaaaa")

	if v := string(a.Peek("foo")); v != "bar" {
		t.Fatalf("unexpected result %q. Expecting %q", v, "bar")
	}
	if v := string(a.Peek("bar")); v != "baz" {
		t.Fatalf("unexpected result %q. Expecting %q", v, "baz")
	}
	if v := string(a.Peek("baz")); v != "aaaaaa" {
		t.Fatalf("unexpected result %q. Expecting %q", v, "aaaaaa")
	}
	if v := string(a.Peek("xxxx")); v != "" {
		t.Fatalf("unexpected result %q. Expecting %q", v, "")
	}

	a.Del("foo")
	a.Del("bar")
	if v := string(a.Peek("foo")); v != "" {
		t.Fatalf("unexpected result %q. Expecting %q", v, "")
	}
	if v := string(a.Peek("bar")); v != "" {
		t.Fatalf("unexpected result %q. Expecting %q", v, "")
	}
	if v := string(a.Peek("baz")); v != "aaaaaa" {
		t.Fatalf("unexpected result %q. Expecting %q", v, "aaaaaa")
	}
	if v := string(a.Peek("xxxx")); v != "" {
		t.Fatalf("unexpected result %q. Expecting %q", v, "")
	}
}

func TestArgsGetAll3NoValue(t *testing.T) {
	var a Args

	a.Add("foo", "bar")
	a.AddNoValue("bar")
	a.Add("baz", "aaaaaa")

	if v := string(a.Peek("foo")); v != "bar" {
		t.Fatalf("unexpected result %q. Expecting %q", v, "bar")
	}
	if v := string(a.Peek("bar")); v != "" {
		t.Fatalf("unexpected result %q. Expecting %q", v, "")
	}
	if v := string(a.Peek("baz")); v != "aaaaaa" {
		t.Fatalf("unexpected result %q. Expecting %q", v, "aaaaaa")
	}
	if v := string(a.Peek("xxxx")); v != "" {
		t.Fatalf("unexpected result %q. Expecting %q", v, "")
	}

	a.Del("foo")
	a.Del("bar")
	if v := string(a.Peek("foo")); v != "" {
		t.Fatalf("unexpected result %q. Expecting %q", v, "")
	}
	if v := string(a.Peek("bar")); v != "" {
		t.Fatalf("unexpected result %q. Expecting %q", v, "")
	}
	if v := string(a.Peek("baz")); v != "aaaaaa" {
		t.Fatalf("unexpected result %q. Expecting %q", v, "aaaaaa")
	}
	if v := string(a.Peek("xxxx")); v != "" {
		t.Fatalf("unexpected result %q. Expecting %q", v, "")
	}
}

func TestArgsGetAll4(t *testing.T) {
	var a Args

	a.Add("foo", "bar")
	a.Add("bar", "baz")
	a.Add("baz", "aaaaaa")

	if v := string(a.Peek("foo")); v != "bar" {
		t.Fatalf("unexpected result %q. Expecting %q", v, "bar")
	}
	if v := string(a.Peek("bar")); v != "baz" {
		t.Fatalf("unexpected result %q. Expecting %q", v, "baz")
	}
	if v := string(a.Peek("baz")); v != "aaaaaa" {
		t.Fatalf("unexpected result %q. Expecting %q", v, "aaaaaa")
	}
	if v := string(a.Peek("xxxx")); v != "" {
		t.Fatalf("unexpected result %q. Expecting %q", v, "")
	}

	a.Del("foo")
	a.Del("bar")
	a.Del("baz")
	if v := string(a.Peek("foo")); v != "" {
		t.Fatalf("unexpected result %q. Expecting %q", v, "")
	}
	if v := string(a.Peek("bar")); v != "" {
		t.Fatalf("unexpected result %q. Expecting %q", v, "")
	}
	if v := string(a.Peek("baz")); v != "" {
		t.Fatalf("unexpected result %q. Expecting %q", v, "")
	}
	if v := string(a.Peek("xxxx")); v != "" {
		t.Fatalf("unexpected result %q. Expecting %q", v, "")
	}
}

func TestArgsGetAll4NoValue(t *testing.T) {
	var a Args

	a.Add("foo", "bar")
	a.AddNoValue("bar")
	a.Add("baz", "aaaaaa")

	if v := string(a.Peek("foo")); v != "bar" {
		t.Fatalf("unexpected result %q. Expecting %q", v, "bar")
	}
	if v := string(a.Peek("bar")); v != "" {
		t.Fatalf("unexpected result %q. Expecting %q", v, "")
	}
	if v := string(a.Peek("baz")); v != "aaaaaa" {
		t.Fatalf("unexpected result %q. Expecting %q", v, "aaaaaa")
	}
	if v := string(a.Peek("xxxx")); v != "" {
		t.Fatalf("unexpected result %q. Expecting %q", v, "")
	}

	a.Del("foo")
	a.Del("bar")
	a.Del("baz")
	if v := string(a.Peek("foo")); v != "" {
		t.Fatalf("unexpected result %q. Expecting %q", v, "")
	}
	if v := string(a.Peek("bar")); v != "" {
		t.Fatalf("unexpected result %q. Expecting %q", v, "")
	}
	if v := string(a.Peek("baz")); v != "" {
		t.Fatalf("unexpected result %q. Expecting %q", v, "")
	}
	if v := string(a.Peek("xxxx")); v != "" {
		t.Fatalf("unexpected result %q. Expecting %q", v, "")
	}
}

func TestArgsGetAll5(t *testing.T) {
	var a Args

	a.Add("foo", "bar")
	a.Add("bar", "baz")
	a.Add("baz", "aaaaaa")

	if v := string(a.Peek("foo")); v != "bar" {
		t.Fatalf("unexpected result %q. Expecting %q", v, "bar")
	}
	if v := string(a.Peek("bar")); v != "baz" {
		t.Fatalf("unexpected result %q. Expecting %q", v, "baz")
	}
	if v := string(a.Peek("baz")); v != "aaaaaa" {
		t.Fatalf("unexpected result %q. Expecting %q", v, "aaaaaa")
	}
	if v := string(a.Peek("xxxx")); v != "" {
		t.Fatalf("unexpected result %q. Expecting %q", v, "")
	}

	a.Del("foo")
	a.Del("bar")
	a.Del("baz")
	a.Del("xxxx")
	if v := string(a.Peek("foo")); v != "" {
		t.Fatalf("unexpected result %q. Expecting %q", v, "")
	}
	if v := string(a.Peek("bar")); v != "" {
		t.Fatalf("unexpected result %q. Expecting %q", v, "")
	}
	if v := string(a.Peek("baz")); v != "" {
		t.Fatalf("unexpected result %q. Expecting %q", v, "")
	}
	if v := string(a.Peek("xxxx")); v != "" {
		t.Fatalf("unexpected result %q. Expecting %q", v, "")
	}
}

func TestArgsGetAll5NoValue(t *testing.T) {
	var a Args

	a.Add("foo", "bar")
	a.AddNoValue("bar")
	a.Add("baz", "aaaaaa")

	if v := string(a.Peek("foo")); v != "bar" {
		t.Fatalf("unexpected result %q. Expecting %q", v, "bar")
	}
	if v := string(a.Peek("bar")); v != "" {
		t.Fatalf("unexpected result %q. Expecting %q", v, "")
	}
	if v := string(a.Peek("baz")); v != "aaaaaa" {
		t.Fatalf("unexpected result