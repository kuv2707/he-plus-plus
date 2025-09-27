package node_types

import "fmt"



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

func (o *PrePostOperatorNode) String(ind string) string {
	return fmt.Sprintf("%s%s %s\n%s", ind, o.opType, o.Op, o.Operand.String(ind+TAB))
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

func (o *InfixOperatorNode) String(ind string) string {
	return ind + o.Op + "\n" + o.Left.String(ind+TAB) + "\n" + o.Right.String(ind+TAB)
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

func (t *TernaryOperatorNode) String(ind string) string {
	return ind + "ternary\n" + t.condition.String(ind+TAB) + "\n" + t.ifTrue.String(ind+TAB) + "\n" + t.ifFalse.String(ind+TAB)
}

func (t *TernaryOperatorNode) Type() TreeNodeType {
	return OPERATOR
}

func NewTernaryNode(condition TreeNode, ifTrue TreeNode, ifFalse TreeNode) *TernaryOperatorNode {
	return &TernaryOperatorNode{condition, ifTrue, ifFalse, NodeMetadata{}}
}
