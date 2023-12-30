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
	"readNumber": {
		exec: nativeReadNumber,
		args: []string{},
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
	value,exists := ctx.variables["a"]
	if !exists {
		interrupt("missing argument to print in function call print")
	}
	switch value.vartype {
	case TYPE_POINTER:
		interrupt("cannot print pointer")
	case TYPE_ARRAY:
		return nativePrintArray(ctx, value)
	}

	val := getValue(value)
	switch value.vartype {
	case TYPE_BOOLEAN:
		if getBool(value) {
			fmt.Print(utils.Colors["GREEN"], "true", utils.Colors["RESET"])
		} else {
			fmt.Print(utils.Colors["RED"], "false", utils.Colors["RESET"])
		}
	case TYPE_NUMBER:
		printNumber(val)
	case TYPE_CHAR:
		fmt.Print(utils.Colors["WHITE"], string(int(val)), utils.Colors["RESET"])
	default:
		fmt.Print(utils.Colors["WHITE"], val, utils.Colors["RESET"])

	}
	return value

}
func printNumber(val float64) {
	fmt.Print(utils.Colors["WHITE"], val, utils.Colors["RESET"])
}

func nativePrintln(ctx *scopeContext) Variable {
	v := nativePrint(ctx)
	fmt.Print("\n")
	return v
}

func nativeReadNumber(ctx *scopeContext) Variable {
	var value float64
	fmt.Scan(&value)
	memaddr := malloc(type_sizes[TYPE_NUMBER], ctx.scopeId, false)
	writeBits(*memaddr, numberByteArray(value))
	v := Variable{memaddr, TYPE_NUMBER}
	ctx.returnValue = &v
	return v
}

func nativePrintArray(ctx *scopeContext, arrvar Variable) Variable {
	a := heapSlice(arrvar.pointer.address, type_sizes[TYPE_NUMBER])
	size := byteArrayToFloat64(a)
	addr := arrvar.pointer.address
	fmt.Print(utils.Colors["BOLDBLUE"] + "[ " + utils.Colors["RESET"])
	addr += type_sizes[TYPE_NUMBER]
	for i := 1; i <= int(size); i++ {
		ptr := byteArrayToPointer(heapSlice(addr, type_sizes[TYPE_POINTER]))
		//todo: somehow retrieve the type of variable pointed to by the pointer ptr - currently assuming number
		num := byteArrayToFloat64(heapSlice(ptr, type_sizes[TYPE_NUMBER]))
		printNumber(num)
		if i != int(size) {
			fmt.Print(", ")
		}
		addr += type_sizes[TYPE_POINTER]
	}
	fmt.Print(utils.Colors["BOLDBLUE"] + " ]" + utils.Colors["RESET"])
	return arrvar
}
