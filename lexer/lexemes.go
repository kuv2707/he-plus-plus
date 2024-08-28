package lexer

// lexeme types
var PUNCTUATION = LexerTokenType("punctuation")
var KEYWORD = LexerTokenType("keyword")
var BRACKET = LexerTokenType("bracket")
var OPERATOR = LexerTokenType("operator")
var IDENTIFIER = LexerTokenType("identifier")
var VALUE = LexerTokenType("value")
var INTEGER = LexerTokenType("int")
var FLOATINGPT = LexerTokenType("floatingpt")
var STRING_LITERAL = LexerTokenType("string_literal")
var BOOLEAN_LITERAL = LexerTokenType("boolean_literal")

// keywords
var IF = "if"
var ELSE_IF = "elseif"
var ELSE = "else"
var LET = "let"
var INT = "int"
var FLOAT = "float"
var BOOLEAN = "bool"
var STRING = "string"
var LOOP = "loop"
var BREAK = "break"
var CONTINUE = "continue"
var RETURN = "return"
var FUNCTION = "function"
var STRUCT = "struct"
var TRUE = "true"
var FALSE = "false"

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
var LSHIFT = "<<"
var RSHIFT = ">>"

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

var Keywords = []string{IF, ELSE_IF, ELSE, LET, INT, FLOAT, BOOLEAN, TRUE, FALSE, STRING, LOOP, BREAK, CONTINUE, RETURN, FUNCTION, STRUCT}

var Operators = []string{ADD, SUB, MUL, DIV, MODULO, LESS, GREATER, NOT, PIPE, AMP, EQ, NEQ, LEQ, GEQ, INC, DEC, ANDAND, OROR, ASSN, HASHTAG, AMP, DOT}

var names = map[string]string{
	IF:           "if",
	ELSE_IF:      "else_if",
	ELSE:         "else",
	LET:          "let",
	INT:          "int",
	FLOAT:        "float",
	BOOLEAN:      "boolean",
	STRING:       "string",
	LOOP:         "loop",
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

func Contains(slice []string, item string) bool {
	for _, element := range slice {
		if element == item {
			return true
		}
	}
	return false
}

func isOperator(c string) bool {
	return Contains(Operators, c)
}

func isDelimiter(c string) bool {
	return c == " " || c == "\n" || c == "\t" || c == "\r"
}

func isQuote(c string) bool {
	return c == "\"" || c == "`"
}

func isBracket(c string) bool {
	return c == LPAREN || c == RPAREN || c == OPEN_PAREN || c == CLOSE_PAREN || c == OPEN_SQUARE || c == CLOSE_SQUARE || c == ANGLE_START || c == ANGLE_END
}

func isPunctuation(c string) bool {
	return c == COMMA || c == SEMICOLON || c == COLON
}

func isKeyword(c string) bool {
	return Contains(Keywords, c)
}

func isDigit(c string) bool {
	return c >= "0" && c <= "9"
}

func getTokenName(c string) string {
	return names[c]
}
