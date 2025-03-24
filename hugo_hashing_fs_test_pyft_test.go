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

package hugofs

import (
	"testing"

	"github.com/cespare/xxhash/v2"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

func TestHashingFs(t *testing.T) {
	fs := afero.NewMemMapFs()

	var filename string
	var checksum uint64

	hashFs := NewHashingFs(fs, &testHashReceiver{onFileClose: func(name string, csum uint64) {
		filename = name
		checksum = csum
	}})

	f, err := hashFs.Create("file.txt")
	require.NoError(t, err)

	_, err = f.Write([]byte("Hello world"))
	require.NoError(t, err)
	require.NoError(t, f.Close())

	require.Equal(t, "file.txt", filename)
	require.Equal(t, xxhash.Sum64String("Hello world"), checksum)
}

type testHashReceiver struct {
	onFileClose func(name string, csum uint64)
}

func (r *testHashReceiver) OnFileClose(name string, csum uint64) {
	r.onFileClose(name, csum)
}