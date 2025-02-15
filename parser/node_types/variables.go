package node_types

type VariableDeclarationNode struct {
	declarations []TreeNode
}

func MakeVariableDeclarationNode() *VariableDeclarationNode {
	return &VariableDeclarationNode{make([]TreeNode, 0)}
}

func (v *VariableDeclarationNode) String(ind string) string {
	ret := ind + "var decl\n"
	for _, decl := range v.declarations {
		ret += decl.String(ind+TAB) + "\n"
	}
	return ret
}

func (v *VariableDeclarationNode) Type() TreeNodeType {
	return VAR_DECL
}

func (v *VariableDeclarationNode) AddDeclaration(decl TreeNode) {
	v.declarations = append(v.declarations, decl)
}

type PointerToType struct {
	UnderlyingType TreeNode
}
