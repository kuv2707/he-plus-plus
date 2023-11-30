package interpreter

import (
	"fmt"
	"math"
	"toylingo/parser"
	"toylingo/utils"
)

func evaluateOperator(node parser.TreeNode, ctx scopeContext) Variable {
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
		return evaluateDMAS(ctx, node, node.Description[0])
	// case "<":
	// 	return evaluateLessThan(ctx, node)
	// case ">":
	// 	return evaluateGreaterThan(ctx, node)
	case "#":
		return evaluatePrint(ctx, node)
	// case "!":
	// 	return evaluateNot(ctx, node)
	// case "||":
	// 	return evaluateOr(ctx, node)
	// case "&&":
	// 	return evaluateAnd(ctx, node)
	// case "==":
	// 	return evaluateEquals(ctx, node)
	// case "!=":
	// 	return evaluateNotEquals(ctx, node)
	// case "<=":
	// 	return evaluateLessThanOrEquals(ctx, node)
	// case ">=":
	// 	return evaluateGreaterThanOrEquals(ctx, node)
	// case "++":
	// 	return evaluateIncrement(ctx, node)
	// case "--":
	// 	return evaluateDecrement(ctx, node)
	default:
		panic("invalid operator " + node.Description)
	}
}

func evaluateAssignment(ctx scopeContext, node parser.TreeNode) Variable {
	variableName := node.Children[0].Description
	variableValue := evaluateExpression(node.Children[1], ctx)
	ctx.variables[variableName] = variableValue
	fmt.Println(variableName, getNumber(variableValue))
	return variableValue
}

func evaluateDMAS(ctx scopeContext, node parser.TreeNode, operator byte) Variable {
	left := evaluateExpression(node.Children[0], ctx)
	right := evaluateExpression(node.Children[1], ctx)
	if left.vartype == "number" && right.vartype == "number" {
		value := math.NaN()
		switch operator {
		case '+':
			value = getValue(left).(float64) + getValue(right).(float64)
		case '-':
			value = getValue(left).(float64) - getValue(right).(float64)
		case '*':
			value = getValue(left).(float64) * getValue(right).(float64)
		case '/':
			value = getValue(left).(float64) / getValue(right).(float64)
		}
		fmt.Println("DMAS", fmt.Sprintf("%c", operator), value)
		//convert value to IEEE 754 format 64-bit floating point number and store in HEAP by a malloc call -> a pointer is returned -> return a Variable with that pointer
		memaddr := malloc(8)
		writeBits(memaddr, int64(math.Float64bits(value)), 8)
		return Variable{memaddr, "number"}

	} else {
		panic("invalid operands to binary +")

	}
}

func evaluatePrint(ctx scopeContext, node parser.TreeNode) Variable {
	value := evaluateExpression(node.Children[0], ctx)
	fmt.Print( utils.Colors["GREEN"], getNumber(value),utils.Colors["RESET"])
	return value
}

func parsePrimary(node parser.TreeNode, ctx scopeContext) Variable {
	val := node.Description
	if utils.IsNumber(val) {
		memaddr := malloc(8)
		writeBits(memaddr, int64(math.Float64bits(utils.StringToNumber(val))), 8)
		return Variable{memaddr, "number"}
	} else {
		return ctx.variables[val]
	}
}

func evaluateExpression(node *parser.TreeNode, ctx scopeContext) Variable {
	switch node.Label {
	case "operator":
		return evaluateOperator(*node, ctx)
	case "literal":
		fallthrough
	case "primary":
		return parsePrimary(*node, ctx)
	default:
		panic("invalid expression " + node.Label)
	}
}
