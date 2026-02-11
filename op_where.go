// Package iterlinq provides lazy LINQ-style sequence operators for Go.
package iterlinq

// Where filters the sequence based on a predicate.
func (s Sequence[T]) Where(predicate func(T) bool) Sequence[T] {
	if s.err != nil {
		return Sequence[T]{err: s.err}
	}
	return Sequence[T]{
		seq: func(yield func(T) bool) error {
			return s.seqOrEmpty()(func(v T) bool {
				if !predicate(v) {
					return true
				}
				return yield(v)
			})
		},
	}
}
