package node_types

import "he++/utils"

// expression related nodes
type NumberNode struct {
	DataBytes []byte
	NumType   string //todo: constrain to "int" and "float" using enum etc
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

func NewNumberNode(dataBytes []byte, numType string, meta *NodeMetadata) *NumberNode {
	return &NumberNode{dataBytes, numType, *meta}
}
