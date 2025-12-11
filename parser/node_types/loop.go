package node_types

import "he++/utils"

var loopSeq = 0

type LoopNode struct {
	Seq         int
	Initializer TreeNode
	Condition   TreeNode
	Updater     TreeNode
	Scope       *ScopeNode
	NodeMetadata
}

func (l *LoopNode) String(p *utils.ASTPrinter) {
	p.PushIndent()
	p.WriteLine(utils.Underline("loop"))
	p.WriteLine("initializer:")
	l.Initializer.String(p)
	p.WriteLine("condition:")
	l.Condition.String(p)
	p.WriteLine("updater:")
	l.Updater.String(p)
	l.Scope.String(p)
	p.PopIndent()
}

func (l *LoopNode) Type() TreeNodeType {
	return LOOP
}

func MakeLoopNode(initializer TreeNode, condition TreeNode, updater TreeNode, scope *ScopeNode, meta *NodeMetadata) *LoopNode {
	loopSeq++
	return &LoopNode{loopSeq, initializer, condition, updater, scope, *meta}
}
