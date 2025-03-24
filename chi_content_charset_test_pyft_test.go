package middleware

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestContentCharset(t *testing.T) {
	h := ContentCharset("utf-8", "latin-1", "")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}))
	r := httptest.NewRequest("POST", "/", strings.NewReader("foo"))
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	if w.Result().StatusCode != 200 {
		t.Fatalf("expected 200, got %d", w.Result().StatusCode)
	}

	for _, ct := range []string{"application/json", "text/plain"} {
		r := httptest.NewRequest("POST", "/", strings.NewReader("foo"))
		r.Header.Set("Content-Type", ct)
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)
		if w.Result().StatusCode != 200 {
			t.Fatalf("expected 200, got %d", w.Result().StatusCode)
		}
	}

	for _, ct := range []string{
		"application/json; charset=latin-1",
		"text/plain; charset=utf-8",
		"application/json; charset=latin-1; other=value",
		"application/json; charset=utf-8",
	} {
		r := httptest.NewRequest("POST", "/", strings.NewReader("foo"))
		r.Header.Set("Content-Type", ct)
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)
		if w.Result().StatusCode != 200 {
			t.Fatalf("expected 200, got %d", w.Result().StatusCode)
		}
	}

	for _, ct := range []string{
		"application/json; charset=latin-2",
		"text/plain; charset=latin-3",
		"application/json; charset=utf-16",
	} {
		r := httptest.NewRequest("POST", "/", strings.NewReader("foo"))
		r.Header.Set("Content-Type", ct)
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)
		if w.Result().StatusCode != http.StatusUnsupportedMediaType {
			t.Fatalf("expected %d, got %d", http.StatusUnsupportedMediaType, w.Result().StatusCode)
		}
	}

	h = ContentCharset("utf-8", "latin-1", "")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}))
	r = httptest.NewRequest("POST", "/", strings.NewReader("foo"))
	w = httptest.NewRecorder()
	h.ServeHTTP(w, r)
	if w.Result().StatusCode != 200 {
		t.Fatalf("expected 200, got %d", w.Result().StatusCode)
	}
}

func TestCharset(t *testing.T) {
	h := ContentCharset("utf-8", "latin-1")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}))
	r := httptest.NewRequest("POST", "/", strings.NewReader("foo"))
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	if w.Result().StatusCode != 415 {
		t.Fatalf("expected 415, got %d", w.Result().StatusCode)
	}
}

func TestContentCharsetEmpty(t *testing.T) {
	h := ContentCharset()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}))
	r := httptest.NewRequest("POST", "/", strings.NewReader("foo"))
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	if w.Result().StatusCode != 200 {
		t.Fatalf("expected 200, got %d", w.Result().StatusCode)
	}

	for _, ct := range []string{
		"application/json; charset=latin-2",
		"text/plain; charset=latin-3",
		"application/json; charset=utf-16",
		"application/json; charset=latin-1",
		"text/plain; charset=utf-8",
		"application/json; charset=latin-1; other=value",
		"application/json; charset=utf-8",
	} {
		r := httptest.NewRequest("POST", "/", strings.NewReader("foo"))
		r.Header.Set("Content-Type", ct)
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)
		if w.Result().StatusCode != 200 {
			t.Fatalf("expected 200, got %d", w.Result().StatusCode)
		}
	}
}

func TestCharsetEmpty(t *testing.T) {
	h := ContentCharset()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}))
	r := httptest.NewRequest("POST", "/", strings.NewReader("foo"))
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	if w.Result().StatusCode != 200 {
		t.Fatalf("expected 200, got %d", w.Result().StatusCode)
	}
}

func TestContentCharsetWithNext(t *testing.T) {
	r := httptest.NewRequest("POST", "/", strings.NewReader("foo"))
	r.Header.Set("Content-Type", "application/json; charset=latin-1")
	w := httptest.NewRecorder()

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("next"))
	})

	h := ContentCharset("utf-8", "latin-1", "")(next)
	h.ServeHTTP(w, r)

	resp := w.Result()
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	if string(body) != "next" {
		t.Fatalf("expected body to be 'next', got %s", string(body))
	}
}