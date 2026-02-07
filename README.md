# Iterlinq

[![Go Reference](https://pkg.go.dev/badge/github.com/trevorstanfield/iterlinq.svg)](https://pkg.go.dev/github.com/trevorstanfield/iterlinq)

Iterlinq is a fluent, method-centric LINQ (Language Integrated Query) implementation for Go, built specifically for the `iter.Seq2[V, error]` semantics introduced in Go 1.23.

## Ethos

The project is guided by three core principles:

1.  **Fluent API**: Logic should read like a sentence. By wrapping standard iterators in a `Sequence[T]` type, we enable chaining (`Where().Select().ToSlice()`) that feels natural and reduces boilerplate.
2.  **First-Class Errors**: In Go, errors are values. Iterlinq treats errors as first-class citizens in the iterator pipeline. If any stage of the pipeline encounters or produces an error, it propagates immediately and safely to the caller.
3.  **Idiomatic & Lightweight**: We lean into Go 1.23's `range` over functions. Iterlinq doesn't try to hide Go's nature; it enhances it with the power of functional composition.

---

## Quick Start

```go
package main

import (
    "fmt"
    "iterlinq"
)

func main() {
    nums := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

    // Filter evens, square them, and collect into a slice
    res, err := iterlinq.FromSlice(nums).
        Where(func(n int) bool { return n%2 == 0 }).
        Select(func(n int) any { return n * n }).
        ToSlice()

    if err != nil {
        panic(err)
    }
    fmt.Println(res) // [4 16 36 64 100]
}
```

---

## Features & Operators

### Filtering & Partitioning
- `Where(predicate)` / `WhereWithError(predicate)`
- `Take(n)`: First `n` elements.
- `Skip(n)`: Bypass first `n` elements.
- `DistinctBy(keySelector)`: Unique elements based on a comparable key.

### 🏗 Projection & Flattening
- `Select(transform)` / `SelectWithError(transform)`
- `SelectMany(transform)`: Flatten nested sequences.

> **Note on Generics**: Due to Go's constraint that methods cannot introduce new type parameters, `Select` and `SelectMany` return `Sequence[any]`. This preserves the fluent method-chaining experience.

### Aggregation & Quantifiers
- `Count()`: Returns the number of elements.
- `Any(predicate)` / `AnyWithError(predicate)`
- `All(predicate)`
- `ToSlice()`: Materializes the sequence.

### Element Operators
- `First` / `FirstOrDefault`
- `Last` / `LastOrDefault`
- `Single` / `SingleOrDefault`
- *(All available with `WithError` variants for complex logic)*

---

## Error Handling

Iterlinq uses `iter.Seq2[T, error]`. Errors can originate from:
1. The **Source**: e.g., a database iterator encountering a connection drop.
2. The **Pipeline**: e.g., a `SelectWithError` transform failing for specific data.

When an error occurs, the iteration stops immediately. No further elements are processed, and the error is bubbled up to the final materializer (like `ToSlice`) or quantifier.

```go
// Example of early exit on error
result, err := iterlinq.FromSlice(data).
    WhereWithError(func(v int) (bool, error) {
        if v < 0 {
            return false, fmt.Errorf("negative value: %d", v)
        }
        return v % 2 == 0, nil
    }).
    ToSlice()
```

---

Run the tests yourself:
```bash
go test ./...
```

---

## License

MIT License. See [LICENSE](LICENSE) for details.
