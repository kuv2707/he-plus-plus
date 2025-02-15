package node_types


// expression related nodes
type NumberNode struct {
	dataBytes []byte
	numType   string
}

func (n *NumberNode) String(ind string) string {
	return ind + string(n.dataBytes)
}

func (n *NumberNode) Type() TreeNodeType {
	return VALUE
}

func NewNumberNode(dataBytes []byte, numType string) *NumberNode {
	return &NumberNode{dataBytes, numType}
}
