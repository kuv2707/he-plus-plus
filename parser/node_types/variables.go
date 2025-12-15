package node_types

import "he++/utils"

type VariableDeclarationNode struct {
	Declarations []TreeNode
	DataT        DataType
	NodeMetadata
}

func MakeVariableDeclarationNode(decls []TreeNode, dt DataType, meta *NodeMetadata) *VariableDeclarationNode {
	return &VariableDeclarationNode{decls, dt, *meta}
}

func (v *VariableDeclarationNode) String(p *utils.ASTPrinter) {
	p.PushIndent()
	p.WriteLine("VarDecl")
	p.WriteLine("DataType: " + utils.Cyan(v.DataT.Text()))
	for _, decl := range v.Declarations {
		decl.String(p)
	}
	p.PopIndent()
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

