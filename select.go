package iterlinq

// Select transforms each element of the sequence using the provided function and returns a sequence of any.
// Note: Methods cannot introduce new type parameters in Go, so this method returns Sequence[any].
// For strong typing, consider using package-level helpers in client code if acceptable.
func (s Sequence[T]) Select(transform func(T) any) Sequence[any] {
	return s.SelectWithError(func(v T) (any, error) { return transform(v), nil })
}

// SelectWithError transforms each element using a function that may return an error.
// If the transform returns an error, iteration stops and the error is yielded.
func (s Sequence[T]) SelectWithError(transform func(T) (any, error)) Sequence[any] {
	return Sequence[any]{
		seq: func(yield func(any, error) bool) {
			for v, err := range s.seq {
				if err != nil {
					yield(nil, err)
					return
				}
				transformed, err := transform(v)
				if err != nil {
					yield(nil, err)
					return
				}
				if !yield(transformed, nil) {
					return
				}
			}
		},
	}
}
