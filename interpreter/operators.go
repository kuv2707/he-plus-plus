package interpreter

import (
	"fmt"
	"math"
	"toylingo/parser"
	"toylingo/utils"
)

func evaluateOperator(node parser.TreeNode, ctx *scopeContext) Variable {
	if node.Label == "literal" {
		return parsePrimary(node, ctx)

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

func evaluateAssignment(ctx *scopeContext, node parser.TreeNode) Variable {
	variableName := node.Children[0].Description
	variableValue := evaluateExpression(node.Children[1], ctx)
	val, ok := ctx.variables[variableName]
	if ok {
		fmt.Println(">>>", variableName, val.pointer, getNumber(variableValue))
		writeBits(val.pointer, int64(math.Float64bits(getNumber(variableValue))), 8)
		return val

	}
	ctx.variables[variableName] = variableValue
	ctx.inScopeVars = append(ctx.inScopeVars, variableName)
	//print inscopevars
	fmt.Println(ctx.scopeType, "inScopeVars", ctx.inScopeVars)
	fmt.Println(">>", variableName, variableValue.pointer, getNumber(variableValue))
	return variableValue
}

func evaluateDMAS(ctx *scopeContext, node parser.TreeNode, operator string) Variable {
	left := evaluateExpression(node.Children[0], ctx)
	right := evaluateExpression(node.Children[1], ctx)
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
		memaddr := malloc(8, ctx.scopeType,true)
		writeBits(memaddr, int64(math.Float64bits(value)), 8)
		fmt.Println("want to free", left.pointer, right.pointer)
		// left or right may be mapped to a variable in ctx.variables, so we must not free them without checking
		// possible solution: parsePrimary should return a copy of the variable, not the variable itself
		freePtr(left.pointer)
		freePtr(right.pointer)
		return Variable{memaddr, "number"}

	} else {
		panic("invalid operands to binary operator " + operator)

	}
}

func evaluatePrint(ctx *scopeContext, node parser.TreeNode) Variable {
	value := evaluateExpression(node.Children[0], ctx)
	fmt.Print(utils.Colors["GREEN"], getNumber(value), utils.Colors["RESET"])
	return value
}

func parsePrimary(node parser.TreeNode, ctx *scopeContext) Variable {
	val := node.Description
	if utils.IsNumber(val) {
		memaddr := malloc(8, ctx.scopeType,true)
		writeBits(memaddr, int64(math.Float64bits(utils.StringToNumber(val))), 8)
		return Variable{memaddr, "number"}
	} else {
		//if val is not a key in ctx.variables, it returns {0,0} why?
		// fmt.Println("parsePrimary", val, ctx.variables[val], getNumber(ctx.variables[val]))
		return copyVariable(ctx.variables[val])
	}
}
//todo: garbage collect after every expression evaluation
func evaluateExpression(node *parser.TreeNode, ctx *scopeContext) Variable {
	// defer gc()
	switch node.Label {
	case "operator":
		ret:= evaluateOperator(*node, ctx)
		gc()
		return ret
	case "literal":
		fallthrough
	case "primary":
		ret:= parsePrimary(*node, ctx)
		gc()
		return ret
	default:
		panic("invalid expression " + node.Label)
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
		fmt.Println("COMP", operator, value)
		//convert value to IEEE 754 format 64-bit floating point number and store in HEAP by a malloc call -> a pointer is returned -> return a Variable with that pointer
		memaddr := malloc(8, ctx.scopeType,true)
		fmt.Println("new allocated address", memaddr)
		val := 0.0
		if value {
			val = 1.0
		}
		writeBits(memaddr, int64(math.Float64bits(val)), 8)
		//todo: freeing left and right causes bugs
		fmt.Println("want to free", left.pointer, right.pointer)
		freePtr(left.pointer)
		freePtr(right.pointer)
		return Variable{memaddr, "number"}

	} else {
		panic("invalid operands to binary operator " + operator)

	}
}
