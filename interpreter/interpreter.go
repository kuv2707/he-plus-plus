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

type ScopeInfo struct{
	shouldContinue bool
	returnVal Variable
}

type Environment struct {
	variables map[string]Variable
	functions map[string]*parser.TreeNode


}

//todo: call stack implement


//returns whether the execution of a scope was interrupted by a break/continue statement
func ExecuteAST(node *parser.TreeNode, env *Environment) ScopeInfo{

MAIN:
	for _, child := range node.Children {

		switch child.Label {
		case "operator":
			executeOperator(child, *env)
		case "scope":
			info:=ExecuteAST(child,env)
			if !info.shouldContinue{
				return info
			}
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
				// child.Children[k].PrintTree("-")
				ExecuteAST(child.Children[k], env)
			}
		case "LOOP":
			treenode:=child.Properties["condition"]
			for {
				verd := executeOperator(treenode, *env)
				if verd.value.(bool) {
					vvv:=ExecuteAST(child.Children[0], env)
					if !vvv.shouldContinue{
						break
					}
				}else{
					break
				}

			}
		case "BREAK":
			// fmt.Println("breaking")
			return ScopeInfo{false,Variable{"",DATA_TYPES["NULL"],"NULL"}}
		case "RETURN":
			returnVal:=executeOperator(child.Children[0],*env)
			// fmt.Println("return",returnVal)
			return ScopeInfo{false,returnVal}
		case "FUNCTION":
			// fmt.Println("function",child.Properties["name"].Label)
			env.functions[child.Properties["name"].Label] = child
		}


	}
	return ScopeInfo{true,Variable{"",DATA_TYPES["NULL"],"NULL"}}
}

func Interpret(node *parser.TreeNode) {
	var env = Environment{make(map[string]Variable), make(map[string]*parser.TreeNode)}
	ExecuteAST(node, &env)
}


func NewCallStackContext(env Environment) Environment{
	e:= Environment{make(map[string]Variable), env.functions}
	for k,v:=range env.variables{
		e.variables[k]=v
	}
	for k,v:=range env.functions{
		e.functions[k]=v
	}
	return e
}