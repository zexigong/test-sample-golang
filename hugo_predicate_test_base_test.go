package predicate_test

import (
	"testing"

	"github.com/your-repo-name/hugo/predicate"
)

func TestAnd(t *testing.T) {
	greaterThanTen := predicate.P[int](func(i int) bool { return i > 10 })
	lessThanTwenty := predicate.P[int](func(i int) bool { return i < 20 })

	andPredicate := greaterThanTen.And(lessThanTwenty)

	tests := []struct {
		value    int
		expected bool
	}{
		{5, false},
		{15, true},
		{25, false},
	}

	for _, test := range tests {
		result := andPredicate(test.value)
		if result != test.expected {
			t.Errorf("And predicate failed for value %d, expected %v, got %v", test.value, test.expected, result)
		}
	}
}

func TestOr(t *testing.T) {
	lessThanTen := predicate.P[int](func(i int) bool { return i < 10 })
	greaterThanTwenty := predicate.P[int](func(i int) bool { return i > 20 })

	orPredicate := lessThanTen.Or(greaterThanTwenty)

	tests := []struct {
		value    int
		expected bool
	}{
		{5, true},
		{15, false},
		{25, true},
	}

	for _, test := range tests {
		result := orPredicate(test.value)
		if result != test.expected {
			t.Errorf("Or predicate failed for value %d, expected %v, got %v", test.value, test.expected, result)
		}
	}
}

func TestNegate(t *testing.T) {
	isEven := predicate.P[int](func(i int) bool { return i%2 == 0 })

	notEven := isEven.Negate()

	tests := []struct {
		value    int
		expected bool
	}{
		{2, false},
		{3, true},
	}

	for _, test := range tests {
		result := notEven(test.value)
		if result != test.expected {
			t.Errorf("Negate predicate failed for value %d, expected %v, got %v", test.value, test.expected, result)
		}
	}
}

func TestFilter(t *testing.T) {
	isEven := predicate.P[int](func(i int) bool { return i%2 == 0 })

	slice := []int{1, 2, 3, 4, 5, 6}
	filtered := isEven.Filter(slice)

	expected := []int{2, 4, 6}

	if len(filtered) != len(expected) {
		t.Errorf("Filter length mismatch, expected %d, got %d", len(expected), len(filtered))
	}

	for i, v := range expected {
		if filtered[i] != v {
			t.Errorf("Filter failed at index %d, expected %d, got %d", i, v, filtered[i])
		}
	}
}

func TestFilterCopy(t *testing.T) {
	isOdd := predicate.P[int](func(i int) bool { return i%2 != 0 })

	slice := []int{1, 2, 3, 4, 5, 6}
	filteredCopy := isOdd.FilterCopy(slice)

	expected := []int{1, 3, 5}

	if len(filteredCopy) != len(expected) {
		t.Errorf("FilterCopy length mismatch, expected %d, got %d", len(expected), len(filteredCopy))
	}

	for i, v := range expected {
		if filteredCopy[i] != v {
			t.Errorf("FilterCopy failed at index %d, expected %d, got %d", i, v, filteredCopy[i])
		}
	}
}