package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRecovererMiddleware(t *testing.T) {
	testCases := []struct {
		name          string
		handlerFunc   http.HandlerFunc
		expectStatus  int
		expectMessage string
	}{
		{
			name: "no panic",
			handlerFunc: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("OK"))
			},
			expectStatus:  http.StatusOK,
			expectMessage: "OK",
		},
		{
			name: "panic with message",
			handlerFunc: func(w http.ResponseWriter, r *http.Request) {
				panic("something went wrong")
			},
			expectStatus:  http.StatusInternalServerError,
			expectMessage: "",
		},
		{
			name: "panic with http.ErrAbortHandler",
			handlerFunc: func(w http.ResponseWriter, r *http.Request) {
				panic(http.ErrAbortHandler)
			},
			expectStatus:  http.StatusOK,
			expectMessage: "OK",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a request to pass to our handler.
			req, err := http.NewRequest("GET", "http://example.com/foo", nil)
			if err != nil {
				t.Fatal(err)
			}

			// Create a ResponseRecorder to record the response.
			rr := httptest.NewRecorder()

			// Create handler with Recoverer middleware.
			handler := Recoverer(tc.handlerFunc)

			// Serve HTTP
			handler.ServeHTTP(rr, req)

			// Check the status code is what we expect.
			if status := rr.Code; status != tc.expectStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tc.expectStatus)
			}

			// Check the response message.
			if strings.TrimSpace(rr.Body.String()) != tc.expectMessage {
				t.Errorf("handler returned unexpected body: got %v want %v",
					rr.Body.String(), tc.expectMessage)
			}
		})
	}
}