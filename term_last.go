package iterlinq

// Last returns the last element of the sequence that satisfies the predicate.
// It enumerates the full sequence, tracking the last matching element.
func (s Sequence[T]) Last(predicate func(T) bool) (T, error) {
	if s.err != nil {
		var zero T
		return zero, s.err
	}
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
	if s.err != nil {
		var zero T
		return zero, s.err
	}
	last, _, err := s.lastInternal(predicate)
	return last, err
}

func (s Sequence[T]) lastInternal(predicate func(T) bool) (T, bool, error) {
	var last T
	found := false
	err := s.seqOrEmpty()(func(v T) bool {
		if predicate(v) {
			last = v
			found = true
		}
		return true
	})
	if err != nil {
		var zero T
		return zero, false, err
	}
	return last, found, nil
}
