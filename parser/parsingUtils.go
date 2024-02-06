package parser

import (
	"fmt"
	"he++/lexer"
	"he++/utils"
	"os"
	"regexp"
)

var i, length int = 0, 0

func next() lexer.TokenType {
	i++
	return tokensArr[i-1]
}
func prev() lexer.TokenType {
	i--
	return tokensArr[i+1]
}

func peek() lexer.TokenType {
	//todo check index before accessing
	return tokensArr[i+1]

}

func index() int {
	return i
}

func maxIndex() int {
	return length
}

func expect(tokenType string) {
	if tokensArr[i].Type != tokenType {
		panic("unexpected token " + tokensArr[i].Ref + " " + tokensArr[i].Ref + " expected " + tokenType)
	}
}

// the opening bracket should be included in the tokens passed, at the first position
func collectTillBalanced(close string, tokens []lexer.TokenType) ([]lexer.TokenType, int) {
	open := tokens[0].Ref
	balance := 1
	for i := 1; i < len(tokens); i++ {
		if tokens[i].Ref == open {
			balance++
		} else if tokens[i].Ref == close {
			balance--
		}
		if balance == 0 {
			return tokens[1:i], i
		}
	}
	if balance != 0 {
		abort(tokens[0].LineNo, "unbalanced parentheses", open)
	}
	return []lexer.TokenType{}, len(tokens)-1
}

func splitTokens(tokens []lexer.TokenType, separator string) [][]lexer.TokenType {
	tokensArr := make([][]lexer.TokenType, 0)
	start := 0
	for i := 0; i < len(tokens); i++ {
		if tokens[i].Ref == separator {
			tokensArr = append(tokensArr, tokens[start:i])
			start = i + 1
		}
	}
	tokensArr = append(tokensArr, tokens[start:])
	return tokensArr
}

// collects tokens till a token of type tokenType is found and consumes but does not include it in returned array
func collectTill(tokenType string) []lexer.TokenType {
	tokens := make([]lexer.TokenType, 0)
	for ; i < len(tokensArr); i++ {
		if tokensArr[i].Type == tokenType {
			break
		}
		tokens = append(tokens, tokensArr[i])
	}
	i++
	return tokens
}

func collectTillIn(tokenType string, tokens []lexer.TokenType) ([]lexer.TokenType, int) {
	toks := make([]lexer.TokenType, 0)
	i := 0
	for ; i < len(tokens); i++ {
		if tokens[i].Type == tokenType {
			break
		}
		toks = append(toks, tokens[i])
	}
	i++
	return toks, i
}

func matchCurrent(tokenType string) bool {
	if index() >= maxIndex() {
		return false
	}
	if tokensArr[i].Type == tokenType {
		return true
	}
	return false
}

func consume(tokenType string) {
	if index() >= maxIndex() {
		return
	}
	if tokensArr[i].Type == tokenType {
		i++
	} else {
		panic("unexpected token" + tokensArr[i].Type + " " + tokensArr[i].Ref + " expected " + tokenType)
	}
}

func seekClosingParen(tokens []lexer.TokenType, bracket string) int {
	balance := 1
	for i := 0; i < len(tokens); i++ {
		if tokens[i].Ref == bracket {
			balance++
		} else if tokens[i].Ref == utils.ClosingBracket(bracket) {
			balance--
		}
		if balance == 0 {
			return i
		}
	}
	fmt.Println(tokens)
	panic("unbalanced parentheses")
}

func (treeNode *TreeNode) PrintTree(space string) {
	color := utils.GetRandomColor()
	fmt.Print(color)
	fmt.Println(space + "{")
	space += "  "
	fmt.Println(space + "label:" + treeNode.Label)
	fmt.Println(space + "desc:" + treeNode.Description)
	for key, val := range treeNode.Properties {
		fmt.Println(space + key + ":")
		val.PrintTree(space + utils.ONETAB)
	}
	if len(treeNode.Children) > 0 {
		fmt.Println(space + "children:\n" + space + "[")
		for _, child := range treeNode.Children {

			child.PrintTree(space + utils.ONETAB)
			fmt.Print(utils.Colors["RESET"])
		}
		fmt.Print(color)
		fmt.Println(space + "]")
	}
	fmt.Println(space[0:len(space)-2] + "}")

}
func printTokensArr(tokens []lexer.TokenType) {
	for i := 0; i < len(tokens); i++ {
		fmt.Print(tokens[i].Ref, " ")
	}
	fmt.Println()
}

func isBalancedExpression(tokens []lexer.TokenType) bool {
	balance := 0
	for i := 0; i < len(tokens); i++ {
		if tokens[i].Ref == "(" {
			balance++
		} else if tokens[i].Ref == ")" {
			balance--
		}
	}
	return balance == 0
}

func abort(lineNo int, k ...interface{}) {
	fmt.Print(utils.Colors["RED"])
	fmt.Print("syntax error at line", fmt.Sprint(lineNo), ": ")
	fmt.Println(k...)
	fmt.Print(utils.Colors["RESET"])
	fmt.Print(utils.Colors["BOLDRED"])
	fmt.Println("execution interrupted")
	fmt.Print(utils.Colors["RESET"])
	os.Exit(1)
}

func isValidVariableName(variableName string) bool {
	regexPattern := `^[a-zA-Z_][a-zA-Z0-9_]*$`
	regExp, err := regexp.Compile(regexPattern)
	if err != nil {
		return false
	}
	return regExp.MatchString(variableName)
}
