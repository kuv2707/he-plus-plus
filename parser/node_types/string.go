package node_types

import (
	"fmt"
	"he++/utils"
)

type StringNode struct {
	DataBytes []byte
	NodeMetadata
}

func (s *StringNode) String(p *utils.ASTPrinter) {
	p.PushIndent()
	p.WriteLine(fmt.Sprintf("\"%s\"", utils.Yellow(string(s.DataBytes))))
	p.PopIndent()
}

func (s *StringNode) Type() TreeNodeType {
	return VALUE
}

func NewStringNode(dataBytes []byte, meta *NodeMetadata) *StringNode {
	return &StringNode{dataBytes, *meta}
}
