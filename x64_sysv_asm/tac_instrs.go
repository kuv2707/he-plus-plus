package x64sysvasm

import (
	"fmt"
	"he++/utils"
	"strings"
)

type ThreeAddressInstr interface {
	String() string
	Labels() []string
}

type TACBaseInstr struct {
	labels []string
}

func (b *TACBaseInstr) Labels() []string {
	return b.labels
}

func LabInstrStr(k ThreeAddressInstr, s string) string {
	return strings.Join(k.Labels(), ":\n") + s
}

type BinaryOpInstr struct {
	TACBaseInstr
	assnTo TACOpArg
	op     string // todo: use enum
	arg1   TACOpArg
	arg2   TACOpArg
}

func (b *BinaryOpInstr) String() string {
	return LabInstrStr(b, fmt.Sprintf("%v = %v %s %v", utils.Cyan(b.assnTo.String()), b.arg1, b.op, b.arg2))
}

type UnaryOpInstr struct {
	TACBaseInstr
	assTo TACOpArg
	op    string // todo: use enum
	arg1  TACOpArg
}

func (u *UnaryOpInstr) String() string {
	return LabInstrStr(u, fmt.Sprintf("%v = %s%v", utils.Cyan(u.assTo.String()), u.op, u.arg1))
}

type JumpInstr struct {
	TACBaseInstr
	jmpToLabel string
}

func (j *JumpInstr) String() string {
	return LabInstrStr(j, fmt.Sprintf("jmp %v", j.jmpToLabel))
}

type CJumpInstr struct {
	TACBaseInstr
	op   string
	argL TACOpArg
	argR TACOpArg
	// to jump to if the op yielded false
	jmpToLabel string
}

func (j *CJumpInstr) String() string {
	return LabInstrStr(j, fmt.Sprintf("cjmp %s %s %s : %v", j.argL, j.op, j.argR, j.jmpToLabel))
}

type ParamInstr struct {
	TACBaseInstr
	arg TACOpArg
}

func (j *ParamInstr) String() string {
	return LabInstrStr(j, fmt.Sprintf("param %v", j.arg))
}

type CallInstr struct {
	TACBaseInstr
	calleeAddr TACOpArg
}

func (j *CallInstr) String() string {
	return LabInstrStr(j, fmt.Sprintf("call %v", j.calleeAddr))
}

type LoadLabelInstr struct {
	TACBaseInstr
	loadeeLabel string
	to          TACOpArg
}

func (j *LoadLabelInstr) String() string {
	return LabInstrStr(j, fmt.Sprintf("%v = load %v", j.to, utils.Green(j.loadeeLabel)))
}
