// Original: https://github.com/ethereum/research/blob/master/mimc_stark/fft.py

package fft

import (
	"math/bits"

	gmcl "github.com/alinush/go-mcl"
	"github.com/hyperproofs/kzg-go/ff"
)

// if not already a power of 2, return the next power of 2
func nextPowOf2(v uint64) uint64 {
	if v == 0 {
		return 1
	}
	return uint64(1) << bits.Len64(v-1)
}

// Expands the power circle for a given root of unity to WIDTH+1 values.
// The first entry will be 1, the last entry will also be 1,
// for convenience when reversing the array (useful for inverses)
func expandRootOfUnity(RootOfUnity *gmcl.Fr) []gmcl.Fr {
	rootz := make([]gmcl.Fr, 2)
	rootz[0] = ff.ONE // some unused number in py code
	rootz[1] = *RootOfUnity
	for i := 1; !rootz[i].IsOne(); {
		rootz = append(rootz, gmcl.Fr{})
		this := &rootz[i]
		i++
		gmcl.FrMul(&rootz[i], this, RootOfUnity)
	}
	return rootz
}

type FFTSettings struct {
	MaxWidth uint64
	// the generator used to get all roots of unity
	RootOfUnity *gmcl.Fr
	// domain, starting and ending with 1 (duplicate!)
	ExpandedRootsOfUnity []gmcl.Fr
	// reverse domain, same as inverse values of domain. Also starting and ending with 1.
	ReverseRootsOfUnity []gmcl.Fr
}

func NewFFTSettings(maxScale uint8) *FFTSettings {
	width := uint64(1) << maxScale
	root := &ff.Scale2RootOfUnity[maxScale]
	rootz := expandRootOfUnity(&ff.Scale2RootOfUnity[maxScale])
	// reverse roots of unity
	rootzReverse := make([]gmcl.Fr, len(rootz), len(rootz))
	copy(rootzReverse, rootz)
	for i, j := uint64(0), uint64(len(rootz)-1); i < j; i, j = i+1, j-1 {
		rootzReverse[i], rootzReverse[j] = rootzReverse[j], rootzReverse[i]
	}

	return &FFTSettings{
		MaxWidth:             width,
		RootOfUnity:          root,
		ExpandedRootsOfUnity: rootz,
		ReverseRootsOfUnity:  rootzReverse,
	}
}
