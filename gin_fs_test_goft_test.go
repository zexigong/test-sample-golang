// Copyright 2017 Manu Martinez-Almeida. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package gin

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOnlyFilesFS(t *testing.T) {
	assert := assert.New(t)

	root := filepath.Join(os.TempDir(), "testOnlyFilesFs")
	filePath := filepath.Join(root, "file")
	dirPath := filepath.Join(root, "dir")

	assert.NoError(os.MkdirAll(dirPath, 0o755))
	assert.NoError(os.WriteFile(filePath, []byte("Gin Web Framework"), 0o666))
	defer os.RemoveAll(root)

	fs := OnlyFilesFS{FileSystem: http.Dir(root)}

	file, err := fs.Open("/file")
	assert.NoError(err)
	defer func() { assert.NoError(file.Close()) }()

	content, err := io.ReadAll(file)
	assert.NoError(err)
	assert.Equal([]byte("Gin Web Framework"), content)

	_, err = fs.Open("/dir")
	assert.Equal(err.Error(), "open /dir: is a directory")
}

func TestDir(t *testing.T) {
	assert := assert.New(t)

	root := filepath.Join(os.TempDir(), "testDir")
	filePath := filepath.Join(root, "file")
	dirPath := filepath.Join(root, "dir")

	assert.NoError(os.MkdirAll(dirPath, 0o755))
	assert.NoError(os.WriteFile(filePath, []byte("Gin Web Framework"), 0o666))
	defer os.RemoveAll(root)

	// listDirectory = true
	{
		fs := Dir(root, true)

		file, err := fs.Open("/file")
		assert.NoError(err)
		defer func() { assert.NoError(file.Close()) }()

		content, err := io.ReadAll(file)
		assert.NoError(err)
		assert.Equal([]byte("Gin Web Framework"), content)

		file, err = fs.Open("/dir")
		assert.NoError(err)
		defer func() { assert.NoError(file.Close()) }()

		_, err = file.Readdir(0)
		assert.NoError(err)
	}

	// listDirectory = false
	{
		fs := Dir(root, false)

		file, err := fs.Open("/file")
		assert.NoError(err)
		defer func() { assert.NoError(file.Close()) }()

		content, err := io.ReadAll(file)
		assert.NoError(err)
		assert.Equal([]byte("Gin Web Framework"), content)

		file, err = fs.Open("/dir")
		assert.NoError(err)
		defer func() { assert.NoError(file.Close()) }()

		files, err := file.Readdir(0)
		assert.NoError(err)
		assert.Nil(files)
	}
}

func TestOnlyFilesFS_Open(t *testing.T) {
	assert := assert.New(t)

	// create a directory with a file in it
	root := filepath.Join(os.TempDir(), "testOnlyFilesFS_Open")
	filePath := filepath.Join(root, "file")
	dirPath := filepath.Join(root, "dir")

	assert.NoError(os.MkdirAll(dirPath, 0o755))
	assert.NoError(os.WriteFile(filePath, []byte("Gin Web Framework"), 0o666))
	defer os.RemoveAll(root)

	fs := OnlyFilesFS{FileSystem: http.Dir(root)}

	// open a file
	file, err := fs.Open("/file")
	assert.NoError(err)
	defer func() { assert.NoError(file.Close()) }()

	content, err := io.ReadAll(file)
	assert.NoError(err)
	assert.Equal([]byte("Gin Web Framework"), content)

	// open a directory
	file, err = fs.Open("/dir")
	assert.Error(err)
	assert.Nil(file)
	assert.True(strings.HasSuffix(err.Error(), "/dir: is a directory"))
}

func TestNeutralizedReaddirFile_Readdir(t *testing.T) {
	assert := assert.New(t)

	// create a directory with a file in it
	root := filepath.Join(os.TempDir(), "testNeutralizedReaddirFile_Readdir")
	filePath := filepath.Join(root, "file")
	dirPath := filepath.Join(root, "dir")

	assert.NoError(os.MkdirAll(dirPath, 0o755))
	assert.NoError(os.WriteFile(filePath, []byte("Gin Web Framework"), 0o666))
	defer os.RemoveAll(root)

	fs := OnlyFilesFS{FileSystem: http.Dir(root)}

	// open a file
	file, err := fs.Open("/file")
	assert.NoError(err)
	defer func() { assert.NoError(file.Close()) }()

	// try to read the directory
	files, err := file.Readdir(0)
	assert.NoError(err)
	assert.Nil(files)
}