package parser

import (
	"encoding/binary"
	"fmt"
	"he++/globals"
	_ "he++/globals"
	"he++/lexer"
	"he++/utils"
	"math"
	"sort"
)

var tokensArr = make([]lexer.TokenType, 0)

func makeTreeNode(label string, children []*TreeNode, description string, LineNo int) *TreeNode {

	return &TreeNode{label, description, children, make(map[string]*TreeNode), LineNo}
}

func ParseTree(tokens *lexer.Node) *TreeNode {

	for node := tokens; node != nil; node = node.Next {
		tokensArr = append(tokensArr, node.Val)
	}
	length = len(tokensArr)
	return parseScope()
}

func parseScope() *TreeNode {
	scopeNode := makeTreeNode("scope", make([]*TreeNode, 0), "scope", -1)
OUT:
	for index() < maxIndex() {
		token := next()
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
			scopeNode.Children = append(scopeNode.Children, makeTreeNode("return", []*TreeNode{parseExpression(collectTill("SEMICOLON"), 0)}, "return", token.LineNo))

		case "BREAK":
			val := []*TreeNode{}
			toks := collectTill("SEMICOLON")
			if len(toks) > 0 {
				val = []*TreeNode{parseExpression(toks, 0)}
			}
			scopeNode.Children = append(scopeNode.Children, makeTreeNode("break", val, "break", token.LineNo))

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
	condBlock := makeTreeNode("conditional_block", nil, "cond", -1)
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
	loopNode := makeTreeNode("loop", nil, "loop", -1)
	loopNode.Properties["condition"] = parseExpression(collectTill("SCOPE_START"), 0)
	loopNode.Properties["body"] = parseScope()
	return loopNode
}

func parseFunction() *TreeNode {
	funcNode := makeTreeNode("function", nil, "func", -1)
	expect("IDENTIFIER")
	funcNode.Description = tokensArr[i].Ref
	next()
	expect("OPEN_PAREN")
	next()
	funcNode.Properties["args"] = parseFormalArgs(collectTill("CLOSE_PAREN"))
	expect("SCOPE_START")
	next()
	funcNode.Properties["body"] = parseScope()

	return funcNode
}

func parseFormalArgs(tokens []lexer.TokenType) *TreeNode {
	argsNode := makeTreeNode("args", nil, "args", -1)
	for i := 0; i < len(tokens); i++ {
		if tokens[i].Type == "COMMA" {
			continue
		}
		argsNode.Children = append(argsNode.Children, makeTreeNode("arg", nil, tokens[i].Ref, tokens[i].LineNo))
	}
	return argsNode
}

func parseActualArgs(tokens []lexer.TokenType) *TreeNode {
	argsNode := makeTreeNode("args", nil, "args", -1)
	argToks := splitTokensBalanced(tokens, "COMMA")
	for i := 0; i < len(argToks); i++ {
		if len(argToks[i]) == 0 {
			continue
		}
		argsNode.Children = append(argsNode.Children, parseExpression(argToks[i], 0))
	}
	return argsNode
}

var precedence = [][]string{
	{"="},
	{"||", "&&"},
	{"==", "!=", "<", ">", "<=", ">="},
	{"+", "-"},
	{"*", "/"},
	{"!", "-", "#", "++", "--"},
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
			_, end := collectTillBalanced(utils.ClosingBracket(tokens[i].Ref), tokens[i:])
			i += end
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
		return makeTreeNode("operator", []*TreeNode{left, right}, op, tokens[opIndex].LineNo)
	} else {
		return parseExpression(tokens, rank+1)
	}
}

func parseUnary(tokens []lexer.TokenType, operators []string) *TreeNode {
	if utils.IsOneOf(tokens[0].Ref, operators) {
		return makeTreeNode("operator", []*TreeNode{parseUnary(tokens[1:], operators)}, tokens[0].Ref, tokens[0].LineNo)
	} else {
		return parsePrimary(tokens)
	}
}

func numberByteArray(value float64) []byte {
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, math.Float64bits(value))
	return bytes
}

func parsePrimary(tokens []lexer.TokenType) *TreeNode {
	if !isBalancedExpression(tokens) {
		abort(tokens[0].LineNo, "unbalanced expression")
	}

	if tokens[0].Type == "OPEN_PAREN" {
		return parseExpression(tokens[1:len(tokens)-1], 0)
	}
	if tokens[0].Type == "OPEN_SQUARE" {
		return parseArray(tokens)
	}

	switch tokens[0].Type {
	case "NUMBER":
		num := StringToNumber(tokens[0].Ref)
		globals.NumMap[tokens[0].Ref] = numberByteArray(num)
		return makeTreeNode("number", nil, tokens[0].Ref, tokens[0].LineNo)
	case "STRING":
		return makeTreeNode("string", nil, tokens[0].Ref, tokens[0].LineNo)
	case "BOOLEAN":
		return makeTreeNode("boolean", nil, tokens[0].Ref, tokens[0].LineNo)
	}

	if isValidVariableName(tokens[0].Ref) {
		name := tokens[0].Ref
		node := makeTreeNode("variable", nil, name, tokens[0].LineNo)
		if len(tokens) == 1 {
			return node
		}
		for i := 1; i < len(tokens); {
			if tokens[i].Ref == "[" {
				toks, end := collectTillBalanced(utils.ClosingBracket(tokens[i].Ref), tokens[i:])
				ind := parseExpression(toks, 0)
				index := makeTreeNode("index", nil, "index", tokens[i].LineNo)
				index.Properties["index"] = ind
				index.Children = append(index.Children, node)
				node = index
				i += end + 1
			} else if tokens[i].Ref == "(" {
				toks, end := collectTillBalanced(utils.ClosingBracket(tokens[i].Ref), tokens[i:])
				args := parseActualArgs(toks)
				call := makeTreeNode("call", nil, name, tokens[i].LineNo)
				call.Properties["args"] = args
				call.Children = append(call.Children, node)
				node = call
				i += end + 1
			} else {
				i++
			}
		}
		return node

	} else {
		abort(tokens[0].LineNo, "invalid variable name "+tokens[0].Ref)
	}

	return makeTreeNode("literal", nil, tokens[0].Ref, tokens[0].LineNo)
}

func parseArray(tokens []lexer.TokenType) *TreeNode {
	arrNode := makeTreeNode("array", nil, "array", tokens[0].LineNo)
	if len(tokens) == 2 {
		return arrNode
	}
	tokens = tokens[1 : len(tokens)-1]
	elems := splitTokensBalanced(tokens, "COMMA")
	for i := 0; i < len(elems); i++ {
		if len(elems[i]) == 0 {
			continue
		}
		arrNode.Children = append(arrNode.Children, parseExpression(elems[i], 0))
	}
	return arrNode
}
