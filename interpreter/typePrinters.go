package interpreter

import (
	"fmt"
)

func nPrintString(value *Pointer) {
	log(cyan(stringValue(value)))
}

func nPrintBoolean(value *Pointer) {
	if booleanValue(value) {
		//todo: replace with variables for true and false string representations
		log(green("true"))
	} else {
		log(red("false"))
	}
}

func nPrintNumber(value *Pointer) {
	log(magenta(fmt.Sprint(numberValue(value))))
}

func nPrintChar(value *Pointer) {
	log(white(string(charValue(value))))

}

func mockPointer(address int, temp bool) *Pointer {
	p := Pointer{address, temp}
	return &p
}

func nativePrintArray(value *Pointer) {
	len := value.getDataLength() / type_sizes[POINTER]
	addr := value.address + PTR_DATA_OFFSET
	// fmt.Println(len, addr)
	log(bold("[ "))
	for i := 0; i < len; i++ {
		address := bytesToInt(heapSlice(addr, type_sizes[POINTER]))
		addr += type_sizes[POINTER]
		pointer := mockPointer(address, true)
		// fmt.Println(pointer)
		printVar(pointer)
		if i < len-1 {
			log(green(" , "))
		}
	}
	log(bold(" ]"))
}

func printVar(value *Pointer) {
	switch value.getDataType() {
	case NULL:
		log(red("<null>"))
	case POINTER:
		log(white(fmt.Sprintf("<pointer#%d>", value.address)))
	case ARRAY:
		nativePrintArray(value)
	case OBJECT:
		printObject(value)
	default:
		f, exists := printersMap[value.getDataType()]
		if !exists {

			log(white(fmt.Sprintf("<var@%d>", value.address)))
			return
		}
		f(value)

	}
}

func printObject(value *Pointer) {
	len := value.getDataLength() / type_sizes[POINTER]
	len = len / 2
	keyAddr := value.address + PTR_DATA_OFFSET
	log(bold("{\n"))
	pushIndent()
	for i := 0; i < len; i++ {
		keyHash := bytesToInt(heapSlice(keyAddr, type_sizes[POINTER]))
		keyAddr += type_sizes[POINTER]
		indentLog(bold(fmt.Sprintf("%d", keyHash)))
		log(green(" : "))
		valueAddr := bytesToInt(heapSlice(keyAddr, type_sizes[POINTER]))
		keyAddr += type_sizes[POINTER]
		valuePointer := mockPointer(valueAddr, true)
		printVar(valuePointer)
		log("\n")

	}
	log(bold("}"))
	popIndent()

}
