package asm_gen

import "he++/tac"

type x86_64Reg string

const (
	None x86_64Reg = "none"
	RDI  x86_64Reg = "rdi"
	RSI  x86_64Reg = "rsi"
	RDX  x86_64Reg = "rdx"
	RCX  x86_64Reg = "rcx"
	R8   x86_64Reg = "r8"
	R9   x86_64Reg = "r9"
	R10  x86_64Reg = "r10"
	R11  x86_64Reg = "r11"
	RBX  x86_64Reg = "rbx"

	XMM0 x86_64Reg = "xmm0"
	XMM1 x86_64Reg = "xmm1"
	XMM2 x86_64Reg = "xmm2"
	XMM3 x86_64Reg = "xmm3"
	XMM4 x86_64Reg = "xmm4"
	XMM5 x86_64Reg = "xmm5"
	XMM6 x86_64Reg = "xmm6"
	XMM7 x86_64Reg = "xmm7"

	RAX x86_64Reg = "rax"

	RBP x86_64Reg = "rbp"
	RSP x86_64Reg = "rsp"
)

func vRegComparator(ftac *tac.FunctionTAC, ra, rb tac.VirtualRegisterNumber) bool {
	// decide priority based on lifetime, use frequency etc.
	lifes := ftac.RegLifetimes()
	l1, ex := lifes[ra]
	if !ex {
		panic("Should have had lifetime info")
	}
	l2, ex := lifes[rb]
	if !ex {
		panic("Should have had lifetime info")
	}
	return l1.End > l2.End
}
