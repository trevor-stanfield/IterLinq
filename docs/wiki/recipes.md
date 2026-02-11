# Recipes

## Recipe: Runtime Source with Fallback (Fluent)
```go
result, err := iterlinq.FromFunc(func(yield func(int) bool) error {
    for _, v := range []int{1, 2, 3} {
        if v == 2 {
            return errors.New("source failure")
        }
        if !yield(v) {
            return nil
        }
    }
    return nil
}).
Recover(func(error) iterlinq.Sequence[int] {
    return iterlinq.FromSlice([]int{9, 10})
}).
ToSlice()
```

## Recipe: Deduplicate by Composite Key (Rigid)
```go
t := rigid.DistinctBy(func(u User) (string, error) {
    return fmt.Sprintf("%s:%d", u.Region, u.ID), nil
})
out, err := rigid.Run(users, t)
```

## Recipe: Typed Multi-Stage Pipeline (Rigid)
```go
base := rigid.From([]int{1,2,3,4,5,6})
step1 := rigid.Then(base, rigid.Where(func(v int) (bool, error) { return v%2 == 0, nil }))
step2 := rigid.Then(step1, rigid.Select(func(v int) (string, error) { return strconv.Itoa(v), nil }))
out, err := step2.Execute()
```

## Recipe: Normalize Errors

### Fluent
```go
out, err := seq.MapError(func(err error) error {
    return fmt.Errorf("domain failure: %w", err)
}).ToSlice()
```

### Rigid
```go
out, err := pipeline.MapError(func(err error) error {
    return fmt.Errorf("domain failure: %w", err)
}).Execute()
```

## Recipe: Assert No Match

### Fluent
```go
v, err := seq.First(func(x T) bool { return x == target })
if errors.Is(err, iterlinq.ErrNoMatch) {
    // handle empty match
}
```

### Rigid
```go
v, err := pipeline.First(func(x T) (bool, error) { return x == target, nil })
if errors.Is(err, rigid.ErrNoMatch) {
    // handle empty match
}
```
