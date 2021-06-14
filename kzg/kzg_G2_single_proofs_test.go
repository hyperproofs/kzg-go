// +build !bignum_pure,!bignum_hol256

package kzg

import (
	"testing"

	gmcl "github.com/alinush/go-mcl"
	"github.com/hyperproofs/kzg-go/ff"
)

func TestKZG2Settings_CheckProofSingle(t *testing.T) {
	s1, s2 := GenerateTestingSetup("1927409816240961209460912649124", 16+1)
	ks := NewKZG2Settings(s2, s1[:2])
	for i := 0; i < len(ks.VK); i++ {
		t.Logf("secret g1 %d: %s", i, ff.StrG1(&ks.VK[i]))
	}

	polynomial := testPoly(1, 2, 3, 4, 7, 7, 7, 7, 13, 13, 13, 13, 13, 13, 13, 13)
	for i := 0; i < len(polynomial); i++ {
		t.Logf("poly coeff %d: %s", i, ff.FrStr(&polynomial[i]))
	}

	commitment := ks.CommitToPoly(polynomial)
	t.Log("commitment\n", ff.StrG2(commitment))

	var dst, x gmcl.Fr
	// ff.AsFr(&dst, 17)
	dst.Random()
	ff.CopyFr(&x, &dst)
	proof, y := ks.ComputeProofSingle(polynomial, &dst)
	t.Log("proof\n", ff.StrG2(proof))

	// ff.AsFr(&x, 17)
	var value gmcl.Fr
	gmcl.FrEvaluatePolynomial(
		&value,
		polynomial,
		&x,
	)

	t.Log("value\n", ff.FrStr(&value))
	t.Log("y\n", ff.FrStr(y))

	if value.IsEqual(y) == false {
		t.Fatal("Evaluations did not match!")
	}

	if !ks.CheckProofSingle(commitment, proof, &x, &value) {
		t.Fatal("could not verify proof")
	}
}
