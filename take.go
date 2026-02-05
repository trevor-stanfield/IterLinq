package iterlinq

// Take returns a specified number of contiguous elements from the start of a sequence.
func (s Sequence[T]) Take(n int) Sequence[T] {
	return Sequence[T]{
		seq: func(yield func(T, error) bool) {
			if n <= 0 {
				return
			}
			count := 0
			for v, err := range s.seq {
				if err != nil {
					yield(v, err)
					return
				}
				if !yield(v, nil) {
					return
				}
				count++
				if count >= n {
					return
				}
			}
		},
	}
}
