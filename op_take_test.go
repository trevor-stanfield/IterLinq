package iterlinq

import "testing"

// TestTake tests the Take operation.
func TestTake(t *testing.T) {
	// Define test cases using testCases.
	testCases := []struct {
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

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := FromSlice(tc.input).Take(tc.take).ToSlice()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(res) != len(tc.expected) {
				t.Fatalf("expected length %d, got %d", len(tc.expected), len(res))
			}
			for i := range res {
				if res[i] != tc.expected[i] {
					t.Errorf("at index %d: expected %d, got %d", i, tc.expected[i], res[i])
				}
			}
		})
	}
}
