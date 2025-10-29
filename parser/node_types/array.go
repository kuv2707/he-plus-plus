package node_types

import "he++/utils"

type ArrayDeclarationNode struct {
	Elems []TreeNode
	DataT DataType
	NodeMetadata
}

func (ad *ArrayDeclarationNode) String(p *utils.ASTPrinter) {
	p.PushIndent()
	p.WriteLine(utils.Underline("ArrayDeclaration:"))
	for _, elem := range ad.Elems {
		elem.String(p)
	}
	p.PopIndent()
}

func (ad *ArrayDeclarationNode) Type() TreeNodeType {
	return ARRAY_DECL
}

func MakeArrayDeclarationNode(elems []TreeNode, dt DataType, meta *NodeMetadata) *ArrayDeclarationNode {
	return &ArrayDeclarationNode{
		Elems:        elems,
		DataT:        dt,
		NodeMetadata: *meta,
	}
}
