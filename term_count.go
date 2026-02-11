package iterlinq

// Count returns the number of elements in the sequence.
// If an element error is encountered during enumeration, it is returned immediately with the partial count.
func (s Sequence[T]) Count() (int, error) {
	if s.err != nil {
		return 0, s.err
	}
	count := 0
	err := s.seqOrEmpty()(func(T) bool {
		count++
		return true
	})
	if err != nil {
		return count, err
	}
	return count, nil
}
