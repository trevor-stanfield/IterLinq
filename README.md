# Iterlinq

[![Go Reference](https://pkg.go.dev/badge/github.com/trevorstanfield/iterlinq.svg)](https://pkg.go.dev/github.com/trevor-stanfield/iterlinq)

Iterlinq is a LINQ-style toolkit for Go with two interface modes:

- `iterlinq`: fluent, lazy sequence pipelines
- `iterlinq/rigid`: strict typed slice-transform pipelines

## Documentation Wiki

Deep documentation now lives in the wiki-style docs set:

- [Wiki Home](./docs/wiki/README.md)
- [Quick Start](./docs/wiki/quick-start.md)
- [Architecture](./docs/wiki/architecture.md)
- [Error Model](./docs/wiki/error-model.md)
- [Fluent API](./docs/wiki/fluent-api.md)
- [Rigid API](./docs/wiki/rigid-api.md)
- [Operators Reference](./docs/wiki/operators-reference.md)
- [Testing and Quality](./docs/wiki/testing-and-quality.md)
- [Migration Guide](./docs/wiki/migration-guide.md)
- [Recipes](./docs/wiki/recipes.md)

## Install

```bash
go get github.com/trevor-stanfield/iterlinq
```

## Quick Example (Fluent)

```go
package main

import (
    "fmt"

    "github.com/trevor-stanfield/iterlinq"
)

func main() {
    out, err := iterlinq.FromSlice([]int{1, 2, 3, 4, 5, 6}).
        Where(func(v int) bool { return v%2 == 0 }).
        Select(func(v int) any { return v * 10 }).
        ToSlice()
    if err != nil {
        panic(err)
    }
    fmt.Println(out) // [20 40 60]
}
```

## Quick Example (Rigid)

```go
package main

import (
    "fmt"
    "strconv"

    "github.com/trevor-stanfield/iterlinq/rigid"
)

func main() {
    t := rigid.Compose(
        rigid.Where(func(v int) (bool, error) { return v%2 == 0, nil }),
        rigid.Select(func(v int) (string, error) { return strconv.Itoa(v), nil }),
    )

    out, err := rigid.Run([]int{1, 2, 3, 4}, t)
    if err != nil {
        panic(err)
    }
    fmt.Println(out) // [2 4]
}
```

## Validate

```bash
go test ./...
```

## License

MIT. See [LICENSE](./LICENSE).
