package node_types

import "he++/utils"

type LoopNode struct {
	initializer TreeNode
	condition   TreeNode
	updater     TreeNode
	scope       *ScopeNode
	NodeMetadata
}

func (l *LoopNode) String(p *utils.ASTPrinter) {
	p.PushIndent()
	p.WriteLine(utils.Underline("loop"))
	p.WriteLine("initializer:")
	l.initializer.String(p)
	p.WriteLine("condition:")
	l.condition.String(p)
	p.WriteLine("updater:")
	l.updater.String(p)
	l.scope.String(p)
	p.PopIndent()
}

func (l *LoopNode) Type() TreeNodeType {
	return LOOP
}

func MakeLoopNode(initializer TreeNode, condition TreeNode, updater TreeNode, scope *ScopeNode, meta *NodeMetadata) *LoopNode {
	return &LoopNode{initializer, condition, updater, scope, *meta}
}
