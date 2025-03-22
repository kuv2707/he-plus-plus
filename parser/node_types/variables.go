package node_types

type VariableDeclarationNode struct {
	Declarations []TreeNode
	DataT    DataType
}

func MakeVariableDeclarationNode() *VariableDeclarationNode {
	return &VariableDeclarationNode{make([]TreeNode, 0), &ErrorType{Message: "Undefined"}}
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
