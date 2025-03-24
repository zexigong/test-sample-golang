// Copyright 2024 The Hugo Authors. All rights reserved.
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

package hugio

import (
	"strings"
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestHasBytesWriter(t *testing.T) {
	c := qt.New(t)

	w := &HasBytesWriter{
		Patterns: []*HasBytesPattern{
			{Pattern: []byte("abc")},
			{Pattern: []byte("def")},
		},
	}

	_, err := w.Write([]byte("1abc2def3"))
	c.Assert(err, qt.IsNil)
	for _, pp := range w.Patterns {
		c.Assert(pp.Match, qt.IsTrue)
	}

	w = &HasBytesWriter{
		Patterns: []*HasBytesPattern{
			{Pattern: []byte("abc")},
			{Pattern: []byte("def")},
		},
	}

	// Test with some input that needs to be buffered internally.
	_, err = w.Write([]byte("1ab"))
	c.Assert(err, qt.IsNil)
	_, err = w.Write([]byte("c2d"))
	c.Assert(err, qt.IsNil)
	_, err = w.Write([]byte("ef3"))
	c.Assert(err, qt.IsNil)
	for _, pp := range w.Patterns {
		c.Assert(pp.Match, qt.IsTrue)
	}

	w = &HasBytesWriter{
		Patterns: []*HasBytesPattern{
			{Pattern: []byte("abc")},
			{Pattern: []byte("def")},
		},
	}

	for i := 0; i < 100; i++ {
		_, err = w.Write([]byte(strings.Repeat("1abc", 100)))
		c.Assert(err, qt.IsNil)
		c.Assert(w.Patterns[0].Match, qt.IsTrue)
	}

	_, err = w.Write([]byte("def"))
	c.Assert(err, qt.IsNil)
	for _, pp := range w.Patterns {
		c.Assert(pp.Match, qt.IsTrue)
	}

}