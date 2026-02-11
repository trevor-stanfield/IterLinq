package iterlinq

// Skip bypasses a specified number of elements in a sequence and then returns the remaining elements.
func (s Sequence[T]) Skip(n int) Sequence[T] {
	if s.err != nil {
		return Sequence[T]{err: s.err}
	}
	return Sequence[T]{
		seq: func(yield func(T) bool) error {
			count := 0
			return s.seqOrEmpty()(func(v T) bool {
				if count < n {
					count++
					return true
				}
				return yield(v)
			})
		},
	}
}
