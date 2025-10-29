package node_types

import "he++/utils"

type StructDefnNode struct {
	Name      string
	StructDef *StructType
	NodeMetadata
}

func (s *StructDefnNode) String(p *utils.ASTPrinter) {
	p.PushIndent()
	p.WriteLine(utils.Underline("StructDef:"))
	p.WriteLine(utils.Green(s.Name))
	p.WriteLine(s.StructDef.Text())
	p.PopIndent()
}

func (s *StructDefnNode) Type() TreeNodeType {
	return STRUCT
}

type StructValueNode struct {
	FieldValues map[string]TreeNode
	NodeMetadata
}

func (s *StructValueNode) String(p *utils.ASTPrinter) {
	p.PushIndent()
	p.WriteLine("{")
	p.PushIndent()
	for k, v := range s.FieldValues {
		p.WriteLine(k + ":")
		v.String(p)
	}
	p.PopIndent()
	p.WriteLine("}")
	p.PopIndent()
}

func (s *StructValueNode) Type() TreeNodeType {
	return STRUCT_VAL
}

func MakeStructValueNode(mp map[string]TreeNode, meta *NodeMetadata) *StructValueNode {
	return &StructValueNode{FieldValues: mp, NodeMetadata: *meta}
}
