package parser

import "fmt"

type TreeNodeType string

const (
	SCOPE       TreeNodeType = "Scope"
	CONDITIONAL TreeNodeType = "Conditional"
	LOOP        TreeNodeType = "Loop"
	FUNCTION    TreeNodeType = "Function"
	STRUCT      TreeNodeType = "Struct"
	OPERATOR    TreeNodeType = "Expression"
	VALUE       TreeNodeType = "Value"
	ARR_IND     TreeNodeType = "Array_Index"
)

const TAB = "  "

type TreeNode interface {
	String(ind string) string
	Type() TreeNodeType
}

type NumberNode struct {
	dataBytes []byte
	numType   string
}

func (n *NumberNode) String(ind string) string {
	return ind + string(n.dataBytes)
}

func (n *NumberNode) Type() TreeNodeType {
	return VALUE
}

func NewNumberNode(dataBytes []byte, numType string) *NumberNode {
	return &NumberNode{dataBytes, numType}
}

type StringNode struct {
	dataBytes []byte
}

func (s *StringNode) String(ind string) string {
	return ind + string(s.dataBytes)
}

func (s *StringNode) Type() TreeNodeType {
	return VALUE
}

func NewStringNode(dataBytes []byte) *StringNode {
	return &StringNode{dataBytes}
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

const (
	PREFIX  = "pre"
	POSTFIX = "post"
)

type PrePostOperatorNode struct {
	opType  string
	op      string
	operand TreeNode
}

func (o *PrePostOperatorNode) String(ind string) string {
	return fmt.Sprintf("%s%s %s\n%s", ind, o.opType, o.op, o.operand.String(ind+TAB))
}

func (o *PrePostOperatorNode) Type() TreeNodeType {
	return OPERATOR
}

func NewPrePostOperatorNode(opType string, op string, operand TreeNode) *PrePostOperatorNode {
	return &PrePostOperatorNode{opType, op, operand}
}

type InfixOperatorNode struct {
	left  TreeNode
	op    string
	right TreeNode
}

func (o *InfixOperatorNode) String(ind string) string {
	return ind + o.op + "\n" + o.left.String(ind+TAB) + "\n" + o.right.String(ind+TAB)
}

func (o *InfixOperatorNode) Type() TreeNodeType {
	return OPERATOR
}

func NewInfixOperatorNode(left TreeNode, op string, right TreeNode) *InfixOperatorNode {
	return &InfixOperatorNode{left, op, right}
}

type FuncCallNode struct {
	callee TreeNode
	args   []TreeNode
}

func (f *FuncCallNode) String(ind string) string {
	args := ""
	for _, arg := range f.args {
		args += arg.String(ind+TAB) + "\n"
	}
	return ind + "call\n" + f.callee.String(ind+TAB) + ind + "args:\n" + args
}

func (f *FuncCallNode) Type() TreeNodeType {
	return FUNCTION
}

func NewFuncCallNode(name TreeNode) *FuncCallNode {
	return &FuncCallNode{name, make([]TreeNode, 0)}
}

func (f *FuncCallNode) arg(arg TreeNode) {
	f.args = append(f.args, arg)
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
