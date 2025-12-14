package utils


import "slices"

type Queue[T any] struct {
    buf  []T
    head int
}

func MakeQueue[T any]() Queue[T] {
	return Queue[T]{buf: nil, head: 0}
}

func (q *Queue[T]) Push(x T) {
    q.buf = append(q.buf, x)
}

func (q *Queue[T]) Empty() bool {
    return q.head >= len(q.buf)
}

func (q *Queue[T]) Pop() T {
    if q.Empty() {
        panic("pop from empty queue")
    }
    x := q.buf[q.head]
    q.head++

    // compaction
    if q.head > 1024 && q.head*2 >= len(q.buf) {
        q.buf = slices.Clone(q.buf[q.head:])
        q.head = 0
    }
    return x
}