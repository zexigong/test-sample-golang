package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRealIP(t *testing.T) {
	t.Parallel()

	r := http.NewServeMux()
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(r.RemoteAddr))
	})

	ts := httptest.NewServer(RealIP(r))
	defer ts.Close()

	req, _ := http.NewRequest("GET", ts.URL, nil)
	req.Header.Set("X-Forwarded-For", "8.8.8.8")

	res, err := ts.Client().Do(req)
	if err != nil {
		t.Fatal(err)
	}

	buf := make([]byte, 20)
	n, _ := res.Body.Read(buf)
	res.Body.Close()

	if got, want := string(buf[:n]), "8.8.8.8"; got != want {
		t.Errorf("got %s want %s", got, want)
	}
}

func TestRealIPHeader(t *testing.T) {
	t.Parallel()

	r := http.NewServeMux()
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(r.RemoteAddr))
	})

	ts := httptest.NewServer(RealIP(r))
	defer ts.Close()

	req, _ := http.NewRequest("GET", ts.URL, nil)
	req.Header.Set("X-Real-IP", "8.8.8.8")

	res, err := ts.Client().Do(req)
	if err != nil {
		t.Fatal(err)
	}

	buf := make([]byte, 20)
	n, _ := res.Body.Read(buf)
	res.Body.Close()

	if got, want := string(buf[:n]), "8.8.8.8"; got != want {
		t.Errorf("got %s want %s", got, want)
	}
}

func TestRealIPHeaderPriority(t *testing.T) {
	t.Parallel()

	r := http.NewServeMux()
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(r.RemoteAddr))
	})

	ts := httptest.NewServer(RealIP(r))
	defer ts.Close()

	req, _ := http.NewRequest("GET", ts.URL, nil)
	req.Header.Set("X-Forwarded-For", "8.8.8.8")
	req.Header.Set("X-Real-IP", "9.9.9.9")

	res, err := ts.Client().Do(req)
	if err != nil {
		t.Fatal(err)
	}

	buf := make([]byte, 20)
	n, _ := res.Body.Read(buf)
	res.Body.Close()

	if got, want := string(buf[:n]), "9.9.9.9"; got != want {
		t.Errorf("got %s want %s", got, want)
	}
}

func TestRealIPHeaderTrueClientIPPriority(t *testing.T) {
	t.Parallel()

	r := http.NewServeMux()
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(r.RemoteAddr))
	})

	ts := httptest.NewServer(RealIP(r))
	defer ts.Close()

	req, _ := http.NewRequest("GET", ts.URL, nil)
	req.Header.Set("X-Forwarded-For", "8.8.8.8")
	req.Header.Set("X-Real-IP", "9.9.9.9")
	req.Header.Set("True-Client-IP", "10.10.10.10")

	res, err := ts.Client().Do(req)
	if err != nil {
		t.Fatal(err)
	}

	buf := make([]byte, 20)
	n, _ := res.Body.Read(buf)
	res.Body.Close()

	if got, want := string(buf[:n]), "10.10.10.10"; got != want {
		t.Errorf("got %s want %s", got, want)
	}
}

func TestRealIPHeaderInvalid(t *testing.T) {
	t.Parallel()

	r := http.NewServeMux()
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(r.RemoteAddr))
	})

	ts := httptest.NewServer(RealIP(r))
	defer ts.Close()

	req, _ := http.NewRequest("GET", ts.URL, nil)
	req.Header.Set("X-Real-IP", "unparseable")

	res, err := ts.Client().Do(req)
	if err != nil {
		t.Fatal(err)
	}

	buf := make([]byte, 20)
	n, _ := res.Body.Read(buf)
	res.Body.Close()

	if got, want := string(buf[:n]), "127.0.0.1"; got != want {
		t.Errorf("got %s want %s", got, want)
	}
}

func TestRealIPHeaderEmpty(t *testing.T) {
	t.Parallel()

	r := http.NewServeMux()
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(r.RemoteAddr))
	})

	ts := httptest.NewServer(RealIP(r))
	defer ts.Close()

	req, _ := http.NewRequest("GET", ts.URL, nil)

	res, err := ts.Client().Do(req)
	if err != nil {
		t.Fatal(err)
	}

	buf := make([]byte, 20)
	n, _ := res.Body.Read(buf)
	res.Body.Close()

	if got, want := string(buf[:n]), "127.0.0.1"; got != want {
		t.Errorf("got %s want %s", got, want)
	}
}

func TestRealIPHeaderWithPort(t *testing.T) {
	t.Parallel()

	r := http.NewServeMux()
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(r.RemoteAddr))
	})

	ts := httptest.NewServer(RealIP(r))
	defer ts.Close()

	req, _ := http.NewRequest("GET", ts.URL, nil)
	req.Header.Set("X-Forwarded-For", "8.8.8.8:1234")

	res, err := ts.Client().Do(req)
	if err != nil {
		t.Fatal(err)
	}

	buf := make([]byte, 20)
	n, _ := res.Body.Read(buf)
	res.Body.Close()

	if got, want := string(buf[:n]), "8.8.8.8"; got != want {
		t.Errorf("got %s want %s", got, want)
	}
}

func TestRealIPHeaderWithPortInvalid(t *testing.T) {
	t.Parallel()

	r := http.NewServeMux()
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(r.RemoteAddr))
	})

	ts := httptest.NewServer(RealIP(r))
	defer ts.Close()

	req, _ := http.NewRequest("GET", ts.URL, nil)
	req.Header.Set("X-Real-IP", "unparseable:1234")

	res, err := ts.Client().Do(req)
	if err != nil {
		t.Fatal(err)
	}

	buf := make([]byte, 20)
	n, _ := res.Body.Read(buf)
	res.Body.Close()

	if got, want := string(buf[:n]), "127.0.0.1"; got != want {
		t.Errorf("got %s want %s", got, want)
	}
}