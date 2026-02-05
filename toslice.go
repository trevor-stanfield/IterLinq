package iterlinq

// ToSlice collects the sequence into a slice and returns it along with any error encountered.
// The enumeration stops at the first error and returns it alongside a nil slice.
func (s Sequence[T]) ToSlice() ([]T, error) {
	var result []T
	for v, err := range s.seq {
		if err != nil {
			return nil, err
		}
		result = append(result, v)
	}
	return result, nil
}
