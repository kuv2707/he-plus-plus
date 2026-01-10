package asm_gen

import "he++/tac"

type x86_64Reg int

const (
	None x86_64Reg = iota
	RDI
	RSI
	RDX
	RCX
	R8
	R9
	R10
	R11
	RBX

	XMM0
	XMM1
	XMM2
	XMM3
	XMM4
	XMM5
	XMM6
	XMM7

	RAX // ret value

	RBP
	RSP
)

type Location struct {
	reg    x86_64Reg
	offset int // non zero offset means stack allocated
}

// following SysV ABI
type AsmGen struct {
}

type FunctionAsm struct {
	VRegMapping         map[tac.VirtualRegisterNumber]Location
	regsInUse           map[tac.VirtualRegisterNumber]tac.Life
	intRegListOrdered   []x86_64Reg
	
	floatRegListOrdered []x86_64Reg
}

func MakeFunctionAsm() FunctionAsm {
	fasm := FunctionAsm{
		VRegMapping:         make(map[tac.VirtualRegisterNumber]Location),
		regsInUse:           make(map[tac.VirtualRegisterNumber]tac.Life),
		intRegListOrdered:   []x86_64Reg{RDI, RSI, RDX, RCX, R8, R9, R10, R11},
		floatRegListOrdered: []x86_64Reg{XMM0, XMM1, XMM2, XMM3, XMM4, XMM5, XMM6, XMM7},
	}
	return fasm
}

func (fasm *FunctionAsm) GenerateAsm(ftac *tac.FunctionTAC) {
	for i, ins := range ftac.Instrs() {
		switch ins.(type) {

		}
	}
}
