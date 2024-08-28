package lexer

import (
	"fmt"
	"he++/utils"
)

type LexerTokenType string

func (m LexerTokenType) String() string {
	return string(m)
}

type LexerToken struct {
	tokenType LexerTokenType
	ref       string
	lineNo    int
}

func NewLexerToken(tokenType LexerTokenType, ref string, lineNo int) LexerToken {
	return LexerToken{tokenType, ref, lineNo}
}

func (m LexerToken) String() string {
	return utils.Blue(string(m.tokenType)) + " " + utils.Yellow(m.ref) + " " + utils.Red(fmt.Sprint(m.lineNo))
}

func (l LexerToken) Text() string {
	return l.ref
}

func (l LexerToken) Type() LexerTokenType {
	return l.tokenType
}


