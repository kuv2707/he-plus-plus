package backend

import (
	// "fmt"
	"he++/utils"
)

// constant folding

func (ftac *FunctionTAC) foldConstants() {
	regNum := make(map[VirtualRegisterNumber]int64) // include size of variable
	// since I currently do not reuse regs used to load Imms, no need to remove
	// entries in numReg
	for i, instr := range ftac.instrs {
		switch v := instr.(type) {
		case *AssignInstr:
			{
				if v.arg.locType == Immediate {
					if _, ex := regNum[v.assnTo.ival]; ex {
						// this instr is useless as the number is already mapped
						// to some other register.
						ftac.instrs[i] = nil
					}
					// we still need to create mapping
					regNum[v.assnTo.ival] = v.arg.ival
				} else {
					fold(&v.arg, regNum)
				}
			}
		case *BinaryOpInstr:
			{
				fold(&v.arg1, regNum)
				fold(&v.arg2, regNum)
			}
		case *UnaryOpInstr:
			{
			}
		case *AllocInstr:
			{
				fold(&v.sizeReg, regNum)
			}
		case *MemStoreInstr:
			{
				fold(&v.storeWhat, regNum)
			}
		case *CJumpInstr:
			{
				fold(&v.argL, regNum)
				fold(&v.argR, regNum)
			}
		case *ParamInstr:
			{
				fold(&v.arg, regNum)
			}
		case *JumpInstr, *LabelPlaceholder, *LoadLabelInstr, *MemLoadInstr, *CallInstr:
			{

			}
		default:
			panic("Not impl for " + v.String())
		}
	}
}

func fold(arg *TACOpArg, regNum map[int64]VirtualRegisterNumber) {
	if arg.locType != VirtualRegister {
		return
	}
	numval, exis := regNum[arg.ival]
	if !exis {
		return
	}
	arg.locType = Immediate
	// fmt.Println("Swap reg", arg.ival, numval)
	arg.ival = numval
}

// eliminating useless instrs and vregs ie,
// those whose data doesn't flow into instructions: `param`, `call`, `store`.
func (ftac *FunctionTAC) Prune() {
	depReg := make(map[VirtualRegisterNumber][]VirtualRegisterNumber)
	usefulRegs := make(map[VirtualRegisterNumber]bool)
	for _, instr := range ftac.instrs {
		switch v := instr.(type) {
		case *AssignInstr:
			{
				addDataFlowEntry(depReg, &v.arg, &v.assnTo)
			}
		case *BinaryOpInstr:
			{
				addDataFlowEntry(depReg, &v.arg1, &v.assnTo)
				addDataFlowEntry(depReg, &v.arg2, &v.assnTo)
			}
		case *UnaryOpInstr:
			{
				addDataFlowEntry(depReg, &v.arg1, &v.assnTo)
			}
		case *AllocInstr:
			{
				addDataFlowEntry(depReg, &v.sizeReg, &v.ptrToAlloc)
			}
		case *MemStoreInstr:
			{
				addDataFlowEntry(depReg, &v.storeWhat, &v.storeAt)
				usefulRegs[v.storeAt.ival] = true
				usefulRegs[v.storeWhat.ival] = true

			}
		case *CJumpInstr:
			{
				usefulRegs[v.argL.ival] = true
				usefulRegs[v.argR.ival] = true
			}
		case *ParamInstr:
			{
				usefulRegs[v.arg.ival] = true
			}
		case *MemLoadInstr:
			{
				addDataFlowEntry(depReg, &v.loadFrom, &v.storeAt)
			}
		case *CallInstr:
			{
				addDataFlowEntry(depReg, &v.retReg, &v.calleeAddr)
			}
		case *JumpInstr, *LabelPlaceholder, *LoadLabelInstr:
			{

			}
		}
	}
	q := utils.MakeQueue[VirtualRegisterNumber]()
	for a := range usefulRegs {
		q.Push(a)
	}
	for !q.Empty() {
		t := q.Pop()
		for _, nb := range depReg[t] {
			if _, ex := usefulRegs[nb]; ex {
				continue
			}
			usefulRegs[nb] = true
			q.Push(nb)
		}
	}

	for i, instr := range ftac.instrs {
		switch v := instr.(type) {
		case *AssignInstr:
			{
				if !usefulRegs[v.assnTo.ival] {
					ftac.instrs[i] = nil
				}
			}
		case *BinaryOpInstr:
			{
				if !usefulRegs[v.assnTo.ival] {
					ftac.instrs[i] = nil
				}
			}
		case *UnaryOpInstr:
			{
				if !usefulRegs[v.assnTo.ival] {
					ftac.instrs[i] = nil
				}
			}
		case *AllocInstr:
			{
				if !usefulRegs[v.ptrToAlloc.ival] {
					ftac.instrs[i] = nil
				}
			}

		case *CJumpInstr:
			{
			}
		case *ParamInstr:
			{
			}
		case *MemLoadInstr:
			{
				if !usefulRegs[v.storeAt.ival] {
					ftac.instrs[i] = nil
				}
			}
		case *CallInstr:
			{
				if !usefulRegs[v.retReg.ival] {
					v.retReg = NOWHERE
					ftac.instrs[i] = v
				}
			}
		case *JumpInstr, *LabelPlaceholder, *LoadLabelInstr:
			{

			}
		}
	}
	pruned := make([]ThreeAddressInstr, 0)
	for _, ins := range ftac.instrs {
		if ins != nil {
			pruned = append(pruned, ins)
		} else {
			// fmt.Println("DEL instr", i)
		}
	}
	ftac.instrs = pruned
}

func addDataFlowEntry(adj map[int64][]int64, dest *TACOpArg, src *TACOpArg) {
	if src.locType != VirtualRegister || dest.locType != VirtualRegister {
		return
	}

	adj[dest.ival] = append(adj[dest.ival], src.ival)
}
