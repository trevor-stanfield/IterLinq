package iterlinq

// Skip bypasses a specified number of elements in a sequence and then returns the remaining elements.
func (s Sequence[T]) Skip(n int) Sequence[T] {
	return Sequence[T]{
		seq: func(yield func(T, error) bool) {
			count := 0
			for v, err := range s.seq {
				if err != nil {
					yield(v, err)
					return
				}
				if count < n {
					count++
					continue
				}
				if !yield(v, nil) {
					return
				}
			}
		},
	}
}
