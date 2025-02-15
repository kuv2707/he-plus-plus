package node_types

type StructField struct {
	offset int
	name   struct {
		a string
	}
}
type StructNode struct {
	Name   string
	Fields map[string]StructField
}

func (s StructNode) String(ind string) string {
	return ""
}

func (s StructNode) Type() TreeNodeType {
	return STRUCT
}
