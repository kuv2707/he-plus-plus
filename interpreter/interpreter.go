package interpreter

import (
	"fmt"
	"toylingo/parser"
)

func Interpret(root *parser.TreeNode) {
	executeScope(root, pushScopeContext("scope_0"))
	printMemoryStats()
}

func executeScope(node *parser.TreeNode, ctx *scopeContext) {
	fmt.Println("entered", ctx.scopeType)
	// printVariableList(ctx.variables)
SCOPE_EXECUTION:
	for i := range node.Children {
		child := node.Children[i]
		switch child.Label {
		case "function":
			ctx.functions[child.Description] = *child

		case "scope":
			executeScope(child, pushScopeContext("scope"))
		case "loop":
			for true {
				variable := evaluateExpressionClean(child.Properties["condition"], ctx)
				if getNumber(variable) == 0 {
					break
				}
				executeScope(child.Properties["body"], pushScopeContext("loop"))
			}
		case "operator":
			fallthrough
		case "call":
			evaluateExpressionClean(child, ctx)
		// case "conditional_block":

		case "return":
			// fmt.Println("---eval ret val")
			// printVariableList(ctx.variables)
			expr := evaluateExpression(child.Children[0], ctx)
			expr.pointer.temp = false
			// fmt.Println("returning", expr.pointer, getNumber(expr))
			ctx.returnValue = &expr
			break SCOPE_EXECUTION
		default:
			fmt.Println("__unknown", child.Label)
		}

	}
	printMemoryStats()
	popScopeContext()
}

func evaluateExpressionClean(node *parser.TreeNode, ctx *scopeContext) Variable {
	variable := evaluateExpression(node, ctx)
	variable.pointer.temp = false
	// fmt.Println("value", getNumber(variable), variable.pointer)
	gc()
	return variable
}
