package glob

import (
	"testing"

	"github.com/gobwas/glob"
	"github.com/stretchr/testify/assert"
)

func TestGetGlob(t *testing.T) {
	tests := []struct {
		pattern string
		input   string
		match   bool
	}{
		{"*.txt", "file.txt", true},
		{"*.txt", "file.md", false},
		{"file.*", "file.txt", true},
		{"file.*", "file.md", true},
		{"file.*", "folder/file.txt", false},
	}

	for _, test := range tests {
		g, err := GetGlob(test.pattern)
		assert.NoError(t, err)
		assert.Equal(t, test.match, g.Match(test.input))
	}
}

func TestOr(t *testing.T) {
	g1, _ := GetGlob("*.txt")
	g2, _ := GetGlob("*.md")
	orGlob := Or(g1, g2)

	assert.True(t, orGlob.Match("file.txt"))
	assert.True(t, orGlob.Match("file.md"))
	assert.False(t, orGlob.Match("file.pdf"))
}

func TestNormalizePath(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"a//b/c", "a/b/c"},
		{"a/b/../c", "a/c"},
		{"./a/b/c/", "a/b/c"},
		{"A/B/C", "a/b/c"},
	}

	for _, test := range tests {
		assert.Equal(t, test.expected, NormalizePath(test.input))
	}
}

func TestNormalizePathNoLower(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"a//b/c", "a/b/c"},
		{"a/b/../c", "a/c"},
		{"./a/b/c/", "a/b/c"},
		{"A/B/C", "A/B/C"},
	}

	for _, test := range tests {
		assert.Equal(t, test.expected, NormalizePathNoLower(test.input))
	}
}

func TestResolveRootDir(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"assets/**/*.json", "assets"},
		{"assets/*.json", "assets"},
		{"assets/", "assets"},
		{"*.json", ""},
	}

	for _, test := range tests {
		assert.Equal(t, test.expected, ResolveRootDir(test.input))
	}
}

func TestFilterGlobParts(t *testing.T) {
	input := []string{"assets", "*.json", "images", "**"}
	expected := []string{"assets", "images"}
	result := FilterGlobParts(input)
	assert.Equal(t, expected, result)
}

func TestHasGlobChar(t *testing.T) {
	assert.True(t, HasGlobChar("*"))
	assert.True(t, HasGlobChar("?"))
	assert.False(t, HasGlobChar("abc"))
}