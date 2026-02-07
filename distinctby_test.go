package iterlinq

import "testing"

// TestDistinctBy tests the DistinctBy operation.
func TestDistinctBy(t *testing.T) {
	type person struct {
		id   int
		name string
	}

	// Define test cases using handlers.
	handlers := []struct {
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

	for _, h := range handlers {
		t.Run(h.name, func(t *testing.T) {
			if h.useInvalid {
				nums := []int{1, 2}
				_, err := FromSlice(nums).DistinctBy(func(n int) any { return []int{n} }).ToSlice()
				if err != h.expectedErr {
					t.Fatalf("expected error %v, got %v", h.expectedErr, err)
				}
				return
			}

			res, err := FromSlice(h.input).DistinctBy(func(p person) any { return p.id }).ToSlice()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(res) != len(h.expected) {
				t.Fatalf("expected length %d, got %d", len(h.expected), len(res))
			}
			for i := range res {
				if res[i].id != h.expected[i] {
					t.Errorf("at index %d: expected ID %d, got %d", i, h.expected[i], res[i].id)
				}
			}
		})
	}
}
