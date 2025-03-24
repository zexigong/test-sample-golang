package glob

import (
	"testing"

	"github.com/gobwas/glob"
	"github.com/stretchr/testify/assert"
)

func TestNormalizeFilenameGlobPattern(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"path/to/file", "/path/to/file"},
		{"/path/to/file", "/path/to/file"},
		{"path\\to\\file", "/path/to/file"},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.expected, normalizeFilenameGlobPattern(tt.input))
	}
}

func TestNewFilenameFilter(t *testing.T) {
	inclusions := []string{"*.go", "docs/*"}
	exclusions := []string{"*.test", "build/*"}

	filter, err := NewFilenameFilter(inclusions, exclusions)
	assert.NoError(t, err)
	assert.NotNil(t, filter)
}

func TestNewFilenameFilter_Empty(t *testing.T) {
	filter, err := NewFilenameFilter(nil, nil)
	assert.NoError(t, err)
	assert.Nil(t, filter)
}

func TestMustNewFilenameFilter(t *testing.T) {
	inclusions := []string{"*.go", "docs/*"}
	exclusions := []string{"*.test", "build/*"}

	assert.NotPanics(t, func() {
		filter := MustNewFilenameFilter(inclusions, exclusions)
		assert.NotNil(t, filter)
	})
}

func TestNewFilenameFilterForInclusionFunc(t *testing.T) {
	filter := NewFilenameFilterForInclusionFunc(func(filename string) bool {
		return filename == "include.txt"
	})

	assert.NotNil(t, filter)
	assert.True(t, filter.Match("include.txt", false))
	assert.False(t, filter.Match("exclude.txt", false))
}

func TestFilenameFilter_Match(t *testing.T) {
	inclusions := []string{"*.go", "docs/*"}
	exclusions := []string{"*.test", "build/*"}

	filter, _ := NewFilenameFilter(inclusions, exclusions)

	assert.True(t, filter.Match("main.go", false))
	assert.False(t, filter.Match("main.test", false))
	assert.True(t, filter.Match("docs/index.md", false))
	assert.False(t, filter.Match("build/main.go", false))
}

func TestFilenameFilter_Append(t *testing.T) {
	filter1 := MustNewFilenameFilter([]string{"*.go"}, nil)
	filter2 := MustNewFilenameFilter(nil, []string{"*.test"})

	combinedFilter := filter1.Append(filter2)

	assert.True(t, combinedFilter.Match("main.go", false))
	assert.False(t, combinedFilter.Match("main.test", false))
}