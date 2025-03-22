package node_types

type ArrayDeclaration struct {
	Elems []TreeNode
	DataT DataType
}

func (ad *ArrayDeclaration) String(ind string) string {
	result := ind + "ArrayDeclaration:\n"
	for _, elem := range ad.Elems {
		result += elem.String(ind + TAB) + "\n"
	}
	return result
}

func (ad *ArrayDeclaration) Type() TreeNodeType {
	return ARRAY_DECL
}
