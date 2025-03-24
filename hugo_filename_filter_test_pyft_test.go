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

	qt "github.com/frankban/quicktest"
)

func TestFilenameFilter_Match(t *testing.T) {
	c := qt.New(t)
	for _, test := range []struct {
		name       string
		inclusions []string
		exclusions []string
		files      map[string]bool
	}{
		{
			name:       "Inclusions",
			inclusions: []string{"*.md"},
			files: map[string]bool{
				"readme.md":  true,
				"readme.txt": false,
			},
		},
		{
			name:       "Inclusions with dir",
			inclusions: []string{"content/*.md"},
			files: map[string]bool{
				"content/readme.md": true,
				"readme.txt":        false,
			},
		},
		{
			name:       "Inclusions with dir and nested",
			inclusions: []string{"content/**/*.md"},
			files: map[string]bool{
				"content/readme.md":             true,
				"content/sub/readme.md":         true,
				"content/sub/sub2/readme.md":    true,
				"content/sub/sub2/readme2.md":   true,
				"content/sub/sub2/README2.md":   true,
				"content/sub/sub2/README2.txt":  false,
				"content/sub/sub2/README2.json": false,
				"content/sub/sub2/README2.yaml": false,
				"content/sub/readme.txt":        false,
				"readme.txt":                    false,
			},
		},
		{
			name:       "Exclusions",
			exclusions: []string{"*.md"},
			files: map[string]bool{
				"readme.md":  false,
				"readme.txt": true,
			},
		},
		{
			name:       "Inclusions and Exclusions",
			inclusions: []string{"*.*"},
			exclusions: []string{"*.md"},
			files: map[string]bool{
				"readme.md":  false,
				"readme.txt": true,
			},
		},
		{
			name:       "Inclusions and Exclusions no ext",
			inclusions: []string{"*"},
			exclusions: []string{"*.md"},
			files: map[string]bool{
				"readme.md":  false,
				"readme.txt": true,
			},
		},
		{
			name:       "Inclusions and Exclusions dir",
			inclusions: []string{"content/**"},
			exclusions: []string{"content/**/*.md"},
			files: map[string]bool{
				"content/readme.md":            false,
				"content/sub/readme.md":        false,
				"content/sub/sub2/readme.md":   false,
				"content/sub/sub2/readme2.md":  false,
				"content/sub/readme.txt":       true,
				"content/sub/sub2/README2.txt": true,
				"content/sub/sub2/README2.md":  false,
				"readme.txt":                   false,
			},
		},
	} {
		c.Run(test.name, func(c *qt.C) {
			filter := MustNewFilenameFilter(test.inclusions, test.exclusions)

			for filename, expect := range test.files {
				c.Assert(filter.Match(filename, false), qt.Equals, expect, qt.Commentf("file: %s", filename))
			}

			c.Assert(func() {
				MustNewFilenameFilter([]string{"{"}, nil)
			}, qt.PanicMatches, "syntax error in pattern")
		})
	}

}

func TestFilenameFilter_MatchWithInclusionFunc(t *testing.T) {
	c := qt.New(t)

	filter := NewFilenameFilterForInclusionFunc(func(filename string) bool {
		return strings.HasSuffix(filename, "txt")
	})

	for _, file := range []string{
		"content/file.txt",
		"content/file.md",
		"file.txt",
		"/file.txt",
		"file.md",
	} {
		c.Assert(filter.Match(file, false), qt.Equals, strings.HasSuffix(file, "txt"))
	}
}