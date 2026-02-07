package iterlinq

// SelectMany flattens a sequence of sequences into a single sequence of any.
func (s Sequence[T]) SelectMany(transform func(T) Sequence[any]) Sequence[any] {
	return Sequence[any]{
		seq: func(yield func(any, error) bool) {
			for v, err := range s.seq {
				if err != nil {
					yield(nil, err)
					return
				}
				inner := transform(v)
				for iv, ierr := range inner.seq {
					if ierr != nil {
						yield(nil, ierr)
						return
					}
					if !yield(iv, nil) {
						return
					}
				}
			}
		},
	}
}
