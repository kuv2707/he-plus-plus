package parser

import (
	"he++/lexer"
)

type Parser struct {
	tokenStream      *TokenStream
	prefixParselets  map[string]func(*Parser) TreeNode
	postfixParselets map[string]func(*Parser, TreeNode) TreeNode

	scopeParselets map[string]func(*Parser) TreeNode
}

func NewParser(tokens []lexer.LexerToken) *Parser {
	ts := NewTokenStream(tokens)
	p := &Parser{ts,
		make(map[string]func(*Parser) TreeNode),
		make(map[string]func(*Parser, TreeNode) TreeNode),
		make(map[string]func(*Parser) TreeNode),
	}
	p.initParselets()
	return p
}

func (p *Parser) initParselets() {
	p.prefixParselets[lexer.FLOATINGPT.String()] = parseFloat
	p.prefixParselets[lexer.STRING_LITERAL.String()] = parseString
	p.prefixParselets[lexer.INTEGER.String()] = parseInteger
	p.prefixParselets[lexer.BOOLEAN_LITERAL.String()] = parseBoolean
	p.prefixParselets[lexer.IDENTIFIER.String()] = parseIdentifier
	p.prefixParselets[lexer.DEC] = parsePrefixOperator
	p.prefixParselets[lexer.INC] = parsePrefixOperator
	p.prefixParselets[lexer.OPEN_PAREN] = parseBracketExpression

	p.postfixParselets[lexer.OPEN_PAREN] = parseFuncCallArgs
	p.postfixParselets[lexer.OPEN_SQUARE] = parseArrayIndex
	p.postfixParselets[lexer.INC] = parsePostfixOperator
	p.postfixParselets[lexer.DEC] = parsePostfixOperator

	p.scopeParselets[lexer.FUNCTION] = parseFunction
	p.scopeParselets[lexer.LET] = parseVariableDeclaration
	p.scopeParselets[lexer.IF] = parseIfStatement
	p.scopeParselets[lexer.FOR] = parseLoopStatement
	p.scopeParselets[lexer.WHILE] = parseLoopStatement
	p.scopeParselets[lexer.RETURN] = parseReturnStatement

}

func (p *Parser) getPrefixParselet(tok lexer.LexerToken) (func(*Parser) TreeNode, bool) {
	prefix, exists := p.prefixParselets[tok.Text()]
	if !exists {
		prefix, exists = p.prefixParselets[tok.Type().String()]
	}
	return prefix, exists
}
