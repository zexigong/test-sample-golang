// Copyright 2021 The Hugo Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package hugofs

import (
	"path/filepath"
	"testing"

	"github.com/gohugoio/hugo/hugofs/glob"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

func TestFilenameFilterFs(t *testing.T) {
	fs := afero.NewMemMapFs()
	err := fs.MkdirAll("a", 0777)
	require.NoError(t, err)
	f, err := fs.Create("a/file1.txt")
	require.NoError(t, err)
	f.Close()
	f, err = fs.Create("a/file2.txt")
	require.NoError(t, err)
	f.Close()
	f, err = fs.Create("a/file3.txt")
	require.NoError(t, err)
	f.Close()
	f, err = fs.Create("a/file4.txt")
	require.NoError(t, err)
	f.Close()
	err = fs.MkdirAll("a/b", 0777)
	require.NoError(t, err)
	f, err = fs.Create("a/b/file5.txt")
	require.NoError(t, err)
	f.Close()

	filter := glob.MustNewFilenameFilter(
		[]string{
			filepath.FromSlash("/a/file1.txt"),
			filepath.FromSlash("/a/file2.txt"),
		}, nil)

	filterFs := newFilenameFilterFs(fs, "", filter)

	{
		fis, err := afero.ReadDir(filterFs, "a")
		require.NoError(t, err)
		require.Equal(t, 2, len(fis))
		require.Equal(t, "file1.txt", fis[0].Name())
		require.Equal(t, "file2.txt", fis[1].Name())
	}

	{
		fis, err := afero.ReadDir(filterFs, "a/b")
		require.NoError(t, err)
		require.Equal(t, 0, len(fis))
	}

	{
		fis, err := afero.ReadDir(filterFs, "c")
		require.Error(t, err)
		require.Nil(t, fis)
	}
}