# Quick Start

## Requirements
- Go `1.25+`
- Module path: `github.com/trevor-stanfield/iterlinq`

## Install
```bash
go get github.com/trevor-stanfield/iterlinq
```

## Fluent Mode (Lazy)
```go
package main

import (
    "fmt"

    "github.com/trevor-stanfield/iterlinq"
)

func main() {
    nums := []int{1, 2, 3, 4, 5, 6}

    out, err := iterlinq.FromSlice(nums).
        Where(func(v int) bool { return v%2 == 0 }).
        Select(func(v int) any { return v * v }).
        ToSlice()
    if err != nil {
        panic(err)
    }

    fmt.Println(out) // [4 16 36]
}
```

## Rigid Mode (Strict Typed)
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
        rigid.Select(func(v int) (string, error) { return strconv.Itoa(v * 10), nil }),
    )

    out, err := rigid.Run([]int{1, 2, 3, 4}, t)
    if err != nil {
        panic(err)
    }

    fmt.Println(out) // [20 40]
}
```

## Validate Locally
```bash
go test ./...
```

## Next
- [Architecture](./architecture.md)
- [Error Model](./error-model.md)
- [Fluent API](./fluent-api.md)
- [Rigid API](./rigid-api.md)
