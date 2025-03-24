// Copyright 2023 The Hugo Authors. All rights reserved.
// Use of this source code is governed by a license that can be found in the LICENSE file.

package types

import (
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestEvictingQueue(t *testing.T) {
	t.Run("Generic", func(t *testing.T) {
		q := NewEvictingQueue[int](2)
		q.Add(1).Add(2).Add(3)

		qt.Assert(t, q.Len(), qt.Equals, 2)
		qt.Assert(t, q.PeekAll(), qt.DeepEquals, []int{3, 2})
		qt.Assert(t, q.PeekAllSet(), qt.DeepEquals, map[int]bool{2: true, 3: true})

		qt.Assert(t, q.Peek(), qt.Equals, 3)
		qt.Assert(t, q.Contains(3), qt.Equals, true)
		qt.Assert(t, q.Contains(4), qt.Equals, false)

		qt.Assert(t, NewEvictingQueue[int](2).PeekAll(), qt.HasLen, 0)
		qt.Assert(t, NewEvictingQueue[int](2).Peek(), qt.Equals, 0)
		qt.Assert(t, NewEvictingQueue[int](2).Contains(1), qt.Equals, false)
	})

	t.Run("String", func(t *testing.T) {
		q := NewEvictingQueue[string](2)
		q.Add("1").Add("2").Add("3")

		qt.Assert(t, q.Len(), qt.Equals, 2)
		qt.Assert(t, q.PeekAll(), qt.DeepEquals, []string{"3", "2"})
		qt.Assert(t, q.PeekAllSet(), qt.DeepEquals, map[string]bool{"2": true, "3": true})

		qt.Assert(t, q.Peek(), qt.Equals, "3")
		qt.Assert(t, q.Contains("3"), qt.Equals, true)
		qt.Assert(t, q.Contains("4"), qt.Equals, false)

		qt.Assert(t, NewEvictingQueue[string](2).PeekAll(), qt.HasLen, 0)
		qt.Assert(t, NewEvictingQueue[string](2).Peek(), qt.Equals, "")
		qt.Assert(t, NewEvictingQueue[string](2).Contains("1"), qt.Equals, false)
	})
}