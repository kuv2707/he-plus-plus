package interpreter

import (
	"fmt"
	"toylingo/parser"
	"toylingo/utils"
)

type funcDef struct {
	exec func(*scopeContext) Variable
	args []string
}

var nativeFunctions = map[string]funcDef{
	"print": {
		exec: nativePrint,
		args: []string{"a"},
	},
	"println": {
		exec: nativePrintln,
		args: []string{"a"},
	},
}

func isNativeFunction(name string) bool {
	_, ok := nativeFunctions[name]
	return ok
}

func addNativeFuncDeclarations(ctx *scopeContext) {
	for k, v := range nativeFunctions {
		name := k
		args := v.args
		fnode := parser.TreeNode{
			Label:       "function",
			Description: name,
			Children:    nil,
			Properties:  nil,
		}
		argnode := parser.TreeNode{
			Label:       "args",
			Description: "args",
			Children:    nil,
			Properties:  nil,
		}
		for _, arg := range args {
			argnode.Children = append(argnode.Children, &parser.TreeNode{
				Label:       "arg",
				Description: arg,
				Children:    nil,
				Properties:  nil,
			})
		}
		fnode.Properties = map[string]*parser.TreeNode{
			"args": &argnode,
		}
		ctx.functions[name] = fnode
	}
}

func nativePrint(ctx *scopeContext) Variable {
	value := ctx.variables["a"]
	val := getValue(value)
	if value.vartype == TYPE_NUMBER {
		fmt.Print(utils.Colors["CYAN"], val, utils.Colors["RESET"])
	} else if value.vartype == "bool" {
		if getBool(value) {
			fmt.Print(utils.Colors["GREEN"], "true", utils.Colors["RESET"])
		} else {
			fmt.Print(utils.Colors["RED"], "false", utils.Colors["RESET"])
		}
	} else if value.vartype == "char" {
		fmt.Print(utils.Colors["WHITE"], string(int(val)), utils.Colors["RESET"])
	} else {
		fmt.Print(utils.Colors["WHITE"], val, utils.Colors["RESET"])
	}
	return value

}
func nativePrintln(ctx *scopeContext) Variable {
	v:=nativePrint(ctx)
	fmt.Print("\n")
	return v
}