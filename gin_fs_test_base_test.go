package gin

import (
	"net/http"
	"os"
	"testing"
	"testing/fstest"
)

func TestOnlyFilesFS_Open(t *testing.T) {
	fs := &OnlyFilesFS{FileSystem: fstest.MapFS{
		"file.txt": &fstest.MapFile{
			Data: []byte("hello, world"),
		},
	}}

	file, err := fs.Open("file.txt")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	defer file.Close()

	buffer := make([]byte, 12)
	_, err = file.Read(buffer)
	if err != nil {
		t.Fatalf("expected no error reading file, got %v", err)
	}
	expected := "hello, world"
	if string(buffer) != expected {
		t.Fatalf("expected %s, got %s", expected, string(buffer))
	}

	_, err = file.Readdir(0)
	if err != nil {
		t.Fatalf("expected no error from Readdir, got %v", err)
	}
}

func TestOnlyFilesFS_OpenNonExistentFile(t *testing.T) {
	fs := &OnlyFilesFS{FileSystem: fstest.MapFS{}}

	_, err := fs.Open("nonexistent.txt")
	if err == nil {
		t.Fatal("expected error for non-existent file, got nil")
	}
}

func TestDir(t *testing.T) {
	// Test with listDirectory = true
	fs := Dir(".", true)
	_, ok := fs.(http.Dir)
	if !ok {
		t.Fatal("expected http.Dir when listDirectory is true")
	}

	// Test with listDirectory = false
	fs = Dir(".", false)
	_, ok = fs.(*OnlyFilesFS)
	if !ok {
		t.Fatal("expected OnlyFilesFS when listDirectory is false")
	}
}