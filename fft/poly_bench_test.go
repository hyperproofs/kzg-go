package fft

import (
	"fmt"
	"testing"

	gmcl "github.com/alinush/go-mcl"
	"github.com/hyperproofs/kzg-go/ff"
)

func BenchmarkPolyMul(b *testing.B) {

	for scale := uint8(4); scale < 15; scale++ {
		n := uint64(1) << scale
		A := make([]gmcl.Fr, n, n)
		B := make([]gmcl.Fr, n, n)
		for i := uint64(0); i < n; i++ {
			A[i] = *(ff.RandomFr())
			B[i] = *(ff.RandomFr())
		}
		b.Run(fmt.Sprintf("scale_%d", scale), func(t *testing.B) {
			b.ResetTimer()
			_ = PolyMul(A, B)
		})
	}
}

func BenchmarkPolyDiv(b *testing.B) {

	for scale := uint8(10); scale < 15; scale++ {
		n := uint64(1) << scale
		A := make([]gmcl.Fr, n, n)
		B := make([]gmcl.Fr, 2, 2)
		for i := uint64(0); i < n; i++ {
			A[i] = *(ff.RandomFr())
		}
		B[0] = *(ff.RandomFr())
		B[1] = *(ff.RandomFr())
		b.Run(fmt.Sprintf("scale_%d", scale), func(t *testing.B) {
			b.ResetTimer()
			_, _ = PolyDiv(A, B)
		})
	}
}
