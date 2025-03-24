package collections

import (
	"reflect"
	"testing"
)

func TestAppend(t *testing.T) {
	tests := []struct {
		to      any
		from    []any
		want    any
		wantErr bool
	}{
		{nil, []any{}, nil, false},
		{[]int{}, []any{1, 2, 3}, []int{1, 2, 3}, false},
		{[]int{1, 2}, []any{3, 4}, []int{1, 2, 3, 4}, false},
		{[]string{}, []any{"a", "b"}, []string{"a", "b"}, false},
		{[]interface{}{1, "a"}, []any{2, "b"}, []interface{}{1, "a", 2, "b"}, false},
		{[]int{1}, []any{[]int{2, 3}}, []int{1, 2, 3}, false},
		{[]interface{}{}, []any{nil, "a"}, []interface{}{nil, "a"}, false},
		{[]int{}, []any{[]string{"a"}}, nil, true},
		{[]string{}, []any{nil}, []interface{}{nil}, false},
		{nil, []any{1, 2, 3}, []interface{}{1, 2, 3}, false},
		{[]int(nil), []any{1}, []int{1}, false},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			got, err := Append(tt.to, tt.from...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Append() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Append() = %v, want %v", got, tt.want)
			}
		})
	}
}