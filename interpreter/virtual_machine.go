package interpreter

import (
	"fmt"
	"os"
)

var MEMSIZE = 1024 * 1024 //1mb
var HEAP = make([]byte, MEMSIZE)
var reserved = make([]bool, MEMSIZE)

var pointers = make(map[int]*Pointer, 0)

/*
a pointer returned by malloc will always have dataLength set to the requested length
and the data region will be zeroed out
type needs to be set by the caller
*/

// todo: this makes the whole language terribly slow
// automatically allocates extra memory for metadata
func malloc(datalen int, temp bool) *Pointer {
	size := datalen + PTR_DATA_OFFSET
	if size > MEMSIZE {
		interrupt(-1, "requested more memory than available", size, ">", MEMSIZE)
	}
	cap := 0
	for i := len(HEAP) - 1; i >= 0; i-- {
		if !reserved[i] {
			cap++
		} else {
			cap = 0
		}
		if cap == size {
			//reserve [i:i+size] and return pointer to i
			for j := i; j < i+size; j++ {
				reserved[j] = true
			}
			p := Pointer{i, temp}
			p.setDataLength(datalen)
			p.setReferenceCount(0)
			pointers[i] = &p
			debug_info("allocated", size, "bytes at", i)
			return &p
		}
	}
	//todo: try to defragment memory and try again
	// printMemoryStats()
	interrupt(-1, "out of memory: failed to allocate", size, "bytes")
	return nil
}

func freePtr(ptr *Pointer) {
	if ptr == NULL_POINTER {
		return
	}
	if ptr.getReferenceCount() > 0 {
		return
	}
	validatePointer(ptr)
	delete(pointers, ptr.address)
	debug_info("freeing", ptr.address)
	if ptr.getDataType() == ARRAY {
		freeArray(ptr)
	} else if ptr.getDataType() == OBJECT {
		freeObject(ptr)
	}
	end := ptr.address + ptr.getDataLength() + PTR_DATA_OFFSET
	for i := ptr.address; i < end; i++ {
		reserved[i] = false
		HEAP[i] = 0
	}
}

func freeArray(ptr *Pointer) {
	len := ptr.getDataLength() / type_sizes[POINTER]
	addr := ptr.address + PTR_DATA_OFFSET
	for i := 0; i < len; i++ {
		address := bytesToInt(heapSlice(addr, type_sizes[POINTER]))
		addr += type_sizes[POINTER]
		if address == 0 {
			continue
		}
		pointer := mockPointer(address, true)
		pointer.changeReferenceCount(false)
		freePtr(pointer)
	}
}

func freeObject(ptr *Pointer) {
	addr := ptr.address + PTR_DATA_OFFSET + type_sizes[POINTER]
	numkeys := ptr.getDataLength() / (type_sizes[POINTER] * 2)
	for i := 0; i < numkeys; i++ {
		valaddr := bytesToInt(heapSlice(addr, type_sizes[POINTER]))
		addr += type_sizes[POINTER] * 2
		if valaddr == 0 {
			continue
		}
		val := mockPointer(valaddr, true)
		val.changeReferenceCount(false)
		freePtr(val)
	}
}

func heapSlice(start int, size int) []byte {
	if start+size > MEMSIZE {
		interrupt(-1, "invalid heap slice access")
	}
	return HEAP[start : start+size]
}

// frees all temp pointers
func gc() {
	for _, ptr := range pointers {
		if ptr.temp {
			debug_info("gc: freeing temp pointer", ptr.address)
			freePtr(ptr)
		}
	}
}

func validatePointer(ptr *Pointer) {
	if !reserved[ptr.address] {
		interrupt(-1, "invalid pointer "+fmt.Sprint(ptr))
	}
}

func pointerAt(address int) *Pointer {
	return pointers[address]
}

func printMemoryStats() {
	if os.Getenv("MEMSTATS") == "0" {
		return
	}
	rvd := 0
	for _, v := range reserved {

		if v {
			rvd = rvd + 1
		}
	}

	fmt.Println("Occupied", rvd, "/", MEMSIZE, "bytes")
}
