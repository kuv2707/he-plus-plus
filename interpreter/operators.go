package interpreter

import (
	"fmt"
	"toylingo/parser"
	"toylingo/utils"
)

func evaluateOperator(node parser.TreeNode, ctx *scopeContext) Variable {
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
	return Variable{}
}

func evaluateExpression(node *parser.TreeNode, ctx *scopeContext) Variable {
	LineNo = node.LineNo
	switch node.Label {
	case "operator":
		return evaluateOperator(*node, ctx)
	case "literal":
		fallthrough
	case "primary":
		return evaluatePrimary(*node, ctx)
	case "call":
		ret := evaluateFuncCall(*node, ctx)
		if ret.pointer == nil {
			interrupt("function " + node.Description + " does not return a value but is expected to")
		}
		return ret
	default:
		node.PrintTree("->")
		panic("invalid expression " + node.Label)
	}
}

func evaluateAssignment(ctx *scopeContext, node parser.TreeNode) Variable {
	if node.Children[0].Description == "index" {
		return assignToArrayIndex(ctx, node)
	}
	variableName := node.Children[0].Description
	variableValue := evaluateExpression(node.Children[1], ctx)
	val, alreadyExists := ctx.variables[variableName]
	if alreadyExists {
		//todo: make a function to copy value from one pointer to another: memcpy
		writeBits(*val.pointer, heapSlice(variableValue.pointer.address, variableValue.pointer.size))
		return val
	}
	ctx.variables[variableName] = variableValue
	variableValue.pointer.temp = false

	return variableValue
}

func assignToArrayIndex(ctx *scopeContext, node parser.TreeNode) Variable {
	arrayVarname := node.Children[0].Properties["array"].Description
	arrayVar, exists := ctx.variables[arrayVarname]
	if !exists {
		interrupt("array " + arrayVarname + " does not exist in current scope")
	}
	if arrayVar.vartype != TYPE_ARRAY {
		interrupt("variable " + arrayVarname + " is not an array")
	}
	indexVar := evaluateExpression(node.Children[0].Properties["index"], ctx)
	index := int(getNumber(indexVar))
	size := arrayVar.pointer.size / type_sizes[TYPE_POINTER]
	if index >= size || index < 0 {
		interrupt("cannot assign to index ", index, " of array ", arrayVarname, " of length", size)
	}
	newval := evaluateExpression(node.Children[1], ctx)
	newval.pointer.temp = false
	pointerToValueBits := arrayVar.pointer.address + type_sizes[TYPE_POINTER]*int(index)
	freePtr(pointers[byteArrayToPointer(heapSlice(pointerToValueBits, type_sizes[TYPE_POINTER]))])
	unsafeWriteBits(pointerToValueBits, pointerByteArray(newval.pointer.address))
	return newval
}

func evaluateLogical(ctx *scopeContext, node parser.TreeNode, operator string) Variable {
	left := evaluateExpression(node.Children[0], ctx)
	right := evaluateExpression(node.Children[1], ctx)
	if left.vartype == "bool" && right.vartype == "bool" {
		value := false
		switch operator {
		case "&&":
			value = getBool(left) && getBool(right)
		case "||":
			value = getBool(left) || getBool(right)
		}
		memaddr := malloc(type_sizes[TYPE_BOOLEAN], ctx.scopeId, true)
		var val byte = 0
		if value {
			val = 1
		}
		writeBits(*memaddr, []byte{val})
		return Variable{memaddr, TYPE_BOOLEAN}
	} else {
		interrupt("invalid operands to binary operator " + operator)
	}
	return Variable{}
}



//todo: cleanup, optimize, simplify
func evaluateDMAS(ctx *scopeContext, node parser.TreeNode, operator string) Variable {
	left := evaluateExpression(node.Children[0], ctx)
	right := evaluateExpression(node.Children[1], ctx)
	if left.vartype == TYPE_NUMBER && right.vartype == TYPE_NUMBER {
		leftVal := getValue(left)
		rightVal := getValue(right)
		value := 0.0
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
		memaddr := malloc(type_sizes[TYPE_NUMBER], ctx.scopeId, true)
		writeBits(*memaddr, numberByteArray(value))
		return Variable{memaddr, TYPE_NUMBER}
	} else if left.vartype == TYPE_STRING && right.vartype == TYPE_STRING {
		if operator != "+" {
			interrupt("Invalid operator", operator, "for string operands")
		}
		leftVal := byteArrayString(heapSlice(left.pointer.address, left.pointer.size))
		rightVal := byteArrayString(heapSlice(right.pointer.address, right.pointer.size))
		newval := leftVal + rightVal
		memaddr := malloc(type_sizes[TYPE_CHAR]*len(newval), ctx.scopeId, true)
		writeBits(*memaddr, stringByteArray(newval))
		return Variable{memaddr, TYPE_STRING}
	} else if left.vartype == TYPE_STRING && right.vartype == TYPE_NUMBER {
		leftVal := byteArrayString(heapSlice(left.pointer.address, left.pointer.size))
		rightVal := getValue(right)
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
		memaddr := malloc(type_sizes[TYPE_CHAR]*len(newval), ctx.scopeId, true)
		writeBits(*memaddr, stringByteArray(newval))
		return Variable{memaddr, TYPE_STRING}
	} else if left.vartype == TYPE_NUMBER && right.vartype == TYPE_STRING{
		leftVal := getValue(left)
		rightVal := byteArrayString(heapSlice(right.pointer.address, right.pointer.size))
		newval := ""
		switch operator {
		case "+":
			newval = fmt.Sprint(leftVal) + rightVal
		default:
			interrupt("number and string operands cannot be used with operator", operator)
		}
		memaddr := malloc(type_sizes[TYPE_CHAR]*len(newval), ctx.scopeId, true)
		writeBits(*memaddr, stringByteArray(newval))
		return Variable{memaddr, TYPE_STRING}

	} else {
		interrupt("invalid operands to binary operator " + operator)
	}
	return Variable{}
}

func evaluateComparison(ctx *scopeContext, node parser.TreeNode, operator string) Variable {
	left := evaluateExpression(node.Children[0], ctx)
	right := evaluateExpression(node.Children[1], ctx)
	if left.vartype == TYPE_NUMBER && right.vartype == TYPE_NUMBER {
		value := false
		leftVal := getValue(left)
		rightVal := getValue(right)
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
		memaddr := malloc(type_sizes[TYPE_BOOLEAN], ctx.scopeId, true)
		var val byte = 0
		if value {
			val = 1
		}
		writeBits(*memaddr, []byte{val})
		return Variable{memaddr, TYPE_BOOLEAN}
	} else {
		interrupt("invalid operands to binary operator " + operator)
	}
	return Variable{}
}

func evaluateUnary(node parser.TreeNode, ctx *scopeContext, operator string) Variable {
	pm := 1.0
	switch operator {
	case "--":
		pm = -1
		fallthrough
	case "++":
		varname := node.Children[0].Description
		varval, exists := ctx.variables[varname]
		if !exists {
			interrupt("cannot increment variable " + varname + " as it does not exist in current scope")
		}
		if varval.vartype != TYPE_NUMBER {
			interrupt("cannot increment variable " + varname + " as it is not a number")
		}
		writeBits(*varval.pointer, numberByteArray(getNumber(varval)+pm)) //todo: maybe optimize?
		return varval
	}
	val := evaluateExpression(node.Children[0], ctx)
	if val.vartype == TYPE_NUMBER {
		switch operator {
		case "+":
			return val
		case "-":
			memaddr := malloc(type_sizes[TYPE_NUMBER], ctx.scopeId, true)
			writeBits(*memaddr, numberByteArray(-getNumber(val)))
			return Variable{memaddr, TYPE_NUMBER}

		default:
			interrupt("invalid unary operator " + operator)
		}
	}
	if val.vartype == "bool" {
		switch operator {
		case "!":
			memaddr := malloc(type_sizes[TYPE_BOOLEAN], ctx.scopeId, true)
			var valb byte = 1
			if !getBool(val) {
				valb = 0
			}
			writeBits(*memaddr, []byte{valb})
			return Variable{memaddr, TYPE_BOOLEAN}
		default:
			interrupt("invalid unary operator " + operator)
		}
	}
	interrupt("invalid operand to unary operator " + operator)
	return Variable{}
}

func evaluatePrimary(node parser.TreeNode, ctx *scopeContext) Variable {
	val := node.Description
	if isCompositeDS(node) {
		//func call or array or object
		return evaluateCompositeDS(node, ctx)
	}
	//todo: redundant to verify type here: done in AST phase
	if utils.IsNumber(val) {
		memaddr := malloc(type_sizes[TYPE_NUMBER], ctx.scopeId, true)
		writeBits(*memaddr, numberByteArray(StringToNumber(val)))
		return Variable{memaddr, TYPE_NUMBER}
	} else if utils.IsBoolean(val) {
		memaddr := malloc(type_sizes[TYPE_BOOLEAN], ctx.scopeId, true)
		var boolnum byte = 0
		if utils.StringToBoolean(val) {
			boolnum = 1
		}
		writeBits(*memaddr, []byte{boolnum})
		return Variable{memaddr, TYPE_BOOLEAN}
	} else if utils.IsString(val) {
		memaddr := malloc(type_sizes[TYPE_CHAR]*(len(val)-2), ctx.scopeId, true)
		writeBits(*memaddr, stringByteArray(val[1:len(val)-1]))
		return Variable{memaddr, TYPE_STRING}
	} else {
		if v, exists := ctx.variables[val]; exists {
			if v.vartype == TYPE_ARRAY {
				//leads to array being freed twice
				return v
			}
			copy := copyVariable(v, ctx.scopeId)
			return copy
		} else {
			LineNo = node.LineNo
			interrupt("variable " + val + " does not exist in current scope")
		}
	}
	return Variable{}
}

func evaluateFuncCall(node parser.TreeNode, ctx *scopeContext) Variable {

	funcNode, exists := ctx.functions[node.Description]
	if !exists {
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
		argValue.pointer.temp = false
		if argValue.vartype != TYPE_ARRAY {
			argValue.pointer.scopeId = newCtx.scopeId
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
	var ret Variable = Variable{}
	//fix to the memory leak bug
	if newCtx.returnValue != nil {
		ret = copyVariable(*newCtx.returnValue, ctx.scopeId)
		freePtr(newCtx.returnValue.pointer)
	}
	return ret
}

func evaluateCompositeDS(node parser.TreeNode, ctx *scopeContext) Variable {
	switch node.Description {
	case "array":
		return evaluateArray(node, ctx)
	case "index":
		return evaluateArrayIndex(node, ctx)
	}
	panic("invalid composite data structure")

}

// returns a Variable with pointer to an array
// the arrangement of an array variable is as follows:
// 1. the first 8 bytes contain the length of the array
// 2. the next 8*n bytes contain the pointer to the first element of the array, or actual elements if it is a primitive

func evaluateArray(node parser.TreeNode, ctx *scopeContext) Variable {
	len := len(node.Children)
	memaddr := malloc(type_sizes[TYPE_POINTER]*len, ctx.scopeId, true)
	addr := memaddr.address
	for i := range node.Children {
		val := evaluateExpression(node.Children[i], ctx)
		val.pointer.temp = false
		unsafeWriteBits(addr, pointerByteArray(val.pointer.address))
		addr += type_sizes[TYPE_POINTER]
	}
	return Variable{memaddr, TYPE_ARRAY}
}

func evaluateArrayIndex(node parser.TreeNode, ctx *scopeContext) Variable {
	// printVariableList(ctx.variables)
	array, yes := ctx.variables[node.Properties["array"].Description]
	if !yes {
		interrupt("array " + node.Properties["array"].Description + " does not exist in current scope")
	}
	index := getNumber(evaluateExpression(node.Properties["index"], ctx))
	size := array.pointer.size
	if int(index) >= size || index < 0 {
		interrupt("array index", index, " out of bounds for length", size)
	}
	ptr := byteArrayToPointer(heapSlice(array.pointer.address+type_sizes[TYPE_POINTER]*int(index), type_sizes[TYPE_POINTER]))
	//assuming number is stored at ptr
	value := byteArrayToFloat64(heapSlice(ptr, type_sizes[TYPE_NUMBER]))
	memaddr := malloc(type_sizes[TYPE_NUMBER], ctx.scopeId, true)
	writeBits(*memaddr, numberByteArray(value))
	return Variable{memaddr, TYPE_NUMBER}
}
