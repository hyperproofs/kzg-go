package fft

import (
	"fmt"
	"testing"

	gmcl "github.com/alinush/go-mcl"
	"github.com/hyperproofs/kzg-go/ff"
)

func benchFFT(scale uint8, inv bool, b *testing.B) {
	fs := NewFFTSettings(scale)
	data := make([]gmcl.Fr, fs.MaxWidth, fs.MaxWidth)
	for i := uint64(0); i < fs.MaxWidth; i++ {
		ff.CopyFr(&data[i], ff.RandomFr())
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		out, err := fs.FFT(data, inv)
		b.StopTimer()
		if err != nil {
			b.Fatal(err)
		}
		if len(out) != len(data) {
			panic("output len doesn't match input")
		}
		b.StartTimer()
	}
}

func BenchmarkFFTSettings_FFT(b *testing.B) {
	for scale := uint8(4); scale < 17; scale++ {
		b.Run(fmt.Sprintf("scale_%d", scale), func(b *testing.B) {
			benchFFT(scale, false, b)
		})
	}
}

func BenchmarkFFTSettings_InvFFT(b *testing.B) {
	for scale := uint8(4); scale < 17; scale++ {
		b.Run(fmt.Sprintf("scale_%d", scale), func(b *testing.B) {
			benchFFT(scale, true, b)
		})
	}
}
