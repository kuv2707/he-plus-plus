package backend

import (
	// "fmt"
	"fmt"
	"he++/utils"
)

// constant and copy propagation
func (ftac *FunctionTAC) PropagateRegs() {
	foldBlacklist := make(map[VirtualRegisterNumber]bool)
	inLoop := false
	// we should not fold those registers which were both written and
	// read inside the loop
	for _, instr := range ftac.instrs {
		if v, ok := instr.(*LabelPlaceholder); ok {
			k, se := hasLoopLabel(v)
			if k {
				inLoop = se
			}
		} else {
			des, _, _ := instr.ThreeAdresses()
			if des.locType == VReg && inLoop {
				foldBlacklist[des.ival] = true
			}
		}
	}

	fmt.Println("blacklist: ", foldBlacklist)

	propagMap := make(map[VirtualRegisterNumber]*TACOpArg) // size of var ka kya krna h

	for _, instr := range ftac.instrs {
		switch v := instr.(type) {
		case *AssignInstr:
			{
				// for cases like r1 = #5, r2 = r1, r3 = r2 + #blabla
				// we want r2 to be replaced by #5, not r1.
				var substitut *TACOpArg = nil
				if v.arg.locType == VReg {
					if k, exists := propagMap[v.arg.ival]; exists {
						// fmt.Println("SUB", v.assnTo, k)
						substitut = k
					} else {
						// fmt.Println("SUBB", v.assnTo, v.arg)
						substitut = &v.arg
					}
				} else {
					substitut = &v.arg
				}
				v.arg = *substitut
				if _, ex := foldBlacklist[v.assnTo.ival]; !ex {
					propagMap[v.assnTo.ival] = substitut
				}
				// fmt.Println(propagMap)
			}
		default:
			dest, arg1, arg2 := v.ThreeAdresses()
			fold(arg1, propagMap)
			fold(arg2, propagMap)
			// todo: assert dest is Vreg
			delete(propagMap, dest.ival)
		}
	}
}

func fold(arg *TACOpArg, regNum map[VirtualRegisterNumber]*TACOpArg) {
	if arg.locType != VReg {
		return
	}
	p, exis := regNum[arg.ival]
	if !exis {
		return
	}
	arg.locType = p.locType
	arg.ival = p.ival
}

// eliminating useless instrs and vregs ie, those whose data doesn't flow
// into side-effect instrs which are `param`, `store`. `cjump` is kept as it
// can affect memory state indirectly.
func (ftac *FunctionTAC) Prune() {
	depReg := make(map[VirtualRegisterNumber]map[VirtualRegisterNumber]bool)
	usefulRegs := make(map[VirtualRegisterNumber]bool)
	for _, instr := range ftac.instrs {
		dest, src1, src2 := instr.ThreeAdresses()
		addDataFlowEntry(depReg, src1, dest)
		addDataFlowEntry(depReg, src2, dest)
		switch v := instr.(type) {
		case *MemStoreInstr:
			{
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
		case *CallInstr:
			{
				usefulRegs[v.calleeAddr.ival] = true
			}
		}
	}
	// Debug: print depReg adjacency list
	for reg, deps := range depReg {
		if len(deps) > 0 {
			print(reg, " : ")
			for k, _ := range deps {
				fmt.Printf("%d, ", k)
			}
			println()
		}
	}

	q := utils.MakeQueue[VirtualRegisterNumber]()
	for a := range usefulRegs {
		q.Push(a)
	}
	for !q.Empty() {
		t := q.Pop()
		for nb := range depReg[t] {
			if _, ex := usefulRegs[nb]; ex {
				continue
			}
			usefulRegs[nb] = true
			q.Push(nb)
		}
	}

	for i, instr := range ftac.instrs {
		switch v := instr.(type) {
		case *CallInstr:
			{
				if !usefulRegs[v.retReg.ival] {
					v.retReg = NOWHERE
					ftac.instrs[i] = v
				}
			}
		default:
			{
				dest, _, _ := instr.ThreeAdresses()
				if dest.locType == VReg && !usefulRegs[dest.ival] {
					ftac.instrs[i] = nil
					fmt.Println("X", dest, i)
				}
			}
		}
	}
	pruned := make([]ThreeAddressInstr, 0)
	for _, ins := range ftac.instrs {
		if ins != nil {
			pruned = append(pruned, ins)
		}
	}
	ftac.instrs = pruned
}

func addDataFlowEntry(adj map[VirtualRegisterNumber]map[VirtualRegisterNumber]bool, edgeSrc *TACOpArg, edgeDest *TACOpArg) {
	if edgeSrc.locType != VReg || edgeDest.locType != VReg {
		return
	}
	if edgeSrc.ival == edgeDest.ival {
		return
	}
	if adj[edgeDest.ival] == nil {
		adj[edgeDest.ival] = make(map[VirtualRegisterNumber]bool)
	}
	adj[edgeDest.ival][edgeSrc.ival] = true
}

func (ftac *FunctionTAC) livenessAnalysis() {
	// iterate backwards and store encountered registers so far in a map
	// remove those instructions whose `dest` is not present in the map
}
