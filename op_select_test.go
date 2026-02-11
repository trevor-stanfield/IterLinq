package iterlinq

import (
	"fmt"
	"testing"
)

// TestSelect tests the Select transformation operation.
func TestSelect(t *testing.T) {
	testCases := []struct {
		name        string
		input       []int
		function    func(int) any
		failAt      int
		expected    []any
		expectedErr string
	}{
		{
			name:     "DoubleNumbers",
			input:    []int{1, 2, 3},
			function: func(n int) any { return n * 2 },
			expected: []any{2, 4, 6},
		},
		{
			name:     "NumbersToStrings",
			input:    []int{1, 2, 3},
			function: func(n int) any { return fmt.Sprint(n) },
			expected: []any{"1", "2", "3"},
		},
		{
			name:        "ErrorInSource",
			input:       []int{1, 2, 3},
			function:    func(n int) any { return n * 2 },
			failAt:      2,
			expectedErr: "error at 2",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			s := FromSlice(tc.input)
			if tc.failAt != 0 {
				s = FromFunc(func(yield func(int) bool) error {
					for _, n := range tc.input {
						if n == tc.failAt {
							return fmt.Errorf("error at %d", n)
						}
						if !yield(n) {
							return nil
						}
					}
					return nil
				})
			}

			res, err := s.Select(tc.function).ToSlice()
			if tc.expectedErr != "" {
				if err == nil || err.Error() != tc.expectedErr {
					t.Fatalf("expected error %q, got %v", tc.expectedErr, err)
				}
				return
			}

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
