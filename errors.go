package iterlinq

import "errors"

// ErrNoMatch is returned when no element satisfies the predicate.
var ErrNoMatch = errors.New("no element satisfies the predicate")

// ErrNonComparableKey is returned when a key selector yields a non-comparable key.
var ErrNonComparableKey = errors.New("key is not comparable")
