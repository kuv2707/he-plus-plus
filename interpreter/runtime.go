package interpreter

import (
	"encoding/binary"
	"fmt"
	"math"
	"os"
	"toylingo/parser"
	"toylingo/utils"
)

var type_sizes = map[string]int{
	"number": 8,
	"char":   1,
	"bool":   1,
}
var LineNo = -1

// returns new variable with pointer to different address but same value is stored in both addresses
func copyVariable(variable Variable, scopeType string) Variable {
	addr := malloc(variable.pointer.size, scopeType, true)
	writeBits(*addr, int64(math.Float64bits(getNumber(variable))), 8)
	return Variable{addr, variable.vartype}
}

func getValue(variable Variable) interface{} {
	switch variable.vartype {
	case "number":
		return getNumber(variable)
	// case "char":
	// 	return getChar(variable)
	case "bool":
		return getBool(variable)
	}
	interrupt("invalid variable type " + variable.vartype)
	return nil
}

//todo:accept a byte array as value
func writeBits(ptr Pointer, value int64, size int) {
	for i := 0; i < size; i++ {
		HEAP[ptr.address+i] = byte(value & 0xFF)
		value = value >> 8
	}
}

func getNumber(variable Variable) float64 {
	ptr := variable.pointer
	validatePointer(*ptr)
	// Take 8 bytes from HEAP starting at ptr.address and convert to float64
	bytes := HEAP[ptr.address : ptr.address+8]
	parsedFloat := math.Float64frombits(binary.LittleEndian.Uint64(bytes))
	return parsedFloat
}

func getBool(variable Variable) bool {
	pointer := variable.pointer
	validatePointer(*pointer)
	// Take 1 byte from HEAP from end side of block pointed to by ptr
	parsedBool := HEAP[pointer.address+pointer.size-1] == 1
	return parsedBool
}

var contextStack = utils.MakeStack()

func pushScopeContext(label string) *scopeContext {
	ctx := scopeContext{label + "_" + fmt.Sprint(contextStack.Len()), make(map[string]Variable), make(map[string]parser.TreeNode), nil}
	if contextStack.IsEmpty() {
		contextStack.Push(ctx)
		return &ctx
	}
	for k, v := range contextStack.Peek().(scopeContext).variables {
		ctx.variables[k] = v
	}
	for k, v := range contextStack.Peek().(scopeContext).functions {
		ctx.functions[k] = v
	}
	contextStack.Push(ctx)
	return &ctx
}

func popScopeContext() {
	if contextStack.IsEmpty() {
		panic("no context to pop")
	}
	ctx := contextStack.Peek().(scopeContext)
	contextStack.Pop()
	for k, v := range ctx.variables {
		// debug_error("freeing?", k, v, "in", ctx.scopeType)
		if v.pointer.scopeId == ctx.scopeType {
			debug_info("freeing", k, v, "in", ctx.scopeType)
			freePtr(v.pointer)
		}
	}
	//free memory of inScopeVars

}

func printStackTrace() {
	s:=contextStack.GetStack()
	for i := range s {
		fmt.Println(s[len(s)-1-i].(scopeContext).scopeType)
	}
}

func interrupt(k ...interface{}) {
	fmt.Print(utils.Colors["RED"])
	fmt.Print("error at line", fmt.Sprint(LineNo), ": ")
	fmt.Println(k...)
	printStackTrace()
	fmt.Print(utils.Colors["RESET"])
	fmt.Print(utils.Colors["BOLDRED"])
	fmt.Println("execution interrupted")
	fmt.Print(utils.Colors["RESET"])
	os.Exit(1)
}
