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
	"os"
	"path/filepath"
	"testing"

	"github.com/gohugoio/hugo/hugofs/glob"

	"github.com/spf13/afero"

	qt "github.com/frankban/quicktest"
)

func TestFilenameFilterFs(t *testing.T) {
	c := qt.New(t)

	c.Run("Basic", func(c *qt.C) {
		fs := NewMem()
		afero.WriteFile(fs, "a.txt", []byte("a"), 0777)
		afero.WriteFile(fs, "b.txt", []byte("b"), 0777)
		afero.WriteFile(fs, "c.txt", []byte("c"), 0777)
		afero.WriteFile(fs, "d.txt", []byte("d"), 0777)

		filter, _ := glob.NewFilenameFilter([]string{"a.txt", "b.txt"}, nil)

		fffs := newFilenameFilterFs(fs, "", filter)

		_, err := fffs.Stat("a.txt")
		c.Assert(err, qt.IsNil)

		_, err = fffs.Stat("c.txt")
		c.Assert(err, qt.Not(qt.IsNil))

		f, err := fffs.Open("b.txt")
		c.Assert(err, qt.IsNil)
		f.Close()

		f, err = fffs.Open("c.txt")
		c.Assert(err, qt.Not(qt.IsNil))
		if f != nil {
			f.Close()
		}

		// This is not supported
		_, err = fffs.OpenFile("a.txt", os.O_APPEND, 0777)
		c.Assert(err, qt.IsNil)

		_, err = fffs.OpenFile("c.txt", os.O_APPEND, 0777)
		c.Assert(err, qt.Not(qt.IsNil))
	})

	c.Run("Dir", func(c *qt.C) {
		fs := NewMem()
		afero.WriteFile(fs, "a.txt", []byte("a"), 0777)
		afero.WriteFile(fs, "b.txt", []byte("b"), 0777)
		afero.WriteFile(fs, "c.txt", []byte("c"), 0777)
		afero.WriteFile(fs, "d.txt", []byte("d"), 0777)
		afero.WriteFile(fs, "e/f.txt", []byte("f"), 0777)

		filter, _ := glob.NewFilenameFilter([]string{"*.txt", "e"}, []string{"c.txt"})

		fffs := newFilenameFilterFs(fs, "", filter)

		f, err := fffs.Open("")
		c.Assert(err, qt.IsNil)

		dir, err := f.Readdir(0)
		c.Assert(err, qt.IsNil)
		c.Assert(dir, qt.HasLen, 3)

		dirnames, err := f.Readdirnames(0)
		c.Assert(err, qt.IsNil)
		c.Assert(dirnames, qt.HasLen, 3)

		f.Close()

		_, err = fffs.Stat("e")
		c.Assert(err, qt.IsNil)

		_, err = fffs.Stat(filepath.FromSlash("e/f.txt"))
		c.Assert(err, qt.IsNil)
	})

}