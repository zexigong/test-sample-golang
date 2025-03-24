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
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGlobCache(t *testing.T) {
	c := &globCache{
		cache: make(map[string]globErr),
	}

	g, err := c.GetGlob("a*b")
	assert.NoError(t, err)
	assert.True(t, g.Match("AB"))
	assert.True(t, g.Match("ab"))

	// From cache
	g, err = c.GetGlob("a*b")
	assert.NoError(t, err)
	assert.True(t, g.Match("AB"))
	assert.True(t, g.Match("ab"))

}

func TestGlob(t *testing.T) {
	g, err := GetGlob("a*b")
	assert.NoError(t, err)
	assert.True(t, g.Match("AB"))
	assert.True(t, g.Match("ab"))

	// Or
	g, err = GetGlob("a*b")
	g = Or(g, MatchesFunc(func(s string) bool {
		return strings.HasPrefix(s, "123")
	}))
	assert.NoError(t, err)
	assert.True(t, g.Match("ab"))
	assert.True(t, g.Match("123"))
	assert.True(t, g.Match("123something"))
	assert.False(t, g.Match("no"))

	// Invalid pattern
	g, err = GetGlob("a[")
	assert.Error(t, err)
	assert.Nil(t, g)
}

func TestNormalizePath(t *testing.T) {
	assert.Equal(t, "abc/cde", NormalizePath("Abc/Cde"))
	assert.Equal(t, "abc/cde", NormalizePath("abc/cde"))
	assert.Equal(t, "abc/cde", NormalizePath("abc/cde/"))
	assert.Equal(t, "abc/cde", NormalizePath("/abc/cde/"))
	assert.Equal(t, "abc/cde", NormalizePath("abc/cde/./"))
	assert.Equal(t, "abc/cde", NormalizePath("abc/cde/./."))
	assert.Equal(t, "", NormalizePath("abc/cde/../"))
	assert.Equal(t, "cde", NormalizePath("abc/cde/../cde"))
}

func TestNormalizePathNoLower(t *testing.T) {
	assert.Equal(t, "Abc/Cde", NormalizePathNoLower("Abc/Cde"))
	assert.Equal(t, "abc/cde", NormalizePathNoLower("abc/cde"))
	assert.Equal(t, "abc/cde", NormalizePathNoLower("abc/cde/"))
	assert.Equal(t, "abc/cde", NormalizePathNoLower("/abc/cde/"))
	assert.Equal(t, "abc/cde", NormalizePathNoLower("abc/cde/./"))
	assert.Equal(t, "abc/cde", NormalizePathNoLower("abc/cde/./."))
	assert.Equal(t, "", NormalizePathNoLower("abc/cde/../"))
	assert.Equal(t, "cde", NormalizePathNoLower("abc/cde/../cde"))
}

func TestResolveRootDir(t *testing.T) {
	assert.Equal(t, "assets/sass", ResolveRootDir("assets/sass/**.scss"))
	assert.Equal(t, "assets/sass", ResolveRootDir("assets/sass/**/*.scss"))
	assert.Equal(t, "assets/sass", ResolveRootDir("assets/sass/.scss"))
	assert.Equal(t, "", ResolveRootDir("assets/**.scss"))
	assert.Equal(t, "", ResolveRootDir("assets"))
	assert.Equal(t, "assets", ResolveRootDir("assets/"))
}

func TestFilterGlobParts(t *testing.T) {
	assert.Equal(t, []string{}, FilterGlobParts([]string{}))
	assert.Equal(t, []string{}, FilterGlobParts([]string{"**"}))
	assert.Equal(t, []string{}, FilterGlobParts([]string{"**", "foo*"}))
	assert.Equal(t, []string{"foo"}, FilterGlobParts([]string{"foo", "foo*"}))
	assert.Equal(t, []string{"foo", "bar"}, FilterGlobParts([]string{"foo", "foo*", "bar"}))
}

func TestHasGlobChar(t *testing.T) {
	assert.False(t, HasGlobChar("a"))
	assert.True(t, HasGlobChar("*"))
	assert.True(t, HasGlobChar("a*a"))
	assert.True(t, HasGlobChar("a*a/**"))
	assert.True(t, HasGlobChar("**"))
}