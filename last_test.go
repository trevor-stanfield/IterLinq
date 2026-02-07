package iterlinq

import (
	"errors"
	"fmt"
	"testing"
)

// TestLast tests the Last family of operations.
func TestLast(t *testing.T) {
	// Define test cases using handlers.
	handlers := []struct {
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

	for _, h := range handlers {
		t.Run(h.name, func(t *testing.T) {
			s := FromSlice(h.input)
			var res int
			var err error

			predicate := func(n int) (bool, error) {
				if h.failAt != 0 && n == h.failAt {
					return false, fmt.Errorf("error at %d", n)
				}
				return n%2 == 0, nil
			}

			if h.useDefault {
				res, err = s.LastOrDefaultWithError(predicate)
			} else {
				res, err = s.LastWithError(predicate)
			}

			if h.expectedErr != nil {
				if err == nil || err.Error() != h.expectedErr.Error() {
					t.Fatalf("expected error %v, got %v", h.expectedErr, err)
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if res != h.expected {
					t.Errorf("expected %d, got %d", h.expected, res)
				}
			}
		})
	}
}
