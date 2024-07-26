package interpreter

import (
	"fmt"
	u "he++/utils"
)

func nPrintString(value *Pointer) {
	u.Log(u.Cyan(stringValue(value)))
}

func nPrintBoolean(value *Pointer) {
	if booleanValue(value) {
		//todo: replace with variables for true and false string representations
		u.Log(u.Green("true"))
	} else {
		u.Log(u.Red("false"))
	}
}

func nPrintNumber(value *Pointer) {
	u.Log(u.Magenta(fmt.Sprint(numberValue(value))))
}

func nPrintChar(value *Pointer) {
	u.Log(u.White(string(charValue(value))))
}

func mockPointer(address int, temp bool) *Pointer {
	p := Pointer{address, temp}
	return &p
}

func nativePrintArray(value *Pointer) {
	len := value.getDataLength() / type_sizes[POINTER]
	addr := value.address + PTR_DATA_OFFSET
	// fmt.Println(len, addr)
	u.Log(u.Bold("[ "))
	for i := 0; i < len; i++ {
		address := bytesToInt(heapSlice(addr, type_sizes[POINTER]))
		addr += type_sizes[POINTER]
		pointer := mockPointer(address, true)
		// fmt.Println(pointer)
		printVar(pointer)
		if i < len-1 {
			u.Log(u.Green(" , "))
		}
	}
	u.Log(u.Bold(" ]"))
}

func printVar(value *Pointer) {
	switch value.getDataType() {
	case NULL:
		u.Log(u.Red("<null>"))
	case POINTER:
		u.Log(u.White(fmt.Sprintf("<pointer#%d>", value.address)))
	case ARRAY:
		nativePrintArray(value)
	case OBJECT:
		printObject(value)
	default:
		f, exists := printersMap[value.getDataType()]
		if !exists {
			u.Log(u.White(fmt.Sprintf("<var@%d>", value.address)))
			return
		}
		f(value)
	}
}

func printObject(value *Pointer) {
	len := value.getDataLength() / type_sizes[POINTER]
	len = len / 2
	keyAddr := value.address + PTR_DATA_OFFSET
	u.Log(u.Bold("{\n"))
	u.PushIndent()
	for i := 0; i < len; i++ {
		keyHash := bytesToInt(heapSlice(keyAddr, type_sizes[POINTER]))
		keyAddr += type_sizes[POINTER]
		u.IndentLog(u.Bold(fmt.Sprintf("%d", keyHash)))
		u.Log(u.Green(" : "))
		valueAddr := bytesToInt(heapSlice(keyAddr, type_sizes[POINTER]))
		keyAddr += type_sizes[POINTER]
		valuePointer := mockPointer(valueAddr, true)
		printVar(valuePointer)
		u.Log("\n")
	}
	u.Log(u.Bold("}"))
	u.PopIndent()
}
