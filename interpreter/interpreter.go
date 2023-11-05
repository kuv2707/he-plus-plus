package interpreter

import (
	"fmt"
	_ "fmt"
	"toylingo/parser"
)

var DATA_TYPES = map[string]string{
	"NUMBER": "NUMBER",
	"STRING": "STRING",
	"BOOL":   "BOOL",
}

type Variable struct {
	name     string
	datatype string
	value    interface{}
}

type Environment struct {
	variables map[string]Variable
}
//returns whether the execution of a scope was interrupted by a break/continue statement
func ExecuteAST(node *parser.TreeNode, env *Environment) bool{

MAIN:
	for _, child := range node.Children {

		switch child.Label {
		case "operator":
			executeOperator(child, *env)
		case "scope":
			ExecuteAST(child,env)
		case "IF":
			k := 0
			for ; child.Properties["condition"+fmt.Sprint(k)] != nil; k++ {
				treenode:=child.Properties["condition"+fmt.Sprint(k)]
				verd := executeOperator(treenode, *env)
				
				if verd.value.(bool) {
					ExecuteAST(child.Children[k], env)
					continue MAIN
				}
			}
			
			if k < len(child.Children) {
				child.Children[k].PrintTree("-")
				ExecuteAST(child.Children[k], env)
			}
		case "LOOP":
			treenode:=child.Properties["condition"]
			for {
				verd := executeOperator(treenode, *env)
				if verd.value.(bool) {
					vvv:=ExecuteAST(child.Children[0], env)
					if !vvv{
						break
					}
				}else{
					break
				}

			}
		case "BREAK":
			fmt.Println("breaking")
			return false
		}


	}
	return true
}

func Interpret(node *parser.TreeNode) {
	var env = Environment{make(map[string]Variable)}
	ExecuteAST(node, &env)
}
