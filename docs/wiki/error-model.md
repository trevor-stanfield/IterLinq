# Error Model

## Design Goal
Errors are first-class and deterministic. A failed stage should stop downstream work immediately.

## Fluent Mode (`iterlinq`)

### Error Sources
1. Attached monadic state (`FromError`, `MapError`).
2. Runtime source failure (`FromFunc`).
3. Operator failure (for example, non-comparable key in `DistinctBy`).

### Propagation Rules
1. If `Sequence.err != nil`, all operators return a new errored sequence without invoking callbacks.
2. Terminal methods return that error immediately.
3. Runtime errors during enumeration terminate evaluation and bubble up.

### Recovery
- `MapError` rewrites error values.
- `Recover` can switch to a replacement sequence.
- `OrElse` is fallback shorthand.

## Rigid Mode (`iterlinq/rigid`)

### Error Sources
1. Pipeline monadic state (`Pipeline.err`).
2. Transformer callback failures.
3. Composition/guard errors (nil transformer or nil handlers).

### Propagation Rules
1. A pipeline in error state short-circuits `Then` and `Execute`.
2. Transformer failures stop execution and return immediately.
3. Terminal methods (`Any`, `First`, etc.) pass callback errors through unchanged.

### Recovery
- `MapError` maps attached/runtime errors.
- `Recover` replaces with another pipeline.
- `OrElse` provides fallback pipeline.

## Sentinel Errors

### Fluent package
- `ErrNoMatch`
- `ErrNonComparableKey`
- `ErrNilErrorMapper`
- `ErrNilErrorHandler`

### Rigid package
- `ErrNoMatch`
- `ErrNilTransformer`
- `ErrNilPredicate`
- `ErrNilSelector`
- `ErrNilErrorMapper`
- `ErrNilErrorHandler`

## Guidance
- Prefer `errors.Is` for sentinel checks.
- Keep recover handlers side-effect free when possible.
- Use `MapError` to normalize external-system errors into domain-specific errors.
