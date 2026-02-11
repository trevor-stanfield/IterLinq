# Iterlinq Wiki

This wiki is the source of truth for architecture, API semantics, error behavior, and usage patterns.

## Start Here
- [Quick Start](./quick-start.md)
- [Architecture](./architecture.md)
- [Error Model](./error-model.md)
- [Fluent API (lazy)](./fluent-api.md)
- [Rigid API (strict typed)](./rigid-api.md)
- [Operators Reference](./operators-reference.md)
- [Testing and Quality](./testing-and-quality.md)
- [Migration Guide](./migration-guide.md)
- [Recipes](./recipes.md)

## Choose Your Interface

| Mode | Package | Execution Style | Type Rigidity | Best For |
|---|---|---|---|---|
| Fluent | `iterlinq` | Lazy sequence pipeline | Medium (`Select` returns `Sequence[any]`) | Runtime pipelines, stream-like composition |
| Rigid | `iterlinq/rigid` | Eager slice pipeline | High (full generic typing across stages) | Compile-time strictness and explicit transform contracts |

## Current Capability Snapshot

### Fluent (`iterlinq`)
- Constructors: `From`, `FromFunc`, `FromSlice`, `FromError`
- Error flow: `Err`, `HasError`, `MapError`, `Recover`, `OrElse`
- Operators: `Where`, `Select`, `SelectMany`, `DistinctBy`, `Take`, `Skip`
- Terminals: `ToSlice`, `Count`, `Any`, `All`, `First`, `FirstOrDefault`, `Last`, `LastOrDefault`

### Rigid (`iterlinq/rigid`)
- Core contracts: `Transformer`, `TransformerFunc`, `Compose`, `Run`
- Pipeline: `From`, `Then`, `Execute`, `ToSlice`
- Error flow: `Err`, `HasError`, `MapError`, `Recover`, `OrElse`
- Operators: `Where`, `Select`, `SelectMany`, `DistinctBy`, `Take`, `Skip`
- Terminals: `Count`, `Any`, `All`, `First`, `FirstOrDefault`, `Last`, `LastOrDefault`

## Conventions
- Errors are first-class and short-circuit by design.
- APIs are immutable-style: each operator returns a new value.
- Tests include unit, property, fuzz, and purity checks.
