package fasthttp

import (
	"testing"
)

func TestRoundUpForSliceCap(t *testing.T) {
	tests := []struct {
		input    int
		expected int
	}{
		{0, 0},
		{-1, 0},
		{1, 1},
		{2, 2},
		{3, 4},
		{4, 4},
		{5, 8},
		{15, 16},
		{16, 16},
		{17, 32},
		{33, 64},
		{65, 128},
		{129, 256},
		{257, 512},
		{513, 1024},
		{1025, 2048},
		{2049, 4096},
		{100 * 1024 * 1024, 100 * 1024 * 1024},
		{100*1024*1024 + 1, 100*1024*1024 + 1},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			result := roundUpForSliceCap(tt.input)
			if result != tt.expected {
				t.Errorf("roundUpForSliceCap(%d) = %d; want %d", tt.input, result, tt.expected)
			}
		})
	}
}