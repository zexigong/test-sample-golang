// Copyright 2019 The Hugo Authors. All rights reserved.
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

func TestAppend(t *testing.T) {
	c := qt.New(t)

	type tp struct {
		in1  any
		in2  []any
		want any
	}

	must := func(r any, err error) any {
		c.Assert(err, qt.IsNil)
		return r
	}

	for _, test := range []tp{
		{nil, nil, nil},
		{nil, []any{"a"}, []any{"a"}},
		{[]string{}, nil, []string{}},
		{[]string{}, []any{"a"}, []string{"a"}},
		{[]string{"a"}, []any{"b"}, []string{"a", "b"}},
		{[]string{"a"}, []any{"b", "c"}, []string{"a", "b", "c"}},
		{[]int{1}, []any{2, 3}, []int{1, 2, 3}},
		{[]int{1}, []any{2}, []int{1, 2}},
		{[]int{1}, []any{nil}, []any{1, nil}},

		// This is a special case, if we have a single argument that is a slice
		// of the same type, we append that slice.
		{[]string{"a"}, []any{[]string{"b", "c"}}, []string{"a", "b", "c"}},

		// Fall back to []interface{} if types doesn't match.
		{[]string{"a"}, []any{[]int{1, 2}}, []any{"a", []int{1, 2}}},
		{[]string{"a"}, []any{[]int{}}, []any{"a", []int{}}},
		{[]string{"a"}, []any{[]string{}}, []string{"a"}},
		{[]string{}, []any{[]string{}}, []string{}},

		// nil
		{[]string{"a"}, []any{nil}, []any{"a", nil}},
		{[]string{"a"}, []any{[]string(nil)}, []string{"a"}},

		// nil slices
		{[]string(nil), nil, []string(nil)},
		{[]string(nil), []any{"a"}, []string{"a"}},
		{[]string(nil), []any{"a", "b"}, []string{"a", "b"}},
		{[]int(nil), []any{1, 2}, []int{1, 2}},
		{[]int(nil), []any{1}, []int{1}},

		// Slices of slices
		{[][]string{}, []any{[]string{"a"}}, [][]string{{"a"}}},
		{[][]string{{"a"}}, []any{[]string{"b"}}, [][]string{{"a"}, {"b"}}},
		{[][]string(nil), []any{[]string{"a"}, []string{"b"}}, [][]string{{"a"}, {"b"}}},

		// Arrays
		{[3]string{}, []any{"a"}, nil},
		{[...]string{}, []any{"a"}, nil},
		{[...]string{"a", "b"}, []any{"c"}, nil},
	} {

		c.Assert(must(Append(test.in1, test.in2...)), qt.DeepEquals, test.want)
	}

}