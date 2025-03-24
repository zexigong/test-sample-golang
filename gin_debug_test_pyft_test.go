// Copyright 2014 Manu Martinez-Almeida. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package gin

import (
	"bytes"
	"html/template"
	"os"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// test function to mock runtime.Version() and return expected version strings
func mockRuntimeVersion(version string) func() {
	oldRuntimeVersion := runtimeVersion
	runtimeVersion = func() string {
		return version
	}
	return func() { runtimeVersion = oldRuntimeVersion }
}

func TestDebugPrint(t *testing.T) {
	resetDebugPrint()

	var w bytes.Buffer
	setup(&w)

	SetMode(DebugMode)
	debugPrint("DEBUG this!")
	assert.Equal(t, "[GIN-debug] DEBUG this!\n", w.String())
	resetDebugPrint()
}

func TestDebugPrintError(t *testing.T) {
	resetDebugPrint()

	var w bytes.Buffer
	setup(&w)

	SetMode(DebugMode)
	debugPrintError(nil)
	assert.Equal(t, "", w.String())

	debugPrintError(assert.AnError)
	assert.Equal(t, "[GIN-debug] [ERROR] assert.AnError general error for testing\n", w.String())
	resetDebugPrint()
}

func TestDebugPrintRoutes(t *testing.T) {
	resetDebugPrint()

	var w bytes.Buffer
	setup(&w)

	SetMode(DebugMode)
	debugPrintRoute("GET", "/path/to/route/:param", HandlersChain{func(c *Context) {}})
	assert.Equal(t, "[GIN-debug] GET    /path/to/route/:param --> github.com/gin-gonic/gin.TestDebugPrintRoutes.func1 (1 handlers)\n", w.String())
	resetDebugPrint()
}

func TestDebugPrintLoadTemplate(t *testing.T) {
	resetDebugPrint()

	var w bytes.Buffer
	setup(&w)

	SetMode(DebugMode)
	tmpl, err := template.New("t").Parse(`{{define "T"}}Hello, {{.}}!{{end}}`)
	assert.NoError(t, err)

	debugPrintLoadTemplate(tmpl)
	assert.Equal(t, "[GIN-debug] Loaded HTML Templates (1): \n\t- t\n\n", w.String())
	resetDebugPrint()
}

func TestDebugPrintWARNINGSetHTMLTemplate(t *testing.T) {
	resetDebugPrint()

	var w bytes.Buffer
	setup(&w)

	SetMode(DebugMode)
	debugPrintWARNINGSetHTMLTemplate()
	assert.Equal(t, "[GIN-debug] [WARNING] Since SetHTMLTemplate() is NOT thread-safe. It should only be called\nat initialization. ie. before any route is registered or the router is listening in a socket:\n\n	router := gin.Default()\n	router.SetHTMLTemplate(template) // << good place\n\n", w.String())
	resetDebugPrint()
}

func TestDebugPrintWARNINGNew(t *testing.T) {
	resetDebugPrint()

	var w bytes.Buffer
	setup(&w)

	SetMode(DebugMode)
	debugPrintWARNINGNew()
	assert.Equal(t, "[GIN-debug] [WARNING] Running in \"debug\" mode. Switch to \"release\" mode in production.\n - using env:\texport GIN_MODE=release\n - using code:\tgin.SetMode(gin.ReleaseMode)\n\n", w.String())
	resetDebugPrint()
}

func TestDebugPrintWARNINGDefault(t *testing.T) {
	resetDebugPrint()

	var w bytes.Buffer
	setup(&w)

	SetMode(DebugMode)
	debugPrintWARNINGDefault()
	assert.Equal(t, "[GIN-debug] [WARNING] Creating an Engine instance with the Logger and Recovery middleware already attached.\n\n", w.String())
	resetDebugPrint()
}

func TestDebugPrintWARNINGDefault_GoVer(t *testing.T) {
	resetDebugPrint()

	// mock runtime.Version() to return go1.20
	defer mockRuntimeVersion("go1.20")()

	var w bytes.Buffer
	setup(&w)

	SetMode(DebugMode)
	debugPrintWARNINGDefault()
	assert.Equal(t, "[GIN-debug] [WARNING] Now Gin requires Go 1.22+.\n\n[GIN-debug] [WARNING] Creating an Engine instance with the Logger and Recovery middleware already attached.\n\n", w.String())
	resetDebugPrint()
}

func TestDebugPrintWARNINGDefault_GoVerNoPrint(t *testing.T) {
	resetDebugPrint()

	// mock runtime.Version() to return go1.22
	defer mockRuntimeVersion("go1.22")()

	var w bytes.Buffer
	setup(&w)

	SetMode(DebugMode)
	debugPrintWARNINGDefault()
	assert.Equal(t, "[GIN-debug] [WARNING] Creating an Engine instance with the Logger and Recovery middleware already attached.\n\n", w.String())
	resetDebugPrint()
}

func TestDebugPrintWARNINGDefault_GoVerDot(t *testing.T) {
	resetDebugPrint()

	// mock runtime.Version() to return go1.20.1
	defer mockRuntimeVersion("go1.20.1")()

	var w bytes.Buffer
	setup(&w)

	SetMode(DebugMode)
	debugPrintWARNINGDefault()
	assert.Equal(t, "[GIN-debug] [WARNING] Now Gin requires Go 1.22+.\n\n[GIN-debug] [WARNING] Creating an Engine instance with the Logger and Recovery middleware already attached.\n\n", w.String())
	resetDebugPrint()
}

func setup(w *bytes.Buffer) {
	ginMode = debugCode
	// default Writer set to os.Stdout
	DefaultWriter = w
	DefaultErrorWriter = w
}

func resetDebugPrint() {
	ginMode = testCode
	DefaultWriter = os.Stdout
	DefaultErrorWriter = os.Stderr
}

func TestDebugPrintFunc(t *testing.T) {
	resetDebugPrint()

	var w bytes.Buffer
	setup(&w)

	SetMode(DebugMode)
	DebugPrintFunc = func(format string, values ...interface{}) {
		w.WriteString("[TEST-debug] " + format + "\n")
	}
	debugPrint("DEBUG this!")
	assert.Equal(t, "[TEST-debug] DEBUG this!\n", w.String())

	resetDebugPrint()
}