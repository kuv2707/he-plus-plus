package node_types

import (
	"he++/lexer"
	"he++/utils"
)

type BooleanNode struct {
	boolVal bool
	NodeMetadata
}

func (b *BooleanNode) String(p *utils.ASTPrinter) {
	p.PushIndent()
	if b.boolVal {
		p.WriteLine(utils.Blue(lexer.TRUE))
	} else {
		p.WriteLine(utils.Blue(lexer.FALSE))
	}
	p.PopIndent()
}

func (b *BooleanNode) Type() TreeNodeType {
	return VALUE
}

func NewBooleanNode(val bool, meta *NodeMetadata) *BooleanNode {
	return &BooleanNode{val, *meta}
}
