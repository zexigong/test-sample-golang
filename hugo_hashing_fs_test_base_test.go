package hugofs

import (
	"bytes"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

// MockReceiver is a mock implementation of FileHashReceiver for testing.
type MockReceiver struct {
	received map[string]uint64
}

func (m *MockReceiver) OnFileClose(name string, checksum uint64) {
	m.received[name] = checksum
}

func TestNewHashingFs(t *testing.T) {
	mockReceiver := &MockReceiver{received: make(map[string]uint64)}
	fs := NewHashingFs(afero.NewMemMapFs(), mockReceiver)

	assert.Equal(t, "hashingFs", fs.Name())
}

func TestHashingFs_CreateAndWrite(t *testing.T) {
	mockReceiver := &MockReceiver{received: make(map[string]uint64)}
	fs := NewHashingFs(afero.NewMemMapFs(), mockReceiver)

	file, err := fs.Create("testfile")
	assert.NoError(t, err)
	defer file.Close()

	content := []byte("Hello, World!")
	n, err := file.Write(content)
	assert.NoError(t, err)
	assert.Equal(t, len(content), n)

	err = file.Close()
	assert.NoError(t, err)

	checksum, ok := mockReceiver.received["testfile"]
	assert.True(t, ok)
	assert.NotZero(t, checksum)
}

func TestHashingFs_OpenFileAndWrite(t *testing.T) {
	mockReceiver := &MockReceiver{received: make(map[string]uint64)}
	fs := NewHashingFs(afero.NewMemMapFs(), mockReceiver)

	// Create and write to the file first
	file, err := fs.Create("testfile")
	assert.NoError(t, err)
	content := []byte("Hello, World!")
	file.Write(content)
	file.Close()

	// Open the file in write mode and write more content
	file, err = fs.OpenFile("testfile", os.O_WRONLY|os.O_APPEND, 0644)
	assert.NoError(t, err)

	moreContent := []byte(" More content.")
	n, err := file.Write(moreContent)
	assert.NoError(t, err)
	assert.Equal(t, len(moreContent), n)

	err = file.Close()
	assert.NoError(t, err)

	checksum, ok := mockReceiver.received["testfile"]
	assert.True(t, ok)
	assert.NotZero(t, checksum)
}