package interpreter

import (
	"fmt"
	"he++/utils"
)

func nPrintString(value *Pointer) {
	fmt.Print(utils.Colors["CYAN"], stringValue(value), utils.Colors["RESET"])
}

func nPrintBoolean(value *Pointer) {
	if booleanValue(value) {
		//todo: replace with variables for true and false string representations
		fmt.Print(utils.Colors["GREEN"], "true", utils.Colors["RESET"])
	} else {
		fmt.Print(utils.Colors["RED"], "false", utils.Colors["RESET"])
	}
}

func nPrintNumber(value *Pointer) {
	fmt.Print(utils.Colors["MAGENTA"], numberValue(value), utils.Colors["RESET"])
}

func nPrintChar(value *Pointer) {
	fmt.Print(utils.Colors["WHITE"], charValue(value), utils.Colors["RESET"])

}

func nativePrintArray(value *Pointer) {
	len := value.getDataLength() / type_sizes[POINTER]
	addr := value.address + PTR_DATA_OFFSET
	// fmt.Println(len, addr)
	fmt.Print(utils.Colors["GREEN"], "[ ")
	for i := 0; i < len; i++ {
		address := bytesToInt(heapSlice(addr, type_sizes[POINTER]))
		addr += type_sizes[POINTER]
		pointer := pointers[address]
		// fmt.Println(pointer)
		printVar(pointer)
		if i < len-1 {
			fmt.Print(utils.Colors["GREEN"], " , ")
		}
	}
	fmt.Print(utils.Colors["GREEN"], " ]", utils.Colors["RESET"])
}

func printVar(value *Pointer) {
	switch value.getDataType() {
	case NULL:
		fmt.Print(utils.Colors["RED"], "<null>", utils.Colors["RESET"])
	case POINTER:
		fmt.Print(utils.Colors["WHITE"], fmt.Sprintf("<pointer#%d>", value.address), utils.Colors["RESET"])
	case ARRAY:
		nativePrintArray(value)
	default:
		f, exists := printersMap[value.getDataType()]
		if !exists {

			fmt.Print(utils.Colors["WHITE"], fmt.Sprintf("<var@%d>", value.address), utils.Colors["RESET"])
			return
		}
		f(value)

	}
}
