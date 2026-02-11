package iterlinq

// ToSlice collects the sequence into a slice and returns it along with any error encountered.
// The enumeration stops at the first error and returns it alongside a nil slice.
func (s Sequence[T]) ToSlice() ([]T, error) {
	if s.err != nil {
		return nil, s.err
	}
	var result []T
	err := s.seqOrEmpty()(func(v T) bool {
		result = append(result, v)
		return true
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}
