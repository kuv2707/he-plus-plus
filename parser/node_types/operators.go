package node_types

import (
	"he++/utils"
)

const (
	PREFIX  = "pre"
	POSTFIX = "post"
)

type PrePostOperatorNode struct {
	opType  string
	Op      string
	Operand TreeNode
	NodeMetadata
}

func (o *PrePostOperatorNode) String(p *utils.ASTPrinter) {
	p.PushIndent()
	p.WriteLine(utils.Bold(o.opType) + " " + utils.Magenta(o.Op))
	o.Operand.String(p)
	p.PopIndent()
}

func (o *PrePostOperatorNode) Type() TreeNodeType {
	return OPERATOR
}

func NewPrePostOperatorNode(opType string, op string, operand TreeNode, meta *NodeMetadata) *PrePostOperatorNode {
	return &PrePostOperatorNode{opType, op, operand, *meta}
}

type InfixOperatorNode struct {
	Left  TreeNode
	Op    string
	Right TreeNode
	NodeMetadata
}

func (o *InfixOperatorNode) String(p *utils.ASTPrinter) {
	p.PushIndent()
	p.WriteLine(utils.Magenta(o.Op))
	o.Left.String(p)
	o.Right.String(p)
	p.PopIndent()
}

func (o *InfixOperatorNode) Type() TreeNodeType {
	return OPERATOR
}

func NewInfixOperatorNode(left TreeNode, op string, right TreeNode, meta *NodeMetadata) *InfixOperatorNode {
	return &InfixOperatorNode{left, op, right, *meta}
}

type TernaryOperatorNode struct {
	condition TreeNode
	ifTrue    TreeNode
	ifFalse   TreeNode
	NodeMetadata
}

func (t *TernaryOperatorNode) String(p *utils.ASTPrinter) {
	p.PushIndent()
	p.WriteLine("ternary")
	t.condition.String(p)
	t.ifTrue.String(p)
	t.ifFalse.String(p)
	p.PopIndent()
}

func (t *TernaryOperatorNode) Type() TreeNodeType {
	return OPERATOR
}

func NewTernaryNode(condition TreeNode, ifTrue TreeNode, ifFalse TreeNode) *TernaryOperatorNode {
	return &TernaryOperatorNode{condition, ifTrue, ifFalse, NodeMetadata{}}
}

type ArrIndNode struct {
	ArrProvider TreeNode
	Indexer     TreeNode
	NodeMetadata
}

func (a *ArrIndNode) String(p *utils.ASTPrinter) {
	p.PushIndent()
	p.WriteLine(utils.Magenta("index"))
	a.Indexer.String(p)
	a.ArrProvider.String(p)
	p.PopIndent()
}

func (f *ArrIndNode) Type() TreeNodeType {
	return ARR_IND
}

func NewArrIndNode(arrProvider TreeNode, indexer TreeNode, meta *NodeMetadata) *ArrIndNode {
	return &ArrIndNode{ArrProvider: arrProvider, Indexer: indexer, NodeMetadata: *meta}
}

type FuncCallNode struct {
	Callee TreeNode
	Args   []TreeNode
	NodeMetadata
}

func (f *FuncCallNode) String(p *utils.ASTPrinter) {
	p.PushIndent()
	p.WriteLine(utils.Magenta("call"))
	p.PushIndent()
	p.WriteLine("callee:")
	f.Callee.String(p)
	p.WriteLine("args:")
	for _, arg := range f.Args {
		arg.String(p)
	}
	p.PopIndent()
	p.PopIndent()
}

func (f *FuncCallNode) Type() TreeNodeType {
	return FUNCTION_CALL
}

func NewFuncCallNode(name TreeNode, meta *NodeMetadata) *FuncCallNode {
	return &FuncCallNode{name, make([]TreeNode, 0), *meta}
}
