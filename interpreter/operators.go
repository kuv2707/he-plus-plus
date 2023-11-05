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
	} else if node.Description == "!"{
		res:=executeNOT(node,env)
		return res
	} else if node.Description == "||"{
		return executeOR(node,env)
	} else if node.Description == "&&"{
		return executeAND(node,env)
	}
	panic("invalid operator")
	// return Variable{"", DATA_TYPES["NULL"], "NULL"}
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
		return Variable{"", DATA_TYPES["NUMBER"], left.value.(float32) + right.value.(float32)}
	} else if left.datatype == DATA_TYPES["STRING"] && right.datatype == DATA_TYPES["STRING"] {
		return Variable{"", DATA_TYPES["STRING"], left.value.(string) + right.value.(string)}
	} else {
		//todo
		return Variable{"", DATA_TYPES["STRING"], utils.StringVal(left.value)+utils.StringVal(right.value)}
	}
}
func executeSubtraction(node *parser.TreeNode, env Environment) Variable {
	left := executeOperator(node.Children[0], env)
	right := executeOperator(node.Children[1], env)
	if left.datatype == DATA_TYPES["NUMBER"] && right.datatype == DATA_TYPES["NUMBER"] {
		return Variable{"", DATA_TYPES["NUMBER"], left.value.(float32) - right.value.(float32)}
	} else {
		//todo
		return Variable{"", DATA_TYPES["NUMBER"], 0}
	}
}
func executeMultiplication(node *parser.TreeNode, env Environment) Variable {
	left := executeOperator(node.Children[0], env)
	right := executeOperator(node.Children[1], env)
	if left.datatype == DATA_TYPES["NUMBER"] && right.datatype == DATA_TYPES["NUMBER"] {
		return Variable{"", DATA_TYPES["NUMBER"], left.value.(float32) * right.value.(float32)}
	} else {
		//todo multiplying strings
		return Variable{"", DATA_TYPES["NUMBER"], 0}
	}
}
func executeDivision(node *parser.TreeNode, env Environment) Variable {
	left := executeOperator(node.Children[0], env)
	right := executeOperator(node.Children[1], env)
	if left.datatype == DATA_TYPES["NUMBER"] && right.datatype == DATA_TYPES["NUMBER"] {
		return Variable{"", DATA_TYPES["NUMBER"], left.value.(float32) / right.value.(float32)}
	} else {
		//todo
		return Variable{"", DATA_TYPES["NUMBER"], 0}
	}
}
func executeGreaterThan(node *parser.TreeNode, env Environment) Variable {
	left := executeOperator(node.Children[0], env)
	right := executeOperator(node.Children[1], env)
	if left.datatype == DATA_TYPES["NUMBER"] && right.datatype == DATA_TYPES["NUMBER"] {
		return Variable{"", DATA_TYPES["BOOL"], left.value.(float32) > right.value.(float32)}
	} else {
		//todo
		return Variable{"", DATA_TYPES["BOOL"], false}
	}
}
func executeLessThan(node *parser.TreeNode, env Environment) Variable {
	left := executeOperator(node.Children[0], env)
	right := executeOperator(node.Children[1], env)
	if left.datatype == DATA_TYPES["NUMBER"] && right.datatype == DATA_TYPES["NUMBER"] {
		return Variable{"", DATA_TYPES["BOOL"], left.value.(float32) < right.value.(float32)}
	} else {
		//todo
		return Variable{"", DATA_TYPES["BOOL"], false}
	}
}

func executeEquals(node *parser.TreeNode, env Environment) Variable {
	left := executeOperator(node.Children[0], env)
	right := executeOperator(node.Children[1], env)
	return Variable{"",DATA_TYPES["BOOL"],utils.StringVal(left.value)==utils.StringVal(right.value)}
}

func executePrimary(node *parser.TreeNode, env Environment) Variable {
	if node.Description == "true" || node.Description == "false" {
		return Variable{"", DATA_TYPES["BOOL"], node.Description == "true"}
	} else if utils.IsNumber(node.Description) {
		return Variable{"", DATA_TYPES["NUMBER"], utils.StringToNumber(node.Description)}
	} else {
		if !utils.ValidVariableName(node.Description){
			// evaluatedStr:=""
			for i:=1;i<len(node.Description)-1;i++{
				if node.Description[i]=='$' && node.Description[i+1]=='{'{

				}
			}
			return Variable{"",DATA_TYPES["STRING"],node.Description[1:len(node.Description)-1]}
		} else{

			return env.variables[node.Description]
		}
	}
}

func executePrint(node *parser.TreeNode, env Environment) Variable {
	val := executeOperator(node.Children[0], env)
	fmt.Print(utils.Colors["GREEN"],val.value,utils.Colors["RESET"])
	return val
}


func executeNOT(node *parser.TreeNode, env Environment) Variable {
	val:= executeOperator(node.Children[0],env)
	val.value=!val.value.(bool)
	return val
}

func executeOR(node *parser.TreeNode, env Environment) Variable {
	left:=executeOperator(node.Children[0],env)
	right:=executeOperator(node.Children[1],env)
	return Variable{"",DATA_TYPES["BOOL"],left.value.(bool)||right.value.(bool)}
}

func executeAND(node *parser.TreeNode, env Environment) Variable {
	left:=executeOperator(node.Children[0],env)
	right:=executeOperator(node.Children[1],env)
	return Variable{"",DATA_TYPES["BOOL"],left.value.(bool)&&right.value.(bool)}
}