package staticanalyzer

import (
	"fmt"
	"he++/lexer"
	nodes "he++/parser/node_types"
	"he++/utils"
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
				a.AddError(v.Range().Start, utils.UndefinedError, fmt.Sprintf("Undefined identifier %s", utils.Green(varname)))
				return &nodes.ErrorType{}
			}
			return norm_tname
		}
	case *nodes.InfixOperatorNode:
		{
			left := computeType(v.Left, a)
			right := computeType(v.Right, a)
			// op := v.Op
			if !left.Equals(right) {
				a.AddError(v.Range().Start, utils.TypeError, fmt.Sprintf("Can't perform %s on types %s and %s", v.Op, utils.Cyan(left.Text()), utils.Cyan(right.Text())))
				return &nodes.ErrorType{}
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
				operandType := computeType(v.Operand, a)
				if ch, ok := operandType.(*nodes.PrefixOfType); ok && ch.Prefix == nodes.PointerOf {
					return ch.OfType
				} else {
					a.AddError(v.Range().Start, utils.TypeError, fmt.Sprintf("Cannot dereference type %s", utils.Cyan(operandType.Text())))
					return &nodes.ErrorType{}
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
	case *nodes.ArrayDeclarationNode:
		{
			expectedType := v.DataT
			for i, elem := range v.Elems {
				typ := computeType(elem, a)
				if !typ.Equals(expectedType) {
					// todo: check possibility of type casting
					a.AddError(elem.Range().Start, utils.TypeError,
						fmt.Sprintf("Element at index %d of type %s cannot be casted to %s", i, utils.Cyan(typ.Text()), utils.Cyan(expectedType.Text())))
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
	case *nodes.ArrIndNode:
		arrType := computeType(v.ArrProvider, a)
		indexerType := computeType(v.Indexer, a)
		indexedValueType, ok := isIndexable(a, arrType, indexerType)
		if !ok {
			a.AddError(v.Range().Start, utils.TypeError, fmt.Sprintf("The type %s cannot be indexed by %s", utils.Cyan(arrType.Text()), utils.Cyan(indexerType.Text())))
			return &nodes.ErrorType{}
		}
		return indexedValueType
	case *nodes.StringNode:
		return &nodes.NamedType{Name: lexer.STRING}
	default:
		a.AddError(v.Range().Start, utils.UndefinedError, fmt.Sprintf("Can't compute type for %T", v))
		return &nodes.ErrorType{}
	}
}
