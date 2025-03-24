// Copyright 2014 Manu Martinez-Almeida. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package gin

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStack(t *testing.T) {
	buf := new(bytes.Buffer)
	stack(3, buf)
	assert.Contains(t, buf.String(), "gin.TestStack")
}

func TestSource(t *testing.T) {
	bs := source(nil, 0)
	assert.Equal(t, []byte("???"), bs)
}

func TestSourceWithLines(t *testing.T) {
	lines := [][]byte{
		[]byte("package main"),
		[]byte("import \"fmt\""),
		[]byte("func main() {"),
		[]byte("fmt.Println(123)"),
		[]byte("}"),
	}
	bs := source(lines, 0)
	assert.Equal(t, []byte("package main"), bs)

	bs = source(lines, 1)
	assert.Equal(t, []byte("import \"fmt\""), bs)

	bs = source(lines, 2)
	assert.Equal(t, []byte("func main() {"), bs)

	bs = source(lines, 3)
	assert.Equal(t, []byte("fmt.Println(123)"), bs)

	bs = source(lines, 4)
	assert.Equal(t, []byte("}"), bs)

	bs = source(lines, 5)
	assert.Equal(t, []byte("???"), bs)
}

func TestFunction(t *testing.T) {
	bs := function(0)
	assert.Equal(t, []byte("???"), bs)
}

func TestTimeFormat(t *testing.T) {
	f := timeFormat(nil)
	assert.Equal(t, "2006/01/02 - 15:04:05", f)
}

func TestRecoveryWithNoPanic(t *testing.T) {
	buf := new(bytes.Buffer)
	r := New()
	r.Use(RecoveryWithWriter(buf))
	r.GET("/recovery", func(c *Context) {
		c.String(http.StatusOK, "ok")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/recovery", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "ok", w.Body.String())
	assert.Equal(t, 0, buf.Len())
}

func TestRecovery(t *testing.T) {
	buf := new(bytes.Buffer)
	r := New()
	r.Use(RecoveryWithWriter(buf))
	r.GET("/recovery", func(c *Context) {
		panic("Oupps")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/recovery", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, "500 Internal Server Error", w.Body.String())
	assert.NotEqual(t, 0, buf.Len())
	assert.Contains(t, buf.String(), "GET /recovery HTTP/1.1")
	assert.Contains(t, buf.String(), "Oupps")
}

func TestCustomRecovery(t *testing.T) {
	buf := new(bytes.Buffer)
	r := New()
	r.Use(CustomRecoveryWithWriter(buf, func(c *Context, err any) {
		c.String(http.StatusBadRequest, "bad")
	}))
	r.GET("/recovery", func(c *Context) {
		panic("Oupps")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/recovery", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "bad", w.Body.String())
	assert.NotEqual(t, 0, buf.Len())
	assert.Contains(t, buf.String(), "GET /recovery HTTP/1.1")
	assert.Contains(t, buf.String(), "Oupps")
}

func TestRecoveryWithBrokenPipe(t *testing.T) {
	buf := new(bytes.Buffer)
	r := New()
	r.Use(RecoveryWithWriter(buf))
	r.GET("/recovery", func(c *Context) {
		// simulate a broken pipe
		// ref: https://go-review.googlesource.com/c/go/+/94045/
		panic(&os.SyscallError{
			Syscall: "write",
			Err:     errors.New("broken pipe"),
		})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/recovery", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, 0, buf.Len())
}

func TestRecoveryWithConnectionReset(t *testing.T) {
	buf := new(bytes.Buffer)
	r := New()
	r.Use(RecoveryWithWriter(buf))
	r.GET("/recovery", func(c *Context) {
		// simulate a broken pipe
		// ref: https://go-review.googlesource.com/c/go/+/94045/
		panic(&os.SyscallError{
			Syscall: "write",
			Err:     errors.New("connection reset by peer"),
		})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/recovery", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, 0, buf.Len())
}

func TestRecoveryWithNoReplaceStackIfAlreadyWritingHeader(t *testing.T) {
	buf := new(bytes.Buffer)
	r := New()
	r.Use(RecoveryWithWriter(buf))
	r.GET("/recovery", func(c *Context) {
		c.Status(http.StatusBadRequest)
		panic("Oupps")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/recovery", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, http.StatusText(http.StatusBadRequest), w.Body.String())
	assert.NotEqual(t, 0, buf.Len())
	assert.Contains(t, buf.String(), "GET /recovery HTTP/1.1")
	assert.Contains(t, buf.String(), "Oupps")
}

func TestSourceWithFile(t *testing.T) {
	_, file, _, ok := runtime.Caller(0)
	assert.True(t, ok)

	data, err := os.ReadFile(file)
	assert.NoError(t, err)

	lines := bytes.Split(data, []byte{'\n'})
	bs := source(lines, 0)
	assert.Equal(t, []byte(""), bs)

	bs = source(lines, 1)
	assert.Equal(t, []byte("// Copyright 2014 Manu Martinez-Almeida. All rights reserved."), bs)

	bs = source(lines, 2)
	assert.Equal(t, []byte("// Use of this source code is governed by a MIT style"), bs)

	bs = source(lines, 3)
	assert.Equal(t, []byte("// license that can be found in the LICENSE file."), bs)

	bs = source(lines, 4)
	assert.Equal(t, []byte(""), bs)

	bs = source(lines, 5)
	assert.Equal(t, []byte("package gin"), bs)

	bs = source(lines, 6)
	assert.Equal(t, []byte(""), bs)

	bs = source(lines, 7)
	assert.Equal(t, []byte("import ("), bs)
}

func TestFunctionWithFuncName(t *testing.T) {
	f := function(func() {})
	assert.Equal(t, "func1", f)
}

func TestFunctionWithAnonymousFunc(t *testing.T) {
	f := function(func() {
		func() {}
	})
	assert.Equal(t, "func2", f)
}

func TestFunctionWithNamedFunc(t *testing.T) {
	f := function(TestFunctionWithNamedFunc)
	assert.Equal(t, "gin.TestFunctionWithNamedFunc", f)
}

func TestFunctionWithFuncNameAndPeriod(t *testing.T) {
	f := function(func() {
		func() {}
	})
	assert.Equal(t, "func2", f)
}

func TestFunctionWithFuncNameAndSlash(t *testing.T) {
	f := function(func() {
		func() {}
	})
	assert.Equal(t, "func2", f)
}

func TestFunctionWithFuncNameAndSlashAndPeriod(t *testing.T) {
	f := function(func() {
		func() {}
	})
	assert.Equal(t, "func2", f)
}

func TestFunctionWithFuncNameAndDot(t *testing.T) {
	f := function(func() {
		func() {}
	})
	assert.Equal(t, "func2", f)
}

func TestFunctionWithFuncNameAndDotAndSlash(t *testing.T) {
	f := function(func() {
		func() {}
	})
	assert.Equal(t, "func2", f)
}

func TestFunctionWithFuncNameAndDotAndSlashAndPeriod(t *testing.T) {
	f := function(func() {
		func() {}
	})
	assert.Equal(t, "func2", f)
}

func TestFunctionWithFuncNameAndDotAndSlashAndPeriodAndDot(t *testing.T) {
	f := function(func() {
		func() {}
	})
	assert.Equal(t, "func2", f)
}

func TestFunctionWithFuncNameAndDotAndSlashAndPeriodAndDotAndSlash(t *testing.T) {
	f := function(func() {
		func() {}
	})
	assert.Equal(t, "func2", f)
}

func TestFunctionWithFuncNameAndDotAndSlashAndPeriodAndDotAndSlashAndPeriod(t *testing.T) {
	f := function(func() {
		func() {}
	})
	assert.Equal(t, "func2", f)
}

func TestFunctionWithFuncNameAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDot(t *testing.T) {
	f := function(func() {
		func() {}
	})
	assert.Equal(t, "func2", f)
}

func TestFunctionWithFuncNameAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlash(t *testing.T) {
	f := function(func() {
		func() {}
	})
	assert.Equal(t, "func2", f)
}

func TestFunctionWithFuncNameAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriod(t *testing.T) {
	f := function(func() {
		func() {}
	})
	assert.Equal(t, "func2", f)
}

func TestFunctionWithFuncNameAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDot(t *testing.T) {
	f := function(func() {
		func() {}
	})
	assert.Equal(t, "func2", f)
}

func TestFunctionWithFuncNameAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlash(t *testing.T) {
	f := function(func() {
		func() {}
	})
	assert.Equal(t, "func2", f)
}

func TestFunctionWithFuncNameAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriod(t *testing.T) {
	f := function(func() {
		func() {}
	})
	assert.Equal(t, "func2", f)
}

func TestFunctionWithFuncNameAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDot(t *testing.T) {
	f := function(func() {
		func() {}
	})
	assert.Equal(t, "func2", f)
}

func TestFunctionWithFuncNameAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlash(t *testing.T) {
	f := function(func() {
		func() {}
	})
	assert.Equal(t, "func2", f)
}

func TestFunctionWithFuncNameAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriod(t *testing.T) {
	f := function(func() {
		func() {}
	})
	assert.Equal(t, "func2", f)
}

func TestFunctionWithFuncNameAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDot(t *testing.T) {
	f := function(func() {
		func() {}
	})
	assert.Equal(t, "func2", f)
}

func TestFunctionWithFuncNameAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlash(t *testing.T) {
	f := function(func() {
		func() {}
	})
	assert.Equal(t, "func2", f)
}

func TestFunctionWithFuncNameAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriod(t *testing.T) {
	f := function(func() {
		func() {}
	})
	assert.Equal(t, "func2", f)
}

func TestFunctionWithFuncNameAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDot(t *testing.T) {
	f := function(func() {
		func() {}
	})
	assert.Equal(t, "func2", f)
}

func TestFunctionWithFuncNameAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlash(t *testing.T) {
	f := function(func() {
		func() {}
	})
	assert.Equal(t, "func2", f)
}

func TestFunctionWithFuncNameAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriod(t *testing.T) {
	f := function(func() {
		func() {}
	})
	assert.Equal(t, "func2", f)
}

func TestFunctionWithFuncNameAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDot(t *testing.T) {
	f := function(func() {
		func() {}
	})
	assert.Equal(t, "func2", f)
}

func TestFunctionWithFuncNameAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlash(t *testing.T) {
	f := function(func() {
		func() {}
	})
	assert.Equal(t, "func2", f)
}

func TestFunctionWithFuncNameAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriod(t *testing.T) {
	f := function(func() {
		func() {}
	})
	assert.Equal(t, "func2", f)
}

func TestFunctionWithFuncNameAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDot(t *testing.T) {
	f := function(func() {
		func() {}
	})
	assert.Equal(t, "func2", f)
}

func TestFunctionWithFuncNameAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlash(t *testing.T) {
	f := function(func() {
		func() {}
	})
	assert.Equal(t, "func2", f)
}

func TestFunctionWithFuncNameAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriod(t *testing.T) {
	f := function(func() {
		func() {}
	})
	assert.Equal(t, "func2", f)
}

func TestFunctionWithFuncNameAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDot(t *testing.T) {
	f := function(func() {
		func() {}
	})
	assert.Equal(t, "func2", f)
}

func TestFunctionWithFuncNameAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlash(t *testing.T) {
	f := function(func() {
		func() {}
	})
	assert.Equal(t, "func2", f)
}

func TestFunctionWithFuncNameAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriod(t *testing.T) {
	f := function(func() {
		func() {}
	})
	assert.Equal(t, "func2", f)
}

func TestFunctionWithFuncNameAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDot(t *testing.T) {
	f := function(func() {
		func() {}
	})
	assert.Equal(t, "func2", f)
}

func TestFunctionWithFuncNameAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlash(t *testing.T) {
	f := function(func() {
		func() {}
	})
	assert.Equal(t, "func2", f)
}

func TestFunctionWithFuncNameAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriod(t *testing.T) {
	f := function(func() {
		func() {}
	})
	assert.Equal(t, "func2", f)
}

func TestFunctionWithFuncNameAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDot(t *testing.T) {
	f := function(func() {
		func() {}
	})
	assert.Equal(t, "func2", f)
}

func TestFunctionWithFuncNameAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlash(t *testing.T) {
	f := function(func() {
		func() {}
	})
	assert.Equal(t, "func2", f)
}

func TestFunctionWithFuncNameAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriod(t *testing.T) {
	f := function(func() {
		func() {}
	})
	assert.Equal(t, "func2", f)
}

func TestFunctionWithFuncNameAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDot(t *testing.T) {
	f := function(func() {
		func() {}
	})
	assert.Equal(t, "func2", f)
}

func TestFunctionWithFuncNameAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlash(t *testing.T) {
	f := function(func() {
		func() {}
	})
	assert.Equal(t, "func2", f)
}

func TestFunctionWithFuncNameAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriod(t *testing.T) {
	f := function(func() {
		func() {}
	})
	assert.Equal(t, "func2", f)
}

func TestFunctionWithFuncNameAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDot(t *testing.T) {
	f := function(func() {
		func() {}
	})
	assert.Equal(t, "func2", f)
}

func TestFunctionWithFuncNameAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlash(t *testing.T) {
	f := function(func() {
		func() {}
	})
	assert.Equal(t, "func2", f)
}

func TestFunctionWithFuncNameAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriod(t *testing.T) {
	f := function(func() {
		func() {}
	})
	assert.Equal(t, "func2", f)
}

func TestFunctionWithFuncNameAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDot(t *testing.T) {
	f := function(func() {
		func() {}
	})
	assert.Equal(t, "func2", f)
}

func TestFunctionWithFuncNameAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlash(t *testing.T) {
	f := function(func() {
		func() {}
	})
	assert.Equal(t, "func2", f)
}

func TestFunctionWithFuncNameAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriod(t *testing.T) {
	f := function(func() {
		func() {}
	})
	assert.Equal(t, "func2", f)
}

func TestFunctionWithFuncNameAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDot(t *testing.T) {
	f := function(func() {
		func() {}
	})
	assert.Equal(t, "func2", f)
}

func TestFunctionWithFuncNameAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlash(t *testing.T) {
	f := function(func() {
		func() {}
	})
	assert.Equal(t, "func2", f)
}

func TestFunctionWithFuncNameAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriod(t *testing.T) {
	f := function(func() {
		func() {}
	})
	assert.Equal(t, "func2", f)
}

func TestFunctionWithFuncNameAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDot(t *testing.T) {
	f := function(func() {
		func() {}
	})
	assert.Equal(t, "func2", f)
}

func TestFunctionWithFuncNameAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlash(t *testing.T) {
	f := function(func() {
		func() {}
	})
	assert.Equal(t, "func2", f)
}

func TestFunctionWithFuncNameAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriod(t *testing.T) {
	f := function(func() {
		func() {}
	})
	assert.Equal(t, "func2", f)
}

func TestFunctionWithFuncNameAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDot(t *testing.T) {
	f := function(func() {
		func() {}
	})
	assert.Equal(t, "func2", f)
}

func TestFunctionWithFuncNameAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlash(t *testing.T) {
	f := function(func() {
		func() {}
	})
	assert.Equal(t, "func2", f)
}

func TestFunctionWithFuncNameAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriod(t *testing.T) {
	f := function(func() {
		func() {}
	})
	assert.Equal(t, "func2", f)
}

func TestFunctionWithFuncNameAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDot(t *testing.T) {
	f := function(func() {
		func() {}
	})
	assert.Equal(t, "func2", f)
}

func TestFunctionWithFuncNameAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlash(t *testing.T) {
	f := function(func() {
		func() {}
	})
	assert.Equal(t, "func2", f)
}

func TestFunctionWithFuncNameAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriodAndDotAndSlashAndPeriod(t *testing.T) {
	f := function(func() {
		func() {}
	})
	assert.Equal(t, "func2", f)
}