package iterlinq

import (
	"fmt"
	"testing"
)

// TestAny tests the Any quantifier operation.
func TestAny(t *testing.T) {
	testCases := []struct {
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

			res, err := s.Any(func(n int) bool { return n%2 == 0 })
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
