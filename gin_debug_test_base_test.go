package gin

import (
	"bytes"
	"errors"
	"html/template"
	"io"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync/atomic"
	"testing"
)

var (
	ginMode          int32
	debugCode        int32 = 0
	DefaultWriter    io.Writer = os.Stdout
	DefaultErrorWriter io.Writer = os.Stderr
)

type HandlersChain []HandlerFunc

func (c HandlersChain) Last() HandlerFunc {
	if len(c) > 0 {
		return c[len(c)-1]
	}
	return nil
}

type HandlerFunc func()

func nameOfFunction(f HandlerFunc) string {
	return runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
}

func TestIsDebugging(t *testing.T) {
	atomic.StoreInt32(&ginMode, debugCode)
	if !IsDebugging() {
		t.Errorf("Expected debugging mode to be true")
	}

	atomic.StoreInt32(&ginMode, 1)
	if IsDebugging() {
		t.Errorf("Expected debugging mode to be false")
	}
}

func TestDebugPrintRoute(t *testing.T) {
	atomic.StoreInt32(&ginMode, debugCode)
	var buf bytes.Buffer
	DefaultWriter = &buf

	DebugPrintRouteFunc = nil

	handlers := HandlersChain{func() {}}
	debugPrintRoute("GET", "/test", handlers)

	expected := "[GIN-debug] GET    /test                     --> github.com/gin-gonic/gin.TestDebugPrintRoute.func1 (1 handlers)\n"
	if buf.String() != expected {
		t.Errorf("Expected %q but got %q", expected, buf.String())
	}
}

func TestDebugPrintLoadTemplate(t *testing.T) {
	atomic.StoreInt32(&ginMode, debugCode)
	var buf bytes.Buffer
	DefaultWriter = &buf

	tmpl := template.Must(template.New("test").Parse("{{define \"T1\"}}Hello{{end}}{{define \"T2\"}}World{{end}}"))
	debugPrintLoadTemplate(tmpl)

	expected := "[GIN-debug] Loaded HTML Templates (3): \n\t- T1\n\t- T2\n\t- test\n\n"
	if buf.String() != expected {
		t.Errorf("Expected %q but got %q", expected, buf.String())
	}
}

func TestDebugPrintWarningDefault(t *testing.T) {
	atomic.StoreInt32(&ginMode, debugCode)
	var buf bytes.Buffer
	DefaultWriter = &buf

	debugPrintWARNINGDefault()

	v, _ := getMinVer(runtime.Version())
	var expected string
	if v < ginSupportMinGoVer {
		expected = "[GIN-debug] [WARNING] Now Gin requires Go 1.22+.\n\n"
	}
	expected += "[GIN-debug] [WARNING] Creating an Engine instance with the Logger and Recovery middleware already attached.\n\n"

	if buf.String() != expected {
		t.Errorf("Expected %q but got %q", expected, buf.String())
	}
}

func TestDebugPrintWarningNew(t *testing.T) {
	atomic.StoreInt32(&ginMode, debugCode)
	var buf bytes.Buffer
	DefaultWriter = &buf

	debugPrintWARNINGNew()

	expected := `[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:	export GIN_MODE=release
 - using code:	gin.SetMode(gin.ReleaseMode)

`
	if buf.String() != expected {
		t.Errorf("Expected %q but got %q", expected, buf.String())
	}
}

func TestDebugPrintWarningSetHTMLTemplate(t *testing.T) {
	atomic.StoreInt32(&ginMode, debugCode)
	var buf bytes.Buffer
	DefaultWriter = &buf

	debugPrintWARNINGSetHTMLTemplate()

	expected := `[GIN-debug] [WARNING] Since SetHTMLTemplate() is NOT thread-safe. It should only be called
at initialization. ie. before any route is registered or the router is listening in a socket:

	router := gin.Default()
	router.SetHTMLTemplate(template) // << good place

`
	if buf.String() != expected {
		t.Errorf("Expected %q but got %q", expected, buf.String())
	}
}

func TestDebugPrintError(t *testing.T) {
	atomic.StoreInt32(&ginMode, debugCode)
	var buf bytes.Buffer
	DefaultErrorWriter = &buf

	err := errors.New("test error")
	debugPrintError(err)

	expected := "[GIN-debug] [ERROR] test error\n"
	if buf.String() != expected {
		t.Errorf("Expected %q but got %q", expected, buf.String())
	}
}

func TestGetMinVer(t *testing.T) {
	tests := []struct {
		version  string
		expected uint64
	}{
		{"go1.20.1", 20},
		{"go1.22", 22},
		{"go1.21.0", 21},
	}

	for _, test := range tests {
		result, err := getMinVer(test.version)
		if err != nil {
			t.Errorf("Unexpected error for version %s: %v", test.version, err)
		}
		if result != test.expected {
			t.Errorf("For version %s, expected %d but got %d", test.version, test.expected, result)
		}
	}
}

func TestGetMinVer_Error(t *testing.T) {
	_, err := getMinVer("invalid.version")
	if err == nil {
		t.Error("Expected error for invalid version string, but got none")
	}
}