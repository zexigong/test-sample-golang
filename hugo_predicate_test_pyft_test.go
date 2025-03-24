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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPredicate(t *testing.T) {
	p := P[int](func(v int) bool {
		return v > 0
	})

	p2 := P[int](func(v int) bool {
		return v < 10
	})

	assert.True(t, p(1))
	assert.False(t, p(-1))
	assert.True(t, p.And(p2)(5))
	assert.False(t, p.And(p2)(-5))
	assert.False(t, p.And(p2)(15))
	assert.True(t, p.Or(p2)(15))
	assert.False(t, p.Or(p2)(-15))
	assert.True(t, p.Negate()(-1))
	assert.False(t, p.Negate()(1))
	assert.Equal(t, []int{1, 2, 3}, p.And(p2).Filter([]int{1, 2, 3, -1, 11}))
	assert.Equal(t, []int{1, 2, 3}, p.And(p2).FilterCopy([]int{1, 2, 3, -1, 11}))
}