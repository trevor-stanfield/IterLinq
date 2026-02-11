package iterlinq

// OrderedSequence is a deferred ordering plan over a Sequence.
// Any non-order operation on OrderedSequence forces ordering before proceeding.
type OrderedSequence[T any] struct {
	sourceSequence Sequence[T]
	orderingKeys   []orderKey[T]
	sorter         Sorter[T]
	shouldReverse  bool
	err            error
}

// MustOrderBy starts an ordered variant that will sort by comparator when forced.
func (sequence Sequence[T]) MustOrderBy(comparator Comparator[T]) OrderedSequence[T] {
	if sequence.err != nil {
		return OrderedSequence[T]{sourceSequence: sequence, err: sequence.err}
	}
	if comparator == nil {
		return OrderedSequence[T]{sourceSequence: sequence, err: ErrNilComparator}
	}
	return OrderedSequence[T]{
		sourceSequence: sequence,
		orderingKeys:   []orderKey[T]{{comparator: comparator}},
	}
}

// MustOrderByDesc starts an ordered variant with descending primary order.
func (sequence Sequence[T]) MustOrderByDesc(comparator Comparator[T]) OrderedSequence[T] {
	if sequence.err != nil {
		return OrderedSequence[T]{sourceSequence: sequence, err: sequence.err}
	}
	if comparator == nil {
		return OrderedSequence[T]{sourceSequence: sequence, err: ErrNilComparator}
	}
	return OrderedSequence[T]{
		sourceSequence: sequence,
		orderingKeys:   []orderKey[T]{{comparator: comparator, isDescending: true}},
	}
}

// ThenBy appends a secondary ascending ordering key.
func (orderedSequence OrderedSequence[T]) ThenBy(comparator Comparator[T]) OrderedSequence[T] {
	if orderedSequence.err != nil {
		return orderedSequence
	}
	if comparator == nil {
		orderedSequence.err = ErrNilComparator
		return orderedSequence
	}
	orderedSequence.orderingKeys = append(cloneOrderingKeys(orderedSequence.orderingKeys), orderKey[T]{comparator: comparator})
	return orderedSequence
}

// ThenByDesc appends a secondary descending ordering key.
func (orderedSequence OrderedSequence[T]) ThenByDesc(comparator Comparator[T]) OrderedSequence[T] {
	if orderedSequence.err != nil {
		return orderedSequence
	}
	if comparator == nil {
		orderedSequence.err = ErrNilComparator
		return orderedSequence
	}
	orderedSequence.orderingKeys = append(cloneOrderingKeys(orderedSequence.orderingKeys), orderKey[T]{comparator: comparator, isDescending: true})
	return orderedSequence
}

// Reverse reverses the final ordered output.
func (orderedSequence OrderedSequence[T]) Reverse() OrderedSequence[T] {
	if orderedSequence.err != nil {
		return orderedSequence
	}
	orderedSequence.shouldReverse = !orderedSequence.shouldReverse
	return orderedSequence
}

// UsingSorter configures a custom sorter used at force time.
// The sorter receives the composed comparator from MustOrderBy/ThenBy.
func (orderedSequence OrderedSequence[T]) UsingSorter(sorter Sorter[T]) OrderedSequence[T] {
	if orderedSequence.err != nil {
		return orderedSequence
	}
	if sorter == nil {
		orderedSequence.err = ErrNilSortFunc
		return orderedSequence
	}
	orderedSequence.sorter = sorter
	return orderedSequence
}

// UsingSortFunc is kept as a compatibility wrapper over UsingSorter.
func (orderedSequence OrderedSequence[T]) UsingSortFunc(sortFunc SortFunc[T]) OrderedSequence[T] {
	return orderedSequence.UsingSorter(sortFunc)
}

// Err returns any attached error state on the ordered sequence.
func (orderedSequence OrderedSequence[T]) Err() error {
	if orderedSequence.err != nil {
		return orderedSequence.err
	}
	return orderedSequence.sourceSequence.err
}

// HasError reports whether the ordered sequence has an attached error.
func (orderedSequence OrderedSequence[T]) HasError() bool {
	return orderedSequence.Err() != nil
}

// ForceOrder materializes the source, applies ordering, and returns a normal Sequence.
func (orderedSequence OrderedSequence[T]) ForceOrder() Sequence[T] {
	if orderedSequence.err != nil {
		return failed[T](orderedSequence.err)
	}
	if orderedSequence.sourceSequence.err != nil {
		return failed[T](orderedSequence.sourceSequence.err)
	}
	return FromFunc(func(yield func(T) bool) error {
		collectedValues, err := orderedSequence.sourceSequence.ToSlice()
		if err != nil {
			return err
		}

		orderedValues := make([]T, len(collectedValues))
		copy(orderedValues, collectedValues)
		if err := orderedSequence.applyDeferredOrdering(orderedValues); err != nil {
			return err
		}

		for _, value := range orderedValues {
			if !yield(value) {
				return nil
			}
		}
		return nil
	})
}

func (orderedSequence OrderedSequence[T]) applyDeferredOrdering(items []T) error {
	if len(orderedSequence.orderingKeys) > 0 {
		if orderedSequence.sorter != nil {
			if err := orderedSequence.sorter(items, composeOrderingComparator(orderedSequence.orderingKeys)); err != nil {
				return err
			}
		} else {
			for keyIndex := len(orderedSequence.orderingKeys) - 1; keyIndex >= 0; keyIndex-- {
				orderingKey := orderedSequence.orderingKeys[keyIndex]
				stableMergeSortUsingComparator(items, func(left, right T) int {
					comparison := orderingKey.comparator(left, right)
					if orderingKey.isDescending {
						comparison = -comparison
					}
					return comparison
				})
			}
		}
	}
	if orderedSequence.shouldReverse {
		reverseInPlace(items)
	}
	return nil
}

func cloneOrderingKeys[T any](orderingKeys []orderKey[T]) []orderKey[T] {
	if len(orderingKeys) == 0 {
		return nil
	}
	clonedKeys := make([]orderKey[T], len(orderingKeys))
	copy(clonedKeys, orderingKeys)
	return clonedKeys
}

// MapError forces ordering before applying Sequence error mapping behavior.
func (orderedSequence OrderedSequence[T]) MapError(mapper func(error) error) Sequence[T] {
	return orderedSequence.ForceOrder().MapError(mapper)
}

// Recover forces ordering before applying Sequence recovery behavior.
func (orderedSequence OrderedSequence[T]) Recover(handler func(error) Sequence[T]) Sequence[T] {
	return orderedSequence.ForceOrder().Recover(handler)
}

// OrElse forces ordering before applying Sequence fallback behavior.
func (orderedSequence OrderedSequence[T]) OrElse(fallback Sequence[T]) Sequence[T] {
	return orderedSequence.ForceOrder().OrElse(fallback)
}

// Where forces ordering before filtering.
func (orderedSequence OrderedSequence[T]) Where(predicate func(T) bool) Sequence[T] {
	return orderedSequence.ForceOrder().Where(predicate)
}

// Select forces ordering before mapping.
func (orderedSequence OrderedSequence[T]) Select(transform func(T) any) Sequence[any] {
	return orderedSequence.ForceOrder().Select(transform)
}

// SelectMany forces ordering before flattening.
func (orderedSequence OrderedSequence[T]) SelectMany(transform func(T) Sequence[any]) Sequence[any] {
	return orderedSequence.ForceOrder().SelectMany(transform)
}

// DistinctBy forces ordering before distinct key selection.
func (orderedSequence OrderedSequence[T]) DistinctBy(keySelector func(T) any) Sequence[T] {
	return orderedSequence.ForceOrder().DistinctBy(keySelector)
}

// Take forces ordering before taking n elements.
func (orderedSequence OrderedSequence[T]) Take(count int) Sequence[T] {
	return orderedSequence.ForceOrder().Take(count)
}

// Skip forces ordering before skipping n elements.
func (orderedSequence OrderedSequence[T]) Skip(count int) Sequence[T] {
	return orderedSequence.ForceOrder().Skip(count)
}

// ToSlice forces ordering and materializes to a slice.
func (orderedSequence OrderedSequence[T]) ToSlice() ([]T, error) {
	return orderedSequence.ForceOrder().ToSlice()
}

// Count forces ordering before counting.
func (orderedSequence OrderedSequence[T]) Count() (int, error) {
	return orderedSequence.ForceOrder().Count()
}

// Any forces ordering before predicate check.
func (orderedSequence OrderedSequence[T]) Any(predicate func(T) bool) (bool, error) {
	return orderedSequence.ForceOrder().Any(predicate)
}

// All forces ordering before predicate check.
func (orderedSequence OrderedSequence[T]) All(predicate func(T) bool) (bool, error) {
	return orderedSequence.ForceOrder().All(predicate)
}

// First forces ordering before first-match query.
func (orderedSequence OrderedSequence[T]) First(predicate func(T) bool) (T, error) {
	return orderedSequence.ForceOrder().First(predicate)
}

// FirstOrDefault forces ordering before first-match-or-default query.
func (orderedSequence OrderedSequence[T]) FirstOrDefault(predicate func(T) bool) (T, error) {
	return orderedSequence.ForceOrder().FirstOrDefault(predicate)
}

// Last forces ordering before last-match query.
func (orderedSequence OrderedSequence[T]) Last(predicate func(T) bool) (T, error) {
	return orderedSequence.ForceOrder().Last(predicate)
}

// LastOrDefault forces ordering before last-match-or-default query.
func (orderedSequence OrderedSequence[T]) LastOrDefault(predicate func(T) bool) (T, error) {
	return orderedSequence.ForceOrder().LastOrDefault(predicate)
}
