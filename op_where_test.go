package iterlinq

import (
	"fmt"
	"testing"
)

// TestWhere tests the Where filtering operation.
func TestWhere(t *testing.T) {
	// Define test cases using testCases.
	testCases := []struct {
		name          string
		input         []int
		failAt        int
		expectedMatch []int
		expectedErr   string
	}{
		{
			name:          "FilterEvens",
			input:         []int{1, 2, 3, 4, 5},
			expectedMatch: []int{2, 4},
		},
		{
			name:          "NoMatch",
			input:         []int{1, 3, 5},
			expectedMatch: []int{},
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

			res, err := s.Where(func(n int) bool {
				return n%2 == 0
			}).ToSlice()

			if tc.expectedErr != "" {
				if err == nil || err.Error() != tc.expectedErr {
					t.Fatalf("expected error %q, got %v", tc.expectedErr, err)
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if len(res) != len(tc.expectedMatch) {
					t.Fatalf("expected length %d, got %d", len(tc.expectedMatch), len(res))
				}
				for i := range res {
					if res[i] != tc.expectedMatch[i] {
						t.Errorf("at index %d: expected %d, got %d", i, tc.expectedMatch[i], res[i])
					}
				}
			}
		})
	}
}
