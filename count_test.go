package iterlinq

import (
	"fmt"
	"testing"
)

// TestCount tests the Count operation.
func TestCount(t *testing.T) {
	// Define test cases using handlers.
	handlers := []struct {
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

	for _, h := range handlers {
		t.Run(h.name, func(t *testing.T) {
			var s Sequence[int]
			if h.failAt != 0 {
				s = From(func(yield func(int, error) bool) {
					for i := 1; i <= h.failAt; i++ {
						if i == h.failAt {
							yield(0, fmt.Errorf("error in sequence"))
							return
						}
						if !yield(i, nil) {
							return
						}
					}
				})
			} else {
				s = FromSlice(h.input)
			}

			count, err := s.Count()

			if h.expectedErr != "" {
				if err == nil || err.Error() != h.expectedErr {
					t.Fatalf("expected error %q, got %v", h.expectedErr, err)
				}
				if count != h.expected {
					t.Errorf("expected partial count %d, got %d", h.expected, count)
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if count != h.expected {
					t.Errorf("expected count %d, got %d", h.expected, count)
				}
			}
		})
	}
}
