package interpreter

import (
	"encoding/binary"
	"fmt"
	"he++/parser"
	"he++/utils"
	"math"
	"os"
)

var type_sizes = map[DataType]int{
	NUMBER:  8,
	CHAR:    1,
	BOOLEAN: 1,
	POINTER: 4,
}

var NULL_POINTER = &Pointer{-1, "", false}

func writeDataContent(ptr *Pointer, value []byte) {
	validatePointer(ptr)
	datalen := ptr.getDataLength()
	//the implication of the following check is that the data length of a pointer cannot be changed
	//meaning that STRING and ARRAY types cannot be resized
	if len(value) != datalen {
		interrupt("invalid data length", len(value), "expected", datalen)
	}
	for i := range value {
		HEAP[ptr.address+5+i] = value[i]
	}
}

// creates a shallow copy of the data pointed to by src into dest
func writeContentFromOnePointerToAnother(dest *Pointer, src *Pointer) {
	validatePointer(dest)
	validatePointer(src)
	// copy metadata
	HEAP[dest.address] = HEAP[src.address]
	for i := 1; i < 5; i++ {
		HEAP[dest.address+i] = HEAP[src.address+i]
	}

	for i := 0; i < src.getDataLength(); i++ {
		HEAP[dest.address+5+i] = HEAP[src.address+5+i]
	}
}

func unsafeWriteBytes(ptr int, value []byte) {
	for i := range value {
		HEAP[ptr+i] = value[i]
	}
}

func byteArrayToFloat64(bytes []byte) float64 {
	return math.Float64frombits(binary.LittleEndian.Uint64(bytes))
}

func byteArrayToPointer(bytes []byte) int {
	return int(binary.LittleEndian.Uint32(bytes))
}

func bytesToInt(bytes []byte) int {
	return int(binary.LittleEndian.Uint32(bytes))
}

func intToBytes(value int) []byte {
	bytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(bytes, uint32(value))
	return bytes
}

func numberByteArray(value float64) []byte {
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, math.Float64bits(value))
	return bytes
}

func stringAsBytes(value string) []byte {
	bytes := []byte(value)
	return bytes
}

func bytesAsString(value []byte) string {
	return string(value)
}

func pointerAsBytes(value int) []byte {
	bytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(bytes, uint32(value))
	return bytes
}

func booleanValue(p *Pointer) bool {
	return HEAP[p.address+PTR_DATA_OFFSET] == 1
}

func numberValue(p *Pointer) float64 {
	return byteArrayToFloat64(HEAP[p.address+PTR_DATA_OFFSET : p.address+PTR_DATA_OFFSET+8])
}

func charValue(p *Pointer) rune {
	return rune(HEAP[p.address+PTR_DATA_OFFSET])
}

func stringValue(p *Pointer) string {
	return bytesAsString(HEAP[p.address+PTR_DATA_OFFSET : p.address+PTR_DATA_OFFSET+p.getDataLength()])
}

var contextStack = utils.MakeStack()

func pushScopeContext(scopetype string, scopename string) *scopeContext {
	ctx := scopeContext{generateId(), scopetype, scopename, make(map[string]*Pointer), make(map[string]parser.TreeNode), NULL_POINTER}
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
		if v.scopeId == ctx.scopeId {
			debug_info("freeing?", k, v, "in", ctx.scopeName)
			debug_info("freeing", k, "in", ctx.scopeName)
			freePtr(v)
		}
	}
	gc()
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
