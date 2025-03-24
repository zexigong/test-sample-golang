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
	"reflect"
	"testing"

	qt "github.com/frankban/quicktest"
)

type testSlicer string

func (ts testSlicer) Slice(items any) (any, error) {
	return []testSlicer{"a", "b"}, nil
}

type testSlicerFail string

func (ts testSlicerFail) Slice(items any) (any, error) {
	return nil, errors.New("failed")
}

func TestSlice(t *testing.T) {
	t.Parallel()

	c := qt.New(t)

	c.Assert(Slice(), qt.DeepEquals, []any{})

	v := testSlicer("a")
	c.Assert(Slice(v, v), qt.DeepEquals, []testSlicer{"a", "b"})

	vf := testSlicerFail("a")
	c.Assert(Slice(vf, vf), qt.DeepEquals, []any{vf, vf})

	c.Assert(Slice(32, 33), qt.DeepEquals, []int{32, 33})
	c.Assert(Slice(32, "33"), qt.DeepEquals, []any{32, "33"})
	c.Assert(Slice(32), qt.DeepEquals, []int{32})
}

func TestStringSliceToInterfaceSlice(t *testing.T) {
	t.Parallel()

	c := qt.New(t)

	for _, test := range []struct {
		slice []string
	}{
		{nil},
		{[]string{}},
		{[]string{"a", "b", "c"}},
	} {
		c.Run(reflect.ValueOf(test.slice).String(), func(c *qt.C) {
			interfaceSlice := StringSliceToInterfaceSlice(test.slice)
			c.Assert(len(interfaceSlice), qt.Equals, len(test.slice))
			for i := 0; i < len(test.slice); i++ {
				c.Assert(interfaceSlice[i], qt.Equals, test.slice[i])
			}
		})
	}
}

func TestSortedStringSliceContains(t *testing.T) {
	t.Parallel()

	c := qt.New(t)

	for _, test := range []struct {
		slice     SortedStringSlice
		contains  []string
		notIn     []string
	}{
		{nil, []string{}, []string{"a", "b"}},
		{SortedStringSlice{"a", "b", "c"}, []string{"a", "b", "c"}, []string{"aa", "bb"}},
	} {
		c.Run(reflect.ValueOf(test.slice).String(), func(c *qt.C) {
			for _, v := range test.contains {
				c.Assert(test.slice.Contains(v), qt.Equals, true)
			}
			for _, v := range test.notIn {
				c.Assert(test.slice.Contains(v), qt.Equals, false)
			}
		})
	}
}

func TestSortedStringSliceCount(t *testing.T) {
	t.Parallel()

	c := qt.New(t)

	for _, test := range []struct {
		slice     SortedStringSlice
		counts    map[string]int
	}{
		{nil, map[string]int{}},
		{SortedStringSlice{"a", "a", "b", "b", "c"}, map[string]int{"a": 2, "b": 2, "c": 1, "d": 0}},
	} {
		c.Run(reflect.ValueOf(test.slice).String(), func(c *qt.C) {
			for k, v := range test.counts {
				c.Assert(test.slice.Count(k), qt.Equals, v)
			}
		})
	}
}