package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRequestID(t *testing.T) {
	mux := http.NewServeMux()
	mux.Handle("/", RequestID(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := GetReqID(r.Context())
		assert.NotEmpty(t, reqID)
	})))

	ts := httptest.NewServer(mux)
	defer ts.Close()

	resp, err := http.Get(ts.URL)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}