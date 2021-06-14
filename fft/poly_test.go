package fft

import (
	"fmt"
	"testing"

	gmcl "github.com/alinush/go-mcl"
	"github.com/hyperproofs/kzg-go/debug"
	"github.com/hyperproofs/kzg-go/ff"
)

func CheckEqualVec(a []gmcl.Fr, b []gmcl.Fr) bool {
	n := len(a)
	if n == len(b) && n > 0 {
		flag := true
		for i := 0; i < n; i++ {
			flag = flag && a[i].IsEqual(&b[i])
		}
		return flag
	}
	return false
}

func printSubTree(M [][][]gmcl.Fr, printMsg string) {
	for i := 0; i < len(M); i++ {
		for j := 0; j < len(M[i]); j++ {
			msg := fmt.Sprintf("%s [%d, %d]", printMsg, i, j)
			debug.DebugFrs(msg, M[i][j])
		}
	}
}

func TestPolyIsPolyZero(t *testing.T) {
	var tests = []struct {
		a    []int64
		want bool
	}{
		{[]int64{1, 2, 3, 4, 7, 8}, false},
		{[]int64{0}, true},
		{[]int64{1, 2, 3, 4, 7, 8, 0, 0}, false},
		{[]int64{0, 0, 0, 0, 0}, true},
	}

	for counter, tt := range tests {
		testname := fmt.Sprintf("%d", counter+1)
		t.Run(testname, func(t *testing.T) {

			aFr := ff.FromInt64Vec(tt.a)
			ans := IsPolyZero(aFr)

			if ans != tt.want {
				t.Errorf("IsPolyZero: Answer did not match with expected.")
			}
		})
	}
}

func TestPolyCondense(t *testing.T) {
	var tests = []struct {
		a, want []int64
	}{
		{[]int64{1, 2, 3, 4, 7, 8}, []int64{1, 2, 3, 4, 7, 8}},
		{[]int64{0}, []int64{0}},
		{[]int64{0, 0, 0, 0, 0}, []int64{0}},
		{[]int64{1, 2, 3, 4, 7, 8, 0, 0, 0, 0, 0}, []int64{1, 2, 3, 4, 7, 8}},
		{[]int64{1, 2, 3, 4, 7, -8, 0, 0, 0, 0, 0}, []int64{1, 2, 3, 4, 7, -8}},
	}

	for counter, tt := range tests {
		testname := fmt.Sprintf("%d", counter+1)
		t.Run(testname, func(t *testing.T) {

			aFr := ff.FromInt64Vec(tt.a)
			wantFr := ff.FromInt64Vec(tt.want)
			ansFr := PolyCondense(aFr)

			flag := CheckEqualVec(wantFr, ansFr)
			if flag == false {
				t.Errorf("PolyCondense: Answer did not match with expected.")
			}
		})
	}
}

func TestPolyAdd(t *testing.T) {
	var tests = []struct {
		a, b, want []int64
	}{
		{
			[]int64{1, 2, 3, 4, 7, 8},
			[]int64{5, 1, 3},
			[]int64{6, 3, 6, 4, 7, 8},
		},
		{
			[]int64{1, 2, 3, 4, 7, 8}, []int64{0}, []int64{1, 2, 3, 4, 7, 8},
		},
		{
			[]int64{1, 2, 3, 4, 7, 8},
			[]int64{-1, -2, -3, -4, -7, -8},
			[]int64{0},
		},
		{
			[]int64{1, 2, 3, 4, 7, 8},
			[]int64{1, 2, 3, 4, 7, -8},
			[]int64{2, 4, 6, 8, 14},
		},
		{
			[]int64{1, 2, 3, 4, 7, 8},
			[]int64{0, 0, 0, 0, 0, -8},
			[]int64{1, 2, 3, 4, 7},
		},
		{
			[]int64{0, 0, 0, 0},
			[]int64{1, 2, 3, 4, 7, -8},
			[]int64{1, 2, 3, 4, 7, -8},
		},
		{
			[]int64{1, 2, 3, 4, 7, -8},
			[]int64{0, 0, 0, 0},
			[]int64{1, 2, 3, 4, 7, -8},
		},
	}

	for counter, tt := range tests {
		testname := fmt.Sprintf("%d", counter+1)
		t.Run(testname, func(t *testing.T) {

			aFr := ff.FromInt64Vec(tt.a)
			bFr := ff.FromInt64Vec(tt.b)
			wantFr := ff.FromInt64Vec(tt.want)
			ansFr := PolyAdd(aFr, bFr)

			// debug.DebugFrs("", wantFr)
			// debug.DebugFrs("", ansFr)

			flag := CheckEqualVec(wantFr, ansFr)
			if flag == false {
				t.Errorf("PolyAdd: Answer did not match with expected.")
			}
		})
	}
}

func TestPolySub(t *testing.T) {
	var tests = []struct {
		a, b, want []int64
	}{
		{
			[]int64{0, 0, 0, 0},
			[]int64{1, 2, 3, 4, 7, -8},
			[]int64{-1, -2, -3, -4, -7, 8},
		},
		{
			[]int64{1, 2, 3, 4, 7, -8},
			[]int64{0, 0, 0, 0},
			[]int64{1, 2, 3, 4, 7, -8},
		},
		{
			[]int64{1, 2, 3, 4, 7, 8},
			[]int64{5, 1, 3},
			[]int64{-4, 1, 0, 4, 7, 8},
		},
		{
			[]int64{1, 2, 3, 4, 7, 8}, []int64{0}, []int64{1, 2, 3, 4, 7, 8},
		},
		{
			[]int64{1, 2, 3, 4, 7, 8},
			[]int64{1, 2, 3, 4, 7, 8},
			[]int64{0},
		},
	}

	for counter, tt := range tests {
		testname := fmt.Sprintf("%d", counter+1)
		t.Run(testname, func(t *testing.T) {

			aFr := ff.FromInt64Vec(tt.a)
			bFr := ff.FromInt64Vec(tt.b)
			wantFr := ff.FromInt64Vec(tt.want)
			ansFr := PolySub(aFr, bFr)

			// debug.DebugFrs("", wantFr)
			// debug.DebugFrs("", ansFr)

			flag := CheckEqualVec(wantFr, ansFr)
			if flag == false {
				t.Errorf("PolySub: Answer did not match with expected.")
			}
		})
	}
}

func TestPolyMul(t *testing.T) {
	var tests = []struct {
		a, b, want []int64
	}{
		{[]int64{1, 2, 3, 4, 7, 8}, []int64{5, 1, 3}, []int64{5, 11, 20, 29, 48, 59, 29, 24}},
		{[]int64{5, 1, 3}, []int64{1, 2, 3, 4, 7}, []int64{5, 11, 20, 29, 48, 19, 21}},
		{[]int64{1, 2, 3, 4, 7, 8}, []int64{0}, []int64{0}},
		{[]int64{0}, []int64{1, 2, 3, 4, 7, 8}, []int64{0}},
		{[]int64{1}, []int64{0}, []int64{0}},
		{[]int64{112}, []int64{2}, []int64{224}},
	}

	for counter, tt := range tests {
		testname := fmt.Sprintf("%d", counter+1)
		t.Run(testname, func(t *testing.T) {

			aFr := ff.FromInt64Vec(tt.a)
			bFr := ff.FromInt64Vec(tt.b)
			wantFr := ff.FromInt64Vec(tt.want)
			ansFr := PolyMul(aFr, bFr)

			// debug.DebugFrs("", wantFr)
			// debug.DebugFrs("", ansFr)

			flag := CheckEqualVec(wantFr, ansFr)
			if flag == false {
				t.Errorf("PolyMul: Answer did not match with expected.")
			}
		})
	}
}

func TestPolyLongDiv(t *testing.T) {

	var tests = []struct {
		a, b, want []int64
	}{
		{[]int64{1, 2, 3, 4}, []int64{5, 1}, []int64{87, -17, 4}},
	}

	for counter, tt := range tests {
		testname := fmt.Sprintf("%d", counter+1)
		t.Run(testname, func(t *testing.T) {

			aFr := ff.FromInt64Vec(tt.a)
			bFr := ff.FromInt64Vec(tt.b)
			wantFr := ff.FromInt64Vec(tt.want)
			ansFr := PolyLongDiv(aFr, bFr)

			flag := CheckEqualVec(wantFr, ansFr)
			if flag == false {
				t.Errorf("PolyLongDiv: Quotient did not match with expected")
			}
		})
	}
}

func TestPolyDiv(t *testing.T) {

	var tests = []struct {
		a, b, qwant, rwant []int64
	}{
		{[]int64{1, 2, 3, 4}, []int64{5, 1}, []int64{87, -17, 4}, []int64{-434}},
		{[]int64{8, 10, -5, 3}, []int64{-3, 2, 1}, []int64{-11, 3}, []int64{-25, 41}},
	}

	for counter, tt := range tests {
		testname := fmt.Sprintf("%d", counter+1)
		t.Run(testname, func(t *testing.T) {

			aFr := ff.FromInt64Vec(tt.a)
			bFr := ff.FromInt64Vec(tt.b)
			qwantFr := ff.FromInt64Vec(tt.qwant)
			rwantFr := ff.FromInt64Vec(tt.rwant)
			qFr, rFr := PolyDiv(aFr, bFr)

			// debug.DebugFrs("", wantFr)
			// debug.DebugFrs("", ansFr)

			var flag bool

			flag = CheckEqualVec(qwantFr, qFr)
			if flag == false {
				t.Errorf("PolyDiv: Quotient did not match with expected.")
			}
			flag = CheckEqualVec(rwantFr, rFr)
			if flag == false {
				t.Errorf("PolyDiv: Remainder did not match with expected.")
			}
		})
	}
}
