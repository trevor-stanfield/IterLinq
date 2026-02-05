package iterlinq

import "testing"

// TestSelectMany tests the SelectMany operation.
func TestSelectMany(t *testing.T) {
	// Define test cases using handlers.
	handlers := []struct {
		name     string
		input    []int
		expected []any
	}{
		{
			name:     "FlattenNumbers",
			input:    []int{1, 2},
			expected: []any{1, 10, 2, 20},
		},
		{
			name:     "EmptyInput",
			input:    []int{},
			expected: []any{},
		},
	}

	for _, h := range handlers {
		t.Run(h.name, func(t *testing.T) {
			res, err := FromSlice(h.input).SelectMany(func(n int) Sequence[any] {
				return FromSlice([]any{n, n * 10})
			}).ToSlice()

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(res) != len(h.expected) {
				t.Fatalf("expected length %d, got %d", len(h.expected), len(res))
			}
			for i := range res {
				if res[i] != h.expected[i] {
					t.Errorf("at index %d: expected %v, got %v", i, h.expected[i], res[i])
				}
			}
		})
	}
}
