// Copyright 2014 Manu Martinez-Almeida. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package gin

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
	SetMode(TestMode)
}

func TestDebugPrint(t *testing.T) {
	var w bytes.Buffer
	setup(&w)

	// debugPrint error
	debugPrintError(nil)
	assert.Equal(t, "", w.String())

	err := fmt.Errorf("this is an error")
	debugPrintError(err)
	assert.Equal(t, "[GIN-debug] [ERROR] this is an error\n", w.String())

	// debugPrint WARNING
	w.Reset()
	debugPrintWARNINGNew()
	assert.Equal(t, "[GIN-debug] [WARNING] Running in \"debug\" mode. Switch to \"release\" mode in production.\n - using env:\texport GIN_MODE=release\n - using code:\tgin.SetMode(gin.ReleaseMode)\n\n", w.String())

	// debugPrint WARNING SetHTMLTemplate
	w.Reset()
	debugPrintWARNINGSetHTMLTemplate()
	assert.Equal(t, "[GIN-debug] [WARNING] Since SetHTMLTemplate() is NOT thread-safe. It should only be called\nat initialization. ie. before any route is registered or the router is listening in a socket:\n\n\trouter := gin.Default()\n\trouter.SetHTMLTemplate(template) // << good place\n\n", w.String())

	// debugPrint WARNING default
	w.Reset()
	debugPrintWARNINGDefault()
	assert.Equal(t, "[GIN-debug] [WARNING] Creating an Engine instance with the Logger and Recovery middleware already attached.\n\n", w.String())

	// debugPrintRoute
	w.Reset()
	debugPrintRoute("GET", "/path/to/route", HandlersChain{func(c *Context) { c.String(200, "hello") }})
	assert.Equal(t, "[GIN-debug] GET    /path/to/route            --> github.com/gin-gonic/gin.TestDebugPrint.func6 (1 handlers)\n", w.String())

	// debugPrintLoadTemplate
	w.Reset()
	templ, err := template.New("test").Delims("{[{", "}]}").Parse(`Hello [{[{.}]}`);
	assert.NoError(t, err)
	debugPrintLoadTemplate(templ)
	assert.Equal(t, "[GIN-debug] Loaded HTML Templates (1): \n\t- test\n\n\n", w.String())
}

func TestDebugPrintRouteFunc(t *testing.T) {
	var w bytes.Buffer
	setup(&w)

	DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
		fmt.Fprintf(DefaultWriter, "[GIN-debug] Custom route %v, %v, %v, %v\n", httpMethod, absolutePath, handlerName, nuHandlers)
	}

	debugPrintRoute("GET", "/path/to/route", HandlersChain{func(c *Context) { c.String(200, "hello") }})
	assert.Equal(t, "[GIN-debug] Custom route GET, /path/to/route, github.com/gin-gonic/gin.TestDebugPrintRouteFunc.func2, 1\n", w.String())
}

func TestDebugPrintFunc(t *testing.T) {
	var w bytes.Buffer
	setup(&w)

	DebugPrintFunc = func(format string, values ...interface{}) {
		fmt.Fprintf(DefaultWriter, "[GIN-debug] Custom log %v, %v", format, values)
	}

	debugPrint("testPrint", "test")
	assert.Equal(t, "[GIN-debug] Custom log testPrint, [test]", w.String())
}

func setup(w io.Writer) {
	DefaultWriter = w
	DefaultErrorWriter = w
	SetMode(DebugMode)
}