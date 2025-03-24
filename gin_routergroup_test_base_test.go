package gin

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func createTestContext() *Context {
	w := httptest.NewRecorder()
	return &Context{
		Writer:  w,
		Request: &http.Request{},
	}
}

func createTestEngine() *Engine {
	return &Engine{
		routes: make(map[string]map[string]HandlersChain),
	}
}

func TestRouterGroupUse(t *testing.T) {
	engine := createTestEngine()
	group := &RouterGroup{engine: engine}

	middleware := func(c *Context) {}
	group.Use(middleware)

	assert.Len(t, group.Handlers, 1)
	assert.Equal(t, middleware, group.Handlers[0])
}

func TestRouterGroupGroup(t *testing.T) {
	engine := createTestEngine()
	group := &RouterGroup{engine: engine, basePath: "/api"}

	subGroup := group.Group("/v1")

	assert.Equal(t, "/api/v1", subGroup.basePath)
	assert.Equal(t, engine, subGroup.engine)
	assert.Empty(t, subGroup.Handlers)
}

func TestRouterGroupHandle(t *testing.T) {
	engine := createTestEngine()
	group := &RouterGroup{engine: engine}

	handler := func(c *Context) {}
	group.Handle(http.MethodGet, "/test", handler)

	handlers, ok := engine.routes[http.MethodGet]["/test"]
	assert.True(t, ok)
	assert.Len(t, handlers, 1)
	assert.Equal(t, handler, handlers[0])
}

func TestRouterGroupPOST(t *testing.T) {
	engine := createTestEngine()
	group := &RouterGroup{engine: engine}

	handler := func(c *Context) {}
	group.POST("/post", handler)

	handlers, ok := engine.routes[http.MethodPost]["/post"]
	assert.True(t, ok)
	assert.Len(t, handlers, 1)
	assert.Equal(t, handler, handlers[0])
}

func TestRouterGroupGET(t *testing.T) {
	engine := createTestEngine()
	group := &RouterGroup{engine: engine}

	handler := func(c *Context) {}
	group.GET("/get", handler)

	handlers, ok := engine.routes[http.MethodGet]["/get"]
	assert.True(t, ok)
	assert.Len(t, handlers, 1)
	assert.Equal(t, handler, handlers[0])
}

func TestRouterGroupAny(t *testing.T) {
	engine := createTestEngine()
	group := &RouterGroup{engine: engine}

	handler := func(c *Context) {}
	group.Any("/any", handler)

	for _, method := range anyMethods {
		handlers, ok := engine.routes[method]["/any"]
		assert.True(t, ok)
		assert.Len(t, handlers, 1)
		assert.Equal(t, handler, handlers[0])
	}
}

func TestRouterGroupStaticFile(t *testing.T) {
	engine := createTestEngine()
	group := &RouterGroup{engine: engine}

	group.StaticFile("/favicon.ico", "./resources/favicon.ico")

	getHandlers, ok := engine.routes[http.MethodGet]["/favicon.ico"]
	assert.True(t, ok)
	assert.NotEmpty(t, getHandlers)

	headHandlers, ok := engine.routes[http.MethodHead]["/favicon.ico"]
	assert.True(t, ok)
	assert.NotEmpty(t, headHandlers)
}

func TestRouterGroupStaticFS(t *testing.T) {
	engine := createTestEngine()
	group := &RouterGroup{engine: engine}
	fs := http.Dir(".")

	group.StaticFS("/static", fs)

	getHandlers, ok := engine.routes[http.MethodGet]["/static/*filepath"]
	assert.True(t, ok)
	assert.NotEmpty(t, getHandlers)

	headHandlers, ok := engine.routes[http.MethodHead]["/static/*filepath"]
	assert.True(t, ok)
	assert.NotEmpty(t, headHandlers)
}