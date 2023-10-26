package interpreter

import (
	_ "fmt"
	"toylingo/parser"
)

var DATA_TYPES = map[string]string{
	"NUMBER": "NUMBER",
	"STRING": "STRING",
	"BOOL":   "BOOL",
	"NULL":   "NULL",
}

type Variable struct {
	name     string
	datatype string
	value    interface{}
}

type Environment struct {
	variables map[string]Variable
}

var env = Environment{make(map[string]Variable)}

func ExecuteAST(node *parser.TreeNode) {

	for _, child := range node.Children {

		switch child.Label {
		case "operator":
			executeOperator(child, env)

		case "if":
			cond := executeOperator(child.Properties["condition"], env)
			if cond.value.(bool) {
				ExecuteAST(child)
			}
		}

	}
}
