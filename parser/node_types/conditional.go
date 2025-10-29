package node_types

import (
	"fmt"
	"he++/utils"
)

type ConditionalBranch struct {
	Condition TreeNode
	Scope     *ScopeNode
}

type IfNode struct {
	Branches []ConditionalBranch
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
	return &IfNode{branches, *meta}
}
