package interpreter

import (
	"fmt"
	"toylingo/parser"
	"toylingo/utils"
)

func evaluateOperator(node parser.TreeNode, ctx *scopeContext) Variable {
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
	} else if node.Description == "!" {
		return evaluateUnary(node, ctx, node.Description)
	}
	interrupt("invalid operator " + node.Description)
	return Variable{}
}

func evaluateExpression(node *parser.TreeNode, ctx *scopeContext) Variable {
	switch node.Label {
	case "operator":
		return evaluateOperator(*node, ctx)
	case "literal":
		fallthrough
	case "primary":
		return evaluatePrimary(*node, ctx)
	case "call":
		ret := evaluateFuncCall(*node, ctx)
		if ret == nil {
			interrupt("function " + node.Description + " does not return a value but is expected to")
		}
		return *ret
	default:
		node.PrintTree("->")
		panic("invalid expression " + node.Label)
	}
}

func evaluateAssignment(ctx *scopeContext, node parser.TreeNode) Variable {
	variableName := node.Children[0].Description
	variableValue := evaluateExpression(node.Children[1], ctx)
	val, alreadyExists := ctx.variables[variableName]
	if alreadyExists {
		//todo: make a function to copy value from one pointer to another: memcpy
		writeBits(*val.pointer, numberByteArray(getNumber(variableValue)))
		return val
	}
	ctx.variables[variableName] = variableValue
	variableValue.pointer.temp = false

	return variableValue
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
	if len(node.Children) > 0 {
		//func call or array or object
		return evaluateCompositeDS(node, ctx)
	}
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
	} else {
		if v, exists := ctx.variables[val]; exists {
			if v.vartype == TYPE_ARRAY {
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

func evaluateFuncCall(node parser.TreeNode, ctx *scopeContext) *Variable {

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
		argValue.pointer.scopeId = newCtx.scopeId
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
	return &ret
}

func evaluateCompositeDS(node parser.TreeNode, ctx *scopeContext) Variable {
	switch node.Description {
	case "array":
		return evaluateArray(node, ctx)
	}
	panic("invalid composite data structure")

}

// returns a Variable with pointer to an array
// the arrangement of an array variable is as follows:
// 1. the first 8 bytes contain the length of the array
// 2. the next 8*n bytes contain the pointer to the first element of the array, or actual elements if it is a primitive

func evaluateArray(node parser.TreeNode, ctx *scopeContext) Variable {
	len := len(node.Children)
	memaddr := malloc(type_sizes[TYPE_POINTER]*len+ type_sizes[TYPE_NUMBER], ctx.scopeId, true)
	addr:=memaddr.address
	unsafeWriteBits(addr, numberByteArray(float64(len)))
	addr+=type_sizes[TYPE_NUMBER]
	for i:=range node.Children{
		val := evaluateExpression(node.Children[i], ctx)
		val.pointer.temp = false
		unsafeWriteBits(addr, pointerByteArray(val.pointer.address))
		addr+=type_sizes[TYPE_POINTER]
	}
	return Variable{memaddr, TYPE_ARRAY}
}
