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

var NULL_POINTER = &Pointer{-1, false}

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
	// copy data
	for i := 0; i < src.getDataLength(); i++ {
		HEAP[dest.address+PTR_DATA_OFFSET+i] = HEAP[src.address+PTR_DATA_OFFSET+i]
	}
}

func unsafeWriteBytes(ptr int, value []byte) {
	for i,v := range value {
		HEAP[ptr+i] = v
	}
}

func byteArrayToFloat64(bytes []byte) float64 {
	return math.Float64frombits(binary.LittleEndian.Uint64(bytes))
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
	ctx := scopeContext{"", scopetype, scopename, make(map[string]*Pointer), make(map[string]parser.TreeNode), NULL_POINTER}
	contextStack.Push(ctx)
	return &ctx
}

// should be used to get a variable instead of raw ctx.variables[name]
// nativefunctions can use ctx.variables[name] directly coz they are guaranteed to exist in the same scope
func findVariable(name string) *Pointer {
	for i := contextStack.Len() - 1; i >= 0; i-- {
		ctx := contextStack.Get(i).(scopeContext)
		ptr, exists := ctx.variables[name]
		if exists {
			return ptr
		}
	}
	return NULL_POINTER
}

func findFunction(name string) *parser.TreeNode {
	for i := contextStack.Len() - 1; i >= 0; i-- {
		ctx := contextStack.Get(i).(scopeContext)
		fn, exists := ctx.functions[name]
		if exists {
			return &fn
		}
	}
	return nil
}

func popScopeContext() {
	if contextStack.IsEmpty() {
		panic("no context to pop")
	}
	ctx := contextStack.Peek().(scopeContext)
	contextStack.Pop()
	for _, v := range ctx.variables {
		freePtr(v)
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
