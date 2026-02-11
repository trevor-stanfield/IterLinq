package iterlinq

import "testing"

// TestDistinctBy tests the DistinctBy operation.
func TestDistinctBy(t *testing.T) {
	type person struct {
		id   int
		name string
	}

	// Define test cases using testCases.
	testCases := []struct {
		name        string
		input       []person
		useInvalid  bool
		expected    []int
		expectedErr error
	}{
		{
			name: "DistinctIDs",
			input: []person{
				{1, "Alice"},
				{2, "Bob"},
				{1, "Alice Duplicate"},
			},
			expected: []int{1, 2},
		},
		{
			name:        "NonComparableKey",
			useInvalid:  true,
			expectedErr: ErrNonComparableKey,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.useInvalid {
				nums := []int{1, 2}
				_, err := FromSlice(nums).DistinctBy(func(n int) any { return []int{n} }).ToSlice()
				if err != tc.expectedErr {
					t.Fatalf("expected error %v, got %v", tc.expectedErr, err)
				}
				return
			}

			res, err := FromSlice(tc.input).DistinctBy(func(p person) any { return p.id }).ToSlice()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(res) != len(tc.expected) {
				t.Fatalf("expected length %d, got %d", len(tc.expected), len(res))
			}
			for i := range res {
				if res[i].id != tc.expected[i] {
					t.Errorf("at index %d: expected ID %d, got %d", i, tc.expected[i], res[i].id)
				}
			}
		})
	}
}
