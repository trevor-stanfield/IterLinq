package iterlinq

// SelectMany flattens a sequence of sequences into a single sequence of any.
func (s Sequence[T]) SelectMany(transform func(T) Sequence[any]) Sequence[any] {
	if s.err != nil {
		return Sequence[any]{err: s.err}
	}
	return Sequence[any]{
		seq: func(yield func(any) bool) error {
			var innerErr error
			err := s.seqOrEmpty()(func(v T) bool {
				inner := transform(v)
				if inner.err != nil {
					innerErr = inner.err
					return false
				}
				stopped := false
				err := inner.seqOrEmpty()(func(iv any) bool {
					if !yield(iv) {
						stopped = true
						return false
					}
					return true
				})
				if err != nil {
					innerErr = err
					return false
				}
				if stopped {
					return false
				}
				return true
			})
			if innerErr != nil {
				return innerErr
			}
			return err
		},
	}
}
