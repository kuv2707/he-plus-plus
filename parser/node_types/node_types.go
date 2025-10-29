package node_types

import (
	"fmt"
	"he++/utils"
	_ "he++/utils"
)

type TreeNodeType string

const (
	SCOPE         TreeNodeType = "Scope"
	CONDITIONAL   TreeNodeType = "Conditional"
	LOOP          TreeNodeType = "Loop"
	FUNCTION      TreeNodeType = "Function"
	FUNCTION_CALL TreeNodeType = "Function_Call"
	STRUCT        TreeNodeType = "Struct_Declaration"
	STRUCT_VAL    TreeNodeType = "Struct_Value"
	OPERATOR      TreeNodeType = "Expression"
	VALUE         TreeNodeType = "Value"
	ARR_IND       TreeNodeType = "Array_Index"
	VAR_DECL      TreeNodeType = "Variable_Declaration"
	RETURN        TreeNodeType = "Return"
	ARRAY_DECL    TreeNodeType = "Array_Declaration"
)

const TAB = "  "

type LineRange struct {
	Start int
	End   int
}

// Color codes:
// Data type names : cyan
// Numbers & bools : blue
// String literals : yellow
// Ident names     : green
// Operators       : magenta

type TreeNode interface {
	String(p *utils.ASTPrinter)
	Type() TreeNodeType
	Range() LineRange
}

type NodeMetadata struct {
	lr LineRange
}

func (m *NodeMetadata) Range() LineRange {
	return m.lr
}

func MakeMetadata(l int, r int) *NodeMetadata {
	return &NodeMetadata{lr: LineRange{Start: l, End: r}}
}

type EmptyPlaceholderNode struct {
	NodeMetadata
}

func (e *EmptyPlaceholderNode) String(p *utils.ASTPrinter) {
	p.PushIndent()
	p.WriteLine("<empty>")
	p.PopIndent()
}

func (e *EmptyPlaceholderNode) Type() TreeNodeType {
	return TreeNodeType("")
}

type StatementsContainer interface {
	AddChild(child TreeNode)
	String(p *utils.ASTPrinter)
}

type ScopeNode struct {
	Children []TreeNode
	NodeMetadata
}

func MakeScopeNode() *ScopeNode {
	return &ScopeNode{make([]TreeNode, 0), NodeMetadata{}}
}

func (s *ScopeNode) String(p *utils.ASTPrinter) {
	p.PushIndent()
	p.WriteLine(utils.Underline("scope"))
	for _, child := range s.Children {
		child.String(p)
	}
	p.PopIndent()
}

func (s *ScopeNode) Type() TreeNodeType {
	return SCOPE
}

func (s *ScopeNode) AddChild(child TreeNode) {
	s.Children = append(s.Children, child)
}

type SourceFileNode struct {
	FilePath string
	Children []TreeNode
	NodeMetadata
	//todo: store exports of this file
}

func MakeSourceFileNode(path string) *SourceFileNode {
	return &SourceFileNode{FilePath: path, Children: make([]TreeNode, 0), NodeMetadata: NodeMetadata{}}
}

func (s *SourceFileNode) String(p *utils.ASTPrinter) {
	p.PushIndent()
	p.WriteLine(fmt.Sprintf("%s %s", utils.Underline("File:"), utils.Underline(utils.BoldWhite(s.FilePath))))
	for _, child := range s.Children {
		child.String(p)
	}
	p.PopIndent()
}

func (s *SourceFileNode) Type() TreeNodeType {
	return SCOPE
}

func (s *SourceFileNode) AddChild(child TreeNode) {
	s.Children = append(s.Children, child)
}
