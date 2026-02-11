package iterlinq

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
	"testing/quick"
)

type Age int
type UserID string

type Event struct {
	ID     int
	Name   string
	Active bool
}

func quickCfg() *quick.Config {
	return &quick.Config{MaxCount: 300}
}

func slicesEqual[T any](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if !reflect.DeepEqual(a[i], b[i]) {
			return false
		}
	}
	return true
}

func TestPropertyFromSliceToSliceRoundTripInt(t *testing.T) {
	prop := func(xs []int) bool {
		got, err := FromSlice(xs).ToSlice()
		return err == nil && slicesEqual(got, xs)
	}
	if err := quick.Check(prop, quickCfg()); err != nil {
		t.Fatalf("property failed: %v", err)
	}
}

func TestPropertyFromSliceToSliceRoundTripString(t *testing.T) {
	prop := func(xs []string) bool {
		got, err := FromSlice(xs).ToSlice()
		return err == nil && slicesEqual(got, xs)
	}
	if err := quick.Check(prop, quickCfg()); err != nil {
		t.Fatalf("property failed: %v", err)
	}
}

func TestPropertyFromSliceToSliceRoundTripCustomStruct(t *testing.T) {
	prop := func(xs []Event) bool {
		got, err := FromSlice(xs).ToSlice()
		return err == nil && slicesEqual(got, xs)
	}
	if err := quick.Check(prop, quickCfg()); err != nil {
		t.Fatalf("property failed: %v", err)
	}
}

func TestPropertyWhereMatchesManualFilter(t *testing.T) {
	prop := func(xs []int) bool {
		want := make([]int, 0, len(xs))
		for _, v := range xs {
			if v%2 == 0 {
				want = append(want, v)
			}
		}
		got, err := FromSlice(xs).Where(func(v int) bool { return v%2 == 0 }).ToSlice()
		return err == nil && slicesEqual(got, want)
	}
	if err := quick.Check(prop, quickCfg()); err != nil {
		t.Fatalf("property failed: %v", err)
	}
}

func TestPropertyWhereCustomTypeMatchesManualFilter(t *testing.T) {
	prop := func(xs []Event) bool {
		want := make([]Event, 0, len(xs))
		for _, e := range xs {
			if e.Active {
				want = append(want, e)
			}
		}
		got, err := FromSlice(xs).Where(func(e Event) bool { return e.Active }).ToSlice()
		return err == nil && slicesEqual(got, want)
	}
	if err := quick.Check(prop, quickCfg()); err != nil {
		t.Fatalf("property failed: %v", err)
	}
}

func TestPropertyTakeMatchesManualPrefix(t *testing.T) {
	prop := func(xs []Age, n uint8) bool {
		k := int(n)
		if k > len(xs) {
			k = len(xs)
		}
		want := append([]Age(nil), xs[:k]...)
		got, err := FromSlice(xs).Take(int(n)).ToSlice()
		return err == nil && slicesEqual(got, want)
	}
	if err := quick.Check(prop, quickCfg()); err != nil {
		t.Fatalf("property failed: %v", err)
	}
}

func TestPropertySkipMatchesManualSuffix(t *testing.T) {
	prop := func(xs []UserID, n uint8) bool {
		k := int(n)
		if k > len(xs) {
			k = len(xs)
		}
		want := append([]UserID(nil), xs[k:]...)
		got, err := FromSlice(xs).Skip(int(n)).ToSlice()
		return err == nil && slicesEqual(got, want)
	}
	if err := quick.Check(prop, quickCfg()); err != nil {
		t.Fatalf("property failed: %v", err)
	}
}

func TestPropertySelectMatchesManualMap(t *testing.T) {
	prop := func(xs []Age) bool {
		want := make([]any, 0, len(xs))
		for _, v := range xs {
			want = append(want, int(v)*3)
		}
		got, err := FromSlice(xs).Select(func(v Age) any { return int(v) * 3 }).ToSlice()
		return err == nil && slicesEqual(got, want)
	}
	if err := quick.Check(prop, quickCfg()); err != nil {
		t.Fatalf("property failed: %v", err)
	}
}

func TestPropertySelectManyFlattensInOrder(t *testing.T) {
	prop := func(xs []int) bool {
		want := make([]any, 0, len(xs)*2)
		for _, v := range xs {
			want = append(want, v, -v)
		}
		got, err := FromSlice(xs).SelectMany(func(v int) Sequence[any] {
			return FromSlice([]any{v, -v})
		}).ToSlice()
		return err == nil && slicesEqual(got, want)
	}
	if err := quick.Check(prop, quickCfg()); err != nil {
		t.Fatalf("property failed: %v", err)
	}
}

func TestPropertyDistinctByFirstWinsAndUnique(t *testing.T) {
	prop := func(xs []Event) bool {
		got, err := FromSlice(xs).DistinctBy(func(e Event) any { return e.ID }).ToSlice()
		if err != nil {
			return false
		}

		firstByID := map[int]Event{}
		order := make([]int, 0, len(xs))
		for _, e := range xs {
			if _, ok := firstByID[e.ID]; !ok {
				firstByID[e.ID] = e
				order = append(order, e.ID)
			}
		}
		want := make([]Event, 0, len(order))
		for _, id := range order {
			want = append(want, firstByID[id])
		}
		return slicesEqual(got, want)
	}
	if err := quick.Check(prop, quickCfg()); err != nil {
		t.Fatalf("property failed: %v", err)
	}
}

func TestPropertyCountMatchesLenWithoutError(t *testing.T) {
	prop := func(xs []Event) bool {
		got, err := FromSlice(xs).Count()
		return err == nil && got == len(xs)
	}
	if err := quick.Check(prop, quickCfg()); err != nil {
		t.Fatalf("property failed: %v", err)
	}
}

func TestPropertyAnyAllDualityWithoutError(t *testing.T) {
	prop := func(xs []int) bool {
		anyEven, err := FromSlice(xs).Any(func(v int) bool { return v%2 == 0 })
		if err != nil {
			return false
		}
		allOdd, err := FromSlice(xs).All(func(v int) bool { return v%2 != 0 })
		if err != nil {
			return false
		}
		return anyEven == !allOdd
	}
	if err := quick.Check(prop, quickCfg()); err != nil {
		t.Fatalf("property failed: %v", err)
	}
}

func TestPropertyFirstEqualsManualFirstMatch(t *testing.T) {
	prop := func(xs []int) bool {
		got, err := FromSlice(xs).First(func(v int) bool { return v%3 == 0 })

		wantFound := false
		want := 0
		for _, v := range xs {
			if v%3 == 0 {
				want = v
				wantFound = true
				break
			}
		}

		if wantFound {
			return err == nil && got == want
		}
		return err == ErrNoMatch
	}
	if err := quick.Check(prop, quickCfg()); err != nil {
		t.Fatalf("property failed: %v", err)
	}
}

func TestPropertyLastEqualsManualLastMatch(t *testing.T) {
	prop := func(xs []int) bool {
		got, err := FromSlice(xs).Last(func(v int) bool { return v%3 == 0 })

		wantFound := false
		want := 0
		for i := len(xs) - 1; i >= 0; i-- {
			if xs[i]%3 == 0 {
				want = xs[i]
				wantFound = true
				break
			}
		}

		if wantFound {
			return err == nil && got == want
		}
		return err == ErrNoMatch
	}
	if err := quick.Check(prop, quickCfg()); err != nil {
		t.Fatalf("property failed: %v", err)
	}
}

func TestPropertyEarliestSequenceErrorWins(t *testing.T) {
	prop := func(xs []int, fail uint8) bool {
		if len(xs) == 0 {
			return true
		}
		i := int(fail) % len(xs)
		wantErr := fmt.Sprintf("boom:%d", i)

		got, err := FromFunc(func(yield func(int) bool) error {
			for idx, v := range xs {
				if idx == i {
					return errors.New(wantErr)
				}
				if !yield(v) {
					return nil
				}
			}
			return nil
		}).Where(func(v int) bool { return v%2 == 0 }).ToSlice()

		return got == nil && err != nil && err.Error() == wantErr
	}
	if err := quick.Check(prop, quickCfg()); err != nil {
		t.Fatalf("property failed: %v", err)
	}
}
