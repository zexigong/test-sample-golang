// Copyright 2017 Manu Martinez-Almeida. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package gin

import (
	"errors"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOnlyFilesFS(t *testing.T) {
	fs := OnlyFilesFS{FileSystem: http.Dir(".")}
	file, err := fs.Open("fs_test.go")
	assert.NoError(t, err)
	assert.NotNil(t, file)

	_, err = file.Readdir(0)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, os.ErrPermission))
}

func TestDir(t *testing.T) {
	fs := Dir(".", false)
	file, err := fs.Open("fs_test.go")
	assert.NoError(t, err)
	assert.NotNil(t, file)

	_, err = file.Readdir(0)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, os.ErrPermission))
}

func TestDir_ListDirectory(t *testing.T) {
	fs := Dir(".", true)
	file, err := fs.Open("fs_test.go")
	assert.NoError(t, err)
	assert.NotNil(t, file)

	files, err := file.Readdir(0)
	assert.NoError(t, err)
	assert.NotNil(t, files)
}

type dummyFile struct {
	http.File
}

func (d dummyFile) Readdir(count int) ([]os.FileInfo, error) {
	return nil, nil
}

type dummyFileSystem struct {
}

func (d dummyFileSystem) Open(name string) (http.File, error) {
	return dummyFile{}, nil
}

func TestOnlyFilesFS_Readdir(t *testing.T) {
	fs := OnlyFilesFS{FileSystem: dummyFileSystem{}}
	file, err := fs.Open("fs_test.go")
	assert.NoError(t, err)
	assert.NotNil(t, file)

	_, err = file.Readdir(0)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, os.ErrPermission))
}