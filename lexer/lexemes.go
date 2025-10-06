package lexer

import "he++/utils"

// lexeme types
var PUNCTUATION = LexerTokenType("punctuation")
var KEYWORD = LexerTokenType("keyword")
var BRACKET = LexerTokenType("bracket")
var OPERATOR = LexerTokenType("operator")
var IDENTIFIER = LexerTokenType("identifier")
var INTEGER = LexerTokenType("int")
var FLOATINGPT = LexerTokenType("floatingpt")
var STRING_LITERAL = LexerTokenType("string_literal")
var BOOLEAN_LITERAL = LexerTokenType("boolean_literal")

// keywords
var IF = "si"
var THEN = "entonces"
var ELSE_IF = "elseif" // note: may not be used
var ELSE = "o"
var LET = "definir"
var INT = "int"
var FLOAT = "float"
var BOOLEAN = "bool"
var CHAR = "char"
var STRING = "string"
var FOR = "para"
var WHILE = "mientras"
var THAT = "que"
var BREAK = "interrumpir"
var CONTINUE = "seguir"
var RETURN = "devolver"
var FUNCTION = "funcion"
var STRUCT = "estructura"
var TRUE = "verdad"
var FALSE = "falso"
var VOID = "vacio"

// symbols
var LPAREN = "{"
var RPAREN = "}"
var COLON = ":"
var SEMICOLON = ";"
var DOT = "."
var OPEN_PAREN = "("
var CLOSE_PAREN = ")"
var COMMA = ","
var EQUALS = "="
var OPEN_SQUARE = "["
var CLOSE_SQUARE = "]"
var ANGLE_START = "<"
var ANGLE_END = ">"

// arithmetic operators
var ADD = "+"
var SUB = "-"
var MUL = "*"
var DIV = "/"
var MODULO = "%"

// bitwise operators
var NOT = "!"
var PIPE = "|"
var AMP = "&"
var LSHIFT = "<<" // todo: not recognized by lexer
var RSHIFT = ">>" // not recognized by lexer

// logical operators
var LESS = "<"
var GREATER = ">"
var EQ = "=="
var NEQ = "!="
var LEQ = "<="
var GEQ = ">="
var ANDAND = "&&"
var OROR = "||"

// unary operators
var INC = "++"
var DEC = "--"
var ASSN = "="
var HASHTAG = "#" // not used anywhere yet

var TERN_IF = "?"

var MATH_COMMA = '_'
var MATH_DOT = '.'

var Keywords = map[string]bool{
	IF:       true,
	ELSE_IF:  true,
	ELSE:     true,
	LET:      true,
	TRUE:     true,
	FALSE:    true,
	STRING:   true,
	FOR:      true,
	WHILE:    true,
	THAT:     true,
	BREAK:    true,
	CONTINUE: true,
	RETURN:   true,
	FUNCTION: true,
	STRUCT:   true,
	VOID:     true,
}

var Operators = map[string]bool{
	ADD:     true,
	SUB:     true,
	MUL:     true,
	DIV:     true,
	MODULO:  true,
	LESS:    true,
	GREATER: true,
	NOT:     true,
	PIPE:    true,
	AMP:     true,
	EQ:      true,
	NEQ:     true,
	LEQ:     true,
	GEQ:     true,
	INC:     true,
	DEC:     true,
	ANDAND:  true,
	OROR:    true,
	ASSN:    true,
	HASHTAG: true,
	DOT:     true,
	TERN_IF: true,
	COMMA:   true,
}

var names = map[string]string{
	IF:           "if",
	ELSE_IF:      "else_if",
	ELSE:         "else",
	LET:          "let",
	INT:          "int",
	FLOAT:        "float",
	BOOLEAN:      "boolean",
	STRING:       "string",
	FOR:          "for",
	WHILE:        "while",
	BREAK:        "break",
	CONTINUE:     "continue",
	RETURN:       "return",
	FUNCTION:     "function",
	STRUCT:       "struct",
	TRUE:         "true",
	FALSE:        "false",
	LPAREN:       "LPAREN",
	RPAREN:       "RPAREN",
	COLON:        "colon",
	SEMICOLON:    "semicolon",
	DOT:          "dot",
	OPEN_PAREN:   "open_paren",
	CLOSE_PAREN:  "close_paren",
	COMMA:        "comma",
	EQUALS:       "equals",
	OPEN_SQUARE:  "open_square",
	CLOSE_SQUARE: "close_square",
	ANGLE_START:  "angle_start",
	ANGLE_END:    "angle_end",
	ADD:          "add",
	SUB:          "sub",
	MUL:          "mul",
	DIV:          "div",
	MODULO:       "mod",
	NOT:          "not",
	PIPE:         "pipe",
	AMP:          "amp",
	LSHIFT:       "lshift",
	RSHIFT:       "rshift",
	LESS:         "less",
	GREATER:      "greater",
	EQ:           "equal",
	NEQ:          "not_equal",
	INC:          "inc",
	DEC:          "dec",
}

var OpTrie = constructTrie(Operators)
var KwTrie = constructTrie(Keywords)

var decimalDigits = map[byte]bool{
	'0': true, '1': true, '2': true, '3': true, '4': true,
	'5': true, '6': true, '7': true, '8': true, '9': true,
}

var binaryDigits = map[byte]bool{
	'0': true, '1': true,
}

var octalDigits = map[byte]bool{
	'0': true, '1': true, '2': true, '3': true, '4': true,
	'5': true, '6': true, '7': true,
}

var hexDigits = map[byte]bool{
	'0': true, '1': true, '2': true, '3': true, '4': true,
	'5': true, '6': true, '7': true, '8': true, '9': true,
	'a': true, 'b': true, 'c': true, 'd': true, 'e': true, 'f': true,
	'A': true, 'B': true, 'C': true, 'D': true, 'E': true, 'F': true,
}

func isDelimiter(c byte) bool {
	return c == ' ' || c == '\n' || c == '\t' || c == '\r'
}

func isQuote(c byte) bool {
	return c == '"' || c == '`'
}

func isBracket(c string) bool {
	return c == LPAREN || c == RPAREN || c == OPEN_PAREN || c == CLOSE_PAREN || c == OPEN_SQUARE || c == CLOSE_SQUARE || c == ANGLE_START || c == ANGLE_END
}

func isPunctuation(c string) bool {
	return c == SEMICOLON
}

func isKeyword(c string) bool {
	return KwTrie.Search(c)
}

func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func isNumberPart(c byte) bool {
	return isDigit(c) || c == byte(MATH_DOT) || c == byte(MATH_COMMA)
}

func constructTrie(elems map[string]bool) utils.Trie {
	trie := utils.MakeTrie()
	for i := range elems {
		trie.Insert(i)
	}
	return trie
}
