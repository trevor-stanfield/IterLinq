package iterlinq

// Any returns true if any element of the sequence satisfies the predicate.
// It stops at the first satisfying element. If an element error is encountered it is returned immediately.
func (s Sequence[T]) Any(predicate func(T) bool) (bool, error) {
	return s.AnyWithError(func(v T) (bool, error) {
		return predicate(v), nil
	})
}

// AnyWithError returns true if any element satisfies the error-aware predicate.
// It stops early on the first match or first error.
func (s Sequence[T]) AnyWithError(predicate func(T) (bool, error)) (bool, error) {
	for v, err := range s.seq {
		if err != nil {
			return false, err
		}
		match, err := predicate(v)
		if err != nil {
			return false, err
		}
		if match {
			return true, nil
		}
	}
	return false, nil
}
