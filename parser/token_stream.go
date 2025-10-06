package parser

import (
	"fmt"
	lexer "he++/lexer"
)

type TokenStream struct {
	tokenChan  <-chan lexer.LexerToken
	currentTok *lexer.LexerToken
	nextTok    *lexer.LexerToken
	endOfToks  bool
}

func NewTokenStream(tokens <-chan lexer.LexerToken) *TokenStream {
	ts := &TokenStream{tokenChan: tokens, currentTok: nil, nextTok: nil, endOfToks: false}
	ts.Consume()
	ts.Consume()
	return ts
}

func (ts *TokenStream) Current() *lexer.LexerToken {
	return ts.currentTok
}

func (ts *TokenStream) HasTokens() bool {
	return ts.currentTok.Type() != ""
}

func (ts *TokenStream) Consume() *lexer.LexerToken {
	curr := ts.currentTok
	ts.currentTok = ts.nextTok
	tok, ok := <-ts.tokenChan
	if !ok {
		ts.endOfToks = true
	}
	ts.nextTok = &tok
	// fmt.Println("Consumed ", tok, ts.endOfToks)
	return curr
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

func (ts *TokenStream) ConsumeOnlyIfType(t lexer.LexerTokenType) *lexer.LexerToken {
	if ts.HasTokens() && ts.Current().Type() == t {
		return ts.Consume()
	}
	parsingError(fmt.Sprintf("Expected %s but got %s", t, ts.Current().Type().String()), ts.Current().LineNo())
	return nil
}

func (ts *TokenStream) LookOneAhead() *lexer.LexerToken {
	return ts.nextTok
}
