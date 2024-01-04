package interpreter

import (
	"fmt"
	"math/rand"
	"toylingo/parser"
	"toylingo/utils"
)

type funcDef struct {
	exec func(*scopeContext) Variable //doesnt need to return variable though
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
		args: []string{"prompt"},
	},
	"len": {
		exec: nativeLen,
		args: []string{"a"},
	},
	"makeArray": {
		exec: nativeMakeArray,
		args: []string{"size"},
	},
	"random": {
		exec: nativeRandom,
		args: []string{},
	},
}

func isNativeFunction(name string) bool {
	_, ok := nativeFunctions[name]
	return ok
}
//todo: call the actual implementation of function with specified arguments instead of having them retrieve from context
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
	value, exists := ctx.variables["a"]
	if !exists {
		interrupt("missing argument to print in function call print")
	}
	switch value.vartype {
	case TYPE_POINTER:
		interrupt("cannot print pointer")
	case TYPE_ARRAY:
		return nativePrintArray(ctx, value)
	case TYPE_STRING:
		fmt.Print(utils.Colors["WHITE"], string(heapSlice(value.pointer.address, value.pointer.size)), utils.Colors["RESET"])
		return value
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
	prompt,exists:=ctx.variables["prompt"]
	if exists && prompt.vartype==TYPE_STRING{
		fmt.Print(utils.Colors["WHITE"], string(heapSlice(prompt.pointer.address, prompt.pointer.size)), utils.Colors["RESET"])
	}
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

func nativeLen(ctx *scopeContext) Variable {
	value, exists := ctx.variables["a"]
	if !exists {
		interrupt("missing argument in function call len")
	}
	switch value.vartype {
	case TYPE_ARRAY:
		a := heapSlice(value.pointer.address, type_sizes[TYPE_NUMBER])
		size := byteArrayToFloat64(a)
		memaddr := malloc(type_sizes[TYPE_NUMBER], ctx.scopeId, true)
		writeBits(*memaddr, numberByteArray(size))
		ctx.returnValue = &Variable{memaddr, TYPE_NUMBER}
	case TYPE_STRING:
		val:=value.pointer.size
		memaddr := malloc(type_sizes[TYPE_NUMBER], ctx.scopeId, true)
		writeBits(*memaddr, numberByteArray(float64(val)))
		ctx.returnValue = &Variable{memaddr, TYPE_NUMBER}
	default:
		interrupt("function len expects array or string as argument")
	}
	return *ctx.returnValue
}

func nativeMakeArray(ctx *scopeContext) Variable {
	value, exists := ctx.variables["size"]
	if !exists {
		interrupt("missing argument in function call len")
	}
	if value.vartype != TYPE_NUMBER {
		interrupt("illegal value for array size")
	}
	return Variable{}
}

func nativeRandom(ctx *scopeContext) Variable {
	memaddr := malloc(type_sizes[TYPE_NUMBER], ctx.scopeId, true)
	writeBits(*memaddr, numberByteArray(rand.Float64()))
	ctx.returnValue = &Variable{memaddr, TYPE_NUMBER}
	return *ctx.returnValue
}
