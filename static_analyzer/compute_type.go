package staticanalyzer

import (
	"fmt"
	"he++/lexer"
	nodes "he++/parser/node_types"
)

var _ = fmt.Println

func computeType(n nodes.TreeNode, a *Analyzer) nodes.DataType {
	switch v := n.(type) {
	case *nodes.BooleanNode:
		{
			// todo: don't depend on lexer?
			return nodes.DataType{Text: lexer.BOOLEAN}
		}
	case *nodes.NumberNode:
		{
			return nodes.DataType{Text: v.NumType}
		}
	case *nodes.IdentifierNode:
		{
			varname := v.Name()
			norm_tname, exists := a.definedSyms[varname]
			if !exists {
				return nodes.DataType{Text: "$$ERROR$$"}
			}
			return norm_tname
		}
	case *nodes.InfixOperatorNode:
		{
			left := computeType(v.Left, a)
			right := computeType(v.Right, a)
			if left != right {
				return nodes.DataType{Text: "$$ERROR$$"}
			}
			return left
		}
	case *nodes.PrePostOperatorNode:
		{
			// for now, `&` changes type
			return nodes.DataType{Text: "ptr_" + computeType(v.Operand, a).Text}
		}
	default:
		return nodes.DataType{Text: "$$ERROR$$"}
	}
}
