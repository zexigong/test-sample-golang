package gin

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParams_Get(t *testing.T) {
	params := Params{
		{Key: "foo", Value: "bar"},
		{Key: "baz", Value: "qux"},
	}

	value, found := params.Get("foo")
	assert.True(t, found)
	assert.Equal(t, "bar", value)

	value, found = params.Get("baz")
	assert.True(t, found)
	assert.Equal(t, "qux", value)

	value, found = params.Get("notfound")
	assert.False(t, found)
	assert.Equal(t, "", value)
}

func TestParams_ByName(t *testing.T) {
	params := Params{
		{Key: "foo", Value: "bar"},
		{Key: "baz", Value: "qux"},
	}

	value := params.ByName("foo")
	assert.Equal(t, "bar", value)

	value = params.ByName("baz")
	assert.Equal(t, "qux", value)

	value = params.ByName("notfound")
	assert.Equal(t, "", value)
}

func TestMethodTrees_Get(t *testing.T) {
	node1 := &node{path: "/path1"}
	node2 := &node{path: "/path2"}
	trees := methodTrees{
		{method: "GET", root: node1},
		{method: "POST", root: node2},
	}

	assert.Equal(t, node1, trees.get("GET"))
	assert.Equal(t, node2, trees.get("POST"))
	assert.Nil(t, trees.get("PUT"))
}

func TestLongestCommonPrefix(t *testing.T) {
	assert.Equal(t, 0, longestCommonPrefix("foo", "bar"))
	assert.Equal(t, 3, longestCommonPrefix("foobar", "foo"))
	assert.Equal(t, 4, longestCommonPrefix("test", "test"))
}

func TestCountParams(t *testing.T) {
	assert.Equal(t, uint16(0), countParams("/static/path"))
	assert.Equal(t, uint16(1), countParams("/:param/path"))
	assert.Equal(t, uint16(2), countParams("/:param1/:param2"))
}

func TestCountSections(t *testing.T) {
	assert.Equal(t, uint16(3), countSections("/section1/section2/section3"))
	assert.Equal(t, uint16(1), countSections("/"))
	assert.Equal(t, uint16(0), countSections(""))
}

func TestNode_AddChild(t *testing.T) {
	root := &node{}
	child1 := &node{path: "/child1"}
	child2 := &node{path: "/child2", wildChild: true}

	root.addChild(child1)
	root.addChild(child2)

	assert.Equal(t, 2, len(root.children))
	assert.Equal(t, child1, root.children[0])
	assert.Equal(t, child2, root.children[1])
}

func TestNode_IncrementChildPrio(t *testing.T) {
	root := &node{}
	child1 := &node{priority: 1}
	child2 := &node{priority: 2}

	root.children = []*node{child1, child2}
	root.indices = "ab"

	pos := root.incrementChildPrio(0)
	assert.Equal(t, 0, pos)
	assert.Equal(t, uint32(2), child1.priority)

	pos = root.incrementChildPrio(1)
	assert.Equal(t, 0, pos)
	assert.Equal(t, uint32(3), child2.priority)
	assert.Equal(t, "ba", root.indices)
}

func TestNode_AddRoute(t *testing.T) {
	root := &node{}
	handlers := HandlersChain{}

	root.addRoute("/test", handlers)
	assert.Equal(t, "/test", root.path)
	assert.Equal(t, handlers, root.handlers)
}

func TestNode_GetValue(t *testing.T) {
	root := &node{}
	handlers := HandlersChain{func(c *Context) {}}
	root.addRoute("/test/:param", handlers)

	params := Params{}
	value := root.getValue("/test/value", &params, nil, true)

	assert.Equal(t, handlers, value.handlers)
	assert.Equal(t, "value", params.ByName("param"))
}

func TestNode_FindCaseInsensitivePath(t *testing.T) {
	root := &node{}
	handlers := HandlersChain{func(c *Context) {}}
	root.addRoute("/test/path", handlers)

	ciPath, found := root.findCaseInsensitivePath("/TEST/PATH", true)
	assert.True(t, found)
	assert.Equal(t, "/test/path", string(ciPath))
}