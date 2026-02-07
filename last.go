package iterlinq

// Last returns the last element of the sequence that satisfies the predicate.
// It enumerates the full sequence, tracking the last matching element.
func (s Sequence[T]) Last(predicate func(T) bool) (T, error) {
	return s.LastWithError(func(v T) (bool, error) {
		return predicate(v), nil
	})
}

// LastWithError returns the last element of the sequence that satisfies the error-aware predicate.
// If no element matches, ErrNoMatch is returned. Any element or predicate error is returned immediately.
func (s Sequence[T]) LastWithError(predicate func(T) (bool, error)) (T, error) {
	last, found, err := s.lastInternal(predicate)
	if err != nil {
		return last, err
	}
	if !found {
		return last, ErrNoMatch
	}
	return last, nil
}

// LastOrDefault returns the last element of the sequence that satisfies the predicate,
// or the zero value of T if no such element is found.
func (s Sequence[T]) LastOrDefault(predicate func(T) bool) (T, error) {
	return s.LastOrDefaultWithError(func(v T) (bool, error) {
		return predicate(v), nil
	})
}

// LastOrDefaultWithError returns the last element that satisfies the error-aware predicate,
// or the zero value of T if none match.
func (s Sequence[T]) LastOrDefaultWithError(predicate func(T) (bool, error)) (T, error) {
	last, _, err := s.lastInternal(predicate)
	return last, err
}

func (s Sequence[T]) lastInternal(predicate func(T) (bool, error)) (T, bool, error) {
	var last T
	found := false
	for v, err := range s.seq {
		if err != nil {
			return v, false, err
		}
		match, err := predicate(v)
		if err != nil {
			return v, false, err
		}
		if match {
			last = v
			found = true
		}
	}
	return last, found, nil
}
