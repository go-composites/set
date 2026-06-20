package Set_test

import (
	"sort"

	Error "github.com/go-composites/error/src"
	Result "github.com/go-composites/result/src"
	Set "github.com/go-composites/set/src"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

// errResult builds a Result that reports HasError() == true, so Each
// short-circuits on it (HasError() is !error.IsNull()).
func errResult() Result.Interface {
	return Result.New(Result.WithError(Error.New("sentinel")))
}

// sortedInts collects a Set's int items into a sorted Go slice so assertions are
// independent of Go's unspecified map-iteration order.
func sortedInts(s Set.Interface) []int {
	out := []int{}
	s.Each(func(item interface{}) Result.Interface {
		out = append(out, item.(int))
		return Result.New()
	})
	sort.Ints(out)
	return out
}

var _ = ginkgo.Describe("Set", func() {
	ginkgo.Describe("New", func() {
		ginkgo.It("returns a non-nil, non-null, empty Set", func() {
			s := Set.New()
			gomega.Expect(s).NotTo(gomega.BeNil())
			gomega.Expect(s.IsNull()).To(gomega.BeFalse())
			gomega.Expect(s.Len()).To(gomega.Equal(0))
			gomega.Expect(s.IsEmpty()).To(gomega.BeTrue())
		})

		ginkgo.It("seeds variadic items, deduplicating", func() {
			s := Set.New(1, 2, 2, 3, 3, 3)
			gomega.Expect(s.Len()).To(gomega.Equal(3))
			gomega.Expect(sortedInts(s)).To(gomega.Equal([]int{1, 2, 3}))
			gomega.Expect(s.IsEmpty()).To(gomega.BeFalse())
		})
	})

	ginkgo.Describe("Add", func() {
		ginkgo.It("inserts items and returns the receiver for chaining", func() {
			s := Set.New()
			ret := s.Add(1).Add(2)
			gomega.Expect(ret).To(gomega.BeIdenticalTo(s))
			gomega.Expect(s.Len()).To(gomega.Equal(2))
		})

		ginkgo.It("is idempotent for a duplicate item", func() {
			s := Set.New().Add(1).Add(1)
			gomega.Expect(s.Len()).To(gomega.Equal(1))
		})
	})

	ginkgo.Describe("Delete", func() {
		ginkgo.It("removes an item and returns the receiver", func() {
			s := Set.New(1, 2)
			ret := s.Delete(1)
			gomega.Expect(ret).To(gomega.BeIdenticalTo(s))
			gomega.Expect(s.Has(1)).To(gomega.BeFalse())
			gomega.Expect(s.Len()).To(gomega.Equal(1))
		})

		ginkgo.It("is a no-op for an absent item", func() {
			s := Set.New(1)
			s.Delete(99)
			gomega.Expect(s.Len()).To(gomega.Equal(1))
		})
	})

	ginkgo.Describe("Has", func() {
		ginkgo.It("is true for a present item", func() {
			s := Set.New(1)
			gomega.Expect(s.Has(1)).To(gomega.BeTrue())
		})

		ginkgo.It("is false for an absent item", func() {
			s := Set.New()
			gomega.Expect(s.Has(1)).To(gomega.BeFalse())
		})
	})

	ginkgo.Describe("Each", func() {
		ginkgo.It("visits every item and returns a clean Result", func() {
			s := Set.New(1, 2, 3)
			count := 0
			res := s.Each(func(item interface{}) Result.Interface {
				count++
				return Result.New()
			})
			gomega.Expect(count).To(gomega.Equal(3))
			gomega.Expect(res).NotTo(gomega.BeNil())
			gomega.Expect(res.HasError()).To(gomega.BeFalse())
		})

		ginkgo.It("short-circuits on the first error Result", func() {
			s := Set.New(1, 2, 3)
			count := 0
			res := s.Each(func(item interface{}) Result.Interface {
				count++
				return errResult()
			})
			gomega.Expect(count).To(gomega.Equal(1))
			gomega.Expect(res.HasError()).To(gomega.BeTrue())
		})
	})

	ginkgo.Describe("ToArray", func() {
		ginkgo.It("materialises the items into an Array", func() {
			s := Set.New(1, 2, 3)
			arr := s.ToArray()
			gomega.Expect(arr).NotTo(gomega.BeNil())
			gomega.Expect(arr.Len()).To(gomega.Equal(3))

			out := []int{}
			arr.Each(func(_ int, item interface{}) Result.Interface {
				out = append(out, item.(int))
				return Result.New()
			})
			sort.Ints(out)
			gomega.Expect(out).To(gomega.Equal([]int{1, 2, 3}))
		})
	})

	ginkgo.Describe("Union", func() {
		ginkgo.It("contains every item from both sets", func() {
			a := Set.New(1, 2, 3)
			b := Set.New(3, 4, 5)
			gomega.Expect(sortedInts(a.Union(b))).To(
				gomega.Equal([]int{1, 2, 3, 4, 5}))
			gomega.Expect(sortedInts(b.Union(a))).To(
				gomega.Equal([]int{1, 2, 3, 4, 5}))
		})
	})

	ginkgo.Describe("Intersection", func() {
		ginkgo.It("contains only the common items", func() {
			a := Set.New(1, 2, 3)
			b := Set.New(3, 4, 5)
			gomega.Expect(sortedInts(a.Intersection(b))).To(
				gomega.Equal([]int{3}))
			gomega.Expect(sortedInts(b.Intersection(a))).To(
				gomega.Equal([]int{3}))
		})
	})

	ginkgo.Describe("Difference", func() {
		ginkgo.It("contains items in the receiver but not the other", func() {
			a := Set.New(1, 2, 3)
			b := Set.New(3, 4, 5)
			gomega.Expect(sortedInts(a.Difference(b))).To(
				gomega.Equal([]int{1, 2}))
			gomega.Expect(sortedInts(b.Difference(a))).To(
				gomega.Equal([]int{4, 5}))
		})
	})

	ginkgo.Describe("IsSubset", func() {
		ginkgo.It("is true when every item is in the other set", func() {
			gomega.Expect(Set.New(1, 2).IsSubset(Set.New(1, 2, 3))).To(
				gomega.BeTrue())
		})

		ginkgo.It("is false when an item is missing from the other set", func() {
			gomega.Expect(Set.New(1, 9).IsSubset(Set.New(1, 2, 3))).To(
				gomega.BeFalse())
		})
	})

	ginkgo.Describe("Equal", func() {
		ginkgo.It("is true for sets with the same items", func() {
			gomega.Expect(Set.New(1, 2, 3).Equal(Set.New(3, 2, 1))).To(
				gomega.BeTrue())
		})

		ginkgo.It("is false for sets of different sizes", func() {
			gomega.Expect(Set.New(1, 2).Equal(Set.New(1, 2, 3))).To(
				gomega.BeFalse())
		})

		ginkgo.It("is false for same-size sets with different items", func() {
			gomega.Expect(Set.New(1, 2, 3).Equal(Set.New(1, 2, 9))).To(
				gomega.BeFalse())
		})
	})

	ginkgo.Describe("Null", func() {
		ginkgo.It("is a Null-Object: IsNull true and inert", func() {
			n := Set.Null()
			gomega.Expect(n).NotTo(gomega.BeNil())
			gomega.Expect(n.IsNull()).To(gomega.BeTrue())
			gomega.Expect(n.Len()).To(gomega.Equal(0))
			gomega.Expect(n.IsEmpty()).To(gomega.BeTrue())

			// Mutators are no-ops that return the receiver.
			gomega.Expect(n.Add(1)).To(gomega.BeIdenticalTo(n))
			gomega.Expect(n.Delete(1)).To(gomega.BeIdenticalTo(n))
			gomega.Expect(n.Len()).To(gomega.Equal(0))

			// Membership always misses.
			gomega.Expect(n.Has(1)).To(gomega.BeFalse())

			// ToArray is empty.
			gomega.Expect(n.ToArray().Len()).To(gomega.Equal(0))

			// Set algebra returns the (inert) null set.
			gomega.Expect(n.Union(Set.New(1)).IsNull()).To(gomega.BeTrue())
			gomega.Expect(n.Intersection(Set.New(1)).IsNull()).To(
				gomega.BeTrue())
			gomega.Expect(n.Difference(Set.New(1)).IsNull()).To(
				gomega.BeTrue())

			// IsSubset of anything is true; Equal holds only for empty sets.
			gomega.Expect(n.IsSubset(Set.New(1))).To(gomega.BeTrue())
			gomega.Expect(n.Equal(Set.New())).To(gomega.BeTrue())
			gomega.Expect(n.Equal(Set.New(1))).To(gomega.BeFalse())

			// Each returns a clean Result without invoking fn.
			called := false
			res := n.Each(func(item interface{}) Result.Interface {
				called = true
				return errResult()
			})
			gomega.Expect(called).To(gomega.BeFalse())
			gomega.Expect(res.HasError()).To(gomega.BeFalse())
		})
	})
})
