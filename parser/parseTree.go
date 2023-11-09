package parser

import (
	"fmt"
	"toylingo/lexer"
	"toylingo/utils"
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
				is_else := tokensArr[i].Type == "ELSE"
				// fmt.Println("is else", is_else)
				if !(is_else) {
					for j := i + 1; j < len(tokensArr); j++ {
						if tokensArr[j].Type == "SCOPE_START" {
							condNode.Properties["condition"+fmt.Sprint(count)] = parseExpression(tokensArr[i+1 : j])
							i = j - 1
							count++
							break
						}
					}
				}
				i++
				scp := parseScope(tokensArr)
				// scp.PrintTree("/")
				condNode.Children = append(condNode.Children, scp)
				if is_else {
					// fmt.Println(tokensArr[i])
					break
				}
				i++
			}
			treeNode.Children = append(treeNode.Children, condNode)
		} else if tokensArr[i].Type == "LOOP" {
			loopnode := makeTreeNode("LOOP", make([]*TreeNode, 0), "loop")
			for j := i + 1; j < len(tokensArr); j++ {
				if tokensArr[j].Type == "SCOPE_START" {
					loopnode.Properties["condition"] = parseExpression(tokensArr[i+1 : j])
					i = j - 1
					break
				}
			}
			i++
			scp := parseScope(tokensArr)
			loopnode.Children = append(loopnode.Children, scp)
			treeNode.Children = append(treeNode.Children, loopnode)

		} else if tokensArr[i].Type == "FUNCTION" {
			funcNode := makeTreeNode("FUNCTION", make([]*TreeNode, 0), "function")
			count:=0
			i++
			for j := i + 1; j < len(tokensArr); j++ {
				if (tokensArr[j].Type == "COMMA" || tokensArr[j].Type == "CLOSE_PAREN"){
					funcNode.Properties["args"+fmt.Sprint(count)] = parseExpression(tokensArr[i+1 : j])
					i = j
					count++
				}
				if tokensArr[j].Type == "CLOSE_PAREN" {
					i = j+1
					break
				}
			}
			name:=tokensArr[i].Ref
			funcNode.Properties["name"] = makeTreeNode(name, nil, "varname")
			i++
			scp := parseScope(tokensArr)
			funcNode.Children = append(funcNode.Children, scp)
			treeNode.Children = append(treeNode.Children, funcNode)
		}else if tokensArr[i].Type == "BREAK" {
			treeNode.Children = append(treeNode.Children, makeTreeNode("BREAK", nil, "break"))
		} else if tokensArr[i].Type == "RETURN" {
			for j := i + 1; j < len(tokensArr); j++ {
				retNode := makeTreeNode("RETURN", make([]*TreeNode, 0), "return")
				if tokensArr[j].Type == "SEMICOLON" {
					retNode.Children = append(retNode.Children, parseExpression(tokensArr[i+1:j]))
					treeNode.Children = append(treeNode.Children,retNode )
					
					i = j
					break
				}

			}
		} else if tokensArr[i].Type == "SCOPE_START" {
			// fmt.Println("scope start")
			scope := parseScope(tokensArr)
			treeNode.Children = append(treeNode.Children, scope)
		} else if tokensArr[i].Type == "SCOPE_END" {
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
	scopeNode := makeTreeNode("scope", make([]*TreeNode, 0), "scope")
	i++
	ParseTree(tokensArr, scopeNode)
	return scopeNode
}

func (treeNode *TreeNode) PrintTree(space string) {
	color:=utils.GetRandomColor()
	fmt.Print(color)
	fmt.Println(space + "{")
	space += "  "
	fmt.Println(space + "desc:" + treeNode.Description)
	fmt.Println(space + "label:" + treeNode.Label)
	for key, val := range treeNode.Properties {
		fmt.Println(space + key + ":")
		val.PrintTree(space + utils.ONETAB)
	}
	fmt.Println(space + "children:\n"+space+"[")
	for _, child := range treeNode.Children {

		child.PrintTree(space + utils.ONETAB)
		fmt.Print(utils.Colors["RESET"])
	}
	fmt.Print(color)
	fmt.Println(space + "]\n" + space[0:len(space)-2] + "}")

}

func parseExpression(tokens []lexer.TokenType) *TreeNode {
	// printTokensArr(tokens)
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
	if eqIndex == -1 {
		return parseLogic(tokens)
	} else {
		right := parseEquality(tokens[eqIndex+1:])
		left := makeTreeNode(tokens[eqIndex-1].Ref, nil, "varname")
		node := makeTreeNode("operator", []*TreeNode{left, right}, "=")
		return node
	}
}

func parseLogic(tokens []lexer.TokenType) *TreeNode {
	// printTokensArr(tokens)
	// fmt.Println("logic", tokens)
	//get index of equality operator
	logicIndex := -1
	logicOp := ""
	for i := 0; i < len(tokens); i++ {
		if tokens[i].Ref == "(" {
			i += seekClosingParen(tokens[i+1:])
			continue
		}
		if tokens[i].Ref == "&&" || tokens[i].Ref == "||" {
			logicIndex = i
			logicOp = tokens[i].Ref
			break
		}
	}
	if logicIndex == -1 {
		return parseComparison(tokens)
	} else {
		right := parseLogic(tokens[logicIndex+1:])
		left := parseComparison(tokens[:logicIndex])
		node := makeTreeNode("operator", []*TreeNode{left, right}, logicOp)
		return node
	}
}

func parseComparison(tokens []lexer.TokenType) *TreeNode {
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
	// printTokensArr(tokens)
	if tokens[0].Type == "OPEN_PAREN" {
		return parseExpression(tokens[1 : len(tokens)-1])
	} else {
		primNode:=makeTreeNode("primary", nil, tokens[0].Ref)
		if len(tokens) > 1 {
			//might be func call
			if tokens[1].Type != "OPEN_PAREN" {
				panic("invalid expression")
			}
			node := makeTreeNode("func_call", make([]*TreeNode, 0), tokens[0].Ref)
			primNode.Children = append(primNode.Children, node)
			//parse args
			balance:=1;last:=2
			for k:=2;k<len(tokens);k++{
				if tokens[k].Type=="OPEN_PAREN"{
					balance++
				}else if tokens[k].Type=="CLOSE_PAREN"{
					balance--
				}

				if tokens[k].Type=="COMMA" && balance==1{
					node.Properties["args"+fmt.Sprint(len(node.Properties))] = parseExpression(tokens[last:k])
					last=k+1
				}
				if balance==0{
					node.Properties["args"+fmt.Sprint(len(node.Properties))]=parseExpression(tokens[last:k])
					break
				}
			}
		}
		return primNode

	}
}

func printTokensArr(tokens []lexer.TokenType) {
	for i := 0; i < len(tokens); i++ {
		fmt.Print(tokens[i].Ref, " ")
	}
	fmt.Println()
}

func seekClosingParen(tokens []lexer.TokenType) int {
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
