// Original: https://github.com/ethereum/research/blob/master/kzg_data_availability/kzg_proofs.py

// +build !bignum_pure,!bignum_hol256

package kzg

import (
	gmcl "github.com/alinush/go-mcl"
	"github.com/hyperproofs/kzg-go/ff"
	"github.com/hyperproofs/kzg-go/fft"
)

// KZG commitment to polynomial in coefficient form
func (ks *KZG2Settings) CommitToPoly(coeffs []gmcl.Fr) *gmcl.G2 {
	var out gmcl.G2
	gmcl.G2MulVec(&out, ks.PK[:len(coeffs)], coeffs)
	return &out
}

// Compute KZG proof for polynomial in coefficient form at position x
func (ks *KZG2Settings) ComputeProofSingle(poly []gmcl.Fr, x *gmcl.Fr) (*gmcl.G2, *gmcl.Fr) {

	divisor := [2]gmcl.Fr{}
	var tmp gmcl.Fr
	ff.CopyFr(&tmp, x)
	gmcl.FrSub(&divisor[0], &ff.ZERO, &tmp)
	ff.CopyFr(&divisor[1], &ff.ONE)

	// quotientPolynomial := fft.PolyLongDiv(poly, divisor[:])

	quotientPolynomial, remainder := fft.PolyDiv(poly, divisor[:])

	var out gmcl.G2
	gmcl.G2MulVec(&out, ks.PK[:len(quotientPolynomial)], quotientPolynomial)
	return &out, &remainder[0]
}

// Check a proof for a KZG commitment for an evaluation f(x) = y
func (ks *KZG2Settings) CheckProofSingle(commitment *gmcl.G2, proof *gmcl.G2, x *gmcl.Fr, y *gmcl.Fr) bool {

	var Q1 gmcl.G2
	gmcl.G2Mul(&Q1, &ks.PK[0], y)
	gmcl.G2Sub(&Q1, commitment, &Q1)

	var P2 gmcl.G1
	gmcl.G1Mul(&P2, &ks.VK[0], x)
	gmcl.G1Sub(&P2, &P2, &ks.VK[1])

	var e gmcl.GT
	gmcl.MillerLoopVec(&e, []gmcl.G1{ks.VK[0], P2}, []gmcl.G2{Q1, *proof})
	gmcl.FinalExp(&e, &e)

	return e.IsOne()
}
