package hugofs

import (
	"os"
	"testing"

	"github.com/gohugoio/hugo/hugofs/glob"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestFilenameFilterFs_Open(t *testing.T) {
	memMapFs := afero.NewMemMapFs()
	basePath := "/"
	_ = afero.WriteFile(memMapFs, "/foo.txt", []byte("foo"), 0644)
	_ = afero.WriteFile(memMapFs, "/bar.txt", []byte("bar"), 0644)

	filter := glob.MustNewFilenameFilter([]string{"/foo.txt"}, nil)
	filterFs := newFilenameFilterFs(memMapFs, basePath, filter)

	file, err := filterFs.Open("/foo.txt")
	assert.NoError(t, err)
	assert.NotNil(t, file)
	file.Close()

	file, err = filterFs.Open("/bar.txt")
	assert.Error(t, err)
	assert.Nil(t, file)
}

func TestFilenameFilterFs_Stat(t *testing.T) {
	memMapFs := afero.NewMemMapFs()
	basePath := "/"
	_ = afero.WriteFile(memMapFs, "/foo.txt", []byte("foo"), 0644)

	filter := glob.MustNewFilenameFilter([]string{"/foo.txt"}, nil)
	filterFs := newFilenameFilterFs(memMapFs, basePath, filter)

	fi, err := filterFs.Stat("/foo.txt")
	assert.NoError(t, err)
	assert.NotNil(t, fi)

	fi, err = filterFs.Stat("/bar.txt")
	assert.Error(t, err)
	assert.Nil(t, fi)
}

func TestFilenameFilterFs_Open_Directory(t *testing.T) {
	memMapFs := afero.NewMemMapFs()
	basePath := "/"
	_ = memMapFs.Mkdir("/dir", 0755)
	_ = afero.WriteFile(memMapFs, "/dir/foo.txt", []byte("foo"), 0644)

	filter := glob.MustNewFilenameFilter([]string{"/dir"}, nil)
	filterFs := newFilenameFilterFs(memMapFs, basePath, filter)

	file, err := filterFs.Open("/dir")
	assert.NoError(t, err)
	assert.NotNil(t, file)
	file.Close()
}

func TestFilenameFilterDir_ReadDir(t *testing.T) {
	memMapFs := afero.NewMemMapFs()
	basePath := "/"
	_ = memMapFs.Mkdir("/dir", 0755)
	_ = afero.WriteFile(memMapFs, "/dir/foo.txt", []byte("foo"), 0644)
	_ = afero.WriteFile(memMapFs, "/dir/bar.txt", []byte("bar"), 0644)

	filter := glob.MustNewFilenameFilter([]string{"/dir/foo.txt"}, nil)
	filterFs := newFilenameFilterFs(memMapFs, basePath, filter)

	dirFile, err := filterFs.Open("/dir")
	assert.NoError(t, err)

	dir, ok := dirFile.(fs.ReadDirFile)
	assert.True(t, ok)

	entries, err := dir.ReadDir(-1)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(entries))
	assert.Equal(t, "foo.txt", entries[0].Name())
}

func TestFilenameFilterFs_Permissions(t *testing.T) {
	memMapFs := afero.NewMemMapFs()
	basePath := "/"
	filter := glob.MustNewFilenameFilter(nil, nil)
	filterFs := newFilenameFilterFs(memMapFs, basePath, filter)

	err := filterFs.Chmod("/foo.txt", 0644)
	assert.Error(t, err)

	err = filterFs.Chtimes("/foo.txt", time.Now(), time.Now())
	assert.Error(t, err)

	err = filterFs.Chown("/foo.txt", 1, 1)
	assert.Error(t, err)

	err = filterFs.Remove("/foo.txt")
	assert.Error(t, err)

	err = filterFs.RemoveAll("/foo.txt")
	assert.Error(t, err)

	err = filterFs.Rename("/foo.txt", "/bar.txt")
	assert.Error(t, err)

	_, err = filterFs.Create("/foo.txt")
	assert.Error(t, err)

	err = filterFs.Mkdir("/foo", 0755)
	assert.Error(t, err)

	err = filterFs.MkdirAll("/foo/bar", 0755)
	assert.Error(t, err)
}