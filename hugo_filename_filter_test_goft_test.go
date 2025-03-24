// Copyright 2021 The Hugo Authors. All rights reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package glob

import (
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestNewFilenameFilter(t *testing.T) {
	t.Parallel()

	c := qt.New(t)

	filter, err := NewFilenameFilter([]string{"a/**", "b/c/**"}, []string{"a/c/exclude", "b/c/exclude"})
	c.Assert(err, qt.IsNil)

	c.Assert(filter.Match("a/b/c/d/file.txt", false), qt.Equals, true)
	c.Assert(filter.Match("a/c/exclude", false), qt.Equals, false)
	c.Assert(filter.Match("a/c/exclude/some/file.txt", false), qt.Equals, false)
	c.Assert(filter.Match("b/c/exclude", false), qt.Equals, false)
	c.Assert(filter.Match("b/c/exclude/file.txt", false), qt.Equals, false)
	c.Assert(filter.Match("b/c/included", false), qt.Equals, true)
	c.Assert(filter.Match("b/c/included/file.txt", false), qt.Equals, true)

	c.Assert(filter.Match("b/c", true), qt.Equals, true)
	c.Assert(filter.Match("b", true), qt.Equals, false)
	c.Assert(filter.Match("b/c/included", true), qt.Equals, true)
	c.Assert(filter.Match("b/c/included/file.txt", false), qt.Equals, true)

	filter = filter.Append(MustNewFilenameFilter(nil, []string{"b/c/included"}))

	c.Assert(filter.Match("b/c/included/file.txt", false), qt.Equals, false)
	c.Assert(filter.Match("b/c/included", true), qt.Equals, false)

	filter = MustNewFilenameFilter([]string{"a/**"}, nil)

	c.Assert(filter.Match("a/b/c/d/file.txt", false), qt.Equals, true)
	c.Assert(filter.Match("a/c/exclude", false), qt.Equals, true)
	c.Assert(filter.Match("b/c/exclude", false), qt.Equals, false)
	c.Assert(filter.Match("b/c/exclude/file.txt", false), qt.Equals, false)

	filter = MustNewFilenameFilter([]string{"a/b/c/d/file.txt"}, nil)

	c.Assert(filter.Match("a/b/c/d/file.txt", false), qt.Equals, true)
	c.Assert(filter.Match("a/b/c/d", true), qt.Equals, true)
	c.Assert(filter.Match("a/b/c", true), qt.Equals, true)
	c.Assert(filter.Match("a/b", true), qt.Equals, true)
	c.Assert(filter.Match("a", true), qt.Equals, true)

	filter = MustNewFilenameFilter(nil, nil)

	c.Assert(filter.Match("a/b/c/d/file.txt", false), qt.Equals, true)

	filter = MustNewFilenameFilter([]string{"a/b/c/d/file.txt"}, []string{"a/b/c/d/file.txt"})

	c.Assert(filter.Match("a/b/c/d/file.txt", false), qt.Equals, false)

	filter = MustNewFilenameFilter(nil, []string{"a/b/c/d/file.txt"})

	c.Assert(filter.Match("a/b/c/d/file.txt", false), qt.Equals, false)

	filter = MustNewFilenameFilter(nil, []string{"a/b/c/d/file.txt"})

	c.Assert(filter.Match("a/b/c/d/file.txt", false), qt.Equals, false)

	filter = NewFilenameFilterForInclusionFunc(func(s string) bool {
		return s == "a/b/c/d/file.txt" || s == "a/b/c/d"
	})

	c.Assert(filter.Match("a/b/c/d/file.txt", false), qt.Equals, true)
	c.Assert(filter.Match("a/b/c/d", true), qt.Equals, true)
	c.Assert(filter.Match("a/b/c", true), qt.Equals, false)

	filter = filter.Append(NewFilenameFilterForInclusionFunc(func(s string) bool {
		return s == "a/b/c"
	}))

	c.Assert(filter.Match("a/b/c/d/file.txt", false), qt.Equals, true)
	c.Assert(filter.Match("a/b/c/d", true), qt.Equals, true)
	c.Assert(filter.Match("a/b/c", true), qt.Equals, true)
}

func TestNewFilenameFilterWindowsPaths(t *testing.T) {
	t.Parallel()

	c := qt.New(t)

	defer func() {
		isWindows = false
	}()
	isWindows = true

	filter, err := NewFilenameFilter([]string{"a/**", "b/c/**"}, []string{"a/c/exclude", "b/c/exclude"})
	c.Assert(err, qt.IsNil)

	c.Assert(filter.Match("a/b/c/d/file.txt", false), qt.Equals, true)
	c.Assert(filter.Match("a\\b\\c\\d\\file.txt", false), qt.Equals, true)
	c.Assert(filter.Match("b\\c\\included\\file.txt", false), qt.Equals, true)
	c.Assert(filter.Match("a/c/exclude", false), qt.Equals, false)
	c.Assert(filter.Match("a\\c\\exclude", false), qt.Equals, false)
	c.Assert(filter.Match("a\\c\\exclude\\some\\file.txt", false), qt.Equals, false)
	c.Assert(filter.Match("b/c/exclude", false), qt.Equals, false)
	c.Assert(filter.Match("b\\c\\exclude", false), qt.Equals, false)
	c.Assert(filter.Match("b\\c\\exclude\\file.txt", false), qt.Equals, false)
	c.Assert(filter.Match("b/c/included", false), qt.Equals, true)
	c.Assert(filter.Match("b/c/included/file.txt", false), qt.Equals, true)

	filter = filter.Append(MustNewFilenameFilter(nil, []string{"b/c/included"}))

	c.Assert(filter.Match("b/c/included/file.txt", false), qt.Equals, false)
	c.Assert(filter.Match("b\\c\\included\\file.txt", false), qt.Equals, false)
	c.Assert(filter.Match("b/c/included", true), qt.Equals, false)

	filter = MustNewFilenameFilter([]string{"a/**"}, nil)

	c.Assert(filter.Match("a/b/c/d/file.txt", false), qt.Equals, true)
	c.Assert(filter.Match("a/c/exclude", false), qt.Equals, true)
	c.Assert(filter.Match("b/c/exclude", false), qt.Equals, false)
	c.Assert(filter.Match("b/c/exclude/file.txt", false), qt.Equals, false)

	filter = MustNewFilenameFilter([]string{"a/b/c/d/file.txt"}, nil)

	c.Assert(filter.Match("a/b/c/d/file.txt", false), qt.Equals, true)
	c.Assert(filter.Match("a/b/c/d", true), qt.Equals, true)
	c.Assert(filter.Match("a/b/c", true), qt.Equals, true)
	c.Assert(filter.Match("a/b", true), qt.Equals, true)
	c.Assert(filter.Match("a", true), qt.Equals, true)

	filter = MustNewFilenameFilter(nil, nil)

	c.Assert(filter.Match("a/b/c/d/file.txt", false), qt.Equals, true)

	filter = MustNewFilenameFilter([]string{"a/b/c/d/file.txt"}, []string{"a/b/c/d/file.txt"})

	c.Assert(filter.Match("a/b/c/d/file.txt", false), qt.Equals, false)

	filter = MustNewFilenameFilter(nil, []string{"a/b/c/d/file.txt"})

	c.Assert(filter.Match("a/b/c/d/file.txt", false), qt.Equals, false)

	filter = MustNewFilenameFilter(nil, []string{"a/b/c/d/file.txt"})

	c.Assert(filter.Match("a/b/c/d/file.txt", false), qt.Equals, false)

	filter = NewFilenameFilterForInclusionFunc(func(s string) bool {
		return s == "a/b/c/d/file.txt" || s == "a/b/c/d"
	})

	c.Assert(filter.Match("a/b/c/d/file.txt", false), qt.Equals, true)
	c.Assert(filter.Match("a/b/c/d", true), qt.Equals, true)
	c.Assert(filter.Match("a/b/c", true), qt.Equals, false)

	filter = filter.Append(NewFilenameFilterForInclusionFunc(func(s string) bool {
		return s == "a/b/c"
	}))

	c.Assert(filter.Match("a/b/c/d/file.txt", false), qt.Equals, true)
	c.Assert(filter.Match("a/b/c/d", true), qt.Equals, true)
	c.Assert(filter.Match("a/b/c", true), qt.Equals, true)
}