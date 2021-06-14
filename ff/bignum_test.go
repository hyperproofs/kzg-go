package ff

import (
	"testing"

	"github.com/alinush/go-mcl"
	gmcl "github.com/alinush/go-mcl"
)

// These are sanity tests, to see if whatever bignum library that is being
// used actually handles dst/arg overlaps well.

func TestInplaceAdd(t *testing.T) {
	aVal := RandomFr()
	bVal := RandomFr()
	aPlusB := new(gmcl.Fr)
	gmcl.FrAdd(aPlusB, aVal, bVal)
	twoA := new(gmcl.Fr)
	gmcl.FrMul(twoA, aVal, &TWO)

	check := func(name string, fn func(a, b *gmcl.Fr) bool) {
		t.Run(name, func(t *testing.T) {
			var a, b gmcl.Fr
			CopyFr(&a, aVal)
			CopyFr(&b, bVal)
			if !fn(&a, &b) {
				t.Error("fail")
			}
		})
	}
	check("dst equals lhs", func(a *gmcl.Fr, b *gmcl.Fr) bool {
		gmcl.FrAdd(a, a, b)
		return a.IsEqual(aPlusB)
	})
	check("dst equals rhs", func(a *gmcl.Fr, b *gmcl.Fr) bool {
		gmcl.FrAdd(b, a, b)
		return b.IsEqual(aPlusB)
	})
	check("dst equals lhs and rhs", func(a *gmcl.Fr, b *gmcl.Fr) bool {
		gmcl.FrAdd(a, a, a)
		return a.IsEqual(twoA)
	})
}

func TestInplaceMul(t *testing.T) {
	aVal := RandomFr()
	bVal := RandomFr()
	aMulB := new(gmcl.Fr)
	gmcl.FrMul(aMulB, aVal, bVal)
	squareA := new(gmcl.Fr)
	gmcl.FrMul(squareA, aVal, aVal)

	check := func(name string, fn func(a, b *gmcl.Fr) bool) {
		t.Run(name, func(t *testing.T) {
			var a, b mcl.Fr
			CopyFr(&a, aVal)
			CopyFr(&b, bVal)
			if !fn(&a, &b) {
				t.Error("fail")
			}
		})
	}
	check("dst equals lhs", func(a *gmcl.Fr, b *gmcl.Fr) bool {
		gmcl.FrMul(a, a, b)
		return a.IsEqual(aMulB)
	})
	check("dst equals rhs", func(a *gmcl.Fr, b *gmcl.Fr) bool {
		gmcl.FrMul(b, a, b)
		return b.IsEqual(aMulB)
	})
	check("dst equals lhs and rhs", func(a *gmcl.Fr, b *gmcl.Fr) bool {
		gmcl.FrMul(a, a, a)
		return a.IsEqual(squareA)
	})
}
