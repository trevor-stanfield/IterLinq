# Testing and Quality

## Test Layers

### Unit Tests
Coverage exists across:
- constructors and core behavior,
- operator semantics,
- terminal semantics,
- error flow behavior,
- rigid mode equivalents.

### Property Tests
Property tests validate broad invariants over many generated inputs:
- round-trips,
- filtering equivalence,
- map equivalence,
- reconstruction laws (`Take` + `Skip`),
- distinct first-occurrence rules.

### Fuzz Tests
Fuzz targets probe unusual and edge-case data:
- empty vs non-empty sequences,
- odd byte patterns,
- recovery and fallback paths,
- rigid mode identity and recover behavior.

### Purity Tests
Purity-oriented tests verify:
- no unexpected input mutation,
- deterministic repeated materialization,
- branch independence from shared source pipelines.

## Commands
```bash
go test ./...
```

Fuzz explicitly for a target:
```bash
go test -run=^$ -fuzz=Fuzz -fuzztime=10s ./...
```

## Quality Expectations for New Features
1. Add operator-level unit tests.
2. Add at least one property or fuzz target when behavior is data-sensitive.
3. Add error-path tests for new callbacks/stages.
4. Update docs in this wiki and `FILES.md` when file layout changes.
