package types

import (
	"encoding/json"
	"html/template"
	"testing"
	"time"

	"github.com/spf13/cast"
	"github.com/stretchr/testify/assert"
)

func TestToDuration(t *testing.T) {
	tests := []struct {
		input    any
		expected time.Duration
	}{
		{1000, 1000 * time.Millisecond},
		{"2s", 2 * time.Second},
		{"invalid", 0},
	}

	for _, test := range tests {
		result := ToDuration(test.input)
		assert.Equal(t, test.expected, result)
	}
}

func TestToDurationE(t *testing.T) {
	tests := []struct {
		input          any
		expected       time.Duration
		expectingError bool
	}{
		{1000, 1000 * time.Millisecond, false},
		{"2s", 2 * time.Second, false},
		{"invalid", 0, true},
	}

	for _, test := range tests {
		result, err := ToDurationE(test.input)
		if test.expectingError {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, test.expected, result)
		}
	}
}

func TestToStringSlicePreserveString(t *testing.T) {
	tests := []struct {
		input    any
		expected []string
	}{
		{"single", []string{"single"}},
		{[]string{"a", "b"}, []string{"a", "b"}},
		{[]int{1, 2}, []string{"1", "2"}},
		{nil, nil},
	}

	for _, test := range tests {
		result := ToStringSlicePreserveString(test.input)
		assert.Equal(t, test.expected, result)
	}
}

func TestToStringSlicePreserveStringE(t *testing.T) {
	tests := []struct {
		input          any
		expected       []string
		expectingError bool
	}{
		{"single", []string{"single"}, false},
		{[]string{"a", "b"}, []string{"a", "b"}, false},
		{[]int{1, 2}, []string{"1", "2"}, false},
		{nil, nil, false},
		{123, nil, true},
	}

	for _, test := range tests {
		result, err := ToStringSlicePreserveStringE(test.input)
		if test.expectingError {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, test.expected, result)
		}
	}
}

func TestTypeToString(t *testing.T) {
	tests := []struct {
		input    any
		expected string
		valid    bool
	}{
		{"string", "string", true},
		{template.HTML("html"), "html", true},
		{123, "", false},
	}

	for _, test := range tests {
		result, valid := TypeToString(test.input)
		assert.Equal(t, test.expected, result)
		assert.Equal(t, test.valid, valid)
	}
}

func TestToString(t *testing.T) {
	tests := []struct {
		input    any
		expected string
	}{
		{"string", "string"},
		{123, "123"},
		{json.RawMessage(`"json"`), `"json"`},
	}

	for _, test := range tests {
		result := ToString(test.input)
		assert.Equal(t, test.expected, result)
	}
}

func TestToStringE(t *testing.T) {
	tests := []struct {
		input          any
		expected       string
		expectingError bool
	}{
		{"string", "string", false},
		{123, "123", false},
		{json.RawMessage(`"json"`), `"json"`, false},
		{[]int{1, 2}, "", true},
	}

	for _, test := range tests {
		result, err := ToStringE(test.input)
		if test.expectingError {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, test.expected, result)
		}
	}
}