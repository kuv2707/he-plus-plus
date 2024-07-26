package utils

import "fmt"

//write code for stack data structure storing generic types

type Stack struct {
	items []interface{}
}

func MakeStack() Stack {
	return Stack{make([]interface{}, 0)}
}

func (s *Stack) Push(item interface{}) {
	s.items = append(s.items, item)
}

func (s *Stack) Pop() interface{} {
	if len(s.items) == 0 {
		return nil
	}
	item := s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]
	return item
}
func (s *Stack) Peek() interface{} {
	if len(s.items) == 0 {
		return nil
	}
	return s.items[len(s.items)-1]
}

func (s *Stack) Get(i int) interface{} {
	if len(s.items) == 0 {
		return nil
	}
	return s.items[i]
}

func (s *Stack) Len() int {
	return len(s.items)
}

func (s *Stack) IsEmpty() bool {
	return len(s.items) == 0
}

func (s *Stack) GetStack() []interface{} {
	return s.items
}

func (s *Stack) PrintStack() bool {
	fmt.Println("---------")
	for i := len(s.items) - 1; i >= 0; i-- {
		fmt.Println(s.items[i])
	}
	fmt.Println("---------")
	return true
}