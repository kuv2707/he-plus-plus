package parser

type TreeNode struct {
	Label       string
	Description string
	Children    []*TreeNode
	Properties  map[string]*TreeNode
	LineNo      int
}
