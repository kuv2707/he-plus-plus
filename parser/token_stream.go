package parser

import (
	"fmt"
	lexer "he++/lexer"
)

type TokenStream struct {
	tokens []lexer.LexerToken
	i      int
}

func NewTokenStream(tokens []lexer.LexerToken) *TokenStream {
	return &TokenStream{tokens, 0}
}

func (ts *TokenStream) Current() lexer.LexerToken {
	return ts.tokens[ts.i]
}

func (ts *TokenStream) HasNext() bool {
	return ts.i < len(ts.tokens)
}

func (ts *TokenStream) Consume() lexer.LexerToken {
	ts.i++
	return ts.tokens[ts.i-1]
}

func (ts *TokenStream) ConsumeIf(t string) {
	if ts.HasNext() && ts.Current().Text() == t {
		ts.Consume()
		return
	}
	parsingError(fmt.Sprintf("Expected %s but got %s", t, ts.Current().Text()), ts.Current().LineNo())
}

func (ts *TokenStream) LookAhead(n int) *lexer.LexerToken {
	if ts.i+n >= len(ts.tokens) {
		return &lexer.LexerToken{}
	}
	return &ts.tokens[ts.i+n]
}
