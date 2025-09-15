package staticanalyzer

import (
	"fmt"
	"he++/lexer"
	nodes "he++/parser/node_types"
	// "runtime/debug"
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
			// op := v.Op
			if !left.Equals(right) {
				return &nodes.ErrorType{Message: fmt.Sprintf("$$Can't perform %s on %s and %s$$", v.Op, left.Text(), right.Text())}
			}
			return left
		}
	case *nodes.PrePostOperatorNode:
		{
			// for now, `&` changes type
			var pref nodes.TypePrefix
			switch v.Op {
			case lexer.AMP:
				pref = nodes.PointerOf
			case lexer.MUL:
				pref = nodes.Dereference
				fmt.Println("->"+v.Operand.String(""))
				if ch, ok := v.Operand.(*nodes.PrePostOperatorNode); ok && ch.Op == lexer.AMP {
					return computeType(ch.Operand, a)
				} else {
					return &nodes.ErrorType{Message: fmt.Sprintf("Cannot dereference type %s", v.String(""))}
				}
			default:
				// pref := nodes.Unknown
			}
			return &nodes.PrefixOfType{Prefix: pref, OfType: computeType(v.Operand, a)}
		}
	case nil:
		{
			return &nodes.VoidType{}
		}
	case *nodes.ArrayDeclaration:
		{
			expectedType := v.DataT
			for i, elem := range v.Elems {
				typ := computeType(elem, a)
				if !typ.Equals(expectedType) {
					// todo: check possibility of type casting
					return &nodes.ErrorType{Message: fmt.Sprintf("%d th element should be of type %s, not %s", i, expectedType.Text(), typ.Text())}
				}
			}
			return &nodes.PrefixOfType{Prefix: nodes.ArrayOf, OfType: v.DataT}
		}
	case *nodes.FuncNode:
		argtypes := make([]nodes.DataType, 0)
		for _, arg := range v.ArgList {
			argtypes = append(argtypes, arg.DataT)
		}
		return &nodes.FuncType{
			ReturnType: v.ReturnType,
			ArgTypes:   argtypes,
		}
	default:
		// fmt.Println("Stack trace:")
		// debug.PrintStack()
		return &nodes.ErrorType{Message: fmt.Sprintf("Can't compute type for %T", v)}
	}
}
