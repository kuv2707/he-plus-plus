package asm_gen

import (
	"fmt"
	"he++/lexer"
	"he++/tac"
	"strings"
)

type x86_64Instr struct {
	instrName string
	params    []string
	labels    []string
}

func (ins x86_64Instr) String() string {
	var sb strings.Builder
	sb.WriteString(strings.Join(ins.labels, ":\n"))
	if len(ins.labels) > 0 {
		sb.WriteString(":\n")
	}
	sb.WriteString(fmt.Sprintf("%s %s", ins.instrName, strings.Join(ins.params, ", ")))
	return sb.String()
}

const (
	MOV  = "mov"
	ADD  = "add"
	SUB  = "sub"
	IMUL = "imul"
	IDIV = "idiv"
	NEG  = "neg"
	CMP  = "cmp"

	JMP = "jmp"
	JE  = "je"
	JNE = "jne"
	JL  = "jl"
	JLE = "jle"
	JG  = "jg"
	JGE = "jge"

	AND = "and"
	OR  = "or"
	XOR = "xor"
	SHL = "shl"
	SHR = "shr"

	CALL = "call"
	RET  = "ret"
)

var compOpsName = map[string]string{
	lexer.LESS:    JL,
	lexer.LEQ:     JLE,
	lexer.GREATER: JG,
	lexer.GEQ:     JGE,
	lexer.EQ:      JE,
	lexer.NEQ:     JNE,
}

func opInstrName(op tac.TACOperator) string {
	switch op {

	// Arithmetic
	case "+":
		return ADD
	case "-":
		return SUB
	case "*":
		return IMUL
	case "/":
		return IDIV
	case "neg":
		return NEG

	// Comparisons (usually paired with CMP + jump)
	case "==", "!=", "<", "<=", ">", ">=":
		return CMP

	// Logical / bitwise
	case "&":
		return AND
	case "|":
		return OR
	case "^":
		return XOR
	case "<<":
		return SHL
	case ">>":
		return SHR

	// Calls
	case "call":
		return CALL
	case "return":
		return RET
	}

	panic("unsupported TAC operator: " + string(op))
}

func OppositeCompOp(s tac.TACOperator) tac.TACOperator {
	type cop = tac.TACOperator
	switch s {
	case cop(lexer.LESS):
		return cop(lexer.GEQ)
	case cop(lexer.LEQ):
		return cop(lexer.GREATER)
	case cop(lexer.GREATER):
		return cop(lexer.LEQ)
	case cop(lexer.GEQ):
		return cop(lexer.LESS)
	case cop(lexer.EQ):
		return cop(lexer.NEQ)
	case cop(lexer.NEQ):
		return cop(lexer.EQ)
	default:
		panic("unsupported comparison operator: " + string(s))
	}
}
