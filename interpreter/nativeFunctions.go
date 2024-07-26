package interpreter

import (
	"fmt"
	"he++/parser"
	"he++/utils"
	"math/rand"
)

type funcDef struct {
	exec func(*ScopeContext)
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
	"len": {
		exec: nativeLen,
		args: []string{"array"},
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

//todo: line no for error message 

func nativeMakeArray(ctx *ScopeContext) {
	value, exists := ctx.variables["size"]
	if !exists {
		interrupt(-1,"Size of array not passed")
	}
	if value.getDataType() != NUMBER {
		interrupt(-1,"invalid argument, expected NUMBER, found", value.getDataType())
	}
	len := int(numberValue(value))
	if len <= 0 {
		interrupt(-1,"length of array must be greater than 0")
	}
	arrptr := malloc(type_sizes[POINTER]*len, false)
	arrptr.setDataType(ARRAY)
	arrptr.setDataLength(len * type_sizes[POINTER])
	ctx.returnValue = arrptr
}

func nativeLen(ctx *ScopeContext) {
	value, exists := ctx.variables["array"]
	if !exists {
		interrupt(-1,"No array is passed to find the length of")
	}
	div := 0
	switch value.getDataType() {
	case ARRAY:
		div = type_sizes[POINTER]
	case STRING:
		div = type_sizes[CHAR]
	default:
		interrupt(-1,"Can only find length of arrays and strings")
	}
	len := value.getDataLength() / div
	memaddr := malloc(type_sizes[NUMBER], false)
	memaddr.setDataType(NUMBER)
	writeDataContent(memaddr, numberByteArray(float64(len)))
	ctx.returnValue = memaddr
}

func isNativeFunction(name string) bool {
	_, ok := nativeFunctions[name]
	return ok
}

// todo: call the actual implementation of function with specified arguments instead of having them retrieve from context
func addNativeFuncDeclarations(ctx *ScopeContext) {
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

var printersMap = map[DataType]func(*Pointer){
	STRING:  nPrintString,
	NUMBER:  nPrintNumber,
	BOOLEAN: nPrintBoolean,
	CHAR:    nPrintChar,
}

func nativePrint(ctx *ScopeContext) {
	value, exists := ctx.variables["arg"]
	if !exists {
		interrupt(-1,"missing argument to print in function call print")
	}
	printVar(value)
}

func nativePrintln(ctx *ScopeContext) {
	nativePrint(ctx)
	fmt.Print("\n")
}

func nativeReadNumber(ctx *ScopeContext) {
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

func nativeRandom(ctx *ScopeContext) {
	memaddr := malloc(type_sizes[NUMBER], false)
	memaddr.setDataType(NUMBER)
	writeDataContent(memaddr, numberByteArray(rand.Float64()))
	ctx.returnValue = memaddr
}
