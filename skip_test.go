package iterlinq

import "testing"

// TestSkip tests the Skip operation.
func TestSkip(t *testing.T) {
	// Define test cases using handlers.
	handlers := []struct {
		name     string
		input    []int
		skip     int
		expected []int
	}{
		{
			name:     "SkipTwo",
			input:    []int{1, 2, 3, 4, 5},
			skip:     2,
			expected: []int{3, 4, 5},
		},
		{
			name:     "SkipAll",
			input:    []int{1, 2},
			skip:     5,
			expected: []int{},
		},
		{
			name:     "SkipNone",
			input:    []int{1, 2},
			skip:     0,
			expected: []int{1, 2},
		},
		{
			name:     "SkipNegative",
			input:    []int{1, 2, 3},
			skip:     -2,
			expected: []int{1, 2, 3},
		},
	}

	for _, h := range handlers {
		t.Run(h.name, func(t *testing.T) {
			res, err := FromSlice(h.input).Skip(h.skip).ToSlice()
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
