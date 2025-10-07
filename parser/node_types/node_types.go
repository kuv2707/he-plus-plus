package node_types

import (
	"he++/utils"
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

type TreeNode interface {
	String(ind string) string
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

func (e *EmptyPlaceholderNode) String(ind string) string {
	return ind + "<empty>"
}

func (e *EmptyPlaceholderNode) Type() TreeNodeType {
	return TreeNodeType("")
}

type BooleanNode struct {
	dataBytes []byte
	NodeMetadata
}

func (b *BooleanNode) String(ind string) string {
	return ind + string(b.dataBytes)
}

func (b *BooleanNode) Type() TreeNodeType {
	return VALUE
}

func NewBooleanNode(dataBytes []byte, meta *NodeMetadata) *BooleanNode {
	return &BooleanNode{dataBytes, *meta}
}

type IdentifierNode struct {
	name string
	NodeMetadata
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

func NewIdentifierNode(name string, meta *NodeMetadata) *IdentifierNode {
	return &IdentifierNode{name, *meta}
}

type ArrIndNode struct {
	indexer     TreeNode
	arrProvider TreeNode
	NodeMetadata
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

func NewArrIndNode(indexer TreeNode, arrProvider TreeNode, meta *NodeMetadata) *ArrIndNode {
	return &ArrIndNode{indexer, arrProvider, *meta}
}

// block related nodes

type StatementsContainer interface {
	AddChild(child TreeNode)
	String(ind string) string
}

type ScopeNode struct {
	Children []TreeNode
	NodeMetadata
}

func MakeScopeNode() *ScopeNode {
	return &ScopeNode{make([]TreeNode, 0), NodeMetadata{}}
}

func (s *ScopeNode) String(ind string) string {
	ret := ind + "scope\n"
	for _, child := range s.Children {
		ret += child.String(ind+TAB) + "\n"
	}
	return utils.RandomColor(ret)
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

func (s *SourceFileNode) String(ind string) string {
	ret := ind + "File: " + s.FilePath + "\n"
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
