package parser

import (
	"he++/lexer"
)

// pratt parsing for expressions

func parseExpression(p *Parser, prec float32) TreeNode {
	t := p.tokenStream
	if !t.HasNext() {
		parsingError("Unexpected end of expression after "+t.LookAhead(-1).Text(), 0)
		return nil
	}
	tok := t.Current()
	prefix, exists := p.getPrefixParselet(tok)
	if !exists {
		panic("Might not be an expression: " + tok.Type().String() + " " + tok.Text())
	}
	leftNode := prefix(p)
	for t.HasNext() {
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
	return NewInfixOperatorNode(leftNode, operator.Text(), rightNode)
}

func parsePostfixOperator(p *Parser, leftNode TreeNode) TreeNode {
	operator := p.tokenStream.Consume()
	return NewPrePostOperatorNode(POSTFIX, operator.Text(), leftNode)
}

func parseFuncCallArgs(p *Parser, leftNode TreeNode) TreeNode {
	p.tokenStream.ConsumeIf(lexer.OPEN_PAREN)
	fcNode := NewFuncCallNode(leftNode)
	for p.tokenStream.HasNext() && p.tokenStream.Current().Text() != lexer.CLOSE_PAREN {
		fcNode.arg(parseExpression(p, 0))
		if p.tokenStream.Current().Text() == lexer.COMMA {
			p.tokenStream.Consume()
		}
	}
	p.tokenStream.Consume()
	return fcNode
}

func parseArrayIndex(p *Parser, leftNode TreeNode) TreeNode {
	p.tokenStream.ConsumeIf(lexer.OPEN_SQUARE)
	indexer := parseExpression(p, 0)
	p.tokenStream.ConsumeIf(lexer.CLOSE_SQUARE)
	arrIndNode := NewArrIndNode(leftNode, indexer)
	if p.tokenStream.HasNext() && p.tokenStream.LookAhead(1).Text() == lexer.OPEN_SQUARE {
		arrIndNode = parseArrayIndex(p, arrIndNode).(*ArrIndNode)
	}
	return arrIndNode
}
