package node_types

import (
	"fmt"
)

type StringNode struct {
	DataBytes []byte
}

func (s *StringNode) String(ind string) string {
	return ind + fmt.Sprintf("\"%s\"", string(s.DataBytes))
}

func (s *StringNode) Type() TreeNodeType {
	return VALUE
}

func NewStringNode(dataBytes []byte) *StringNode {
	return &StringNode{dataBytes}
}
