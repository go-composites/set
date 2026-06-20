package main

import (
	"fmt"
	"sort"

	Result "github.com/go-composites/result/src"
	Set "github.com/go-composites/set/src"
)

func main() {
	a := Set.New(1, 2, 3)
	a.Add(3).Add(4) // Add returns the Set, so calls chain; 3 is idempotent.

	fmt.Printf("Len = %d\n", a.Len())
	fmt.Printf("Has(2) = %t\n", a.Has(2))
	fmt.Printf("IsEmpty = %t\n", a.IsEmpty())

	b := Set.New(3, 4, 5)
	fmt.Printf("Union        = %v\n", sortedInts(a.Union(b)))
	fmt.Printf("Intersection = %v\n", sortedInts(a.Intersection(b)))
	fmt.Printf("Difference   = %v\n", sortedInts(a.Difference(b)))
	fmt.Printf("IsSubset     = %t\n", Set.New(1, 2).IsSubset(a))
	fmt.Printf("Equal        = %t\n", a.Equal(Set.New(4, 3, 2, 1)))

	a.Delete(1)
	fmt.Printf("after Delete(1), Len = %d\n", a.Len())
}

// sortedInts collects a Set's int items into a sorted slice for a stable print
// (Set iteration order is unspecified — Go map order).
func sortedInts(s Set.Interface) []int {
	out := []int{}
	s.Each(func(item interface{}) Result.Interface {
		out = append(out, item.(int))
		return Result.New()
	})
	sort.Ints(out)
	return out
}
