package iterlinq

import "iter"

type seqFunc[T any] func(yield func(T) bool) error

// Sequence is a wrapper around iter.Seq[T] with monadic error state.
type Sequence[T any] struct {
	seq seqFunc[T]
	err error
}

func emptySeq[T any]() seqFunc[T] {
	return func(yield func(T) bool) error { return nil }
}

func failed[T any](err error) Sequence[T] {
	if err == nil {
		return Sequence[T]{seq: emptySeq[T]()}
	}
	return Sequence[T]{err: err}
}

func (s Sequence[T]) seqOrEmpty() seqFunc[T] {
	if s.seq == nil {
		return emptySeq[T]()
	}
	return s.seq
}

// From creates a Sequence from an existing iter.Seq.
// If seq is nil, it returns an empty sequence to avoid panics when ranging.
func From[T any](seq iter.Seq[T]) Sequence[T] {
	if seq == nil {
		return Sequence[T]{
			seq: emptySeq[T](),
		}
	}
	return Sequence[T]{
		seq: func(yield func(T) bool) error {
			for v := range seq {
				if !yield(v) {
					return nil
				}
			}
			return nil
		},
	}
}

// FromFunc creates a Sequence from a function that yields values and returns an optional error.
func FromFunc[T any](fn func(yield func(T) bool) error) Sequence[T] {
	if fn == nil {
		return Sequence[T]{seq: emptySeq[T]()}
	}
	return Sequence[T]{seq: fn}
}

// FromError creates a failed sequence that short-circuits subsequent operators.
func FromError[T any](err error) Sequence[T] {
	return failed[T](err)
}

// FromSlice creates a Sequence from a slice.
func FromSlice[T any](s []T) Sequence[T] {
	return Sequence[T]{
		seq: func(yield func(T) bool) error {
			for _, v := range s {
				if !yield(v) {
					return nil
				}
			}
			return nil
		},
	}
}
