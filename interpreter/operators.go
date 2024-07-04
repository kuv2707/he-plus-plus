package interpreter

import (
	"fmt"
	"he++/globals"
	"he++/parser"
	"he++/utils"
	"math"
)

func evaluateOperator(node parser.TreeNode, ctx *ScopeContext) *Pointer {
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
	interrupt(node.LineNo, "invalid operator "+node.Description)
	return NULL_POINTER
}

func evaluateExpression(node *parser.TreeNode, ctx *ScopeContext) *Pointer {
	ret := NULL_POINTER
	switch node.Label {
	case "operator":
		ret = evaluateOperator(*node, ctx)
	case "literal":
		fallthrough
	case "number":
		ret = evaluatePrimary(*node, ctx)
	case "boolean":
		ret = evaluatePrimary(*node, ctx)
	case "string":
		ret = evaluatePrimary(*node, ctx)
	case "variable":
		ret = evaluateVariable(*node, ctx)
	case "array":
		ret = evaluateArray(*node, ctx)
	case "index":
		ind := evaluateArrayIndex(*node, ctx)
		ret = pointers[bytesToInt(heapSlice(ind, type_sizes[POINTER]))]
	case "object":
		ret = evaluateObject(*node, ctx)
	case "call":
		ret = evaluateFuncCall(*node, ctx)
		if ret == NULL_POINTER {
			interrupt(node.LineNo, "function "+node.Description+" does not return a value but is expected to")
		}
	default:
		node.PrintTree("->")
		panic("invalid expression " + node.Label)
	}
	return ret
}

func evaluateAssignment(ctx *ScopeContext, node parser.TreeNode) *Pointer {
	variableValue := evaluateExpression(node.Children[1], ctx)
	if node.Children[0].Description == "index" {
		addr := evaluateArrayIndex(*node.Children[0], ctx)
		variableValue.temp = false
		old := mockPointer(bytesToInt(heapSlice(addr, 4)), true)
		old.changeReferenceCount(false)
		freePtr(old)
		unsafeWriteBytes(addr, intToBytes(variableValue.address))
		variableValue.setReferenceCount(1)
		return variableValue
	}
	variableName := node.Children[0].Description
	val, _ := findVariable(variableName)
	if !val.isNull() {
		if val.getDataType() != variableValue.getDataType() {
			interrupt(node.LineNo, "cannot assign", variableValue.getDataType().String(), "to", val.getDataType().String())
		}
		writeContentFromOnePointerToAnother(val, variableValue)
		return val
	}
	variableValue.temp = false
	variableValue.changeReferenceCount(true)
	ctx.variables[variableName] = variableValue
	return ctx.variables[variableName]
}

func evaluateLogical(ctx *ScopeContext, node parser.TreeNode, operator string) *Pointer {
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
		ptr := malloc(type_sizes[BOOLEAN], true)
		ptr.setDataType(BOOLEAN)
		var val byte = 0
		if value {
			val = 1
		}
		writeDataContent(ptr, []byte{val})
		return ptr
	} else {
		interrupt(node.LineNo, "invalid operands to logical operator", operator, ":", left.getDataType().String(), right.getDataType().String())
	}
	return NULL_POINTER
}

func evaluateDMAS(ctx *ScopeContext, node parser.TreeNode, operator string) *Pointer {
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
		ptr := malloc(type_sizes[NUMBER], true)
		ptr.setDataType(NUMBER)
		writeDataContent(ptr, numberByteArray(value))
		return ptr
	} else if left.getDataType() == STRING && right.getDataType() == STRING {
		if operator != "+" {
			interrupt(node.LineNo, "Invalid operator", operator, "for string operands")
		}
		leftVal := stringValue(left)
		rightVal := stringValue(right)
		newval := leftVal + rightVal
		ptr := malloc(type_sizes[CHAR]*len(newval), true)
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
			interrupt(node.LineNo, "string and number operands cannot be used with operator", operator)
		}
		ptr := malloc(type_sizes[CHAR]*len(newval), true)
		ptr.setDataType(STRING)
		writeDataContent(ptr, stringAsBytes(newval))
		return ptr
	} else if left.getDataType() == NUMBER && right.getDataType() == STRING {
		leftVal := numberValue(left)
		rightVal := stringValue(right)
		newval := ""
		switch operator {
		case "+":
			newval = fmt.Sprint(leftVal) + rightVal
		default:
			interrupt(node.LineNo, "number and string operands cannot be used with operator", operator)
		}
		ptr := malloc(type_sizes[CHAR]*len(newval), true)
		ptr.setDataType(STRING)
		writeDataContent(ptr, stringAsBytes(newval))
		return ptr

	} else {
		interrupt(node.LineNo, "invalid operands to arithmetic operator", operator, ":", left.getDataType().String(), right.getDataType().String())
	}
	return NULL_POINTER
}

func evaluateComparison(ctx *ScopeContext, node parser.TreeNode, operator string) *Pointer {
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
		ptr := malloc(type_sizes[BOOLEAN], true)
		ptr.setDataType(BOOLEAN)
		var val byte = 0
		if value {
			val = 1
		}
		writeDataContent(ptr, []byte{val})
		return ptr
	} else {
		interrupt(node.LineNo, "invalid operands to relational operator", operator, ":", left.getDataType().String(), right.getDataType().String())
	}
	return NULL_POINTER
}

func evaluateUnary(node parser.TreeNode, ctx *ScopeContext, operator string) *Pointer {
	pm := 1.0
	switch operator {
	case "--":
		pm = -1
		fallthrough
	case "++":
		varname := node.Children[0].Description
		varval, _ := findVariable(varname)
		if varval.isNull() {
			interrupt(node.LineNo, "cannot increment variable "+varname+" as it does not exist in current scope")
		}
		if varval.getDataType() != NUMBER {
			interrupt(node.LineNo, "cannot increment variable "+varname+" as it is not a number")
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
			ptr := malloc(type_sizes[NUMBER], true)
			ptr.setDataType(NUMBER)
			writeDataContent(ptr, numberByteArray(-numberValue(val)))
			return ptr

		default:
			interrupt(node.LineNo, "invalid unary operator "+operator)
		}
	} else if val.getDataType() == BOOLEAN {
		switch operator {
		case "!":
			memaddr := malloc(type_sizes[BOOLEAN], true)
			var valb byte = 1
			if !booleanValue(val) {
				valb = 0
			}
			writeDataContent(memaddr, []byte{valb})
			return memaddr
		default:
			interrupt(node.LineNo, "invalid unary operator "+operator)
		}
	}
	interrupt(node.LineNo, "invalid operand to unary operator", operator, ":", val.getDataType().String())
	return NULL_POINTER
}

// todo: use regex and simplify checks
// todo: this makes the whole language terribly slow!
// shift checking what kind of primary it is to the AST phase
func evaluatePrimary(node parser.TreeNode, ctx *ScopeContext) *Pointer {
	val := node.Description
	switch node.Label {
	case "number":
		ptr := malloc(type_sizes[NUMBER], true)
		ptr.setDataType(NUMBER)
		num := globals.NumMap[val]
		writeDataContent(ptr, num)
		return ptr
	case "boolean":
		ptr := malloc(type_sizes[BOOLEAN], true)
		ptr.setDataType(BOOLEAN)
		var boolnum byte = 0
		if utils.StringToBoolean(val) {
			boolnum = 1
		}
		writeDataContent(ptr, []byte{boolnum})
		return ptr
	case "string":

		ptr := malloc(type_sizes[CHAR]*(len(val)-2), true)
		ptr.setDataType(STRING)
		writeDataContent(ptr, stringAsBytes(val[1:len(val)-1]))
		return ptr
	}

	return NULL_POINTER
}

func evaluateVariable(node parser.TreeNode, ctx *ScopeContext) *Pointer {
	val := node.Description
	if v, _ := findVariable(val); !v.isNull() {
		if isCompositeType(v.getDataType()) {
			return v
		}
		cln := v.clone()
		return cln
	} else {
		interrupt(node.LineNo, "variable "+val+" does not exist in current scope")
	}
	return NULL_POINTER
}

func evaluateFuncCall(node parser.TreeNode, ctx *ScopeContext) *Pointer {

	function := findFunction(node.Description)
	funcNode := function
	if function == nil {
		interrupt(node.LineNo, "function "+node.Description+" does not exist in current scope")
	}
	newCtx := pushScopeContext(TYPE_FUNCTION, node.Description)
	actualArgs := node.Properties["args"]
	if len(actualArgs.Children) != len(funcNode.Properties["args"].Children) {
		interrupt(node.LineNo, "invalid number of arguments in function call "+funcNode.Description, "expected", len(funcNode.Properties["args"].Children), "got", len(actualArgs.Children))
	}
	debug_info("alloc for ", funcNode.Description)
	for i := 0; i < len(funcNode.Properties["args"].Children); i++ {
		argName := funcNode.Properties["args"].Children[i].Description
		argNode := actualArgs.Children[i]
		if argNode == nil {
			interrupt(node.LineNo, "missing argument "+argName+" in function call "+funcNode.Description)
		}
		argValue := evaluateExpression(argNode, newCtx)
		argValue.changeReferenceCount(true)
		argValue.temp = false
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
	debug_info("exited", funcNode.Description)
	return newCtx.returnValue
}

/*
structure of an array in memory:
bit 0: data type
bit 1-4: length (number of elements*sizeof(pointer type))
bit 5-8: address of first element
...
*/
func evaluateArray(node parser.TreeNode, ctx *ScopeContext) *Pointer {

	len := len(node.Children)
	arrptr := malloc(type_sizes[POINTER]*len, false)
	arrptr.setDataType(ARRAY)
	arrptr.setDataLength(len * type_sizes[POINTER])
	addr := arrptr.address + PTR_DATA_OFFSET
	for i := 0; i < len; i++ {
		ptri := evaluateExpression(node.Children[i], ctx)
		ptri.temp = false
		ptri.changeReferenceCount(true)
		unsafeWriteBytes(addr, intToBytes(ptri.address))
		addr += type_sizes[POINTER]
	}

	arrptr.temp = true

	return arrptr
}

/*
b[i] returns a pointer to the part of the array b which stores the address of the ith element of the array
*/
func evaluateArrayIndex(node parser.TreeNode, ctx *ScopeContext) int {
	varname := node.Children[0].Description

	ptr := evaluateExpression(node.Children[0], ctx)
	if ptr.getDataType() != ARRAY {
		interrupt(node.LineNo, "variable", varname, "is not an array")
		return NULL_POINTER.address
	}
	index := evaluateExpression(node.Properties["index"], ctx)
	if index.getDataType() != NUMBER {
		interrupt(node.LineNo, "array cannot be indexed by", index.getDataType().String())
		return NULL_POINTER.address
	}
	indexNo := int(numberValue(index))
	if indexNo >= ptr.getDataLength()/type_sizes[POINTER] || indexNo < 0 {
		interrupt(node.LineNo, "index", indexNo, "out of range for array length", ptr.getDataLength()/type_sizes[POINTER])
		return NULL_POINTER.address
	}
	addr := ptr.address + PTR_DATA_OFFSET + indexNo*type_sizes[POINTER]
	return addr
}

func evaluateObject(node parser.TreeNode, ctx *ScopeContext) *Pointer {
	numkeys := len(node.Children)
	objptr := malloc(type_sizes[POINTER]*numkeys*2, false)
	objptr.setDataType(OBJECT)
	objptr.setDataLength(numkeys * 2 * type_sizes[POINTER])
	addr := objptr.address + PTR_DATA_OFFSET
	for i := 0; i < numkeys; i++ {
		key := node.Children[i].Properties["key"].Description
		value := evaluateExpression(node.Children[i].Properties["value"], ctx)
		value.temp = false
		value.changeReferenceCount(true)
		unsafeWriteBytes(addr, intToBytes(globals.HashString(key)))
		unsafeWriteBytes(addr+type_sizes[POINTER], intToBytes(value.address))
		addr += type_sizes[POINTER] * 2
	}
	objptr.temp = true
	return objptr
}
