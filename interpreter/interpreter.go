package interpreter

import "toylingo/parser"

func Interpret(root *parser.TreeNode){
	executeScope(root,pushScopeContext("root"))
}

func executeScope(node *parser.TreeNode,ctx scopeContext ) {
	
	for i :=range node.Children{
		child:=node.Children[i]
		switch child.Label{
		case "function":
			ctx.functions[child.Description]=*child
		case "loop":
			for getNumber(evaluateExpression(child.Children[0],ctx))!=0{
				executeScope(child.Children[1],pushScopeContext("loop"))
			}
		case "operator":
			evaluateOperator(*child,ctx)
		// case "conditional_block":

		}
	}
}

