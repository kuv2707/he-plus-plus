package staticanalyzer

import (
	"fmt"
	"he++/lexer"
	nodes "he++/parser/node_types"
	"he++/utils"
)

type PrimitiveType struct {
	numBytes int
	typeName nodes.DataType
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
var INT_DATATYPE = nodes.NamedType{
	Name: lexer.INT,
	DataTypeMetaData: nodes.DataTypeMetaData{
		TypeSize:    4,
		Tid:         nodes.UniqueTypeId(),
		Fundamental: true,
	},
}
var FLOAT_DATATYPE = nodes.NamedType{
	Name: lexer.FLOAT,
	DataTypeMetaData: nodes.DataTypeMetaData{
		TypeSize:    4,
		Tid:         nodes.UniqueTypeId(),
		Fundamental: true,
	},
}

var BYTE_DATATYPE = nodes.NamedType{
	Name: "BYTE",
	DataTypeMetaData: nodes.DataTypeMetaData{
		TypeSize:    1,
		Tid:         nodes.UniqueTypeId(),
		Fundamental: true,
	},
}

var BOOLEAN_DATATYPE = nodes.NamedType{
	Name: lexer.BOOLEAN,
	DataTypeMetaData: nodes.DataTypeMetaData{
		TypeSize:    1,
		Tid:         nodes.UniqueTypeId(),
		Fundamental: true,
	},
}

func addFundamentalDefinitions(a *Analyzer) {
	a.definedTypes[lexer.INT] = utils.MakeStack[nodes.DataType](&INT_DATATYPE)
	a.definedTypes[lexer.FLOAT] = utils.MakeStack[nodes.DataType](&FLOAT_DATATYPE)
	a.definedTypes[lexer.BOOLEAN] = utils.MakeStack[nodes.DataType](&BOOLEAN_DATATYPE)
	a.definedTypes[lexer.VOID] = utils.MakeStack[nodes.DataType](&nodes.VOID_DATATYPE)

	// a.operatorTypeRelations[&INT_DATATYPE][&INT_DATATYPE] = &INT_DATATYPE
}

func isBooleanType(typ nodes.DataType) bool {
	// todo: also consider aliases, when supported
	nt, ok := typ.(*nodes.NamedType)
	return ok && nt.Name == lexer.BOOLEAN
}

func isErrorType(typ nodes.DataType) bool {
	_, ok := typ.(*nodes.ErrorType)
	return ok
}

type OperatorSignature struct {
	Left  nodes.DataType
	Right nodes.DataType
	Ret   nodes.DataType
}

var BasicArithmeticOpSigs = []OperatorSignature{
	{Left: &INT_DATATYPE, Right: &INT_DATATYPE, Ret: &INT_DATATYPE},
	{Left: &FLOAT_DATATYPE, Right: &FLOAT_DATATYPE, Ret: &FLOAT_DATATYPE},
	{Left: &INT_DATATYPE, Right: &FLOAT_DATATYPE, Ret: &FLOAT_DATATYPE},
	{Left: &FLOAT_DATATYPE, Right: &INT_DATATYPE, Ret: &FLOAT_DATATYPE},
}

var RelationOpSigs = []OperatorSignature{
	{Left: &INT_DATATYPE, Right: &INT_DATATYPE, Ret: &BOOLEAN_DATATYPE},
	{Left: &FLOAT_DATATYPE, Right: &FLOAT_DATATYPE, Ret: &BOOLEAN_DATATYPE},
	{Left: &INT_DATATYPE, Right: &FLOAT_DATATYPE, Ret: &BOOLEAN_DATATYPE},
	{Left: &FLOAT_DATATYPE, Right: &INT_DATATYPE, Ret: &BOOLEAN_DATATYPE},
}

var LogicalOpSigs = []OperatorSignature{
	{Left: &BOOLEAN_DATATYPE, Right: &BOOLEAN_DATATYPE, Ret: &BOOLEAN_DATATYPE},
}

// todo: shift inside Analyzer to make modifiable by src code
var OperatorRules = map[string][]OperatorSignature{
	lexer.ADD:     BasicArithmeticOpSigs,
	lexer.SUB:     BasicArithmeticOpSigs,
	lexer.MUL:     BasicArithmeticOpSigs,
	lexer.DIV:     BasicArithmeticOpSigs,
	lexer.ASSN:    BasicArithmeticOpSigs,
	lexer.LESS:    RelationOpSigs,
	lexer.GREATER: RelationOpSigs,
	lexer.LEQ:     RelationOpSigs,
	lexer.GEQ:     RelationOpSigs,
	lexer.EQ:      RelationOpSigs,
	lexer.ANDAND:  LogicalOpSigs,
	lexer.OROR:    LogicalOpSigs,
}

func (a *Analyzer) operatorReturnType(op string, lval nodes.DataType, rval nodes.DataType, ls int) nodes.DataType {
	// todo: make more sophisticated by considering operand types in
	// computing the operator return type
	// and having a way in the language to define return types for
	// any operator with any operand (op overloading)
	recs, exists := OperatorRules[op]
	if !exists {
		a.AddError(ls, utils.UndefinedError, fmt.Sprintf("Operator %s undefined for types %s and %s", utils.Magenta(op), utils.Cyan(lval.Text()), utils.Cyan(rval.Text())))
	} else {
		for _, rec := range recs {
			if rec.Left.Equals(lval) && rec.Right.Equals(rval) {
				return rec.Ret
			}
		}
	}
	return &ERROR_TYPE
}
