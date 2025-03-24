package types

import (
	"testing"
)

func TestNewEvictingQueue(t *testing.T) {
	q := NewEvictingQueue[int](3)
	if q == nil {
		t.Fatal("Expected non-nil queue")
	}
	if q.Len() != 0 {
		t.Fatalf("Expected length 0, got %d", q.Len())
	}
}

func TestEvictingQueue_Add(t *testing.T) {
	q := NewEvictingQueue[int](3)
	q.Add(1).Add(2).Add(3)
	if q.Len() != 3 {
		t.Fatalf("Expected length 3, got %d", q.Len())
	}

	q.Add(4)
	if q.Len() != 3 {
		t.Fatalf("Expected length 3 after eviction, got %d", q.Len())
	}

	if !q.Contains(4) || q.Contains(1) {
		t.Fatal("Expected queue to contain 4 and not contain 1 after eviction")
	}
}

func TestEvictingQueue_Contains(t *testing.T) {
	q := NewEvictingQueue[int](3)
	q.Add(1).Add(2)
	if !q.Contains(1) {
		t.Fatal("Expected queue to contain 1")
	}
	if q.Contains(3) {
		t.Fatal("Expected queue not to contain 3")
	}
}

func TestEvictingQueue_Peek(t *testing.T) {
	q := NewEvictingQueue[int](3)
	if q.Peek() != 0 {
		t.Fatal("Expected zero value when peeking empty queue")
	}
	q.Add(1).Add(2)
	if q.Peek() != 2 {
		t.Fatalf("Expected to peek 2, got %d", q.Peek())
	}
}

func TestEvictingQueue_PeekAll(t *testing.T) {
	q := NewEvictingQueue[int](3)
	if vals := q.PeekAll(); len(vals) != 0 {
		t.Fatalf("Expected empty slice, got %v", vals)
	}
	q.Add(1).Add(2).Add(3)
	expected := []int{3, 2, 1}
	if vals := q.PeekAll(); !equal(vals, expected) {
		t.Fatalf("Expected %v, got %v", expected, vals)
	}
}

func TestEvictingQueue_PeekAllSet(t *testing.T) {
	q := NewEvictingQueue[int](3)
	q.Add(1).Add(2).Add(3)
	set := q.PeekAllSet()
	if len(set) != 3 {
		t.Fatalf("Expected set of length 3, got %d", len(set))
	}
	for _, v := range []int{1, 2, 3} {
		if !set[v] {
			t.Fatalf("Expected set to contain %d", v)
		}
	}
}

// Helper function to check equality of two slices
func equal(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}