package interpreter

import (
	"fmt"
	"os"
)

// import "fmt"

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
			pointers[i] = &p
			//debug_info("allocated", size, "data bytes at", i)
			return &p
		}
	}
	//todo: try to defragment memory and try again
	printMemoryStats()
	interrupt(-1, "out of memory: failed to allocate", size, "bytes")
	return nil
}

func freePtr(ptr *Pointer) {
	validatePointer(ptr)
	delete(pointers, ptr.address)
	// fmt.Print("freeing ")
	// ptr.print()
	//debug_info("freeing", ptr.address)
	end := ptr.address + ptr.getDataLength() + PTR_DATA_OFFSET
	cnt := 0
	for i := ptr.address; i < end; i++ {
		reserved[i] = false
		HEAP[i] = 0
		cnt++
	}
}

func heapSlice(start int, size int) []byte {
	if start+size > MEMSIZE {
		interrupt(-1, "invalid heap slice access")
	}
	return HEAP[start : start+size]
}

// should not be needed
// func freeAll() {
// 	debug_error("freeing all pointers")
// 	for _, ptr := range pointers {
// 		freePtr(ptr)
// 	}
// }

// frees all temp pointers
func gc() {
	for _, ptr := range pointers {
		if ptr.temp {
			freePtr(ptr)
		}
	}
}

func validatePointer(ptr *Pointer) {
	if !reserved[ptr.address] {
		interrupt(-1, "invalid pointer "+fmt.Sprint(ptr))
	}
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

	//debug_info("Occupied", rvd, "/", MEMSIZE, "bytes")
}
