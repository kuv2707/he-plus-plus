package parser

import (
	"fmt"
	"toylingo/lexer"
	"toylingo/utils"
	_ "toylingo/utils"
)

type TreeNode struct {
	Label       string
	Description string
	Children    []*TreeNode
	Properties  map[string]*TreeNode
}

var KEYWORDS = []string{"IF", "ELSE IF", "ELSE", "FUNCTION", "SCOPE_END", "LET"}

func makeTreeNode(label string, children []*TreeNode, description string) *TreeNode {
	return &TreeNode{label, description, children, make(map[string]*TreeNode)}
}

func ParseTreeM(tokens *lexer.Node) *TreeNode {
	tokensArr := make([]lexer.TokenType, 0)
	for node := tokens; node != nil; node = node.Next {
		tokensArr = append(tokensArr, node.Val)
	}
	treeNode := makeTreeNode("root", make([]*TreeNode, 0), "root node")
	ParseTree(tokensArr, treeNode)
	return treeNode
}

var i int = 0

func ParseTree(tokensArr []lexer.TokenType, treeNode *TreeNode) {
	for ; i < len(tokensArr); i++ {

		if tokensArr[i].Type == "LET" {
			for j := i + 1; j < len(tokensArr); j++ {
				if tokensArr[j].Type == "SEMICOLON" {
					treeNode.Children = append(treeNode.Children, parseExpression(tokensArr[i+1:j]))
					i = j
					break
				}
			}
		} else if utils.IsOneOfArr(tokensArr[i].Type, []string{"IF", "ELSE IF", "ELSE"}) {
			condNode := makeTreeNode(tokensArr[i].Type, make([]*TreeNode, 0), "conditional_block")
			count := 0
			for (i < len(tokensArr)) && (tokensArr[i].Type == "IF" || tokensArr[i].Type == "ELSE IF" || tokensArr[i].Type == "ELSE") {
				is_else:=tokensArr[i].Type == "ELSE"
				if !(is_else) {
					for j := i + 1; j < len(tokensArr); j++ {
						if tokensArr[j].Type == "SCOPE_START" {
							condNode.Properties["condition"+fmt.Sprint(count)] = parseExpression(tokensArr[i+1 : j])
							i = j-1
							count++
							break
						}
					}
				}
				i++
				scp := parseScope(tokensArr)
				// scp.PrintTree("/")
				condNode.Children = append(condNode.Children, scp)
				if (is_else) {
					break
				}
				i++
			}
			treeNode.Children = append(treeNode.Children, condNode)
		} else if tokensArr[i].Type == "SCOPE_START" {
			// fmt.Println("scope start")
			scope := parseScope(tokensArr)
			treeNode.Children = append(treeNode.Children, scope)
		} else if tokensArr[i].Type == "SCOPE_END" {

			// fmt.Println("scope end")
			return
		} else {
			for j := i + 1; j < len(tokensArr); j++ {
				if tokensArr[j].Type == "SEMICOLON" {
					treeNode.Children = append(treeNode.Children, parseExpression(tokensArr[i:j]))
					i = j
					break
				}
			}
		}

	}
	return
}

func parseScope(tokensArr []lexer.TokenType) *TreeNode {
	// fmt.Println("scope start")
	// printTokensArr(tokensArr[i:])
	scopeNode := makeTreeNode("scope", make([]*TreeNode, 0), "scope")

	i++
	// fmt.Println("calling parsetree")
	ParseTree(tokensArr, scopeNode)
	// fmt.Println("back from parsetree")

	return scopeNode
}

func (treeNode *TreeNode) PrintTree(space string) {
	fmt.Println(space + "{")
	space += "  "
	fmt.Println(space + "desc:" + treeNode.Description)
	fmt.Println(space + "label:" + treeNode.Label)
	// if treeNode.Label == "IF" {
	// 	fmt.Println(space + "if block condition:")
	// 	treeNode.Properties["condition0"].PrintTree(space + "      ")
	// }
	fmt.Println(space + "children: [")
	for _, child := range treeNode.Children {

		child.PrintTree(space + "  ")
	}
	fmt.Println(space + "]\n" + space[0:len(space)-2] + "}")

}

func parseExpression(tokens []lexer.TokenType) *TreeNode {

	return parseEquality(tokens)
}
func parseEquality(tokens []lexer.TokenType) *TreeNode {
	// printTokensArr(tokens)
	// fmt.Println("equality", tokens)
	//get index of equality operator
	eqIndex := -1
	for i := 0; i < len(tokens); i++ {
		if tokens[i].Ref == "(" {
			i += seekClosingParen(tokens[i+1:])
			continue
		}
		if tokens[i].Ref == "=" {
			eqIndex = i
			break
		}
	}
	//varname is tokens[eqIndex-1]
	if eqIndex == -1 {
		return parseComparison(tokens)
	} else {
		right := parseEquality(tokens[eqIndex+1:])
		left := makeTreeNode(tokens[eqIndex-1].Ref, nil, "varname")
		node := makeTreeNode("operator", []*TreeNode{left, right}, "=")
		return node
	}
}
func parseComparison(tokens []lexer.TokenType) *TreeNode {
	// fmt.Println("comparison", tokens)
	//find index of first comp operator < > <= >=
	compIndex := -1
	compOp := ""
	for i := 0; i < len(tokens); i++ {
		if tokens[i].Ref == "(" {
			i += seekClosingParen(tokens[i+1:])
			continue
		}
		if tokens[i].Ref == "==" || tokens[i].Ref == "!=" || tokens[i].Ref == "<=" || tokens[i].Ref == ">=" || tokens[i].Ref == "<" || tokens[i].Ref == ">" {
			compIndex = i
			compOp = tokens[i].Ref
			break
		}
	}
	//it is malformed expression if compindex is =0
	if compIndex > 0 {

		left := parseTerm(tokens[:compIndex])
		right := parseComparison(tokens[compIndex+1:])
		node := makeTreeNode("operator", []*TreeNode{left, right}, compOp)
		return node

	} else {

		right := parseTerm(tokens[compIndex+1:])
		return right
	}

}

func parseTerm(tokens []lexer.TokenType) *TreeNode {
	// fmt.Println("term", tokens)
	opIndex := -1
	op := ""
	for i := 0; i < len(tokens); i++ {
		if tokens[i].Ref == "(" {
			i += seekClosingParen(tokens[i+1:])
			continue
		}
		if tokens[i].Ref == "+" || tokens[i].Ref == "-" {
			opIndex = i
			op = tokens[i].Ref
			break
		}
	}
	if opIndex > 0 {
		left := parseFactor(tokens[:opIndex])
		right := parseTerm(tokens[opIndex+1:])
		node := makeTreeNode("operator", []*TreeNode{left, right}, op)
		return node

	} else {
		right := parseFactor(tokens[opIndex+1:])
		return right
	}
}

func parseFactor(tokens []lexer.TokenType) *TreeNode {
	// fmt.Println("factor", tokens)
	opIndex := -1
	op := ""
	for i := 0; i < len(tokens); i++ {
		if tokens[i].Ref == "(" {
			i += seekClosingParen(tokens[i+1:])
			continue
		}
		if tokens[i].Ref == "*" || tokens[i].Ref == "/" {
			opIndex = i
			op = tokens[i].Ref
			break
		}
	}
	if opIndex > 0 {
		left := parseUnary(tokens[:opIndex])
		right := parseFactor(tokens[opIndex+1:])
		node := makeTreeNode("operator", []*TreeNode{left, right}, op)
		return node
	} else {

		right := parseUnary(tokens[opIndex+1:])
		return right
	}
}

func parseUnary(tokens []lexer.TokenType) *TreeNode {
	opIndex := -1
	op := ""
	for i := 0; i < len(tokens); i++ {
		if tokens[i].Ref == "(" {
			i += seekClosingParen(tokens[i+1:])
			continue
		}
		if tokens[i].Ref == "!" || tokens[i].Ref == "-" || tokens[i].Ref == "#" {
			opIndex = i
			op = tokens[i].Ref
			break
		}
	}
	if op == "" {
		node := parsePrimary(tokens)
		return node
	} else {
		node := makeTreeNode("operator", []*TreeNode{parseUnary(tokens[opIndex+1:])}, op)
		return node
	}
}

func parsePrimary(tokens []lexer.TokenType) *TreeNode {
	// fmt.Println("prim", (tokens))
	if tokens[0].Type == "OPEN_PAREN" {
		return parseExpression(tokens[1 : len(tokens)-1])
	} else {
		return makeTreeNode("primary", nil, tokens[0].Ref)

	}
}

func printTokensArr(tokens []lexer.TokenType) {
	for i := 0; i < len(tokens); i++ {
		fmt.Print(tokens[i].Ref, " ")
	}
	fmt.Println()
}

func seekClosingParen(tokens []lexer.TokenType) int {
	// fmt.Println("seeking for ",tokens)
	balance := 1
	for i := 0; i < len(tokens); i++ {
		if tokens[i].Ref == "(" {
			balance++
		} else if tokens[i].Ref == ")" {
			balance--
		}
		if balance == 0 {
			return i
		}
	}
	panic("unbalanced parentheses")
}
