// +build !bignum_pure,!bignum_hol256,!bignum_kilic,!bignum_hbls

package ff

import (
	"fmt"
	"strings"

	gmcl "github.com/alinush/go-mcl"
)

var ZERO_G1 gmcl.G1

var GenG1 gmcl.G1
var GenG2 gmcl.G2

var ZeroG1 gmcl.G1
var ZeroG2 gmcl.G2

// Herumi BLS doesn't offer these points to us, so we have to work around it by declaring them ourselves.
func initG1G2() {
	GenG1.X.SetString("3685416753713387016781088315183077757961620795782546409894578378688607592378376318836054947676345821548104185464507", 10)
	GenG1.Y.SetString("1339506544944476473020471379941921221584933875938349620426543736416511423956333506472724655353366534992391756441569", 10)
	GenG1.Z.SetInt64(1)

	GenG2.X.D[0].SetString("352701069587466618187139116011060144890029952792775240219908644239793785735715026873347600343865175952761926303160", 10)
	GenG2.X.D[1].SetString("3059144344244213709971259814753781636986470325476647558659373206291635324768958432433509563104347017837885763365758", 10)
	GenG2.Y.D[0].SetString("1985150602287291935568054521177171638300868978215655730859378665066344726373823718423869104263333984641494340347905", 10)
	GenG2.Y.D[1].SetString("927553665492332455747201965776037880757740193453592970025027978793976877002675564980949289727957565575433344219582", 10)
	GenG2.Z.D[0].SetInt64(1)
	GenG2.Z.D[1].Clear()

	ZeroG1.X.SetInt64(1)
	ZeroG1.Y.SetInt64(1)
	ZeroG1.Z.SetInt64(0)

	ZeroG2.X.D[0].SetInt64(1)
	ZeroG2.X.D[1].SetInt64(0)
	ZeroG2.Y.D[0].SetInt64(1)
	ZeroG2.Y.D[1].SetInt64(0)
	ZeroG2.Z.D[0].SetInt64(0)
	ZeroG2.Z.D[1].SetInt64(0)
}
func CopyG1(dst *gmcl.G1, v *gmcl.G1) {
	*dst = *v
}

func StrG1(v *gmcl.G1) string {
	return (v).GetString(10)
}

func StrG2(v *gmcl.G2) string {
	return (v).GetString(10)
}

// // e(a1^(-1), a2) * e(b1,  b2) = 1_T
// func PairingsVerify(a1 *gmcl.G1, a2 *gmcl.G2, b1 *gmcl.G1, b2 *gmcl.G2) bool {
// 	var tmp gmcl.GT
// 	gmcl.Pairing(&tmp, (*gmcl.G1)(a1), (*gmcl.G2)(a2))
// 	//fmt.Println("tmp", tmp.GetString(10))
// 	var tmp2 gmcl.GT
// 	gmcl.Pairing(&tmp2, (*gmcl.G1)(b1), (*gmcl.G2)(b2))

// 	// invert left pairing
// 	var tmp3 gmcl.GT
// 	gmcl.GTInv(&tmp3, &tmp)

// 	// multiply the two
// 	var tmp4 gmcl.GT
// 	gmcl.GTMul(&tmp4, &tmp3, &tmp2)

// 	// final exp.
// 	var tmp5 gmcl.GT
// 	gmcl.FinalExp(&tmp5, &tmp4)

// 	// = 1_T
// 	return tmp5.IsOne()

// 	// TODO, alternatively use the equal check (faster or slower?):
// 	////fmt.Println("tmp2", tmp2.GetString(10))
// 	//return tmp.IsEqual(&tmp2)
// }

// e(a1^(-1), a2) * e(b1,  b2) = 1_T
func PairingsVerify(a1 *gmcl.G1, a2 *gmcl.G2, b1 *gmcl.G1, b2 *gmcl.G2) bool {

	var a1Inv gmcl.G1
	a1Inv.Clear()
	gmcl.G1Sub(&a1Inv, &a1Inv, a1)

	P := []gmcl.G1{a1Inv, *b1}
	Q := []gmcl.G2{*a2, *b2}

	var tmp gmcl.GT

	gmcl.MillerLoopVec(&tmp, P, Q)
	gmcl.FinalExp(&tmp, &tmp)

	return tmp.IsOne()

	// TODO, alternatively use the equal check (faster or slower?):
	////fmt.Println("tmp2", tmp2.GetString(10))
	//return tmp.IsEqual(&tmp2)
}

func DebugG1s(msg string, values []gmcl.G1) {
	var out strings.Builder
	for i := range values {
		out.WriteString(fmt.Sprintf("%s %d: %s\n", msg, i, StrG1(&values[i])))
	}
	fmt.Println(out.String())
}
