package lexer

type TokenType struct {
	Type string
	Ref  string //like int,float etc for DATATYPE
}

type Node struct {
	Val  TokenType
	Next *Node
}
