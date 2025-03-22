package node_types


type StringNode struct {
	DataBytes []byte
}

func (s *StringNode) String(ind string) string {
	return ind + string(s.DataBytes)
}

func (s *StringNode) Type() TreeNodeType {
	return VALUE
}

func NewStringNode(dataBytes []byte) *StringNode {
	return &StringNode{dataBytes}
}
