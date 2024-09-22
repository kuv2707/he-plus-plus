package parser

import (
	"fmt"
	"he++/globals"
	"he++/lexer"
	"os"
)

func debug(a ...any) {
	fmt.Println(a...)
}

func isPostfixOperator(op string) bool {
	return op == lexer.INC || op == lexer.DEC || op == lexer.OPEN_PAREN || op == lexer.OPEN_SQUARE
}

func parsingError(msg string, lineNo int) {
	fmt.Println(globals.Red(fmt.Sprintf("Parsing error at line %d: %s\n", lineNo, msg)))
	os.Exit(1)
}

func Contains(arr []interface{}, e interface{}) bool {
	for i := range arr {
		if arr[i] == e {
			return true
		}
	}
	return false
}
