package iterlinq

// Count returns the number of elements in the sequence.
// If an element error is encountered during enumeration, it is returned immediately with the partial count.
func (s Sequence[T]) Count() (int, error) {
	count := 0
	for _, err := range s.seq {
		if err != nil {
			return count, err
		}
		count++
	}
	return count, nil
}
