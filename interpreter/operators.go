package interpreter

import (
	// "fmt"
	"fmt"
	"toylingo/parser"
	"toylingo/utils"
)

func executeOperator(node *parser.TreeNode, env Environment) Variable {
	// fmt.Println(node.Description)

	if node.Label == "primary" {
		return executePrimary(node, env)
	}

	if node.Description == "=" {
		return executeAssignment(node, env)
	} else if node.Description == "+" {
		return executeAddition(node, env)
	} else if node.Description == "-" {
		return executeSubtraction(node, env)
	} else if node.Description == "*" {
		return executeMultiplication(node, env)
	} else if node.Description == "/" {
		return executeDivision(node, env)
	} else if node.Description == ">" {
		return executeGreaterThan(node, env)
	} else if node.Description == "<" {
		return executeLessThan(node, env)
	} else if node.Description == "#" {
		return executePrint(node, env)
	} else if node.Description == "==" {
		return executeEquals(node, env)
	} else if node.Description == "!=" {
		res:=executeEquals(node,env)
		res.value=!res.value.(bool)
		return res
	} else if node.Description == "<=" {
	} else if node.Description == ">=" {
	}
	return Variable{"", DATA_TYPES["NULL"], "NULL"}
}

func executeAssignment(node *parser.TreeNode, env Environment) Variable {
	varname := node.Children[0].Label
	varval := executeOperator(node.Children[1], env)
	env.variables[varname] = varval

	return varval
}
func executeAddition(node *parser.TreeNode, env Environment) Variable {
	left := executeOperator(node.Children[0], env)
	right := executeOperator(node.Children[1], env)
	if left.datatype == DATA_TYPES["NUMBER"] && right.datatype == DATA_TYPES["NUMBER"] {
		return Variable{"", DATA_TYPES["NUMBER"], left.value.(int) + right.value.(int)}
	} else if left.datatype == DATA_TYPES["STRING"] && right.datatype == DATA_TYPES["STRING"] {
		return Variable{"", DATA_TYPES["STRING"], left.value.(string) + right.value.(string)}
	} else {
		//todo
		return Variable{"", DATA_TYPES["NUMBER"], 0}
	}
}
func executeSubtraction(node *parser.TreeNode, env Environment) Variable {
	left := executeOperator(node.Children[0], env)
	right := executeOperator(node.Children[1], env)
	if left.datatype == DATA_TYPES["NUMBER"] && right.datatype == DATA_TYPES["NUMBER"] {
		return Variable{"", DATA_TYPES["NUMBER"], left.value.(int) - right.value.(int)}
	} else {
		//todo
		return Variable{"", DATA_TYPES["NUMBER"], 0}
	}
}
func executeMultiplication(node *parser.TreeNode, env Environment) Variable {
	left := executeOperator(node.Children[0], env)
	right := executeOperator(node.Children[1], env)
	if left.datatype == DATA_TYPES["NUMBER"] && right.datatype == DATA_TYPES["NUMBER"] {
		return Variable{"", DATA_TYPES["NUMBER"], left.value.(int) * right.value.(int)}
	} else {
		//todo
		return Variable{"", DATA_TYPES["NUMBER"], 0}
	}
}
func executeDivision(node *parser.TreeNode, env Environment) Variable {
	left := executeOperator(node.Children[0], env)
	right := executeOperator(node.Children[1], env)
	if left.datatype == DATA_TYPES["NUMBER"] && right.datatype == DATA_TYPES["NUMBER"] {
		return Variable{"", DATA_TYPES["NUMBER"], left.value.(int) / right.value.(int)}
	} else {
		//todo
		return Variable{"", DATA_TYPES["NUMBER"], 0}
	}
}
func executeGreaterThan(node *parser.TreeNode, env Environment) Variable {
	left := executeOperator(node.Children[0], env)
	right := executeOperator(node.Children[1], env)
	if left.datatype == DATA_TYPES["NUMBER"] && right.datatype == DATA_TYPES["NUMBER"] {
		return Variable{"", DATA_TYPES["BOOL"], left.value.(int) > right.value.(int)}
	} else {
		//todo
		return Variable{"", DATA_TYPES["BOOL"], false}
	}
}
func executeLessThan(node *parser.TreeNode, env Environment) Variable {
	left := executeOperator(node.Children[0], env)
	right := executeOperator(node.Children[1], env)
	if left.datatype == DATA_TYPES["NUMBER"] && right.datatype == DATA_TYPES["NUMBER"] {
		return Variable{"", DATA_TYPES["BOOL"], left.value.(int) < right.value.(int)}
	} else {
		//todo
		return Variable{"", DATA_TYPES["BOOL"], false}
	}
}

func executeEquals(node *parser.TreeNode, env Environment) Variable {
	left := executeOperator(node.Children[0], env)
	right := executeOperator(node.Children[1], env)
	if left.datatype == DATA_TYPES["NUMBER"] && right.datatype == DATA_TYPES["NUMBER"] {
		return Variable{"", DATA_TYPES["BOOL"], left.value.(int) == right.value.(int)}
	} else {
		//todo
		return Variable{"", DATA_TYPES["BOOL"], false}
	}
}

func executePrimary(node *parser.TreeNode, env Environment) Variable {
	if node.Description == "true" || node.Description == "false" {
		return Variable{"", DATA_TYPES["BOOL"], node.Description == "true"}
	} else if utils.IsNumber(node.Description) {
		return Variable{"", DATA_TYPES["NUMBER"], utils.StringToInt(node.Description)}
	} else {
		if utils.InQuotes(node.Description){
			return Variable{"",DATA_TYPES["STRING"],node.Description}
		} else{

			return env.variables[node.Description]
		}
	}
}

func executePrint(node *parser.TreeNode, env Environment) Variable {
	val := executeOperator(node.Children[0], env)
	fmt.Println(val.value)
	return val
}
