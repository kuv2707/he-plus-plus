package asm_gen

import (
	"fmt"
	"he++/tac"
)

type x86_64Reg struct {
	name_1 string
	name_4 string
	name_8 string
}

func (reg *x86_64Reg) NameForSize(size int) string {
	switch size {
	case 1:
		return reg.name_1
	case 4:
		return reg.name_4
	case 8:
		return reg.name_8
	default:
		return "ERR_" + reg.name_8
	}
}

func (reg *x86_64Reg) String() string {
	return fmt.Sprintf("reg_%s", reg.name_8)
}

// Predefined registers (as pointers)
var (
	None = &x86_64Reg{"", "", "none"}

	RDI = &x86_64Reg{"dil", "edi", "rdi"}
	RSI = &x86_64Reg{"sil", "esi", "rsi"}
	RDX = &x86_64Reg{"dl", "edx", "rdx"}
	RCX = &x86_64Reg{"cl", "ecx", "rcx"}

	R8  = &x86_64Reg{"r8b", "r8d", "r8"}
	R9  = &x86_64Reg{"r9b", "r9d", "r9"}
	R10 = &x86_64Reg{"r10b", "r10d", "r10"}
	R11 = &x86_64Reg{"r11b", "r11d", "r11"}

	RBX = &x86_64Reg{"bl", "ebx", "rbx"}
	RAX = &x86_64Reg{"al", "eax", "rax"}

	RBP = &x86_64Reg{"bpl", "ebp", "rbp"}
	RSP = &x86_64Reg{"spl", "esp", "rsp"}

	// XMM registers (no sub-widths)
	XMM0 = &x86_64Reg{"", "", "xmm0"}
	XMM1 = &x86_64Reg{"", "", "xmm1"}
	XMM2 = &x86_64Reg{"", "", "xmm2"}
	XMM3 = &x86_64Reg{"", "", "xmm3"}
	XMM4 = &x86_64Reg{"", "", "xmm4"}
	XMM5 = &x86_64Reg{"", "", "xmm5"}
	XMM6 = &x86_64Reg{"", "", "xmm6"}
	XMM7 = &x86_64Reg{"", "", "xmm7"}
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
