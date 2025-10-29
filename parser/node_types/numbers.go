package node_types

import "he++/utils"

// expression related nodes
type NumberNode struct {
	DataBytes []byte
	NumType   string //todo: constrain to "int" and "float" using enum etc
	NodeMetadata
}

func (n *NumberNode) String(ind string) string {
	return ind + utils.Blue(string(n.DataBytes))
}

func (n *NumberNode) Type() TreeNodeType {
	return VALUE
}

func NewNumberNode(dataBytes []byte, numType string, meta *NodeMetadata) *NumberNode {
	return &NumberNode{dataBytes, numType, *meta}
}
