package middleware

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
)

func TestThrottle(t *testing.T) {
	limit := 10
	mux := chi.NewRouter()
	mux.Use(Throttle(limit))
	mux.Get("/", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Second * 1)
		w.Write([]byte("."))
	})

	for i := 0; i < limit; i++ {
		r := httptest.NewRequest("GET", "/", nil)
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, r)
		if w.Code != http.StatusOK {
			t.Fatalf("expecting %d got %d", http.StatusOK, w.Code)
		}
	}

	r := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)
	if w.Code != http.StatusTooManyRequests {
		t.Fatalf("expecting %d got %d", http.StatusTooManyRequests, w.Code)
	}
}

func TestThrottleBacklog(t *testing.T) {
	limit := 10
	backlogLimit := 20
	mux := chi.NewRouter()
	mux.Use(ThrottleBacklog(limit, backlogLimit, time.Second*3))
	mux.Get("/", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Second * 1)
		w.Write([]byte("."))
	})

	var pending uint32

	for i := 0; i < limit+backlogLimit; i++ {
		go func() {
			r := httptest.NewRequest("GET", "/", nil)
			w := httptest.NewRecorder()
			mux.ServeHTTP(w, r)
			if w.Code != http.StatusOK {
				t.Errorf("expecting %d got %d", http.StatusOK, w.Code)
			}
			atomic.AddUint32(&pending, ^uint32(0))
		}()
		atomic.AddUint32(&pending, 1)
	}

	time.Sleep(time.Second * 1)

	r := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)
	if w.Code != http.StatusTooManyRequests {
		t.Errorf("expecting %d got %d", http.StatusTooManyRequests, w.Code)
	}
	time.Sleep(time.Second * 2)

	for {
		if atomic.LoadUint32(&pending) == 0 {
			break
		}
	}
}

func TestThrottleWithOpts(t *testing.T) {
	run := func(opts ThrottleOpts) {
		if opts.Limit == 0 {
			return
		}
		mux := chi.NewRouter()
		mux.Use(ThrottleWithOpts(opts))
		mux.Get("/", func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(time.Second * 1)
			w.Write([]byte("."))
		})

		var pending uint32

		for i := 0; i < opts.Limit+opts.BacklogLimit; i++ {
			go func() {
				r := httptest.NewRequest("GET", "/", nil)
				w := httptest.NewRecorder()
				mux.ServeHTTP(w, r)
				if w.Code != http.StatusOK {
					t.Errorf("expecting %d got %d", http.StatusOK, w.Code)
				}
				atomic.AddUint32(&pending, ^uint32(0))
			}()
			atomic.AddUint32(&pending, 1)
		}

		time.Sleep(time.Second * 1)

		r := httptest.NewRequest("GET", "/", nil)
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, r)
		if w.Code != opts.StatusCode {
			t.Errorf("expecting %d got %d", opts.StatusCode, w.Code)
		}
		if opts.RetryAfterFn != nil {
			expected := strconv.Itoa(int(opts.RetryAfterFn(false).Seconds()))
			if w.Header().Get("Retry-After") != expected {
				t.Errorf("expecting %s got %s", expected, w.Header().Get("Retry-After"))
			}
		}
		time.Sleep(time.Second * 2)

		for {
			if atomic.LoadUint32(&pending) == 0 {
				break
			}
		}
	}

	tests := []struct {
		opts ThrottleOpts
	}{
		{ThrottleOpts{}},
		{ThrottleOpts{Limit: 10}},
		{ThrottleOpts{Limit: 10, BacklogLimit: 20}},
		{ThrottleOpts{Limit: 10, StatusCode: http.StatusTeapot}},
		{ThrottleOpts{Limit: 10, RetryAfterFn: func(_ bool) time.Duration { return time.Second * 10 }}},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%+v", test.opts), func(t *testing.T) {
			run(test.opts)
		})
	}
}