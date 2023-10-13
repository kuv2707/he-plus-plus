package lexer

type TokenType struct{
	Type string
	Ref string //like int,float etc for DATATYPE
}

var dict=map[string]string{
	"SCOPE_START":"{",
	"SCOPE_END":"}",
	"COLON":":",
	"SEMICOLON":";",
	"LET":"let",
	"INTEGER":"int",
	"FLOAT":"float",
	"BOOLEAN":"bool",
	"STRING":"string",
	"DOT":".",
	"OPEN_PAREN":"(",
	"CLOSE_PAREN":")",
}