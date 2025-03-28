package cli

import (
	"bytes"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandleExitCoder_nil(t *testing.T) {
	exitCode := 0
	called := false

	OsExiter = func(rc int) {
		if !called {
			exitCode = rc
			called = true
		}
	}

	defer func() { OsExiter = fakeOsExiter }()

	HandleExitCoder(nil)

	assert.Equal(t, 0, exitCode)
	assert.False(t, called)
}

func TestHandleExitCoder_ExitCoder(t *testing.T) {
	exitCode := 0
	called := false

	OsExiter = func(rc int) {
		if !called {
			exitCode = rc
			called = true
		}
	}

	defer func() { OsExiter = fakeOsExiter }()

	HandleExitCoder(Exit("galactic perimeter breach", 9))

	assert.Equal(t, 9, exitCode)
	assert.True(t, called)
}

func TestHandleExitCoder_ErrorExitCoder(t *testing.T) {
	exitCode := 0
	called := false

	OsExiter = func(rc int) {
		if !called {
			exitCode = rc
			called = true
		}
	}

	defer func() { OsExiter = fakeOsExiter }()

	HandleExitCoder(Exit(errors.New("galactic perimeter breach"), 9))

	assert.Equal(t, 9, exitCode)
	assert.True(t, called)
}

func TestHandleExitCoder_MultiErrorWithExitCoder(t *testing.T) {
	exitCode := 0
	called := false

	OsExiter = func(rc int) {
		if !called {
			exitCode = rc
			called = true
		}
	}

	defer func() { OsExiter = fakeOsExiter }()

	exitErr := Exit("galactic perimeter breach", 9)
	exitErr2 := Exit("last ExitCoder", 11)

	err := newMultiError(errors.New("wowsa"), errors.New("egad"), exitErr, exitErr2)
	HandleExitCoder(err)

	assert.Equal(t, 11, exitCode)
	assert.True(t, called)
}

type exitFormatter struct {
	code int
}

func (f *exitFormatter) Format(s fmt.State, verb rune) {
	_, _ = s.Write([]byte("some other special"))
}

func (f *exitFormatter) ExitCode() int {
	return f.code
}

func (f *exitFormatter) Error() string {
	return fmt.Sprintf("my special error code %d", f.code)
}

func TestHandleExitCoder_ErrorFormatter(t *testing.T) {
	exitCode := 0
	called := false

	OsExiter = func(rc int) {
		if !called {
			exitCode = rc
			called = true
		}
	}

	oldWriter := ErrWriter
	var buf bytes.Buffer
	ErrWriter = &buf
	defer func() {
		OsExiter = fakeOsExiter
		ErrWriter = oldWriter
	}()

	exitErr := Exit("galactic perimeter breach", 9)
	exitErr2 := Exit("last ExitCoder", 11)
	exitErr3 := &exitFormatter{code: 12}

	// add some recursion for multi error to fix test coverage
	err := newMultiError(errors.New("wowsa"), errors.New("egad"), exitErr3, newMultiError(exitErr, exitErr2))
	HandleExitCoder(err)

	assert.Equal(t, 11, exitCode)
	assert.True(t, called)
	assert.Contains(t, buf.String(), "some other special")
}

func TestHandleExitCoder_MultiErrorWithoutExitCoder(t *testing.T) {
	exitCode := 0
	called := false

	OsExiter = func(rc int) {
		if !called {
			exitCode = rc
			called = true
		}
	}

	defer func() { OsExiter = fakeOsExiter }()

	err := newMultiError(errors.New("wowsa"), errors.New("egad"))
	HandleExitCoder(err)

	assert.Equal(t, 1, exitCode)
	assert.True(t, called)
}

// make a stub to not import pkg/errors
type ErrorWithFormat struct {
	error
}

func NewErrorWithFormat(m string) *ErrorWithFormat {
	return &ErrorWithFormat{error: errors.New(m)}
}

func (f *ErrorWithFormat) Format(s fmt.State, verb rune) {
	fmt.Fprintf(s, "This the format: %v", f.error)
}

func TestHandleExitCoder_ErrorWithFormat(t *testing.T) {
	called := false

	OsExiter = func(int) {
		if !called {
			called = true
		}
	}
	ErrWriter = &bytes.Buffer{}

	defer func() {
		OsExiter = fakeOsExiter
		ErrWriter = fakeErrWriter
	}()

	err := Exit(NewErrorWithFormat("I am formatted"), 1)
	HandleExitCoder(err)

	assert.True(t, called)
	assert.Equal(t, ErrWriter.(*bytes.Buffer).String(), "This the format: I am formatted\n")
}

func TestHandleExitCoder_MultiErrorWithFormat(t *testing.T) {
	called := false

	OsExiter = func(int) {
		if !called {
			called = true
		}
	}
	ErrWriter = &bytes.Buffer{}

	defer func() { OsExiter = fakeOsExiter }()

	err := newMultiError(NewErrorWithFormat("err1"), NewErrorWithFormat("err2"))
	HandleExitCoder(err)

	assert.True(t, called)
	assert.Equal(t, ErrWriter.(*bytes.Buffer).String(), "This the format: err1\nThis the format: err2\n")
}

func TestMultiErrorErrorsCopy(t *testing.T) {
	errList := []error{
		errors.New("foo"),
		errors.New("bar"),
		errors.New("baz"),
	}
	me := newMultiError(errList...)
	assert.Equal(t, errList, me.Errors())
}

func TestErrRequiredFlags_Error(t *testing.T) {
	missingFlags := []string{"flag1", "flag2"}
	err := &errRequiredFlags{missingFlags: missingFlags}
	expectedMsg := "Required flags \"flag1, flag2\" not set"
	assert.Equal(t, expectedMsg, err.Error())

	missingFlags = []string{"flag1"}
	err = &errRequiredFlags{missingFlags: missingFlags}
	expectedMsg = "Required flag \"flag1\" not set"
	assert.Equal(t, expectedMsg, err.Error())
}
