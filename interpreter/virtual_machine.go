package interpreter

var MEMSIZE = 1024  //1kb
var HEAP = make([]byte, MEMSIZE)
var reserved = make([]bool,MEMSIZE)

func malloc(size int) Pointer {
	cap:=0
	for i := len(HEAP)-1; i >=0; i++ {
		if HEAP[i]==0 {
			cap++
		}else{
			cap=0
		}
		if cap==size {
			//reserve [i:i+size] and return pointer to i
			for j:=i;j<i+size;j++ {
				reserved[j]=true
			}
			return Pointer{i,size}
		}
	}
	//todo: try to defragment memory and try again

	panic("out of memory")
}

func free(ptr Pointer) {
	for i:=ptr.address;i<ptr.address+ptr.size;i++ {
		reserved[i]=false
		HEAP[i]=0
	}
}
