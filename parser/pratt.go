package parser

import (
	"he++/lexer"
)

// pratt parsing for expressions

func getPrecedence(op string) float32 {
	switch op {
	case lexer.DOT, lexer.AMP:
		return 3
	case lexer.OPEN_PAREN, lexer.OPEN_SQUARE:
		return 2.9
	case lexer.NOT:
		return 2.8
	case lexer.INC, lexer.DEC:
		return 2.8
	case lexer.MODULO:
		return 2.7
	case lexer.PIPE:
		return 2.6
	case lexer.DIV, lexer.MUL:
		return 2
	case lexer.ADD, lexer.SUB:
		return 1
	case lexer.EQ, lexer.NEQ, lexer.GREATER, lexer.LESS, lexer.LEQ, lexer.GEQ:
		return 0.5
	case lexer.ANDAND, lexer.OROR:
		return 0.4
	case lexer.TERN_IF:
		return 0.3
	case lexer.ASSN:
		return 0.1

	}
	return 0
}

func parseExpression(p *Parser, prec float32) TreeNode {
	t := p.tokenStream
	if !t.HasTokens() {
		parsingError("Unexpected end of expression after "+t.LookAhead(-1).Text(), 0)
		return nil
	}
	tok := t.Current()
	prefix, exists := p.getPrefixParselet(tok)
	if !exists {
		panic("Might not be an expression: " + tok.Type().String() + " " + tok.Text())
	}
	leftNode := prefix(p)
	for t.HasTokens() {
		opSymbol := t.Current().Text()
		if !(getPrecedence(opSymbol) > prec) {
			break
		}
		if isPostfixOperator(opSymbol) {
			// two operators in a row means this one is a postfix
			leftNode = p.postfixParselets[opSymbol](p, leftNode)
		} else {
			leftNode = parseInfixOperator(p, leftNode)
		}
	}
	return leftNode
}

func parseBracketExpression(p *Parser) TreeNode {
	p.tokenStream.Consume()
	expr := parseExpression(p, 0)
	if p.tokenStream.Current().Text() != lexer.CLOSE_PAREN {
		parsingError("Expected closing parenthesis", p.tokenStream.Current().LineNo())
	}
	p.tokenStream.Consume()
	return expr
}

func parseInteger(p *Parser) TreeNode {
	return NewNumberNode([]byte(p.tokenStream.Consume().Text()), "int")
}

func parseFloat(p *Parser) TreeNode {
	return NewNumberNode([]byte(p.tokenStream.Consume().Text()), "float")
}

func parseString(p *Parser) TreeNode {
	return NewStringNode([]byte(p.tokenStream.Consume().Text()))
}

func parseBoolean(p *Parser) TreeNode {
	tok := p.tokenStream.Consume()
	truth := tok.Text() == lexer.TRUE
	if truth {
		return NewBooleanNode([]byte(lexer.TRUE))
	}
	if tok.Text() != lexer.FALSE {
		parsingError("Expected boolean value", tok.LineNo())
	}
	return NewBooleanNode([]byte(lexer.FALSE))
}

func parseIdentifier(p *Parser) TreeNode {
	name := p.tokenStream.Consume().Text()
	return NewIdentifierNode(name)
}

func parsePrefixOperator(p *Parser) TreeNode {
	operator := p.tokenStream.Consume()
	operand := parseExpression(p, 0)
	return NewPrePostOperatorNode(PREFIX, operator.Text(), operand)
}

func parseInfixOperator(p *Parser, leftNode TreeNode) TreeNode {
	operator := p.tokenStream.Consume()
	rightNode := parseExpression(p, getPrecedence(operator.Text()))
	if operator.Text() == lexer.TERN_IF {
		p.tokenStream.ConsumeOnlyIf(lexer.COLON)
		ternElse := parseExpression(p, getPrecedence(operator.Text()))
		return NewTernaryNode(leftNode, rightNode, ternElse)
	}
	return NewInfixOperatorNode(leftNode, operator.Text(), rightNode)
}

func parsePostfixOperator(p *Parser, leftNode TreeNode) TreeNode {
	operator := p.tokenStream.Consume()
	return NewPrePostOperatorNode(POSTFIX, operator.Text(), leftNode)
}

func parseFuncCallArgs(p *Parser, leftNode TreeNode) TreeNode {
	p.tokenStream.ConsumeOnlyIf(lexer.OPEN_PAREN)
	fcNode := NewFuncCallNode(leftNode)
	for p.tokenStream.HasTokens() && p.tokenStream.Current().Text() != lexer.CLOSE_PAREN {
		fcNode.arg(parseExpression(p, 0))
		if p.tokenStream.Current().Text() == lexer.COMMA {
			p.tokenStream.Consume()
		}
	}
	p.tokenStream.Consume()
	return fcNode
}

func parseArrayIndex(p *Parser, leftNode TreeNode) TreeNode {
	p.tokenStream.ConsumeOnlyIf(lexer.OPEN_SQUARE)
	indexer := parseExpression(p, 0)
	p.tokenStream.ConsumeOnlyIf(lexer.CLOSE_SQUARE)
	arrIndNode := NewArrIndNode(leftNode, indexer)
	if p.tokenStream.HasTokens() && p.tokenStream.LookAhead(1).Text() == lexer.OPEN_SQUARE {
		arrIndNode = parseArrayIndex(p, arrIndNode).(*ArrIndNode)
	}
	return arrIndNode
}
