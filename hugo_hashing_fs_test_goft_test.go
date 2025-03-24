// Copyright 2019 The Hugo Authors. All rights reserved.
// Use of this source code is governed by a license that can be found in the LICENSE file.

package hugofs

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fileHashReceiver struct {
	hashes map[string]uint64
}

func (r *fileHashReceiver) OnFileClose(name string, checksum uint64) {
	if r.hashes == nil {
		r.hashes = make(map[string]uint64)
	}
	r.hashes[filepath.Base(name)] = checksum
}

func TestHashingFs(t *testing.T) {
	r := &fileHashReceiver{}
	fs := NewHashingFs(afero.NewMemMapFs(), r)

	f1, err := fs.Create("f1.txt")
	require.NoError(t, err)
	f2, err := fs.Create("f2.txt")
	require.NoError(t, err)

	_, err = f1.Write([]byte("f1 content"))
	require.NoError(t, err)
	_, err = f2.Write([]byte("f2 content"))
	require.NoError(t, err)

	require.NoError(t, f1.Close())
	require.NoError(t, f2.Close())

	assert.Equal(t, uint64(0x5cbacb1d), r.hashes["f1.txt"])
	assert.Equal(t, uint64(0x3aa8ff90), r.hashes["f2.txt"])

	f1, err = fs.OpenFile("f1.txt", os.O_RDWR, 0)
	require.NoError(t, err)

	_, err = f1.Write([]byte("f1 content"))
	require.NoError(t, err)

	require.NoError(t, f1.Close())
	assert.Equal(t, uint64(0x5cbacb1d), r.hashes["f1.txt"])
}