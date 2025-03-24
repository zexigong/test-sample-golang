// Copyright 2014 Manu Martinez-Almeida. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package gin

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBasicAuth(t *testing.T) {
	accounts := Accounts{
		"admin": "password",
		"foo":   "bar",
	}
	ba := BasicAuth(accounts)
	c, _ := CreateTestContext(httptest.NewRecorder())

	// Test correct credentials
	authValue := base64.StdEncoding.EncodeToString([]byte("admin:password"))
	c.Request, _ = http.NewRequest("GET", "/", nil)
	c.Request.Header.Set("Authorization", fmt.Sprintf("Basic %s", authValue))
	ba(c)
	assert.False(t, c.IsAborted())
	assert.Equal(t, "admin", c.MustGet(AuthUserKey).(string))

	// Test incorrect credentials
	c, _ = CreateTestContext(httptest.NewRecorder())
	authValue = base64.StdEncoding.EncodeToString([]byte("foo:bar"))
	c.Request, _ = http.NewRequest("GET", "/", nil)
	c.Request.Header.Set("Authorization", fmt.Sprintf("Basic %s", authValue))
	ba(c)
	assert.False(t, c.IsAborted())
	assert.Equal(t, "foo", c.MustGet(AuthUserKey).(string))
}

func TestBasicAuth401(t *testing.T) {
	accounts := Accounts{
		"admin": "password",
		"foo":   "bar",
	}
	ba := BasicAuth(accounts)
	c, _ := CreateTestContext(httptest.NewRecorder())

	// Test incorrect credentials
	authValue := base64.StdEncoding.EncodeToString([]byte("admin:password1"))
	c.Request, _ = http.NewRequest("GET", "/", nil)
	c.Request.Header.Set("Authorization", fmt.Sprintf("Basic %s", authValue))
	ba(c)
	assert.True(t, c.IsAborted())
	assert.Equal(t, http.StatusUnauthorized, c.Writer.Status())

	// Test another incorrect credentials
	c, _ = CreateTestContext(httptest.NewRecorder())
	authValue = base64.StdEncoding.EncodeToString([]byte("admin1:password"))
	c.Request, _ = http.NewRequest("GET", "/", nil)
	c.Request.Header.Set("Authorization", fmt.Sprintf("Basic %s", authValue))
	ba(c)
	assert.True(t, c.IsAborted())
	assert.Equal(t, http.StatusUnauthorized, c.Writer.Status())

	// Test empty Authorization header
	c, _ = CreateTestContext(httptest.NewRecorder())
	c.Request, _ = http.NewRequest("GET", "/", nil)
	ba(c)
	assert.True(t, c.IsAborted())
	assert.Equal(t, http.StatusUnauthorized, c.Writer.Status())

	// Test empty Authorization header
	c, _ = CreateTestContext(httptest.NewRecorder())
	c.Request, _ = http.NewRequest("GET", "/", nil)
	c.Request.Header.Set("Authorization", "")
	ba(c)
	assert.True(t, c.IsAborted())
	assert.Equal(t, http.StatusUnauthorized, c.Writer.Status())

	// Test empty credentials
	c, _ = CreateTestContext(httptest.NewRecorder())
	authValue = base64.StdEncoding.EncodeToString([]byte(":"))
	c.Request, _ = http.NewRequest("GET", "/", nil)
	c.Request.Header.Set("Authorization", fmt.Sprintf("Basic %s", authValue))
	ba(c)
	assert.True(t, c.IsAborted())
	assert.Equal(t, http.StatusUnauthorized, c.Writer.Status())

	// Test empty user
	c, _ = CreateTestContext(httptest.NewRecorder())
	authValue = base64.StdEncoding.EncodeToString([]byte(":password"))
	c.Request, _ = http.NewRequest("GET", "/", nil)
	c.Request.Header.Set("Authorization", fmt.Sprintf("Basic %s", authValue))
	ba(c)
	assert.True(t, c.IsAborted())
	assert.Equal(t, http.StatusUnauthorized, c.Writer.Status())

	// Test empty password
	c, _ = CreateTestContext(httptest.NewRecorder())
	authValue = base64.StdEncoding.EncodeToString([]byte("admin:"))
	c.Request, _ = http.NewRequest("GET", "/", nil)
	c.Request.Header.Set("Authorization", fmt.Sprintf("Basic %s", authValue))
	ba(c)
	assert.True(t, c.IsAborted())
	assert.Equal(t, http.StatusUnauthorized, c.Writer.Status())
}

func TestBasicAuthPairs(t *testing.T) {
	accounts := Accounts{
		"admin": "password",
		"foo":   "bar",
	}
	pairs := processAccounts(accounts)
	user, found := pairs.searchCredential(authorizationHeader("admin", "password"))
	assert.True(t, found)
	assert.Equal(t, "admin", user)

	user, found = pairs.searchCredential(authorizationHeader("foo", "bar"))
	assert.True(t, found)
	assert.Equal(t, "foo", user)

	_, found = pairs.searchCredential(authorizationHeader("admin", ""))
	assert.False(t, found)

	_, found = pairs.searchCredential(authorizationHeader("", "password"))
	assert.False(t, found)

	_, found = pairs.searchCredential(authorizationHeader("", ""))
	assert.False(t, found)

	_, found = pairs.searchCredential("")
	assert.False(t, found)

	_, found = pairs.searchCredential("foo")
	assert.False(t, found)

	_, found = pairs.searchCredential("Basic ")
	assert.False(t, found)
}

func TestBasicAuthRealm(t *testing.T) {
	accounts := Accounts{
		"admin": "password",
		"foo":   "bar",
	}
	realm := "My Realm"
	ba := BasicAuthForRealm(accounts, realm)
	c, _ := CreateTestContext(httptest.NewRecorder())

	// Test correct credentials
	authValue := base64.StdEncoding.EncodeToString([]byte("admin:password"))
	c.Request, _ = http.NewRequest("GET", "/", nil)
	c.Request.Header.Set("Authorization", fmt.Sprintf("Basic %s", authValue))
	ba(c)
	assert.False(t, c.IsAborted())
	assert.Equal(t, "admin", c.MustGet(AuthUserKey).(string))

	// Test incorrect credentials
	c, _ = CreateTestContext(httptest.NewRecorder())
	authValue = base64.StdEncoding.EncodeToString([]byte("admin:password1"))
	c.Request, _ = http.NewRequest("GET", "/", nil)
	c.Request.Header.Set("Authorization", fmt.Sprintf("Basic %s", authValue))
	ba(c)
	assert.True(t, c.IsAborted())
	assert.Equal(t, http.StatusUnauthorized, c.Writer.Status())
	assert.Equal(t, fmt.Sprintf("Basic realm=%q", realm), c.Writer.Header().Get("WWW-Authenticate"))
}

func TestBasicAuthRealmEmpty(t *testing.T) {
	accounts := Accounts{
		"admin": "password",
		"foo":   "bar",
	}
	ba := BasicAuthForRealm(accounts, "")
	c, _ := CreateTestContext(httptest.NewRecorder())

	// Test correct credentials
	authValue := base64.StdEncoding.EncodeToString([]byte("admin:password"))
	c.Request, _ = http.NewRequest("GET", "/", nil)
	c.Request.Header.Set("Authorization", fmt.Sprintf("Basic %s", authValue))
	ba(c)
	assert.False(t, c.IsAborted())
	assert.Equal(t, "admin", c.MustGet(AuthUserKey).(string))

	// Test incorrect credentials
	c, _ = CreateTestContext(httptest.NewRecorder())
	authValue = base64.StdEncoding.EncodeToString([]byte("admin:password1"))
	c.Request, _ = http.NewRequest("GET", "/", nil)
	c.Request.Header.Set("Authorization", fmt.Sprintf("Basic %s", authValue))
	ba(c)
	assert.True(t, c.IsAborted())
	assert.Equal(t, http.StatusUnauthorized, c.Writer.Status())
	assert.Equal(t, `Basic realm="Authorization Required"`, c.Writer.Header().Get("WWW-Authenticate"))
}

func TestBasicAuthForProxy(t *testing.T) {
	accounts := Accounts{
		"admin": "password",
		"foo":   "bar",
	}
	ba := BasicAuthForProxy(accounts, "")
	c, _ := CreateTestContext(httptest.NewRecorder())

	// Test correct credentials
	authValue := base64.StdEncoding.EncodeToString([]byte("admin:password"))
	c.Request, _ = http.NewRequest("GET", "/", nil)
	c.Request.Header.Set("Proxy-Authorization", fmt.Sprintf("Basic %s", authValue))
	ba(c)
	assert.False(t, c.IsAborted())
	assert.Equal(t, "admin", c.MustGet(AuthProxyUserKey).(string))

	// Test incorrect credentials
	c, _ = CreateTestContext(httptest.NewRecorder())
	authValue = base64.StdEncoding.EncodeToString([]byte("admin:password1"))
	c.Request, _ = http.NewRequest("GET", "/", nil)
	c.Request.Header.Set("Proxy-Authorization", fmt.Sprintf("Basic %s", authValue))
	ba(c)
	assert.True(t, c.IsAborted())
	assert.Equal(t, http.StatusProxyAuthRequired, c.Writer.Status())
	assert.Equal(t, `Basic realm="Proxy Authorization Required"`, c.Writer.Header().Get("Proxy-Authenticate"))
}

func TestBasicAuthForProxyRealm(t *testing.T) {
	accounts := Accounts{
		"admin": "password",
		"foo":   "bar",
	}
	realm := "My Proxy Realm"
	ba := BasicAuthForProxy(accounts, realm)
	c, _ := CreateTestContext(httptest.NewRecorder())

	// Test correct credentials
	authValue := base64.StdEncoding.EncodeToString([]byte("admin:password"))
	c.Request, _ = http.NewRequest("GET", "/", nil)
	c.Request.Header.Set("Proxy-Authorization", fmt.Sprintf("Basic %s", authValue))
	ba(c)
	assert.False(t, c.IsAborted())
	assert.Equal(t, "admin", c.MustGet(AuthProxyUserKey).(string))

	// Test incorrect credentials
	c, _ = CreateTestContext(httptest.NewRecorder())
	authValue = base64.StdEncoding.EncodeToString([]byte("admin:password1"))
	c.Request, _ = http.NewRequest("GET", "/", nil)
	c.Request.Header.Set("Proxy-Authorization", fmt.Sprintf("Basic %s", authValue))
	ba(c)
	assert.True(t, c.IsAborted())
	assert.Equal(t, http.StatusProxyAuthRequired, c.Writer.Status())
	assert.Equal(t, fmt.Sprintf("Basic realm=%q", realm), c.Writer.Header().Get("Proxy-Authenticate"))
}

func TestBasicAuthEmptyAccounts(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Calling BasicAuth() with an empty list of accounts should cause panic")
		}
	}()
	_ = BasicAuth(Accounts{})
}

func TestBasicAuthUserEmpty(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Empty user in the list of accounts should cause panic")
		}
	}()
	accounts := Accounts{
		"": "password",
	}
	_ = BasicAuth(accounts)
}

func TestBasicAuthSecretEmpty(t *testing.T) {
	accounts := Accounts{
		"admin": "password",
		"foo":   "",
	}
	ba := BasicAuth(accounts)
	c, _ := CreateTestContext(httptest.NewRecorder())

	// Test correct credentials
	authValue := base64.StdEncoding.EncodeToString([]byte("admin:password"))
	c.Request, _ = http.NewRequest("GET", "/", nil)
	c.Request.Header.Set("Authorization", fmt.Sprintf("Basic %s", authValue))
	ba(c)
	assert.False(t, c.IsAborted())
	assert.Equal(t, "admin", c.MustGet(AuthUserKey).(string))

	c, _ = CreateTestContext(httptest.NewRecorder())
	authValue = base64.StdEncoding.EncodeToString([]byte("foo:"))
	c.Request, _ = http.NewRequest("GET", "/", nil)
	c.Request.Header.Set("Authorization", fmt.Sprintf("Basic %s", authValue))
	ba(c)
	assert.False(t, c.IsAborted())
	assert.Equal(t, "foo", c.MustGet(AuthUserKey).(string))
}

func TestBasicAuthUserAndSecretEmpty(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Empty user in the list of accounts should cause panic")
		}
	}()
	accounts := Accounts{
		"": "",
	}
	_ = BasicAuth(accounts)
}

func TestBasicAuthPairString(t *testing.T) {
	accounts := Accounts{
		"admin": "password",
		"foo":   "bar",
	}
	pairs := processAccounts(accounts)
	expected := "admin:password,foo:bar"
	assert.Equal(t, expected, pairs.String())
}

func TestBasicAuthPairStringSingle(t *testing.T) {
	accounts := Accounts{
		"admin": "password",
	}
	pairs := processAccounts(accounts)
	expected := "admin:password"
	assert.Equal(t, expected, pairs.String())
}

func TestBasicAuthPairStringEmpty(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Empty list of accounts should cause panic")
		}
	}()
	accounts := Accounts{}
	pairs := processAccounts(accounts)
	expected := ""
	assert.Equal(t, expected, pairs.String())
}

func TestBasicAuthPairUser(t *testing.T) {
	accounts := Accounts{
		"admin": "password",
		"foo":   "bar",
	}
	pairs := processAccounts(accounts)
	user := pairs.searchUser(authorizationHeader("admin", "password"))
	assert.Equal(t, "admin", user)

	user = pairs.searchUser(authorizationHeader("foo", "bar"))
	assert.Equal(t, "foo", user)

	user = pairs.searchUser(authorizationHeader("admin", ""))
	assert.Equal(t, "", user)

	user = pairs.searchUser(authorizationHeader("", "password"))
	assert.Equal(t, "", user)

	user = pairs.searchUser(authorizationHeader("", ""))
	assert.Equal(t, "", user)

	user = pairs.searchUser("")
	assert.Equal(t, "", user)

	user = pairs.searchUser("foo")
	assert.Equal(t, "", user)

	user = pairs.searchUser("Basic ")
	assert.Equal(t, "", user)
}

func TestBasicAuthEmptySecret(t *testing.T) {
	accounts := Accounts{
		"admin": "",
	}
	ba := BasicAuth(accounts)
	c, _ := CreateTestContext(httptest.NewRecorder())

	// Test correct credentials
	authValue := base64.StdEncoding.EncodeToString([]byte("admin:"))
	c.Request, _ = http.NewRequest("GET", "/", nil)
	c.Request.Header.Set("Authorization", fmt.Sprintf("Basic %s", authValue))
	ba(c)
	assert.False(t, c.IsAborted())
	assert.Equal(t, "admin", c.MustGet(AuthUserKey).(string))
}

func TestBasicAuthWithRealm(t *testing.T) {
	accounts := Accounts{
		"admin": "password",
	}
	ba := BasicAuthForRealm(accounts, "My Realm")
	w := httptest.NewRecorder()
	c, _ := CreateTestContext(w)

	// Test incorrect credentials
	authValue := base64.StdEncoding.EncodeToString([]byte("admin:password1"))
	c.Request, _ = http.NewRequest("GET", "/", nil)
	c.Request.Header.Set("Authorization", fmt.Sprintf("Basic %s", authValue))
	ba(c)
	require.True(t, c.IsAborted())
	require.Equal(t, http.StatusUnauthorized, c.Writer.Status())
	require.Equal(t, `Basic realm="My Realm"`, w.Header().Get("WWW-Authenticate"))
}

func TestBasicAuthWithEmptyRealm(t *testing.T) {
	accounts := Accounts{
		"admin": "password",
	}
	ba := BasicAuthForRealm(accounts, "")
	w := httptest.NewRecorder()
	c, _ := CreateTestContext(w)

	// Test incorrect credentials
	authValue := base64.StdEncoding.EncodeToString([]byte("admin:password1"))
	c.Request, _ = http.NewRequest("GET", "/", nil)
	c.Request.Header.Set("Authorization", fmt.Sprintf("Basic %s", authValue))
	ba(c)
	require.True(t, c.IsAborted())
	require.Equal(t, http.StatusUnauthorized, c.Writer.Status())
	require.Equal(t, `Basic realm="Authorization Required"`, w.Header().Get("WWW-Authenticate"))
}