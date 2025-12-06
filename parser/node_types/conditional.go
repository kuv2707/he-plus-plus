package node_types

import (
	"fmt"
	"he++/utils"
)

type ConditionalBranch struct {
	Condition TreeNode
	Scope     *ScopeNode
}

var ifSeq int = 0
type IfNode struct {
	Branches []ConditionalBranch
	Exhaustive bool
	Seq int
	NodeMetadata
}

func (i *IfNode) String(p *utils.ASTPrinter) {
	p.PushIndent()
	p.WriteLine(utils.Underline(fmt.Sprintf("ConditionalBranches(%d):", len(i.Branches))))

	for idx, branch := range i.Branches {
		p.WriteLine("branch <" + fmt.Sprint(idx) + ">")
		branch.Condition.String(p)
		branch.Scope.String(p)
	}
	p.PopIndent()
}

func (i *IfNode) Type() TreeNodeType {
	return CONDITIONAL
}

func MakeIfNode(branches []ConditionalBranch, meta *NodeMetadata) *IfNode {
	ifSeq++
	return &IfNode{branches, false, ifSeq, *meta}
}
