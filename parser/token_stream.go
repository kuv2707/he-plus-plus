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

func (ts *TokenStream) CurrentIndex() int {
	return ts.i
}

func (ts *TokenStream) Current() lexer.LexerToken {
	return ts.tokens[ts.i]
}

func (ts *TokenStream) HasTokens() bool {
	return ts.i < len(ts.tokens)
}

func (ts *TokenStream) Consume() *lexer.LexerToken {
	ts.i++
	return &ts.tokens[ts.i-1]
}

func (ts *TokenStream) ConsumeOnlyIf(t string) *lexer.LexerToken {
	if ts.HasTokens() && ts.Current().Text() == t {
		return ts.Consume()
	}
	parsingError(fmt.Sprintf("Expected %s but got %s", t, ts.Current().Text()), ts.Current().LineNo())
	return nil
}

func (ts *TokenStream) ConsumeIf(t string) *lexer.LexerToken {
	if ts.HasTokens() && ts.Current().Text() == t {
		return ts.Consume()
	}
	return nil
}

func (ts *TokenStream) ConsumeIfType(t lexer.LexerTokenType) *lexer.LexerToken {
	if ts.HasTokens() && ts.Current().Type() == t {
		return ts.Consume()
	}
	parsingError(fmt.Sprintf("Expected %s but got %s", t, ts.Current().Type().String()), ts.Current().LineNo())
	return nil
}

func (ts *TokenStream) LookAhead(n int) *lexer.LexerToken {
	if ts.i+n >= len(ts.tokens) {
		return &lexer.LexerToken{}
	}
	return &ts.tokens[ts.i+n]
}

func (ts *TokenStream) Unread(n int) {
	if n < 0 {
		n = -n
	}
	ts.i -= n
	if ts.i < 0 {
		parsingError(fmt.Sprintf("Exhausted tokens while going back by %d", n), ts.Current().LineNo())
	}
}
