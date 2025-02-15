package node_types

type IfNode struct {
	condition TreeNode
	ifScope   TreeNode
	elseScope TreeNode
}

func (i *IfNode) String(ind string) string {
	return ind + "if\n" + i.condition.String(ind+TAB) + "\n" + ind + "then\n" + i.ifScope.String(ind+TAB) + ind + "else\n" + i.elseScope.String(ind+TAB)
}

func (i *IfNode) Type() TreeNodeType {
	return CONDITIONAL
}

func MakeIfNode(condition TreeNode, ifScope TreeNode, elseScope TreeNode) *IfNode {
	return &IfNode{condition, ifScope, elseScope}
}
