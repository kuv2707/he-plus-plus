package node_types

import "he++/utils"

type ArrayDeclarationNode struct {
	SizeNode  TreeNode
	Elems []TreeNode
	DataT DataType // this is the type of the array elements not the array itself
	NodeMetadata
}

func (ad *ArrayDeclarationNode) String(p *utils.ASTPrinter) {
	p.PushIndent()
	p.WriteLine(utils.Underline("ArrayDeclaration:"))
	p.PushIndent()
	p.WriteLine("size")
	ad.SizeNode.String(p)
	p.PopIndent()
	for _, elem := range ad.Elems {
		elem.String(p)
	}
	p.PopIndent()
}

func (ad *ArrayDeclarationNode) Type() TreeNodeType {
	return ARRAY_DECL
}

func MakeArrayDeclarationNode(size TreeNode, elems []TreeNode, dt DataType, meta *NodeMetadata) *ArrayDeclarationNode {
	return &ArrayDeclarationNode{
		SizeNode:         size,
		Elems:        elems,
		DataT:        dt,
		NodeMetadata: *meta,
	}
}
