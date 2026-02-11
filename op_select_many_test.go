package iterlinq

import "testing"

// TestSelectMany tests the SelectMany operation.
func TestSelectMany(t *testing.T) {
	// Define test cases using testCases.
	testCases := []struct {
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

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := FromSlice(tc.input).SelectMany(func(n int) Sequence[any] {
				return FromSlice([]any{n, n * 10})
			}).ToSlice()

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(res) != len(tc.expected) {
				t.Fatalf("expected length %d, got %d", len(tc.expected), len(res))
			}
			for i := range res {
				if res[i] != tc.expected[i] {
					t.Errorf("at index %d: expected %v, got %v", i, tc.expected[i], res[i])
				}
			}
		})
	}
}
