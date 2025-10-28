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

func checkInt(n nodes.TreeNode) bool {
	nn, ok := n.(*nodes.NumberNode)
	if !ok {
		return false
	}
	// todo: remove hardcoded str
	return nn.NumType == "int"
}

func checkFloat(n nodes.TreeNode) bool {
	nn, ok := n.(*nodes.NumberNode)
	if !ok {
		return false
	}
	// todo: remove hardcoded str
	return nn.NumType == "float"
}

func checkBool(n nodes.TreeNode) bool {
	_, ok := n.(*nodes.BooleanNode)
	return ok
}


func addInbuiltDefinitions(a *Analyzer) {
	definedTypes := make(map[string]utils.Stack[nodes.DataType], 0)

	definedTypes[lexer.INT] = utils.MakeStack[nodes.DataType](&nodes.NamedType{Name: lexer.INT})
	definedTypes[lexer.FLOAT] = utils.MakeStack[nodes.DataType](&nodes.NamedType{Name: lexer.FLOAT})
	definedTypes[lexer.BOOLEAN] = utils.MakeStack[nodes.DataType](&nodes.NamedType{Name: lexer.BOOLEAN})
}
