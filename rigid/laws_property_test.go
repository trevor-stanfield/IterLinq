package rigid

import (
	"reflect"
	"testing"
	"testing/quick"
)

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

func TestPropertyIdentityRoundTrip(t *testing.T) {
	prop := func(xs []int) bool {
		got, err := Run(xs, Identity[int]())
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
		got, err := Run(xs, Where(func(v int) (bool, error) { return v%2 == 0, nil }))
		return err == nil && slicesEqual(got, want)
	}
	if err := quick.Check(prop, quickCfg()); err != nil {
		t.Fatalf("property failed: %v", err)
	}
}

func TestPropertySelectMatchesManualMap(t *testing.T) {
	prop := func(xs []int) bool {
		want := make([]string, 0, len(xs))
		for _, v := range xs {
			want = append(want, string(rune(v)))
		}
		got, err := Run(xs, Select(func(v int) (string, error) { return string(rune(v)), nil }))
		return err == nil && slicesEqual(got, want)
	}
	if err := quick.Check(prop, quickCfg()); err != nil {
		t.Fatalf("property failed: %v", err)
	}
}

func TestPropertyTakeSkipReconstruct(t *testing.T) {
	prop := func(xs []int, n uint8) bool {
		k := int(n)
		head, err := Run(xs, Take[int](k))
		if err != nil {
			return false
		}
		tail, err := Run(xs, Skip[int](k))
		if err != nil {
			return false
		}
		reconstructed := append(append([]int{}, head...), tail...)
		return slicesEqual(reconstructed, xs)
	}
	if err := quick.Check(prop, quickCfg()); err != nil {
		t.Fatalf("property failed: %v", err)
	}
}
