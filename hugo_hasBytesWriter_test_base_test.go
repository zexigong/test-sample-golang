package hugio_test

import (
	"bytes"
	"testing"

	"github.com/gohugoio/hugo/hugio"
)

func TestHasBytesWriter(t *testing.T) {
	tests := []struct {
		name      string
		patterns  [][]byte
		input     []byte
		expectDone bool
	}{
		{
			name: "Single pattern match",
			patterns: [][]byte{
				[]byte("hello"),
			},
			input: []byte("hello world"),
			expectDone: true,
		},
		{
			name: "Multiple patterns match",
			patterns: [][]byte{
				[]byte("foo"),
				[]byte("bar"),
			},
			input: []byte("foo bar baz"),
			expectDone: true,
		},
		{
			name: "No pattern match",
			patterns: [][]byte{
				[]byte("foo"),
				[]byte("bar"),
			},
			input: []byte("baz qux"),
			expectDone: false,
		},
		{
			name: "Partial pattern match",
			patterns: [][]byte{
				[]byte("foo"),
				[]byte("bar"),
			},
			input: []byte("foo baz"),
			expectDone: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			patterns := make([]*hugio.HasBytesPattern, len(tt.patterns))
			for i, p := range tt.patterns {
				patterns[i] = &hugio.HasBytesPattern{Pattern: p}
			}

			writer := &hugio.HasBytesWriter{
				Patterns: patterns,
			}

			n, err := writer.Write(tt.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if n != len(tt.input) {
				t.Fatalf("expected %d bytes written, got %d", len(tt.input), n)
			}

			if writer.done != tt.expectDone {
				t.Fatalf("expected done: %v, got %v", tt.expectDone, writer.done)
			}

			for i, p := range patterns {
				if bytes.Contains(tt.input, p.Pattern) != p.Match {
					t.Errorf("pattern %d match expected %v, got %v", i, bytes.Contains(tt.input, p.Pattern), p.Match)
				}
			}
		})
	}
}