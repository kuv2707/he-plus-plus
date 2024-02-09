package parser

import (
	"fmt"
	"he++/lexer"
	"he++/utils"
	"math"
	"os"
	"regexp"
	"strings"
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
	return []lexer.TokenType{}, len(tokens) - 1
}

func splitTokensBalanced(tokens []lexer.TokenType, separator string) [][]lexer.TokenType {
	tokensArr := make([][]lexer.TokenType, 0)
	start := 0
	balance := utils.MakeStack()
	for i := 0; i < len(tokens); i++ {
		if utils.IsOpenBracket(tokens[i].Ref) {
			balance.Push(utils.ClosingBracket(tokens[i].Ref))
		} else if tokens[i].Ref == balance.Peek() {
			balance.Pop()
		}
		if tokens[i].Type == separator && balance.IsEmpty() {
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



// todo: move to parsing phase for faster execution
func StringToNumber(str string) float64 {
	base := 10
	num := ""
	if len(str) < 2 {
		num = str
	} else {

		switch str[0:2] {
		case "0x":
			base = 16
			num = str[2:]
		case "0b":
			base = 2
			num = str[2:]
		case "0o":
			base = 8
			num = str[2:]
		default:
			num = str
		}
	}
	parsedNum := 0.0
	dotsep := strings.Split(num, ".")
	if len(dotsep) > 2 {
		panic("invalid number " + str)
	} else if len(dotsep) == 2 {
		parsedNum = parseNumber(dotsep[0], base) + parseFraction(dotsep[1], base)
	} else if len(dotsep) == 1 {
		parsedNum = parseNumber(dotsep[0], base)
	} else {
		panic("invalid number " + str)
	}
	return parsedNum
}

func parseNumber(num string, base int) float64 {
	parsedNum := 0.0
	l := len(num)
	for i := 0; i < l; i++ {
		parsedNum += float64(numVal(num[i])) * math.Pow(float64(base), float64(l-i-1))
	}
	return parsedNum
}
func parseFraction(num string, base int) float64 {
	parsedNum := 0.0
	l := len(num)
	for i := 0; i < l; i++ {
		parsedNum += float64(numVal(num[i])) * math.Pow(float64(base), float64(-(i+1)))
	}
	return parsedNum
}
func numVal(c byte) int {
	if c >= '0' && c <= '9' {
		return int(c - '0')
	}
	if c >= 'a' && c <= 'z' {
		return int(c - 'a' + 10)
	}
	if c >= 'A' && c <= 'Z' {
		return int(c - 'A' + 10)
	}
	panic("invalid number")
	return -1
}
