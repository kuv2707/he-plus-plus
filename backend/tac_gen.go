package backend

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"he++/lexer"
	"he++/parser/node_types"
	"he++/utils"
)

/**
Sample 3AC: Adding even elements in array:
s = 0
for i from 0 to n-1 {
    if A[i] % 2 == 0 {
        s += A[i]
    }
}
--
i = 0
s = 0
L1:
if i >= n goto L2
t1 = A[i]
t2 = t1 % 2
if t2 != 0 goto L3
s = s + t1
L3:
i = i + 1
goto L1
L2:
*/

type SymDef struct {
	reg VirtualRegisterNumber
	dt  DataCategory
}

type DataSectionAllocEntry struct {
	label  string
	nBytes TACOpArg
}

type FunctionTAC struct {
	fname             string
	regCnt            int64
	instrs            []ThreeAddressInstr
	symTable          map[VirtualRegisterNumber]SymDef // todo: typeAsm stores the assembly-level type for instr selection, like int or float ...
	nameToReg         map[string]VirtualRegisterNumber
	dataSectionAllocs []DataSectionAllocEntry
}

func (ft *FunctionTAC) assignVirtualReg(name string) VirtualRegisterNumber {
	ft.regCnt++
	ft.nameToReg[name] = ft.regCnt
	ft.symTable[ft.regCnt] = SymDef{ft.regCnt, VOID}
	return ft.regCnt
}

func (ftac *FunctionTAC) assignRegDataCategory(r VirtualRegisterNumber, dc DataCategory) {
	symDef := ftac.symTable[r]
	symDef.dt = dc
	ftac.symTable[r] = symDef
}

func (ftac *FunctionTAC) getDataRegCategory(r VirtualRegisterNumber) DataCategory {
	if symDef, exists := ftac.symTable[r]; exists {
		return symDef.dt
	}
	return VOID
}

type TACHandler struct {
	ast       *node_types.SourceFileNode
	tacBlocks map[string]FunctionTAC
}

func NewTACGen(ast *node_types.SourceFileNode) *TACHandler {
	// assumes the AST is well-shaped
	return &TACHandler{ast, make(map[string]FunctionTAC)}
}

func (ag *TACHandler) GenerateTac() {
	for _, ch := range ag.ast.Children {
		switch v := ch.(type) {
		case *node_types.FuncNode:
			ftac := FunctionTAC{
				fname:             v.Name,
				regCnt:            0, // first reg gets 1 since inc before assn
				instrs:            nil,
				symTable:          make(map[VirtualRegisterNumber]SymDef),
				nameToReg:         make(map[string]VirtualRegisterNumber),
				dataSectionAllocs: make([]DataSectionAllocEntry, 0)}
			// todo: load func args
			ftac.genScopeTAC(v.Scope)
			// constant and copy propag would cause some data flow chains to be dangling
			// which will then be eliminated by pruning, so we want purning to be after propag.
			ftac.printInstrs()
			ftac.Optimize()

			ftac.printInstrs()

			// fmt.Println("\n.data alloc entries")
			// for i, k := range ftac.dataSectionAllocs {
			// 	fmt.Printf("%d) %v", i, k)
			// }

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
				ftac.assignRegDataCategory(ret, dataCategoryForType(v.DataT))
				r := ftac.genExprTAC(dcl.Right)
				ftac.emitInstr(&AssignInstr{assnTo: &VRegArg{ret}, arg: &VRegArg{r}})
			}
		}
	case *node_types.IfNode:
		{
			ifEnd := fmt.Sprintf("cond_%d_end", v.Seq)
			for i, branch := range v.Branches {
				ftac.emitInstr(placeholderWithLabels(fmt.Sprintf("cond_%d_brch_%d", v.Seq, i)))
				if i == len(v.Branches)-1 {
					ftac.genScopeTAC(branch.Scope)
					continue
				}
				ftac.genExprTAC(branch.Condition)
				lastInstr := ftac.instrs[len(ftac.instrs)-1]
				condInstr := genCondInstr(lastInstr, fmt.Sprintf("cond_%d_brch_%d", v.Seq, i+1)) // jmp to next condn
				ftac.instrs[len(ftac.instrs)-1] = &condInstr
				ftac.genScopeTAC(branch.Scope)
				ftac.emitInstr(&JumpInstr{jmpToLabel: ifEnd})
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

			ftac.instrs[len(ftac.instrs)-1] = &LoopBoundary{
				loopNo:       v.Seq,
				startEnd:     true,
				TACBaseInstr: TACBaseInstr{labels: []string{loopStartLabel}},
			}

			ftac.emitInstr(&condInstr)
			ftac.genScopeTAC(v.Scope)
			ftac.genNodeTAC(v.Updater)
			ftac.emitInstr(&JumpInstr{jmpToLabel: loopStartLabel})
			ftac.emitInstr(&LoopBoundary{
				loopNo:       v.Seq,
				startEnd:     false,
				TACBaseInstr: TACBaseInstr{labels: []string{loopEndLabel}},
			})
		}
	case *node_types.InfixOperatorNode:
		{
			ftac.genExprTAC(v)
		}
	case *node_types.FuncCallNode:
		{
			ftac.genExprTAC(v)
		}
	case *node_types.ScopeNode:
		{
			ftac.genScopeTAC(v)
		}
	case *node_types.ReturnNode:
		{

		}
	case *node_types.EmptyPlaceholderNode:
		{
			// no hacer nada
		}
	default:
		ftac.genExprTAC(v)
	}
}

func (ftac *FunctionTAC) genExprTAC(n node_types.TreeNode) VirtualRegisterNumber {
	switch v := n.(type) {
	case *node_types.NumberNode:
		{
			// todo: 16 vs 32 vs 64 bit distinction.
			if v.NumType == node_types.INT_NUM {
				var num int64
				binary.Read(bytes.NewReader(v.RawNumBytes), binary.BigEndian, &num)
				if reg, exists := ftac.nameToReg[fmt.Sprint(num)]; exists {
					return reg
				}
				numReg := ftac.assignVirtualReg(fmt.Sprint(num))
				ftac.emitInstr(&AssignInstr{assnTo: &VRegArg{numReg}, arg: &ImmIntArg{num}})
				ftac.assignRegDataCategory(numReg, I64)
				return numReg
			} else {
				var num float64
				binary.Read(bytes.NewReader(v.RawNumBytes), binary.BigEndian, &num)
				if reg, exists := ftac.nameToReg[fmt.Sprint(num)]; exists {
					return reg
				}
				numReg := ftac.assignVirtualReg(fmt.Sprint(num))
				ftac.emitInstr(&AssignInstr{assnTo: &VRegArg{numReg}, arg: &ImmFloatArg{num}})
				ftac.assignRegDataCategory(numReg, I64)
				return numReg
			}
		}
	case *node_types.InfixOperatorNode:
		{
			switch v.Op {
			case lexer.ASSN:
				{
					right := ftac.genExprTAC(v.Right)
					if vl, ok := v.Left.(*node_types.ArrIndNode); ok {
						indexedElemAddrReg, _ := ftac.getArrIndPointingAt(vl)
						ftac.emitInstr(&MemStoreInstr{
							storeAt:   &VRegArg{indexedElemAddrReg},
							storeWhat: &VRegArg{right},
							numBytes:  vl.DataType.Size()}, // todo: compute reqd size and remove hardcoding
						)
						return right // todo: decide semantics of a <binop> b = c
					} else {
						left := ftac.genExprTAC(v.Left)
						ftac.emitInstr(&AssignInstr{assnTo: &VRegArg{left}, arg: &VRegArg{right}})
						return left
					}
				}
			default:
				{
					retReg := ftac.assignVirtualReg("")
					ftac.assignRegDataCategory(retReg, dataCategoryForType(v.ResultDT))
					left := ftac.genExprTAC(v.Left)
					// todo: short circuiting for || and &&
					right := ftac.genExprTAC(v.Right)
					ftac.emitInstr(&BinaryOpInstr{
						assnTo: &VRegArg{retReg},
						op:     v.Op,
						arg1:   &VRegArg{left},
						arg2:   &VRegArg{right},
					})
					return retReg
				}
			}

		}
	case *node_types.PrePostOperatorNode:
		{
			retReg := ftac.assignVirtualReg("")
			ftac.assignRegDataCategory(retReg, dataCategoryForType(v.ResultDT))
			operand := ftac.genExprTAC(v.Operand)
			switch v.Op {
			case lexer.SUB:
				{
					ftac.emitInstr(&UnaryOpInstr{assnTo: &VRegArg{retReg}, op: v.Op, arg1: &VRegArg{operand}})
				}
			default:
				panic("Not impl for prepost op" + v.Op)
			}
			return retReg
		}
	case *node_types.IdentifierNode:
		{
			reg, ok := ftac.nameToReg[v.Name()]
			if !ok {
				// treat this name as a label
				retReg := ftac.assignVirtualReg(v.Name()) // for subsequent use
				ftac.assignRegDataCategory(retReg, dataCategoryForType(v.DataT))
				ftac.emitInstr(&LoadLabelInstr{
					loadeeLabel: v.Name(),
					to:          &VRegArg{retReg},
				})
				return retReg
			}

			return reg
		}
	case *node_types.ArrayDeclarationNode: // todo: WIP
		{
			// we can use an instr like : r1 = ALLOC <sizeofarray_bytes> for arrays and structs
			// whether to use data section or stack for allocating space is a later decision

			arrSizeReg := ftac.genExprTAC(v.SizeNode)
			elemSizeBytes := v.DataT.Size()
			reqBytesReg := ftac.assignVirtualReg("")
			ftac.assignRegDataCategory(reqBytesReg, ftac.getDataRegCategory(arrSizeReg))
			ftac.emitInstr(&BinaryOpInstr{
				assnTo: &VRegArg{reqBytesReg},
				op:     lexer.MUL,
				arg1:   &VRegArg{arrSizeReg},
				arg2:   &ImmIntArg{int64(elemSizeBytes)}})

			// todo: data section allocation should only be in case of const array, else alloc on stack
			// ftac.dataSectionAllocs = append(ftac.dataSectionAllocs, DataSectionAllocEntry{label, TACOpArg{VirtualRegister, reqBytesReg}})
			arrPtr := ftac.assignVirtualReg("")
			ftac.assignRegDataCategory(arrPtr, PTR)
			ftac.emitInstr(&AllocInstr{allocType: STACK_ALLOC,
				sizeReg:    &VRegArg{reqBytesReg},
				ptrToAlloc: &VRegArg{arrPtr},
			})
			memLocReg := ftac.assignVirtualReg("")
			ftac.assignRegDataCategory(memLocReg, PTR)
			ftac.emitInstr(&AssignInstr{assnTo: &VRegArg{memLocReg}, arg: &VRegArg{arrPtr}})
			for _, entry := range v.Elems {
				ftac.emitInstr(&BinaryOpInstr{op: lexer.ADD,
					assnTo: &VRegArg{memLocReg},
					arg1:   &VRegArg{memLocReg},
					arg2:   &ImmIntArg{int64(elemSizeBytes)},
				})
				storeVal := ftac.genExprTAC(entry)
				ftac.emitInstr(&MemStoreInstr{
					storeAt:   &VRegArg{memLocReg},
					storeWhat: &VRegArg{storeVal},
					numBytes:  v.DataT.Size(),
				})
			}
			return arrPtr
		}
	case *node_types.ArrIndNode:
		{
			indexedElemAddrReg, indexedElemType := ftac.getArrIndPointingAt(v)
			ftac.assignRegDataCategory(indexedElemAddrReg, dataCategoryForType(indexedElemType))
			ftac.emitInstr(&MemLoadInstr{
				loadFrom: &VRegArg{indexedElemAddrReg},
				storeAt:  &VRegArg{indexedElemAddrReg},
				numBytes: indexedElemType.Size(),
			})
			return indexedElemAddrReg
		}
	case *node_types.FuncCallNode:
		{
			for _, arg := range v.Args {
				areg := ftac.genExprTAC(arg)
				ftac.emitInstr(&ParamInstr{arg: &VRegArg{areg}})
			}

			retReg := ftac.assignVirtualReg("")
			ftac.assignRegDataCategory(retReg, dataCategoryForType(v.CalleeT.ReturnType))
			callAddr := ftac.genExprTAC(v.Callee)
			// if callee is a static label, have a separate instr for it instead of assigning
			// label to vreg and then calling the vreg
			ftac.emitInstr(&CallInstr{retReg: &VRegArg{retReg}, calleeAddr: &VRegArg{callAddr}})
			return retReg
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
	return 0
}

func genCondInstr(fromInstr ThreeAddressInstr, jmpToLabel string) CJumpInstr {
	condInstr := CJumpInstr{}
	switch vv := fromInstr.(type) {
	case *BinaryOpInstr:
		{
			condInstr.op = vv.op
			condInstr.argL = vv.arg1
			condInstr.argR = vv.arg2
			condInstr.jmpToLabel = jmpToLabel
		}
	case *UnaryOpInstr:
		{
			// todo: not op
			panic("! operator not impl in if-node tac gen")
		}
	}
	return condInstr
}

func (ftac *FunctionTAC) getArrIndPointingAt(v *node_types.ArrIndNode) (VirtualRegisterNumber, node_types.DataType) {
	arrBaseAddrReg := ftac.genExprTAC(v.ArrProvider)
	indVarReg := ftac.genExprTAC(v.Indexer)
	indReg := ftac.assignVirtualReg("")
	ftac.emitInstr(&AssignInstr{assnTo: &VRegArg{indReg}, arg: &VRegArg{indVarReg}})
	sizeBytes := v.DataType.Size()

	ftac.emitInstr(&BinaryOpInstr{
		assnTo: &VRegArg{indReg},
		op:     lexer.MUL,
		arg1:   &VRegArg{indReg},
		arg2:   &ImmIntArg{int64(sizeBytes)},
	})
	ftac.emitInstr(&BinaryOpInstr{
		assnTo: &VRegArg{indReg},
		op:     lexer.ADD,
		arg1:   &VRegArg{arrBaseAddrReg},
		arg2:   &VRegArg{indReg},
	})
	ftac.assignRegDataCategory(indReg, PTR)
	return indReg, v.DataType
}
