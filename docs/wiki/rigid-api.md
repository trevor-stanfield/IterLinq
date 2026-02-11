# Rigid API (`iterlinq/rigid`)

This mode emphasizes compile-time type rigidity for transform chains.

## Core Types

### `Transformer[In, Out]`
```go
type Transformer[In, Out any] interface {
    Transform([]In) ([]Out, error)
}
```

### `Pipeline[In, Out]`
Stores source + composed transform + monadic error state.

## Construction and Composition

### `From(source []T) Pipeline[T, T]`
Creates a pipeline with an identity transform.

### `Compose(ab Transformer[A,B], bc Transformer[B,C]) Transformer[A,C]`
Composes typed transforms directly.

### `Then(p Pipeline[A,B], next Transformer[B,C]) Pipeline[A,C]`
Appends a stage to a pipeline.

### `Run(in []In, t Transformer[In, Out]) ([]Out, error)`
Executes a transformer without creating a pipeline.

## Operators
All operator callbacks are error-aware and typed.

- `Where(func(T) (bool, error)) Transformer[T, T]`
- `Select(func(T) (R, error)) Transformer[T, R]`
- `SelectMany(func(T) ([]R, error)) Transformer[T, R]`
- `DistinctBy(func(T) (K, error)) Transformer[T, T]` (`K comparable`)
- `Take(n int) Transformer[T, T]`
- `Skip(n int) Transformer[T, T]`

## Ordered Pipeline Variant (Semi-Deferred)

### Start and Build
- `Pipeline[In, Out].MustOrderBy(func(Out, Out) int) OrderedPipeline[In, Out]`
- `Pipeline[In, Out].MustOrderByDesc(func(Out, Out) int) OrderedPipeline[In, Out]`
- `OrderedPipeline[In, Out].ThenBy(func(Out, Out) int) OrderedPipeline[In, Out]`
- `OrderedPipeline[In, Out].ThenByDesc(func(Out, Out) int) OrderedPipeline[In, Out]`
- `OrderedPipeline[In, Out].Reverse() OrderedPipeline[In, Out]`
- `OrderedPipeline[In, Out].UsingSorter(func([]Out, func(Out, Out) int) error) OrderedPipeline[In, Out]`
- `OrderedPipeline[In, Out].UsingSortFunc(func([]Out, func(Out, Out) int) error) OrderedPipeline[In, Out]` (compatibility alias)
- `OrderedPipeline[In, Out].ForceOrder() Pipeline[In, Out]`

### Force-Then-Transform Helpers
- `ThenOrdered(OrderedPipeline[A, B], Transformer[B, C]) Pipeline[A, C]`
- `SelectOrdered(OrderedPipeline[A, B], func(B) (C, error)) Pipeline[A, C]`
- `SelectManyOrdered(OrderedPipeline[A, B], func(B) ([]C, error)) Pipeline[A, C]`
- `DistinctByOrdered(OrderedPipeline[A, B], func(B) (K, error)) Pipeline[A, B]`

`OrderedPipeline` keeps order conditions deferred until forced. Non-order methods on `OrderedPipeline` (like `Where`, `Take`, terminals) force ordering before proceeding.

## Pipeline Terminals
- `Execute() ([]Out, error)`
- `ToSlice() ([]Out, error)`
- `Count() (int, error)`
- `Any(func(Out) (bool, error)) (bool, error)`
- `All(func(Out) (bool, error)) (bool, error)`
- `First(func(Out) (bool, error)) (Out, error)`
- `FirstOrDefault(func(Out) (bool, error)) (Out, error)`
- `Last(func(Out) (bool, error)) (Out, error)`
- `LastOrDefault(func(Out) (bool, error)) (Out, error)`

## Pipeline Error Flow
- `Err() error`
- `HasError() bool`
- `MapError(func(error) error) Pipeline[In, Out]`
- `Recover(func(error) Pipeline[In, Out]) Pipeline[In, Out]`
- `OrElse(fallback Pipeline[In, Out]) Pipeline[In, Out]`

## Example
```go
base := rigid.From([]int{1,2,3,4,5,6})
step1 := rigid.Then(base, rigid.Where(func(v int) (bool, error) { return v > 2, nil }))
step2 := rigid.Then(step1, rigid.Select(func(v int) (string, error) { return strconv.Itoa(v), nil }))
out, err := step2.Execute()
```

## When to Use
Use `rigid` when:
- you need stage-by-stage typed outputs,
- you want callback-level error-aware signatures,
- you prefer explicit slice transforms over lazy sequence semantics.

Use fluent `iterlinq` when pipeline readability and lazy sequence composition are the priority.
