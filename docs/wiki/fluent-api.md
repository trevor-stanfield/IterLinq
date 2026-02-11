# Fluent API (`iterlinq`)

## Constructors

### `From(seq iter.Seq[T]) Sequence[T]`
Wraps a standard `iter.Seq[T]` source.

### `FromFunc(fn func(yield func(T) bool) error) Sequence[T]`
Builds a lazy source that can fail during enumeration.

### `FromSlice(s []T) Sequence[T]`
Builds a sequence from a slice.

### `FromError(err error) Sequence[T]`
Creates an errored sequence that short-circuits until recovered.

## Error Flow

### `Err() error`
Returns attached monadic error.

### `HasError() bool`
True if sequence is in error state.

### `MapError(mapper func(error) error) Sequence[T]`
Maps an attached or runtime error.

### `Recover(handler func(error) Sequence[T]) Sequence[T]`
Replaces an errored sequence.

### `OrElse(fallback Sequence[T]) Sequence[T]`
Fallback convenience wrapper over `Recover`.

## Operators

### Filtering and Projection
- `Where(func(T) bool) Sequence[T]`
- `Select(func(T) any) Sequence[any]`
- `SelectMany(func(T) Sequence[any]) Sequence[any]`
- `DistinctBy(func(T) any) Sequence[T]`

### Slicing
- `Take(n int) Sequence[T]`
- `Skip(n int) Sequence[T]`

### Ordered Variant (Semi-Deferred)
- `MustOrderBy(func(T, T) int) OrderedSequence[T]`
- `MustOrderByDesc(func(T, T) int) OrderedSequence[T]`
- `ThenBy(func(T, T) int) OrderedSequence[T]`
- `ThenByDesc(func(T, T) int) OrderedSequence[T]`
- `Reverse() OrderedSequence[T]`
- `UsingSorter(func([]T, func(T, T) int) error) OrderedSequence[T]`
- `UsingSortFunc(func([]T, func(T, T) int) error) OrderedSequence[T]` (compatibility alias)
- `ForceOrder() Sequence[T]`

`OrderedSequence` accumulates ordering conditions lazily.
Calling any non-order operation (for example `Where`, `Select`, `Take`, or terminals) forces ordering first and then proceeds.

Default force behavior uses stable multi-pass merge sorting from the last `ThenBy` key to the first key (LSD-style key priority).
If `UsingSorter` is set, your sorter is used instead and receives the composed comparator.

## Terminals

### Materialization
- `ToSlice() ([]T, error)`

### Quantifiers/Aggregates
- `Count() (int, error)`
- `Any(func(T) bool) (bool, error)`
- `All(func(T) bool) (bool, error)`

### Element Selection
- `First(func(T) bool) (T, error)`
- `FirstOrDefault(func(T) bool) (T, error)`
- `Last(func(T) bool) (T, error)`
- `LastOrDefault(func(T) bool) (T, error)`

## Notes on Typing
Because methods cannot introduce new type parameters in Go, `Select` and `SelectMany` return `Sequence[any]`.

When strict transform typing matters, use [`rigid`](./rigid-api.md).

## Example
```go
res, err := iterlinq.FromSlice([]int{1,2,3,4,5,6}).
    Where(func(v int) bool { return v > 2 }).
    Select(func(v int) any { return v * 10 }).
    ToSlice()
```

## Related Docs
- [Error Model](./error-model.md)
- [Operators Reference](./operators-reference.md)
- [Rigid API](./rigid-api.md)
