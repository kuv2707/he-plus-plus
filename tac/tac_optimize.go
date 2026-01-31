package tac

import (
	"fmt"
	"he++/utils"
)

type TACContext struct {
	regLifetimes  map[VirtualRegisterNumber]Life
	loopLifetimes map[int]Life
	loopWritelog  map[int]map[VirtualRegisterNumber]bool
}

func (ftac *FunctionTAC) Optimize() {
	ctx := ftac.livenessAnalysis()
	ftac.PropagateRegs(&ctx)
	eliminatedRegs := ftac.Prune()

	ftac.eliminateNilInstrs()
	ftac.removeRedundantInstrs()
	ctx = ftac.livenessAnalysis()
	fmt.Println("Reglifetimes:", ctx.regLifetimes)
	fmt.Println("looplifetimes:", ctx.loopLifetimes)
	fmt.Println("Eliminated regs: ", eliminatedRegs)
	ftac.ctx = ctx
}

// constant and copy propagation
func (ftac *FunctionTAC) PropagateRegs(ctx *TACContext) {
	loopStack := utils.MakeStack(-1)

	propagMap := make(map[VirtualRegisterNumber]*TACOpArg) // size of var ka kya krna h

	// for cases like r1 = #5, r2 = r1, r3 = r2 + #blabla
	// we want r2 to be replaced by #5, not r1.
	for i := 0; i < len(ftac.instrs); i++ {
		switch v := ftac.instrs[i].(type) {
		case *LoopBoundary:
			{
				if v.StartEnd {
					loopStack.Push(v.loopNo)
				} else {
					loopStack.Pop()
				}
			}
		default:
			dest, arg1, arg2 := v.ThreeAdresses()
			lno, _ := loopStack.Peek()
			if canPropagateTo(arg1, *lno, ctx) {
				fold(arg1, propagMap)
			}
			if canPropagateTo(arg2, *lno, ctx) {
				fold(arg2, propagMap)
			}

			simpInstr := ftac.simplifyInstr(ftac.instrs[i])
			if assn, ok := simpInstr.(*AssignInstr); ok {
				assto := assn.assnTo.(*VRegArg)
				propagMap[assto.RegNo] = &assn.arg

			} else {
				assto, ok := (*dest).(*VRegArg)
				if ok {
					delete(propagMap, assto.RegNo)
					/**
					if r1 = r2; r2 = 15; r3 = r1 + 3
					then in 3rd instr we want r1 to be the value in r2 before 2nd instr
					but since we've mapped r1 -> r2, we make 3rd instr: r3 = r2 + 1 == 17
					*/
					// currently it seems that the way we generate instructions avoids this failure mode

					// todo: remove those entries in the map which map to assto.regNo too!
				}

			}
			ftac.instrs[i] = simpInstr
		}
	}
}

func fold(arg *TACOpArg, propagMap map[VirtualRegisterNumber]*TACOpArg) {
	if vregarg, ok := (*arg).(*VRegArg); ok {
		if replace, ex := propagMap[vregarg.RegNo]; ex {
			// assign replace to arg
			*arg = *replace
		}
	}
}

func canPropagateTo(arg *TACOpArg, currLoop int, ctx *TACContext) bool {
	loopLife, ex := ctx.loopLifetimes[currLoop]
	if !ex {
		return true
	}
	if varg, ok := (*arg).(*VRegArg); ok {
		// idea: we should not fold those registers satisfying all of:
		// lifetime extends beyond the loop
		// written to inside the loop
		// read from inside the loop.
		regLife, ex := ctx.regLifetimes[varg.RegNo]
		if !ex {
			panic("lifetime info should have existed for arg " + varg.String())
		}
		if regLife.Start < loopLife.Start || regLife.End > loopLife.End {
			if hasBeenWritten(varg.RegNo, currLoop, ctx.loopWritelog) {
				return false
			}
		}
	} else {
		return false // can only propagate to a vreg arg
	}

	return true
}

// Dead code elimination
// eliminating useless instrs and vregs ie, those whose data doesn't flow
// into side-effect instrs which are `param`, `store`. `cjump` is kept as it
// can affect memory state indirectly.
func (ftac *FunctionTAC) Prune() map[*VRegArg]bool {
	depReg := make(map[VirtualRegisterNumber]map[VirtualRegisterNumber]bool) // edge directed from dest to srcs
	usefulRegs := make(map[VirtualRegisterNumber]bool)
	markUsefulReg := func(arg TACOpArg) {
		if varg, ok := arg.(*VRegArg); ok {
			usefulRegs[varg.RegNo] = true
		}
	}
	for _, instr := range ftac.instrs {
		dest, src1, src2 := instr.ThreeAdresses()
		addDataFlowEntry(depReg, *src1, *dest)
		addDataFlowEntry(depReg, *src2, *dest)
		switch v := instr.(type) {
		case *MemStoreInstr:
			{
				markUsefulReg(v.StoreAt)
				markUsefulReg(v.StoreWhat)

			}
		case *CJumpInstr:
			{
				markUsefulReg(v.argL)
				markUsefulReg(v.argR)
			}
		case *ParamInstr:
			{
				markUsefulReg(v.arg)
			}
		case *CallInstr:
			{
				markUsefulReg(v.calleeAddr)
			}
		case *FuncRetInstr:
			{
				markUsefulReg(v.retReg)
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

	fmt.Println("Initial useful regs:", usefulRegs)

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

	eliminatedRegs := make(map[*VRegArg]bool)
	for i, instr := range ftac.instrs {
		switch v := instr.(type) {
		case *CallInstr:
			{
				k := v.retReg.(*VRegArg) // should always be valid
				if !usefulRegs[k.RegNo] {
					v.retReg = NOWHERE
					ftac.instrs[i] = v // maybe unnecessary
				}
			}
		default:
			{
				dest, _, _ := instr.ThreeAdresses()
				if v, ok := (*dest).(*VRegArg); ok && !usefulRegs[v.RegNo] {
					ftac.instrs[i] = nil
					eliminatedRegs[v] = true
				}
			}
		}
	}
	return eliminatedRegs
}

func addDataFlowEntry(adj map[VirtualRegisterNumber]map[VirtualRegisterNumber]bool, edgeSrc TACOpArg,
	edgeDest TACOpArg) {
	vesrc, ok1 := edgeSrc.(*VRegArg)
	vedst, ok2 := edgeDest.(*VRegArg)
	if !ok1 || !ok2 {
		return
	}
	if vesrc.RegNo == vedst.RegNo {
		return
	}

	if adj[vedst.RegNo] == nil {
		adj[vedst.RegNo] = make(map[VirtualRegisterNumber]bool)
	}
	adj[vedst.RegNo][vesrc.RegNo] = true
}

type Life struct {
	Start int
	End   int
}

func (ftac *FunctionTAC) livenessAnalysis() TACContext {
	// iterate backwards and store encountered registers so far in a map
	// remove those instructions whose `dest` is not present in the map

	regLifetimes := make(map[VirtualRegisterNumber]Life)
	loopLifetimes := make(map[int]Life)
	loopWritelog := make(map[int]map[VirtualRegisterNumber]bool)
	updateLiveness := func(a *TACOpArg, i int) {
		arg, ok := (*a).(*VRegArg)
		if !ok {
			return
		}
		if k, ex := regLifetimes[arg.RegNo]; !ex {
			regLifetimes[arg.RegNo] = Life{i, i}
		} else {
			k.End = max(k.End, i)
			k.Start = min(k.Start, i)
			regLifetimes[arg.RegNo] = k
		}

	}
	registerWriteInLoop := func(reg VirtualRegisterNumber, lno int) {
		if loopWritelog[lno] == nil {
			loopWritelog[lno] = make(map[VirtualRegisterNumber]bool)
		}
		loopWritelog[lno][reg] = true
	}
	loopStack := utils.MakeStack(-1)
	for i, instr := range utils.Backwards(ftac.instrs) {
		dest, arg1, arg2 := instr.ThreeAdresses()
		updateLiveness(dest, i)
		updateLiveness(arg1, i)
		updateLiveness(arg2, i)
		if v, ok := (*dest).(*VRegArg); ok {
			// register it as being written to in this loop
			lno, _ := loopStack.Peek() // assert that lno is not zero value
			if *lno != -1 {
				registerWriteInLoop(v.RegNo, *lno)
			}
		}
		if v, ok := instr.(*LoopBoundary); ok {
			if k, ex := loopLifetimes[v.loopNo]; !ex {
				loopLifetimes[v.loopNo] = Life{i, i}
				loopStack.Push(v.loopNo)

			} else {
				k.Start = min(k.Start, i)
				loopLifetimes[v.loopNo] = k
				loopStack.Pop()
			}
		}
	}

	return TACContext{
		regLifetimes:  regLifetimes,
		loopLifetimes: loopLifetimes,
		loopWritelog:  loopWritelog,
	}
}

func hasBeenWritten(reg VirtualRegisterNumber, lno int, wlog map[int]map[VirtualRegisterNumber]bool) bool {
	if wlog[lno] == nil {
		return false
	}
	return wlog[lno][reg]
}

func (ftac *FunctionTAC) eliminateNilInstrs() {
	pruned := make([]ThreeAddressInstr, 0)
	for _, ins := range ftac.instrs {
		if ins != nil {
			pruned = append(pruned, ins)
		}
	}
	ftac.instrs = pruned
}

func (ftac *FunctionTAC) removeRedundantInstrs() {
	pruned := make([]ThreeAddressInstr, 0)
	for i, ins := range ftac.instrs {
		if v, ok := ins.(*LabelPlaceholder); ok {
			if i == len(ftac.instrs) {
				pruned = append(pruned, ins)
				continue
			}
			labs := ftac.instrs[i+1].Labels()
			ftac.instrs[i+1].setLabels(append(labs, v.labels...))
		} else {
			pruned = append(pruned, ins)
		}
	}
	ftac.instrs = pruned

}
