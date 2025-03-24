package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
)

func TestCharset(t *testing.T) {
	r := chi.NewRouter()

	r.Use(ContentCharset("utf-8"))

	r.Post("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("foo"))
	})

	ts := httptest.NewServer(r)
	defer ts.Close()

	if _, err := http.Post(ts.URL, "text/plain; charset=utf-8", nil); err != nil {
		t.Fatalf("expect request to succeed, but got err=%v", err)
	}

	r1, err := http.Post(ts.URL, "text/plain; charset=latin1", nil)
	if err != nil {
		t.Fatalf("expect request to succeed, but got err=%v", err)
	}
	if r1.StatusCode != http.StatusUnsupportedMediaType {
		t.Fatalf("expect request to fail with 415, but got %d", r1.StatusCode)
	}
}