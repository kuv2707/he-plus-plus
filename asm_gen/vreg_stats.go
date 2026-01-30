package asm_gen

import (
	"fmt"
	"he++/tac"
	"he++/utils"
	"slices"
)

type Event struct {
	pos    int
	kind   bool
	vregNo tac.VirtualRegisterNumber
}

func (fasm *FunctionAsm) createVregMapping() {
	lives := fasm.ftac.RegLifetimes()
	events := getLifeEventList(lives)
	intRegListOrdered := utils.MakeStack(R10, R9, R8, RCX, RDX, RSI, RDI)
	// floatRegListOrdered := utils.MakeStack(XMM7, XMM6, XMM5, XMM4, XMM3, XMM2, XMM1, XMM0)
	regComp := func(a, b tac.VirtualRegisterNumber) bool {
		return vRegComparator(fasm.ftac, a, b)
	}
	vregsUsingRealRegs := utils.MakeHeap(regComp)
	// floatActiveUse := utils.MakeHeap(regComp)
	mapping := make(map[tac.VirtualRegisterNumber]Location)
	for _, evt := range events {
		if fasm.ftac.GetDataRegCategory(evt.vregNo).IsFloating() {
			// todo
		} else {
			if evt.kind {
				reg, ex := intRegListOrdered.Pop()
				if ex {
					mapping[evt.vregNo] = Location{reg: reg, offset: 0}
				} else {
					vregsUsingRealRegs.Push(evt.vregNo)
					leastImpReg, _ := vregsUsingRealRegs.Pop()
					// we spill this reg to stack
					if leastImpReg != evt.vregNo {
						mappedReg, ex := mapping[leastImpReg] // supposed to be reg loc, not stack loc
						if !ex {
							panic(fmt.Sprintf("Se esperaba un mapping para %d", leastImpReg))
						} else {
							mapping[evt.vregNo] = mappedReg
						}
					}
					// spill leastImpReg
					mapping[leastImpReg] = Location{reg: RSP, offset: fasm.stackFrameSize}
					fasm.stackFrameSize += fasm.ftac.GetDataRegCategory(leastImpReg).SizeBytes()
				}
			} else {
				// reclaim the reg
				mappedReg, ex := mapping[evt.vregNo]
				if ex {
					if mappedReg.offset == 0 {
						intRegListOrdered.Push(mappedReg.reg)
					} // else no reclamation
				} else {
					// logical error
					panic(fmt.Sprintf("Se esperaba un mapping para %d", evt.vregNo))
				}
			}
		}
	}
	fasm.VRegMapping = mapping
}

func getLifeEventList(lives map[tac.VirtualRegisterNumber]tac.Life) []Event {
	events := make([]Event, 0)
	for vregNo, life := range lives {
		events = append(events, Event{
			pos:    life.Start,
			kind:   true,
			vregNo: vregNo,
		})
		events = append(events, Event{
			pos:    life.End,
			kind:   false,
			vregNo: vregNo,
		})
	}

	slices.SortFunc(events, func(a, b Event) int {
		if a.pos == b.pos {
			// ensures life end is put after life begin if same pos
			if a.kind {
				return -1
			} else {
				return 1
			}
		}
		return a.pos - b.pos
	})
	return events
}
