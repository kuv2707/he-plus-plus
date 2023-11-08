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
		res := executeEquals(node, env)
		res.value = !res.value.(bool)
		return res
	} else if node.Description == "<=" {
	} else if node.Description == ">=" {
	} else if node.Description == "!" {
		res := executeNOT(node, env)
		return res
	} else if node.Description == "||" {
		return executeOR(node, env)
	} else if node.Description == "&&" {
		return executeAND(node, env)
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
		return Variable{"", DATA_TYPES["STRING"], utils.StringVal(left.value) + utils.StringVal(right.value)}
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
	} else if left.datatype == DATA_TYPES["STRING"] && right.datatype == DATA_TYPES["STRING"] {
		panic("invalid operation")
		return Variable{"", DATA_TYPES["STRING"], left.value.(string) + right.value.(string)}
	} else {
		//todo
		if left.datatype == DATA_TYPES["STRING"] && right.datatype == DATA_TYPES["NUMBER"] {
			res := ""
			for i := 0; i < int(right.value.(float32)); i++ {
				res += left.value.(string)
			}
			return Variable{"", DATA_TYPES["STRING"], res}
		} else {
			res := ""
			for i := 0; i < int(left.value.(float32)); i++ {
				res += right.value.(string)
			}
			return Variable{"", DATA_TYPES["STRING"], res}
		}
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
	return Variable{"", DATA_TYPES["BOOL"], utils.StringVal(left.value) == utils.StringVal(right.value)}
}

func executePrimary(node *parser.TreeNode, env Environment) Variable {
	fmt.Println("hereee")
	if node.Description == "true" || node.Description == "false" {
		return Variable{"", DATA_TYPES["BOOL"], node.Description == "true"}
	} else if utils.IsNumber(node.Description) {
		return Variable{"", DATA_TYPES["NUMBER"], utils.StringToNumber(node.Description)}
	} else {
		
		if !utils.ValidVariableName(node.Description) {
			// evaluatedStr:=""
			for i := 1; i < len(node.Description)-1; i++ {
				if node.Description[i] == '$' && node.Description[i+1] == '{' {

				}
			}
			return Variable{"", DATA_TYPES["STRING"], node.Description[1 : len(node.Description)-1]}
		} else {
			if len(node.Children) > 0 {
				//function call
				fmt.Println("function call")
				params:=make([]Variable,0)
				for key, childNode := range node.Children[0].Properties {
					if key[0:4] == "args" {
						params=append(params,executeOperator(childNode,env))
					}
				}
				funcNode:=env.functions[node.Description]
				newenv:=NewCallStackContext(env)
				for i, param := range params {
					newenv.variables[funcNode.Properties["args"+fmt.Sprint(i)].Description]=param
				}
				fmt.Println("calling function",node.Description)
				k:=ExecuteAST(funcNode,&newenv)
				// fmt.Println(k)
				return k.returnVal

			}
			if _, ok := env.variables[node.Description]; !ok {
				panic("variable not defined: " + node.Description)
			}
			return env.variables[node.Description]
		}
	}
}

func executePrint(node *parser.TreeNode, env Environment) Variable {
	val := executeOperator(node.Children[0], env)
	fmt.Print(utils.Colors["GREEN"], val.value, utils.Colors["RESET"])
	return val
}

func executeNOT(node *parser.TreeNode, env Environment) Variable {
	val := executeOperator(node.Children[0], env)
	val.value = !val.value.(bool)
	return val
}

func executeOR(node *parser.TreeNode, env Environment) Variable {
	left := executeOperator(node.Children[0], env)
	right := executeOperator(node.Children[1], env)
	return Variable{"", DATA_TYPES["BOOL"], left.value.(bool) || right.value.(bool)}
}

func executeAND(node *parser.TreeNode, env Environment) Variable {
	left := executeOperator(node.Children[0], env)
	right := executeOperator(node.Children[1], env)
	return Variable{"", DATA_TYPES["BOOL"], left.value.(bool) && right.value.(bool)}
}
