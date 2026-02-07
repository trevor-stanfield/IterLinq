package iterlinq

import (
	"fmt"
	"testing"
)

// TestSelect tests the Select and SelectWithError transformation operations.
func TestSelect(t *testing.T) {
	// Define test cases using handlers.
	handlers := []struct {
		name        string
		input       []int
		failAt      int
		expected    []any
		expectedErr string
	}{
		{
			name:     "DoubleNumbers",
			input:    []int{1, 2, 3},
			expected: []any{2, 4, 6},
		},
		{
			name:        "ErrorInTransform",
			input:       []int{1, 2, 3},
			failAt:      2,
			expectedErr: "error at 2",
		},
	}

	for _, h := range handlers {
		t.Run(h.name, func(t *testing.T) {
			s := FromSlice(h.input)
			var res []any
			var err error

			if h.failAt != 0 {
				res, err = s.SelectWithError(func(n int) (any, error) {
					if n == h.failAt {
						return nil, fmt.Errorf("error at %d", n)
					}
					return n * 2, nil
				}).ToSlice()
			} else {
				res, err = s.Select(func(n int) any {
					return n * 2
				}).ToSlice()
			}

			if h.expectedErr != "" {
				if err == nil || err.Error() != h.expectedErr {
					t.Fatalf("expected error %q, got %v", h.expectedErr, err)
				}
			} else {
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
			}
		})
	}
}
