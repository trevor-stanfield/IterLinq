package iterlinq

// Any returns true if any element of the sequence satisfies the predicate.
// It stops at the first satisfying element. If an element error is encountered it is returned immediately.
func (s Sequence[T]) Any(predicate func(T) bool) (bool, error) {
	if s.err != nil {
		return false, s.err
	}
	found := false
	err := s.seqOrEmpty()(func(v T) bool {
		if predicate(v) {
			found = true
			return false
		}
		return true
	})
	if err != nil {
		return false, err
	}
	return found, nil
}
