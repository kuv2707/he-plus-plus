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

type TACLocationType byte
type VirtualRegisterNumber = int64

const (
	VirtualRegister TACLocationType = 0
	Immediate       TACLocationType = 3
	Null            TACLocationType = 7
)

type TACOpArg struct {
	locType TACLocationType
	ival    int64
}

func (arg TACOpArg) String() string {
	return fmt.Sprintf("%s%d", utils.Yellow(arg.locType.String()), arg.ival)
}

var NOWHERE = TACOpArg{Null, 0}

func (t TACLocationType) String() string {
	switch t {
	case VirtualRegister:
		return "R"
	case Immediate:
		return "#"
	}
	panic("Undefined TACLocationType")
}

type ThreeAddressOpCode string

const (
	LOAD  ThreeAddressOpCode = "LOAD"
	STORE ThreeAddressOpCode = "STORE"
	ADD   ThreeAddressOpCode = "ADD"
	SUB   ThreeAddressOpCode = "SUB"
	MUL   ThreeAddressOpCode = "MUL"
	DIV   ThreeAddressOpCode = "DIV"
	JMP   ThreeAddressOpCode = "JMP"
	PARAM ThreeAddressOpCode = "PARAM"
	CALL  ThreeAddressOpCode = "CALL"
)

type SymDef = struct {
	reg     VirtualRegisterNumber
	typeAsm int
}

type DataSectionAllocEntry struct {
	label  string
	nBytes TACOpArg
}

type FunctionTAC struct {
	fname             string
	regCnt            int64
	instrs            []ThreeAddressInstr
	symTable          map[string]SymDef // todo: typeAsm stores the assembly-level type for instr selection, like int or float ...
	dataSectionAllocs []DataSectionAllocEntry
}

func (ft *FunctionTAC) assignVirtualReg(name string) VirtualRegisterNumber {
	ft.regCnt++
	if name == "" {
		name = fmt.Sprint("r", ft.regCnt)
	}
	ft.symTable[name] = SymDef{ft.regCnt, 0}

	return ft.regCnt
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
				symTable:          make(map[string]SymDef),
				dataSectionAllocs: make([]DataSectionAllocEntry, 0)}
			// todo: load func args
			ftac.genScopeTAC(v.Scope)
			ftac.condenseTAC()
			fmt.Println("TAC for func", v.Name)
			for _, k := range ftac.instrs {
				fmt.Println(k)
			}
			fmt.Println()
			fmt.Println(".data alloc entries")
			for _, k := range ftac.dataSectionAllocs {
				fmt.Println(k)
			}
		default:
			panic(fmt.Sprintf("%T not supported for asm gen yet", ch))
		}
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
				ftac.genExprTAC(dcl.Right, ret)
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
				ftac.genExprTAC(branch.Condition, NOWHERE.ival)
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
			loopLabel := fmt.Sprintf("loop_%d", v.Seq)
			loopEndLabel := loopLabel + "_end"
			ftac.genNodeTAC(v.Initializer)

			ftac.genExprTAC(v.Condition, NOWHERE.ival)
			lastInstr := ftac.instrs[len(ftac.instrs)-1]
			condInstr := genCondInstr(lastInstr, loopEndLabel)
			ftac.instrs[len(ftac.instrs)-1] = &condInstr

			ftac.emitInstr(placeholderWithLabels(loopLabel))
			ftac.genScopeTAC(v.Scope)
			ftac.genNodeTAC(v.Updater)
			ftac.emitInstr(&JumpInstr{jmpToLabel: loopLabel})
			ftac.emitInstr(placeholderWithLabels(loopEndLabel))
		}
	case *node_types.InfixOperatorNode:
		{
			ftac.genExprTAC(v, NOWHERE.ival)
		}
	case *node_types.FuncCallNode:
		{
			ftac.genExprTAC(v, NOWHERE.ival)
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
		ftac.genExprTAC(v, NOWHERE.ival)
	}
}

func (ftac *FunctionTAC) genExprTAC(n node_types.TreeNode, vreg VirtualRegisterNumber) VirtualRegisterNumber {
	if vreg == NOWHERE.ival {
		vreg = ftac.assignVirtualReg("")
	}
	switch v := n.(type) {
	case *node_types.NumberNode:
		{
			var num int64
			// var f float32
			binary.Read(bytes.NewReader(v.RawNumBytes), binary.BigEndian, &num)
			if reg, exists := ftac.symTable[fmt.Sprint(num)]; exists {
				return reg.reg
			}
			numReg := ftac.assignVirtualReg(fmt.Sprint(num))
			// todo: currently we read floats as ints too - impl type tracking of registers (float or int)
			ftac.emitInstr(&UnaryOpInstr{assTo: TACOpArg{VirtualRegister, numReg}, op: "", arg1: TACOpArg{Immediate, num}})
			ftac.emitInstr(&UnaryOpInstr{assTo: TACOpArg{VirtualRegister, vreg}, op: "", arg1: TACOpArg{VirtualRegister, numReg}})
			return vreg
		}
	case *node_types.InfixOperatorNode:
		{
			switch v.Op {
			case lexer.ASSN:
				{
					right := ftac.genExprTAC(v.Right, NOWHERE.ival)
					if vl, ok := v.Left.(*node_types.ArrIndNode); ok {
						indexedElemAddrReg, _ := ftac.getArrIndPointingAt(vl, vreg)
						ftac.emitInstr(&MemStoreInstr{storeAt: TACOpArg{VirtualRegister, indexedElemAddrReg}, storeWhat: TACOpArg{VirtualRegister, right}, numBytes: 4}) // todo: compute reqd size and remove hardcoding
					} else {
						left := ftac.genExprTAC(v.Left, NOWHERE.ival)
						ftac.emitInstr(&UnaryOpInstr{assTo: TACOpArg{VirtualRegister, left}, op: "", arg1: TACOpArg{VirtualRegister, right}})
						return left
					}
				}
			default:
				{
					left := ftac.genExprTAC(v.Left, NOWHERE.ival)
					// todo: short circuiting for || and &&
					right := ftac.genExprTAC(v.Right, NOWHERE.ival)
					ftac.emitInstr(&BinaryOpInstr{
						assnTo: TACOpArg{VirtualRegister, vreg},
						op:     v.Op,
						arg1:   TACOpArg{VirtualRegister, left},
						arg2:   TACOpArg{VirtualRegister, right},
					})
				}
			}

			return vreg
		}
	case *node_types.PrePostOperatorNode:
		{
			operand := ftac.genExprTAC(v.Operand, vreg)
			switch v.Op {
			case lexer.SUB:
				{
					ftac.emitInstr(&UnaryOpInstr{assTo: TACOpArg{VirtualRegister, vreg}, op: v.Op, arg1: TACOpArg{VirtualRegister, operand}})
				}
			default:
				panic("Not impl for prepost op" + v.Op)
			}
			return vreg
		}
	case *node_types.IdentifierNode:
		{
			sym, ok := ftac.symTable[v.Name()]
			if !ok {
				// treat this name as a label
				ftac.emitInstr(&LoadLabelInstr{
					loadeeLabel: v.Name(),
					to:          TACOpArg{VirtualRegister, vreg},
				})
				return vreg
			}

			return sym.reg
		}
	case *node_types.ArrayDeclarationNode: // todo: WIP
		{
			// we can use an instr like : r1 = ALLOC <sizeofarray_bytes> for arrays and structs
			// whether to use data section or stack for allocating space is a later decision

			arrSizeReg := ftac.genExprTAC(v.SizeNode, NOWHERE.ival)
			elemSizeBytes := v.DataT.Size()
			reqBytesReg := ftac.assignVirtualReg("")
			ftac.emitInstr(&BinaryOpInstr{
				assnTo: TACOpArg{
					VirtualRegister,
					reqBytesReg,
				},
				op: lexer.MUL,
				arg1: TACOpArg{
					VirtualRegister,
					arrSizeReg,
				},
				arg2: TACOpArg{Immediate, int64(elemSizeBytes)}})
			// label := fmt.Sprintf("%s_%d", ftac.fname, vreg) // todo: maybe a more systematic naming
			// todo: data section allocation should only be in case of const array, else alloc on stack
			// ftac.dataSectionAllocs = append(ftac.dataSectionAllocs, DataSectionAllocEntry{label, TACOpArg{VirtualRegister, reqBytesReg}})
			arrPtr := vreg
			ftac.emitInstr(&AllocInstr{allocType: STACK_ALLOC,
				sizeReg:    TACOpArg{VirtualRegister, reqBytesReg},
				ptrToAlloc: TACOpArg{VirtualRegister, arrPtr},
			})
			for i, entry := range v.Elems {
				byteOffset := elemSizeBytes * i
				memLocReg := ftac.assignVirtualReg("")
				ftac.emitInstr(&BinaryOpInstr{op: lexer.ADD, assnTo: TACOpArg{VirtualRegister, memLocReg}, arg1: TACOpArg{VirtualRegister, arrPtr}, arg2: TACOpArg{Immediate, int64(byteOffset)}})
				storeVal := ftac.genExprTAC(entry, NOWHERE.ival)
				ftac.emitInstr(&MemStoreInstr{
					storeAt:   TACOpArg{VirtualRegister, memLocReg},
					storeWhat: TACOpArg{VirtualRegister, storeVal},
					numBytes:  8,
				}) // todo: track size of var
			}
			return arrPtr
		}
	case *node_types.ArrIndNode:
		{
			indexedElemAddrReg, sizeBytes := ftac.getArrIndPointingAt(v, vreg)
			ftac.emitInstr(&MemLoadInstr{loadFrom: TACOpArg{VirtualRegister, indexedElemAddrReg}, storeAt: TACOpArg{VirtualRegister, indexedElemAddrReg}, numBytes: sizeBytes})
			return indexedElemAddrReg
		}
	case *node_types.FuncCallNode:
		{
			callAddr := ftac.genExprTAC(v.Callee, NOWHERE.ival)

			for _, arg := range v.Args {
				areg := ftac.genExprTAC(arg, NOWHERE.ival)
				ftac.emitInstr(&ParamInstr{arg: TACOpArg{locType: VirtualRegister, ival: areg}})
			}
			ftac.emitInstr(&CallInstr{retReg: TACOpArg{locType: VirtualRegister, ival: vreg}, calleeAddr: TACOpArg{locType: VirtualRegister, ival: callAddr}})
			return vreg
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
	return NOWHERE.ival
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

func (ftac *FunctionTAC) getArrIndPointingAt(v *node_types.ArrIndNode, vreg VirtualRegisterNumber) (VirtualRegisterNumber, int) {
	arrBaseAddrReg := ftac.genExprTAC(v.ArrProvider, vreg)
	indVarReg := ftac.genExprTAC(v.Indexer, NOWHERE.ival)
	indReg := ftac.assignVirtualReg("")
	ftac.emitInstr(&UnaryOpInstr{assTo: TACOpArg{VirtualRegister, indReg}, op: "", arg1: TACOpArg{VirtualRegister, indVarReg}})
	sizeBytes := v.DataType.Size()

	ftac.emitInstr(&BinaryOpInstr{
		assnTo: TACOpArg{VirtualRegister, indReg},
		op:     lexer.MUL,
		arg1:   TACOpArg{VirtualRegister, indReg},
		arg2:   TACOpArg{Immediate, int64(sizeBytes)},
	})
	ftac.emitInstr(&BinaryOpInstr{
		assnTo: TACOpArg{VirtualRegister, indReg},
		op:     lexer.ADD,
		arg1:   TACOpArg{VirtualRegister, arrBaseAddrReg},
		arg2:   TACOpArg{VirtualRegister, indReg},
	})
	return indReg, sizeBytes
}
