package node_types

type VariableDeclarationNode struct {
	declarations []TreeNode
	data_type    DataType
}

func MakeVariableDeclarationNode() *VariableDeclarationNode {
	return &VariableDeclarationNode{make([]TreeNode, 0), DataType{}}
}

func (v *VariableDeclarationNode) String(ind string) string {
	ret := ind + "var decl\n"
	for _, decl := range v.declarations {
		ret += decl.String(ind+TAB) + "\n"
	}
	ret += ind + "DataType: " + v.data_type.Text
	return ret
}

func (v *VariableDeclarationNode) Type() TreeNodeType {
	return VAR_DECL
}

func (v *VariableDeclarationNode) AddDeclaration(decl TreeNode) {
	v.declarations = append(v.declarations, decl)
}

func (v *VariableDeclarationNode) SetDataType(dt DataType) {
	v.data_type = dt
}

type PointerToType struct {
	UnderlyingType TreeNode
}
