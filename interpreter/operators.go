package interpreter

import (
	"fmt"
	"he++/parser"
	"he++/utils"
	"math"
)

func evaluateOperator(node parser.TreeNode, ctx *ScopeContext) *Pointer {
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

func evaluateExpression(node *parser.TreeNode, ctx *ScopeContext) *Pointer {
	LineNo = node.LineNo
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
		ret = evaluatePrimary(*node, ctx)
	case "array":
		ret = evaluateArray(*node, ctx)
	case "index":
		ind := evaluateArrayIndex(*node, ctx)
		ret = pointers[bytesToInt(heapSlice(ind, type_sizes[POINTER]))]
	case "call":
		ret = evaluateFuncCall(*node, ctx)
		if ret == NULL_POINTER {
			interrupt("function " + node.Description + " does not return a value but is expected to")
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
		// fmt.Println(stringValue(pointers[bytesToInt(heapSlice(addr, 4))]))
		// variableValue.print()
		variableValue.temp = false
		unsafeWriteBytes(addr, intToBytes(variableValue.address))
		// fmt.Println(stringValue(pointers[bytesToInt(heapSlice(addr, 4))]))
		return variableValue
	}
	variableName := node.Children[0].Description
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
		interrupt("invalid operands to logical operator", operator, ":", left.getDataType().String(), right.getDataType().String())
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
			interrupt("Invalid operator", operator, "for string operands")
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
			interrupt("string and number operands cannot be used with operator", operator)
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
			interrupt("number and string operands cannot be used with operator", operator)
		}
		ptr := malloc(type_sizes[CHAR]*len(newval), true)
		ptr.setDataType(STRING)
		writeDataContent(ptr, stringAsBytes(newval))
		return ptr

	} else {
		interrupt("invalid operands to arithmetic operator", operator, ":", left.getDataType().String(), right.getDataType().String())
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
		interrupt("invalid operands to relational operator", operator, ":", left.getDataType().String(), right.getDataType().String())
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
			ptr := malloc(type_sizes[NUMBER], true)
			ptr.setDataType(NUMBER)
			writeDataContent(ptr, numberByteArray(-numberValue(val)))
			return ptr

		default:
			interrupt("invalid unary operator " + operator)
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
			interrupt("invalid unary operator " + operator)
		}
	}
	interrupt("invalid operand to unary operator", operator, ":", val.getDataType().String())
	return NULL_POINTER
}

// todo: use regex and simplify checks
// todo: this makes the whole language terribly slow!
// shift checking what kind of primary it is to the AST phase
func evaluatePrimary(node parser.TreeNode, ctx *ScopeContext) *Pointer {
	LineNo = node.LineNo
	val := node.Description
	switch node.Label {
	case "number":
		ptr := malloc(type_sizes[NUMBER], true)
		ptr.setDataType(NUMBER)
		writeDataContent(ptr, numberByteArray(StringToNumber(val)))
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
	case "variable":
		if v := findVariable(val); !v.isNull() {
			/*
				array also gets copied, but the elements referred to in it, are not. so if a is an array and we do b=a, then make changes
			*/
			if v.getDataType() == ARRAY {
				return v
			}
			return v.clone()
		} else {
			LineNo = node.LineNo
			interrupt("variable " + val + " does not exist in current scope")
		}
	}

	// if isCompositeDS(node) {
	// 	//func call or array or object
	// 	return evaluateCompositeDS(node, ctx)
	// }

	return NULL_POINTER
}

func evaluateFuncCall(node parser.TreeNode, ctx *ScopeContext) *Pointer {

	function := findFunction(node.Description)
	funcNode := function
	if function == nil {
		interrupt("function " + node.Description + " does not exist in current scope")
	}
	newCtx := pushScopeContext(TYPE_FUNCTION, node.Description)
	lastValidLine := LineNo
	actualArgs := node.Properties["args"]
	if len(actualArgs.Children) != len(funcNode.Properties["args"].Children) {
		interrupt("invalid number of arguments in function call "+funcNode.Description, "expected", len(funcNode.Properties["args"].Children), "got", len(actualArgs.Children))
	}
	for i := 0; i < len(funcNode.Properties["args"].Children); i++ {
		argName := funcNode.Properties["args"].Children[i].Description
		argNode := actualArgs.Children[i]
		if argNode == nil {
			LineNo = lastValidLine
			interrupt("missing argument " + argName + " in function call " + funcNode.Description)
		}
		lastValidLine = argNode.LineNo
		argValue := evaluateExpression(argNode, newCtx)
		argValue.temp = false
		newCtx.variables[argName] = argValue.clone()
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

/*
structure of an array in memory:
bit 0: data type
bit 1-4: length (number of elements*sizeof(pointer type))
bit 5-8: address of first element
...
*/
func evaluateArray(node parser.TreeNode, ctx *ScopeContext) *Pointer {
	LineNo = node.LineNo
	len := len(node.Children)
	arrptr := malloc(type_sizes[POINTER]*len, false)
	arrptr.setDataType(ARRAY)
	arrptr.setDataLength(len * type_sizes[POINTER])
	addr := arrptr.address + PTR_DATA_OFFSET
	for i := 0; i < len; i++ {
		ptri := evaluateExpression(node.Children[i], ctx)
		ptri.temp = false
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
	LineNo = node.LineNo
	ptr := evaluateExpression(node.Children[0], ctx)
	if ptr.getDataType() != ARRAY {
		interrupt("variable", varname, "is not an array")
		return NULL_POINTER.address
	}
	index := evaluateExpression(node.Properties["index"], ctx)
	if index.getDataType() != NUMBER {
		interrupt("array cannot be indexed by", index.getDataType().String())
		return NULL_POINTER.address
	}
	indexNo := int(numberValue(index))
	if indexNo >= ptr.getDataLength()/type_sizes[POINTER] || indexNo < 0 {
		interrupt("index", indexNo, "out of range for array length", ptr.getDataLength()/type_sizes[POINTER])
		return NULL_POINTER.address
	}
	addr := ptr.address + PTR_DATA_OFFSET + indexNo*type_sizes[POINTER]
	return addr
}
