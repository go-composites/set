package Set

import (
	Array "github.com/go-composites/array/src"
	Result "github.com/go-composites/result/src"
)

// Interface is the public contract of a Set composite — an unordered collection
// of unique items, the go-composites sibling of Array (with which it shares the
// Each/ToArray grammar). Items are arbitrary comparable values (a Set is backed
// by a Go map, so each item must be comparable, exactly like a Dictionary key).
// Membership tests return a plain bool, fallible iteration returns a Result, and
// every method honours the Null-Object invariant (never nil).
type Interface interface {
	Add(item interface{}) Interface
	Delete(item interface{}) Interface
	Has(item interface{}) bool
	Len() int
	IsEmpty() bool
	Each(fn func(item interface{}) Result.Interface) Result.Interface
	ToArray() Array.Interface
	Union(other Interface) Interface
	Intersection(other Interface) Interface
	Difference(other Interface) Interface
	IsSubset(other Interface) bool
	Equal(other Interface) bool
	IsNull() bool
}

type data struct {
	value map[interface{}]struct{}
}

// New creates a Set seeded with the given items (deduplicated). Items must be
// comparable, since the Set is backed by a Go map.
func New(items ...interface{}) Interface {
	d := &data{
		value: make(map[interface{}]struct{}),
	}
	for _, item := range items {
		d.value[item] = struct{}{}
	}
	return d
}

// Add inserts item (idempotent — a no-op when already present) and returns the
// receiver so calls chain.
func (d *data) Add(item interface{}) Interface {
	d.value[item] = struct{}{}
	return d
}

// Delete removes item (a no-op when absent) and returns the receiver so calls
// chain.
func (d *data) Delete(item interface{}) Interface {
	delete(d.value, item)
	return d
}

// Has reports whether item is a member of the Set.
func (d *data) Has(item interface{}) bool {
	_, ok := d.value[item]
	return ok
}

// Len returns the number of items in the Set.
func (d *data) Len() int {
	return len(d.value)
}

// IsEmpty reports whether the Set has no items.
func (d *data) IsEmpty() bool {
	return len(d.value) == 0
}

// Each iterates over the items, invoking fn for each. It short-circuits and
// returns the first Result for which HasError() is true; on a full pass it
// returns a fresh Result.New(). Iteration order is unspecified (Go map order).
func (d *data) Each(
	fn func(item interface{}) Result.Interface,
) Result.Interface {
	for item := range d.value {
		if result := fn(item); result.HasError() {
			return result
		}
	}
	return Result.New()
}

// ToArray materialises the Set into a go-composites Array. Order is unspecified
// (Go map order).
func (d *data) ToArray() Array.Interface {
	arr := Array.New()
	for item := range d.value {
		arr.Push(item)
	}
	return arr
}

// Union returns a new Set containing every item present in this Set or in other
// (Ruby's `|`).
func (d *data) Union(other Interface) Interface {
	result := New()
	for item := range d.value {
		result.Add(item)
	}
	other.Each(func(item interface{}) Result.Interface {
		result.Add(item)
		return Result.New()
	})
	return result
}

// Intersection returns a new Set containing only the items present in both this
// Set and other (Ruby's `&`).
func (d *data) Intersection(other Interface) Interface {
	result := New()
	for item := range d.value {
		if other.Has(item) {
			result.Add(item)
		}
	}
	return result
}

// Difference returns a new Set containing the items present in this Set but not
// in other (Ruby's `-`).
func (d *data) Difference(other Interface) Interface {
	result := New()
	for item := range d.value {
		if !other.Has(item) {
			result.Add(item)
		}
	}
	return result
}

// IsSubset reports whether every item of this Set is also in other.
func (d *data) IsSubset(other Interface) bool {
	for item := range d.value {
		if !other.Has(item) {
			return false
		}
	}
	return true
}

// Equal reports whether this Set and other contain exactly the same items.
func (d *data) Equal(other Interface) bool {
	return d.Len() == other.Len() && d.IsSubset(other)
}

// IsNull reports that this is a real (non-null) Set.
func (d *data) IsNull() bool {
	return false
}

// null is the Null-Object variant of a Set: an empty, immutable placeholder
// that honours the full Interface without ever being nil. Mutating methods are
// no-ops that return the receiver; queries are empty/false/zero.
type null struct{}

// Null returns the Null-Object Set.
func Null() Interface {
	return &null{}
}

func (n *null) Add(item interface{}) Interface { return n }

func (n *null) Delete(item interface{}) Interface { return n }

func (n *null) Has(item interface{}) bool { return false }

func (n *null) Len() int { return 0 }

func (n *null) IsEmpty() bool { return true }

func (n *null) Each(
	fn func(item interface{}) Result.Interface,
) Result.Interface {
	return Result.New()
}

func (n *null) ToArray() Array.Interface { return Array.New() }

func (n *null) Union(other Interface) Interface { return n }

func (n *null) Intersection(other Interface) Interface { return n }

func (n *null) Difference(other Interface) Interface { return n }

func (n *null) IsSubset(other Interface) bool { return true }

func (n *null) Equal(other Interface) bool { return other.IsEmpty() }

// IsNull reports that this is the null Set.
func (n *null) IsNull() bool { return true }
