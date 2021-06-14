package main

import (
	"fmt"

	gmcl "github.com/alinush/go-mcl"
	// "github.com/hyperproofs/kzg-go/ff"
)

func main() {
	fmt.Println("Hello, World!")
	gmcl.InitFromString("bls12-381")
}
