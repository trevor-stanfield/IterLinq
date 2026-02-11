package iterlinq

import (
	"fmt"
	"testing"
)

// TestCount tests the Count operation.
func TestCount(t *testing.T) {
	// Define test cases using testCases.
	testCases := []struct {
		name        string
		input       []int
		failAt      int
		expected    int
		expectedErr string
	}{
		{
			name:     "FiveElements",
			input:    []int{1, 2, 3, 4, 5},
			expected: 5,
		},
		{
			name:     "EmptySequence",
			input:    []int{},
			expected: 0,
		},
		{
			name:        "ErrorPropagation",
			failAt:      3,
			expected:    2,
			expectedErr: "error in sequence",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var s Sequence[int]
			if tc.failAt != 0 {
				s = FromFunc(func(yield func(int) bool) error {
					for i := 1; i <= tc.failAt; i++ {
						if i == tc.failAt {
							return fmt.Errorf("error in sequence")
						}
						if !yield(i) {
							return nil
						}
					}
					return nil
				})
			} else {
				s = FromSlice(tc.input)
			}

			count, err := s.Count()

			if tc.expectedErr != "" {
				if err == nil || err.Error() != tc.expectedErr {
					t.Fatalf("expected error %q, got %v", tc.expectedErr, err)
				}
				if count != tc.expected {
					t.Errorf("expected partial count %d, got %d", tc.expected, count)
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if count != tc.expected {
					t.Errorf("expected count %d, got %d", tc.expected, count)
				}
			}
		})
	}
}
