package tac

import (
	"fmt"
	"he++/utils"
	"strings"
)

type ThreeAddressInstr interface {
	String() string
	Labels() []string
	ThreeAdresses() (dest, src1, src2 *TACOpArg)
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
	return LabInstrStr(b, fmt.Sprintf("%v = %v %s %v", b.assnTo, b.arg1, b.op, b.arg2))
}

func (b *BinaryOpInstr) ThreeAdresses() (*TACOpArg, *TACOpArg, *TACOpArg) {
	return &b.assnTo, &b.arg1, &b.arg2
}

type UnaryOpInstr struct {
	TACBaseInstr
	assnTo TACOpArg
	op     string // todo: use enum
	arg1   TACOpArg
}

func (u *UnaryOpInstr) String() string {
	return LabInstrStr(u, fmt.Sprintf("%v = %s %v", u.assnTo, u.op, u.arg1))
}

func (u *UnaryOpInstr) ThreeAdresses() (*TACOpArg, *TACOpArg, *TACOpArg) {
	return &u.assnTo, &u.arg1, &NOWHERE
}

// special kind of UnaryOpInstr
type AssignInstr struct {
	TACBaseInstr
	assnTo TACOpArg
	arg    TACOpArg
}

func (u *AssignInstr) String() string {
	return LabInstrStr(u, fmt.Sprintf("%v = %v", u.assnTo, u.arg))
}

func (a *AssignInstr) ThreeAdresses() (*TACOpArg, *TACOpArg, *TACOpArg) {
	return &a.assnTo, &a.arg, &NOWHERE
}

// Special kind of CJumpInstr
type JumpInstr struct {
	TACBaseInstr
	jmpToLabel string
}

func (j *JumpInstr) String() string {
	return LabInstrStr(j, fmt.Sprintf("%s %v", utils.BoldCyan("jmp"), utils.BoldGreen(j.jmpToLabel)))
}

func (j *JumpInstr) ThreeAdresses() (*TACOpArg, *TACOpArg, *TACOpArg) {
	return &NOWHERE, &NOWHERE, &NOWHERE
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

func (j *CJumpInstr) ThreeAdresses() (*TACOpArg, *TACOpArg, *TACOpArg) {
	return &NOWHERE, &j.argL, &j.argR
}

type ParamInstr struct {
	TACBaseInstr
	arg TACOpArg
}

func (j *ParamInstr) String() string {
	return LabInstrStr(j, fmt.Sprintf("%s %v", utils.BoldCyan("param"), j.arg))
}

func (p *ParamInstr) ThreeAdresses() (*TACOpArg, *TACOpArg, *TACOpArg) {
	return &NOWHERE, &p.arg, &NOWHERE
}

type CallInstr struct {
	TACBaseInstr
	calleeAddr TACOpArg
	retReg     TACOpArg
}

func (j *CallInstr) String() string {
	return LabInstrStr(j, fmt.Sprintf("%s = %s %v", j.retReg, utils.BoldCyan("call"), j.calleeAddr))
}

func (c *CallInstr) ThreeAdresses() (*TACOpArg, *TACOpArg, *TACOpArg) {
	return &c.retReg, &c.calleeAddr, &NOWHERE
}

type LoadLabelInstr struct {
	TACBaseInstr
	loadeeLabel string
	to          TACOpArg
}

func (j *LoadLabelInstr) String() string {
	return LabInstrStr(j, fmt.Sprintf("%s %v [%v]", utils.BoldCyan("load"), j.to, utils.BoldGreen(j.loadeeLabel)))
}

func (l *LoadLabelInstr) ThreeAdresses() (*TACOpArg, *TACOpArg, *TACOpArg) {
	return &l.to, &NOWHERE, &NOWHERE
}

type LabelPlaceholder struct {
	TACBaseInstr
}

func (l *LabelPlaceholder) String() string {
	return LabInstrStr(l, fmt.Sprintf(" PLACEHOLDER"))
}

func (l *LabelPlaceholder) ThreeAdresses() (*TACOpArg, *TACOpArg, *TACOpArg) {
	return &NOWHERE, &NOWHERE, &NOWHERE
}

func placeholderWithLabels(s ...string) *LabelPlaceholder {
	return &LabelPlaceholder{TACBaseInstr{labels: s}}

}

type LoopBoundary struct {
	TACBaseInstr
	loopNo   int
	startEnd bool
}

func (l *LoopBoundary) String() string {
	return LabInstrStr(l, fmt.Sprintf(" loop_boundary_%d", l.loopNo))
}

func (l *LoopBoundary) ThreeAdresses() (*TACOpArg, *TACOpArg, *TACOpArg) {
	return &NOWHERE, &NOWHERE, &NOWHERE
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

func (a *AllocInstr) ThreeAdresses() (*TACOpArg, *TACOpArg, *TACOpArg) {
	return &a.ptrToAlloc, &a.sizeReg, &NOWHERE
}

type MemStoreInstr struct {
	TACBaseInstr
	storeAt   TACOpArg
	storeWhat TACOpArg
	numBytes  int
}

func (m *MemStoreInstr) String() string {
	return fmt.Sprintf("%s [%v], %v (%d bytes)", utils.BoldCyan("store"), m.storeAt, m.storeWhat, m.numBytes)
}

func (m *MemStoreInstr) ThreeAdresses() (*TACOpArg, *TACOpArg, *TACOpArg) {
	return &NOWHERE, &m.storeAt, &m.storeWhat
}

type MemLoadInstr struct {
	TACBaseInstr
	loadFrom TACOpArg
	storeAt  TACOpArg
	numBytes int
}

func (m *MemLoadInstr) String() string {
	return fmt.Sprintf("%v = %s %v, %d bytes", m.storeAt, utils.BoldCyan("loadfrom"), m.loadFrom, m.numBytes)
}

func (m *MemLoadInstr) ThreeAdresses() (*TACOpArg, *TACOpArg, *TACOpArg) {
	return &m.storeAt, &m.loadFrom, &NOWHERE
}

type FuncRetInstr struct {
	TACBaseInstr
	retReg TACOpArg
}

func (f *FuncRetInstr) String() string {
	return LabInstrStr(f, fmt.Sprintf("%s %v", utils.BoldCyan("ret"), f.retReg))
}

func (f *FuncRetInstr) ThreeAdresses() (*TACOpArg, *TACOpArg, *TACOpArg) {
	return &NOWHERE, &f.retReg, &NOWHERE
}

type FuncArgRecvInstr struct {
	TACBaseInstr
	argNo    int
	recvInto TACOpArg
}

func (f *FuncArgRecvInstr) String() string {
	return LabInstrStr(f, fmt.Sprintf("%v = %s %d", f.recvInto, utils.BoldCyan("arg"), f.argNo))
}

func (f *FuncArgRecvInstr) ThreeAdresses() (*TACOpArg, *TACOpArg, *TACOpArg) {
	return &f.recvInto, &NOWHERE, &NOWHERE
}

var LOOP_START_PREFIX = "loop_start_"
var LOOP_END_PREFIX = "loop_end_"

func hasLoopLabel(v ThreeAddressInstr) (bool, bool) {
	labels := v.Labels()
	hasStart := false
	hasEnd := false

	for _, label := range labels {
		if strings.HasPrefix(label, LOOP_START_PREFIX) {
			hasStart = true
		}
		if strings.HasPrefix(label, LOOP_END_PREFIX) {
			hasEnd = true
		}
	}
	// todo assert both aren't simultaneously true
	if !hasStart && !hasEnd {
		return false, false
	}

	return true, hasStart
}
