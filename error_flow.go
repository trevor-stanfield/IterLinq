package iterlinq

// Err returns the monadic error currently attached to the sequence, if any.
func (s Sequence[T]) Err() error {
	return s.err
}

// HasError reports whether the sequence is in an error state.
func (s Sequence[T]) HasError() bool {
	return s.err != nil
}

// MapError maps an error carried by this sequence (or produced during enumeration).
// If mapper returns nil, the error is cleared and iteration terminates.
func (s Sequence[T]) MapError(mapper func(error) error) Sequence[T] {
	if mapper == nil {
		return failed[T](ErrNilErrorMapper)
	}
	if s.err != nil {
		return failed[T](mapper(s.err))
	}
	return Sequence[T]{
		seq: func(yield func(T) bool) error {
			err := s.seqOrEmpty()(yield)
			if err == nil {
				return nil
			}
			return mapper(err)
		},
	}
}

// Recover handles an error carried by this sequence (or produced during enumeration)
// by switching to a replacement sequence.
func (s Sequence[T]) Recover(handler func(error) Sequence[T]) Sequence[T] {
	if handler == nil {
		return failed[T](ErrNilErrorHandler)
	}
	if s.err != nil {
		return handler(s.err)
	}
	return Sequence[T]{
		seq: func(yield func(T) bool) error {
			err := s.seqOrEmpty()(yield)
			if err == nil {
				return nil
			}
			recovered := handler(err)
			if recovered.err != nil {
				return recovered.err
			}
			return recovered.seqOrEmpty()(yield)
		},
	}
}

// OrElse replaces an errored sequence with the provided fallback sequence.
func (s Sequence[T]) OrElse(fallback Sequence[T]) Sequence[T] {
	return s.Recover(func(error) Sequence[T] { return fallback })
}
