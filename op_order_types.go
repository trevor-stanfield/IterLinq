package iterlinq

// Comparator compares two values:
//   - negative if left < right
//   - zero if left == right
//   - positive if left > right
type Comparator[T any] func(left, right T) int

// CompareFunc is kept as a compatibility alias.
type CompareFunc[T any] = Comparator[T]

// Sorter sorts items in-place using the provided comparator.
// Implementations should honor the comparator and return an error on failure.
type Sorter[T any] func(items []T, comparator Comparator[T]) error

// SortFunc is kept as a compatibility alias.
type SortFunc[T any] = Sorter[T]

type orderKey[T any] struct {
	comparator   Comparator[T]
	isDescending bool
}

func composeOrderingComparator[T any](orderingKeys []orderKey[T]) Comparator[T] {
	return func(left, right T) int {
		for _, key := range orderingKeys {
			comparison := key.comparator(left, right)
			if key.isDescending {
				comparison = -comparison
			}
			if comparison != 0 {
				return comparison
			}
		}
		return 0
	}
}

func stableMergeSortUsingComparator[T any](items []T, comparator Comparator[T]) {
	if len(items) < 2 {
		return
	}
	buffer := make([]T, len(items))
	copy(buffer, items)
	stableMergeSortIntoBuffers[T](items, buffer, comparator)
}

func stableMergeSortIntoBuffers[T any](destination, source []T, comparator Comparator[T]) {
	if len(destination) < 2 {
		return
	}
	middleIndex := len(destination) / 2
	leftDestination, rightDestination := destination[:middleIndex], destination[middleIndex:]
	leftSource, rightSource := source[:middleIndex], source[middleIndex:]

	// Swap source/destination roles recursively for stable merge sorting.
	stableMergeSortIntoBuffers[T](leftSource, leftDestination, comparator)
	stableMergeSortIntoBuffers[T](rightSource, rightDestination, comparator)

	leftIndex, rightIndex, destinationIndex := 0, 0, 0
	for leftIndex < len(leftSource) && rightIndex < len(rightSource) {
		if comparator(leftSource[leftIndex], rightSource[rightIndex]) <= 0 {
			destination[destinationIndex] = leftSource[leftIndex]
			leftIndex++
		} else {
			destination[destinationIndex] = rightSource[rightIndex]
			rightIndex++
		}
		destinationIndex++
	}
	for leftIndex < len(leftSource) {
		destination[destinationIndex] = leftSource[leftIndex]
		leftIndex++
		destinationIndex++
	}
	for rightIndex < len(rightSource) {
		destination[destinationIndex] = rightSource[rightIndex]
		rightIndex++
		destinationIndex++
	}
	copy(source, destination)
}

func reverseInPlace[T any](items []T) {
	for leftIndex, rightIndex := 0, len(items)-1; leftIndex < rightIndex; leftIndex, rightIndex = leftIndex+1, rightIndex-1 {
		items[leftIndex], items[rightIndex] = items[rightIndex], items[leftIndex]
	}
}
