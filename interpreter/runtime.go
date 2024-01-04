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
	TYPE_NUMBER:  8,
	TYPE_CHAR:    1,
	TYPE_BOOLEAN: 1,
	TYPE_POINTER: 4,
}
var LineNo = -1

// returns new variable with pointer to different address but same value is stored in both addresses
func copyVariable(variable Variable, sid string) Variable {
	addr := malloc(variable.pointer.size, sid, true)
	writeBits(*addr, heapSlice(variable.pointer.address, variable.pointer.size))
	return Variable{addr, variable.vartype}
}

// returns the number equivalent of the variable
func getValue(variable Variable) float64 {
	switch variable.vartype {
	case TYPE_NUMBER:
		return getNumber(variable)
	// case "char":
	// 	return getChar(variable)
	case TYPE_BOOLEAN:
		b := getBool(variable)
		if b {
			return 1
		}
		return 0
	//DOUBT: shouldnt expose pointer like this right? just return 0
	case TYPE_ARRAY:
		fallthrough
	case TYPE_STRING:
		return float64(variable.pointer.address)
	}

	interrupt("invalid variable type " + variable.vartype)
	return 0
}

//todo:accept a byte array as value
func writeBits(ptr Pointer, value []byte) {
	validatePointer(ptr)
	for i := range value {
		HEAP[ptr.address+i] = value[i]
	}
}

func unsafeWriteBits(ptr int, value []byte) {
	for i := range value {
		HEAP[ptr+i] = value[i]
	}
}

func getNumber(variable Variable) float64 {
	if variable.vartype != TYPE_NUMBER {
		interrupt("invalid number type " + variable.vartype)
	}
	ptr := variable.pointer
	validatePointer(*ptr)
	// Take 8 bytes from HEAP starting at ptr.address and convert to float64
	bytes := HEAP[ptr.address : ptr.address+8]
	parsedFloat := byteArrayToFloat64(bytes)
	return parsedFloat
}

func byteArrayToFloat64(bytes []byte) float64 {
	return math.Float64frombits(binary.LittleEndian.Uint64(bytes))
}

func byteArrayToPointer(bytes []byte) int {
	return int(binary.LittleEndian.Uint32(bytes))
}

func numberByteArray(value float64) []byte {
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, math.Float64bits(value))
	return bytes
}

func stringByteArray(value string) []byte {
	bytes := []byte(value)
	return bytes
}

func byteArrayString(value []byte) string {
	return string(value)
}

func pointerByteArray(value int) []byte {
	bytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(bytes, uint32(value))
	return bytes
}

func getBool(variable Variable) bool {
	pointer := variable.pointer
	validatePointer(*pointer)
	// Take 1 byte from HEAP from end side of block pointed to by ptr
	parsedBool := HEAP[pointer.address+pointer.size-1] == 1
	return parsedBool
}

var contextStack = utils.MakeStack()

func pushScopeContext(scopetype string, scopename string) *scopeContext {
	ctx := scopeContext{generateId(), scopetype, scopename, make(map[string]Variable), make(map[string]parser.TreeNode), nil}
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
		if v.pointer.scopeId == ctx.scopeId {
			debug_info("freeing", k, v.pointer, v.vartype, "in", ctx.scopeName)
			if v.vartype == TYPE_ARRAY {
				freeArrPtr(v.pointer)
			} else {
				freePtr(v.pointer)
			}
		}
	}
	//free memory of inScopeVars

}

func freeArrPtr(ptr *Pointer) {
	for i := type_sizes[TYPE_NUMBER]; i < ptr.size; i += type_sizes[TYPE_POINTER] {
		p:=byteArrayToPointer(heapSlice(ptr.address+i, type_sizes[TYPE_POINTER]))
		freePtr(pointers[p])
	}
	freePtr(ptr)
}

func getScopeContext(depth int) scopeContext {
	return contextStack.Get(contextStack.Len() - 1 - depth).(scopeContext)

}

func printStackTrace() {
	s := contextStack.GetStack()
	for i := range s {
		fmt.Println(s[len(s)-1-i].(scopeContext).scopeName)
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
