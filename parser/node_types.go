package parser

type Node interface{
	IsNode()
}
type TreeNode struct {
	Label       string
	Description string
	Children    []*TreeNode
	Properties map[string]*TreeNode
}
