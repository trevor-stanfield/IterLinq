package iterlinq

// Where filters the sequence based on a predicate.
func (s Sequence[T]) Where(predicate func(T) bool) Sequence[T] {
	return s.WhereWithError(func(v T) (bool, error) {
		return predicate(v), nil
	})
}

// WhereWithError filters the sequence based on a predicate that can return an error.
func (s Sequence[T]) WhereWithError(predicate func(T) (bool, error)) Sequence[T] {
	return Sequence[T]{
		seq: func(yield func(T, error) bool) {
			for v, err := range s.seq {
				if err != nil {
					yield(v, err)
					return
				}
				match, err := predicate(v)
				if err != nil {
					yield(v, err)
					return
				}
				if match {
					if !yield(v, nil) {
						return
					}
				}
			}
		},
	}
}
