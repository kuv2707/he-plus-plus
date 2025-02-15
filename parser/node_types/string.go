package node_types


type StringNode struct {
	dataBytes []byte
}

func (s *StringNode) String(ind string) string {
	return ind + string(s.dataBytes)
}

func (s *StringNode) Type() TreeNodeType {
	return VALUE
}

func NewStringNode(dataBytes []byte) *StringNode {
	return &StringNode{dataBytes}
}
