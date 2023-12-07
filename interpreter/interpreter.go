package interpreter

import (
	"fmt"
	"toylingo/parser"
)

func Interpret(root *parser.TreeNode){
	executeScope(root,pushScopeContext("scope_0"),0)
}

func executeScope(node *parser.TreeNode,ctx *scopeContext,depth int ) {
	
	for i :=range node.Children{
		child:=node.Children[i]
		switch child.Label{
		case "function":
			ctx.functions[child.Description]=*child

		case "scope":
			executeScope(child,pushScopeContext(fmt.Sprintf("scope_%d",depth+1)),depth+1)
		case "loop":
			for getNumber(evaluateExpression(child.Properties["condition"],ctx))!=0{
				executeScope(child.Properties["body"],pushScopeContext(fmt.Sprintf("loop_%d",depth+1)),depth+1)
			}
		case "operator":
			evaluateOperator(*child,ctx)
		// case "conditional_block":

		}
	}
	println("exiting ",ctx.scopeType)
	
	popScopeContext()
}

