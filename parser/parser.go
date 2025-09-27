package parser

import (
	"he++/lexer"
	nodes "he++/parser/node_types"
)

type Parser struct {
	tokenStream      *TokenStream
	prefixParselets  map[string]func(*Parser) nodes.TreeNode
	postfixParselets map[string]func(*Parser, nodes.TreeNode) nodes.TreeNode

	scopeParselets map[string]func(*Parser) nodes.TreeNode
}

func NewParser(tokens []lexer.LexerToken) *Parser {
	ts := NewTokenStream(tokens)
	p := &Parser{ts,
		make(map[string]func(*Parser) nodes.TreeNode),
		make(map[string]func(*Parser, nodes.TreeNode) nodes.TreeNode),
		make(map[string]func(*Parser) nodes.TreeNode),
	}
	p.initParselets()
	return p
}

func (p *Parser) initParselets() {
	p.prefixParselets[lexer.FLOATINGPT.String()] = parseFloat
	p.prefixParselets[lexer.STRING_LITERAL.String()] = parseString
	p.prefixParselets[lexer.INTEGER.String()] = parseInteger
	p.prefixParselets[lexer.IDENTIFIER.String()] = parseIdentifier
	p.prefixParselets[lexer.TRUE] = parseBoolean
	p.prefixParselets[lexer.FALSE] = parseBoolean
	p.prefixParselets[lexer.DEC] = parsePrefixOperator
	p.prefixParselets[lexer.INC] = parsePrefixOperator
	p.prefixParselets[lexer.AMP] = parsePrefixOperator
	p.prefixParselets[lexer.MUL] = parsePrefixOperator
	p.prefixParselets[lexer.OPEN_PAREN] = parseBracketExpression
	p.prefixParselets[lexer.OPEN_SQUARE] = parseArrayDeclaration
	p.prefixParselets[lexer.LPAREN] = parseStructValue

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
	p.scopeParselets[lexer.STRUCT] = parseStructDefn

}

func (p *Parser) getPrefixParselet(tok lexer.LexerToken) (func(*Parser) nodes.TreeNode, bool) {
	prefix, exists := p.prefixParselets[tok.Text()]
	if !exists {
		prefix, exists = p.prefixParselets[tok.Type().String()]
	}
	return prefix, exists
}
