# File Index

This repository uses filename prefixes so related logic groups together alphabetically.

## Core
- `doc.go`: package-level docs
- `core_sequence.go`: `Sequence` type and constructors (`From`, `FromFunc`, `FromSlice`, `FromError`)
- `error_defs.go`: shared sentinel errors
- `error_flow.go`: monadic error-flow helpers (`Err`, `HasError`, `MapError`, `Recover`, `OrElse`)
- `internal_must_collect_sequence.go`: internal/experimental helper type

## Operators (lazy sequence transforms)
- `op_where.go`
- `op_select.go`
- `op_select_many.go`
- `op_distinct_by.go`
- `op_take.go`
- `op_skip.go`
- `op_order_types.go`
- `op_ordered_sequence.go`

## Terminals (materialization/quantifiers/elements)
- `term_to_slice.go`
- `term_count.go`
- `term_any.go`
- `term_all.go`
- `term_first.go`
- `term_last.go`

## Test Grouping
- `core_*_test.go`: core type/constructor behavior
- `error_*_test.go`: error-flow behavior
- `op_*_test.go`: operator behavior
- `term_*_test.go`: terminal behavior
- `laws_property_test.go`: property-based tests
- `laws_fuzz_test.go`: fuzz targets
- `laws_purity_test.go`: purity and immutability checks

## Optional Strict Typed Interface
- `rigid/`: optional strictly typed pipeline/transformer API
  - includes ordered pipeline support in `rigid/order_types.go` and `rigid/ordered_pipeline.go`

## Documentation Wiki
- `README.md`: top-level documentation hub
- `docs/wiki/README.md`: wiki home
- `docs/wiki/quick-start.md`: installation and quick examples
- `docs/wiki/architecture.md`: design and execution model
- `docs/wiki/fluent-api.md`: lazy fluent API reference
- `docs/wiki/rigid-api.md`: strict typed API reference
- `docs/wiki/error-model.md`: error propagation and recovery behavior
- `docs/wiki/operators-reference.md`: operator and terminal reference
- `docs/wiki/testing-and-quality.md`: testing strategy and commands
- `docs/wiki/migration-guide.md`: migration notes from previous designs
- `docs/wiki/recipes.md`: practical usage patterns
