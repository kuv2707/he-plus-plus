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
			switch v.NumType {
			case nodes.FLOAT_NUMBER:
				return &FLOAT_DATATYPE
			case nodes.INT32_NUMBER:
				return &INT_DATATYPE
			default:
				return &ERROR_TYPE
			}
		}
	case *nodes.IdentifierNode:
		{
			varname := v.Name()
			s, exists := a.GetSym(varname)
			if !exists {
				a.AddError(v.Range().Start, utils.UndefinedError, fmt.Sprintf("Undefined identifier %s", utils.Green(varname)))
				return &ERROR_TYPE
			}
			return s
		}
	case *nodes.InfixOperatorNode:
		{
			left := computeType(v.Left, a)
			right := computeType(v.Right, a)
			// op := v.Op
			if !left.Equals(right) {
				a.AddError(v.Range().Start, utils.TypeError, fmt.Sprintf("Can't perform %s on types %s and %s", v.Op, utils.Cyan(left.Text()), utils.Cyan(right.Text())))
				return &ERROR_TYPE
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
					return &ERROR_TYPE
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
			if !a.verifyType(arg.DataT) {
				a.AddError(v.Range().Start, utils.UndefinedError, fmt.Sprintf("Arg type %s is undefined or depends on an undefined type", utils.Cyan(arg.DataT.Text())))
			}
			argtypes = append(argtypes, arg.DataT)
		}
		if !a.verifyType(v.ReturnType) {
			a.AddError(v.Range().Start, utils.UndefinedError, fmt.Sprintf("Return type %s is undefined or depends on an undefined type", utils.Cyan(v.ReturnType.Text())))
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
			return &ERROR_TYPE
		}
		return indexedValueType
	case *nodes.StringNode:
		return &nodes.NamedType{Name: lexer.STRING}
	case *nodes.FuncCallNode:
		funcType := computeType(v.Callee, a)
		ftyp, ok := funcType.(*nodes.FuncType)
		if !ok {
			a.AddError(
				v.Range().Start,
				utils.TypeError,
				fmt.Sprintf("Type is not callable: %s", utils.Cyan(funcType.Text())),
			)
			return &ERROR_TYPE
		}
		funcInf := utils.MakeASTPrinter()
		v.Callee.String(&funcInf)
		funcNameTreeStr := funcInf.Builder.String()

		if len(v.Args) != len(ftyp.ArgTypes) {
			a.AddError(
				v.Range().Start,
				utils.TypeError,
				fmt.Sprintf("Function %s expects %s parameters, but supplied %s", utils.Blue(funcNameTreeStr), utils.Yellow(fmt.Sprint(len(ftyp.ArgTypes))), utils.Yellow(fmt.Sprint(len(v.Args)))),
			)
			return &ERROR_TYPE
		}
		for i, k := range v.Args {
			expT := ftyp.ArgTypes[i]
			passedT := computeType(k, a)

			if !expT.Equals(passedT) {
				a.AddError(
					v.Range().Start,
					utils.TypeError,
					fmt.Sprintf("%d th parameter to function %s should be of type %s, not %s", i, utils.Blue(funcNameTreeStr), utils.Cyan(expT.Text()), utils.Cyan(passedT.Text())),
				)
			}
		}
		// todo: this return type isn't being verified to be defined.
		// verification should be done when storing the typedef from the func node
		return ftyp.ReturnType
	default:
		a.AddError(v.Range().Start, utils.UndefinedError, fmt.Sprintf("Can't compute type for %T", v))
		return &ERROR_TYPE
	}
}

// types are dependent on one another, forming a graph.
// todo: handle loops
func (a *Analyzer) verifyType(dt nodes.DataType) bool {
	switch v := dt.(type) {
	case *nodes.ErrorType:
		return false // or true?
	case *nodes.FuncType:
		ok := true
		for i := range v.ArgTypes {
			ok = ok && a.verifyType(v.ArgTypes[i])
		}
		ok = ok && a.verifyType(v.ReturnType)
		return ok
	case *nodes.NamedType:
		_, exists := a.GetType(v.Name)
		return exists
	case *nodes.PrefixOfType:
		// todo: check if this type can be prefixed with this prefix
		return a.verifyType(v.OfType)
	case *nodes.StructType:
		ok := true
		for i := range v.Fields {
			ok = ok && a.verifyType(v.Fields[i].Type)
		}
		return ok
	case *nodes.UnspecifiedType:
		return false
	case *nodes.VoidType:
		return true
	default:
		a.AddError(-1, utils.UndefinedError, fmt.Sprintf("Can't verify type %T", v))
		return false
	}
}
