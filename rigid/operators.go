package rigid

// Where filters a collection based on an error-aware predicate.
func Where[T any](predicate func(T) (bool, error)) Transformer[T, T] {
	if predicate == nil {
		return failed[T, T](ErrNilPredicate)
	}
	return TransformerFunc[T, T](func(in []T) ([]T, error) {
		out := make([]T, 0, len(in))
		for _, v := range in {
			match, err := predicate(v)
			if err != nil {
				return nil, err
			}
			if match {
				out = append(out, v)
			}
		}
		return out, nil
	})
}

// Select maps each element using an error-aware selector.
func Select[T, R any](selector func(T) (R, error)) Transformer[T, R] {
	if selector == nil {
		return failed[T, R](ErrNilSelector)
	}
	return TransformerFunc[T, R](func(in []T) ([]R, error) {
		out := make([]R, 0, len(in))
		for _, v := range in {
			mapped, err := selector(v)
			if err != nil {
				return nil, err
			}
			out = append(out, mapped)
		}
		return out, nil
	})
}

// SelectMany maps each element to a slice and flattens the result.
func SelectMany[T, R any](selector func(T) ([]R, error)) Transformer[T, R] {
	if selector == nil {
		return failed[T, R](ErrNilSelector)
	}
	return TransformerFunc[T, R](func(in []T) ([]R, error) {
		out := make([]R, 0)
		for _, v := range in {
			mapped, err := selector(v)
			if err != nil {
				return nil, err
			}
			out = append(out, mapped...)
		}
		return out, nil
	})
}

// Take keeps the first n elements.
func Take[T any](n int) Transformer[T, T] {
	return TransformerFunc[T, T](func(in []T) ([]T, error) {
		if n <= 0 || len(in) == 0 {
			return []T{}, nil
		}
		if n > len(in) {
			n = len(in)
		}
		out := make([]T, n)
		copy(out, in[:n])
		return out, nil
	})
}

// Skip drops the first n elements.
func Skip[T any](n int) Transformer[T, T] {
	return TransformerFunc[T, T](func(in []T) ([]T, error) {
		if n <= 0 {
			out := make([]T, len(in))
			copy(out, in)
			return out, nil
		}
		if n >= len(in) {
			return []T{}, nil
		}
		out := make([]T, len(in)-n)
		copy(out, in[n:])
		return out, nil
	})
}

// DistinctBy keeps the first element for each key.
func DistinctBy[T any, K comparable](selector func(T) (K, error)) Transformer[T, T] {
	if selector == nil {
		return failed[T, T](ErrNilSelector)
	}
	return TransformerFunc[T, T](func(in []T) ([]T, error) {
		seen := make(map[K]struct{}, len(in))
		out := make([]T, 0, len(in))
		for _, v := range in {
			k, err := selector(v)
			if err != nil {
				return nil, err
			}
			if _, ok := seen[k]; ok {
				continue
			}
			seen[k] = struct{}{}
			out = append(out, v)
		}
		return out, nil
	})
}
