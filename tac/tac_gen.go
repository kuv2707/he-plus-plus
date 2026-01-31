package tac

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"he++/lexer"
	"he++/parser/node_types"
	"he++/utils"
)

type DataSectionAllocEntry struct {
	label  string
	nBytes TACOpArg
}

type FunctionTAC struct {
	fname             string
	regCnt            int64
	instrs            []ThreeAddressInstr
	nameToReg         map[string]VirtualRegisterNumber
	dataSectionAllocs []DataSectionAllocEntry
	allocCnt          int
	ctx               TACContext
}

func (ft *FunctionTAC) Instrs() []ThreeAddressInstr {
	return ft.instrs
}

func (ft *FunctionTAC) assignVirtualReg(name string) VirtualRegisterNumber {
	ft.regCnt++
	ft.nameToReg[name] = ft.regCnt
	return ft.regCnt
}

func (ftac *FunctionTAC) RegLifetimes() map[VirtualRegisterNumber]Life {
	return ftac.ctx.regLifetimes
}

type TACHandler struct {
	ast       *node_types.SourceFileNode
	TacBlocks map[string]*FunctionTAC
}

func NewTACGen(ast *node_types.SourceFileNode) *TACHandler {
	// assumes the AST is well-shaped
	return &TACHandler{ast, make(map[string]*FunctionTAC)}
}

func (ag *TACHandler) GenerateTac() {
	for _, ch := range ag.ast.Children {
		switch v := ch.(type) {
		case *node_types.FuncNode:
			ftac := FunctionTAC{
				fname:             v.Name,
				regCnt:            0, // first reg gets 1 since inc before assn
				instrs:            nil,
				nameToReg:         make(map[string]VirtualRegisterNumber),
				dataSectionAllocs: make([]DataSectionAllocEntry, 0)}
			// todo: load func args
			ftac.loadFuncArgs(v.ArgList)
			ftac.genScopeTAC(v.Scope)

			ftac.Optimize()

			ftac.printInstrs()

			// fmt.Println("\n.data alloc entries")
			// for i, k := range ftac.dataSectionAllocs {
			// 	fmt.Printf("%d) %v", i, k)
			// }
			ag.TacBlocks[ftac.fname] = &ftac
		default:
			panic(fmt.Sprintf("%T not supported for asm gen yet", ch))
		}
	}
}

func (ftac *FunctionTAC) printInstrs() {
	fmt.Printf("TAC for func %s\n", utils.BoldGreen(ftac.fname))
	for i, k := range ftac.instrs {
		fmt.Printf("%d) %s\n", i, k)
	}
}

func (ag *FunctionTAC) emitInstr(tai ThreeAddressInstr) {
	ag.instrs = append(ag.instrs, tai)
}

func (ftac *FunctionTAC) loadFuncArgs(arglist []node_types.FuncArg) {
	for i := range arglist {
		argReg := ftac.assignVirtualReg(arglist[i].Name)
		dc := dataCategoryForType(arglist[i].DataT)

		ftac.emitInstr(&FuncArgRecvInstr{argNo: i, recvInto: &VRegArg{argReg, dc}})
	}
}

func (ftac *FunctionTAC) genScopeTAC(scp *node_types.ScopeNode) {
	for _, k := range scp.Children {
		ftac.genNodeTAC(k)
	}
}

func (ftac *FunctionTAC) genNodeTAC(k node_types.TreeNode) {
	switch v := k.(type) {
	case *node_types.VariableDeclarationNode:
		{
			for _, d := range v.Declarations {
				dcl := d.(*node_types.InfixOperatorNode)
				vname := dcl.Left.(*node_types.IdentifierNode).Name()
				ret := ftac.assignVirtualReg(vname)
				dc := dataCategoryForType(v.DataT)
				r := ftac.genExprTAC(dcl.Right)
				ftac.emitInstr(&AssignInstr{assnTo: &VRegArg{ret, dc}, arg: r})
			}
		}
	case *node_types.IfNode:
		{
			ifEnd := fmt.Sprintf("cond_%d_brch_%d", v.Seq, len(v.Branches))
			for i, branch := range v.Branches {
				ftac.emitInstr(placeholderWithLabels(fmt.Sprintf("cond_%d_brch_%d", v.Seq, i)))
				if bn, ok := branch.Condition.(*node_types.BooleanNode); ok {
					if bn.BoolVal {
						ftac.genScopeTAC(branch.Scope)
						break // since the next branches are dead code
					} else {
						// skip this dead code
					}
				}
				ftac.genExprTAC(branch.Condition)
				lastInstr := ftac.instrs[len(ftac.instrs)-1]
				condInstr := genCondInstr(lastInstr, fmt.Sprintf("cond_%d_brch_%d", v.Seq, i+1)) // jmp to next condn
				ftac.instrs[len(ftac.instrs)-1] = &condInstr
				ftac.genScopeTAC(branch.Scope)
				if i < len(v.Branches)-1 {
					ftac.emitInstr(&JumpInstr{JmpToLabel: ifEnd})
				}
			}
			ftac.emitInstr(placeholderWithLabels(ifEnd))
		}
	case *node_types.LoopNode:
		{
			// if any vreg is assigned inside a loop scope, it should not be folded.
			loopStartLabel := fmt.Sprintf("%s%d", LOOP_START_PREFIX, v.Seq)
			loopEndLabel := fmt.Sprintf("%s%d", LOOP_END_PREFIX, v.Seq)
			ftac.genNodeTAC(v.Initializer)

			ftac.genExprTAC(v.Condition)
			lastInstr := ftac.instrs[len(ftac.instrs)-1]
			condInstr := genCondInstr(lastInstr, loopEndLabel)
			condInstr.setLabels([]string{loopStartLabel})
			ftac.instrs[len(ftac.instrs)-1] = &LoopBoundary{
				loopNo:   v.Seq,
				StartEnd: true,
			}

			ftac.emitInstr(&condInstr)
			ftac.genScopeTAC(v.Scope)
			ftac.genNodeTAC(v.Updater)
			ftac.emitInstr(&JumpInstr{JmpToLabel: loopStartLabel})
			ftac.emitInstr(&LoopBoundary{
				loopNo:       v.Seq,
				StartEnd:     false,
				TACBaseInstr: TACBaseInstr{labels: []string{loopEndLabel}},
			})
		}
	case *node_types.ScopeNode:
		{
			ftac.genScopeTAC(v)
		}
	case *node_types.ReturnNode:
		{
			rValReg := ftac.genExprTAC(v.Value)
			ftac.emitInstr(&FuncRetInstr{retReg: rValReg})
		}
	case *node_types.EmptyPlaceholderNode:
		{
			// no hacer nada
		}
	default:
		ftac.genExprTAC(v)
	}
}

func (ftac *FunctionTAC) genExprTAC(n node_types.TreeNode) TACOpArg {
	switch v := n.(type) {
	case *node_types.NumberNode:
		{
			// todo: 16 vs 32 vs 64 bit distinction.
			if v.NumType == node_types.INT_NUM {
				var num int64
				binary.Read(bytes.NewReader(v.RawNumBytes), binary.BigEndian, &num)
				return &ImmIntArg{num, I64}
			} else {
				var num float64
				binary.Read(bytes.NewReader(v.RawNumBytes), binary.BigEndian, &num)
				return &ImmFloatArg{num, F64}
			}
		}
	case *node_types.InfixOperatorNode:
		{
			switch v.Op {
			case lexer.ASSN:
				{
					right := ftac.genExprTAC(v.Right)
					if vl, ok := v.Left.(*node_types.ArrIndNode); ok {
						indexedElemAddrArg, elemType := ftac.getMemLocationPointingAt(vl)
						ftac.emitInstr(&MemStoreInstr{
							StoreAt:   indexedElemAddrArg,
							StoreWhat: right,
							NumBytes:  elemType.Size()},
						)
						return right // todo: decide semantics of a <binop> b = c
					} else {
						left := ftac.genExprTAC(v.Left)
						ftac.emitInstr(&AssignInstr{assnTo: left, arg: right})
						return left
					}
				}
			default:
				{
					retArg := &VRegArg{ftac.assignVirtualReg(""), dataCategoryForType(v.ResultDT)}

					// todo: short circuiting for || and &&
					left := ftac.genExprTAC(v.Left)
					right := ftac.genExprTAC(v.Right)
					ftac.emitInstr(&BinaryOpInstr{
						assnTo: retArg,
						op:     TACOperator(v.Op),
						arg1:   left,
						arg2:   right,
					})
					return retArg
				}
			}

		}
	case *node_types.PrePostOperatorNode:
		{
			retArg := &VRegArg{ftac.assignVirtualReg(""), dataCategoryForType(v.ResultDT)}
			operand := ftac.genExprTAC(v.Operand)
			switch v.Op {
			case lexer.SUB:
				{
					ftac.emitInstr(&UnaryOpInstr{assnTo: retArg, op: v.Op, arg1: operand})
				}
			default:
				panic("Not impl for prepost op" + v.Op)
			}
			return retArg
		}
	case *node_types.IdentifierNode:
		{
			reg, ok := ftac.nameToReg[v.Name()]
			dc := dataCategoryForType(v.DataT)
			if !ok {
				// treat this name as a label
				retArg := &VRegArg{ftac.assignVirtualReg(v.Name()), dc}

				ftac.emitInstr(&LoadLabelInstr{
					loadeeLabel: v.Name(),
					to:          retArg,
				})
				return retArg
			}

			return &VRegArg{reg, dc}
		}
	case *node_types.ArrayDeclarationNode: // todo: WIP
		{
			// we can use an instr like : r1 = ALLOC <sizeofarray_bytes> for arrays and structs
			// whether to use data section or stack for allocating space is a later decision

			arrSizeArg := ftac.genExprTAC(v.SizeNode)
			elemSizeBytes := v.DataT.Size()
			reqBytesArg := &VRegArg{ftac.assignVirtualReg(""), I64}

			ftac.emitInstr(&BinaryOpInstr{
				assnTo: reqBytesArg,
				op:     TACOperator(lexer.MUL),
				arg1:   arrSizeArg,
				arg2:   &ImmIntArg{int64(elemSizeBytes), I64}})

			arrPtr := &VRegArg{ftac.assignVirtualReg(""), PTR}

			ftac.emitInstr(&AllocInstr{AllocType: STACK_ALLOC,
				SizeReg:    reqBytesArg,
				PtrToAlloc: arrPtr,
				AllocNo:    ftac.allocCnt,
			})
			ftac.allocCnt++

			memLocArg := &VRegArg{ftac.assignVirtualReg(""), PTR}
			ftac.emitInstr(&AssignInstr{assnTo: memLocArg, arg: arrPtr})
			for _, entry := range v.Elems {
				ftac.emitInstr(&BinaryOpInstr{op: TACOperator(lexer.ADD),
					assnTo: memLocArg,
					arg1:   memLocArg,
					arg2:   &ImmIntArg{int64(elemSizeBytes), I64},
				})
				storeVal := ftac.genExprTAC(entry)
				ftac.emitInstr(&MemStoreInstr{
					StoreAt:   memLocArg,
					StoreWhat: storeVal,
					NumBytes:  v.DataT.Size(),
				})
			}
			return arrPtr
		}
	case *node_types.ArrIndNode:
		{
			bytePosArg, indexedElemType := ftac.getMemLocationPointingAt(v)
			indexedElemArg := &VRegArg{ftac.assignVirtualReg(""), dataCategoryForType(indexedElemType)}

			ftac.emitInstr(&MemLoadInstr{
				LoadFrom: bytePosArg,
				StoreAt:  indexedElemArg,
				NumBytes: indexedElemType.Size(),
			})
			return indexedElemArg
		}
	case *node_types.FuncCallNode:
		{
			for _, arg := range v.Args {
				areg := ftac.genExprTAC(arg)
				ftac.emitInstr(&ParamInstr{arg: areg})
			}

			retArg := &VRegArg{ftac.assignVirtualReg(""), dataCategoryForType(v.CalleeT.ReturnType)}
			callAddr := ftac.genExprTAC(v.Callee)
			// if callee is a static label, have a separate instr for it instead of assigning
			// label to vreg and then calling the vreg
			ftac.emitInstr(&CallInstr{retReg: retArg, calleeAddr: callAddr})
			return retArg
		}
	case *node_types.EmptyPlaceholderNode:
		{
			// nada que hacer
		}
	default:
		p := utils.ASTPrinter{}
		v.String(&p)
		panic(fmt.Sprintf("TAC gen Not implemented for %s", p.Builder.String()))
	}
	return &NULLOpArg{}
}

func genCondInstr(fromInstr ThreeAddressInstr, jmpToLabel string) CJumpInstr {
	condInstr := CJumpInstr{}
	switch vv := fromInstr.(type) {
	case *BinaryOpInstr:
		{
			condInstr.Op = vv.op
			condInstr.argL = vv.arg1
			condInstr.argR = vv.arg2
			condInstr.JmpToLabel = jmpToLabel
		}
	case *UnaryOpInstr:
		{
			// todo: not op
			panic("! operator not impl in if-node tac gen")
		}
	}
	return condInstr
}

func (ftac *FunctionTAC) getMemLocationPointingAt(v *node_types.ArrIndNode) (TACOpArg, node_types.DataType) {
	arrBaseAddrArg := ftac.genExprTAC(v.ArrProvider)
	indVarArg := ftac.genExprTAC(v.Indexer)
	indArg := &VRegArg{ftac.assignVirtualReg(""), indVarArg.Category()}
	ftac.emitInstr(&AssignInstr{assnTo: indArg, arg: indVarArg})
	sizeBytes := v.DataType.Size()
	byteOffsetArg := &VRegArg{ftac.assignVirtualReg(""), indArg.Category()}
	ftac.emitInstr(&BinaryOpInstr{
		assnTo: byteOffsetArg,
		op:     TACOperator(lexer.MUL),
		arg1:   indArg,
		arg2:   &ImmIntArg{int64(sizeBytes), I64},
	})

	bytePosArg := &VRegArg{ftac.assignVirtualReg(""), PTR}
	ftac.emitInstr(&BinaryOpInstr{
		assnTo: bytePosArg,
		op:     TACOperator(lexer.ADD),
		arg1:   arrBaseAddrArg,
		arg2:   &VRegArg{byteOffsetArg.RegNo, PTR},
	})
	return bytePosArg, v.DataType
}
