package rigid

import (
	"errors"
	"slices"
	"sort"
	"strconv"
	"testing"
)

type orderedRecord struct {
	id    int
	name  string
	score int
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

func compareNames(left, right orderedRecord) int {
	switch {
	case left.name < right.name:
		return -1
	case left.name > right.name:
		return 1
	default:
		return 0
	}
}

func compareScores(left, right orderedRecord) int {
	return compareIntegers(left.score, right.score)
}

func TestOrderedPipeline_DefaultSortHonorsThenByPriority(t *testing.T) {
	records := []orderedRecord{
		{id: 1, name: "bob", score: 3},
		{id: 2, name: "alice", score: 4},
		{id: 3, name: "bob", score: 1},
		{id: 4, name: "alice", score: 2},
	}

	orderedRecords, err := From(records).
		MustOrderBy(compareNames).
		ThenBy(compareScores).
		Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedIDs := []int{4, 2, 3, 1}
	actualIDs := make([]int, 0, len(orderedRecords))
	for _, record := range orderedRecords {
		actualIDs = append(actualIDs, record.id)
	}
	if !slices.Equal(actualIDs, expectedIDs) {
		t.Fatalf("unexpected order: got %v, want %v", actualIDs, expectedIDs)
	}
}

func TestOrderedPipeline_ThenByDesc(t *testing.T) {
	records := []orderedRecord{
		{id: 1, name: "bob", score: 3},
		{id: 2, name: "alice", score: 4},
		{id: 3, name: "bob", score: 1},
		{id: 4, name: "alice", score: 2},
	}

	orderedRecords, err := From(records).
		MustOrderBy(compareNames).
		ThenByDesc(compareScores).
		Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedIDs := []int{2, 4, 1, 3}
	actualIDs := make([]int, 0, len(orderedRecords))
	for _, record := range orderedRecords {
		actualIDs = append(actualIDs, record.id)
	}
	if !slices.Equal(actualIDs, expectedIDs) {
		t.Fatalf("unexpected order: got %v, want %v", actualIDs, expectedIDs)
	}
}

func TestOrderedPipeline_ReverseForcesBeforeNextOp(t *testing.T) {
	orderedValues, err := From([]int{3, 1, 2}).
		MustOrderBy(compareIntegers).
		Reverse().
		Take(2).
		Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !slices.Equal(orderedValues, []int{3, 2}) {
		t.Fatalf("unexpected result: got %v, want [3 2]", orderedValues)
	}
}

func TestOrderedPipeline_UsingSortFuncReceivesComparator(t *testing.T) {
	sorterWasCalled := false
	sorter := func(items []int, comparator Comparator[int]) error {
		sorterWasCalled = true
		sort.SliceStable(items, func(leftIndex, rightIndex int) bool {
			return comparator(items[leftIndex], items[rightIndex]) < 0
		})
		return nil
	}

	orderedValues, err := From([]int{3, 1, 2}).
		MustOrderBy(compareIntegers).
		UsingSorter(sorter).
		Execute()
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

func TestOrderedPipeline_UsingSortFuncErrorPropagates(t *testing.T) {
	boom := errors.New("boom")
	_, err := From([]int{3, 1, 2}).
		MustOrderBy(compareIntegers).
		UsingSorter(func(items []int, comparator Comparator[int]) error {
			return boom
		}).
		Execute()
	if !errors.Is(err, boom) {
		t.Fatalf("expected boom, got %v", err)
	}
}

func TestOrderedPipeline_NilComparatorOrSorter(t *testing.T) {
	testCases := []struct {
		name            string
		orderedPipeline OrderedPipeline[int, int]
		expectedError   error
	}{
		{
			name:            "must order by nil comparator",
			orderedPipeline: From([]int{1}).MustOrderBy(nil),
			expectedError:   ErrNilComparator,
		},
		{
			name:            "then by nil comparator",
			orderedPipeline: From([]int{1}).MustOrderBy(compareIntegers).ThenBy(nil),
			expectedError:   ErrNilComparator,
		},
		{
			name:            "then by desc nil comparator",
			orderedPipeline: From([]int{1}).MustOrderBy(compareIntegers).ThenByDesc(nil),
			expectedError:   ErrNilComparator,
		},
		{
			name:            "using nil sorter",
			orderedPipeline: From([]int{1}).MustOrderBy(compareIntegers).UsingSorter(nil),
			expectedError:   ErrNilSortFunc,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			_, err := testCase.orderedPipeline.Execute()
			if !errors.Is(err, testCase.expectedError) {
				t.Fatalf("expected %v, got %v", testCase.expectedError, err)
			}
		})
	}
}

func TestOrderedPipeline_EarliestErrorWins(t *testing.T) {
	var nilTransformer Transformer[int, int]
	pipeline := Then(From([]int{1, 2, 3}), nilTransformer)
	_, err := pipeline.MustOrderBy(nil).Execute()
	if !errors.Is(err, ErrNilTransformer) {
		t.Fatalf("expected original error, got %v", err)
	}
}

func TestSelectOrdered_ForcesOrderBeforeProjection(t *testing.T) {
	orderedPipeline := From([]int{3, 1, 2}).MustOrderBy(compareIntegers)
	pipeline := SelectOrdered(orderedPipeline, func(value int) (string, error) {
		return strconv.Itoa(value), nil
	})

	values, err := pipeline.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !slices.Equal(values, []string{"1", "2", "3"}) {
		t.Fatalf("unexpected result: got %v, want [1 2 3]", values)
	}
}

func TestDistinctByOrdered_ForcesOrderBeforeDistinct(t *testing.T) {
	orderedPipeline := From([]int{3, 1, 2, 4}).MustOrderBy(compareIntegers).Reverse()
	pipeline := DistinctByOrdered(orderedPipeline, func(value int) (int, error) {
		return value % 2, nil
	})

	values, err := pipeline.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !slices.Equal(values, []int{4, 3}) {
		t.Fatalf("unexpected result: got %v, want [4 3]", values)
	}
}
