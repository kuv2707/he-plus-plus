package interpreter

import "toylingo/parser"

//primitive
var TYPE_NUMBER = "number"
var TYPE_CHAR = "char"//unused for now
var TYPE_BOOLEAN = "bool"
var TYPE_POINTER = "pointer"

//composite
var TYPE_ARRAY = "array"

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
	scopeId     string
	scopeTyp    string
	scopeName   string
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
