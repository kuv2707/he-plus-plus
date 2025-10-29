package node_types

import "he++/utils"

type IdentifierNode struct {
	name string
	NodeMetadata
}

func (i *IdentifierNode) Name() string {
	return i.name
}

func (i *IdentifierNode) String(p *utils.ASTPrinter) {
	p.PushIndent()
	p.WriteLine(utils.Green(i.name))
	p.PopIndent()
}

func (i *IdentifierNode) Type() TreeNodeType {
	return VALUE
}

func NewIdentifierNode(name string, meta *NodeMetadata) *IdentifierNode {
	return &IdentifierNode{name, *meta}
}
