package iterlinq

// All returns true if all elements satisfy the predicate. It stops early on the first false
// or on the first encountered sequence error.
func (s Sequence[T]) All(predicate func(T) bool) (bool, error) {
	if s.err != nil {
		return false, s.err
	}
	all := true
	err := s.seqOrEmpty()(func(v T) bool {
		if !predicate(v) {
			all = false
			return false
		}
		return true
	})
	if err != nil {
		return false, err
	}
	return all, nil
}
