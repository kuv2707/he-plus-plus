package node_types


// expression related nodes
type NumberNode struct {
	DataBytes []byte
	NumType   string //todo: constrain to "int" and "float"
}

func (n *NumberNode) String(ind string) string {
	return ind + string(n.DataBytes)
}

func (n *NumberNode) Type() TreeNodeType {
	return VALUE
}

func NewNumberNode(dataBytes []byte, numType string) *NumberNode {
	return &NumberNode{dataBytes, numType}
}
