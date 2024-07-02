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

var NULL_POINTER = &Pointer{0, false}

func writeDataContent(ptr *Pointer, value []byte) {
	validatePointer(ptr)
	datalen := ptr.getDataLength()
	//the implication of the following check is that the data length of a pointer cannot be changed
	//meaning that STRING and ARRAY types cannot be resized
	if len(value) != datalen {
		interrupt(-1, "invalid data length", len(value), "expected", datalen)
	}
	for i := range value {
		HEAP[ptr.address+PTR_DATA_OFFSET+i] = value[i]
	}
}

// creates a shallow copy of the data pointed to by src into dest
func writeContentFromOnePointerToAnother(dest *Pointer, src *Pointer) {
	validatePointer(dest)
	validatePointer(src)
	//todo: merge both
	// copy metadata -- should we?
	// HEAP[dest.address] = HEAP[src.address]
	// for i := 1; i < PTR_DATA_OFFSET; i++ {
	// 	HEAP[dest.address+i] = HEAP[src.address+i]
	// }
	// copy data
	for i := 0; i < src.getDataLength(); i++ {
		HEAP[dest.address+PTR_DATA_OFFSET+i] = HEAP[src.address+PTR_DATA_OFFSET+i]
	}
}

func unsafeWriteBytes(ptr int, value []byte) {
	for i, v := range value {
		HEAP[ptr+i] = v
	}
}

func byteArrayToFloat64(bytes []byte) float64 {
	return math.Float64frombits(binary.LittleEndian.Uint64(bytes))
}

func bytesToInt(bytes []byte) int {
	return int(binary.LittleEndian.Uint32(bytes))
}

func bytesToInt16(bytes []byte) int {
	return int(binary.LittleEndian.Uint16(bytes))
}

func intToBytes(value int) []byte {
	bytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(bytes, uint32(value))
	return bytes
}

func int16ToBytes(value int) []byte {
	bytes := make([]byte, 2)
	binary.LittleEndian.PutUint16(bytes, uint16(value))
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

func makeScopeContext(scopetype string, scopename string) ScopeContext {
	ctx := ScopeContext{scopename, scopetype, 0, make(map[string]*Pointer), make(map[string]parser.TreeNode), NULL_POINTER}
	return ctx
}

func pushToContextStack(ctx ScopeContext) {
	contextStack.Push(ctx)
}

func pushScopeContext(scopetype string, scopename string) *ScopeContext {
	k := makeScopeContext(scopetype, scopename)
	contextStack.Push(k)
	return &k
}

// should be used to get a variable instead of raw ctx.variables[name]
// nativefunctions can use ctx.variables[name] directly coz they are guaranteed to exist in the same scope
// returns the pointer and the scope context in which the variable was found
func findVariable(name string) (*Pointer, *ScopeContext) {
	for i := contextStack.Len() - 1; i >= 0; i-- {
		ctx := contextStack.Get(i).(ScopeContext)
		ptr, exists := ctx.variables[name]
		if exists {
			return ptr, &ctx
		}
	}
	return NULL_POINTER, nil
}

func findFunction(name string) *parser.TreeNode {
	for i := contextStack.Len() - 1; i >= 0; i-- {
		ctx := contextStack.Get(i).(ScopeContext)
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

	ctx := contextStack.Peek().(ScopeContext)
	if ctx.scopeId == "root" && os.Getenv("REPL") == "1" {
		return
	}
	for _, v := range ctx.variables {
		v.changeReferenceCount(false)
		freePtr(v)
	}
	gc()
	contextStack.Pop()
}

func getScopeContext(depth int) ScopeContext {
	return contextStack.Get(contextStack.Len() - 1 - depth).(ScopeContext)

}

func printStackTrace() {
	s := contextStack.GetStack()
	for i := range s {
		ctx := s[len(s)-1-i].(ScopeContext)
		fmt.Println(ctx.scopeId, "line", ctx.currentLine)
	}
}

func interrupt(lineNo int, k ...interface{}) {
	if lineNo < 0 {
		// lineNo = getScopeContext(0).currentLine
	}
	fmt.Print(utils.Colors["RED"])
	fmt.Print("error at line ", fmt.Sprint(lineNo), ": ")
	fmt.Println(k...)
	printStackTrace()
	fmt.Print(utils.Colors["RESET"])
	fmt.Print(utils.Colors["BOLDRED"])
	fmt.Println("execution interrupted")
	fmt.Print(utils.Colors["RESET"])
	os.Exit(1)
}
