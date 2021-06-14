// +build !bignum_pure,!bignum_hol256,!bignum_kilic,!bignum_hbls

package ff

import (
	gmcl "github.com/alinush/go-mcl"
)

func init() {
	gmcl.InitFromString("bls12-381")
	initGlobals()
	ZERO_G1.Clear()
	initG1G2()
}

func SetFr(dst *gmcl.Fr, v string) {
	if err := dst.SetString(v, 10); err != nil {
		panic(err)
	}
}

func RandomFr() *gmcl.Fr {
	var out gmcl.Fr
	out.Random()
	return &out
}

// FrTo32 serializes a fr number to 32 bytes. Encoded little-endian.
func FrTo32(src *gmcl.Fr) (v [32]byte) {
	b := src.Serialize()
	last := len(b) - 1
	// reverse endianness, Herumi outputs big-endian bytes
	for i := 0; i < 16; i++ {
		b[i], b[last-i] = b[last-i], b[i]
	}
	copy(v[:], b)
	return
}

func CopyFr(dst *gmcl.Fr, v *gmcl.Fr) {
	*dst = *v
}

func AsFr(dst *gmcl.Fr, i uint64) {
	dst.SetInt64(int64(i))
}

func FrStr(b *gmcl.Fr) string {
	if b == nil {
		return "<nil>"
	}
	return b.GetString(10)
}

func IntAsFr(dst *gmcl.Fr, i int64) {
	dst.SetInt64(i)
}

func FromInt64Vec(in []int64) []gmcl.Fr {
	n := len(in)
	dst := make([]gmcl.Fr, n, n)
	for i := 0; i < n; i++ {
		(&dst[i]).SetInt64(in[i])
	}
	return dst
}

func MulVecFr(a, b []gmcl.Fr) []gmcl.Fr {

	n := len(a)
	if n == len(b) && n > 0 {
		result := make([]gmcl.Fr, n, n)
		for i := 0; i < n; i++ {
			gmcl.FrMul(&result[i], &a[i], &b[i])
		}
		return result
	}
	result := make([]gmcl.Fr, 0)
	return result
}
