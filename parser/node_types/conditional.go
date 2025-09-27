package node_types

import "fmt"

type ConditionalBranch struct {
	Condition TreeNode
	Scope      *ScopeNode
}

type IfNode struct {
	Branches []ConditionalBranch
	NodeMetadata
}

func (i *IfNode) String(ind string) string {
	result := ind + fmt.Sprintf("conditional branches (%d):", len(i.Branches))
	for idx, branch := range i.Branches {
		result += ind + "branch <" + fmt.Sprint(idx) + ">\n"
		result += branch.Condition.String(ind+TAB) + "\n"
		result += branch.Scope.String(ind + TAB)
	}
	return result
}

func (i *IfNode) Type() TreeNodeType {
	return CONDITIONAL
}

func MakeIfNode(branches []ConditionalBranch, meta *NodeMetadata) *IfNode {
	return &IfNode{branches, *meta}
}
