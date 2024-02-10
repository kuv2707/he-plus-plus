package interpreter

import (
	"fmt"
)

/*
Pointer schema:
From the base address pointed to by the pointer:
* 1 byte: data type
* 4 bytes: data length  -  for primitive types, the data length is known beforehand, but for arrays and structs, it is not known and is hence included in the metadata
* n bytes: data
*/
type Pointer struct {
	address int
	temp    bool
}

// todo: a sophisticated struct storing function name, args, return type etc, and phase out the parser.TreeNode
type Function struct {
	name string
}

var PTR_DATA_OFFSET = 5

func (p *Pointer) getDataType() DataType {
	return DataType(HEAP[p.address])
}

func (p *Pointer) getDataLength() int {
	return bytesToInt(HEAP[p.address+1 : p.address+5])
}

func (p *Pointer) setDataLength(length int) {
	bts := intToBytes(length)
	for k := range bts {
		HEAP[p.address+1+k] = bts[k]
	}

}

func (p *Pointer) setDataType(dt DataType) {
	HEAP[p.address] = byte(dt)
}

func (p *Pointer) isNull() bool {
	return p == NULL_POINTER
}

// prints the region in hex form
func (p *Pointer) print() {
	fmt.Print(p.address, " ", p.temp, " ")
	fmt.Print(" ", p.getDataType())
	datalen := p.getDataLength()
	fmt.Print(" ", datalen, " ")
	for i := 0; i < datalen; i += 2 {
		fmt.Printf("%x ", HEAP[p.address+PTR_DATA_OFFSET+i:p.address+PTR_DATA_OFFSET+i+2])
	}
	fmt.Println()
}

func (p *Pointer) clone() *Pointer {
	if p == NULL_POINTER {
		return NULL_POINTER
	}
	// p.print()
	//allocate new memory
	newptr := malloc(p.getDataLength(), true)
	//copy metadata
	newptr.setDataType(p.getDataType())
	newptr.setDataLength(p.getDataLength())
	//copy data
	for i := 0; i < p.getDataLength(); i++ {
		HEAP[newptr.address+PTR_DATA_OFFSET+i] = HEAP[p.address+PTR_DATA_OFFSET+i]
	}
	//debug_info("cloned", p.address, "to", newptr.address, "datalen", p.getDataLength())
	return newptr
}
