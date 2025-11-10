package staticanalyzer

import (
	"he++/lexer"
	nodes "he++/parser/node_types"
	"he++/utils"
)

type DataTypeInfo interface {
}

type PrimitiveType struct {
	numBytes int
	typeName nodes.DataType
}

type Aliastype struct {
	AliasTo DataTypeInfo
}

func checkIntNode(n nodes.TreeNode) bool {
	nn, ok := n.(*nodes.NumberNode)
	if !ok {
		return false
	}
	// todo: remove hardcoded str
	return nn.NumType == "int"
}

func checkFloatNode(n nodes.TreeNode) bool {
	nn, ok := n.(*nodes.NumberNode)
	if !ok {
		return false
	}
	// todo: remove hardcoded str
	return nn.NumType == "float"
}

func checkBoolNode(n nodes.TreeNode) bool {
	_, ok := n.(*nodes.BooleanNode)
	return ok
}

var ERROR_TYPE = nodes.ErrorType{}
var INT_DATATYPE = nodes.NamedType{Name: lexer.INT}
var FLOAT_DATATYPE = nodes.NamedType{Name: lexer.FLOAT}
var BOOLEAN_DATATYPE = nodes.NamedType{Name: lexer.BOOLEAN}

func addInbuiltDefinitions(a *Analyzer) {
	k := utils.MakeStack[nodes.DataType](&INT_DATATYPE)
	a.definedTypes[lexer.INT] = &k
	k = utils.MakeStack[nodes.DataType](&FLOAT_DATATYPE)
	a.definedTypes[lexer.FLOAT] = &k
	k = utils.MakeStack[nodes.DataType](&BOOLEAN_DATATYPE)
	a.definedTypes[lexer.BOOLEAN] = &k
}

func isBooleanType(typ nodes.DataType) bool {
	// todo: also consider aliases, when supported
	// todo
	return true
}

func (a *Analyzer) operatorReturnType(op string, lval nodes.DataType, rval nodes.DataType) nodes.DataType {
	// todo: make more sophisticated by considering operand types
	// and having a way in the language to define return types for
	// any operator with any operand (op overloading)
	switch op {
	case "+":
		return &INT_DATATYPE
	case "-":
		return &INT_DATATYPE
	case "*":
		return &INT_DATATYPE
	case "/":
		return &INT_DATATYPE
	}
	return &ERROR_TYPE
}
