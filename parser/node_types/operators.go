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
}

func (o *PrePostOperatorNode) String(ind string) string {
	return fmt.Sprintf("%s%s %s\n%s", ind, o.opType, o.Op, o.Operand.String(ind+TAB))
}

func (o *PrePostOperatorNode) Type() TreeNodeType {
	return OPERATOR
}

func NewPrePostOperatorNode(opType string, op string, operand TreeNode) *PrePostOperatorNode {
	return &PrePostOperatorNode{opType, op, operand}
}

type InfixOperatorNode struct {
	Left  TreeNode
	Op    string
	Right TreeNode
}

func (o *InfixOperatorNode) String(ind string) string {
	return ind + o.Op + "\n" + o.Left.String(ind+TAB) + "\n" + o.Right.String(ind+TAB)
}

func (o *InfixOperatorNode) Type() TreeNodeType {
	return OPERATOR
}

func NewInfixOperatorNode(left TreeNode, op string, right TreeNode) *InfixOperatorNode {
	return &InfixOperatorNode{left, op, right}
}

type TernaryOperatorNode struct {
	condition TreeNode
	ifTrue    TreeNode
	ifFalse   TreeNode
}

func (t *TernaryOperatorNode) String(ind string) string {
	return ind + "ternary\n" + t.condition.String(ind+TAB) + "\n" + t.ifTrue.String(ind+TAB) + "\n" + t.ifFalse.String(ind+TAB)
}

func (t *TernaryOperatorNode) Type() TreeNodeType {
	return OPERATOR
}

func NewTernaryNode(condition TreeNode, ifTrue TreeNode, ifFalse TreeNode) *TernaryOperatorNode {
	return &TernaryOperatorNode{condition, ifTrue, ifFalse}
}
