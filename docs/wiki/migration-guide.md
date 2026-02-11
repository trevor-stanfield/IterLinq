# Migration Guide

## From Earlier `Seq2`-Style Design

If you previously modeled errors as streamed values, use these replacements:

- Old: source emits `(value, error)` pairs
- New: source yields values and returns error from `FromFunc`

### Before (conceptual)
```go
// old style
From(seq2Source)
```

### After
```go
FromFunc(func(yield func(T) bool) error {
    // yield values
    // return error if source fails
})
```

## From Callback `WithError` Variants
The library now favors monadic error flow and explicit source/pipeline error handling.

- Use `FromFunc` to inject runtime source failures.
- Use `MapError`, `Recover`, or `OrElse` to handle failures in-chain.

## Moving to Strict Typed Mode
When `Sequence[any]` from fluent `Select` is too loose:

1. Move transform-heavy code to `iterlinq/rigid`.
2. Replace fluent chain with `Compose`/`Then` typed stages.
3. Keep fluent mode for areas where readability and laziness are more important than static stage typing.

## Common Refactor Patterns

### Pattern: fluent to rigid for type safety
```go
// fluent
iterlinq.FromSlice(xs).Where(...).Select(...).ToSlice()

// rigid
rigid.Run(xs, rigid.Compose(rigid.Where(...), rigid.Select(...)))
```

### Pattern: runtime fallback
```go
// fluent
seq.Recover(func(err error) iterlinq.Sequence[T] { ... })

// rigid
pipeline.Recover(func(err error) rigid.Pipeline[A, B] { ... })
```

## Validation Checklist
- Run `go test ./...`
- Ensure error assertions use `errors.Is` where applicable.
- Confirm docs and examples reference current module path.
