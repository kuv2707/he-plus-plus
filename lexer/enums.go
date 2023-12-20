package lexer

type TokenType struct {
	Type   string
	Ref    string
	LineNo int
}

type Node struct {
	Val  TokenType
	Next *Node
}
