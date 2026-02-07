package iterlinq

import (
	"fmt"
	"testing"
)

// TestAll tests the All quantifier operation.
func TestAll(t *testing.T) {
	// Define test cases using handlers.
	handlers := []struct {
		name        string
		input       []int
		failAt      int
		expected    bool
		expectedErr string
	}{
		{
			name:     "AllPositive",
			input:    []int{1, 2, 3},
			expected: true,
		},
		{
			name:     "SomeNegative",
			input:    []int{1, -2, 3},
			expected: false,
		},
		{
			name:        "ErrorInjection",
			input:       []int{1, 2, 3},
			failAt:      2,
			expectedErr: "error at 2",
		},
	}

	for _, h := range handlers {
		t.Run(h.name, func(t *testing.T) {
			s := FromSlice(h.input)
			res, err := s.All(func(n int) (bool, error) {
				if h.failAt != 0 && n == h.failAt {
					return false, fmt.Errorf("error at %d", n)
				}
				return n > 0, nil
			})

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
