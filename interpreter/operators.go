package interpreter

import (
	"fmt"
	"math"
	"toylingo/parser"
	"toylingo/utils"
)

func evaluateOperator(node parser.TreeNode, ctx *scopeContext) Variable {
	if node.Label == "literal" {
		return evaluatePrimary(node, ctx)

	}
	switch node.Description {
	case "=":
		return evaluateAssignment(ctx, node)
	case "+":
		fallthrough
	case "-":
		fallthrough
	case "*":
		fallthrough
	case "/":
		return evaluateDMAS(ctx, node, node.Description)
	case "<":
		fallthrough
	case ">":
		fallthrough
	case "<=":
		fallthrough
	case ">=":
		fallthrough
	case "==":
		return evaluateComparison(ctx, node, node.Description)
	case "#":
		return evaluatePrint(ctx, node)

	}
	panic("invalid operator " + node.Description)
}

// todo: garbage collect after every expression evaluation
func evaluateExpression(node *parser.TreeNode, ctx *scopeContext) Variable {
	switch node.Label {
	case "operator":
		ret := evaluateOperator(*node, ctx)
		return ret
	case "literal":
		fallthrough
	case "primary":
		ret := evaluatePrimary(*node, ctx)
		return ret
	case "call":
		ret := evaluateFuncCall(*node, ctx)
		return ret
	default:
		node.PrintTree("")
		panic("invalid expression " + node.Label)
	}
}
func evaluateAssignment(ctx *scopeContext, node parser.TreeNode) Variable {
	variableName := node.Children[0].Description
	variableValue := evaluateExpression(node.Children[1], ctx)
	val, ok := ctx.variables[variableName]
	if ok {
		// fmt.Println("updating existing var",variableName)
		// fmt.Println(">>>", variableName, val.pointer, getNumber(variableValue))
		//todo: make a function to copy value from one pointer to another
		writeBits(*val.pointer, int64(math.Float64bits(getNumber(variableValue))), 8)
		return val

	}
	ctx.variables[variableName] = variableValue
	variableValue.pointer.temp = false

	return variableValue
}

func evaluateDMAS(ctx *scopeContext, node parser.TreeNode, operator string) Variable {
	// fmt.Println("eval DMAS", node.Description)
	left := evaluateExpression(node.Children[0], ctx)
	right := evaluateExpression(node.Children[1], ctx)
	// fmt.Println("DMAS args: ", getNumber(left), getNumber(right))
	if left.vartype == "number" && right.vartype == "number" {
		value := math.NaN()
		switch operator {
		case "+":
			value = getValue(left).(float64) + getValue(right).(float64)
		case "-":
			value = getValue(left).(float64) - getValue(right).(float64)
		case "*":
			value = getValue(left).(float64) * getValue(right).(float64)
		case "/":
			value = getValue(left).(float64) / getValue(right).(float64)
		}
		// fmt.Println("DMAS", operator, value)
		memaddr := malloc(8, ctx.scopeType, true)
		writeBits(*memaddr, int64(math.Float64bits(value)), 8)
		// fmt.Println("want to free", left.pointer, right.pointer)
		// left or right may be mapped to a variable in ctx.variables, so we must not free them without checking
		// possible solution: evaluatePrimary should return a copy of the variable, not the variable itself
		// freePtr(left.pointer)
		// freePtr(right.pointer)
		return Variable{memaddr, "number"}

	} else {
		panic("invalid operands to binary operator " + operator)

	}
}

func evaluatePrint(ctx *scopeContext, node parser.TreeNode) Variable {
	// node.PrintTree("")
	value := evaluateExpression(node.Children[0], ctx)
	fmt.Println("printval", value.pointer)
	fmt.Print(utils.Colors["GREEN"], getNumber(value), utils.Colors["RESET"])
	return value
}

func evaluatePrimary(node parser.TreeNode, ctx *scopeContext) Variable {
	val := node.Description
	if len(node.Children) > 0 {
		//func call or array or object
		return evaluateCompositeDS(node, ctx)
	}
	if utils.IsNumber(val) {
		memaddr := malloc(8, ctx.scopeType, true)
		writeBits(*memaddr, int64(math.Float64bits(utils.StringToNumber(val))), 8)
		return Variable{memaddr, "number"}
	} else {
		//if val is not a key in ctx.variables, it returns {0,0} why?
		// fmt.Println("evaluatePrimary", val, ctx.variables[val], getNumber(ctx.variables[val]))
		copy := copyVariable(ctx.variables[val])
		// fmt.Println("evaluatePrimary", val, copy, getNumber(copy))
		return copy
	}
}

func evaluateComparison(ctx *scopeContext, node parser.TreeNode, operator string) Variable {
	left := evaluateExpression(node.Children[0], ctx)
	right := evaluateExpression(node.Children[1], ctx)
	if left.vartype == "number" && right.vartype == "number" {
		value := false
		switch operator {
		case "<":
			value = getValue(left).(float64) < getValue(right).(float64)
		case ">":
			value = getValue(left).(float64) > getValue(right).(float64)
		case "<=":
			value = getValue(left).(float64) <= getValue(right).(float64)
		case ">=":
			value = getValue(left).(float64) >= getValue(right).(float64)
		case "==":
			value = getValue(left).(float64) == getValue(right).(float64)
		case "!=":
			value = getValue(left).(float64) != getValue(right).(float64)

		}
		// fmt.Println("COMP", operator, value)
		//convert value to IEEE 754 format 64-bit floating point number and store in HEAP by a malloc call -> a pointer is returned -> return a Variable with that pointer
		memaddr := malloc(8, ctx.scopeType, true)
		val := 0.0
		if value {
			val = 1.0
		}
		writeBits(*memaddr, int64(math.Float64bits(val)), 8)
		//todo: freeing left and right causes bugs
		// fmt.Println("want to free", left.pointer, right.pointer)
		// freePtr(left.pointer)
		// freePtr(right.pointer)
		return Variable{memaddr, "number"}

	} else {
		panic("invalid operands to binary operator " + operator)

	}
}

// functions are not scoped as of now
// if a function is defined in a previously executed scope, it can be called in the current scope, even though that scope has been popped from the stack
func evaluateFuncCall(node parser.TreeNode, ctx *scopeContext) Variable {
	funcNode := ctx.functions[node.Description]
	newCtx := pushScopeContext("function")
	for i := 0; i < len(funcNode.Properties["args"].Children); i++ {
		argName := funcNode.Properties["args"].Children[i].Description
		argValue := evaluateExpression(node.Properties["args"+fmt.Sprint(i)], ctx)
		argValue.pointer.temp = false
		newCtx.variables[argName] = argValue
	}
	fmt.Println("calling", funcNode.Description)
	executeScope(funcNode.Properties["body"], newCtx)
	return *newCtx.returnValue
}

func evaluateCompositeDS(node parser.TreeNode, ctx *scopeContext) Variable {
	// switch node.Description{
	// case "array":
	// 	return evaluateArray(node, ctx)
	// }
	panic("invalid composite data structure")

}
//returns a Variable with pointer to an array
// func evaluateArray(node parser.TreeNode, ctx *scopeContext) Variable {
// 	len:=len(node.Children)
	
// }