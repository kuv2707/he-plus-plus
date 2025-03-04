package parser

import (
	"he++/lexer"
	nodes "he++/parser/node_types"
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

func parseExpression(p *Parser, prec float32) nodes.TreeNode {
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

func parseBracketExpression(p *Parser) nodes.TreeNode {
	p.tokenStream.Consume()
	expr := parseExpression(p, 0)
	if p.tokenStream.Current().Text() != lexer.CLOSE_PAREN {
		parsingError("Expected closing parenthesis", p.tokenStream.Current().LineNo())
	}
	p.tokenStream.Consume()
	return expr
}

func parseInteger(p *Parser) nodes.TreeNode {
	return nodes.NewNumberNode([]byte(p.tokenStream.Consume().Text()), "int")
}

func parseFloat(p *Parser) nodes.TreeNode {
	return nodes.NewNumberNode([]byte(p.tokenStream.Consume().Text()), "float")
}

func parseString(p *Parser) nodes.TreeNode {
	return nodes.NewStringNode([]byte(p.tokenStream.Consume().Text()))
}

func parseBoolean(p *Parser) nodes.TreeNode {
	tok := p.tokenStream.Consume()
	truth := tok.Text() == lexer.TRUE
	if truth {
		return nodes.NewBooleanNode([]byte(lexer.TRUE))
	}
	if tok.Text() != lexer.FALSE {
		parsingError("Expected boolean value", tok.LineNo())
	}
	return nodes.NewBooleanNode([]byte(lexer.FALSE))
}

func parseIdentifier(p *Parser) nodes.TreeNode {
	name := p.tokenStream.Consume().Text()
	return nodes.NewIdentifierNode(name)
}

func parsePrefixOperator(p *Parser) nodes.TreeNode {
	operator := p.tokenStream.Consume()
	operand := parseExpression(p, 0)
	return nodes.NewPrePostOperatorNode(nodes.PREFIX, operator.Text(), operand)
}

func parseInfixOperator(p *Parser, leftNode nodes.TreeNode) nodes.TreeNode {
	operator := p.tokenStream.Consume()
	rightNode := parseExpression(p, getPrecedence(operator.Text()))
	if operator.Text() == lexer.TERN_IF {
		p.tokenStream.ConsumeOnlyIf(lexer.COLON)
		ternElse := parseExpression(p, getPrecedence(operator.Text()))
		return nodes.NewTernaryNode(leftNode, rightNode, ternElse)
	}
	return nodes.NewInfixOperatorNode(leftNode, operator.Text(), rightNode)
}

func parsePostfixOperator(p *Parser, leftNode nodes.TreeNode) nodes.TreeNode {
	operator := p.tokenStream.Consume()
	return nodes.NewPrePostOperatorNode(nodes.POSTFIX, operator.Text(), leftNode)
}

func parseFuncCallArgs(p *Parser, leftNode nodes.TreeNode) nodes.TreeNode {
	p.tokenStream.ConsumeOnlyIf(lexer.OPEN_PAREN)
	fcNode := nodes.NewFuncCallNode(leftNode)
	for p.tokenStream.HasTokens() && p.tokenStream.Current().Text() != lexer.CLOSE_PAREN {
		fcNode.Arg(parseExpression(p, 0))
		if p.tokenStream.Current().Text() == lexer.COMMA {
			p.tokenStream.Consume()
		}
	}
	p.tokenStream.Consume()
	return fcNode
}

func parseArrayIndex(p *Parser, leftNode nodes.TreeNode) nodes.TreeNode {
	p.tokenStream.ConsumeOnlyIf(lexer.OPEN_SQUARE)
	indexer := parseExpression(p, 0)
	p.tokenStream.ConsumeOnlyIf(lexer.CLOSE_SQUARE)
	arrIndNode := nodes.NewArrIndNode(leftNode, indexer)
	if p.tokenStream.HasTokens() && p.tokenStream.LookAhead(1).Text() == lexer.OPEN_SQUARE {
		arrIndNode = parseArrayIndex(p, arrIndNode).(*nodes.ArrIndNode)
	}
	return arrIndNode
}

func parseArrayDeclaration(p *Parser) nodes.TreeNode {
	p.tokenStream.ConsumeOnlyIf(lexer.OPEN_SQUARE)
	p.tokenStream.ConsumeOnlyIf(lexer.CLOSE_SQUARE)
	p.tokenStream.ConsumeOnlyIf(lexer.LPAREN)
	elems := make([]nodes.TreeNode, 0)
	for p.tokenStream.Current().Text() != lexer.RPAREN {
		k := parseExpression(p, 0.0)
		elems = append(elems, k)
		p.tokenStream.ConsumeIf(lexer.COMMA)
	}
	
	return nodes.ArrayDeclaration{
		Elems: elems,
	}
}
