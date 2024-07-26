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



func (m TokenType) String() string {
	return  "" + m.Ref
}