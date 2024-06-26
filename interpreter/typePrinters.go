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
		log("can't print objects yet")
	default:
		f, exists := printersMap[value.getDataType()]
		if !exists {

			log(white(fmt.Sprintf("<var@%d>", value.address)))
			return
		}
		f(value)

	}
}
