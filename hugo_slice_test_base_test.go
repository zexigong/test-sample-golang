package collections

import (
	"reflect"
	"testing"
)

type mockSlicer struct {
	sliceFunc func(items any) (any, error)
}

func (m mockSlicer) Slice(items any) (any, error) {
	return m.sliceFunc(items)
}

func TestSlice(t *testing.T) {
	tests := []struct {
		name     string
		input    []any
		expected any
	}{
		{
			name:     "Empty input",
			input:    []any{},
			expected: []any{},
		},
		{
			name:     "Single type input",
			input:    []any{1, 2, 3},
			expected: []int{1, 2, 3},
		},
		{
			name:     "Mixed type input",
			input:    []any{1, "two", 3},
			expected: []any{1, "two", 3},
		},
		{
			name: "Slicer interface input",
			input: []any{
				mockSlicer{sliceFunc: func(items any) (any, error) {
					return []string{"a", "b"}, nil
				}},
			},
			expected: []string{"a", "b"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Slice(tt.input...)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestStringSliceToInterfaceSlice(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected []any
	}{
		{
			name:     "Empty slice",
			input:    []string{},
			expected: []any{},
		},
		{
			name:     "Non-empty slice",
			input:    []string{"a", "b"},
			expected: []any{"a", "b"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := StringSliceToInterfaceSlice(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestSortedStringSliceContains(t *testing.T) {
	ss := SortedStringSlice{"a", "b", "c"}

	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "Contains string",
			input:    "b",
			expected: true,
		},
		{
			name:     "Does not contain string",
			input:    "d",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ss.Contains(tt.input)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestSortedStringSliceCount(t *testing.T) {
	ss := SortedStringSlice{"a", "b", "b", "c"}

	tests := []struct {
		name     string
		input    string
		expected int
	}{
		{
			name:     "Count existing string",
			input:    "b",
			expected: 2,
		},
		{
			name:     "Count non-existing string",
			input:    "d",
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ss.Count(tt.input)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}