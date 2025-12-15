package node_types

import (
	"fmt"
	"he++/utils"
)

type NumberType string

const (
	FLOAT_NUM NumberType = "float"
	INT_NUM   NumberType = "int"
)

// expression related nodes
type NumberNode struct {
	RawNumBytes []byte
	NumType     NumberType
	NodeMetadata
}

func (n *NumberNode) String(p *utils.ASTPrinter) {
	p.PushIndent()
	p.WriteLine(utils.Blue(fmt.Sprint(n.RawNumBytes)))
	p.PopIndent()
}

func (n *NumberNode) Type() TreeNodeType {
	return VALUE
}

func NewNumberNode(RawNumBytes []byte, numType NumberType, meta *NodeMetadata) *NumberNode {
	return &NumberNode{RawNumBytes, numType, *meta}
}
