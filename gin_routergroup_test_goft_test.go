// Copyright 2014 Manu Martinez-Almeida. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package gin

import (
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRouterGroupBasic(t *testing.T) {
	router := New()
	assert.Equal(t, "/", router.BasePath())
	emptyRouter := router.Group("")
	assert.Equal(t, "/", emptyRouter.BasePath())
	// equal to empty
	rootRouter := router.Group("/")
	assert.Equal(t, "/", rootRouter.BasePath())
	// equal to empty
	routerA := router.Group("/a")
	assert.Equal(t, "/a", routerA.BasePath())
	routerAB := routerA.Group("/b")
	assert.Equal(t, "/a/b", routerAB.BasePath())
	routerB := router.Group("/b")
	assert.Equal(t, "/b", routerB.BasePath())
	routerC := router.Group("/c")
	assert.Equal(t, "/c", routerC.BasePath())

	// re-use a same path
	routerC2 := router.Group("/c")
	assert.Equal(t, "/c", routerC2.BasePath())

	// re-use a same path
	routerC3 := routerC2.Group("/")
	assert.Equal(t, "/c", routerC3.BasePath())
}

func TestRouterGroupBasicHandle(t *testing.T) {
	router := New()
	emptyRouter := router.Group("")
	// equal to empty
	rootRouter := router.Group("/")
	routerA := router.Group("/a")
	routerAB := routerA.Group("/b")
	routerB := router.Group("/b")

	routerC := router.Group("/c")
	routerC.GET("/a", func(c *Context) {})
	routerC.GET("/a/b", func(c *Context) {})
	routerC2 := router.Group("/c")
	routerC2.GET("/a", func(c *Context) {})
	routerC2.GET("/a/b", func(c *Context) {})
	routerC3 := routerC2.Group("/")
	routerC3.GET("/a", func(c *Context) {})
	routerC3.GET("/a/b", func(c *Context) {})
	routerC3.GET("/", func(c *Context) {})

	emptyRouter.GET("/a", func(c *Context) {})
	emptyRouter.GET("/a/b", func(c *Context) {})
	rootRouter.GET("/a", func(c *Context) {})
	rootRouter.GET("/a/b", func(c *Context) {})
	routerA.GET("/a", func(c *Context) {})
	routerA.GET("/a/b", func(c *Context) {})
	routerAB.GET("/a", func(c *Context) {})
	routerAB.GET("/a/b", func(c *Context) {})
	routerB.GET("/a", func(c *Context) {})
	routerB.GET("/a/b", func(c *Context) {})
}

func TestRouterGroupBasicInvalidStatic(t *testing.T) {
	router := New()
	assert.Panics(t, func() {
		router.StaticFile("/:file", "/")
	})
	assert.Panics(t, func() {
		router.StaticFile("/*file", "/")
	})
	assert.Panics(t, func() {
		router.Static("/:file", "/")
	})
	assert.Panics(t, func() {
		router.Static("/*file", "/")
	})
	assert.Panics(t, func() {
		router.StaticFileFS("/:file", "/", Dir("/", false))
	})
	assert.Panics(t, func() {
		router.StaticFileFS("/*file", "/", Dir("/", false))
	})
	assert.Panics(t, func() {
		router.StaticFS("/:file", Dir("/", false))
	})
	assert.Panics(t, func() {
		router.StaticFS("/*file", Dir("/", false))
	})
}

func TestRouterGroupBasicStaticFile(t *testing.T) {
	router := New()
	router.StaticFile("/file", ".")

	router.StaticFileFS("/file2", ".", Dir(".", false))

	performRequest(router, "GET", "/file")
	performRequest(router, "HEAD", "/file")
	performRequest(router, "GET", "/file2")
	performRequest(router, "HEAD", "/file2")
}

func TestRouterGroupBasicStatic(t *testing.T) {
	router := New()
	router.Static("/folder", ".")

	router.StaticFS("/folder2", Dir(".", false))

	performRequest(router, "GET", "/folder/gin.go")
	performRequest(router, "HEAD", "/folder/gin.go")
	performRequest(router, "GET", "/folder2/gin.go")
	performRequest(router, "HEAD", "/folder2/gin.go")
}

func TestRouterGroupInvalidStatic(t *testing.T) {
	router := New()
	assert.Panics(t, func() {
		router.StaticFile("/:file", "/")
	})
	assert.Panics(t, func() {
		router.StaticFile("/*file", "/")
	})
	assert.Panics(t, func() {
		router.Static("/:file", "/")
	})
	assert.Panics(t, func() {
		router.Static("/*file", "/")
	})
	assert.Panics(t, func() {
		router.StaticFileFS("/:file", "/", Dir("/", false))
	})
	assert.Panics(t, func() {
		router.StaticFileFS("/*file", "/", Dir("/", false))
	})
	assert.Panics(t, func() {
		router.StaticFS("/:file", Dir("/", false))
	})
	assert.Panics(t, func() {
		router.StaticFS("/*file", Dir("/", false))
	})
}

func TestRouterGroupStaticListingDir(t *testing.T) {
	router := New()
	router.StaticFS("/fs", Dir("./", true))

	w := performRequest(router, "GET", "/fs")

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "gin.go")
}

func TestRouterGroupStaticNoListing(t *testing.T) {
	router := New()
	router.Static("/", ".")

	// print debug info
	w := performRequest(router, "GET", "/")

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.NotContains(t, w.Body.String(), "gin.go")
}

func TestRouterGroupStaticFSNoListing(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir(".", false))

	w := performRequest(router, "GET", "/")

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.NotContains(t, w.Body.String(), "gin.go")
}

func TestRouterGroupStaticFSNoListingWithIndex(t *testing.T) {
	router := New()
	router.StaticFS("/", Dir("./fixtures", false))

	w := performRequest(router, "GET", "/")

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "index.html")
}

func TestRouterGroupStaticFSNoListingEmptyDir(t *testing.T) {
	err := os.Mkdir("./fixtures/test", 0o755)
	assert.NoError(t, err)

	router := New()
	router.StaticFS("/fs", Dir("./fixtures/test", false))

	w := performRequest(router, "GET", "/fs")

	assert.Equal(t, http.StatusNotFound, w.Code)

	err = os.Remove("./fixtures/test")
	assert.NoError(t, err)
}

func TestRouterGroupStaticFSNoListingNonExisting(t *testing.T) {
	router := New()
	router.StaticFS("/fs", Dir("./fixtures/test", false))

	w := performRequest(router, "GET", "/fs")

	assert.Equal(t, http.StatusNotFound, w.Code)
}