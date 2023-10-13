package parser

import (
	"fmt"
	"toylingo/lexer"
)

type TreeNode struct {
	Val      string
	Children []*TreeNode
}

func ParseTree(tokens *lexer.Node) *TreeNode {
	tokensArr := make([]lexer.TokenType, 0)
	for node := tokens; node != nil; node = node.Next {
		tokensArr = append(tokensArr, node.Val)
	}
	treeNode := &TreeNode{"root", make([]*TreeNode, 0)}
outer:
	for i := 0; i < len(tokensArr); i++ {
		if tokensArr[i].Type == "LET" {
			for j := i + 1; j < len(tokensArr); j++ {
				if tokensArr[j].Type == "SEMICOLON" {
					// treeNode.Children = append(treeNode.Children, parseExpression(tokensArr[i+1:j]))
					fmt.Println(tokensArr[i+1 : j])
					PrintTree(parseExpression(tokensArr[i+1:j]), "")
					break outer
				}
			}
		}

	}
	return treeNode
}
//prints prefix notation of tree
func PrintTree(treeNode *TreeNode, space string) {
	fmt.Print( treeNode.Val)
	for _, child := range treeNode.Children {

		PrintTree(child, space+"_")
	}

}

func parseExpression(tokens []lexer.TokenType) *TreeNode {
	fmt.Println("parseExpression", tokens)
	for i := 0; i < len(tokens); i++ {
		// fmt.Print(tokens[i], " ")
	}
	// fmt.Println()
	return parseEquality(tokens)
}
func parseEquality(tokens []lexer.TokenType) *TreeNode {
	fmt.Println("equality", tokens)
	//get index of equality operator
	eqIndex := -1
	for i := 0; i < len(tokens); i++ {
		if tokens[i].Ref == "=" {
			eqIndex = i
			break
		}
	}
	//varname is tokens[eqIndex-1]
	right := parseComparison(tokens[eqIndex+1:])
	node := TreeNode{"=", []*TreeNode{&TreeNode{tokens[eqIndex-1].Ref, nil}, right}}
	return &node
}
func parseComparison(tokens []lexer.TokenType) *TreeNode {
	fmt.Println("comparison", tokens)
	//find index of first comp operator < > <= >=
	compIndex := -1
	compOp := ""
	for i := 0; i < len(tokens); i++ {
		if tokens[i].Ref == "==" || tokens[i].Ref == "!=" {
			compIndex = i
			compOp = tokens[i].Ref
			break
		}
	}
	if compIndex > 0 {

		left := parseTerm(tokens[:compIndex])
		right := parseComparison(tokens[compIndex+1:])
		node := TreeNode{compOp, []*TreeNode{left, right}}
		return &node

	} else {

		right := parseTerm(tokens[compIndex+1:])
		node := TreeNode{compOp, []*TreeNode{right}}
		return &node
	}

}

func parseTerm(tokens []lexer.TokenType) *TreeNode {
	fmt.Println("term", tokens)
	opIndex := -1
	op := ""
	for i := 0; i < len(tokens); i++ {
		if tokens[i].Ref == "+" || tokens[i].Ref == "-" {
			opIndex = i
			op = tokens[i].Ref
			break
		}
	}
	if opIndex > 0 {
		left := parseFactor(tokens[:opIndex])
		right := parseTerm(tokens[opIndex+1:])
		node := TreeNode{op, []*TreeNode{left, right}}
		return &node

	} else {
		right := parseFactor(tokens[opIndex+1:])
		node := TreeNode{op, []*TreeNode{right}}
		return &node
	}
}

func parseFactor(tokens []lexer.TokenType) *TreeNode {
	fmt.Println("factor", tokens)
	opIndex := -1
	op := ""
	for i := 0; i < len(tokens); i++ {
		if tokens[i].Ref == "*" || tokens[i].Ref == "/" {
			opIndex = i
			op = tokens[i].Ref
			break
		}
	}
	if opIndex > 0 {
		left := parseUnary(tokens[:opIndex])
		right := parseFactor(tokens[opIndex+1:])
		node := TreeNode{op, []*TreeNode{left,right}}
		return &node
	} else {

		right := parseUnary(tokens[opIndex+1:])
		node := TreeNode{op, []*TreeNode{right}}
		return &node
	}
}

func parseUnary(tokens []lexer.TokenType) *TreeNode {
	opIndex := -1
	op := ""
	for i := 0; i < len(tokens); i++ {
		if tokens[i].Ref == "!" || tokens[i].Ref == "-" {
			opIndex = i
			op = tokens[i].Ref
			break
		}
	}
	if op == "" {
		node := TreeNode{op, []*TreeNode{parsePrimary(tokens[opIndex+1:])}}
		return &node
	} else {
		node := TreeNode{op, []*TreeNode{parseUnary(tokens[opIndex+1:])}}
		return &node
	}
}

func parsePrimary(tokens []lexer.TokenType) *TreeNode {
	// fmt.Println("prim",len(tokens))
	return &TreeNode{tokens[0].Ref, nil}
}
