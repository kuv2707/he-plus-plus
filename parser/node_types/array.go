package node_types

type ArrayDeclarationNode struct {
	Elems []TreeNode
	DataT DataType
	NodeMetadata
}

func (ad *ArrayDeclarationNode) String(ind string) string {
	result := ind + "ArrayDeclaration:\n"
	for _, elem := range ad.Elems {
		result += elem.String(ind+TAB) + "\n"
	}
	return result
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
