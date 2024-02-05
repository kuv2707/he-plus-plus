package interpreter

import (
	"fmt"
	"he++/parser"
	"he++/utils"
	"math"
)

func evaluateOperator(node parser.TreeNode, ctx *scopeContext) *Pointer {
	LineNo = node.LineNo
	if node.Label == "literal" {
		return evaluatePrimary(node, ctx)
	}
	if utils.IsOneOf(node.Description, []string{"+", "-", "*", "/"}) {
		if len(node.Children) == 1 {
			return evaluateUnary(node, ctx, node.Description)
		} else {

			return evaluateDMAS(ctx, node, node.Description)
		}
	} else if utils.IsOneOf(node.Description, []string{"<", ">", "<=", ">=", "==", "!="}) {
		return evaluateComparison(ctx, node, node.Description)
	} else if node.Description == "=" {
		return evaluateAssignment(ctx, node)
	} else if utils.IsOneOf(node.Description, []string{"&&", "||"}) {
		return evaluateLogical(ctx, node, node.Description)
	} else if utils.IsOneOf(node.Description, []string{"++", "--", "!"}) {
		return evaluateUnary(node, ctx, node.Description)
	}
	interrupt("invalid operator " + node.Description)
	return NULL_POINTER
}

func evaluateExpression(node *parser.TreeNode, ctx *scopeContext) *Pointer {
	LineNo = node.LineNo
	ret := NULL_POINTER
	switch node.Label {
	case "operator":
		ret = evaluateOperator(*node, ctx)
	case "literal":
		fallthrough
	case "primary":
		ret = evaluatePrimary(*node, ctx)
	case "call":
		ret = evaluateFuncCall(*node, ctx)
		if ret == NULL_POINTER {
			interrupt("function " + node.Description + " does not return a value but is expected to")
		}
		ret.scopeId = ctx.scopeId
	default:
		node.PrintTree("->")
		panic("invalid expression " + node.Label)
	}
	return ret
}

func evaluateAssignment(ctx *scopeContext, node parser.TreeNode) *Pointer {
	// if node.Children[0].Description == "index" {
	// 	return assignToArrayIndex(ctx, node)
	// }
	variableName := node.Children[0].Description
	variableValue := evaluateExpression(node.Children[1], ctx)
	val := findVariable(variableName)
	// fmt.Println(variableName, " exists:", !val.isNull())
	if !val.isNull() {
		if val.getDataType() != variableValue.getDataType() {
			interrupt("cannot assign", variableValue.getDataType().String(), "to", val.getDataType().String())
		}
		writeContentFromOnePointerToAnother(val, variableValue)
		freePtr(variableValue)
	} else {

		variableValue.temp = false
		ctx.variables[variableName] = variableValue
	}
	debug_info("assigned", variableName, "to", variableValue)
	return variableValue
}

// func assignToArrayIndex(ctx *scopeContext, node parser.TreeNode) Variable {
// 	arrayVarname := node.Children[0].Properties["array"].Description
// 	arrayVar, exists := ctx.variables[arrayVarname]
// 	if !exists {
// 		interrupt("array " + arrayVarname + " does not exist in current scope")
// 	}
// 	if arrayVar.vartype != TYPE_ARRAY {
// 		interrupt("variable " + arrayVarname + " is not an array")
// 	}
// 	indexVar := evaluateExpression(node.Children[0].Properties["index"], ctx)
// 	index := int(getNumber(indexVar))
// 	size := arrayVar.pointer.size / type_sizes[TYPE_POINTER]
// 	if index >= size || index < 0 {
// 		interrupt("cannot assign to index ", index, " of array ", arrayVarname, " of length", size)
// 	}
// 	newval := evaluateExpression(node.Children[1], ctx)
// 	newval.pointer.temp = false
// 	pointerToValueBytes := arrayVar.pointer.address + type_sizes[TYPE_POINTER]*int(index)
// 	freePtr(pointers[byteArrayToPointer(heapSlice(pointerToValueBytes, type_sizes[TYPE_POINTER]))])
// 	unsafeWriteBytes(pointerToValueBytes, pointerAsBytes(newval.pointer.address))
// 	return newval
// }

func evaluateLogical(ctx *scopeContext, node parser.TreeNode, operator string) *Pointer {
	left := evaluateExpression(node.Children[0], ctx)
	right := evaluateExpression(node.Children[1], ctx)
	if left.getDataType() == BOOLEAN && right.getDataType() == BOOLEAN {
		value := false
		switch operator {
		case "&&":
			value = booleanValue(left) && booleanValue(right)
		case "||":
			value = booleanValue(left) || booleanValue(right)
		}
		ptr := malloc(type_sizes[BOOLEAN], ctx.scopeId, true)
		ptr.setDataType(BOOLEAN)
		var val byte = 0
		if value {
			val = 1
		}
		writeDataContent(ptr, []byte{val})
		return ptr
	} else {
		interrupt("invalid operands for binary operator " + operator)
	}
	return NULL_POINTER
}

func evaluateDMAS(ctx *scopeContext, node parser.TreeNode, operator string) *Pointer {
	left := evaluateExpression(node.Children[0], ctx)
	right := evaluateExpression(node.Children[1], ctx)
	if left.getDataType() == NUMBER && right.getDataType() == NUMBER {
		leftVal := numberValue(left)
		rightVal := numberValue(right)
		value := math.NaN()
		switch operator {
		case "+":
			value = leftVal + rightVal
		case "-":
			value = leftVal - rightVal
		case "*":
			value = leftVal * rightVal
		case "/":
			value = leftVal / rightVal
		}
		ptr := malloc(type_sizes[NUMBER], ctx.scopeId, true)
		ptr.setDataType(NUMBER)
		writeDataContent(ptr, numberByteArray(value))
		return ptr
	} else if left.getDataType() == STRING && right.getDataType() == STRING {
		if operator != "+" {
			interrupt("Invalid operator", operator, "for string operands")
		}
		leftVal := stringValue(left)
		rightVal := stringValue(right)
		newval := leftVal + rightVal
		ptr := malloc(type_sizes[CHAR]*len(newval), ctx.scopeId, true)
		ptr.setDataType(STRING)
		writeDataContent(ptr, stringAsBytes(newval))
		return ptr
	} else if left.getDataType() == STRING && right.getDataType() == NUMBER {
		leftVal := stringValue(left)
		rightVal := numberValue(right)
		newval := ""
		switch operator {
		case "+":
			newval = leftVal + fmt.Sprint(rightVal)
		case "*":
			for rightVal > 0 {
				newval += leftVal
				rightVal -= 1
			}
		default:
			interrupt("string and number operands cannot be used with operator", operator)
		}
		ptr := malloc(type_sizes[CHAR]*len(newval), ctx.scopeId, true)
		ptr.setDataType(STRING)
		writeDataContent(ptr, stringAsBytes(newval))
		return ptr
	} else if left.getDataType() == NUMBER && right.getDataType() == NUMBER {
		leftVal := numberValue(left)
		rightVal := stringValue(right)
		newval := ""
		switch operator {
		case "+":
			newval = fmt.Sprint(leftVal) + rightVal
		default:
			interrupt("number and string operands cannot be used with operator", operator)
		}
		ptr := malloc(type_sizes[CHAR]*len(newval), ctx.scopeId, true)
		writeDataContent(ptr, stringAsBytes(newval))
		return ptr

	} else {
		interrupt("invalid operands for binary operator " + operator)
	}
	return NULL_POINTER
}

func evaluateComparison(ctx *scopeContext, node parser.TreeNode, operator string) *Pointer {
	left := evaluateExpression(node.Children[0], ctx)
	left.temp = false
	right := evaluateExpression(node.Children[1], ctx)
	defer func() {
		left.temp = true
		right.temp = true
	}()
	if left.getDataType() == NUMBER && right.getDataType() == NUMBER {
		leftVal := numberValue(left)
		rightVal := numberValue(right)
		value := false
		switch operator {
		case "<":
			value = leftVal < rightVal
		case ">":
			value = leftVal > rightVal
		case "<=":
			value = leftVal <= rightVal
		case ">=":
			value = leftVal >= rightVal
		case "==":
			value = leftVal == rightVal
		case "!=":
			value = leftVal != rightVal
		}
		ptr := malloc(type_sizes[BOOLEAN], ctx.scopeId, true)
		ptr.setDataType(BOOLEAN)
		var val byte = 0
		if value {
			val = 1
		}
		writeDataContent(ptr, []byte{val})
		return ptr
	} else {
		interrupt("invalid operands for binary operator " + operator)
	}
	return NULL_POINTER
}

func evaluateUnary(node parser.TreeNode, ctx *scopeContext, operator string) *Pointer {
	pm := 1.0
	switch operator {
	case "--":
		pm = -1
		fallthrough
	case "++":
		varname := node.Children[0].Description
		varval := findVariable(varname)
		if varval.isNull() {
			interrupt("cannot increment variable " + varname + " as it does not exist in current scope")
		}
		if varval.getDataType() != NUMBER {
			interrupt("cannot increment variable " + varname + " as it is not a number")
		}
		val := numberValue(varval)
		val += pm
		writeDataContent(varval, numberByteArray(val))
		return varval
	}

	val := evaluateExpression(node.Children[0], ctx)
	if val.getDataType() == NUMBER {
		switch operator {
		case "+":
			return val
		case "-":
			ptr := malloc(type_sizes[NUMBER], ctx.scopeId, true)
			ptr.setDataType(NUMBER)
			writeDataContent(ptr, numberByteArray(-numberValue(val)))
			return ptr

		default:
			interrupt("invalid unary operator " + operator)
		}
	} else if val.getDataType() == BOOLEAN {
		switch operator {
		case "!":
			memaddr := malloc(type_sizes[BOOLEAN], ctx.scopeId, true)
			var valb byte = 1
			if !booleanValue(val) {
				valb = 0
			}
			writeDataContent(memaddr, []byte{valb})
			return memaddr
		default:
			interrupt("invalid unary operator " + operator)
		}
	}
	interrupt("invalid operand", val.getDataType().String(), " to unary operator "+operator)
	return NULL_POINTER
}

// todo: use regex and simplify checks
func evaluatePrimary(node parser.TreeNode, ctx *scopeContext) *Pointer {
	val := node.Description
	// if isCompositeDS(node) {
	// 	//func call or array or object
	// 	return evaluateCompositeDS(node, ctx)
	// }
	//todo: redundant to verify type here: done in AST phase
	if utils.IsNumber(val) {
		ptr := malloc(type_sizes[NUMBER], ctx.scopeId, true)
		ptr.setDataType(NUMBER)
		writeDataContent(ptr, numberByteArray(StringToNumber(val)))
		return ptr
	} else if utils.IsBoolean(val) {
		ptr := malloc(type_sizes[BOOLEAN], ctx.scopeId, true)
		ptr.setDataType(BOOLEAN)
		var boolnum byte = 0
		if utils.StringToBoolean(val) {
			boolnum = 1
		}
		writeDataContent(ptr, []byte{boolnum})
		return ptr
	} else if utils.IsString(val) {
		ptr := malloc(type_sizes[CHAR]*(len(val)-2), ctx.scopeId, true)
		ptr.setDataType(STRING)
		writeDataContent(ptr, stringAsBytes(val[1:len(val)-1]))
		return ptr
	} else {
		if v := findVariable(val); !v.isNull() {
			return v.clone()
		} else {
			LineNo = node.LineNo
			interrupt("variable " + val + " does not exist in current scope")
		}
	}
	return NULL_POINTER
}

func evaluateFuncCall(node parser.TreeNode, ctx *scopeContext) *Pointer {

	function := findFunction(node.Description)
	funcNode := function
	if function == nil {
		interrupt("function " + node.Description + " does not exist in current scope")
	}
	newCtx := pushScopeContext(TYPE_FUNCTION, node.Description)
	lastValidLine := LineNo
	for i := 0; i < len(funcNode.Properties["args"].Children); i++ {
		argName := funcNode.Properties["args"].Children[i].Description
		argNode := node.Properties["args"+fmt.Sprint(i)]
		if argNode == nil {
			LineNo = lastValidLine
			interrupt("missing argument " + argName + " in function call " + funcNode.Description)
		}
		lastValidLine = argNode.LineNo
		argValue := evaluateExpression(argNode, newCtx)
		argValue.temp = false
		if argValue.getDataType() != ARRAY {
			argValue.scopeId = newCtx.scopeId
		}
		newCtx.variables[argName] = argValue
	}
	debug_info("calling", funcNode.Description)
	nfunc, nfexists := nativeFunctions[funcNode.Description]
	body, bexists := funcNode.Properties["body"]
	if bexists {
		executeScope(body, newCtx)
	} else if nfexists {
		nfunc.exec(newCtx)
		popScopeContext()
	} else {
		panic("SEVERE: internal logical error. func definition should have been present in either of the maps")
	}
	return newCtx.returnValue
}

// func evaluateCompositeDS(node parser.TreeNode, ctx *scopeContext) Variable {
// 	switch node.Description {
// 	case "array":
// 		return evaluateArray(node, ctx)
// 	case "index":
// 		return evaluateArrayIndex(node, ctx)
// 	}
// 	panic("invalid composite data structure")

// }

// returns a Variable with pointer to an array
// the arrangement of an array variable is as follows:
// 1. the first 8 bytes contain the length of the array
// 2. the next 8*n bytes contain the pointer to the first element of the array, or actual elements if it is a primitive

// func evaluateArray(node parser.TreeNode, ctx *scopeContext) Variable {
// 	len := len(node.Children)
// 	memaddr := malloc(type_sizes[TYPE_POINTER]*len, ctx.scopeId, true)
// 	addr := memaddr.address
// 	for i := range node.Children {
// 		val := evaluateExpression(node.Children[i], ctx)
// 		val.pointer.temp = false
// 		unsafeWriteBytes(addr, pointerAsBytes(val.pointer.address))
// 		addr += type_sizes[TYPE_POINTER]
// 	}
// 	return Variable{memaddr, TYPE_ARRAY}
// }

// func evaluateArrayIndex(node parser.TreeNode, ctx *scopeContext) Variable {
// 	// printVariableList(ctx.variables)
// 	array, yes := ctx.variables[node.Properties["array"].Description]
// 	if !yes {
// 		interrupt("array " + node.Properties["array"].Description + " does not exist in current scope")
// 	}
// 	index := getNumber(evaluateExpression(node.Properties["index"], ctx))
// 	size := array.pointer.size
// 	if int(index) >= size || index < 0 {
// 		interrupt("array index", index, " out of bounds for length", size)
// 	}
// 	ptr := byteArrayToPointer(heapSlice(array.pointer.address+type_sizes[TYPE_POINTER]*int(index), type_sizes[TYPE_POINTER]))
// 	//assuming number is stored at ptr
// 	value := byteArrayToFloat64(heapSlice(ptr, type_sizes[TYPE_NUMBER]))
// 	memaddr := malloc(type_sizes[TYPE_NUMBER], ctx.scopeId, true)
// 	writeBytes(*memaddr, numberByteArray(value))
// 	return Variable{memaddr, TYPE_NUMBER}
// }
