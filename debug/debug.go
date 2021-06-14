package debug

import (
	"fmt"
	"strings"

	gmcl "github.com/alinush/go-mcl"
	"github.com/hyperproofs/kzg-go/ff"
)

func DebugFrPtrs(msg string, values []*gmcl.Fr) {
	var out strings.Builder
	out.WriteString("---")
	out.WriteString(msg)
	out.WriteString("---\n")
	for i := range values {
		out.WriteString(fmt.Sprintf("#%4d: %s\n", i, ff.FrStr(values[i])))
	}
	fmt.Println(out.String())
}

func DebugFrs(msg string, values []gmcl.Fr) {
	fmt.Println("---------------------------")
	var out strings.Builder
	for i := range values {
		out.WriteString(fmt.Sprintf("%s %d: %s\n", msg, i, ff.FrStr(&values[i])))
	}
	fmt.Print(out.String())
}
