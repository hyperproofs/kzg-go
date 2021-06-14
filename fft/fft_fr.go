package fft

import (
	"fmt"

	gmcl "github.com/alinush/go-mcl"
	"github.com/hyperproofs/kzg-go/ff"
)

func (fs *FFTSettings) simpleFT(vals []gmcl.Fr, valsOffset uint64, valsStride uint64, rootsOfUnity []gmcl.Fr, rootsOfUnityStride uint64, out []gmcl.Fr) {
	l := uint64(len(out))
	var v gmcl.Fr
	var tmp gmcl.Fr
	var last gmcl.Fr
	for i := uint64(0); i < l; i++ {
		jv := &vals[valsOffset]
		r := &rootsOfUnity[0]
		gmcl.FrMul(&v, jv, r)
		ff.CopyFr(&last, &v)

		for j := uint64(1); j < l; j++ {
			jv := &vals[valsOffset+j*valsStride]
			r := &rootsOfUnity[((i*j)%l)*rootsOfUnityStride]
			gmcl.FrMul(&v, jv, r)
			ff.CopyFr(&tmp, &last)
			gmcl.FrAdd(&last, &tmp, &v)
		}
		ff.CopyFr(&out[i], &last)
	}
}

func (fs *FFTSettings) _fft(vals []gmcl.Fr, valsOffset uint64, valsStride uint64, rootsOfUnity []gmcl.Fr, rootsOfUnityStride uint64, out []gmcl.Fr) {
	if len(out) <= 1 { // if the value count is small, run the unoptimized version instead. // TODO tune threshold.
		fs.simpleFT(vals, valsOffset, valsStride, rootsOfUnity, rootsOfUnityStride, out)
		return
	}

	half := uint64(len(out)) >> 1
	// L will be the left half of out
	fs._fft(vals, valsOffset, valsStride<<1, rootsOfUnity, rootsOfUnityStride<<1, out[:half])
	// R will be the right half of out
	fs._fft(vals, valsOffset+valsStride, valsStride<<1, rootsOfUnity, rootsOfUnityStride<<1, out[half:]) // just take even again

	var yTimesRoot gmcl.Fr
	var x, y gmcl.Fr
	for i := uint64(0); i < half; i++ {
		// temporary copies, so that writing to output doesn't conflict with input
		ff.CopyFr(&x, &out[i])
		ff.CopyFr(&y, &out[i+half])
		root := &rootsOfUnity[i*rootsOfUnityStride]
		gmcl.FrMul(&yTimesRoot, &y, root)
		gmcl.FrAdd(&out[i], &x, &yTimesRoot)
		gmcl.FrSub(&out[i+half], &x, &yTimesRoot)
	}
}

func (fs *FFTSettings) FFT(vals []gmcl.Fr, inv bool) ([]gmcl.Fr, error) {
	n := uint64(len(vals))
	if n > fs.MaxWidth {
		return nil, fmt.Errorf("got %d values but only have %d roots of unity", n, fs.MaxWidth)
	}
	n = nextPowOf2(n)
	// We make a copy so we can mutate it during the work.
	valsCopy := make([]gmcl.Fr, n, n)
	for i := 0; i < len(vals); i++ {
		ff.CopyFr(&valsCopy[i], &vals[i])
	}
	for i := uint64(len(vals)); i < n; i++ {
		ff.CopyFr(&valsCopy[i], &ff.ZERO)
	}
	out := make([]gmcl.Fr, n, n)
	if err := fs.InplaceFFT(valsCopy, out, inv); err != nil {
		return nil, err
	}
	return out, nil
}

func (fs *FFTSettings) InplaceFFT(vals []gmcl.Fr, out []gmcl.Fr, inv bool) error {
	n := uint64(len(vals))
	if n > fs.MaxWidth {
		return fmt.Errorf("got %d values but only have %d roots of unity", n, fs.MaxWidth)
	}
	if !ff.IsPowerOfTwo(n) {
		return fmt.Errorf("got %d values but not a power of two", n)
	}
	if inv {
		var invLen gmcl.Fr
		ff.AsFr(&invLen, n)
		gmcl.FrInv(&invLen, &invLen)
		rootz := fs.ReverseRootsOfUnity[:fs.MaxWidth]
		stride := fs.MaxWidth / n

		fs._fft(vals, 0, 1, rootz, stride, out)
		var tmp gmcl.Fr
		for i := 0; i < len(out); i++ {
			gmcl.FrMul(&tmp, &out[i], &invLen)
			ff.CopyFr(&out[i], &tmp) // TODO: depending on Fr implementation, allow to directly write back to an input
		}
		return nil
	} else {
		rootz := fs.ExpandedRootsOfUnity[:fs.MaxWidth]
		stride := fs.MaxWidth / n
		// Regular FFT
		fs._fft(vals, 0, 1, rootz, stride, out)
		return nil
	}
}

// rearrange Fr elements in reverse bit order. Supports 2**31 max element count.
func ReverseBitOrderFr(values []gmcl.Fr) {
	if len(values) > (1 << 31) {
		panic("list too large")
	}
	var tmp gmcl.Fr
	reverseBitOrder(uint32(len(values)), func(i, j uint32) {
		ff.CopyFr(&tmp, &values[i])
		ff.CopyFr(&values[i], &values[j])
		ff.CopyFr(&values[j], &tmp)
	})
}

// rearrange Fr ptr elements in reverse bit order. Supports 2**31 max element count.
func ReverseBitOrderFrPtr(values []*gmcl.Fr) {
	if len(values) > (1 << 31) {
		panic("list too large")
	}
	reverseBitOrder(uint32(len(values)), func(i, j uint32) {
		values[i], values[j] = values[j], values[i]
	})
}
