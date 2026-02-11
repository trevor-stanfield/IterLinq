# Operators Reference

## Fluent Operators (`iterlinq`)

### `Where(predicate)`
- Type: `Sequence[T] -> Sequence[T]`
- Behavior: keeps values where predicate is true.
- Complexity: O(n)

### `Select(transform)`
- Type: `Sequence[T] -> Sequence[any]`
- Behavior: maps each value.
- Complexity: O(n)

### `SelectMany(transform)`
- Type: `Sequence[T] -> Sequence[any]`
- Behavior: maps each value to a sequence and flattens.
- Complexity: O(total flattened values)

### `DistinctBy(keySelector)`
- Type: `Sequence[T] -> Sequence[T]`
- Behavior: keeps first value per key.
- Complexity: O(n) average, hash-map backed.
- Error: fails if key is not comparable.

### `Take(n)`
- Type: `Sequence[T] -> Sequence[T]`
- Behavior: first `n` values.
- Complexity: O(min(n, len))

### `Skip(n)`
- Type: `Sequence[T] -> Sequence[T]`
- Behavior: drops first `n` values.
- Complexity: O(n + remaining)

### Ordered Semi-Deferred Variant
- `MustOrderBy(comparator)` / `MustOrderByDesc(comparator)` create `OrderedSequence[T]`.
- `ThenBy(comparator)` / `ThenByDesc(comparator)` append secondary keys.
- `Reverse()` reverses final ordered output.
- `UsingSorter(sorter)` installs custom forced-sort implementation.
- `UsingSortFunc(sorter)` remains available as a compatibility alias.
- Any non-order call on `OrderedSequence` forces ordering first.
- Default forced ordering: stable multi-pass merge sorting from last key to first key (LSD key priority).

## Fluent Terminals

### `ToSlice()`
Materializes the sequence.

### `Count()`
Counts values post-pipeline.

### `Any(predicate)` / `All(predicate)`
Quantifier checks with short-circuit behavior.

### `First(...)`, `FirstOrDefault(...)`, `Last(...)`, `LastOrDefault(...)`
Element queries with `ErrNoMatch` semantics for non-default variants.

## Rigid Operators (`iterlinq/rigid`)
Equivalent conceptual operators exist with typed signatures and error-aware callbacks.

### `Where(func(T) (bool, error))`
### `Select(func(T) (R, error))`
### `SelectMany(func(T) ([]R, error))`
### `DistinctBy(func(T) (K, error))`
### `Take(n)`
### `Skip(n)`

### Ordered Semi-Deferred Variant
- `Pipeline.MustOrderBy(comparator)` / `Pipeline.MustOrderByDesc(comparator)` produce `OrderedPipeline`.
- `ThenBy`, `ThenByDesc`, `Reverse`, `UsingSorter`, `ForceOrder` mirror fluent semantics.
- `UsingSortFunc` remains available as a compatibility alias.
- For type-changing stages, use `ThenOrdered`, `SelectOrdered`, `SelectManyOrdered`, `DistinctByOrdered`.

## Rigid Terminals
- `Execute` / `ToSlice`
- `Count`
- `Any` / `All`
- `First` / `FirstOrDefault`
- `Last` / `LastOrDefault`

## Operator Selection Guidance
- Use fluent mode for readability and lazy pipelines.
- Use rigid mode for stronger stage typing and error signatures.
