package iterlinq

import (
	"fmt"
	"testing"
)

// TestAny tests the Any and AnyWithError quantifier operations.
func TestAny(t *testing.T) {
	// Define test cases using handlers.
	handlers := []struct {
		name        string
		input       []int
		failAt      int
		expected    bool
		expectedErr string
	}{
		{
			name:     "HasEven",
			input:    []int{1, 3, 5, 6},
			expected: true,
		},
		{
			name:     "NoneEven",
			input:    []int{1, 3, 5},
			expected: false,
		},
		{
			name:        "ErrorInjection",
			input:       []int{1, 3, 5},
			failAt:      3,
			expectedErr: "error at 3",
		},
	}

	for _, h := range handlers {
		t.Run(h.name, func(t *testing.T) {
			s := FromSlice(h.input)
			var res bool
			var err error

			if h.failAt != 0 {
				res, err = s.AnyWithError(func(n int) (bool, error) {
					if n == h.failAt {
						return false, fmt.Errorf("error at %d", n)
					}
					return n%2 == 0, nil
				})
			} else {
				res, err = s.Any(func(n int) bool {
					return n%2 == 0
				})
			}

			if h.expectedErr != "" {
				if err == nil || err.Error() != h.expectedErr {
					t.Fatalf("expected error %q, got %v", h.expectedErr, err)
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if res != h.expected {
					t.Errorf("expected %v, got %v", h.expected, res)
				}
			}
		})
	}
}
