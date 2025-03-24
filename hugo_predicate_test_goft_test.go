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

package predicate

import (
	"fmt"
	"reflect"
	"testing"
)

func TestPredicate(t *testing.T) {
	isOdd := P[int](func(v int) bool {
		return v%2 != 0
	})
	isEven := isOdd.Negate()

	isEvenAndGreaterThanFive := isEven.And(func(v int) bool {
		return v > 5
	})

	isOddOrGreaterThanFive := isOdd.Or(func(v int) bool {
		return v > 5
	})

	for i, test := range []struct {
		predicate P[int]
		input     int
		expected  bool
	}{
		{isOdd, 1, true},
		{isOdd, 2, false},
		{isEven, 2, true},
		{isEven, 3, false},
		{isEvenAndGreaterThanFive, 2, false},
		{isEvenAndGreaterThanFive, 7, false},
		{isEvenAndGreaterThanFive, 6, true},
		{isOddOrGreaterThanFive, 2, false},
		{isOddOrGreaterThanFive, 7, true},
		{isOddOrGreaterThanFive, 5, true},
	} {
		if got := test.predicate(test.input); got != test.expected {
			t.Errorf("[%d] got %v; want %v", i, got, test.expected)
		}
	}
}

func TestPredicateFilter(t *testing.T) {
	isOdd := P[int](func(v int) bool {
		return v%2 != 0
	})
	isEven := isOdd.Negate()

	for _, test := range []struct {
		filter   func([]int) []int
		input    []int
		expected []int
	}{
		{isOdd.Filter, []int{1, 2, 3, 4}, []int{1, 3}},
		{isEven.Filter, []int{1, 2, 3, 4}, []int{2, 4}},
		{isOdd.FilterCopy, []int{1, 2, 3, 4}, []int{1, 3}},
		{isEven.FilterCopy, []int{1, 2, 3, 4}, []int{2, 4}},
	} {
		name := fmt.Sprintf("%s-%v", nameOfFunc(test.filter), test.input)

		t.Run(name, func(t *testing.T) {
			if got := test.filter(test.input); !reflect.DeepEqual(got, test.expected) {
				t.Errorf("got %v; want %v", got, test.expected)
			}
		})
	}
}

func nameOfFunc(f interface{}) string {
	return reflect.ValueOf(f).Type().Name()
}