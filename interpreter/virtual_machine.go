package interpreter

import "fmt"

// import "fmt"

var MEMSIZE = 1024 //1kb
var HEAP = make([]byte, MEMSIZE)
var reserved = make([]bool, MEMSIZE)

var pointers = make(map[Pointer]interface{}, 0)

func malloc(size int) Pointer {

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
			p := Pointer{i, size}
			pointers[p] = true
			return p
		}
	}
	//todo: try to defragment memory and try again

	panic("out of memory")
}

func freePtr(ptr Pointer) {
	delete(pointers, ptr)
	fmt.Printf("freeing %d to %d\n", ptr.address, ptr.address+ptr.size)
	for i := ptr.address; i < ptr.address+ptr.size; i++ {
		reserved[i] = false
		HEAP[i] = 0
	}
}

func freeAll() {
	for ptr := range pointers {
		freePtr(ptr)
	}
}
