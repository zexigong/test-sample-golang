// Copyright 2018 The Hugo Authors. All rights reserved.
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

package collections

import (
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestSlice(t *testing.T) {
	c := qt.New(t)

	c.Assert(Slice(), qt.DeepEquals, []any{})
	c.Assert(Slice(nil, nil), qt.DeepEquals, []any{nil, nil})
	c.Assert(Slice(32, 34, 45), qt.DeepEquals, []int{32, 34, 45})
	c.Assert(Slice("32", "34", "45"), qt.DeepEquals, []string{"32", "34", "45"})
	c.Assert(Slice("32", 45), qt.DeepEquals, []any{"32", 45})

	var s Slicer
	c.Assert(Slice(s), qt.DeepEquals, []any{nil})
	c.Assert(Slice(SortedStringSlice{"b", "a"}), qt.DeepEquals, SortedStringSlice{"a", "b"})

}

func TestSortedStringSliceContains(t *testing.T) {
	c := qt.New(t)

	ss := SortedStringSlice{"a", "b", "b", "c"}

	c.Assert(ss.Contains("a"), qt.Equals, true)
	c.Assert(ss.Contains("b"), qt.Equals, true)
	c.Assert(ss.Contains("c"), qt.Equals, true)
	c.Assert(ss.Contains("d"), qt.Equals, false)
}

func TestSortedStringSliceCount(t *testing.T) {
	c := qt.New(t)

	ss := SortedStringSlice{"a", "b", "b", "b", "c"}

	c.Assert(ss.Count("a"), qt.Equals, 1)
	c.Assert(ss.Count("b"), qt.Equals, 3)
	c.Assert(ss.Count("c"), qt.Equals, 1)
	c.Assert(ss.Count("d"), qt.Equals, 0)
}