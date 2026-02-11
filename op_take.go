package iterlinq

// Take returns a specified number of contiguous elements from the start of a sequence.
func (s Sequence[T]) Take(n int) Sequence[T] {
	if s.err != nil {
		return Sequence[T]{err: s.err}
	}
	return Sequence[T]{
		seq: func(yield func(T) bool) error {
			if n <= 0 {
				return nil
			}
			count := 0
			return s.seqOrEmpty()(func(v T) bool {
				if !yield(v) {
					return false
				}
				count++
				if count >= n {
					return false
				}
				return true
			})
		},
	}
}
