package iterlinq

import (
	"errors"
	"fmt"
	"testing"
)

// TestLast tests the Last family of operations.
func TestLast(t *testing.T) {
	testCases := []struct {
		name        string
		input       []int
		failAt      int
		useDefault  bool
		expected    int
		expectedErr error
	}{
		{
			name:     "MatchFound",
			input:    []int{1, 2, 3, 2},
			expected: 2,
		},
		{
			name:        "NoMatch",
			input:       []int{1, 3},
			expectedErr: ErrNoMatch,
		},
		{
			name:       "NoMatchDefault",
			input:      []int{1, 3},
			useDefault: true,
			expected:   0,
		},
		{
			name:        "ErrorInjection",
			input:       []int{1, 2, 3},
			failAt:      2,
			expectedErr: errors.New("error at 2"),
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

			var (
				res int
				err error
			)
			if tc.useDefault {
				res, err = s.LastOrDefault(func(n int) bool { return n%2 == 0 })
			} else {
				res, err = s.Last(func(n int) bool { return n%2 == 0 })
			}

			if tc.expectedErr != nil {
				if err == nil || err.Error() != tc.expectedErr.Error() {
					t.Fatalf("expected error %v, got %v", tc.expectedErr, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if res != tc.expected {
				t.Errorf("expected %d, got %d", tc.expected, res)
			}
		})
	}
}
