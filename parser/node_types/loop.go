package node_types

type LoopNode struct {
	initializer TreeNode
	condition   TreeNode
	updater     TreeNode
	scope       *ScopeNode
	NodeMetadata
}

func (l *LoopNode) String(ind string) string {
	ret := ind + "loop\n"
	ret += l.initializer.String(ind+TAB) + "\n"
	ret += l.condition.String(ind+TAB) + "\n"
	ret += l.updater.String(ind+TAB) + "\n"
	ret += l.scope.String(ind+TAB) + "\n"
	return ret
}

func (l *LoopNode) Type() TreeNodeType {
	return LOOP
}

func MakeLoopNode(initializer TreeNode, condition TreeNode, updater TreeNode, scope *ScopeNode, meta *NodeMetadata) *LoopNode {
	return &LoopNode{initializer, condition, updater, scope, *meta}
}
