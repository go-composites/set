<p align="center"><img src="https://raw.githubusercontent.com/go-composites/brand/main/social/go-composites.png" alt="go-composites/set" width="720"></p>

# set

[![ci](https://github.com/go-composites/set/actions/workflows/ci.yml/badge.svg)](https://github.com/go-composites/set/actions/workflows/ci.yml)

A Set composite for Composition-Oriented Programming — an unordered collection
of unique items, the sibling of [`array`](https://github.com/go-composites/array)
and [`dictionary`](https://github.com/go-composites/dictionary). A `Set` wraps a
`map[interface{}]struct{}` (arbitrary **comparable** items, like a Dictionary
key) behind a small `Interface` that follows the go-composites grammar:

- **Never nil / Null-Object**: every constructor and method returns a real
  object; `Null()` provides an inert variant and `IsNull()` distinguishes it.
- **Result-based errors**: fallible iteration returns a
  [`Result`](https://github.com/go-composites/result) — `Each` short-circuits on
  the first `Result` whose `HasError()` is true. No panics, no bare nils.
- **Composite returns**: `ToArray()` materialises into an
  [`Array`](https://github.com/go-composites/array); set algebra (`Union`,
  `Intersection`, `Difference`) returns fresh `Set`s.

Items must be comparable, since a `Set` is backed by a Go map.

## Install

```sh
go get github.com/go-composites/set@main
```

## Usage

```go
package main

import (
	"fmt"

	Result "github.com/go-composites/result/src"
	Set "github.com/go-composites/set/src"
)

func main() {
	a := Set.New(1, 2, 3)
	a.Add(3).Add(4) // Add returns the Set, so calls chain; 3 is idempotent.

	fmt.Println(a.Len())   // 4
	fmt.Println(a.Has(2))  // true
	fmt.Println(a.IsEmpty()) // false

	b := Set.New(3, 4, 5)
	_ = a.Union(b)        // {1,2,3,4,5}
	_ = a.Intersection(b) // {3,4}
	_ = a.Difference(b)   // {1,2}

	fmt.Println(Set.New(1, 2).IsSubset(a))   // true
	fmt.Println(a.Equal(Set.New(4, 3, 2, 1))) // true

	// Each short-circuits on the first Result whose HasError() is true.
	a.Each(func(item interface{}) Result.Interface {
		fmt.Println(item)
		return Result.New()
	})

	// ToArray materialises into a go-composites Array (order unspecified).
	_ = a.ToArray()

	a.Delete(1)
}
```

### API

| Method | Returns | Notes |
| --- | --- | --- |
| `New(items...)` | `Set.Interface` | variadic, deduplicated |
| `Null()` | `Set.Interface` | inert Null-Object; `IsNull()` is `true` |
| `Add(item)` | `Set.Interface` | idempotent; returns the receiver (chainable) |
| `Delete(item)` | `Set.Interface` | no-op when absent; chainable |
| `Has(item)` | `bool` | membership test |
| `Len()` | `int` | number of items |
| `IsEmpty()` | `bool` | `true` when there are no items |
| `Each(fn)` | `Result.Interface` | iterate; short-circuit on `HasError()` |
| `ToArray()` | `Array.Interface` | materialise into an Array (unspecified order) |
| `Union(other)` | `Set.Interface` | Ruby `\|` — items in either set |
| `Intersection(other)` | `Set.Interface` | Ruby `&` — items in both sets |
| `Difference(other)` | `Set.Interface` | Ruby `-` — items in receiver, not other |
| `IsSubset(other)` | `bool` | every item is also in `other` |
| `Equal(other)` | `bool` | same items |
| `IsNull()` | `bool` | `false` for a real Set |

## License

BSD-3-Clause — see [LICENSE](LICENSE).
