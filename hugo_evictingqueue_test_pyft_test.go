// Copyright 2017-present The Hugo Authors. All rights reserved.
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

package types

import (
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestEvictingQueue(t *testing.T) {
	c := qt.New(t)

	q := NewEvictingQueue[string](3)
	q.Add("a").Add("b").Add("c").Add("d")

	c.Assert(q.Len(), qt.Equals, 3)
	c.Assert(q.Peek(), qt.Equals, "d")

	q.Add("b").Add("e")

	c.Assert(q.Peek(), qt.Equals, "e")
	c.Assert(q.Len(), qt.Equals, 3)
	c.Assert(q.PeekAll(), qt.DeepEquals, []string{"e", "b", "d"})

	c.Assert(q.Contains("d"), qt.Equals, true)
	c.Assert(q.Contains("a"), qt.Equals, false)

	set := q.PeekAllSet()
	c.Assert(set["e"], qt.Equals, true)
	c.Assert(set["b"], qt.Equals, true)
	c.Assert(set["d"], qt.Equals, true)

	c.Assert(q.Contains("b"), qt.Equals, true)

	q = nil
	c.Assert(q.Len(), qt.Equals, 0)
	c.Assert(q.Peek(), qt.Equals, "")
	c.Assert(q.PeekAll(), qt.IsNil)
	c.Assert(q.Contains("b"), qt.Equals, false)
}