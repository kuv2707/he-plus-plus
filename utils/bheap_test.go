// heap_test.go
package utils

import "testing"

func TestHeapPushMinHeap(t *testing.T) {
	h := MakeHeap(func(a, b int) bool { return a < b })

	input := []int{5, 3, 8, 1, 4}
	for _, v := range input {
		h.Push(v)
	}

	if h.Peek() != 1 {
		t.Fatalf("expected min element 1, got %d", h.Peek())
	}
}

func TestHeapPopMinHeap(t *testing.T) {
	h := MakeHeap(func(a, b int) bool { return a < b })

	input := []int{7, 2, 6, 1, 5}
	for _, v := range input {
		h.Push(v)
	}

	expected := []int{1, 2, 5, 6, 7}
	for _, exp := range expected {
		v, exists := h.Pop()
		if !exists || v != exp {
			t.Fatalf("expected %d, got %d", exp, v)
		}
	}
}

func TestHeapInternalArrayInvariant(t *testing.T) {
	h := MakeHeap(func(a, b int) bool { return a < b })

	h.arr = []int{1, 3, 2, 7, 6, 4}

	for i := 1; i < h.Size(); i++ {
		p := parent(i)
		if !h.comp(h.arr[p], h.arr[i]) {
			t.Fatalf(
				"heap invariant violated at index %d (parent=%d, child=%d)",
				i, h.arr[p], h.arr[i],
			)
		}
	}
}

func TestHeapifyFixesViolation(t *testing.T) {
	h := MakeHeap(func(a, b int) bool { return a < b })

	// violate heap property at root
	h.arr = []int{9, 2, 3, 4, 5}
	h.heapify(0)

	if h.arr[0] != 2 {
		t.Fatalf("expected root 2 after heapify, got %d", h.arr[0])
	}
}

func TestHeapSizeAndEmptyPop(t *testing.T) {
	h := MakeHeap(func(a, b int) bool { return a < b })

	h.Push(3)
	h.Push(1)

	if h.Size() != 2 {
		t.Fatalf("expected size 2, got %d", h.Size())
	}

	h.Pop()
	h.Pop()

	if h.Size() != 0 {
		t.Fatalf("expected empty heap, got size %d", h.Size())
	}
}
