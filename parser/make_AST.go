package parser

import (
	"encoding/binary"
	"fmt"
	"he++/globals"
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

		case "STRUCT":
			scopeNode.Children = append(scopeNode.Children, parseStruct())

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
	fmt.Println(">", argToks)
	// return argsNode
	for i := 0; i < len(argToks); i++ {
		if len(argToks[i]) == 0 {
			continue
		}
		fmt.Println(">>>", argToks[i])
		argsNode.Children = append(argsNode.Children, parseExpression(argToks[i], 0))
	}
	if len(argToks) == 0 {
		fmt.Println("loser", argToks, i)
	}
	return argsNode
}

var precedence = [][]string{
	{"="},
	{"||", "&&"},
	{"==", "!=", "<", ">", "<=", ">="},
	{"+", "-"},
	{"*", "/"},
	{"."},
	{"!", "-", "#", "++", "--"}, //unary
}

func parseExpression(tokens []lexer.TokenType, rank int) *TreeNode {
	// fmt.Println("::",tokens)
	if rank != len(precedence)-1 {
		if rank == 5 {
		}
		return parseBinary(tokens, precedence[rank], rank)
	} else {
		return parseUnary(tokens, precedence[rank])
	}
}

func parseBinary(tokens []lexer.TokenType, operators []string, rank int) *TreeNode {
	opIndex := -1
	op := ""
	for i := len(tokens) - 1; i >= 0; i-- {
		if tokens[i].Ref == ")" || tokens[i].Ref == "]" || tokens[i].Ref == "}" {
			_, end := collectTillBalancedReverse(utils.OpeningBracket(tokens[i].Ref), tokens[0:i+1])
			i -= end - 1
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
	// fmt.Println("here", tokens)
	if !isBalancedExpression(tokens) {
		abort(tokens[0].LineNo, "unbalanced expression")
	}

	if tokens[0].Type == "OPEN_PAREN" {
		return parseExpression(tokens[1:len(tokens)-1], 0)
	}
	if tokens[0].Type == "OPEN_SQUARE" {
		return parseArray(tokens)
	}
	if tokens[0].Type == "SCOPE_START" {
		return parseObject(tokens)
	}

	dataval, valid := parseDataValue(tokens[0])
	if valid {
		return dataval
	}

	if isValidVariableName(tokens[0].Ref) {
		name := tokens[0].Ref
		node := makeTreeNode("variable", nil, name, tokens[0].LineNo)
		if len(tokens) == 1 {
			return node
		}
		for i := len(tokens) - 1; i > 0; {
			if tokens[i].Ref == "]" {
				toks, end := collectTillBalancedReverse(utils.OpeningBracket(tokens[i].Ref), tokens[0:i+1])
				ind := parseExpression(toks, 0)
				index := makeTreeNode("index", nil, "index", tokens[i].LineNo)
				index.Properties["index"] = ind
				index.Children = append(index.Children, node)
				node = index
				i -= end - 1
			} else if tokens[i].Ref == ")" {
				toks, end := collectTillBalancedReverse(utils.OpeningBracket(tokens[i].Ref), tokens[0:i+1])
				fmt.Println("call", toks, end)
				args := parseActualArgs(toks)
				call := makeTreeNode("call", nil, name, tokens[i].LineNo)
				call.Properties["args"] = args
				call.Children = append(call.Children, node)
				node = call
				if end > 1{
					i -= end - 1
				} else {
					// for empty arg list, end = 1
					i--
				}
			} else {
				i--
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

func parseObject(tokens []lexer.TokenType) *TreeNode {
	fmt.Println("objj",tokens)
	tokens = tokens[1 : len(tokens)-1]
	objNode := makeTreeNode("object", nil, "object", tokens[0].LineNo)
	kvps := splitTokensBalanced(tokens, "COMMA")
	for _, kvp := range kvps {
		objNode.Children = append(objNode.Children, parseKeyValuePair(kvp))
	}
	return objNode
}

func parseDataValue(token lexer.TokenType) (*TreeNode, bool) {
	switch token.Type {
	case "NUMBER":
		num := StringToNumber(token.Ref)
		globals.NumMap[token.Ref] = numberByteArray(num)
		return makeTreeNode("number", nil, token.Ref, token.LineNo), true
	case "STRING":
		return makeTreeNode("string", nil, token.Ref, token.LineNo), true
	case "BOOLEAN":
		return makeTreeNode("boolean", nil, token.Ref, token.LineNo), true
	}
	return nil, false
}

func parseKeyValuePair(tokens []lexer.TokenType) *TreeNode {
	fmt.Println("kvp", tokens)
	kvp := makeTreeNode("key_value", nil, "key_val", tokens[0].LineNo)
	globals.NumMap[tokens[0].Ref] = numberByteArray(float64(globals.HashString(tokens[0].Ref)))
	kvp.Properties["key"] = makeTreeNode("key", nil, tokens[0].Ref, tokens[0].LineNo)
	kvp.Properties["value"] = parseExpression(tokens[2:], 0)
	return kvp
}

func parseStruct() *TreeNode {
	if !globals.BeginsWithCapital(tokensArr[i].Ref) {
		abort(tokensArr[i].LineNo, "struct name must begin with a capital letter")
	}
	structNode := makeTreeNode("struct", nil, "struct", -1)
	expect("IDENTIFIER")
	structNode.Description = tokensArr[i].Ref
	next()
	expect("SCOPE_START")
	next()
	for matchCurrent("IDENTIFIER") {
		structNode.Children = append(structNode.Children, makeTreeNode("field", nil, tokensArr[i].Ref, tokensArr[i].LineNo))
		next()
		if matchCurrent("SCOPE_END") {
			break
		}
		consume("COMMA")
	}
	consume("SCOPE_END")
	return structNode
}
