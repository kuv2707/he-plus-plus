package node_types

type StructField struct {
	Name      string
	FieldType TreeNode
}
type StructNode struct {
	Name   string
	Fields map[string]StructField
}

func (s StructNode) String(ind string) string {
	result := ind + "StructNode: " + s.Name + "\n"
	for _, field := range s.Fields {
		result += ind + "Field: " + field.Name + ", Type: " + field.FieldType.String(ind+" ") + "\n"
	}
	return result
}

func (s StructNode) Type() TreeNodeType {
	return STRUCT
}
