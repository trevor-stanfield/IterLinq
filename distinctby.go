package iterlinq

import "reflect"

// DistinctBy returns a sequence that contains the first occurrence of each element based on the provided key selector.
// The key must be comparable; if a non-comparable key is produced, iteration stops and returns ErrNonComparableKey.
func (s Sequence[T]) DistinctBy(key func(T) any) Sequence[T] {
	return Sequence[T]{
		seq: func(yield func(T, error) bool) {
			seen := make(map[any]struct{})
			for v, err := range s.seq {
				if err != nil {
					yield(v, err)
					return
				}
				k := key(v)
				rt := reflect.TypeOf(k)
				if rt == nil || !rt.Comparable() {
					yield(v, ErrNonComparableKey)
					return
				}
				if _, ok := seen[k]; !ok {
					seen[k] = struct{}{}
					if !yield(v, nil) {
						return
					}
				}
			}
		},
	}
}
