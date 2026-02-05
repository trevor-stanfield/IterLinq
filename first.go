package iterlinq

// First returns the first element of the sequence that satisfies the predicate.
// It short-circuits on the first match. If an element error is encountered, it is returned immediately.
// If no element matches, it returns ErrNoMatch.
func (s Sequence[T]) First(predicate func(T) bool) (T, error) {
	return s.FirstWithError(func(v T) (bool, error) {
		return predicate(v), nil
	})
}

// FirstWithError returns the first element of the sequence that satisfies the error-aware predicate.
// The iteration stops on the first match or the first error from the sequence or the predicate.
// If no element matches, it returns ErrNoMatch.
func (s Sequence[T]) FirstWithError(predicate func(T) (bool, error)) (T, error) {
	for v, err := range s.seq {
		if err != nil {
			return v, err
		}
		match, err := predicate(v)
		if err != nil {
			return v, err
		}
		if match {
			return v, nil
		}
	}
	var zero T
	return zero, ErrNoMatch
}

// FirstOrDefault returns the first element of the sequence that satisfies the predicate,
// or the zero value of T if no such element is found. Any element or predicate error is returned immediately.
func (s Sequence[T]) FirstOrDefault(predicate func(T) bool) (T, error) {
	return s.FirstOrDefaultWithError(func(v T) (bool, error) {
		return predicate(v), nil
	})
}

// FirstOrDefaultWithError returns the first element that satisfies the error-aware predicate,
// or the zero value of T if none match. Any sequence or predicate error is returned immediately.
func (s Sequence[T]) FirstOrDefaultWithError(predicate func(T) (bool, error)) (T, error) {
	for v, err := range s.seq {
		if err != nil {
			return v, err
		}
		match, err := predicate(v)
		if err != nil {
			return v, err
		}
		if match {
			return v, nil
		}
	}
	var zero T
	return zero, nil
}
