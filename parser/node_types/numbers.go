package node_types

import "he++/utils"

type NumberType string

const (
	FLOAT_NUMBER NumberType = "float"
	INT32_NUMBER NumberType = "int32"
)

// expression related nodes
type NumberNode struct {
	DataBytes []byte
	NumType   NumberType
	NodeMetadata
}

func (n *NumberNode) String(p *utils.ASTPrinter) {
	p.PushIndent()
	p.WriteLine(utils.Blue(string(n.DataBytes)))
	p.PopIndent()
}

func (n *NumberNode) Type() TreeNodeType {
	return VALUE
}

func NewNumberNode(dataBytes []byte, numType NumberType, meta *NodeMetadata) *NumberNode {
	return &NumberNode{dataBytes, numType, *meta}
}
