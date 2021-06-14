// +build !bignum_pure,!bignum_hol256

package kzg

import (
	gmcl "github.com/alinush/go-mcl"
)

type KZG1Settings struct {

	// setup values
	// [b.multiply(b.G1, pow(s, i, MODULUS)) for i in range(WIDTH+1)],
	PK []gmcl.G1
	// [b.multiply(b.G2, pow(s, i, MODULUS)) for i in range(WIDTH+1)],
	VK []gmcl.G2
}

func NewKZG1Settings(pk []gmcl.G1, vk []gmcl.G2) *KZG1Settings {
	if len(vk) != 2 {
		panic("VK size has to be 2!")
	}
	ks := &KZG1Settings{

		PK: pk,
		VK: vk,
	}

	return ks
}

type KZG2Settings struct {

	// setup values
	// [b.multiply(b.G1, pow(s, i, MODULUS)) for i in range(WIDTH+1)],
	PK []gmcl.G2
	// [b.multiply(b.G2, pow(s, i, MODULUS)) for i in range(WIDTH+1)],
	VK []gmcl.G1
}

func NewKZG2Settings(pk []gmcl.G2, vk []gmcl.G1) *KZG2Settings {
	if len(vk) != 2 {
		panic("VK size has to be 2!")
	}

	ks := &KZG2Settings{
		PK: pk,
		VK: vk,
	}

	return ks
}
