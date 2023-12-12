package interpreter

import (
	"fmt"
	"toylingo/parser"
)

func Interpret(root *parser.TreeNode) {
	executeScope(root, pushScopeContext("scope_0"), 0)
	printMemoryStats()
}

func executeScope(node *parser.TreeNode, ctx *scopeContext, depth int) {
	fmt.Println("entered", ctx.scopeType)
	for variable:=range ctx.variables{
		fmt.Println("variable",variable,ctx.variables[variable],getNumber(ctx.variables[variable]))
	}
	for i := range node.Children {
		child := node.Children[i]
		switch child.Label {
		case "function":
			ctx.functions[child.Description] = *child

		case "scope":
			executeScope(child, pushScopeContext(fmt.Sprintf("scope_%d", depth+1)), depth+1)
		case "loop":
			for true{
				variable := evaluateExpressionClean(child.Properties["condition"], ctx)
				if getNumber(variable)==0{
					break
				}
				executeScope(child.Properties["body"], pushScopeContext(fmt.Sprintf("loop_%d", depth+1)), depth+1)
			}
		case "operator":
			evaluateExpressionClean(child, ctx)
			// case "conditional_block":

		}
	}
	println("exiting ", ctx.scopeType)
	printMemoryStats()
	popScopeContext()
}


func evaluateExpressionClean(node *parser.TreeNode, ctx *scopeContext) Variable {
	variable := evaluateExpression(node, ctx)
	variable.pointer.temp = false
	fmt.Println("value", getNumber(variable), variable.pointer)
	gc()
	return variable
}