package parser

import (
	"fmt"
	_"toylingo/globals"
	"toylingo/lexer"
	"toylingo/utils"
	_"sort"
)

var tokensArr = make([]lexer.TokenType, 0)


func makeTreeNode(label string, children []*TreeNode, description string) *TreeNode {
	
	return &TreeNode{label, description, children, make(map[string]*TreeNode)}
}

func ParseTree(tokens *lexer.Node) *TreeNode {

	for node := tokens; node != nil; node = node.Next {
		tokensArr = append(tokensArr, node.Val)
	}
	length = len(tokensArr)
	return parseScope()
}

func parseScope() *TreeNode {
	scopeNode := makeTreeNode("scope", make([]*TreeNode, 0), "scope")
	OUT:
	for index() < maxIndex() {
		token := next()
		switch token.Type {
		case "LET":
			scopeNode.Children = append(scopeNode.Children, parseLet())

		case "IF":
			scopeNode.Children = append(scopeNode.Children,parseConditionalBlock())

		case "LOOP":
			scopeNode.Children = append(scopeNode.Children, parseLoop())

		case "FUNCTION":
			scopeNode.Children = append(scopeNode.Children, parseFunction())

		case "SCOPE_END":
			break OUT


		default:
			prev()
			scopeNode.Children= append(scopeNode.Children,parseExpression(collectTill("SEMICOLON")))
		}

	}
	//sort children such that all function definitions are at the beginning
	fmt.Println("Fdafa",len(scopeNode.Children))
	// sort.SliceStable(scopeNode.Children,func(i,j int) bool{
	// 	fmt.Println("sorting",scopeNode.Children[i].Label,scopeNode.Children[j].Label)
		
	// })

	return scopeNode
}

func parseLet() *TreeNode {
	tokens := collectTill("SEMICOLON")
	fmt.Println("let tokens:", tokens)
	expNode := parseExpression(tokens)
	

	return expNode
}



func parseConditionalBlock() *TreeNode {
	condBlock:=makeTreeNode("conditional_block",nil,"cond")
	condNode := parseExpression(collectTill("SCOPE_START"))
	condBlock.Properties["condition0"]=condNode
	condBlock.Properties["ifnode0"]=parseScope()
	count:=1
	for matchCurrent("ELSE IF") {
		next()
		condNode := parseExpression(collectTill("SCOPE_START"))
		condBlock.Properties["condition"+fmt.Sprint(count)]=condNode
		condBlock.Properties["ifnode"+fmt.Sprint(count)]=parseScope()
		count++
	}
	if matchCurrent("ELSE") {
		next()
		condBlock.Properties["else"]=parseScope()
	}
	return condBlock
}

func parseLoop() *TreeNode {
	loopNode:=makeTreeNode("loop",nil,"loop")
	loopNode.Properties["condition"]=parseExpression(collectTill("SCOPE_START"))
	loopNode.Properties["body"]=parseScope()
	return loopNode
}

func parseFunction() *TreeNode {
	funcNode:=makeTreeNode("function",nil,"func")
	expect("OPEN_PAREN")
	next()
	funcNode.Properties["args"]=parseArgs(collectTill("CLOSE_PAREN"))
	funcNode.Properties["body"]=parseScope()
	
	return funcNode
}

func parseArgs(tokens []lexer.TokenType) *TreeNode {
	argsNode:=makeTreeNode("args",nil,"args")
	for i:=0;i<len(tokens);i++ {
		if tokens[i].Type=="COMMA" {
			continue
		}
		argsNode.Children=append(argsNode.Children,makeTreeNode("arg",nil,tokens[i].Ref))
	}
	return argsNode
}

func parseExpression(tokens []lexer.TokenType) *TreeNode {
	precedence := [][]string{
		{"="},
		{"||", "&&"},
		{"==", "!=", "<", ">", "<=", ">="},
		{"+", "-"},
		{"*", "/"},
		{"!", "-", "#"},
	}
	utils.DoNothing(precedence)
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

	} else if opIndex == 0 {
		right := parseFactor(tokens[opIndex:])
		return right
	} else {
		return parseFactor(tokens[opIndex+1:])
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
	op:=""
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
	if opIndex == -1 {
		node := parsePrimary(tokens)
		return node
	} else {
		node := makeTreeNode("operator", []*TreeNode{parseUnary(tokens[opIndex+1:])}, op)
		return node
	}
}

func parsePrimary(tokens []lexer.TokenType) *TreeNode {
	if tokens[0].Type == "OPEN_PAREN" {
		return parseExpression(tokens[1 : len(tokens)-1])
	} else {
		primNode := makeTreeNode("primary", nil, tokens[0].Ref)
		if len(tokens) > 1 {
			//might be func call
			if tokens[1].Type != "OPEN_PAREN" {
				fmt.Println("invalid expression")
				printTokensArr(tokens)
				panic("invalid expression")
			}
			node := makeTreeNode("func_call", make([]*TreeNode, 0), tokens[0].Ref)
			primNode.Children = append(primNode.Children, node)
			//parse args
			balance := 1
			last := 2
			for k := 2; k < len(tokens); k++ {
				if tokens[k].Type == "OPEN_PAREN" {
					balance++
				} else if tokens[k].Type == "CLOSE_PAREN" {
					balance--
				}

				if tokens[k].Type == "COMMA" && balance == 1 {
					node.Properties["args"+fmt.Sprint(len(node.Properties))] = parseExpression(tokens[last:k])
					last = k + 1
				}
				if balance == 0 {
					node.Properties["args"+fmt.Sprint(len(node.Properties))] = parseExpression(tokens[last:k])
					break
				}
			}
		}
		return primNode

	}
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
