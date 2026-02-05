package iterlinq

// All returns true if all elements satisfy the predicate. It stops early on the first false
// or on the first encountered error. If a predicate error occurs, it is returned.
func (s Sequence[T]) All(predicate func(T) (bool, error)) (bool, error) {
	for v, err := range s.seq {
		if err != nil {
			return false, err
		}
		match, err := predicate(v)
		if err != nil {
			return false, err
		}
		if !match {
			return false, nil
		}
	}
	return true, nil
}
