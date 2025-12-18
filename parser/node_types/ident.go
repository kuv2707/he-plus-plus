package node_types

import "he++/utils"

type IdentifierNode struct {
	name string
	DataT DataType
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

func (i *IdentifierNode) ChangeName(name string) {
	i.name = name
}

func NewIdentifierNode(name string, meta *NodeMetadata) *IdentifierNode {
	return &IdentifierNode{name, nil, *meta}
}
