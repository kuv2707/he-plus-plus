package backend 

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
	return utils.BoldGreen(strings.Join(k.Labels(), ":\n")) + s
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
	return LabInstrStr(j, fmt.Sprintf("jmp %v", utils.BoldGreen(j.jmpToLabel)))
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
	return LabInstrStr(j, fmt.Sprintf("jmp_if_false %s %s %s , %v", j.argL, j.op, j.argR, utils.BoldGreen(j.jmpToLabel)))
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
	retReg     TACOpArg
}

func (j *CallInstr) String() string {
	return LabInstrStr(j, fmt.Sprintf("%s = call %v", j.retReg, j.calleeAddr))
}

type LoadLabelInstr struct {
	TACBaseInstr
	loadeeLabel string
	to          TACOpArg
}

func (j *LoadLabelInstr) String() string {
	return LabInstrStr(j, fmt.Sprintf("%v = load %v", j.to, utils.BoldGreen(j.loadeeLabel)))
}

type LabelPlaceholder struct {
	TACBaseInstr
}

func (l *LabelPlaceholder) String() string {
	return LabInstrStr(l, fmt.Sprintf(" PLACEHOLDER"))
}

func placeholderWithLabels(s ...string) *LabelPlaceholder {
	return &LabelPlaceholder{TACBaseInstr{labels: s}}
}

type AllocType byte

const (
	STACK_ALLOC AllocType = 's'
	HEAP_ALLOC  AllocType = 'h'
)

type AllocInstr struct {
	TACBaseInstr
	allocType  AllocType
	sizeReg    TACOpArg
	ptrToAlloc TACOpArg
}

func (a *AllocInstr) String() string {
	return LabInstrStr(a, fmt.Sprintf("%v = alloc(%c, %v)", a.ptrToAlloc, a.allocType, a.sizeReg))
}

type MemStoreInstr struct {
	TACBaseInstr
	storeAt  TACOpArg
	storeWhat TACOpArg
	numBytes int
}

func (m *MemStoreInstr) String() string {
	return fmt.Sprintf("store [%v], %v (%d bytes)", m.storeAt, m.storeWhat, m.numBytes)
}

type MemLoadInstr struct {
	TACBaseInstr
	loadFrom  TACOpArg
	storeAt TACOpArg
	numBytes int
}

func (m *MemLoadInstr) String() string {
	return fmt.Sprintf("%v = load %v, %d bytes",m.storeAt, m.loadFrom, m.numBytes)
}
