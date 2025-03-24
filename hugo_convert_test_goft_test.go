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

package types

import (
	"encoding/json"
	"fmt"
	"html/template"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTypeToString(t *testing.T) {
	for i, test := range []struct {
		v      any
		expect string
	}{
		{"string", "string"},
		{template.HTML("html"), "html"},
		{template.CSS("css"), "css"},
		{template.HTMLAttr("htmlattr"), "htmlattr"},
		{template.JS("js"), "js"},
		{template.JSStr("jsstr"), "jsstr"},
		{template.URL("url"), "url"},
		{template.Srcset("srcset"), "srcset"},
	} {
		result, _ := TypeToString(test.v)
		assert.Equal(t, test.expect, result, fmt.Sprintf("Test %d", i))
	}
}

func TestToString(t *testing.T) {
	assert := assert.New(t)

	for i, test := range []struct {
		v      any
		expect string
	}{
		{"string", "string"},
		{template.HTML("html"), "html"},
		{template.CSS("css"), "css"},
		{template.HTMLAttr("htmlattr"), "htmlattr"},
		{template.JS("js"), "js"},
		{template.JSStr("jsstr"), "jsstr"},
		{template.URL("url"), "url"},
		{template.Srcset("srcset"), "srcset"},
		{json.RawMessage("raw"), "raw"},
	} {
		result, _ := ToStringE(test.v)
		assert.Equal(test.expect, result, fmt.Sprintf("Test %d", i))
	}
}

func TestToDuration(t *testing.T) {
	assert := assert.New(t)

	for i, test := range []struct {
		v      any
		expect time.Duration
	}{
		{"5s", 5 * time.Second},
		{5, 5 * time.Millisecond},
		{int32(5), 5 * time.Millisecond},
		{int64(5), 5 * time.Millisecond},
		{uint32(5), 5 * time.Millisecond},
		{uint64(5), 5 * time.Millisecond},
	} {
		result, err := ToDurationE(test.v)
		assert.NoError(err, fmt.Sprintf("Test %d", i))
		assert.Equal(test.expect, result, fmt.Sprintf("Test %d", i))
	}
}

func TestToStringSlicePreserveString(t *testing.T) {
	assert := assert.New(t)

	for i, test := range []struct {
		v      any
		expect []string
	}{
		{"str", []string{"str"}},
		{[]string{"a", "b"}, []string{"a", "b"}},
		{[]int{1, 2}, []string{"1", "2"}},
		{[]uint64{1, 2}, []string{"1", "2"}},
		{[]interface{}{1, "a"}, []string{"1", "a"}},
		{nil, nil},
	} {
		result, err := ToStringSlicePreserveStringE(test.v)
		assert.NoError(err, fmt.Sprintf("Test %d", i))
		assert.Equal(test.expect, result, fmt.Sprintf("Test %d", i))
	}
}