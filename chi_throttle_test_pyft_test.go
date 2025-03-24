package middleware

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func TestThrottle(t *testing.T) {
	t.Parallel()

	c := chi.NewRouter()
	c.Use(Throttle(2))
	c.Get("/", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Millisecond * 10)
		w.Write([]byte("."))
	})

	ts := httptest.NewServer(c)
	defer ts.Close()

	type want struct {
		respCode    int
		respBody    string
		contentType string
	}

	tests := []struct {
		name string
		want want
	}{
		{
			"allowed",
			want{
				respCode:    http.StatusOK,
				respBody:    ".",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			"allowed",
			want{
				respCode:    http.StatusOK,
				respBody:    ".",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			"rejected",
			want{
				respCode:    http.StatusTooManyRequests,
				respBody:    "Server capacity exceeded.\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			"rejected",
			want{
				respCode:    http.StatusTooManyRequests,
				respBody:    "Server capacity exceeded.\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			"rejected",
			want{
				respCode:    http.StatusTooManyRequests,
				respBody:    "Server capacity exceeded.\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			"allowed",
			want{
				respCode:    http.StatusOK,
				respBody:    ".",
				contentType: "text/plain; charset=utf-8",
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			resp, err := ts.Client().Get(ts.URL)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if resp.StatusCode != tt.want.respCode {
				t.Errorf("got %d, want %d", resp.StatusCode, tt.want.respCode)
			}

			if ct := resp.Header.Get("Content-Type"); ct != tt.want.contentType {
				t.Errorf("got %q, want %q", ct, tt.want.contentType)
			}

			got, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if string(got) != tt.want.respBody {
				t.Errorf("got %q, want %q", got, tt.want.respBody)
			}
		})
	}
}

func TestThrottleBacklog(t *testing.T) {
	t.Parallel()

	c := chi.NewRouter()
	c.Use(ThrottleBacklog(1, 2, time.Second))
	c.Get("/", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Millisecond * 10)
		w.Write([]byte("."))
	})

	ts := httptest.NewServer(c)
	defer ts.Close()

	type want struct {
		respCode    int
		respBody    string
		contentType string
	}

	tests := []struct {
		name string
		want want
	}{
		{
			"allowed",
			want{
				respCode:    http.StatusOK,
				respBody:    ".",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			"allowed",
			want{
				respCode:    http.StatusOK,
				respBody:    ".",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			"allowed",
			want{
				respCode:    http.StatusOK,
				respBody:    ".",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			"rejected",
			want{
				respCode:    http.StatusTooManyRequests,
				respBody:    "Server capacity exceeded.\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			"rejected",
			want{
				respCode:    http.StatusTooManyRequests,
				respBody:    "Server capacity exceeded.\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			"rejected",
			want{
				respCode:    http.StatusTooManyRequests,
				respBody:    "Server capacity exceeded.\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			"allowed",
			want{
				respCode:    http.StatusOK,
				respBody:    ".",
				contentType: "text/plain; charset=utf-8",
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			resp, err := ts.Client().Get(ts.URL)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if resp.StatusCode != tt.want.respCode {
				t.Errorf("got %d, want %d", resp.StatusCode, tt.want.respCode)
			}

			if ct := resp.Header.Get("Content-Type"); ct != tt.want.contentType {
				t.Errorf("got %q, want %q", ct, tt.want.contentType)
			}

			got, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if string(got) != tt.want.respBody {
				t.Errorf("got %q, want %q", got, tt.want.respBody)
			}
		})
	}
}

func TestThrottleBacklogTimeout(t *testing.T) {
	t.Parallel()

	c := chi.NewRouter()
	c.Use(ThrottleBacklog(1, 1, time.Millisecond*10))
	c.Get("/", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Millisecond * 30)
		w.Write([]byte("."))
	})

	ts := httptest.NewServer(c)
	defer ts.Close()

	type want struct {
		respCode    int
		respBody    string
		contentType string
	}

	tests := []struct {
		name string
		want want
	}{
		{
			"allowed",
			want{
				respCode:    http.StatusOK,
				respBody:    ".",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			"rejected",
			want{
				respCode:    http.StatusTooManyRequests,
				respBody:    "Timed out while waiting for a pending request to complete.\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			"rejected",
			want{
				respCode:    http.StatusTooManyRequests,
				respBody:    "Server capacity exceeded.\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			"rejected",
			want{
				respCode:    http.StatusTooManyRequests,
				respBody:    "Server capacity exceeded.\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			"rejected",
			want{
				respCode:    http.StatusTooManyRequests,
				respBody:    "Server capacity exceeded.\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			"allowed",
			want{
				respCode:    http.StatusOK,
				respBody:    ".",
				contentType: "text/plain; charset=utf-8",
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			resp, err := ts.Client().Get(ts.URL)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if resp.StatusCode != tt.want.respCode {
				t.Errorf("got %d, want %d", resp.StatusCode, tt.want.respCode)
			}

			if ct := resp.Header.Get("Content-Type"); ct != tt.want.contentType {
				t.Errorf("got %q, want %q", ct, tt.want.contentType)
			}

			got, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if string(got) != tt.want.respBody {
				t.Errorf("got %q, want %q", got, tt.want.respBody)
			}
		})
	}
}

func TestThrottleWithOpts(t *testing.T) {
	t.Parallel()

	c := chi.NewRouter()
	c.Use(ThrottleWithOpts(ThrottleOpts{
		Limit:          1,
		BacklogLimit:   1,
		BacklogTimeout: time.Millisecond * 10,
		RetryAfterFn: func(ctxDone bool) time.Duration {
			if ctxDone {
				return 1
			}
			return 2
		},
	}))
	c.Get("/", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Millisecond * 30)
		w.Write([]byte("."))
	})

	ts := httptest.NewServer(c)
	defer ts.Close()

	type want struct {
		respCode    int
		respBody    string
		retryAfter  string
		contentType string
	}

	tests := []struct {
		name string
		want want
	}{
		{
			"allowed",
			want{
				respCode:    http.StatusOK,
				respBody:    ".",
				retryAfter:  "",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			"rejected",
			want{
				respCode:    http.StatusTooManyRequests,
				respBody:    "Timed out while waiting for a pending request to complete.\n",
				retryAfter:  "1",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			"rejected",
			want{
				respCode:    http.StatusTooManyRequests,
				respBody:    "Server capacity exceeded.\n",
				retryAfter:  "2",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			"rejected",
			want{
				respCode:    http.StatusTooManyRequests,
				respBody:    "Server capacity exceeded.\n",
				retryAfter:  "2",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			"rejected",
			want{
				respCode:    http.StatusTooManyRequests,
				respBody:    "Server capacity exceeded.\n",
				retryAfter:  "2",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			"allowed",
			want{
				respCode:    http.StatusOK,
				respBody:    ".",
				retryAfter:  "",
				contentType: "text/plain; charset=utf-8",
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			resp, err := ts.Client().Get(ts.URL)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if resp.StatusCode != tt.want.respCode {
				t.Errorf("got %d, want %d", resp.StatusCode, tt.want.respCode)
			}

			if ct := resp.Header.Get("Content-Type"); ct != tt.want.contentType {
				t.Errorf("got %q, want %q", ct, tt.want.contentType)
			}

			got, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if string(got) != tt.want.respBody {
				t.Errorf("got %q, want %q", got, tt.want.respBody)
			}

			if ra := resp.Header.Get("Retry-After"); ra != tt.want.retryAfter {
				t.Errorf("got %q, want %q", ra, tt.want.retryAfter)
			}
		})
	}
}

func TestThrottleWithOptsPanic(t *testing.T) {
	t.Parallel()

	opts := []ThrottleOpts{
		{Limit: -1},
		{Limit: 0},
		{Limit: 1, BacklogLimit: -1},
	}
	for _, opt := range opts {
		opt := opt
		t.Run(fmt.Sprintf("%+v", opt), func(t *testing.T) {
			t.Parallel()
			assertPanic(t, func() { ThrottleWithOpts(opt) })
		})
	}
}

func assertPanic(t *testing.T, fn func()) {
	t.Helper()

	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected to panic")
		}
	}()
	fn()
}

func TestThrottleWithOptsBacklogTimeoutWithContext(t *testing.T) {
	t.Parallel()

	c := chi.NewRouter()
	c.Use(ThrottleWithOpts(ThrottleOpts{
		Limit:          1,
		BacklogLimit:   1,
		BacklogTimeout: time.Millisecond * 10,
		StatusCode:     http.StatusInternalServerError,
	}))
	c.Get("/", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Millisecond * 30)
		w.Write([]byte("."))
	})

	ts := httptest.NewServer(c)
	defer ts.Close()

	type want struct {
		respCode    int
		respBody    string
		contentType string
	}

	tests := []struct {
		name string
		want want
	}{
		{
			"allowed",
			want{
				respCode:    http.StatusOK,
				respBody:    ".",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			"rejected",
			want{
				respCode:    http.StatusInternalServerError,
				respBody:    "Timed out while waiting for a pending request to complete.\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			"rejected",
			want{
				respCode:    http.StatusInternalServerError,
				respBody:    "Server capacity exceeded.\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			"rejected",
			want{
				respCode:    http.StatusInternalServerError,
				respBody:    "Server capacity exceeded.\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			"rejected",
			want{
				respCode:    http.StatusInternalServerError,
				respBody:    "Server capacity exceeded.\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			"allowed",
			want{
				respCode:    http.StatusOK,
				respBody:    ".",
				contentType: "text/plain; charset=utf-8",
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			resp, err := ts.Client().Get(ts.URL)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if resp.StatusCode != tt.want.respCode {
				t.Errorf("got %d, want %d", resp.StatusCode, tt.want.respCode)
			}

			if ct := resp.Header.Get("Content-Type"); ct != tt.want.contentType {
				t.Errorf("got %q, want %q", ct, tt.want.contentType)
			}

			got, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if string(got) != tt.want.respBody {
				t.Errorf("got %q, want %q", got, tt.want.respBody)
			}
		})
	}
}

func Test_throttler_setRetryAfterHeaderIfNeeded(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		retryAfterFn func(ctxDone bool) time.Duration
		ctxDone     bool
		want        string
	}{
		{
			name:        "nil retryAfterFn",
			retryAfterFn: nil,
			ctxDone:     false,
			want:        "",
		},
		{
			name:        "retryAfterFn returns 1 second",
			retryAfterFn: func(ctxDone bool) time.Duration { return 1 * time.Second },
			ctxDone:     false,
			want:        "1",
		},
		{
			name:        "retryAfterFn returns 2 seconds",
			retryAfterFn: func(ctxDone bool) time.Duration { return 2 * time.Second },
			ctxDone:     true,
			want:        "2",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			rr := httptest.NewRecorder()

			tr := throttler{retryAfterFn: tt.retryAfterFn}
			tr.setRetryAfterHeaderIfNeeded(rr, tt.ctxDone)

			if got := rr.Header().Get("Retry-After"); got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestThrottleBacklogTimeoutWithContext(t *testing.T) {
	t.Parallel()

	c := chi.NewRouter()
	c.Use(middleware.Timeout(time.Millisecond * 50))
	c.Use(ThrottleBacklog(1, 1, time.Second))
	c.Get("/", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Millisecond * 100)
		w.Write([]byte("."))
	})

	ts := httptest.NewServer(c)
	defer ts.Close()

	type want struct {
		respCode    int
		respBody    string
		contentType string
	}

	tests := []struct {
		name string
		want want
	}{
		{
			"allowed",
			want{
				respCode:    http.StatusOK,
				respBody:    ".",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			"rejected",
			want{
				respCode:    http.StatusGatewayTimeout,
				respBody:    "Service Timeout\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			"rejected",
			want{
				respCode:    http.StatusTooManyRequests,
				respBody:    "Server capacity exceeded.\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			"rejected",
			want{
				respCode:    http.StatusTooManyRequests,
				respBody:    "Server capacity exceeded.\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			"rejected",
			want{
				respCode:    http.StatusTooManyRequests,
				respBody:    "Server capacity exceeded.\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			"allowed",
			want{
				respCode:    http.StatusOK,
				respBody:    ".",
				contentType: "text/plain; charset=utf-8",
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			resp, err := ts.Client().Get(ts.URL)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if resp.StatusCode != tt.want.respCode {
				t.Errorf("got %d, want %d", resp.StatusCode, tt.want.respCode)
			}

			if ct := resp.Header.Get("Content-Type"); ct != tt.want.contentType {
				t.Errorf("got %q, want %q", ct, tt.want.contentType)
			}

			got, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if string(got) != tt.want.respBody {
				t.Errorf("got %q, want %q", got, tt.want.respBody)
			}
		})
	}
}

func TestThrottleWithOptsBacklogTimeoutWithContextNoStatusCode(t *testing.T) {
	t.Parallel()

	c := chi.NewRouter()
	c.Use(middleware.Timeout(time.Millisecond * 50))
	c.Use(ThrottleWithOpts(ThrottleOpts{
		Limit:          1,
		BacklogLimit:   1,
		BacklogTimeout: time.Millisecond * 10,
	}))
	c.Get("/", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Millisecond * 30)
		w.Write([]byte("."))
	})

	ts := httptest.NewServer(c)
	defer ts.Close()

	type want struct {
		respCode    int
		respBody    string
		contentType string
	}

	tests := []struct {
		name string
		want want
	}{
		{
			"allowed",
			want{
				respCode:    http.StatusOK,
				respBody:    ".",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			"rejected",
			want{
				respCode:    http.StatusGatewayTimeout,
				respBody:    "Service Timeout\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			"rejected",
			want{
				respCode:    http.StatusTooManyRequests,
				respBody:    "Server capacity exceeded.\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			"rejected",
			want{
				respCode:    http.StatusTooManyRequests,
				respBody:    "Server capacity exceeded.\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			"rejected",
			want{
				respCode:    http.StatusTooManyRequests,
				respBody:    "Server capacity exceeded.\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			"allowed",
			want{
				respCode:    http.StatusOK,
				respBody:    ".",
				contentType: "text/plain; charset=utf-8",
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			resp, err := ts.Client().Get(ts.URL)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if resp.StatusCode != tt.want.respCode {
				t.Errorf("got %d, want %d", resp.StatusCode, tt.want.respCode)
			}

			if ct := resp.Header.Get("Content-Type"); ct != tt.want.contentType {
				t.Errorf("got %q, want %q", ct, tt.want.contentType)
			}

			got, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if string(got) != tt.want.respBody {
				t.Errorf("got %q, want %q", got, tt.want.respBody)
			}
		})
	}
}

func TestThrottleWithOptsBacklogTimeoutWithContextWithStatusCode(t *testing.T) {
	t.Parallel()

	c := chi.NewRouter()
	c.Use(middleware.Timeout(time.Millisecond * 50))
	c.Use(ThrottleWithOpts(ThrottleOpts{
		Limit:          1,
		BacklogLimit:   1,
		BacklogTimeout: time.Millisecond * 10,
		StatusCode:     http.StatusInternalServerError,
	}))
	c.Get("/", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Millisecond * 30)
		w.Write([]byte("."))
	})

	ts := httptest.NewServer(c)
	defer ts.Close()

	type want struct {
		respCode    int
		respBody    string
		contentType string
	}

	tests := []struct {
		name string
		want want
	}{
		{
			"allowed",
			want{
				respCode:    http.StatusOK,
				respBody:    ".",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			"rejected",
			want{
				respCode:    http.StatusGatewayTimeout,
				respBody:    "Service Timeout\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			"rejected",
			want{
				respCode:    http.StatusInternalServerError,
				respBody:    "Server capacity exceeded.\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			"rejected",
			want{
				respCode:    http.StatusInternalServerError,
				respBody:    "Server capacity exceeded.\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			"rejected",
			want{
				respCode:    http.StatusInternalServerError,
				respBody:    "Server capacity exceeded.\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			"allowed",
			want{
				respCode:    http.StatusOK,
				respBody:    ".",
				contentType: "text/plain; charset=utf-8",
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			resp, err := ts.Client().Get(ts.URL)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if resp.StatusCode != tt.want.respCode {
				t.Errorf("got %d, want %d", resp.StatusCode, tt.want.respCode)
			}

			if ct := resp.Header.Get("Content-Type"); ct != tt.want.contentType {
				t.Errorf("got %q, want %q", ct, tt.want.contentType)
			}

			got, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if string(got) != tt.want.respBody {
				t.Errorf("got %q, want %q", got, tt.want.respBody)
			}
		})
	}
}

func TestThrottleWithOptsPanicInvalidBacklogTimeout(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		opts ThrottleOpts
	}{
		{
			"negative limit",
			ThrottleOpts{Limit: -1},
		},
		{
			"zero limit",
			ThrottleOpts{Limit: 0},
		},
		{
			"negative backlog limit",
			ThrottleOpts{Limit: 1, BacklogLimit: -1},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assertPanic(t, func() { ThrottleWithOpts(tt.opts) })
		})
	}
}

func Test_throttler_setRetryAfterHeaderIfNeeded(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		retryAfterFn func(ctxDone bool) time.Duration
		ctxDone     bool
		want        string
	}{
		{
			name:        "nil retryAfterFn",
			retryAfterFn: nil,
			ctxDone:     false,
			want:        "",
		},
		{
			name:        "retryAfterFn returns 1 second",
			retryAfterFn: func(ctxDone bool) time.Duration { return 1 * time.Second },
			ctxDone:     false,
			want:        "1",
		},
		{
			name:        "retryAfterFn returns 2 seconds",
			retryAfterFn: func(ctxDone bool) time.Duration { return 2 * time.Second },
			ctxDone:     true,
			want:        "2",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			rr := httptest.NewRecorder()

			tr := throttler{retryAfterFn: tt.retryAfterFn}
			tr.setRetryAfterHeaderIfNeeded(rr, tt.ctxDone)

			if got := rr.Header().Get("Retry-After"); got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestThrottleBacklogTimeoutWithContext(t *testing.T) {
	t.Parallel()

	c := chi.NewRouter()
	c.Use(middleware.Timeout(time.Millisecond * 50))
	c.Use(ThrottleBacklog(1, 1, time.Second))
	c.Get("/", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Millisecond * 100)
		w.Write([]byte("."))
	})

	ts := httptest.NewServer(c)
	defer ts.Close()

	type want struct {
		respCode    int
		respBody    string
		contentType string
	}

	tests := []struct {
		name string
		want want
	}{
		{
			"allowed",
			want{
				respCode:    http.StatusOK,
				respBody:    ".",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			"rejected",
			want{
				respCode:    http.StatusGatewayTimeout,
				respBody:    "Service Timeout\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			"rejected",
			want{
				respCode:    http.StatusTooManyRequests,
				respBody:    "Server capacity exceeded.\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			"rejected",
			want{
				respCode:    http.StatusTooManyRequests,
				respBody:    "Server capacity exceeded.\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			"rejected",
			want{
				respCode:    http.StatusTooManyRequests,
				respBody:    "Server capacity exceeded.\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			"allowed",
			want{
				respCode:    http.StatusOK,
				respBody:    ".",
				contentType: "text/plain; charset=utf-8",
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			resp, err := ts.Client().Get(ts.URL)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if resp.StatusCode != tt.want.respCode {
				t.Errorf("got %d, want %d", resp.StatusCode, tt.want.respCode)
			}

			if ct := resp.Header.Get("Content-Type"); ct != tt.want.contentType {
				t.Errorf("got %q, want %q", ct, tt.want.contentType)
			}

			got, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if string(got) != tt.want.respBody {
				t.Errorf("got %q, want %q", got, tt.want.respBody)
			}
		})
	}
}

func TestThrottleWithOptsBacklogTimeoutWithContextNoStatusCode(t *testing.T) {
	t.Parallel()

	c := chi.NewRouter()
	c.Use(middleware.Timeout(time.Millisecond * 50))
	c.Use(ThrottleWithOpts(ThrottleOpts{
		Limit:          1,
		BacklogLimit:   1,
		BacklogTimeout: time.Millisecond * 10,
	}))
	c.Get("/", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Millisecond * 30)
		w.Write([]byte("."))
	})

	ts := httptest.NewServer(c)
	defer ts.Close()

	type want struct {
		respCode    int
		respBody    string
		contentType string
	}

	tests := []struct {
		name string
		want want
	}{
		{
			"allowed",
			want{
				respCode:    http.StatusOK,
				respBody:    ".",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			"rejected",
			want{
				respCode:    http.StatusGatewayTimeout,
				respBody:    "Service Timeout\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			"rejected",
			want{
				respCode:    http.StatusTooManyRequests,
				respBody:    "Server capacity exceeded.\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			"rejected",
			want{
				respCode:    http.StatusTooManyRequests,
				respBody:    "Server capacity exceeded.\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			"rejected",
			want{
				respCode:    http.StatusTooManyRequests,
				respBody:    "Server capacity exceeded.\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			"allowed",
			want{
				respCode:    http.StatusOK,
				respBody:    ".",
				contentType: "text/plain; charset=utf-8",
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			resp, err := ts.Client().Get(ts.URL)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if resp.StatusCode != tt.want.respCode {
				t.Errorf("got %d, want %d", resp.StatusCode, tt.want.respCode)
			}

			if ct := resp.Header.Get("Content-Type"); ct != tt.want.contentType {
				t.Errorf("got %q, want %q", ct, tt.want.contentType)
			}

			got, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if string(got) != tt.want.respBody {
				t.Errorf("got %q, want %q", got, tt.want.respBody)
			}
		})
	}
}

func TestThrottleWithOptsBacklogTimeoutWithContextWithStatusCode(t *testing.T) {
	t.Parallel()

	c := chi.NewRouter()
	c.Use(middleware.Timeout(time.Millisecond * 50))
	c.Use(ThrottleWithOpts(ThrottleOpts{
		Limit:          1,
		BacklogLimit:   1,
		BacklogTimeout: time.Millisecond * 10,
		StatusCode:     http.StatusInternalServerError,
	}))
	c.Get("/", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Millisecond * 30)
		w.Write([]byte("."))
	})

	ts := httptest.NewServer(c)
	defer ts.Close()

	type want struct {
		respCode    int
		respBody    string
		contentType string
	}

	tests := []struct {
		name string
		want want
	}{
		{
			"allowed",
			want{
				respCode:    http.StatusOK,
				respBody:    ".",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			"rejected",
			want{
				respCode:    http.StatusGatewayTimeout,
				respBody:    "Service Timeout\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			"rejected",
			want{
				respCode:    http.StatusInternalServerError,
				respBody:    "Server capacity exceeded.\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			"rejected",
			want{
				respCode:    http.StatusInternalServerError,
				respBody:    "Server capacity exceeded.\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			"rejected",
			want{
				respCode:    http.StatusInternalServerError,
				respBody:    "Server capacity exceeded.\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			"allowed",
			want{
				respCode:    http.StatusOK,
				respBody:    ".",
				contentType: "text/plain; charset=utf-8",
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			resp, err := ts.Client().Get(ts.URL)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if resp.StatusCode != tt.want.respCode {
				t.Errorf("got %d, want %d", resp.StatusCode, tt.want.respCode)
			}

			if ct := resp.Header.Get("Content-Type"); ct != tt.want.contentType {
				t.Errorf("got %q, want %q", ct, tt.want.contentType)
			}

			got, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if string(got) != tt.want.respBody {
				t.Errorf("got %q, want %q", got, tt.want.respBody)
			}
		})
	}
}

func Test_throttler_setRetryAfterHeaderIfNeeded(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		retryAfterFn func(ctxDone bool) time.Duration
		ctxDone     bool
		want        string
	}{
		{
			name:        "nil retryAfterFn",
			retryAfterFn: nil,
			ctxDone:     false,
			want:        "",
		},
		{
			name:        "retryAfterFn returns 1 second",
			retryAfterFn: func(ctxDone bool) time.Duration { return 1 * time.Second },
			ctxDone:     false,
			want:        "1",
		},
		{
			name:        "retryAfterFn returns 2 seconds",
			retryAfterFn: func(ctxDone bool) time.Duration { return 2 * time.Second },
			ctxDone:     true,
			want:        "2",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			rr := httptest.NewRecorder()

			tr := throttler{retryAfterFn: tt.retryAfterFn}
			tr.setRetryAfterHeaderIfNeeded(rr, tt.ctxDone)

			if got := rr.Header().Get("Retry-After"); got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestThrottleBacklogTimeoutWithContext(t *testing.T) {
	t.Parallel()

	c := chi.NewRouter()
	c.Use(middleware.Timeout(time.Millisecond * 50))
	c.Use(ThrottleBacklog(1, 1, time.Second))
	c.Get("/", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Millisecond * 100)
		w.Write([]byte("."))
	})

	ts := httptest.NewServer(c)
	defer ts.Close()

	type want struct {
		respCode    int
		respBody    string
		contentType string
	}

	tests := []struct {
		name string
		want want
	}{
		{
			"allowed",
			want{
				respCode:    http.StatusOK,
				respBody:    ".",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			"rejected",
			want{
				respCode:    http.StatusGatewayTimeout,
				respBody:    "Service Timeout\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			"rejected",
			want{
				respCode:    http.StatusTooManyRequests,
				respBody:    "Server capacity exceeded.\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			"rejected",
			want{
				respCode:    http.StatusTooManyRequests,
				respBody:    "Server capacity exceeded.\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			"rejected",
			want{
				respCode:    http.StatusTooManyRequests,
				respBody:    "Server capacity exceeded.\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			"allowed",
			want{
				respCode:    http.StatusOK,
				respBody:    ".",
				contentType: "text/plain; charset=utf-8",
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			resp, err := ts.Client().Get(ts.URL)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if resp.StatusCode != tt.want.respCode {
				t.Errorf("got %d, want %d", resp.StatusCode, tt.want.respCode)
			}

			if ct := resp.Header.Get("Content-Type"); ct != tt.want.contentType {
				t.Errorf("got %q, want %q", ct, tt.want.contentType)
			}

			got, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if string(got) != tt.want.respBody {
				t.Errorf("got %q, want %q", got, tt.want.respBody)
			}
		})
	}
}

func TestThrottleWithOptsBacklogTimeoutWithContextNoStatusCode(t *testing.T) {
	t.Parallel()

	c := chi.NewRouter()
	c.Use(middleware.Timeout(time.Millisecond * 50))
	c.Use(ThrottleWithOpts(ThrottleOpts{
		Limit:          1,
		BacklogLimit:   1,
		BacklogTimeout: time.Millisecond * 10,
	}))
	c.Get("/", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Millisecond * 30)
		w.Write([]byte("."))
	})

	ts := httptest.NewServer(c)
	defer ts.Close()

	type want struct {
		respCode    int
		respBody    string
		contentType string
	}

	tests := []struct {
		name string
		want want
	}{
		{
			"allowed",
			want{
				respCode:    http.StatusOK,
				respBody:    ".",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			"rejected",
			want{
				respCode:    http.StatusGatewayTimeout,
				respBody:    "Service Timeout\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			"rejected",
			want{
				respCode:    http.StatusTooManyRequests,
				respBody:    "Server capacity exceeded.\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			"rejected",
			want{
				respCode:    http.StatusTooManyRequests,
				respBody:    "Server capacity exceeded.\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			"rejected",
			want{
				respCode:    http.StatusTooManyRequests,
				respBody:    "Server capacity exceeded.\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			"allowed",
			want{
				respCode:    http.StatusOK,
				respBody:    ".",
				contentType: "text/plain; charset=utf-8",
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			resp, err := ts.Client().Get(ts.URL)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if resp.StatusCode != tt.want.respCode {
				t.Errorf("got %d, want %d", resp.StatusCode, tt.want.respCode)
			}

			if ct := resp.Header.Get("Content-Type"); ct != tt.want.contentType {
				t.Errorf("got %q, want %q", ct, tt.want.contentType)
			}

			got, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if string(got) != tt.want.respBody {
				t.Errorf("got %q, want %q", got, tt.want.respBody)
			}
		})
	}
}

func TestThrottleWithOptsBacklogTimeoutWithContextWithStatusCode(t *testing.T) {
	t.Parallel()

	c := chi.NewRouter()
	c.Use(middleware.Timeout(time.Millisecond * 50))
	c.Use(ThrottleWithOpts(ThrottleOpts{
		Limit:          1,
		BacklogLimit:   1,
		BacklogTimeout: time.Millisecond * 10,
		StatusCode:     http.StatusInternalServerError,
	}))
	c.Get("/", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Millisecond * 30)
		w.Write([]byte("."))
	})

	ts := httptest.NewServer(c)
	defer ts.Close()

	type want struct {
		respCode    int
		respBody    string
		contentType string
	}

	tests := []struct {
		name string
		want want
	}{
		{
			"allowed",
			want{
				respCode:    http.StatusOK,
				respBody:    ".",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			"rejected",
			want{
				respCode:    http.StatusGatewayTimeout,
				respBody:    "Service Timeout\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			"rejected",
			want{
				respCode:    http.StatusInternalServerError,
				respBody:    "Server capacity exceeded.\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			"rejected",
			want{
				respCode:    http.StatusInternalServerError,
				respBody:    "Server capacity exceeded.\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			"rejected",
			want{
				respCode:    http.StatusInternalServerError,
				respBody:    "Server capacity exceeded.\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			"allowed",
			want{
				respCode:    http.StatusOK,
				respBody:    ".",
				contentType: "text/plain; charset=utf-8",
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			resp, err := ts.Client().Get(ts.URL)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if resp.StatusCode != tt.want.respCode {
				t.Errorf("got %d, want %d", resp.StatusCode, tt.want.respCode)
			}

			if ct := resp.Header.Get("Content-Type"); ct != tt.want.contentType {
				t.Errorf("got %q, want %q", ct, tt.want.contentType)
			}

			got, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if string(got) != tt.want.respBody {
				t.Errorf("got %q, want %q", got, tt.want.respBody)
			}
		})
	}
}

func TestThrottleBacklogTimeoutWithContext(t *testing.T) {
	t.Parallel()

	c := chi.NewRouter()
	c.Use(middleware.Timeout(time.Millisecond * 50))
	c.Use(ThrottleBacklog(1, 1, time.Second))
	c.Get("/", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Millisecond * 100)
		w.Write([]byte("."))
	})

	ts := httptest.NewServer(c)
	defer ts.Close()

	type want struct {
		respCode    int
		respBody    string
		contentType string
	}

	tests := []struct {
		name string
		want want
	}{
		{
			"allowed",
			want{
				respCode:    http.StatusOK,
				respBody:    ".",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			"rejected",
			want{
				respCode:    http.StatusGatewayTimeout,
				respBody:    "Service Timeout\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			"rejected",
			want{
				respCode:    http.StatusTooManyRequests,
				respBody:    "Server capacity exceeded.\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			"rejected",
			want{
				respCode:    http.StatusTooManyRequests,
				respBody:    "Server capacity exceeded.\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			"rejected",
			want{
				respCode:    http.StatusTooManyRequests,
				respBody:    "Server capacity exceeded.\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			"allowed",
			want{
				respCode:    http.StatusOK,
				respBody:    ".",
				contentType: "text/plain; charset=utf-8",
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			resp, err := ts.Client().Get(ts.URL)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if resp.StatusCode != tt.want.respCode {
				t.Errorf("got %d, want %d", resp.StatusCode, tt.want.respCode)
			}

			if ct := resp.Header.Get("Content-Type"); ct != tt.want.contentType {
				t.Errorf("got %q, want %q", ct, tt.want.contentType)
			}

			got, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if string(got) != tt.want.respBody {
				t.Errorf("got %q, want %q", got, tt.want.respBody)
			}
		})
	}
}

func TestThrottleWithOptsBacklogTimeoutWithContextNoStatusCode(t *testing.T) {
	t.Parallel()

	c := chi.NewRouter()
	c.Use(middleware.Timeout(time.Millisecond * 50))
	c.Use(ThrottleWithOpts(ThrottleOpts{
		Limit:          1,
		BacklogLimit:   1,
		BacklogTimeout: time.Millisecond * 10,
	}))
	c.Get("/", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Millisecond * 30)
		w.Write([]byte("."))
	})

	ts := httptest.NewServer(c)
	defer ts.Close()

	type want struct {
		respCode    int
		respBody    string
		contentType string
	}

	tests := []struct {
		name string
		want want
	}{
		{
			"allowed",
			want{
				respCode:    http.StatusOK,
				respBody:    ".",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			"rejected",
			want{
				respCode:    http.StatusGatewayTimeout,
				respBody:    "Service Timeout\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			"rejected",
			want{
				respCode:    http.StatusTooManyRequests,
				respBody:    "Server capacity exceeded.\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			"rejected",
			want{
				respCode:    http.StatusTooManyRequests,
				respBody:    "Server capacity exceeded.\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			"rejected",
			want{
				respCode:    http.StatusTooManyRequests,
				respBody:    "Server capacity exceeded.\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			"allowed",
			want{
				respCode:    http.StatusOK,
				respBody:    ".",
				contentType: "text/plain; charset=utf-8",
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			resp, err := ts.Client().Get(ts.URL)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if resp.StatusCode != tt.want.respCode {
				t.Errorf("got %d, want %d", resp.StatusCode, tt.want.respCode)
			}

			if ct := resp.Header.Get("Content-Type"); ct != tt.want.contentType {
				t.Errorf("got %q, want %q", ct, tt.want.contentType)
			}

			got, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if string(got) != tt.want.respBody {
				t.Errorf("got %q, want %q", got, tt.want.respBody)
			}
		})
	}
}

func TestThrottleWithOptsBacklogTimeoutWithContextWithStatusCode(t *testing.T) {
	t.Parallel()

	c := chi.NewRouter()
	c.Use(middleware.Timeout(time.Millisecond * 50))
	c.Use(ThrottleWithOpts(ThrottleOpts{
		Limit:          1,
		BacklogLimit:   1,
		BacklogTimeout: time.Millisecond * 10,
		StatusCode:     http.StatusInternalServerError,
	}))
	c.Get("/", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Millisecond * 30)
		w.Write([]byte("."))
	})

	ts := httptest.NewServer(c)
	defer ts.Close()

	type want struct {
		respCode    int
		respBody    string
		contentType string
	}

	tests := []struct {
		name string
		want want
	}{
		{
			"allowed",
			want{
				respCode:    http.StatusOK,
				respBody:    ".",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			"rejected",
			want{
				respCode:    http.StatusGatewayTimeout,
				respBody:    "Service Timeout\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			"rejected",
			want{
				respCode:    http.StatusInternalServerError,
				respBody:    "Server capacity exceeded.\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			"rejected",
			want{
				respCode:    http.StatusInternalServerError,
				respBody:    "Server capacity exceeded.\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			"rejected",
			want{
				respCode:    http.StatusInternalServerError,
				respBody:    "Server capacity exceeded.\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			"allowed",
			want{
				respCode:    http.StatusOK,
				respBody:    ".",
				contentType: "text/plain; charset=utf-8",
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			resp, err := ts.Client().Get(ts.URL)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if resp.StatusCode != tt.want.respCode {
				t.Errorf("got %d, want %d", resp.StatusCode, tt.want.respCode)
			}

			if ct := resp.Header.Get("Content-Type"); ct != tt.want.contentType {
				t.Errorf("got %q, want %q", ct, tt.want.contentType)
			}

			got, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if string(got) != tt.want.respBody {
				t.Errorf("got %q, want %q", got, tt.want.respBody)
			}
		})
	}
}

func TestThrottleWithOptsPanicInvalidBacklogTimeout(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		opts ThrottleOpts
	}{
		{
			"negative limit",
			ThrottleOpts{Limit: -1},
		},
		{
			"zero limit",
			ThrottleOpts{Limit: 0},
		},
		{
			"negative backlog limit",
			ThrottleOpts{Limit: 1, BacklogLimit: -1},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assertPanic(t, func() { ThrottleWithOpts(tt.opts) })
		})
	}
}

func Test_throttler_setRetryAfterHeaderIfNeeded(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		retryAfterFn func(ctxDone bool) time.Duration
		ctxDone     bool
		want        string
	}{
		{
			name:        "nil retryAfterFn",
			retryAfterFn: nil,
			ctxDone:     false,
			want:        "",
		},
		{
			name:        "retryAfterFn returns 1 second",
			retryAfterFn: func(ctxDone bool) time.Duration { return 1 * time.Second },
			ctxDone:     false,
			want:        "1",
		},
		{
			name:        "retryAfterFn returns 2 seconds",
			retryAfterFn: func(ctxDone bool) time.Duration { return 2 * time.Second },
			ctxDone:     true,
			want:        "2",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			rr := httptest.NewRecorder()

			tr := throttler{retryAfterFn: tt.retryAfterFn}
			tr.setRetryAfterHeaderIfNeeded(rr, tt.ctxDone)

			if got := rr.Header().Get("Retry-After"); got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestThrottleBacklogTimeoutWithContext(t *testing.T) {
	t.Parallel()

	c := chi.NewRouter()
	c.Use(middleware.Timeout(time.Millisecond * 50))
	c.Use(ThrottleBacklog(1, 1, time.Second))
	c.Get("/", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Millisecond * 100)
		w.Write([]byte("."))
	})

	ts := httptest.NewServer(c)
	defer ts.Close()

	type want struct {
		respCode    int
		respBody    string
		contentType string
	}

	tests := []struct {
		name string
		want want
	}{
		{
			"allowed",
			want{
				respCode:    http.StatusOK,
				respBody:    ".",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			"rejected",
			want{
				respCode:    http.StatusGatewayTimeout,
				respBody:    "Service Timeout\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			"rejected",
			want{
				respCode:    http.StatusTooManyRequests,
				respBody:    "Server capacity exceeded.\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			"rejected",
			want{
				respCode:    http.StatusTooManyRequests,
				respBody:    "Server capacity exceeded.\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			"rejected",
			want{
				respCode:    http.StatusTooManyRequests,
				respBody:    "Server capacity exceeded.\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			"allowed",
			want{
				respCode:    http.StatusOK,
				respBody:    ".",
				contentType: "text/plain; charset=utf-8",
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			resp, err := ts.Client().Get(ts.URL)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if resp.StatusCode != tt.want.respCode {
				t.Errorf("got %d, want %d", resp.StatusCode, tt.want.respCode)
			}

			if ct := resp.Header.Get("Content-Type"); ct != tt.want.contentType {
				t.Errorf("got %q, want %q", ct, tt.want.contentType)
			}

			got, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if string(got) != tt.want.respBody {
				t.Errorf("got %q, want %q", got, tt.want.respBody)
			}
		})
	}
}

func TestThrottleWithOptsBacklogTimeoutWithContextNoStatusCode(t *testing.T) {
	t.Parallel()

	c := chi.NewRouter()
	c.Use(middleware.Timeout(time.Millisecond * 50))
	c.Use(ThrottleWithOpts(ThrottleOpts{
		Limit:          1,
		BacklogLimit:   1,
		BacklogTimeout: time.Millisecond * 10,
	}))
	c.Get("/", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Millisecond * 30)
		w.Write([]byte("."))
	})

	ts := httptest.NewServer(c)
	defer ts.Close()

	type want struct {
		respCode    int
		respBody    string
		contentType string
	}

	tests := []struct {
		name string
		want want
	}{
		{
			"allowed",
			want{
				respCode:    http.StatusOK,
				respBody:    ".",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			"rejected",
			want{
				respCode:    http.StatusGatewayTimeout,
				respBody:    "Service Timeout\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			"rejected",
			want{
				respCode:    http.StatusTooManyRequests,
				respBody:    "Server capacity exceeded.\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			"rejected",
			want{
				respCode:    http.StatusTooManyRequests,
				respBody:    "Server capacity exceeded.\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			"rejected",
			want{
				respCode:    http.StatusTooManyRequests,
				respBody:    "Server capacity exceeded.\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			"allowed",
			want{
				respCode:    http.StatusOK,
				respBody:    ".",
				contentType: "text/plain; charset=utf-8",
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			resp, err := ts.Client().Get(ts.URL)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if resp.StatusCode != tt.want.respCode {
				t.Errorf("got %d, want %d", resp.StatusCode, tt.want.respCode)
			}

			if ct := resp.Header.Get("Content-Type"); ct != tt.want.contentType {
				t.Errorf("got %q, want %q", ct, tt.want.contentType)
			}

			got, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if string(got) != tt.want.respBody {
				t.Errorf("got %q, want %q", got, tt.want.respBody)
			}
		})
	}
}

func TestThrottleWithOptsBacklogTimeoutWithContextWithStatusCode(t *testing.T) {
	t.Parallel()

	c := chi.NewRouter()
	c.Use(middleware.Timeout(time.Millisecond * 50))
	c.Use(ThrottleWithOpts(ThrottleOpts{
		Limit:          1,
		BacklogLimit:   1,
		BacklogTimeout: time.Millisecond * 10,
		StatusCode:     http.StatusInternalServerError,
	}))
	c.Get("/", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Millisecond * 30)
		w.Write([]byte("."))
	})

	ts := httptest.NewServer(c)
	defer ts.Close()

	type want struct {
		respCode    int
		respBody    string
		contentType string
	}

	tests := []struct {
		name string
		want want
	}{
		{
			"allowed",
			want{
				respCode:    http.StatusOK,
				respBody:    ".",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			"rejected",
			want{
				respCode:    http.StatusGatewayTimeout,
				respBody:    "Service Timeout\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			"rejected",
			want{
				respCode:    http.StatusInternalServerError,
				respBody:    "Server capacity exceeded.\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			"rejected",
			want{
				respCode:    http.StatusInternalServerError,
				respBody:    "Server capacity exceeded.\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			"rejected",
			want{
				respCode:    http.StatusInternalServerError,
				respBody:    "Server capacity exceeded.\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			"allowed",
			want{
				respCode:    http.StatusOK,
				respBody:    ".",
				contentType: "text/plain; charset=utf-8",
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			resp, err := ts.Client().Get(ts.URL)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if resp.StatusCode != tt.want.respCode {
				t.Errorf("got %d, want %d", resp.StatusCode, tt.want.respCode)
			}

			if ct := resp.Header.Get("Content-Type"); ct != tt.want.contentType {
				t.Errorf("got %q, want %q", ct, tt.want.contentType)
			}

			got, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if string(got) != tt.want.respBody {
				t.Errorf("got %q, want %q", got, tt.want.respBody)
			}
		})
	}
}

func TestThrottleWithOptsPanicInvalidBacklogTimeout(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		opts ThrottleOpts
	}{
		{
			"negative limit",
			ThrottleOpts{Limit: -1},
		},
		{
			"zero limit",
			ThrottleOpts{Limit: 0},
		},
		{
			"negative backlog limit",
			ThrottleOpts{Limit: 1, BacklogLimit: -1},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assertPanic(t, func() { ThrottleWithOpts(tt.opts) })
		})
	}
}

func Test_throttler_setRetryAfterHeaderIfNeeded(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		retryAfterFn func(ctxDone bool) time.Duration
		ctxDone     bool
		want        string
	}{
		{
			name:        "nil retryAfterFn",
			retryAfterFn: nil,
			ctxDone:     false,
			want:        "",
		},
		{
			name:        "retryAfterFn returns 1 second",
			retryAfterFn: func(ctxDone bool) time.Duration { return 1 * time.Second },
			ctxDone:     false,
			want:        "1",
		},
		{
			name:        "retryAfterFn returns 2 seconds",
			retryAfterFn: func(ctxDone bool) time.Duration { return 2 * time.Second },
			ctxDone:     true,
			want:        "2",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			rr := httptest.NewRecorder()

			tr := throttler{retryAfterFn: tt.retryAfterFn}
			tr.setRetryAfterHeaderIfNeeded(rr, tt.ctxDone)

			if got := rr.Header().Get("Retry-After"); got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestThrottleBacklogTimeoutWithContext(t *testing.T) {
	t.Parallel()

	c := chi.NewRouter()
	c.Use(middleware.Timeout(time.Millisecond * 50))
	c.Use(ThrottleBacklog(1, 1, time.Second))
	c.Get("/", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Millisecond * 100)
		w.Write([]byte("."))
	})

	ts := httptest.NewServer(c)
	defer ts.Close()

	type want struct {
		respCode    int
		respBody    string
		contentType string
	}

	tests := []struct {
		name string