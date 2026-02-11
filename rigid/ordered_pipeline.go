package rigid

// OrderedPipeline stores deferred ordering over a pipeline output stage.
// Any non-order operation on OrderedPipeline forces ordering first.
type OrderedPipeline[In, Out any] struct {
	basePipeline  Pipeline[In, Out]
	orderingKeys  []orderKey[Out]
	sorter        Sorter[Out]
	shouldReverse bool
	err           error
}

// MustOrderBy starts an ordered pipeline by the supplied comparator.
func (pipeline Pipeline[In, Out]) MustOrderBy(comparator Comparator[Out]) OrderedPipeline[In, Out] {
	if pipeline.err != nil {
		return OrderedPipeline[In, Out]{basePipeline: pipeline, err: pipeline.err}
	}
	if comparator == nil {
		return OrderedPipeline[In, Out]{basePipeline: pipeline, err: ErrNilComparator}
	}
	return OrderedPipeline[In, Out]{
		basePipeline: pipeline,
		orderingKeys: []orderKey[Out]{{comparator: comparator}},
	}
}

// MustOrderByDesc starts an ordered pipeline with descending primary order.
func (pipeline Pipeline[In, Out]) MustOrderByDesc(comparator Comparator[Out]) OrderedPipeline[In, Out] {
	if pipeline.err != nil {
		return OrderedPipeline[In, Out]{basePipeline: pipeline, err: pipeline.err}
	}
	if comparator == nil {
		return OrderedPipeline[In, Out]{basePipeline: pipeline, err: ErrNilComparator}
	}
	return OrderedPipeline[In, Out]{
		basePipeline: pipeline,
		orderingKeys: []orderKey[Out]{{comparator: comparator, isDescending: true}},
	}
}

// ThenBy appends a secondary ascending ordering key.
func (orderedPipeline OrderedPipeline[In, Out]) ThenBy(comparator Comparator[Out]) OrderedPipeline[In, Out] {
	if orderedPipeline.err != nil {
		return orderedPipeline
	}
	if comparator == nil {
		orderedPipeline.err = ErrNilComparator
		return orderedPipeline
	}
	orderedPipeline.orderingKeys = append(cloneOrderingKeys(orderedPipeline.orderingKeys), orderKey[Out]{comparator: comparator})
	return orderedPipeline
}

// ThenByDesc appends a secondary descending ordering key.
func (orderedPipeline OrderedPipeline[In, Out]) ThenByDesc(comparator Comparator[Out]) OrderedPipeline[In, Out] {
	if orderedPipeline.err != nil {
		return orderedPipeline
	}
	if comparator == nil {
		orderedPipeline.err = ErrNilComparator
		return orderedPipeline
	}
	orderedPipeline.orderingKeys = append(cloneOrderingKeys(orderedPipeline.orderingKeys), orderKey[Out]{comparator: comparator, isDescending: true})
	return orderedPipeline
}

// Reverse reverses the final ordered output.
func (orderedPipeline OrderedPipeline[In, Out]) Reverse() OrderedPipeline[In, Out] {
	if orderedPipeline.err != nil {
		return orderedPipeline
	}
	orderedPipeline.shouldReverse = !orderedPipeline.shouldReverse
	return orderedPipeline
}

// UsingSorter configures a custom sorter used at force time.
func (orderedPipeline OrderedPipeline[In, Out]) UsingSorter(sorter Sorter[Out]) OrderedPipeline[In, Out] {
	if orderedPipeline.err != nil {
		return orderedPipeline
	}
	if sorter == nil {
		orderedPipeline.err = ErrNilSortFunc
		return orderedPipeline
	}
	orderedPipeline.sorter = sorter
	return orderedPipeline
}

// UsingSortFunc is kept as a compatibility wrapper over UsingSorter.
func (orderedPipeline OrderedPipeline[In, Out]) UsingSortFunc(sortFunc SortFunc[Out]) OrderedPipeline[In, Out] {
	return orderedPipeline.UsingSorter(sortFunc)
}

// Err returns any attached error state on the ordered pipeline.
func (orderedPipeline OrderedPipeline[In, Out]) Err() error {
	if orderedPipeline.err != nil {
		return orderedPipeline.err
	}
	return orderedPipeline.basePipeline.err
}

// HasError reports whether the ordered pipeline has an attached error.
func (orderedPipeline OrderedPipeline[In, Out]) HasError() bool {
	return orderedPipeline.Err() != nil
}

// ForceOrder appends ordering into the pipeline and returns a regular pipeline.
func (orderedPipeline OrderedPipeline[In, Out]) ForceOrder() Pipeline[In, Out] {
	if orderedPipeline.err != nil {
		return Pipeline[In, Out]{err: orderedPipeline.err}
	}
	if orderedPipeline.basePipeline.err != nil {
		return Pipeline[In, Out]{err: orderedPipeline.basePipeline.err}
	}
	if orderedPipeline.basePipeline.trans == nil {
		return Pipeline[In, Out]{err: ErrNilTransformer}
	}
	return Pipeline[In, Out]{
		source: orderedPipeline.basePipeline.source,
		trans: Compose(
			orderedPipeline.basePipeline.trans,
			createOrderedTransformer(orderedPipeline.orderingKeys, orderedPipeline.sorter, orderedPipeline.shouldReverse),
		),
	}
}

func createOrderedTransformer[T any](orderingKeys []orderKey[T], sorter Sorter[T], shouldReverse bool) Transformer[T, T] {
	clonedOrderingKeys := cloneOrderingKeys(orderingKeys)
	return TransformerFunc[T, T](func(input []T) ([]T, error) {
		output := make([]T, len(input))
		copy(output, input)
		if len(clonedOrderingKeys) > 0 {
			if sorter != nil {
				if err := sorter(output, composeOrderingComparator(clonedOrderingKeys)); err != nil {
					return nil, err
				}
			} else {
				for keyIndex := len(clonedOrderingKeys) - 1; keyIndex >= 0; keyIndex-- {
					orderingKey := clonedOrderingKeys[keyIndex]
					stableMergeSortUsingComparator(output, func(left, right T) int {
						comparison := orderingKey.comparator(left, right)
						if orderingKey.isDescending {
							comparison = -comparison
						}
						return comparison
					})
				}
			}
		}
		if shouldReverse {
			reverseInPlace(output)
		}
		return output, nil
	})
}

// ThenOrdered forces ordering, then appends the next transformer.
func ThenOrdered[A, B, C any](orderedPipeline OrderedPipeline[A, B], nextTransformer Transformer[B, C]) Pipeline[A, C] {
	return Then(orderedPipeline.ForceOrder(), nextTransformer)
}

// SelectOrdered forces ordering, then applies Select.
func SelectOrdered[A, B, C any](orderedPipeline OrderedPipeline[A, B], selector func(B) (C, error)) Pipeline[A, C] {
	return ThenOrdered(orderedPipeline, Select(selector))
}

// SelectManyOrdered forces ordering, then applies SelectMany.
func SelectManyOrdered[A, B, C any](orderedPipeline OrderedPipeline[A, B], selector func(B) ([]C, error)) Pipeline[A, C] {
	return ThenOrdered(orderedPipeline, SelectMany(selector))
}

// DistinctByOrdered forces ordering, then applies DistinctBy.
func DistinctByOrdered[A, B any, K comparable](orderedPipeline OrderedPipeline[A, B], selector func(B) (K, error)) Pipeline[A, B] {
	return ThenOrdered(orderedPipeline, DistinctBy(selector))
}

// Where forces ordering before filtering.
func (orderedPipeline OrderedPipeline[In, Out]) Where(predicate func(Out) (bool, error)) Pipeline[In, Out] {
	return ThenOrdered(orderedPipeline, Where(predicate))
}

// Take forces ordering before taking n elements.
func (orderedPipeline OrderedPipeline[In, Out]) Take(count int) Pipeline[In, Out] {
	return ThenOrdered(orderedPipeline, Take[Out](count))
}

// Skip forces ordering before skipping n elements.
func (orderedPipeline OrderedPipeline[In, Out]) Skip(count int) Pipeline[In, Out] {
	return ThenOrdered(orderedPipeline, Skip[Out](count))
}

// MapError forces ordering before error mapping.
func (orderedPipeline OrderedPipeline[In, Out]) MapError(mapper func(error) error) Pipeline[In, Out] {
	return orderedPipeline.ForceOrder().MapError(mapper)
}

// Recover forces ordering before recovery.
func (orderedPipeline OrderedPipeline[In, Out]) Recover(handler func(error) Pipeline[In, Out]) Pipeline[In, Out] {
	return orderedPipeline.ForceOrder().Recover(handler)
}

// OrElse forces ordering before fallback.
func (orderedPipeline OrderedPipeline[In, Out]) OrElse(fallback Pipeline[In, Out]) Pipeline[In, Out] {
	return orderedPipeline.ForceOrder().OrElse(fallback)
}

// Execute forces ordering and materializes the pipeline.
func (orderedPipeline OrderedPipeline[In, Out]) Execute() ([]Out, error) {
	return orderedPipeline.ForceOrder().Execute()
}

// ToSlice forces ordering and materializes the pipeline.
func (orderedPipeline OrderedPipeline[In, Out]) ToSlice() ([]Out, error) {
	return orderedPipeline.ForceOrder().ToSlice()
}

// Count forces ordering before counting.
func (orderedPipeline OrderedPipeline[In, Out]) Count() (int, error) {
	return orderedPipeline.ForceOrder().Count()
}

// Any forces ordering before predicate check.
func (orderedPipeline OrderedPipeline[In, Out]) Any(predicate func(Out) (bool, error)) (bool, error) {
	return orderedPipeline.ForceOrder().Any(predicate)
}

// All forces ordering before predicate check.
func (orderedPipeline OrderedPipeline[In, Out]) All(predicate func(Out) (bool, error)) (bool, error) {
	return orderedPipeline.ForceOrder().All(predicate)
}

// First forces ordering before first-match query.
func (orderedPipeline OrderedPipeline[In, Out]) First(predicate func(Out) (bool, error)) (Out, error) {
	return orderedPipeline.ForceOrder().First(predicate)
}

// FirstOrDefault forces ordering before first-match-or-default query.
func (orderedPipeline OrderedPipeline[In, Out]) FirstOrDefault(predicate func(Out) (bool, error)) (Out, error) {
	return orderedPipeline.ForceOrder().FirstOrDefault(predicate)
}

// Last forces ordering before last-match query.
func (orderedPipeline OrderedPipeline[In, Out]) Last(predicate func(Out) (bool, error)) (Out, error) {
	return orderedPipeline.ForceOrder().Last(predicate)
}

// LastOrDefault forces ordering before last-match-or-default query.
func (orderedPipeline OrderedPipeline[In, Out]) LastOrDefault(predicate func(Out) (bool, error)) (Out, error) {
	return orderedPipeline.ForceOrder().LastOrDefault(predicate)
}
