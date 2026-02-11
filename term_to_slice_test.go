package iterlinq

import (
	"fmt"
	"testing"
)

// TestToSlice tests the materialization of a sequence into a slice.
func TestToSlice(t *testing.T) {
	// Define test cases using testCases.
	testCases := []struct {
		name        string
		input       []int
		failAt      int
		expectedLen int
		expectedErr string
	}{
		{
			name:        "BasicSlice",
			input:       []int{1, 2, 3},
			expectedLen: 3,
		},
		{
			name:        "ErrorPropagation",
			failAt:      2,
			expectedErr: "error in sequence",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var s Sequence[int]
			if tc.failAt != 0 {
				s = FromFunc(func(yield func(int) bool) error {
					if !yield(1) {
						return nil
					}
					return fmt.Errorf("error in sequence")
				})
			} else {
				s = FromSlice(tc.input)
			}

			res, err := s.ToSlice()

			if tc.expectedErr != "" {
				if err == nil || err.Error() != tc.expectedErr {
					t.Fatalf("expected error %q, got %v", tc.expectedErr, err)
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if len(res) != tc.expectedLen {
					t.Errorf("expected length %d, got %d", tc.expectedLen, len(res))
				}
			}
		})
	}
}
