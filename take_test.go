package iterlinq

import "testing"

// TestTake tests the Take operation.
func TestTake(t *testing.T) {
	// Define test cases using handlers.
	handlers := []struct {
		name     string
		input    []int
		take     int
		expected []int
	}{
		{
			name:     "TakeThree",
			input:    []int{1, 2, 3, 4, 5},
			take:     3,
			expected: []int{1, 2, 3},
		},
		{
			name:     "TakeMoreThanAvailable",
			input:    []int{1, 2},
			take:     5,
			expected: []int{1, 2},
		},
		{
			name:     "TakeNone",
			input:    []int{1, 2},
			take:     0,
			expected: []int{},
		},
		{
			name:     "TakeNegative",
			input:    []int{1, 2, 3},
			take:     -1,
			expected: []int{},
		},
	}

	for _, h := range handlers {
		t.Run(h.name, func(t *testing.T) {
			res, err := FromSlice(h.input).Take(h.take).ToSlice()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(res) != len(h.expected) {
				t.Fatalf("expected length %d, got %d", len(h.expected), len(res))
			}
			for i := range res {
				if res[i] != h.expected[i] {
					t.Errorf("at index %d: expected %d, got %d", i, h.expected[i], res[i])
				}
			}
		})
	}
}
