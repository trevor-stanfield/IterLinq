package rigid

import (
	"errors"
	"testing"
)

func TestTerminalCountAnyAll(t *testing.T) {
	p := Then(From([]int{1, 2, 3, 4, 5}), Where(func(v int) (bool, error) { return v > 1, nil }))

	count, err := p.Count()
	if err != nil || count != 4 {
		t.Fatalf("count mismatch: count=%d err=%v", count, err)
	}

	anyEven, err := p.Any(func(v int) (bool, error) { return v%2 == 0, nil })
	if err != nil || !anyEven {
		t.Fatalf("any mismatch: any=%v err=%v", anyEven, err)
	}

	allPositive, err := p.All(func(v int) (bool, error) { return v > 0, nil })
	if err != nil || !allPositive {
		t.Fatalf("all mismatch: all=%v err=%v", allPositive, err)
	}
}

func TestTerminalFirstLastFamily(t *testing.T) {
	p := From([]int{1, 2, 3, 4, 5, 6})

	first, err := p.First(func(v int) (bool, error) { return v > 3, nil })
	if err != nil || first != 4 {
		t.Fatalf("first mismatch: first=%d err=%v", first, err)
	}

	last, err := p.Last(func(v int) (bool, error) { return v%2 == 0, nil })
	if err != nil || last != 6 {
		t.Fatalf("last mismatch: last=%d err=%v", last, err)
	}

	if _, err := p.First(func(v int) (bool, error) { return v > 100, nil }); !errors.Is(err, ErrNoMatch) {
		t.Fatalf("expected ErrNoMatch from First, got %v", err)
	}
	if _, err := p.Last(func(v int) (bool, error) { return v > 100, nil }); !errors.Is(err, ErrNoMatch) {
		t.Fatalf("expected ErrNoMatch from Last, got %v", err)
	}

	firstDef, err := p.FirstOrDefault(func(v int) (bool, error) { return v > 100, nil })
	if err != nil || firstDef != 0 {
		t.Fatalf("first default mismatch: val=%d err=%v", firstDef, err)
	}
	lastDef, err := p.LastOrDefault(func(v int) (bool, error) { return v > 100, nil })
	if err != nil || lastDef != 0 {
		t.Fatalf("last default mismatch: val=%d err=%v", lastDef, err)
	}
}

func TestTerminalPredicateErrorPropagation(t *testing.T) {
	boom := errors.New("boom")
	p := From([]int{1, 2, 3})

	if _, err := p.Any(func(v int) (bool, error) {
		if v == 2 {
			return false, boom
		}
		return false, nil
	}); !errors.Is(err, boom) {
		t.Fatalf("expected boom from Any, got %v", err)
	}

	if _, err := p.First(func(v int) (bool, error) {
		if v == 2 {
			return false, boom
		}
		return false, nil
	}); !errors.Is(err, boom) {
		t.Fatalf("expected boom from First, got %v", err)
	}
}
