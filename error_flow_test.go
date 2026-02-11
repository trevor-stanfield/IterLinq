package iterlinq

import (
	"errors"
	"fmt"
	"testing"
)

func TestFromErrorShortCircuitsOperators(t *testing.T) {
	boom := errors.New("boom")
	called := false

	s := FromError[int](boom).
		Where(func(v int) bool {
			called = true
			return v > 0
		}).
		Take(3)

	if s.err != boom {
		t.Fatalf("expected monadic error %v, got %v", boom, s.err)
	}
	if s.seq != nil {
		t.Fatalf("expected nil seq on errored sequence")
	}

	got, err := s.ToSlice()
	if err != boom {
		t.Fatalf("expected error %v, got %v", boom, err)
	}
	if got != nil {
		t.Fatalf("expected nil materialized slice, got %v", got)
	}
	if called {
		t.Fatalf("expected predicate not to be called on short-circuited sequence")
	}
}

func TestTerminalsReturnMonadicError(t *testing.T) {
	boom := errors.New("boom")
	s := FromError[int](boom)
	if !s.HasError() || s.Err() != boom {
		t.Fatalf("expected monadic error inspection methods to reflect boom")
	}

	if _, err := s.Count(); err != boom {
		t.Fatalf("count expected %v, got %v", boom, err)
	}
	if _, err := s.Any(func(v int) bool { return v%2 == 0 }); err != boom {
		t.Fatalf("any expected %v, got %v", boom, err)
	}
	if _, err := s.All(func(v int) bool { return v%2 == 0 }); err != boom {
		t.Fatalf("all expected %v, got %v", boom, err)
	}
	if _, err := s.First(func(v int) bool { return true }); err != boom {
		t.Fatalf("first expected %v, got %v", boom, err)
	}
	if _, err := s.Last(func(v int) bool { return true }); err != boom {
		t.Fatalf("last expected %v, got %v", boom, err)
	}
}

func TestMapErrorOnMonadicError(t *testing.T) {
	s := FromError[int](errors.New("raw")).
		MapError(func(err error) error {
			return fmt.Errorf("wrapped: %w", err)
		})

	_, err := s.ToSlice()
	if err == nil || err.Error() != "wrapped: raw" {
		t.Fatalf("expected wrapped error, got %v", err)
	}
}

func TestMapErrorOnStreamError(t *testing.T) {
	s := FromFunc(func(yield func(int) bool) error {
		if !yield(1) {
			return nil
		}
		return errors.New("raw")
	}).MapError(func(err error) error {
		return fmt.Errorf("mapped: %w", err)
	})

	_, err := s.ToSlice()
	if err == nil || err.Error() != "mapped: raw" {
		t.Fatalf("expected mapped error, got %v", err)
	}
}

func TestRecoverOnMonadicError(t *testing.T) {
	got, err := FromError[int](errors.New("boom")).
		Recover(func(error) Sequence[int] {
			return FromSlice([]int{7, 8, 9})
		}).
		ToSlice()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 3 || got[0] != 7 || got[1] != 8 || got[2] != 9 {
		t.Fatalf("unexpected recovered values: %v", got)
	}
}

func TestRecoverOnStreamError(t *testing.T) {
	got, err := FromFunc(func(yield func(int) bool) error {
		if !yield(1) {
			return nil
		}
		if !yield(2) {
			return nil
		}
		return errors.New("boom")
	}).Recover(func(error) Sequence[int] {
		return FromSlice([]int{9, 10})
	}).ToSlice()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 4 || got[0] != 1 || got[1] != 2 || got[2] != 9 || got[3] != 10 {
		t.Fatalf("unexpected recovered values: %v", got)
	}
}

func TestOrElseOnError(t *testing.T) {
	got, err := FromError[int](errors.New("boom")).
		OrElse(FromSlice([]int{3, 4})).
		ToSlice()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 || got[0] != 3 || got[1] != 4 {
		t.Fatalf("unexpected values from OrElse: %v", got)
	}
}

func TestNilMapErrorAndRecoverHandlers(t *testing.T) {
	if _, err := FromSlice([]int{1}).MapError(nil).ToSlice(); err != ErrNilErrorMapper {
		t.Fatalf("expected ErrNilErrorMapper, got %v", err)
	}
	if _, err := FromSlice([]int{1}).Recover(nil).ToSlice(); err != ErrNilErrorHandler {
		t.Fatalf("expected ErrNilErrorHandler, got %v", err)
	}
}
