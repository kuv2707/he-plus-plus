package x64sysvasm

import (
	"bytes"
	"encoding/binary"
	"fmt"
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
	return fmt.Sprintf("(%s, %d}", arg.locType.String(), arg.ival)
}

var NOWHERE = TACOpArg{Null, 0}

func (t TACLocationType) String() string {
	switch t {
	case VirtualRegister:
		return "VReg"
	case Immediate:
		return "Imm"
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
	nBytes int
}

type FunctionTAC struct {
	fname             string
	regCnt            int64
	instrs            []ThreeAddressInstr
	symTable          map[string]*utils.Stack[SymDef] // todo: typeAsm stores the assembly-level type for instr selection, like int or float ...
	dataSectionAllocs []DataSectionAllocEntry
}

func (ft *FunctionTAC) assignVirtualReg(name string) VirtualRegisterNumber {
	ft.regCnt++
	if name == "" {
		name = fmt.Sprint("r", ft.regCnt)
	}
	st, ok := ft.symTable[name]
	if !ok {
		st = utils.MakeStack[SymDef]()
		ft.symTable[name] = st
	}
	st.Push(SymDef{ft.regCnt, 0})
	fmt.Println("assn ", name, ft.regCnt)
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
			tb := FunctionTAC{
				fname:             v.Name,
				regCnt:            0, // first reg gets 1 since inc before assn
				instrs:            nil,
				symTable:          make(map[string]*utils.Stack[SymDef]),
				dataSectionAllocs: make([]DataSectionAllocEntry, 0)}
			// todo: load func args
			tb.genScopeTAC(v.Scope)
			fmt.Println("TAC for func", v.Name)
			for _, k := range tb.instrs {
				fmt.Println(k)
			}
			fmt.Println()
			fmt.Println(".data alloc entries")
			for _, k := range tb.dataSectionAllocs {
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
		switch v := k.(type) {
		case *node_types.VariableDeclarationNode:
			{
				for _, d := range v.Declarations {
					dcl := d.(*node_types.InfixOperatorNode)
					vname := dcl.Left.(*node_types.IdentifierNode).Name()
					vreg := ftac.assignVirtualReg(vname)
					ret := ftac.genExprTAC(dcl.Right)
					if ret != NOWHERE.ival {
						ftac.emitInstr(&UnaryOpInstr{assTo: TACOpArg{VirtualRegister, vreg}, op: "", arg1: TACOpArg{VirtualRegister, ret}})
					} else {
						// todo: set to 0 value
					}
				}
			}
		case *node_types.IfNode:
			{
				ifEnd := fmt.Sprintf("cond_%d_end", v.Seq)
				for i, branch := range v.Branches {
					if i == len(v.Branches)-1 {
						ftac.genScopeTAC(branch.Scope)
						continue
					}
					ftac.genExprTAC(branch.Condition)
					lastInstr := ftac.instrs[len(ftac.instrs)-1]
					condInstr := CJumpInstr{}
					switch vv := lastInstr.(type) {
					case *BinaryOpInstr:
						{
							condInstr.op = vv.op
							condInstr.argL = vv.arg1
							condInstr.argR = vv.arg2
							condInstr.jmpToLabel = fmt.Sprintf("cond_%d_brch_%d", v.Seq, i+1) // jmp to next condn
						}
					case *UnaryOpInstr:
						{
							// todo: not op
							panic("! operator not impl in if-node tac gen")
						}
					}
					ftac.instrs[len(ftac.instrs)-1] = &condInstr
					ftac.genScopeTAC(branch.Scope)
					ftac.emitInstr(&JumpInstr{jmpToLabel: ifEnd})

				}
			}

		case *node_types.ScopeNode:
			{
				ftac.genScopeTAC(v)
			}
		}
	}
}

func (ftac *FunctionTAC) genExprTAC(n node_types.TreeNode) VirtualRegisterNumber {

	vreg := ftac.assignVirtualReg("")
	switch v := n.(type) {
	case *node_types.NumberNode:
		{
			var num int64
			// var f float32
			binary.Read(bytes.NewReader(v.RawNumBytes), binary.BigEndian, &num)
			fmt.Println("--->", num)
			// todo: currently we read floats as ints too - impl type tracking of registers (float or int)
			ftac.emitInstr(&UnaryOpInstr{assTo: TACOpArg{VirtualRegister, vreg}, op: "", arg1: TACOpArg{Immediate, num}})
			return vreg
		}
	case *node_types.InfixOperatorNode:
		{
			left := ftac.genExprTAC(v.Left)
			// todo: short circuiting for || and &&
			right := ftac.genExprTAC(v.Right)
			ftac.emitInstr(&BinaryOpInstr{
				assnTo: TACOpArg{VirtualRegister, vreg},
				op:     v.Op,
				arg1:   TACOpArg{VirtualRegister, left},
				arg2:   TACOpArg{VirtualRegister, right}})
			return vreg
		}
	case *node_types.IdentifierNode:
		{
			symStk, ok := ftac.symTable[v.Name()]
			if !ok {
				// treat this name as a label
				ftac.emitInstr(&LoadLabelInstr{
					loadeeLabel: v.Name(),
					to:          TACOpArg{VirtualRegister, vreg},
				})
				return vreg
			}
			e, ok := symStk.Peek()
			if !ok {
				panic("Forgot to push reg no for sym!")
			}
			return e.reg
		}
	case *node_types.ArrayDeclarationNode: // todo: WIP
		{
			// we can use an instr like : r1 = ALLOC <sizeofarray_bytes>
			// whether to use data section or stack for allocating space is a later decision

			nelems := len(v.Elems)
			label := fmt.Sprintf("%s_%d", ftac.fname, vreg) // todo: maybe a more systematic naming
			ftac.dataSectionAllocs = append(ftac.dataSectionAllocs, DataSectionAllocEntry{label, nelems * v.DataT.Size()})

		}
	case *node_types.FuncCallNode:
		{
			callAddr := ftac.genExprTAC(v.Callee)

			for _, arg := range v.Args {
				areg := ftac.genExprTAC(arg)
				ftac.emitInstr(&ParamInstr{arg: TACOpArg{locType: VirtualRegister, ival: areg}})
			}
			ftac.emitInstr(&CallInstr{calleeAddr: TACOpArg{locType: VirtualRegister, ival: callAddr}})
		}
	default:
		panic(fmt.Sprintf("genxprTAC Not implemented for %T", v))
	}
	return NOWHERE.ival
}
