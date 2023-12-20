package interpreter

import "toylingo/parser"

var TYPE_NUMBER = "number"
var TYPE_CHAR = "char"
var TYPE_BOOLEAN = "bool"

type Reason string

type Variable struct {
	pointer *Pointer
	vartype string
}

type Callable interface {
	call(ctx *scopeContext, args []parser.TreeNode) (Reason, *Variable)
}

// func(v Variable) call(ctx *scopeContext, args []parser.TreeNode) (Reason, *Variable) {
// }

type scopeContext struct {
	scopeType   string
	variables   map[string]Variable
	functions   map[string]parser.TreeNode
	returnValue *Variable
}

type Pointer struct {
	address int
	size    int
	scopeId string
	temp    bool
}