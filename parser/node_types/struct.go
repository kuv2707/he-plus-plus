package node_types

type StructDefnNode struct {
	Name      string
	StructDef *StructType
	NodeMetadata
}

func (s *StructDefnNode) String(ind string) string {
	result := ind + "StructDef: " + s.Name + " " + s.StructDef.Text()
	return result
}

func (s *StructDefnNode) Type() TreeNodeType {
	return STRUCT
}



type StructValueNode struct {
	FieldValues map[string]TreeNode
	NodeMetadata
}

func (s *StructValueNode) String(ind string) string {
	ret := ind + "{\n"
	for k, v := range s.FieldValues {
		ret += ind + k + ":\n" + v.String(ind+TAB) + "\n"
	}
	return ret
}

func (s *StructValueNode) Type() TreeNodeType {
	return STRUCT_VAL
}

func MakeStructValueNode(mp map[string]TreeNode, meta *NodeMetadata) *StructValueNode {
	return &StructValueNode{FieldValues: mp, NodeMetadata: *meta}
}
