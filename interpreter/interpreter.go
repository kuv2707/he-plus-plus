package interpreter

import (
	"fmt"
	"strings"
	"toylingo/parser"
	"toylingo/utils"

	"github.com/gofrs/uuid"
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

type ScopeInfo struct {
	resumeIntoScope string
	returnVal       Variable
}

type Environment struct {
	variables map[string]Variable
	functions map[string]*parser.TreeNode
	callstack utils.Stack
}

func (env *Environment) findFromTop(name string) string {
	for i := env.callstack.Len() - 1; i >= 0; i-- {
		if strings.Contains(env.callstack.Get(i).(string), name) {
			return env.callstack.Get(i - 1).(string)
		}
	}
	return ""
}

//todo: call stack implement

// returns whether the execution of a scope was interrupted by a break/continue statement
func ExecuteAST(node *parser.TreeNode, env *Environment) ScopeInfo {
	//assign an id to this scope
	scopeId := node.Description + "_" + uuid.Must(uuid.NewV4()).String()
	env.callstack.Push(scopeId)

MAIN:
	for _, child := range node.Children {

		switch child.Label {
		case "operator":
			executeOperator(child, *env)
		case "primary":
			executePrimary(child, *env)
		case "scope":
			info := ExecuteAST(child, env)
			if !(info.resumeIntoScope == scopeId) {
				env.callstack.Pop()
				return info
			}
		case "IF":
			k := 0
			for ; child.Properties["condition"+fmt.Sprint(k)] != nil; k++ {
				treenode := child.Properties["condition"+fmt.Sprint(k)]
				verd := executeOperator(treenode, *env)
				if verd.value.(bool) {
					info := ExecuteAST(child.Children[k], env)
					if !(info.resumeIntoScope == scopeId) {
						env.callstack.Pop()
						return info
					}
					continue MAIN
				}
			}

			if k < len(child.Children) {
				// child.Children[k].PrintTree("-")
				info := ExecuteAST(child.Children[k], env)
				if !(info.resumeIntoScope == scopeId) {
						env.callstack.Pop()
						return info
					}
					continue MAIN
			}
		case "LOOP":
			treenode := child.Properties["condition"]
			for {
				verd := executeOperator(treenode, *env)
				if verd.value.(bool) {
					info := ExecuteAST(child.Children[0], env)
					if info.resumeIntoScope == ""{
						continue
					}
					if !(info.resumeIntoScope == scopeId) {
						env.callstack.Pop()
						return info
					} else {
						continue MAIN
					}
					
				} else {
					break
				}

			}
		case "BREAK":
			scp := env.findFromTop("loop_scope")
			// env.callstack.PrintStack()
			// fmt.Println("breaking to ", scp)
			return ScopeInfo{scp, Variable{"", DATA_TYPES["NULL"], "NULL"}}
		case "RETURN":
			returnVal := executeOperator(child.Children[0], *env)
			// fmt.Println("return",returnVal)
			return ScopeInfo{env.findFromTop("function_scope"), returnVal}
		case "FUNCTION":
			// fmt.Println("function",child.Properties["name"].Label)
			env.functions[child.Properties["name"].Label] = child
		}

	}
	env.callstack.Pop()
	
	return ScopeInfo{"", Variable{"", DATA_TYPES["NULL"], "NULL"}}
}

func Interpret(node *parser.TreeNode) {
	var env = Environment{make(map[string]Variable), make(map[string]*parser.TreeNode), utils.Stack{}}
	ExecuteAST(node, &env)
}

func NewCallStackContext(env Environment) Environment {
	e := Environment{make(map[string]Variable), env.functions, env.callstack}
	for k, v := range env.variables {
		e.variables[k] = v
	}
	for k, v := range env.functions {
		e.functions[k] = v
	}
	return e
}
