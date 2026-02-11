package iterlinq

import "reflect"

// DistinctBy returns a sequence that contains the first occurrence of each element based on the provided key selector.
// The key must be comparable; if a non-comparable key is produced, iteration stops and returns ErrNonComparableKey.
func (s Sequence[T]) DistinctBy(key func(T) any) Sequence[T] {
	if s.err != nil {
		return Sequence[T]{err: s.err}
	}
	return Sequence[T]{
		seq: func(yield func(T) bool) error {
			seen := make(map[any]struct{})
			var keyErr error
			err := s.seqOrEmpty()(func(v T) bool {
				k := key(v)
				rt := reflect.TypeOf(k)
				if rt == nil || !rt.Comparable() {
					keyErr = ErrNonComparableKey
					return false
				}
				if _, ok := seen[k]; !ok {
					seen[k] = struct{}{}
					return yield(v)
				}
				return true
			})
			if keyErr != nil {
				return keyErr
			}
			return err
		},
	}
}
