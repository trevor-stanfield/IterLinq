package iterlinq

// Select transforms each element of the sequence using the provided function and returns a sequence of any.
// Note: Methods cannot introduce new type parameters in Go, so this method returns Sequence[any].
// For strong typing, consider using package-level helpers in client code if acceptable.
func (s Sequence[T]) Select(transform func(T) any) Sequence[any] {
	if s.err != nil {
		return Sequence[any]{err: s.err}
	}
	return Sequence[any]{
		seq: func(yield func(any) bool) error {
			return s.seqOrEmpty()(func(v T) bool {
				return yield(transform(v))
			})
		},
	}
}
