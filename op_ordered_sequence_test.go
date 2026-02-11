package iterlinq

import (
	"errors"
	"slices"
	"sort"
	"testing"
)

type orderedPerson struct {
	id        int
	name      string
	birthYear int
}

func compareIntegers(left, right int) int {
	switch {
	case left < right:
		return -1
	case left > right:
		return 1
	default:
		return 0
	}
}

func compareNames(left, right orderedPerson) int {
	switch {
	case left.name < right.name:
		return -1
	case left.name > right.name:
		return 1
	default:
		return 0
	}
}

func compareBirthYears(left, right orderedPerson) int {
	return compareIntegers(left.birthYear, right.birthYear)
}

func TestOrderedSequence_DefaultSortHonorsThenByPriority(t *testing.T) {
	people := []orderedPerson{
		{id: 1, name: "bob", birthYear: 1990},
		{id: 2, name: "alice", birthYear: 1991},
		{id: 3, name: "bob", birthYear: 1985},
		{id: 4, name: "alice", birthYear: 1980},
	}

	orderedPeople, err := FromSlice(people).
		MustOrderBy(compareNames).
		ThenBy(compareBirthYears).
		ToSlice()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedIDs := []int{4, 2, 3, 1}
	actualIDs := make([]int, 0, len(orderedPeople))
	for _, person := range orderedPeople {
		actualIDs = append(actualIDs, person.id)
	}
	if !slices.Equal(actualIDs, expectedIDs) {
		t.Fatalf("unexpected order: got %v, want %v", actualIDs, expectedIDs)
	}
}

func TestOrderedSequence_ThenByDesc(t *testing.T) {
	people := []orderedPerson{
		{id: 1, name: "bob", birthYear: 1990},
		{id: 2, name: "alice", birthYear: 1991},
		{id: 3, name: "bob", birthYear: 1985},
		{id: 4, name: "alice", birthYear: 1980},
	}

	orderedPeople, err := FromSlice(people).
		MustOrderBy(compareNames).
		ThenByDesc(compareBirthYears).
		ToSlice()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedIDs := []int{2, 4, 1, 3}
	actualIDs := make([]int, 0, len(orderedPeople))
	for _, person := range orderedPeople {
		actualIDs = append(actualIDs, person.id)
	}
	if !slices.Equal(actualIDs, expectedIDs) {
		t.Fatalf("unexpected order: got %v, want %v", actualIDs, expectedIDs)
	}
}

func TestOrderedSequence_ReverseForcesBeforeNextOp(t *testing.T) {
	orderedValues, err := FromSlice([]int{3, 1, 2}).
		MustOrderBy(compareIntegers).
		Reverse().
		Take(2).
		ToSlice()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !slices.Equal(orderedValues, []int{3, 2}) {
		t.Fatalf("unexpected result: got %v, want [3 2]", orderedValues)
	}
}

func TestOrderedSequence_UsingSortFuncReceivesComparator(t *testing.T) {
	sorterWasCalled := false
	sorter := func(items []int, comparator Comparator[int]) error {
		sorterWasCalled = true
		sort.SliceStable(items, func(leftIndex, rightIndex int) bool {
			return comparator(items[leftIndex], items[rightIndex]) < 0
		})
		return nil
	}

	orderedValues, err := FromSlice([]int{3, 1, 2}).
		MustOrderBy(compareIntegers).
		UsingSorter(sorter).
		ToSlice()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !sorterWasCalled {
		t.Fatalf("expected custom sorter to be called")
	}
	if !slices.Equal(orderedValues, []int{1, 2, 3}) {
		t.Fatalf("unexpected result: got %v, want [1 2 3]", orderedValues)
	}
}

func TestOrderedSequence_UsingSortFuncErrorPropagates(t *testing.T) {
	boom := errors.New("boom")
	_, err := FromSlice([]int{3, 1, 2}).
		MustOrderBy(compareIntegers).
		UsingSorter(func(items []int, comparator Comparator[int]) error {
			return boom
		}).
		ToSlice()
	if !errors.Is(err, boom) {
		t.Fatalf("expected boom, got %v", err)
	}
}

func TestOrderedSequence_NilComparatorOrSorter(t *testing.T) {
	testCases := []struct {
		name            string
		orderedSequence OrderedSequence[int]
		expectedError   error
	}{
		{
			name:            "must order by nil comparator",
			orderedSequence: FromSlice([]int{1}).MustOrderBy(nil),
			expectedError:   ErrNilComparator,
		},
		{
			name:            "then by nil comparator",
			orderedSequence: FromSlice([]int{1}).MustOrderBy(compareIntegers).ThenBy(nil),
			expectedError:   ErrNilComparator,
		},
		{
			name:            "then by desc nil comparator",
			orderedSequence: FromSlice([]int{1}).MustOrderBy(compareIntegers).ThenByDesc(nil),
			expectedError:   ErrNilComparator,
		},
		{
			name:            "using nil sorter",
			orderedSequence: FromSlice([]int{1}).MustOrderBy(compareIntegers).UsingSorter(nil),
			expectedError:   ErrNilSortFunc,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			_, err := testCase.orderedSequence.ToSlice()
			if !errors.Is(err, testCase.expectedError) {
				t.Fatalf("expected %v, got %v", testCase.expectedError, err)
			}
		})
	}
}

func TestOrderedSequence_EarliestErrorWins(t *testing.T) {
	boom := errors.New("boom")
	_, err := FromError[int](boom).MustOrderBy(nil).ToSlice()
	if !errors.Is(err, boom) {
		t.Fatalf("expected original error, got %v", err)
	}
}
