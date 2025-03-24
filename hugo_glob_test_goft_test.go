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

package glob

import (
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestGlobCache(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	gc := &globCache{
		isWindows: false,
		cache:     make(map[string]globErr),
	}

	g, err := gc.GetGlob("ABC")
	c.Assert(err, qt.IsNil)
	c.Assert(g.Match("abC"), qt.Equals, true)

	g, err = gc.GetGlob("DEF")
	c.Assert(err, qt.IsNil)
	c.Assert(g.Match("deF"), qt.Equals, true)

	g, err = gc.GetGlob("ABC")
	c.Assert(err, qt.IsNil)
	c.Assert(g.Match("abC"), qt.Equals, true)

	c.Assert(len(gc.cache), qt.Equals, 2)

	gc = &globCache{
		isWindows: true,
		cache:     make(map[string]globErr),
	}

	g, err = gc.GetGlob("ABC/def")
	c.Assert(err, qt.IsNil)
	c.Assert(g.Match("AbC\\deF"), qt.Equals, true)
}

func TestGlobOr(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	g1, _ := GetGlob("ABC")
	g2, _ := GetGlob("DEF")
	g3, _ := GetGlob("GHI")

	g := Or(g1, g2, g3)

	c.Assert(g.Match("abc"), qt.Equals, true)
	c.Assert(g.Match("def"), qt.Equals, true)
	c.Assert(g.Match("ghi"), qt.Equals, true)
	c.Assert(g.Match("jkl"), qt.Equals, false)
}

func TestNormalizePath(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	c.Assert(NormalizePath(`a\b\c`), qt.Equals, "a/b/c")
	c.Assert(NormalizePath(`a\b\c\`), qt.Equals, "a/b/c")
	c.Assert(NormalizePath(`.\a\b\c\.`), qt.Equals, "a/b/c")

	c.Assert(NormalizePathNoLower(`a\b\c`), qt.Equals, "a/b/c")
	c.Assert(NormalizePathNoLower(`a\b\c\`), qt.Equals, "a/b/c")
	c.Assert(NormalizePathNoLower(`.\a\b\c\.`), qt.Equals, "a/b/c")
}

func TestResolveRootDir(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	c.Assert(ResolveRootDir("a/b/c"), qt.Equals, "a/b/c")
	c.Assert(ResolveRootDir("a/b/c/*.json"), qt.Equals, "a/b/c")
	c.Assert(ResolveRootDir("a/b/*/c/*.json"), qt.Equals, "a/b")
	c.Assert(ResolveRootDir("assets/**.json"), qt.Equals, "assets")
	c.Assert(ResolveRootDir("**.json"), qt.Equals, "")
}

func TestFilterGlobParts(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	c.Assert(FilterGlobParts([]string{"a", "b", "c"}), qt.DeepEquals, []string{"a", "b", "c"})
	c.Assert(FilterGlobParts([]string{"a", "*", "c"}), qt.DeepEquals, []string{"a", "c"})
	c.Assert(FilterGlobParts([]string{"a", "*", "c", "**"}), qt.DeepEquals, []string{"a", "c"})
	c.Assert(FilterGlobParts([]string{"*", "*", "*"}), qt.DeepEquals, []string{})
}

func TestHasGlobChar(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	c.Assert(HasGlobChar("a"), qt.Equals, false)
	c.Assert(HasGlobChar("*"), qt.Equals, true)
	c.Assert(HasGlobChar("a*b"), qt.Equals, true)
	c.Assert(HasGlobChar("a*b*c"), qt.Equals, true)
	c.Assert(HasGlobChar("a*b*c*"), qt.Equals, true)
}