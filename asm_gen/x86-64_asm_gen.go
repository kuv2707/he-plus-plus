package asm_gen

import (
	"fmt"
	"he++/tac"
	"he++/utils"
)

type Location struct {
	reg    x86_64Reg
	offset int // non zero offset means stack allocated
}

func (l Location) String() string {
	if l.offset != 0 {
		return fmt.Sprintf("[%s - %d]", l.reg, l.offset)
	}
	return string(l.reg)
}

// following SysV ABI
type AsmGen struct {
	tacHandler *tac.TACHandler
}

func NewAsmGen(tacHandler *tac.TACHandler) AsmGen {
	return AsmGen{
		tacHandler: tacHandler,
	}
}

func (ag *AsmGen) GenerateAsm() {
	for fname, ftac := range ag.tacHandler.TacBlocks {
		fasm := MakeFunctionAsm(ftac)
		fasm.GenerateAsm()
		fmt.Println("asm for", fname)
		for _, v := range fasm.instrs {
			fmt.Println(v)
		}
		fmt.Println()
	}
}

type FunctionAsm struct {
	VRegMapping map[tac.VirtualRegisterNumber]Location
	// regsInUse           map[tac.VirtualRegisterNumber]tac.Life
	intRegListOrdered   []x86_64Reg
	floatRegListOrdered []x86_64Reg
	intReplaceList      utils.Heap[tac.VirtualRegisterNumber]
	floatReplaceList    utils.Heap[tac.VirtualRegisterNumber]
	ftac                *tac.FunctionTAC
	instrs              []x86_64Instr
	stackFrameSize      int
}

func MakeFunctionAsm(ftac *tac.FunctionTAC) FunctionAsm {
	regComp := func(a, b tac.VirtualRegisterNumber) bool {
		return vRegComparator(ftac, a, b)
	}
	fasm := FunctionAsm{
		VRegMapping: make(map[tac.VirtualRegisterNumber]Location),
		// regsInUse:           make(map[tac.VirtualRegisterNumber]tac.Life),
		intRegListOrdered:   []x86_64Reg{RDI, RSI, RDX, RCX, R8, R9, R10, R11},
		floatRegListOrdered: []x86_64Reg{XMM0, XMM1, XMM2, XMM3, XMM4, XMM5, XMM6, XMM7},
		// these heaps store vregs currently mapped to real regs but which may be spilled
		// to stack if some other vreg is better suited to a real reg
		intReplaceList:   utils.MakeHeap(regComp),
		floatReplaceList: utils.MakeHeap(regComp),
		ftac:             ftac,
		instrs:           make([]x86_64Instr, 0),
		stackFrameSize:   0,
	}
	return fasm
}

func (fasm *FunctionAsm) emitInstr(ins x86_64Instr) {
	fasm.instrs = append(fasm.instrs, ins)
}

func (fasm *FunctionAsm) spill(vreg tac.VirtualRegisterNumber) (x86_64Reg, Location) {
	dataSize := fasm.ftac.GetDataRegCategory(vreg).SizeBytes()
	fasm.stackFrameSize += dataSize
	stLoc := Location{RSP, fasm.stackFrameSize}
	freedReg := None
	if reg, ex := fasm.VRegMapping[vreg]; ex {
		// the mapping here should be to reg and not to stack.
		freedReg = reg.reg
	}
	fasm.VRegMapping[vreg] = stLoc
	return freedReg, stLoc
}

func (fasm *FunctionAsm) getLocationFor(vreg tac.VirtualRegisterNumber,
	heap *utils.Heap[tac.VirtualRegisterNumber]) Location {
	// vreg doesn't have a mapping
	// we take the weakest vreg and spill it to stack. And use its vacated
	// reg for this vreg. If this vreg is the weakest, it gets a space in the
	// stack
	heap.Push(vreg)
	toSpill, ex := heap.Pop()
	if !ex {
		panic("Err in heap")
	}
	freedReg, stLoc := fasm.spill(toSpill)
	if toSpill == vreg {
		// mapping added in .spill
		return stLoc
	}
	loc := Location{freedReg, 0}
	fasm.VRegMapping[vreg] = loc
	return loc
}

func (fasm *FunctionAsm) RequestLocation(vreg *tac.VRegArg) Location {
	// linear reg alloc
	// if a mapping already exists for this vreg, just return that
	if loc, ex := fasm.VRegMapping[vreg.RegNo]; ex {
		return loc
	}

	if fasm.ftac.GetDataRegCategory(vreg.RegNo).IsFloating() {
		// todo
	} else {
		if len(fasm.intRegListOrdered) == 0 {
			loc := fasm.getLocationFor(vreg.RegNo, &fasm.intReplaceList)
			fasm.VRegMapping[vreg.RegNo] = loc
			return loc
		} else {
			reg := fasm.intRegListOrdered[0]
			fasm.intRegListOrdered = fasm.intRegListOrdered[1:]
			return Location{reg, 0}
		}
	}
	return Location{None, -1}
}

func (fasm *FunctionAsm) GenerateAsm() {
	for _, ins := range fasm.ftac.Instrs() {
		switch v := ins.(type) {
		case *tac.AssignInstr:
			fasm.genAsmForAssign(v)
		case *tac.BinaryOpInstr:
			fasm.genAsmForBinary(v)
		case *tac.CJumpInstr:
			fasm.genAsmForCJump(v)
		}
	}
}

func (fasm *FunctionAsm) instrParam(arg tac.TACOpArg) string {
	switch v := arg.(type) {
	case *tac.ImmIntArg:
		return fmt.Sprintf("%d", v.Num())
	case *tac.ImmFloatArg:
		return fmt.Sprintf("%f", v.Num())
	case *tac.VRegArg:
		loc := fasm.RequestLocation(v).String()
		fmt.Println(arg, " : ", loc)
		return loc
	case *tac.NULLOpArg:
		return "<NULL>"
	}
	panic("switch not exhaustive!")
}

func (fasm *FunctionAsm) genAsmForAssign(v *tac.AssignInstr) {
	vregTo, vregArg, _ := v.ThreeAdresses()
	fasm.emitInstr(x86_64Instr{instrName: MOV,
		params: []string{fasm.instrParam(*vregTo), fasm.instrParam(*vregArg)},
		labels: v.Labels(),
	})
}

func (fasm *FunctionAsm) genAsmForBinary(v *tac.BinaryOpInstr) {
	vregTo, vregA1, vregA2 := v.ThreeAdresses()
	to := fasm.instrParam(*vregTo)
	fasm.emitInstr(x86_64Instr{instrName: MOV,
		params: []string{to, fasm.instrParam(*vregA1)},
		labels: v.Labels(),
	})
	fasm.emitInstr(x86_64Instr{
		instrName: opInstrName(v.Operator()),
		params:    []string{to, fasm.instrParam(*vregA2)},
	})
}

func (fasm *FunctionAsm) genAsmForCJump(v *tac.CJumpInstr) {
	_, argL, argR := v.ThreeAdresses()
	op := OppositeCompOp(v.Op)

	fasm.emitInstr(x86_64Instr{
		instrName: CMP,
		params:    []string{fasm.instrParam(*argL), fasm.instrParam(*argR)},
		labels:    v.Labels(),
	})
	fasm.emitInstr(x86_64Instr{
		instrName: compOpsName[string(op)],
		params:    []string{v.JmpToLabel},
	})
}
