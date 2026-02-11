package rigid

import (
	"errors"
	"reflect"
	"strconv"
	"testing"
)

func TestComposeRunTypedPipeline(t *testing.T) {
	p := Compose(
		Where(func(v int) (bool, error) { return v%2 == 0, nil }),
		Select(func(v int) (string, error) { return strconv.Itoa(v * 10), nil }),
	)

	got, err := Run([]int{1, 2, 3, 4, 5}, p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []string{"20", "40"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected output: got=%v want=%v", got, want)
	}
}

func TestPipelineThenExecute(t *testing.T) {
	base := From([]int{1, 2, 3, 4, 5, 6})
	step1 := Then(base, Where(func(v int) (bool, error) { return v > 2, nil }))
	step2 := Then(step1, DistinctBy(func(v int) (int, error) { return v % 3, nil }))
	step3 := Then(step2, Select(func(v int) (int, error) { return v * v, nil }))

	got, err := step3.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []int{9, 16, 25}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected output: got=%v want=%v", got, want)
	}
}

func TestTakeSkipSelectMany(t *testing.T) {
	pipe := Compose(
		Compose(Skip[int](2), Take[int](3)),
		SelectMany(func(v int) ([]int, error) { return []int{v, -v}, nil }),
	)

	got, err := Run([]int{1, 2, 3, 4, 5, 6}, pipe)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []int{3, -3, 4, -4, 5, -5}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected output: got=%v want=%v", got, want)
	}
}

func TestIdentityCopies(t *testing.T) {
	in := []int{1, 2, 3}
	got, err := Identity[int]().Transform(in)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(got, in) {
		t.Fatalf("unexpected output: got=%v want=%v", got, in)
	}
	got[0] = 99
	if in[0] == 99 {
		t.Fatalf("identity should return a copy, but input was mutated")
	}
}

func TestErrorPropagation(t *testing.T) {
	boom := errors.New("boom")
	pipe := Compose(
		Where(func(v int) (bool, error) {
			if v == 3 {
				return false, boom
			}
			return true, nil
		}),
		Select(func(v int) (int, error) { return v * 2, nil }),
	)

	got, err := Run([]int{1, 2, 3, 4}, pipe)
	if !errors.Is(err, boom) {
		t.Fatalf("expected boom, got %v", err)
	}
	if got != nil {
		t.Fatalf("expected nil result on error, got %v", got)
	}
}

func TestNilGuards(t *testing.T) {
	var nilTransformer Transformer[int, int]
	if _, err := Run([]int{1}, nilTransformer); !errors.Is(err, ErrNilTransformer) {
		t.Fatalf("expected ErrNilTransformer, got %v", err)
	}

	if _, err := Where[int](nil).Transform([]int{1}); !errors.Is(err, ErrNilPredicate) {
		t.Fatalf("expected ErrNilPredicate, got %v", err)
	}

	if _, err := Select[int, int](nil).Transform([]int{1}); !errors.Is(err, ErrNilSelector) {
		t.Fatalf("expected ErrNilSelector, got %v", err)
	}

	if _, err := DistinctBy[int, int](nil).Transform([]int{1}); !errors.Is(err, ErrNilSelector) {
		t.Fatalf("expected ErrNilSelector, got %v", err)
	}
}
