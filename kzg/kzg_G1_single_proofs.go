// Original: https://github.com/ethereum/research/blob/master/kzg_data_availability/kzg_proofs.py

// +build !bignum_pure,!bignum_hol256

package kzg

import (
	"github.com/alinush/go-mcl"
	gmcl "github.com/alinush/go-mcl"
	"github.com/hyperproofs/kzg-go/ff"
	"github.com/hyperproofs/kzg-go/fft"
)

// KZG commitment to polynomial in coefficient form
func (ks *KZG1Settings) CommitToPoly(coeffs []gmcl.Fr) *gmcl.G1 {
	var out gmcl.G1
	gmcl.G1MulVec(&out, ks.PK[:len(coeffs)], coeffs)
	return &out
}

// Compute KZG proof for polynomial in coefficient form at position x
func (ks *KZG1Settings) ComputeProofSingle(poly []gmcl.Fr, x *gmcl.Fr) (*gmcl.G1, *gmcl.Fr) {

	divisor := [2]gmcl.Fr{}
	var tmp gmcl.Fr
	ff.CopyFr(&tmp, x)
	gmcl.FrSub(&divisor[0], &ff.ZERO, &tmp)
	ff.CopyFr(&divisor[1], &ff.ONE)

	// quotientPolynomial := fft.PolyLongDiv(poly, divisor[:])

	quotientPolynomial, remainder := fft.PolyDiv(poly, divisor[:])

	var out gmcl.G1
	gmcl.G1MulVec(&out, ks.PK[:len(quotientPolynomial)], quotientPolynomial)
	return &out, &remainder[0]
}

// // Check a proof for a KZG commitment for an evaluation f(x) = y
// func (ks *KZG1Settings) CheckProofSingle(commitment *gmcl.G1, proof *gmcl.G1, x *gmcl.Fr, y *gmcl.Fr) bool {
// 	// Verify the pairing equation
// 	var xG2 gmcl.G2
// 	// gmcl.G2Mul(&xG2, &ff.GenG2, x)
// 	gmcl.G2Mul(&xG2, &ks.VK[0], x)
// 	var sMinuxX gmcl.G2
// 	gmcl.G2Sub(&sMinuxX, &ks.VK[1], &xG2)
// 	var yG1 gmcl.G1
// 	// gmcl.G1Mul(&yG1, &ff.GenG1, y)
// 	gmcl.G1Mul(&yG1, &ks.PK[0], y)
// 	var commitmentMinusY gmcl.G1
// 	gmcl.G1Sub(&commitmentMinusY, commitment, &yG1)

// 	// This trick may be applied in the BLS-lib specific code:
// 	//
// 	// e([commitment - y], [1]) = e([proof],  [s - x])
// 	//    equivalent to
// 	// e([commitment - y]^(-1), [1]) * e([proof],  [s - x]) = 1_T
// 	//
// 	return ff.PairingsVerify(&commitmentMinusY, &ks.VK[0], proof, &sMinuxX)
// }

// Check a proof for a KZG commitment for an evaluation f(x) = y
func (ks *KZG1Settings) CheckProofSingle(commitment *gmcl.G1, proof *gmcl.G1, x *gmcl.Fr, y *gmcl.Fr) bool {

	var g1Tmp, g2Tmp, P1 mcl.G1
	mcl.G1Mul(&g1Tmp, proof, x)
	mcl.G1Add(&P1, commitment, &g1Tmp)
	mcl.G1Mul(&g2Tmp, &ks.PK[0], y)
	mcl.G1Sub(&P1, &P1, &g2Tmp)

	var e1, e2 mcl.GT
	mcl.Pairing(&e1, &P1, &ks.VK[0])
	mcl.Pairing(&e2, proof, &ks.VK[1])
	return e1.IsEqual(&e2)
}
