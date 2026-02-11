package iterlinq

import (
	"fmt"
	"testing"
)

// TestAll tests the All quantifier operation.
func TestAll(t *testing.T) {
	testCases := []struct {
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

			res, err := s.All(func(n int) bool { return n > 0 })
			if tc.expectedErr != "" {
				if err == nil || err.Error() != tc.expectedErr {
					t.Fatalf("expected error %q, got %v", tc.expectedErr, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if res != tc.expected {
				t.Errorf("expected %v, got %v", tc.expected, res)
			}
		})
	}
}
