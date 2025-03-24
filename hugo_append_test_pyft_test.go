// Copyright 2019 The Hugo Authors. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package collections

import (
	"reflect"
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestAppend(t *testing.T) {
	t.Parallel()

	c := qt.New(t)

	for _, test := range []struct {
		to     any
		from   []any
		expect any
	}{
		{[]int{1}, []any{2, 3}, []int{1, 2, 3}},
		{nil, []any{2, 3}, []any{2, 3}},
		{nil, []any{"foo", "bar"}, []any{"foo", "bar"}},
		{nil, []any{[]string{"foo", "bar"}}, []string{"foo", "bar"}},
		{[]string{"a"}, []any{[]string{"foo", "bar"}}, []string{"a", "foo", "bar"}},
		{nil, []any{[]string{"foo"}, []string{"bar"}}, []any{[]string{"foo"}, []string{"bar"}}},
		{[]string{"a"}, []any{[]string{"foo"}, []string{"bar"}}, []any{"a", []string{"foo"}, []string{"bar"}}},
		{[]any{"a"}, []any{[]string{"foo"}, []string{"bar"}}, []any{"a", []string{"foo"}, []string{"bar"}}},
		{nil, nil, nil},
		{nil, []any{nil}, []any{nil}},
		{[]any{"a"}, []any{nil}, []any{"a", nil}},
	} {
		to := test.to
		from := test.from
		expect := test.expect

		result, err := Append(to, from...)
		c.Assert(err, qt.IsNil)
		c.Assert(result, qt.DeepEquals, expect)
	}

	// errors
	for _, test := range []struct {
		to   any
		from []any
	}{
		{nil, []any{[]int{32}, "foo"}},
		{[]string{"a"}, []any{[]string{"foo", "bar"}, 32}},
		{[]int{1}, []any{"foo"}},
		{[]int{1}, []any{[]int{1}, []string{"foo"}}},
	} {
		to := test.to
		from := test.from
		_, err := Append(to, from...)
		c.Assert(err, qt.Not(qt.IsNil))
	}
}

func TestAppendNilSlices(t *testing.T) {
	t.Parallel()

	c := qt.New(t)

	// https://github.com/gohugoio/hugo/issues/6589
	var nilSlice []interface{}
	slice, err := Append(nilSlice, nilSlice)
	c.Assert(err, qt.IsNil)
	c.Assert(slice, qt.DeepEquals, []interface{}{nil})
}

func TestAppendToInterfaceSlice(t *testing.T) {
	t.Parallel()

	c := qt.New(t)

	for _, test := range []struct {
		slice1 any
		slice2 any
		expect any
	}{
		{[]any{1}, []any{2, 3}, []any{1, 2, 3}},
		{[]any{1}, []int{2, 3}, []any{1, 2, 3}},
		{[]int{1}, []int{2, 3}, []any{1, 2, 3}},
		{[]any{1}, []int{2, 3}, []any{1, 2, 3}},
		{[]any{1}, []string{"foo", "bar"}, []any{1, "foo", "bar"}},
		{[]any{1}, []any{"foo", "bar"}, []any{1, "foo", "bar"}},
		{[]string{"a"}, []any{[]string{"foo", "bar"}}, []any{"a", []string{"foo", "bar"}}},
		{[]any{"a"}, []any{[]string{"foo"}, []string{"bar"}}, []any{"a", []string{"foo"}, []string{"bar"}}},
		{[]string{"a"}, []any{[]string{"foo"}, []string{"bar"}}, []any{"a", []string{"foo"}, []string{"bar"}}},
		{[]any{"a"}, []any{[]string{"foo"}, []string{"bar"}}, []any{"a", []string{"foo"}, []string{"bar"}}},
		{[]int{1}, nil, []any{1}},
		{nil, []int{1}, []any{1}},
		{nil, nil, nil},
	} {
		slice1 := reflect.ValueOf(test.slice1)
		slice2 := reflect.ValueOf(test.slice2)
		expect := test.expect

		result, err := appendToInterfaceSliceFromValues(slice1, slice2)
		c.Assert(err, qt.IsNil)
		c.Assert(result, qt.DeepEquals, expect)
	}
}

func TestIndirect(t *testing.T) {
	t.Parallel()

	c := qt.New(t)

	v, isNil := indirect(reflect.ValueOf(nil))
	c.Assert(isNil, qt.Equals, true)
	c.Assert(v.IsValid(), qt.Equals, false)

	var (
		p   *int
		i   interface{}
		i2  interface{} = p
		i3  interface{} = i2
		i4  interface{} = i3
		i5  interface{} = i4
		i6  interface{} = i5
		i7  interface{} = i6
		i8  interface{} = i7
		i9  interface{} = i8
		i10 interface{} = i9
	)

	v, isNil = indirect(reflect.ValueOf(i10))
	c.Assert(isNil, qt.Equals, true)
	c.Assert(v.IsValid(), qt.Equals, false)
}