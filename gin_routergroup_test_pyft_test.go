// Copyright 2014 Manu Martinez-Almeida. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package gin

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRouterGroupBasic(t *testing.T) {
	router := New()

	v1 := router.Group("/v1")
	assert.Equal(t, "/v1", v1.BasePath())
	assert.Equal(t, v1.Handlers, HandlersChain(nil))
	assert.Equal(t, v1.engine, router)

	v2 := v1.Group("/v2")
	assert.Equal(t, "/v1/v2", v2.BasePath())
	assert.Equal(t, v2.Handlers, HandlersChain(nil))
	assert.Equal(t, v2.engine, router)
}

func TestRouterGroupBasicHandle(t *testing.T) {
	router := New()

	router.Handle("GET", "/handler", func(c *Context) {
		c.String(http.StatusOK, "handler")
	})

	v1 := router.Group("/v1")
	v1.Handle("GET", "/handler", func(c *Context) {
		c.String(http.StatusOK, "v1 handler")
	})

	v2 := v1.Group("/v2")
	v2.Handle("GET", "/handler", func(c *Context) {
		c.String(http.StatusOK, "v2 handler")
	})

	w := performRequest(router, "GET", "/handler")
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "handler", w.Body.String())

	w = performRequest(router, "GET", "/v1/handler")
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "v1 handler", w.Body.String())

	w = performRequest(router, "GET", "/v1/v2/handler")
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "v2 handler", w.Body.String())
}

func TestRouterGroupBasicHandleInvalidMethod(t *testing.T) {
	router := New()
	assert.Panics(t, func() {
		router.Handle(" GET", "/handler", func(c *Context) {
			c.String(http.StatusOK, "handler")
		})
	})
	assert.Panics(t, func() {
		router.Handle("GET ", "/handler", func(c *Context) {
			c.String(http.StatusOK, "handler")
		})
	})
	assert.Panics(t, func() {
		router.Handle("", "/handler", func(c *Context) {
			c.String(http.StatusOK, "handler")
		})
	})
	assert.Panics(t, func() {
		router.Handle("PO ST", "/handler", func(c *Context) {
			c.String(http.StatusOK, "handler")
		})
	})
	assert.Panics(t, func() {
		router.Handle("1GET", "/handler", func(c *Context) {
			c.String(http.StatusOK, "handler")
		})
	})
	assert.Panics(t, func() {
		router.Handle("PATCh", "/handler", func(c *Context) {
			c.String(http.StatusOK, "handler")
		})
	})
}

func TestRouterGroupBasicHandleNoHandlers(t *testing.T) {
	router := New()
	router.Handle("POST", "/handler")
	router.Handle("DELETE", "/handler", nil)

	v1 := router.Group("/v1")
	v1.Handle("POST", "/handler")
	v1.Handle("DELETE", "/handler", nil)

	v2 := v1.Group("/v2")
	v2.Handle("POST", "/handler")
	v2.Handle("DELETE", "/handler", nil)
}

func TestRouterGroupBasicHandleHead(t *testing.T) {
	router := New()
	router.Handle("HEAD", "/handler", func(c *Context) {
		c.String(http.StatusOK, "handler")
	})

	v1 := router.Group("/v1")
	v1.Handle("HEAD", "/handler", func(c *Context) {
		c.String(http.StatusOK, "v1 handler")
	})

	v2 := v1.Group("/v2")
	v2.Handle("HEAD", "/handler", func(c *Context) {
		c.String(http.StatusOK, "v2 handler")
	})

	w := performRequest(router, "HEAD", "/handler")
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "handler", w.Body.String())

	w = performRequest(router, "HEAD", "/v1/handler")
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "v1 handler", w.Body.String())

	w = performRequest(router, "HEAD", "/v1/v2/handler")
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "v2 handler", w.Body.String())
}

func TestRouterGroupBasicHandleAny(t *testing.T) {
	router := New()

	router.Any("/handler", func(c *Context) {
		c.String(http.StatusOK, "handler")
	})

	v1 := router.Group("/v1")
	v1.Any("/handler", func(c *Context) {
		c.String(http.StatusOK, "v1 handler")
	})

	v2 := v1.Group("/v2")
	v2.Any("/handler", func(c *Context) {
		c.String(http.StatusOK, "v2 handler")
	})

	for _, method := range anyMethods {
		w := performRequest(router, method, "/handler")
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "handler", w.Body.String())

		w = performRequest(router, method, "/v1/handler")
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "v1 handler", w.Body.String())

		w = performRequest(router, method, "/v1/v2/handler")
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "v2 handler", w.Body.String())
	}
}

func TestRouterGroupBasicHandleStaticFile(t *testing.T) {
	router := New()

	router.StaticFile("/test", "./README.md")

	v1 := router.Group("/v1")
	v1.StaticFile("/test", "./LICENSE")

	v2 := v1.Group("/v2")
	v2.StaticFile("/test", "./HISTORY.md")

	w := performRequest(router, "GET", "/test")
	w2 := performRequest(router, "HEAD", "/test")
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, w.Header().Get("Content-Length"), w2.Header().Get("Content-Length"))
	assert.Equal(t, w.Header().Get("Content-Type"), w2.Header().Get("Content-Type"))

	w = performRequest(router, "GET", "/v1/test")
	w2 = performRequest(router, "HEAD", "/v1/test")
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, w.Header().Get("Content-Length"), w2.Header().Get("Content-Length"))
	assert.Equal(t, w.Header().Get("Content-Type"), w2.Header().Get("Content-Type"))

	w = performRequest(router, "GET", "/v1/v2/test")
	w2 = performRequest(router, "HEAD", "/v1/v2/test")
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, w.Header().Get("Content-Length"), w2.Header().Get("Content-Length"))
	assert.Equal(t, w.Header().Get("Content-Type"), w2.Header().Get("Content-Type"))
}

func TestRouterGroupBasicHandleStaticFileFS(t *testing.T) {
	router := New()

	router.StaticFileFS("/test", "./README.md", Dir{".", false})

	v1 := router.Group("/v1")
	v1.StaticFileFS("/test", "./LICENSE", Dir{".", false})

	v2 := v1.Group("/v2")
	v2.StaticFileFS("/test", "./HISTORY.md", Dir{".", false})

	w := performRequest(router, "GET", "/test")
	w2 := performRequest(router, "HEAD", "/test")
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, w.Header().Get("Content-Length"), w2.Header().Get("Content-Length"))
	assert.Equal(t, w.Header().Get("Content-Type"), w2.Header().Get("Content-Type"))

	w = performRequest(router, "GET", "/v1/test")
	w2 = performRequest(router, "HEAD", "/v1/test")
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, w.Header().Get("Content-Length"), w2.Header().Get("Content-Length"))
	assert.Equal(t, w.Header().Get("Content-Type"), w2.Header().Get("Content-Type"))

	w = performRequest(router, "GET", "/v1/v2/test")
	w2 = performRequest(router, "HEAD", "/v1/v2/test")
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, w.Header().Get("Content-Length"), w2.Header().Get("Content-Length"))
	assert.Equal(t, w.Header().Get("Content-Type"), w2.Header().Get("Content-Type"))
}

func TestRouterGroupBasicHandleStaticFileFSNoListing(t *testing.T) {
	router := New()

	router.StaticFileFS("/test", "./README.md", Dir{".", true})

	v1 := router.Group("/v1")
	v1.StaticFileFS("/test", "./LICENSE", Dir{".", true})

	v2 := v1.Group("/v2")
	v2.StaticFileFS("/test", "./HISTORY.md", Dir{".", true})

	w := performRequest(router, "GET", "/test")
	w2 := performRequest(router, "HEAD", "/test")
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, w.Header().Get("Content-Length"), w2.Header().Get("Content-Length"))
	assert.Equal(t, w.Header().Get("Content-Type"), w2.Header().Get("Content-Type"))

	w = performRequest(router, "GET", "/v1/test")
	w2 = performRequest(router, "HEAD", "/v1/test")
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, w.Header().Get("Content-Length"), w2.Header().Get("Content-Length"))
	assert.Equal(t, w.Header().Get("Content-Type"), w2.Header().Get("Content-Type"))

	w = performRequest(router, "GET", "/v1/v2/test")
	w2 = performRequest(router, "HEAD", "/v1/v2/test")
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, w.Header().Get("Content-Length"), w2.Header().Get("Content-Length"))
	assert.Equal(t, w.Header().Get("Content-Type"), w2.Header().Get("Content-Type"))
}

func TestRouterGroupBasicHandleStatic(t *testing.T) {
	router := New()

	router.Static("/test", ".")

	v1 := router.Group("/v1")
	v1.Static("/test", "..")

	v2 := v1.Group("/v2")
	v2.Static("/test", ".")

	w := performRequest(router, "GET", "/test/routergroup_test.go")
	w2 := performRequest(router, "HEAD", "/test/routergroup_test.go")
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, w.Header().Get("Content-Length"), w2.Header().Get("Content-Length"))
	assert.Equal(t, w.Header().Get("Content-Type"), w2.Header().Get("Content-Type"))

	w = performRequest(router, "GET", "/v1/test/gin")
	w2 = performRequest(router, "HEAD", "/v1/test/gin")
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, w.Header().Get("Content-Length"), w2.Header().Get("Content-Length"))
	assert.Equal(t, w.Header().Get("Content-Type"), w2.Header().Get("Content-Type"))

	w = performRequest(router, "GET", "/v1/v2/test/routergroup_test.go")
	w2 = performRequest(router, "HEAD", "/v1/v2/test/routergroup_test.go")
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, w.Header().Get("Content-Length"), w2.Header().Get("Content-Length"))
	assert.Equal(t, w.Header().Get("Content-Type"), w2.Header().Get("Content-Type"))
}

func TestRouterGroupBasicHandleStaticFS(t *testing.T) {
	router := New()

	router.StaticFS("/test", Dir{".", false})

	v1 := router.Group("/v1")
	v1.StaticFS("/test", Dir{"..", false})

	v2 := v1.Group("/v2")
	v2.StaticFS("/test", Dir{".", false})

	w := performRequest(router, "GET", "/test/routergroup_test.go")
	w2 := performRequest(router, "HEAD", "/test/routergroup_test.go")
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, w.Header().Get("Content-Length"), w2.Header().Get("Content-Length"))
	assert.Equal(t, w.Header().Get("Content-Type"), w2.Header().Get("Content-Type"))

	w = performRequest(router, "GET", "/v1/test/gin")
	w2 = performRequest(router, "HEAD", "/v1/test/gin")
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, w.Header().Get("Content-Length"), w2.Header().Get("Content-Length"))
	assert.Equal(t, w.Header().Get("Content-Type"), w2.Header().Get("Content-Type"))

	w = performRequest(router, "GET", "/v1/v2/test/routergroup_test.go")
	w2 = performRequest(router, "HEAD", "/v1/v2/test/routergroup_test.go")
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, w.Header().Get("Content-Length"), w2.Header().Get("Content-Length"))
	assert.Equal(t, w.Header().Get("Content-Type"), w2.Header().Get("Content-Type"))
}

func TestRouterGroupBasicHandleStaticFSNoListing(t *testing.T) {
	router := New()

	router.StaticFS("/test", Dir{".", true})

	v1 := router.Group("/v1")
	v1.StaticFS("/test", Dir{"..", true})

	v2 := v1.Group("/v2")
	v2.StaticFS("/test", Dir{".", true})

	w := performRequest(router, "GET", "/test/routergroup_test.go")
	w2 := performRequest(router, "HEAD", "/test/routergroup_test.go")
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, w.Header().Get("Content-Length"), w2.Header().Get("Content-Length"))
	assert.Equal(t, w.Header().Get("Content-Type"), w2.Header().Get("Content-Type"))

	w = performRequest(router, "GET", "/v1/test/gin")
	w2 = performRequest(router, "HEAD", "/v1/test/gin")
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, w.Header().Get("Content-Length"), w2.Header().Get("Content-Length"))
	assert.Equal(t, w.Header().Get("Content-Type"), w2.Header().Get("Content-Type"))

	w = performRequest(router, "GET", "/v1/v2/test/routergroup_test.go")
	w2 = performRequest(router, "HEAD", "/v1/v2/test/routergroup_test.go")
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, w.Header().Get("Content-Length"), w2.Header().Get("Content-Length"))
	assert.Equal(t, w.Header().Get("Content-Type"), w2.Header().Get("Content-Type"))
}

func TestRouterGroupStaticListing(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/")
	assert.Contains(t, w.Body.String(), "gin.go")
	assert.Contains(t, w.Body.String(), "gin_test.go")
	assert.Contains(t, w.Body.String(), "routergroup.go")
	assert.Contains(t, w.Body.String(), "routergroup_test.go")
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRouterGroupStaticNoListing(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", true})

	w := performRequest(router, "GET", "/")
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNoListing(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", true})

	w := performRequest(router, "GET", "/")
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSListing(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/")
	assert.Contains(t, w.Body.String(), "gin.go")
	assert.Contains(t, w.Body.String(), "gin_test.go")
	assert.Contains(t, w.Body.String(), "routergroup.go")
	assert.Contains(t, w.Body.String(), "routergroup_test.go")
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRouterGroupStaticCustomFSNoListing(t *testing.T) {
	router := New()
	router.StaticFS("/", http.Dir("./"))

	w := performRequest(router, "GET", "/")
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticCustomFSListing(t *testing.T) {
	router := New()
	router.StaticFS("/", http.Dir("./"))

	w := performRequest(router, "GET", "/")
	assert.Contains(t, w.Body.String(), "gin.go")
	assert.Contains(t, w.Body.String(), "gin_test.go")
	assert.Contains(t, w.Body.String(), "routergroup.go")
	assert.Contains(t, w.Body.String(), "routergroup_test.go")
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRouterGroupStaticFSNotFound(t *testing.T) {
	router := New()
	router.NoRoute(func(c *Context) {
		c.String(http.StatusNotFound, "custom not found")
	})
	router.StaticFS("/", Dir{"./", true})

	w := performRequest(router, "GET", "/nonexistent")
	assert.Equal(t, "custom not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundServingIndex(t *testing.T) {
	router := New()
	router.NoRoute(func(c *Context) {
		c.String(http.StatusNotFound, "custom not found")
	})
	router.StaticFS("/", Dir{"./", true})

	w := performRequest(router, "GET", "/nonexistent/")
	assert.Equal(t, "custom not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandler(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", true})

	w := performRequest(router, "GET", "/nonexistent/")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithoutServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", true})

	w := performRequest(router, "GET", "/nonexistent")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent/")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithoutServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent/")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithoutServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent/")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithoutServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent/")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithoutServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent/")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithoutServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent/")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithoutServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent/")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithoutServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent/")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithoutServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent/")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithoutServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent/")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithoutServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent/")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithoutServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent/")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithoutServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent/")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithoutServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent/")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithoutServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent/")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithoutServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent/")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithoutServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent/")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithoutServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent/")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithoutServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent/")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithoutServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent/")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithoutServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent/")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithoutServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent/")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithoutServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent/")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithoutServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent/")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithoutServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent/")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithoutServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent/")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithoutServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent/")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithoutServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent/")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithoutServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent/")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithoutServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent/")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithoutServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent/")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithoutServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent/")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithoutServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent/")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithoutServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent/")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithoutServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent/")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithoutServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent/")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithoutServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent/")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithoutServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent/")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithoutServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent/")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithoutServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent/")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithoutServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent/")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithoutServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent/")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithoutServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent/")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithoutServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent/")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithoutServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent/")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithoutServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent/")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithoutServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent/")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithoutServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent/")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithoutServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent/")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithoutServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent/")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithoutServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent/")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithoutServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent/")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithoutServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent/")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithoutServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent/")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithoutServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent/")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithoutServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent/")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithoutServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent/")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithoutServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent/")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithoutServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent/")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithoutServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent/")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithoutServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent/")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithoutServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent/")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithoutServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent/")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithoutServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent/")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithoutServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent/")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithoutServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent/")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithoutServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent/")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithoutServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent/")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithoutServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent/")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithoutServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent/")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithoutServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent/")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithoutServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent/")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithoutServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent/")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithoutServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent/")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithoutServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent/")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithoutServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent/")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithoutServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent/")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithoutServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent/")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithoutServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent/")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithoutServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w := performRequest(router, "GET", "/nonexistent/")
	assert.Equal(t, "404 page not found", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterGroupStaticFSNotFoundWithoutNotFoundHandlerWithoutServingIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir{"./", false})

	w