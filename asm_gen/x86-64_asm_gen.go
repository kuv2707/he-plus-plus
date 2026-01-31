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
		return fmt.Sprintf("%s - %d", l.reg, l.offset)
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
		for i := range fasm.instrs {
			fmt.Println(fasm.instrs[i])
		}
		fmt.Println()
	}
}

type FunctionAsm struct {
	VRegMapping         map[tac.VirtualRegisterNumber]Location
	intRegListOrdered   []x86_64Reg
	floatRegListOrdered []x86_64Reg
	ftac                *tac.FunctionTAC
	instrs              []x86_64Instr
	stackFrameSize      int
}

var TEMPREG = R11

func MakeFunctionAsm(ftac *tac.FunctionTAC) FunctionAsm {

	fasm := FunctionAsm{
		VRegMapping: make(map[tac.VirtualRegisterNumber]Location),

		ftac:           ftac,
		instrs:         make([]x86_64Instr, 0),
		stackFrameSize: 0,
	}
	return fasm
}

func (fasm *FunctionAsm) emitInstr(ins x86_64Instr) {
	fasm.instrs = append(fasm.instrs, ins)
}

func (fasm *FunctionAsm) GenerateAsm() {
	fasm.createVregMapping()
	for k, v := range fasm.VRegMapping {
		fmt.Printf("VR: %v, Loc: %s\n", utils.Red(fmt.Sprint(k)), utils.Cyan(v.String()))
	}
	instrs := fasm.ftac.Instrs()
	for i := range instrs {
		switch v := instrs[i].(type) {
		case *tac.AssignInstr:
			fasm.genAsmForAssign(v)
		case *tac.BinaryOpInstr:
			fasm.genAsmForBinary(v)
		case *tac.JumpInstr:
			fasm.genAsmForJump(v)
		case *tac.CJumpInstr:
			fasm.genAsmForCJump(v)
		case *tac.MemStoreInstr:
			fasm.genAsmForMemStore(v)
		case *tac.MemLoadInstr:
			fasm.genAsmForMemLoad(v)
		case *tac.AllocInstr:
			fasm.genAsmForAlloc(v)
		case *tac.CallInstr:
			fasm.genAsmForCall(v)
		case *tac.FuncArgRecvInstr:
			fasm.genAsmForFuncArgRecv(v)
		case *tac.LoopBoundary:
			// ignore
		default:
			fmt.Println("Not impl for", instrs[i])
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
		loc, ex := fasm.VRegMapping[v.RegNo]
		if !ex {
			panic(fmt.Sprintf("Se esperaba un mapping para %s", v.String()))
		}
		// fmt.Println(arg, " : ", loc)
		return loc.String()
	default:
		return "<NULL>"
	}
}

func (fasm *FunctionAsm) genAsmForAssign(v *tac.AssignInstr) {
	vregTo, vregArg, _ := v.ThreeAdresses()
	p1, p2 := fasm.instrParam(*vregTo), fasm.instrParam(*vregArg)
	if p1 == p2 {
		return
	}
	fasm.emitInstr(x86_64Instr{instrName: MOV,
		params: []string{p1, p2},
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

func (fasm *FunctionAsm) genAsmForJump(v *tac.JumpInstr) {
	fasm.emitInstr(x86_64Instr{
		instrName: JMP,
		params:    []string{v.JmpToLabel},
		labels:    v.Labels(),
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

func (fasm *FunctionAsm) genAsmForMemStore(v *tac.MemStoreInstr) {
	fmt.Println(v)
	sat := fasm.VRegMapping[(v.StoreAt.(*tac.VRegArg)).RegNo]
	swhat := fasm.instrParam(v.StoreWhat)
	dest := sat.String()
	if sat.offset != 0 {
		fasm.emitInstr(x86_64Instr{
			instrName: MOV,
			params:    []string{string(TEMPREG), sat.String()},
			labels:    v.Labels(),
		})

		fasm.emitInstr(x86_64Instr{
			instrName: MOV,
			params:    []string{fmt.Sprintf("[%s]", TEMPREG), swhat},
			labels:    v.Labels(),
		})
	} else {
		fasm.emitInstr(x86_64Instr{
			instrName: MOV,
			params:    []string{fmt.Sprintf("[%s]", dest), swhat},
		})
	}
}

func (fasm *FunctionAsm) genAsmForMemLoad(v *tac.MemLoadInstr) {
	lfrom := fasm.VRegMapping[(v.LoadFrom.(*tac.VRegArg)).RegNo]
	sat := fasm.VRegMapping[(v.StoreAt.(*tac.VRegArg)).RegNo]

	if sat.offset != 0 {
		fasm.emitInstr(x86_64Instr{
			instrName: MOV,
			params:    []string{string(TEMPREG), fmt.Sprintf("[%s]", lfrom)},
			labels:    v.Labels(),
		})
		fasm.emitInstr(x86_64Instr{
			instrName: MOV,
			params:    []string{sat.String(), string(TEMPREG)},
		})
	} else {
		fasm.emitInstr(x86_64Instr{
			instrName: MOV,
			params:    []string{sat.String(), fmt.Sprintf("[%s]", lfrom)},
			labels:    v.Labels(),
		})
	}
}

func (fasm *FunctionAsm) genAsmForFuncArgRecv(v *tac.FuncArgRecvInstr) {
	// needs to be set by caller
}

func (fasm *FunctionAsm) genAsmForAlloc(v *tac.AllocInstr) {
	if intSize, ok := v.SizeReg.(*tac.ImmIntArg); ok {
		fasm.emitInstr(x86_64Instr{
			instrName: LEA,
			params:    []string{fasm.instrParam(v.PtrToAlloc), fmt.Sprintf("alloc_%d_size_%d", v.AllocNo, intSize.Num())},
			labels:    v.Labels(),
		})
	} else {
		// todo
	}
}

func (fasm *FunctionAsm) genAsmForCall(v *tac.CallInstr) {

}
