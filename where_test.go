package iterlinq

import (
	"fmt"
	"testing"
)

// TestWhere tests the Where and WhereWithError filtering operations.
func TestWhere(t *testing.T) {
	// Define test cases using handlers.
	handlers := []struct {
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

	for _, h := range handlers {
		t.Run(h.name, func(t *testing.T) {
			s := FromSlice(h.input)

			var res []int
			var err error

			if h.failAt != 0 {
				res, err = s.WhereWithError(func(n int) (bool, error) {
					if n == h.failAt {
						return false, fmt.Errorf("error at %d", n)
					}
					return true, nil
				}).ToSlice()
			} else {
				res, err = s.Where(func(n int) bool {
					return n%2 == 0
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
				if len(res) != len(h.expectedMatch) {
					t.Fatalf("expected length %d, got %d", len(h.expectedMatch), len(res))
				}
				for i := range res {
					if res[i] != h.expectedMatch[i] {
						t.Errorf("at index %d: expected %d, got %d", i, h.expectedMatch[i], res[i])
					}
				}
			}
		})
	}
}
