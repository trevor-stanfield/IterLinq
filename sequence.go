package iterlinq

import "iter"

// Sequence is a wrapper around iter.Seq2[T, error] to provide LINQ-like methods.
type Sequence[T any] struct {
	seq iter.Seq2[T, error]
}

// From creates a Sequence from an existing iter.Seq2.
// If seq is nil, it returns an empty sequence to avoid panics when ranging.
func From[T any](seq iter.Seq2[T, error]) Sequence[T] {
	if seq == nil {
		return Sequence[T]{
			seq: func(yield func(T, error) bool) {},
		}
	}
	return Sequence[T]{seq: seq}
}

// FromSlice creates a Sequence from a slice.
func FromSlice[T any](s []T) Sequence[T] {
	return Sequence[T]{
		seq: func(yield func(T, error) bool) {
			for _, v := range s {
				if !yield(v, nil) {
					return
				}
			}
		},
	}
}
