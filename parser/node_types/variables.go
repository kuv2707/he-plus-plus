package node_types

type VariableDeclarationNode struct {
	Declarations []TreeNode
	DataT    DataType
	NodeMetadata
}

func MakeVariableDeclarationNode(decls []TreeNode, dt DataType, meta *NodeMetadata) *VariableDeclarationNode {
	return &VariableDeclarationNode{decls, dt, *meta}
}

func (v *VariableDeclarationNode) String(ind string) string {
	ret := ind + "var decl\n"
	for _, decl := range v.Declarations {
		ret += decl.String(ind+TAB) + "\n"
	}
	ret += ind + "DataType: " + v.DataT.Text()
	return ret
}

func (v *VariableDeclarationNode) Type() TreeNodeType {
	return VAR_DECL
}

func (v *VariableDeclarationNode) AddDeclaration(decl TreeNode) {
	v.Declarations = append(v.Declarations, decl)
}

func (v *VariableDeclarationNode) SetDataType(dt DataType) {
	v.DataT = dt
}

type PointerToType struct {
	UnderlyingType TreeNode
}
