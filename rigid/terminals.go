package rigid

// Count returns the number of elements after executing the pipeline.
func (p Pipeline[In, Out]) Count() (int, error) {
	values, err := p.Execute()
	if err != nil {
		return 0, err
	}
	return len(values), nil
}

// Any returns true if any element satisfies the predicate.
func (p Pipeline[In, Out]) Any(predicate func(Out) (bool, error)) (bool, error) {
	if predicate == nil {
		return false, ErrNilPredicate
	}
	values, err := p.Execute()
	if err != nil {
		return false, err
	}
	for _, v := range values {
		match, err := predicate(v)
		if err != nil {
			return false, err
		}
		if match {
			return true, nil
		}
	}
	return false, nil
}

// All returns true if all elements satisfy the predicate.
func (p Pipeline[In, Out]) All(predicate func(Out) (bool, error)) (bool, error) {
	if predicate == nil {
		return false, ErrNilPredicate
	}
	values, err := p.Execute()
	if err != nil {
		return false, err
	}
	for _, v := range values {
		match, err := predicate(v)
		if err != nil {
			return false, err
		}
		if !match {
			return false, nil
		}
	}
	return true, nil
}

// First returns the first element that satisfies the predicate.
func (p Pipeline[In, Out]) First(predicate func(Out) (bool, error)) (Out, error) {
	if predicate == nil {
		var zero Out
		return zero, ErrNilPredicate
	}
	values, err := p.Execute()
	if err != nil {
		var zero Out
		return zero, err
	}
	for _, v := range values {
		match, err := predicate(v)
		if err != nil {
			var zero Out
			return zero, err
		}
		if match {
			return v, nil
		}
	}
	var zero Out
	return zero, ErrNoMatch
}

// FirstOrDefault returns the first matching element or the zero value.
func (p Pipeline[In, Out]) FirstOrDefault(predicate func(Out) (bool, error)) (Out, error) {
	if predicate == nil {
		var zero Out
		return zero, ErrNilPredicate
	}
	values, err := p.Execute()
	if err != nil {
		var zero Out
		return zero, err
	}
	for _, v := range values {
		match, err := predicate(v)
		if err != nil {
			var zero Out
			return zero, err
		}
		if match {
			return v, nil
		}
	}
	var zero Out
	return zero, nil
}

// Last returns the last element that satisfies the predicate.
func (p Pipeline[In, Out]) Last(predicate func(Out) (bool, error)) (Out, error) {
	if predicate == nil {
		var zero Out
		return zero, ErrNilPredicate
	}
	values, err := p.Execute()
	if err != nil {
		var zero Out
		return zero, err
	}
	var (
		last  Out
		found bool
	)
	for _, v := range values {
		match, err := predicate(v)
		if err != nil {
			var zero Out
			return zero, err
		}
		if match {
			last = v
			found = true
		}
	}
	if !found {
		var zero Out
		return zero, ErrNoMatch
	}
	return last, nil
}

// LastOrDefault returns the last matching element or the zero value.
func (p Pipeline[In, Out]) LastOrDefault(predicate func(Out) (bool, error)) (Out, error) {
	if predicate == nil {
		var zero Out
		return zero, ErrNilPredicate
	}
	values, err := p.Execute()
	if err != nil {
		var zero Out
		return zero, err
	}
	var (
		last  Out
		found bool
	)
	for _, v := range values {
		match, err := predicate(v)
		if err != nil {
			var zero Out
			return zero, err
		}
		if match {
			last = v
			found = true
		}
	}
	if !found {
		var zero Out
		return zero, nil
	}
	return last, nil
}
