package rigid

import (
	"errors"
	"testing"
)

func TestPipelineErrorStateInspection(t *testing.T) {
	var nilTransformer Transformer[int, int]
	p := Then(From([]int{1}), nilTransformer)
	if !p.HasError() {
		t.Fatalf("expected pipeline error state")
	}
	if !errors.Is(p.Err(), ErrNilTransformer) {
		t.Fatalf("expected ErrNilTransformer, got %v", p.Err())
	}
}

func TestPipelineMapError(t *testing.T) {
	boom := errors.New("boom")
	p := Then(From([]int{1, 2, 3}), Where(func(v int) (bool, error) {
		if v == 2 {
			return false, boom
		}
		return true, nil
	})).MapError(func(err error) error {
		return errors.New("mapped: " + err.Error())
	})

	_, err := p.Execute()
	if err == nil || err.Error() != "mapped: boom" {
		t.Fatalf("expected mapped error, got %v", err)
	}
}

func TestPipelineRecoverAndOrElse(t *testing.T) {
	boom := errors.New("boom")
	base := Then(From([]int{1, 2, 3}), Where(func(v int) (bool, error) {
		if v == 2 {
			return false, boom
		}
		return true, nil
	}))

	recovered := base.Recover(func(err error) Pipeline[int, int] {
		if !errors.Is(err, boom) {
			t.Fatalf("unexpected error passed to recover: %v", err)
		}
		return From([]int{9, 10})
	})
	got, err := recovered.Execute()
	if err != nil {
		t.Fatalf("unexpected recover err: %v", err)
	}
	if len(got) != 2 || got[0] != 9 || got[1] != 10 {
		t.Fatalf("unexpected recovered values: %v", got)
	}

	fallback := From([]int{7, 8})
	got, err = base.OrElse(fallback).Execute()
	if err != nil {
		t.Fatalf("unexpected orElse err: %v", err)
	}
	if len(got) != 2 || got[0] != 7 || got[1] != 8 {
		t.Fatalf("unexpected orElse values: %v", got)
	}
}

func TestPipelineNilErrorHandlers(t *testing.T) {
	if _, err := From([]int{1}).MapError(nil).Execute(); !errors.Is(err, ErrNilErrorMapper) {
		t.Fatalf("expected ErrNilErrorMapper, got %v", err)
	}
	if _, err := From([]int{1}).Recover(nil).Execute(); !errors.Is(err, ErrNilErrorHandler) {
		t.Fatalf("expected ErrNilErrorHandler, got %v", err)
	}
}
