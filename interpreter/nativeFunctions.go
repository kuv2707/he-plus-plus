package interpreter

import (
	"fmt"
	"he++/parser"
	"he++/utils"
	"math/rand"
)

type funcDef struct {
	exec func(*scopeContext)
	args []string
}

var nativeFunctions = map[string]funcDef{
	"print": {
		exec: nativePrint,
		args: []string{"arg"},
	},
	"println": {
		exec: nativePrintln,
		args: []string{"arg"},
	},
	"readNumber": {
		exec: nativeReadNumber,
		args: []string{"prompt"},
	},
	// "len": {
	// 	exec: nativeLen,
	// 	args: []string{"a"},
	// },
	// "makeArray": {
	// 	exec: nativeMakeArray,
	// 	args: []string{"size"},
	// },
	"random": {
		exec: nativeRandom,
		args: []string{},
	},
}

func isNativeFunction(name string) bool {
	_, ok := nativeFunctions[name]
	return ok
}

// todo: call the actual implementation of function with specified arguments instead of having them retrieve from context
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

func nativePrint(ctx *scopeContext) {
	value, exists := ctx.variables["arg"]
	if !exists {
		interrupt("missing argument to print in function call print")
	}
	switch value.getDataType() {
	case POINTER:
		interrupt("cannot print pointer")
	// case ARRAY:
	// 	nativePrintArray(ctx, value)
	case STRING:
		fmt.Print(utils.Colors["WHITE"], stringValue(value), utils.Colors["RESET"])
	

	case BOOLEAN:
		if booleanValue(value) {
			//todo: replace with variables for true and false string representations
			fmt.Print(utils.Colors["GREEN"], "true", utils.Colors["RESET"])
		} else {
			fmt.Print(utils.Colors["RED"], "false", utils.Colors["RESET"])
		}
	case NUMBER:
		printNumber(numberValue(value))
	case CHAR:
		fmt.Print(utils.Colors["WHITE"], charValue(value), utils.Colors["RESET"])
	default:
		fmt.Print(utils.Colors["WHITE"], fmt.Sprintf("<var@%#x>", value.address), utils.Colors["RESET"])

	}
}
func printNumber(val float64) {
	fmt.Print(utils.Colors["WHITE"], val, utils.Colors["RESET"])
}

func nativePrintln(ctx *scopeContext) {
	nativePrint(ctx)
	fmt.Print("\n")
}

func nativeReadNumber(ctx *scopeContext) {
	prompt, exists := ctx.variables["prompt"]
	if exists && prompt.getDataType() == STRING {
		fmt.Print(utils.Colors["WHITE"], stringValue(prompt), utils.Colors["RESET"])
	}
	var value float64
	fmt.Scan(&value)
	ptr := malloc(type_sizes[NUMBER], false)
	ptr.setDataType(NUMBER)
	writeDataContent(ptr, numberByteArray(value))
	ctx.returnValue = ptr
}

// func nativePrintArray(ctx *scopeContext, arrvar Variable) Variable {
// 	fmt.Print(utils.Colors["BOLDBLUE"] + "[ " + utils.Colors["RESET"])
// 	for i := 0; i < int(arrvar.pointer.size); i += type_sizes[TYPE_POINTER] {
// 		ptr := byteArrayToPointer(heapSlice(arrvar.pointer.address+i, type_sizes[TYPE_POINTER]))
// 		//todo: somehow retrieve the type of variable pointed to by the pointer ptr - currently assuming number
// 		num := byteArrayToFloat64(heapSlice(ptr, type_sizes[TYPE_NUMBER]))
// 		printNumber(num)
// 		if i+type_sizes[TYPE_POINTER] < int(arrvar.pointer.size) {
// 			fmt.Print(", ")
// 		}

// 	}
// 	fmt.Print(utils.Colors["BOLDBLUE"] + " ]" + utils.Colors["RESET"])
// 	return arrvar
// }

// func nativeLen(ctx *scopeContext) Variable {
// 	value, exists := ctx.variables["a"]
// 	if !exists {
// 		interrupt("missing argument in function call len")
// 	}
// 	val := 0
// 	switch value.vartype {
// 	case TYPE_ARRAY:
// 		val = value.pointer.size / type_sizes[TYPE_POINTER]
// 	case TYPE_STRING:
// 		val = value.pointer.size
// 	default:
// 		interrupt("function len expects array or string as argument")
// 	}
// 	memaddr := malloc(type_sizes[TYPE_NUMBER], ctx.scopeId, true)
// 	writeBytes(*memaddr, numberByteArray(float64(val)))
// 	ctx.returnValue = &Variable{memaddr, TYPE_NUMBER}
// 	return *ctx.returnValue
// }

// func nativeMakeArray(ctx *scopeContext) Variable {
// 	value, exists := ctx.variables["size"]
// 	if !exists {
// 		interrupt("missing argument in function call len")
// 	}
// 	if value.vartype != TYPE_NUMBER {
// 		interrupt("illegal value for array size")
// 	}
// 	size := int(getValue(value))
// 	if size < 0 {
// 		interrupt("illegal value for array size")
// 	}
// 	fmt.Println("making array of size", size)
// 	memaddr := malloc(size*type_sizes[TYPE_POINTER], ctx.scopeId, true)
// 	ctx.returnValue = &Variable{memaddr, TYPE_ARRAY}
// 	return *ctx.returnValue
// }

func nativeRandom(ctx *scopeContext) {
	memaddr := malloc(type_sizes[NUMBER], true)
	writeDataContent(memaddr, numberByteArray(rand.Float64()))
	ctx.returnValue = memaddr
}
