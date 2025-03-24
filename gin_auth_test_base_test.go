package gin

import (
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestBasicAuth(t *testing.T) {
	accounts := Accounts{
		"admin": "password",
		"user":  "secret",
	}

	router := gin.New()
	router.Use(BasicAuth(accounts))
	router.GET("/protected", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello, %s", c.MustGet(AuthUserKey))
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/protected", nil)

	// No credentials
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Equal(t, "Basic realm=\"Authorization Required\"", w.Header().Get("WWW-Authenticate"))

	// Invalid credentials
	req.SetBasicAuth("admin", "wrongpassword")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	// Valid credentials
	req.SetBasicAuth("admin", "password")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "Hello, admin", w.Body.String())
}

func TestBasicAuthForRealm(t *testing.T) {
	accounts := Accounts{
		"admin": "password",
		"user":  "secret",
	}

	router := gin.New()
	router.Use(BasicAuthForRealm(accounts, "My Realm"))
	router.GET("/protected", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello, %s", c.MustGet(AuthUserKey))
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/protected", nil)

	// No credentials
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Equal(t, "Basic realm=\"My Realm\"", w.Header().Get("WWW-Authenticate"))

	// Invalid credentials
	req.SetBasicAuth("user", "wrongsecret")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	// Valid credentials
	req.SetBasicAuth("user", "secret")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "Hello, user", w.Body.String())
}

func TestBasicAuthForProxy(t *testing.T) {
	accounts := Accounts{
		"proxyadmin": "proxypass",
		"proxyuser":  "proxysecret",
	}

	router := gin.New()
	router.Use(BasicAuthForProxy(accounts, "Proxy Realm"))
	router.GET("/proxy-protected", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello, %s", c.MustGet(AuthProxyUserKey))
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/proxy-protected", nil)

	// No credentials
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusProxyAuthRequired, w.Code)
	assert.Equal(t, "Basic realm=\"Proxy Realm\"", w.Header().Get("Proxy-Authenticate"))

	// Invalid credentials
	auth := "Basic " + base64.StdEncoding.EncodeToString([]byte("proxyuser:wrongsecret"))
	req.Header.Set("Proxy-Authorization", auth)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusProxyAuthRequired, w.Code)

	// Valid credentials
	auth = "Basic " + base64.StdEncoding.EncodeToString([]byte("proxyuser:proxysecret"))
	req.Header.Set("Proxy-Authorization", auth)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "Hello, proxyuser", w.Body.String())
}