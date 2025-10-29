package node_types

import (
	"fmt"
	"he++/utils"
)

type StringNode struct {
	DataBytes []byte
	NodeMetadata
}

func (s *StringNode) String(ind string) string {
	return ind + fmt.Sprintf("\"%s\"", utils.Yellow(string(s.DataBytes)))
}

func (s *StringNode) Type() TreeNodeType {
	return VALUE
}

func NewStringNode(dataBytes []byte, meta *NodeMetadata) *StringNode {
	return &StringNode{dataBytes, *meta}
}
