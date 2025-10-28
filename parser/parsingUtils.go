package parser

import (
	"fmt"
	"he++/utils"
	"he++/lexer"
)

func debug(a ...any) {
	fmt.Println(a...)
}

func isPostfixOperator(op string) bool {
	return op == lexer.INC || op == lexer.DEC || op == lexer.OPEN_PAREN || op == lexer.OPEN_SQUARE
}

func parsingError(msg string, lineNo int) {
	panic(utils.Red(fmt.Sprintf("Parsing error at line %d: %s\n", lineNo, msg)))
}

func Contains(arr []interface{}, e interface{}) bool {
	for i := range arr {
		if arr[i] == e {
			return true
		}
	}
	return false
}
