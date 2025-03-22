package staticanalyzer

import (
	"he++/lexer"
	nodes "he++/parser/node_types"
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

// func checkChar(n nodes.TreeNode) bool {
// 	switch v := n.(type) {
// 	case *nodes.CharNode:
// 		return true
// 	default:
// 		return false
// 	}
// }

func getPrimitiveTypeDefns() map[string]DataTypeInfo {
	definedTypes := make(map[string]DataTypeInfo, 0)
	definedTypes[lexer.INT] = &PrimitiveType{
		numBytes: 4,
	}
	definedTypes[lexer.FLOAT] = &PrimitiveType{
		numBytes: 4,
	}
	definedTypes[lexer.BOOLEAN] = &PrimitiveType{
		numBytes: 1,
	}
	// definedTypes[lexer.CHAR] = &PrimitiveType{
	// 	numBytes: 1,
	// 	validator: checkChar,
	// }
	return definedTypes
}
