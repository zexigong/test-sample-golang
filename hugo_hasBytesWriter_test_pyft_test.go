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
	"testing"
)

func TestHasBytesWriter(t *testing.T) {
	patterns := []*HasBytesPattern{
		{Pattern: []byte("123")},
		{Pattern: []byte("abc")},
	}

	w := &HasBytesWriter{
		Patterns: patterns,
	}

	n, err := w.Write([]byte("123"))
	if err != nil {
		t.Fatal(err)
	}
	if n != 3 {
		t.Fatalf("Got %d expected %d", n, 3)
	}
	if !patterns[0].Match {
		t.Fatalf("Got false expected true")
	}

	w = &HasBytesWriter{
		Patterns: patterns,
	}

	n, err = w.Write([]byte("a"))
	if err != nil {
		t.Fatal(err)
	}
	if n != 1 {
		t.Fatalf("Got %d expected %d", n, 1)
	}
	if patterns[1].Match {
		t.Fatalf("Got true expected false")
	}

	n, err = w.Write([]byte("bc"))
	if err != nil {
		t.Fatal(err)
	}
	if n != 2 {
		t.Fatalf("Got %d expected %d", n, 2)
	}
	if !patterns[1].Match {
		t.Fatalf("Got false expected true")
	}
}