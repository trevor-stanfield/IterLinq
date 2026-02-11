package iterlinq

// First returns the first element of the sequence that satisfies the predicate.
// It short-circuits on the first match. If an element error is encountered, it is returned immediately.
// If no element matches, it returns ErrNoMatch.
func (s Sequence[T]) First(predicate func(T) bool) (T, error) {
	if s.err != nil {
		var zero T
		return zero, s.err
	}
	var (
		found bool
		first T
	)
	err := s.seqOrEmpty()(func(v T) bool {
		if predicate(v) {
			first = v
			found = true
			return false
		}
		return true
	})
	if err != nil {
		var zero T
		return zero, err
	}
	if found {
		return first, nil
	}
	var zero T
	return zero, ErrNoMatch
}

// FirstOrDefault returns the first element of the sequence that satisfies the predicate,
// or the zero value of T if no such element is found. Any sequence error is returned immediately.
func (s Sequence[T]) FirstOrDefault(predicate func(T) bool) (T, error) {
	if s.err != nil {
		var zero T
		return zero, s.err
	}
	var (
		found bool
		first T
	)
	err := s.seqOrEmpty()(func(v T) bool {
		if predicate(v) {
			first = v
			found = true
			return false
		}
		return true
	})
	if err != nil {
		var zero T
		return zero, err
	}
	if found {
		return first, nil
	}
	var zero T
	return zero, nil
}
