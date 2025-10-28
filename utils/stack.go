package utils

import "fmt"

type Stack[T any] struct {
	items []T
	zero  T
}

func MakeStack[T any](items... T) Stack[T] {
	var zero T
	return Stack[T]{items, zero}
}

func (s *Stack[T]) Push(item T) {
	s.items = append(s.items, item)
}

func (s *Stack[T]) Pop() (T, bool) {
	if len(s.items) == 0 {
		return s.zero, false
	}
	item := s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]
	return item, true
}
func (s *Stack[T]) Peek() (T, bool) {
	if len(s.items) == 0 {
		return s.zero, false
	}
	return s.items[len(s.items)-1], true
}

func (s *Stack[T]) Len() int {
	return len(s.items)
}

func (s *Stack[T]) IsEmpty() bool {
	return len(s.items) == 0
}

func (s *Stack[T]) GetStackItems() []T {
	return s.items
}

func (s *Stack[T]) PrintStack() bool {
	fmt.Println("---------")
	for i := len(s.items) - 1; i >= 0; i-- {
		fmt.Println(s.items[i])
	}
	fmt.Println("---------")
	return true
}
