package node_types

import "he++/utils"

type BooleanNode struct {
	dataBytes []byte
	NodeMetadata
}

func (b *BooleanNode) String(p *utils.ASTPrinter) {
	p.PushIndent()
	p.WriteLine(utils.Blue(string(b.dataBytes)))
	p.PopIndent()
}

func (b *BooleanNode) Type() TreeNodeType {
	return VALUE
}

func NewBooleanNode(dataBytes []byte, meta *NodeMetadata) *BooleanNode {
	return &BooleanNode{dataBytes, *meta}
}
