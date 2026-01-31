package staticanalyzer

import (
	"fmt"
	"he++/lexer"
	nodes "he++/parser/node_types"
	"he++/utils"
	// "runtime/debug"
)

// this function partners with checkExpression for checking that area of the tree.
// Assume that the function is not idempotent.
func (a *Analyzer) computeType(n nodes.TreeNode) nodes.DataType {
	switch v := n.(type) {
	case *nodes.BooleanNode, *nodes.NumberNode, *nodes.IdentifierNode, *nodes.InfixOperatorNode:
		return a.checkExpression(v)

	case *nodes.PrePostOperatorNode:
		{
			dt := nodes.NONE
			switch v.Op {
			case lexer.AMP:
				dt = &nodes.PrefixOfType{
					Prefix: nodes.PointerOf,
					OfType: a.computeType(v.Operand),
					DataTypeMetaData: nodes.DataTypeMetaData{
						TypeSize: nodes.POINTER_SIZE,
						Tid:      nodes.UniqueTypeId(),
					},
				}
			case lexer.MUL:
				operandType := a.computeType(v.Operand)
				if ch, ok := operandType.(*nodes.PrefixOfType); ok && ch.Prefix == nodes.PointerOf {
					dt = ch.OfType
				} else {
					a.AddError(v.Range().Start, utils.TypeError, fmt.Sprintf("Cannot dereference type %s", utils.Cyan(operandType.Text())))
					dt = ERROR_TYPE
				}
			case lexer.SUB:
				operandType := a.computeType(v.Operand)
				if !isNumericType(operandType) {
					a.AddError(v.Range().Start, utils.TypeError, fmt.Sprintf("Cannot negate value of type %s", utils.Cyan(operandType.Text())))
				}
				dt = operandType

			default:
				dt = ERROR_TYPE
			}
			v.ResultDT = dt
			return dt
		}
	case nil:
		{
			return nodes.VOID_DATATYPE
		}
	case *nodes.ArrayDeclarationNode:
		{
			sizeType := a.computeType(v.SizeNode)
			if !isNumericType(sizeType) {
				a.AddError(v.SizeNode.Range().Start, utils.TypeError, fmt.Sprintf("Size of array should be numeric"))
			}
			a.verifyAndNormalize(&v.DataT)
			expectedType := v.DataT
			for i, elem := range v.Elems {
				typ := a.computeType(elem)
				if !typ.Equals(expectedType) {
					// todo: check possibility of type casting
					a.AddError(elem.Range().Start, utils.TypeError,
						fmt.Sprintf("Element at index %d of type %s cannot be casted to %s", i, utils.Cyan(typ.Text()), utils.Cyan(expectedType.Text())))
				}
			}
			return &nodes.PrefixOfType{Prefix: nodes.ArrayOf, OfType: v.DataT, DataTypeMetaData: nodes.DataTypeMetaData{TypeSize: nodes.POINTER_SIZE, // todo: consider number of elems in type def
				Tid: nodes.UniqueTypeId(),
			}}
		}
	case *nodes.FuncNode:
		argtypes := make([]nodes.DataType, 0)
		for i := range v.ArgList {
			if !a.verifyAndNormalize(&v.ArgList[i].DataT) {
				a.AddError(v.Range().Start, utils.UndefinedError, fmt.Sprintf("Arg type %s is undefined or depends on an undefined type", utils.Cyan(v.ArgList[i].DataT.Text())))
			}
			argtypes = append(argtypes, v.ArgList[i].DataT)
		}
		if !a.verifyAndNormalize(&v.ReturnType) {
			a.AddError(v.Range().Start, utils.UndefinedError, fmt.Sprintf("Return type %s is undefined or depends on an undefined type", utils.Cyan(v.ReturnType.Text())))
		}
		return &nodes.FuncType{
			ReturnType:       v.ReturnType,
			ArgTypes:         argtypes,
			DataTypeMetaData: nodes.DataTypeMetaData{TypeSize: nodes.POINTER_SIZE, Tid: nodes.UniqueTypeId()},
		}
	case *nodes.ArrIndNode:
		arrType := a.computeType(v.ArrProvider)
		indexerType := a.computeType(v.Indexer)
		indexedValueType, ok := isIndexable(a, arrType, indexerType)
		v.DataType = indexedValueType
		if !ok {
			a.AddError(v.Range().Start, utils.TypeError, fmt.Sprintf("The type %s cannot be indexed by %s", utils.Cyan(arrType.Text()), utils.Cyan(indexerType.Text())))
			return ERROR_TYPE
		}
		return indexedValueType
	case *nodes.StringNode:
		return &nodes.PrefixOfType{Prefix: nodes.ArrayOf, OfType: &BYTE_DATATYPE, DataTypeMetaData: nodes.DataTypeMetaData{TypeSize: nodes.POINTER_SIZE, Tid: nodes.UniqueTypeId()}}
	case *nodes.FuncCallNode:
		funcType := a.computeType(v.Callee)
		ftyp, ok := funcType.(*nodes.FuncType)
		if !ok {
			a.AddError(
				v.Range().Start,
				utils.TypeError,
				fmt.Sprintf("Type is not callable: %s", utils.Cyan(funcType.Text())),
			)
			return ERROR_TYPE
		}
		v.CalleeT = ftyp
		funcInf := utils.MakeASTPrinter()
		v.Callee.String(&funcInf)
		funcNameTreeStr := funcInf.Builder.String()

		if len(v.Args) != len(ftyp.ArgTypes) {
			a.AddError(
				v.Range().Start,
				utils.TypeError,
				fmt.Sprintf("Function %s expects %s parameters, but supplied %s", utils.Blue(funcNameTreeStr), utils.Yellow(fmt.Sprint(len(ftyp.ArgTypes))), utils.Yellow(fmt.Sprint(len(v.Args)))),
			)
			return ERROR_TYPE
		}
		for i, k := range v.Args {
			expT := ftyp.ArgTypes[i]
			passedT := a.computeType(k)

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
		return ERROR_TYPE
	}
}

// types are dependent on one another, forming a graph.
// todo: handle loops in type definition
// todo: replace NamedType with the actual DataType objects the name points to
func (a *Analyzer) verifyAndNormalize(dt *nodes.DataType) bool {
	switch v := (*dt).(type) {
	case *nodes.ErrorType:
		return false
	case *nodes.FuncType:
		ok := true
		for i := range v.ArgTypes {
			ok = ok && a.verifyAndNormalize(&v.ArgTypes[i])
		}
		ok = ok && a.verifyAndNormalize(&v.ReturnType)
		return ok
	case *nodes.NamedType:
		baseT, exists := a.GetType(v.Name)
		if exists {
			*dt = baseT
		}
		return exists
	case *nodes.PrefixOfType:
		// todo: check if this type can be prefixed with this prefix
		return a.verifyAndNormalize(&v.OfType)
	case *nodes.StructType:
		ok := true
		for i := range v.Fields {
			ok = ok && a.verifyAndNormalize(&v.Fields[i].Type)
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
