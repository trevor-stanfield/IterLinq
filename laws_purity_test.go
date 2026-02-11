package iterlinq

import (
	"reflect"
	"testing"
)

// TestPipelineDoesNotMutateInputSlice verifies operators do not mutate source slices.
func TestPipelineDoesNotMutateInputSlice(t *testing.T) {
	input := []int{5, 1, 3, 1, 2, 4}
	original := append([]int(nil), input...)

	_, err := FromSlice(input).
		Where(func(n int) bool { return n > 1 }).
		Skip(1).
		Take(3).
		DistinctBy(func(n int) any { return n }).
		ToSlice()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !reflect.DeepEqual(input, original) {
		t.Fatalf("input slice was mutated: got %v want %v", input, original)
	}
}

// TestRepeatedEnumerationIsDeterministic verifies repeated materialization is stable.
func TestRepeatedEnumerationIsDeterministic(t *testing.T) {
	seq := FromSlice([]int{1, 2, 3, 4, 5, 6}).
		Where(func(n int) bool { return n%2 == 0 }).
		Select(func(n int) any { return n * n })

	first, err := seq.ToSlice()
	if err != nil {
		t.Fatalf("unexpected error on first materialization: %v", err)
	}
	second, err := seq.ToSlice()
	if err != nil {
		t.Fatalf("unexpected error on second materialization: %v", err)
	}

	if !reflect.DeepEqual(first, second) {
		t.Fatalf("materializations differ: first=%v second=%v", first, second)
	}
}

// TestDerivedSequencesAreIndependent verifies branching from a base sequence is immutable.
func TestDerivedSequencesAreIndependent(t *testing.T) {
	base := FromSlice([]int{1, 2, 3, 4, 5, 6})
	evens := base.Where(func(n int) bool { return n%2 == 0 })
	odds := base.Where(func(n int) bool { return n%2 != 0 })

	evenVals, err := evens.ToSlice()
	if err != nil {
		t.Fatalf("unexpected error from evens: %v", err)
	}
	oddVals, err := odds.ToSlice()
	if err != nil {
		t.Fatalf("unexpected error from odds: %v", err)
	}
	baseVals, err := base.ToSlice()
	if err != nil {
		t.Fatalf("unexpected error from base: %v", err)
	}

	if !reflect.DeepEqual(evenVals, []int{2, 4, 6}) {
		t.Fatalf("unexpected evens: %v", evenVals)
	}
	if !reflect.DeepEqual(oddVals, []int{1, 3, 5}) {
		t.Fatalf("unexpected odds: %v", oddVals)
	}
	if !reflect.DeepEqual(baseVals, []int{1, 2, 3, 4, 5, 6}) {
		t.Fatalf("base sequence was affected: %v", baseVals)
	}
}
