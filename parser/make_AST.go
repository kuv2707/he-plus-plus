package parser

import (
	"fmt"
	"sort"
	_ "toylingo/globals"
	"toylingo/lexer"
	"toylingo/utils"
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
		fmt.Println("parsing",token)
		switch token.Type {
		case "LET":
			scopeNode.Children = append(scopeNode.Children, parseLet())

		case "IF":
			scopeNode.Children = append(scopeNode.Children, parseConditionalBlock())

		case "LOOP":
			scopeNode.Children = append(scopeNode.Children, parseLoop())

		case "FUNCTION":
			scopeNode.Children = append(scopeNode.Children, parseFunction())
		
		case "SCOPE_START":
			scopeNode.Children = append(scopeNode.Children, parseScope())

		case "SCOPE_END":
			break OUT
		
		case "RETURN":
			scopeNode.Children = append(scopeNode.Children, makeTreeNode("return", []*TreeNode{parseExpression(collectTill("SEMICOLON"), 0)}, "return"))

		default:
			prev()
			scopeNode.Children = append(scopeNode.Children, parseExpression(collectTill("SEMICOLON"), 0))
		}

	}
	//sort children such that all function definitions are at the beginning (hoisting)
	sort.SliceStable(scopeNode.Children, func(i, j int) bool {
		return scopeNode.Children[i].Label == "function"
	})

	return scopeNode
}

func parseLet() *TreeNode {
	tokens := collectTill("SEMICOLON")
	expNode := parseExpression(tokens, 0)

	return expNode
}

func parseConditionalBlock() *TreeNode {
	condBlock := makeTreeNode("conditional_block", nil, "cond")
	condNode := parseExpression(collectTill("SCOPE_START"), 0)
	condBlock.Properties["condition0"] = condNode
	condBlock.Properties["ifnode0"] = parseScope()
	count := 1
	for matchCurrent("ELSE IF") {
		next()
		condNode := parseExpression(collectTill("SCOPE_START"), 0)
		condBlock.Properties["condition"+fmt.Sprint(count)] = condNode
		condBlock.Properties["ifnode"+fmt.Sprint(count)] = parseScope()
		count++
	}
	if matchCurrent("ELSE") {
		next()
		consume("SCOPE_START")
		
		condBlock.Properties["else"] = parseScope()
	}
	return condBlock
}

func parseLoop() *TreeNode {
	loopNode := makeTreeNode("loop", nil, "loop")
	loopNode.Properties["condition"] = parseExpression(collectTill("SCOPE_START"), 0)
	loopNode.Properties["body"] = parseScope()
	return loopNode
}

func parseFunction() *TreeNode {
	funcNode := makeTreeNode("function", nil, "func")
	expect("OPEN_PAREN")
	next()
	funcNode.Properties["args"] = parseFormalArgs(collectTill("CLOSE_PAREN"))
	expect("IDENTIFIER")
	funcNode.Description = tokensArr[i].Ref
	next()
	expect("SCOPE_START")
	next()
	funcNode.Properties["body"] = parseScope()

	return funcNode
}

func parseFormalArgs(tokens []lexer.TokenType) *TreeNode {
	argsNode := makeTreeNode("args", nil, "args")
	for i := 0; i < len(tokens); i++ {
		if tokens[i].Type == "COMMA" {
			continue
		}
		argsNode.Children = append(argsNode.Children, makeTreeNode("arg", nil, tokens[i].Ref))
	}
	return argsNode
}

func parseActualArgs(tokens []lexer.TokenType) *TreeNode {
	argsNode := makeTreeNode("args", nil, "args")
	coll := make([]lexer.TokenType, 0)
	balance := 0
	for i := 0; i < len(tokens); i++ {
		if tokens[i].Type == "OPEN_PAREN" {
			balance++
		} else if tokens[i].Type == "CLOSE_PAREN" {
			balance--
		}

		if tokens[i].Type == "COMMA" && balance == 0 {
			argsNode.Children = append(argsNode.Children, parseExpression(coll, 0))
			coll = make([]lexer.TokenType, 0)
			continue
		}
		coll = append(coll, tokens[i])
	}
	if len(coll) > 0 {
		argsNode.Children = append(argsNode.Children, parseExpression(coll, 0))
	}
	return argsNode
}

var precedence = [][]string{
	{"="},
	{"||", "&&"},
	{"==", "!=", "<", ">", "<=", ">="},
	{"+", "-"},
	{"*", "/"},
	{"!", "-", "#"},
}

func parseExpression(tokens []lexer.TokenType, rank int) *TreeNode {
	if rank != len(precedence)-1 {
		return parseBinary(tokens, precedence[rank], rank)
	} else {
		return parseUnary(tokens, precedence[rank])
	}
}

func parseBinary(tokens []lexer.TokenType, operators []string, rank int) *TreeNode {
	opIndex := -1
	op := ""
	for i := 0; i < len(tokens); i++ {
		if tokens[i].Ref == "(" || tokens[i].Ref == "[" || tokens[i].Ref == "{" {
			i += seekClosingParen(tokens[i+1:], tokens[i].Ref)
			continue
		}
		if utils.IsOneOf(tokens[i].Ref, operators) {
			opIndex = i
			op = tokens[i].Ref
			break
		}
	}
	if opIndex > 0 {
		left := parseBinary(tokens[:opIndex], operators, rank)
		right := parseBinary(tokens[opIndex+1:], operators, rank)
		return makeTreeNode("operator", []*TreeNode{left, right}, op)
	} else {
		return parseExpression(tokens, rank+1)
	}
}

func parseUnary(tokens []lexer.TokenType, operators []string) *TreeNode {
	if utils.IsOneOf(tokens[0].Ref, operators) {
		return makeTreeNode("operator", []*TreeNode{parseUnary(tokens[1:], operators)}, tokens[0].Ref)
	} else {
		return parsePrimary(tokens)
	}
}

func parsePrimary(tokens []lexer.TokenType) *TreeNode {
	if tokens[0].Type == "OPEN_PAREN" {
		return parseExpression(tokens[1:len(tokens)-1], 0)
	}
	if tokens[0].Type == "OPEN_SQUARE" {
		return parseArray(tokens[1:len(tokens)-1])
	}
	if utils.IsLiteral(tokens[0].Type) {
		return makeTreeNode("literal", nil, tokens[0].Ref)
	}

	if !utils.IsValidVariableName(tokens[0].Ref) {
		panic("invalid variable name " + tokens[0].Ref)
	}
	primNode := makeTreeNode("primary", nil, tokens[0].Ref)
	if len(tokens) == 1 {
		return primNode
	}
	if !utils.IsOpenBracket(tokens[1].Ref) {
		printTokensArr(tokens)
		panic("invalid expression "+fmt.Sprint(tokens))
	}
	if tokens[1].Type == "OPEN_PAREN" {
		fmt.Println("call",tokens)
		node := makeTreeNode("call", make([]*TreeNode, 0), tokens[0].Ref)
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
				node.Properties["args"+fmt.Sprint(len(node.Properties))] = parseExpression(tokens[last:k], 0)
				last = k + 1
			}
			if balance == 0 {
				node.Properties["args"+fmt.Sprint(len(node.Properties))] = parseExpression(tokens[last:k], 0)
				break
			}
		}
		if balance != 0 {
			panic("syntax error in function call: unclosed parenthesis")
		}
		return node

	}
	if tokens[1].Type == "OPEN_SQUARE" {
		node:=makeTreeNode("index", []*TreeNode{parseExpression(tokens[2:len(tokens)-1], 0)}, tokens[0].Ref)
		primNode.Children = append(primNode.Children, node)
		if tokens[len(tokens)-1].Type!="CLOSE_SQUARE" {
			panic("syntax error in array index: unclosed square bracket")
		}
	}

	return primNode

}


func parseArray(tokens []lexer.TokenType) *TreeNode {
	arrNode := makeTreeNode("primary", nil, "array")
	fmt.Println("parsing array", tokens)
	balance := 0
	last := 0
	for k := 0; k < len(tokens); k++ {
		if utils.IsOpenBracket(tokens[k].Ref) {
			balance++
		} else if utils.IsCloseBracket(tokens[k].Ref) {
			balance--
		}

		if tokens[k].Type == "COMMA" && balance == 0 {
			ch:=parseExpression(tokens[last:k], 0)
			// ch.PrintTree("")
			arrNode.Children = append(arrNode.Children, ch)
			last = k + 1
		}
		
	}
	if last<len(tokens) {
		arrNode.Children = append(arrNode.Children, parseExpression(tokens[last:], 0))
	}

	return arrNode
}