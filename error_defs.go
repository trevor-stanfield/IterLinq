package iterlinq

import "errors"

// ErrNoMatch is returned when no element satisfies the predicate.
var ErrNoMatch = errors.New("no element satisfies the predicate")

// ErrNonComparableKey is returned when a key selector yields a non-comparable key.
var ErrNonComparableKey = errors.New("key is not comparable")

// ErrNilErrorMapper is returned when MapError is called with a nil mapper.
var ErrNilErrorMapper = errors.New("error mapper cannot be nil")

// ErrNilErrorHandler is returned when Recover is called with a nil handler.
var ErrNilErrorHandler = errors.New("error handler cannot be nil")

// ErrNilComparator is returned when an ordering comparator is nil.
var ErrNilComparator = errors.New("comparator cannot be nil")

// ErrNilSortFunc is returned when UsingSorter or UsingSortFunc is called with a nil sorter.
var ErrNilSortFunc = errors.New("sort function cannot be nil")
