package main

import (
	"errors"
	"fmt"
	"iterlinq"
)

func main() {
	// 1) Basic filter + map + materialize
	nums := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	evens := iterlinq.FromSlice(nums).
		Where(func(n int) bool { return n%2 == 0 })
	squaredEvens := evens.Select(func(n int) any { return n * n })
	res, err := squaredEvens.ToSlice()
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}
	fmt.Printf("1) Squared evens: %v\n", res)

	// 2) Quantifiers and element operators
	firstEven, err := iterlinq.FromSlice(nums).First(func(n int) bool { return n%2 == 0 })
	fmt.Printf("2) First even: %v (err=%v)\n", firstEven, err)
	lastOver5, err := iterlinq.FromSlice(nums).Last(func(n int) bool { return n > 5 })
	fmt.Printf("   Last > 5: %v (err=%v)\n", lastOver5, err)
	anyDivBy7, err := iterlinq.FromSlice(nums).Any(func(n int) bool { return n%7 == 0 })
	fmt.Printf("   Any divisible by 7? %v (err=%v)\n", anyDivBy7, err)
	allPositive, err := iterlinq.FromSlice(nums).All(func(n int) (bool, error) { return n > 0, nil })
	fmt.Printf("   All positive? %v (err=%v)\n", allPositive, err)
	cnt, err := iterlinq.FromSlice(nums).Count()
	fmt.Printf("   Count: %d (err=%v)\n", cnt, err)

	// 3) DistinctBy on structs
	type user struct {
		id   int
		name string
	}
	users := []user{{1, "ann"}, {2, "bob"}, {1, "ann-dupe"}, {3, "cara"}}
	distinctUsers := iterlinq.FromSlice(users).DistinctBy(func(u user) any { return u.id })
	du, _ := distinctUsers.Select(func(v user) any { return v }).ToSlice()
	fmt.Printf("3) Distinct users by id: %v\n", du)

	// 4) SelectMany to flatten
	words := []string{"go", "linq"}
	chars := iterlinq.FromSlice(words).SelectMany(func(w string) iterlinq.Sequence[any] {
		return iterlinq.FromSlice([]rune(w)).Select(func(r rune) any { return string(r) })
	})
	ch, _ := chars.ToSlice()
	fmt.Printf("4) Chars: %v\n", ch)

	// 5) Error-aware operators demonstration
	// Stop on encountering a 5 using WhereWithError
	stopped, werr := iterlinq.FromSlice(nums).
		WhereWithError(func(n int) (bool, error) {
			if n == 5 {
				return false, errors.New("stop at 5")
			}
			return n%2 == 1, nil
		}).ToSlice()
	fmt.Printf("5) WhereWithError stopped: vals=%v err=%v\n", stopped, werr)

	// 6) Safe From with nil iterator (edge case): treated as empty
	var nilSeq iterlinq.Sequence[int] = iterlinq.From[int](nil)
	nilVals, nerr := nilSeq.ToSlice()
	fmt.Printf("6) From(nil) -> slice=%v err=%v\n", nilVals, nerr)

	// 7) Skip/Take including edge cases
	seg, _ := iterlinq.FromSlice(nums).Skip(3).Take(4).ToSlice()
	fmt.Printf("7) Skip 3 then Take 4: %v\n", seg)
	negTake, _ := iterlinq.FromSlice(nums).Take(-1).ToSlice()
	fmt.Printf("   Take(-1): %v\n", negTake)
	negSkip, _ := iterlinq.FromSlice(nums).Skip(-2).ToSlice()
	fmt.Printf("   Skip(-2): %v\n", negSkip)
}
