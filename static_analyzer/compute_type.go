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
			return &nodes.NamedType{Name: lexer.BOOLEAN}
		}
	case *nodes.NumberNode:
		{
			return &nodes.NamedType{Name: v.NumType}
		}
	case *nodes.IdentifierNode:
		{
			varname := v.Name()
			norm_tname, exists := a.definedSyms[varname]
			if !exists {
				return &nodes.ErrorType{Message: "UNDEFINED_TYPE"}
			}
			return norm_tname
		}
	case *nodes.InfixOperatorNode:
		{
			left := computeType(v.Left, a)
			right := computeType(v.Right, a)
			if left != right {
				return &nodes.ErrorType{Message: fmt.Sprintf("$$Can't perform %s on %s and %s$$", v.Op, left.Text(), right.Text())}
			}
			return left
		}
	case *nodes.PrePostOperatorNode:
		{
			// for now, `&` changes type
			return &nodes.PrefixOfType{Prefix: nodes.PointerOf, OfType: computeType(v.Operand, a)}
		}
	case nil:
		{
			return &nodes.VoidType{}
		}
	case *nodes.ArrayDeclaration:
		{
			return &nodes.PrefixOfType{Prefix: nodes.ArrayOf, OfType: v.DataT}
		}
	default:
		return &nodes.ErrorType{Message: fmt.Sprintf("Can't compute type for %T", v)}
	}
}
