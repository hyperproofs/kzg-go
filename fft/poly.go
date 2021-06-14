package fft

import (
	"math/bits"

	gmcl "github.com/alinush/go-mcl"
	"github.com/hyperproofs/kzg-go/ff"
)

// Returns true if polynomial A is a zero polynomial.
func IsPolyZero(a []gmcl.Fr) bool {

	n := len(a)
	if n == 0 {
		panic("IsPolyZero Error")
	}
	var flag bool
	flag = true
	for i := 0; i < n && flag == true; i++ {
		flag = flag && a[i].IsZero()
	}
	return flag
}

// Returns true if polynomial A is a equal to polynomial B.
func IsPolyEqual(a []gmcl.Fr, b []gmcl.Fr) bool {

	p := PolyCondense(a)
	q := PolyCondense(b)

	if len(p) == 0 {
		panic("IsPolyEqual: P is zero.")
	}
	if len(q) == 0 {
		panic("IsPolyEqual: Q is zero.")
	}
	if len(p) != len(q) {
		return false
	}

	var flag bool
	flag = true

	for i := 0; i < len(p) && flag == true; i++ {
		flag = flag && p[i] == q[i]
	}
	return flag
}

//  Removes extraneous zero entries from in vector representation of polynomial.
//  Example - Degree-4 Polynomial: [0, 1, 2, 3, 4, 0, 0, 0, 0] -> [0, 1, 2, 3, 4]
//  Note: Simplest condensed form is a zero polynomial of vector form: [0]
func PolyCondense(a []gmcl.Fr) []gmcl.Fr {
	n := len(a)
	if n == 0 {
		panic("PolyCondense Error")
	}

	i := n
	for i > 1 {
		if a[i-1].IsZero() != true {
			break
		}
		i--
	}
	return a[:i]
}

// Computes the standard polynomial addition, polynomial A + polynomial B, and stores result in polynomial C.
func PolyAdd(a []gmcl.Fr, b []gmcl.Fr) []gmcl.Fr {

	if IsPolyZero(a) {
		return PolyCondense(b)
	}

	if IsPolyZero(b) {
		return PolyCondense(a)
	}

	aLen := len(a)
	bLen := len(b)
	n := ff.Max(aLen, bLen)
	c := make([]gmcl.Fr, n, n)

	for i := 0; i < n; i++ {
		if i < aLen {
			gmcl.FrAdd(&c[i], &c[i], &a[i])
		}
		if i < bLen {
			gmcl.FrAdd(&c[i], &c[i], &b[i])
		}
	}
	c = PolyCondense(c)
	return c
}

// Computes the standard polynomial subtraction, polynomial A - polynomial B, and stores result in polynomial C.
func PolySub(a []gmcl.Fr, b []gmcl.Fr) []gmcl.Fr {

	if IsPolyZero(b) {
		return a
	}

	aLen := len(a)
	bLen := len(b)
	n := ff.Max(aLen, bLen)
	c := make([]gmcl.Fr, n, n)

	for i := 0; i < n; i++ {
		if i < aLen {
			gmcl.FrAdd(&c[i], &c[i], &a[i])
		}
		if i < bLen {
			gmcl.FrSub(&c[i], &c[i], &b[i])
		}
	}
	c = PolyCondense(c)
	return c
}

// Compute a(x) * b(x)
func PolyMul(a []gmcl.Fr, b []gmcl.Fr) []gmcl.Fr {
	if IsPolyZero(a) || IsPolyZero(b) {
		return []gmcl.Fr{ff.ZERO}
	}

	aLen := len(a)
	bLen := len(b)
	if aLen == bLen && aLen == 1 {
		c := make([]gmcl.Fr, 1, 1)
		gmcl.FrMul(&c[0], &a[0], &b[0])
		return c
	}
	n := uint64(2 * ff.Max(aLen, bLen))
	n = nextPowOf2(n)

	var padding []gmcl.Fr

	padding = make([]gmcl.Fr, n-uint64(aLen), n-uint64(aLen))
	a = append(a, padding...)

	padding = make([]gmcl.Fr, n-uint64(bLen), n-uint64(bLen))
	b = append(b, padding...)

	l := uint8(bits.Len64(n)) - 1 // n = 8 => 3 or 4?
	fs := NewFFTSettings(l)

	evalsA, _ := fs.FFT(a, false)
	evalsB, _ := fs.FFT(b, false)

	res, _ := fs.FFT(ff.MulVecFr(evalsA, evalsB), true)
	res = res[:(aLen + bLen - 1)]
	res = PolyCondense(res)
	return res
}

// Invert the divisor, then multiply
func polyFactorDiv(dst *gmcl.Fr, a *gmcl.Fr, b *gmcl.Fr) {
	// TODO: use divmod instead.
	var tmp gmcl.Fr
	gmcl.FrInv(&tmp, b)
	gmcl.FrMul(dst, &tmp, a)
}

// Long polynomial division for two polynomials in coefficient form
// Need to check divide by zero
func PolyLongDiv(A []gmcl.Fr, B []gmcl.Fr) []gmcl.Fr {
	if IsPolyZero(B) == true {
		panic("PolyLongDiv: Cannot divide by zero polynomial.")
	}
	a := make([]gmcl.Fr, len(A), len(A))
	for i := 0; i < len(a); i++ {
		ff.CopyFr(&a[i], &A[i])
	}
	aPos := len(a) - 1
	bPos := len(B) - 1
	diff := aPos - bPos
	out := make([]gmcl.Fr, diff+1, diff+1)
	for diff >= 0 {
		quot := &out[diff]
		polyFactorDiv(quot, &a[aPos], &B[bPos])
		var tmp, tmp2 gmcl.Fr
		for i := bPos; i >= 0; i-- {
			// In steps: a[diff + i] -= b[i] * quot
			// tmp =  b[i] * quot
			gmcl.FrMul(&tmp, quot, &B[i])
			// tmp2 = a[diff + i] - tmp
			gmcl.FrSub(&tmp2, &a[diff+i], &tmp)
			// a[diff + i] = tmp2
			ff.CopyFr(&a[diff+i], &tmp2)
		}
		aPos -= 1
		diff -= 1
	}
	return out
}

// Computes q(x) and r(x) s.t. a(x) = q(x) * b(x) + r(x)
func PolyDiv(A []gmcl.Fr, B []gmcl.Fr) ([]gmcl.Fr, []gmcl.Fr) {
	if IsPolyZero(B) == true {
		panic("PolyDiv: Cannot divide by zero polynomial.")
	}

	if len(B) > len(A) {
		panic("PolyDiv: Deg(B) should be <= Ded(A)")
	}

	a := make([]gmcl.Fr, len(A), len(A))
	for i := 0; i < len(a); i++ {
		ff.CopyFr(&a[i], &A[i])
	}
	aPos := len(a) - 1
	bPos := len(B) - 1
	diff := aPos - bPos
	out := make([]gmcl.Fr, diff+1, diff+1)
	for diff >= 0 {
		quot := &out[diff]
		polyFactorDiv(quot, &a[aPos], &B[bPos])
		var tmp, tmp2 gmcl.Fr
		for i := bPos; i >= 0; i-- {
			// In steps: a[diff + i] -= b[i] * quot
			// tmp =  b[i] * quot
			gmcl.FrMul(&tmp, quot, &B[i])
			// tmp2 = a[diff + i] - tmp
			gmcl.FrSub(&tmp2, &a[diff+i], &tmp)
			// a[diff + i] = tmp2
			ff.CopyFr(&a[diff+i], &tmp2)
		}
		aPos -= 1
		diff -= 1
	}
	out = PolyCondense(out)
	a = PolyCondense(a)
	return out, a
}
