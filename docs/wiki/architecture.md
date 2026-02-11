# Architecture

## Mental Model
The project provides two interfaces over the same core philosophy: immutable-style composition with explicit errors.

1. `iterlinq` (fluent mode): lazy and sequence-oriented.
2. `iterlinq/rigid` (strict mode): eager and slice-transform oriented.

## Fluent Mode (`iterlinq`)

### Core Type
`Sequence[T]` carries:
- A lazy iterator function (`seq`) that yields values.
- A monadic error state (`err`).

### Execution
- Operators (`Where`, `Select`, `Take`, etc.) compose and return new `Sequence` values.
- Terminal operations (`ToSlice`, `Count`, `First`, etc.) trigger evaluation.
- If `err` is set, all operators short-circuit.

### Error Channel
Errors can come from:
- Monadic state (`FromError`, `MapError`),
- Source function failures (`FromFunc`),
- Operator runtime checks (`DistinctBy` key comparability).

## Rigid Mode (`iterlinq/rigid`)

### Core Contracts
- `Transformer[In, Out]`: deterministic `[]In -> ([]Out, error)` transform.
- `Pipeline[In, Out]`: source plus composed transformer chain.

### Execution
- `Compose`/`Then` combine typed stages.
- `Execute` materializes output.
- `Pipeline.err` short-circuits when set.

### Why Separate Modes
- Fluent mode prioritizes ergonomic chaining.
- Rigid mode prioritizes compile-time transform typing.

## Cross-Cutting Guarantees
- No package-level mutable state.
- Errors stop work early.
- Input slices are not intentionally mutated by library operators.

## Related Docs
- [Error Model](./error-model.md)
- [Fluent API](./fluent-api.md)
- [Rigid API](./rigid-api.md)
