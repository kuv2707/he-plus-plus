package utils

type Comparator[K any] func(a, b K) bool

type Heap[T any] struct {
	comp Comparator[T]
	arr  []T
}

func MakeHeap[T any]() Heap[T] {
	return Heap[T]{}
}

func (h *Heap[T]) Push(k T) {
	h.arr = append(h.arr, k)
	i := h.Size() - 1
	for i != 0 && !h.comp(h.arr[parent(i)], h.arr[i]) {
		h.swapInds(i, parent(i))
		i = parent(i)
	}
}

func (h *Heap[T]) Pop() (T, bool) {
	var zero T
	if h.Size() == 0 {
		return zero, false
	}
	ret := h.arr[0]
	h.swapInds(0, h.Size()-1)
	h.arr = h.arr[:h.Size()-1]
	if h.Size() > 0 {
		h.heapify(0)
	}
	return ret, true
}

func (h *Heap[T]) heapify(i int) {
	l := left(i)
	r := right(i)
	min := i
	if l < h.Size() && h.comp(h.arr[l], h.arr[min]) {
		min = l
	}
	if r < h.Size() && h.comp(h.arr[r], h.arr[min]) {
		min = r
	}

	if min == i {
		return
	}

	h.swapInds(i, min)
	h.heapify(min)
}

func (h *Heap[T]) Size() int {
	return len(h.arr)
}

func (h *Heap[T]) Peek() T {
	return h.arr[0]
}

func (h *Heap[T]) SetComparator(c Comparator[T]) {
	h.comp = c
}

func (h *Heap[T]) swapInds(i, min int) {
	t := h.arr[min]
	h.arr[min] = h.arr[i]
	h.arr[i] = t
}

func parent(i int) int {
	return (i - 1) / 2
}

func left(i int) int {
	return 2*i + 1
}

func right(i int) int {
	return 2*i + 2
}
