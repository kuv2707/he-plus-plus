package interpreter

import "fmt"

// import "fmt"

var MEMSIZE = 1024 //1kb
var HEAP = make([]byte, MEMSIZE)
var reserved = make([]bool, MEMSIZE)

var pointers = make(map[int]*Pointer, 0)

func malloc(datalen int, scid string, temp bool) *Pointer {
	size := datalen + PTR_DATA_OFFSET
	if size > MEMSIZE {
		interrupt("requested more memory than available", size, ">", MEMSIZE)
	}
	debug_info("requested to malloc", size, "bytes for", scid)
	cap := 0
	for i := len(HEAP) - 1; i >= 0; i-- {
		if HEAP[i] == 0 && !reserved[i] {
			cap++
		} else {
			cap = 0
		}
		if cap == size {
			//reserve [i:i+size] and return pointer to i
			for j := i; j < i+size; j++ {
				reserved[j] = true
			}
			p := Pointer{i, scid, temp}
			p.setDataLength(datalen)
			pointers[i] = &p
			// debug_info("allocated", size, "data bytes at", i, "for", scid)
			return &p
		}
	}
	//todo: try to defragment memory and try again

	interrupt("out of memory")
	return nil
}

func freePtr(ptr *Pointer) {
	validatePointer(ptr)
	delete(pointers, ptr.address)
	end := ptr.address + ptr.getDataLength() + PTR_DATA_OFFSET
	for i := ptr.address; i < end; i++ {
		reserved[i] = false
	}
}

func heapSlice(start int, size int) []byte {
	if start+size > MEMSIZE {
		interrupt("invalid heap slice access")
	}
	return HEAP[start : start+size]
}

func freeAll() {
	debug_error("freeing all pointers")
	for _, ptr := range pointers {
		freePtr(ptr)
	}
}

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
		interrupt("invalid pointer " + fmt.Sprint(ptr))
	}
}

func printMemoryStats() {
	rvd := 0
	for _, v := range reserved {
		if v {
			rvd = rvd + 1
		}
	}
	debug_info("Occupied", rvd, "/", MEMSIZE, "bytes")
}
