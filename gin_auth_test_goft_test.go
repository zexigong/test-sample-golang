// Copyright 2014 Manu Martinez-Almeida. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package gin

import (
	"bytes"
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBasicAuth(t *testing.T) {
	for _, tt := range []struct {
		accounts Accounts
		name     string
	}{
		{Accounts{"admin": "password"}, "single basic account"},
		{Accounts{"admin": "password", "test": "password"}, "multiple basic accounts"},
	} {
		// build fresh accounts for each iteration, since processAccount modifies the accounts
		accounts := tt.accounts
		accountsProcessed := processAccounts(accounts)
		assert.Len(t, accountsProcessed, len(accounts), tt.name)
		for _, account := range accountsProcessed {
			assert.NotEmpty(t, account.value, tt.name)
			assert.NotEmpty(t, account.user, tt.name)
		}
	}
}

func TestBasicAuthFails(t *testing.T) {
	for _, tt := range []struct {
		accounts Accounts
		name     string
	}{
		{Accounts{"": "password"}, "empty user name"},
		{Accounts{"admin": ""}, "empty password"},
	} {
		// build fresh accounts for each iteration, since processAccount modifies the accounts
		accounts := tt.accounts
		assert.Panics(t, func() {
			processAccounts(accounts)
		}, tt.name)
	}
}

func TestBasicAuthSearchCredential(t *testing.T) {
	assert := assert.New(t)

	accounts := Accounts{"admin": "password", "test": "password"}
	accountsProcessed := processAccounts(accounts)

	for _, tt := range []struct {
		authValue string
		user      string
		found     bool
	}{
		// Test credentials that exist
		{authorizationHeader("admin", "password"), "admin", true},
		{authorizationHeader("test", "password"), "test", true},

		// Test credentials that don't exist
		{authorizationHeader("admin", ""), "", false},
		{authorizationHeader("", "password"), "", false},
		{authorizationHeader("", ""), "", false},
		{authorizationHeader("admin", "wrong_password"), "", false},
		{authorizationHeader("test", "wrong_password"), "", false},
		{authorizationHeader("wrong_user", "wrong_password"), "", false},
		{"", "", false},
		{"Basic", "", false},
	} {
		user, found := accountsProcessed.searchCredential(tt.authValue)
		assert.Equal(tt.user, user)
		assert.Equal(tt.found, found)
	}
}

func TestBasicAuthAuthorizationHeader(t *testing.T) {
	assert := assert.New(t)

	expectedHeader := "Basic " + base64.StdEncoding.EncodeToString(bytes.NewBufferString("admin:password").Bytes())
	assert.Equal(expectedHeader, authorizationHeader("admin", "password"))

	expectedHeader = "Basic " + base64.StdEncoding.EncodeToString(bytes.NewBufferString("test:password").Bytes())
	assert.Equal(expectedHeader, authorizationHeader("test", "password"))
}

func TestBasicAuthMiddleware(t *testing.T) {
	router := New()
	accounts := Accounts{"admin": "password", "test": "password"}
	router.Use(BasicAuth(accounts))
	router.GET("/foo", func(c *Context) {
		c.String(200, "bar")
	})

	// Test a request without authorization
	w := performRequest(router, "GET", "/foo")
	assert.Equal(t, w.Code, 401)
	assert.Equal(t, w.HeaderMap.Get("WWW-Authenticate"), "Basic realm=\"Authorization Required\"")

	// Test a request with invalid credentials
	w = performRequest(router, "GET", "/foo", authorizationHeader("admin", ""))
	assert.Equal(t, w.Code, 401)
	assert.Equal(t, w.HeaderMap.Get("WWW-Authenticate"), "Basic realm=\"Authorization Required\"")

	// Test valid request
	w = performRequest(router, "GET", "/foo", authorizationHeader("admin", "password"))
	assert.Equal(t, w.Code, 200)
	assert.Equal(t, w.Body.String(), "bar")
}

func TestBasicAuthForRealmMiddleware(t *testing.T) {
	router := New()
	accounts := Accounts{"admin": "password", "test": "password"}
	router.Use(BasicAuthForRealm(accounts, "My Custom Realm"))
	router.GET("/foo", func(c *Context) {
		c.String(200, "bar")
	})

	// Test a request without authorization
	w := performRequest(router, "GET", "/foo")
	assert.Equal(t, w.Code, 401)
	assert.Equal(t, w.HeaderMap.Get("WWW-Authenticate"), "Basic realm=\"My Custom Realm\"")

	// Test a request with invalid credentials
	w = performRequest(router, "GET", "/foo", authorizationHeader("admin", ""))
	assert.Equal(t, w.Code, 401)
	assert.Equal(t, w.HeaderMap.Get("WWW-Authenticate"), "Basic realm=\"My Custom Realm\"")

	// Test valid request
	w = performRequest(router, "GET", "/foo", authorizationHeader("admin", "password"))
	assert.Equal(t, w.Code, 200)
	assert.Equal(t, w.Body.String(), "bar")
}

func TestBasicAuthForProxyMiddleware(t *testing.T) {
	router := New()
	accounts := Accounts{"admin": "password", "test": "password"}
	router.Use(BasicAuthForProxy(accounts, "My Custom Realm"))
	router.GET("/foo", func(c *Context) {
		c.String(200, "bar")
	})

	// Test a request without authorization
	w := performRequest(router, "GET", "/foo")
	assert.Equal(t, w.Code, 407)
	assert.Equal(t, w.HeaderMap.Get("Proxy-Authenticate"), "Basic realm=\"My Custom Realm\"")

	// Test a request with invalid credentials
	w = performRequest(router, "GET", "/foo", authorizationHeader("admin", ""))
	assert.Equal(t, w.Code, 407)
	assert.Equal(t, w.HeaderMap.Get("Proxy-Authenticate"), "Basic realm=\"My Custom Realm\"")

	// Test valid request
	w = performRequest(router, "GET", "/foo", authorizationHeader("admin", "password"))
	assert.Equal(t, w.Code, 200)
	assert.Equal(t, w.Body.String(), "bar")
}

func performRequest(r *Engine, method, path string, header ...string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, nil)
	for _, h := range header {
		req.Header.Set("Authorization", h)
		req.Header.Set("Proxy-Authorization", h)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}