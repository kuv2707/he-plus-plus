package parser

import (
	"fmt"
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
	fmt.Printf("\033[31mParsing error at line %d: %s\033[0m\n", lineNo, msg)
	os.Exit(1)
}

func getPrecedence(op string) float32 {
	switch op {
	case lexer.DOT:
		return 3
	case lexer.OPEN_PAREN, lexer.OPEN_SQUARE:
		return 2.9
	case lexer.NOT:
		return 2.8
	case lexer.INC, lexer.DEC:
		return 2.8
	case lexer.DIV, lexer.MUL:
		return 2
	case lexer.ADD, lexer.SUB:
		return 1
	case lexer.EQ, lexer.NEQ, lexer.GREATER, lexer.LESS, lexer.LEQ, lexer.GEQ:
		return 0.5
	case lexer.ASSN:
		return 0.1

	}
	return 0
}
