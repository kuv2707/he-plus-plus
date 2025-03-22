package node_types

import (
	"he++/globals"
)

type TreeNodeType string

const (
	SCOPE         TreeNodeType = "Scope"
	CONDITIONAL   TreeNodeType = "Conditional"
	LOOP          TreeNodeType = "Loop"
	FUNCTION      TreeNodeType = "Function"
	FUNCTION_CALL TreeNodeType = "Function_Call"
	STRUCT        TreeNodeType = "Struct"
	OPERATOR      TreeNodeType = "Expression"
	VALUE         TreeNodeType = "Value"
	ARR_IND       TreeNodeType = "Array_Index"
	VAR_DECL      TreeNodeType = "Variable_Declaration"
	RETURN        TreeNodeType = "Return"
	ARRAY_DECL    TreeNodeType = "Array_Declaration"
)

const TAB = "  "

type TreeNode interface {
	String(ind string) string
	Type() TreeNodeType
}

type EmptyPlaceholderNode struct {
}

func (e *EmptyPlaceholderNode) String(ind string) string {
	return ind + "<empty>"
}

func (e *EmptyPlaceholderNode) Type() TreeNodeType {
	return TreeNodeType("")
}

type BooleanNode struct {
	dataBytes []byte
}

func (b *BooleanNode) String(ind string) string {
	return ind + string(b.dataBytes)
}

func (b *BooleanNode) Type() TreeNodeType {
	return VALUE
}

func NewBooleanNode(dataBytes []byte) *BooleanNode {
	return &BooleanNode{dataBytes}
}

type IdentifierNode struct {
	name string
}

func (i *IdentifierNode) Name() string {
	return i.name
}

func (i *IdentifierNode) String(ind string) string {
	return ind + i.name
}

func (i *IdentifierNode) Type() TreeNodeType {
	return VALUE
}

func NewIdentifierNode(name string) *IdentifierNode {
	return &IdentifierNode{name}
}

type ArrIndNode struct {
	indexer     TreeNode
	arrProvider TreeNode
}

func (a *ArrIndNode) String(ind string) string {
	ret := ind + " index\n"
	ret += a.indexer.String(ind + TAB)
	ret += "\n"
	ret += a.arrProvider.String(ind + TAB)
	return ret
}

func (f *ArrIndNode) Type() TreeNodeType {
	return ARR_IND
}

func NewArrIndNode(indexer TreeNode, arrProvider TreeNode) *ArrIndNode {
	return &ArrIndNode{indexer, arrProvider}
}

// block related nodes

type StatementsContainer interface {
	AddChild(child TreeNode)
	String(ind string) string
}

type ScopeNode struct {
	Children []TreeNode
}

func MakeScopeNode() *ScopeNode {
	return &ScopeNode{make([]TreeNode, 0)}
}

func (s *ScopeNode) String(ind string) string {
	ret := ind + "scope\n"
	for _, child := range s.Children {
		ret += child.String(ind+TAB) + "\n"
	}
	return globals.RandomColor(ret)
}

func (s *ScopeNode) Type() TreeNodeType {
	return SCOPE
}

func (s *ScopeNode) AddChild(child TreeNode) {
	s.Children = append(s.Children, child)
}

type SourceFileNode struct {
	fileName string
	filePath string
	Children []TreeNode
	//todo: store exports of this file
}

func MakeSourceFileNode() *SourceFileNode {
	return &SourceFileNode{Children: make([]TreeNode, 0)}
}

func (s *SourceFileNode) String(ind string) string {
	ret := ind + "source file\n"
	for _, child := range s.Children {
		ret += child.String(ind+TAB) + "\n"
	}
	return ret
}

func (s *SourceFileNode) Type() TreeNodeType {
	return SCOPE
}

func (s *SourceFileNode) AddChild(child TreeNode) {
	s.Children = append(s.Children, child)
}
