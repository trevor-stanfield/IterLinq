package iterlinq

import "testing"

// TestSequence tests Sequence creation and basic materialization.
func TestSequence(t *testing.T) {
	// Define test cases using handlers.
	handlers := []struct {
		name     string
		input    []int
		expected int
	}{
		{
			name:     "SliceWithThree",
			input:    []int{1, 2, 3},
			expected: 3,
		},
		{
			name:     "EmptySlice",
			input:    []int{},
			expected: 0,
		},
		{
			name:     "NilSlice",
			input:    nil,
			expected: 0,
		},
	}

	for _, h := range handlers {
		t.Run(h.name, func(t *testing.T) {
			s := FromSlice(h.input)
			res, err := s.ToSlice()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(res) != h.expected {
				t.Errorf("expected length %d, got %d", h.expected, len(res))
			}
		})
	}

	// Explicitly verify From(nil) yields an empty sequence (no panic, empty result)
	var zero Sequence[int] = From[int](nil)
	vals, err := zero.ToSlice()
	if err != nil {
		t.Fatalf("unexpected error from From(nil): %v", err)
	}
	if len(vals) != 0 {
		t.Fatalf("expected empty slice from From(nil), got %v", vals)
	}
}
