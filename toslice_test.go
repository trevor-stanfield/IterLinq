package iterlinq

import (
	"fmt"
	"testing"
)

// TestToSlice tests the materialization of a sequence into a slice.
func TestToSlice(t *testing.T) {
	// Define test cases using handlers.
	handlers := []struct {
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

	for _, h := range handlers {
		t.Run(h.name, func(t *testing.T) {
			var s Sequence[int]
			if h.failAt != 0 {
				s = From(func(yield func(int, error) bool) {
					yield(1, nil)
					yield(0, fmt.Errorf("error in sequence"))
				})
			} else {
				s = FromSlice(h.input)
			}

			res, err := s.ToSlice()

			if h.expectedErr != "" {
				if err == nil || err.Error() != h.expectedErr {
					t.Fatalf("expected error %q, got %v", h.expectedErr, err)
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if len(res) != h.expectedLen {
					t.Errorf("expected length %d, got %d", h.expectedLen, len(res))
				}
			}
		})
	}
}
