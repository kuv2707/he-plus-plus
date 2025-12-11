package parser

import (
	"bytes"
	"encoding/binary"
	"fmt"
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
		parsingError("Unexpected end of file while parsing expression", 0)
		return nil
	}
	tok := t.Current()
	prefix, exists := p.getPrefixParselet(*tok)
	if !exists {
		panic(fmt.Sprintf("Might not be an expression: %s %s %d", tok.Type().String(), tok.Text(), tok.LineNo()))
	}
	leftNode := prefix(p)
	for t.HasTokens() {
		opSymbol := t.Current().Text()
		if getPrecedence(opSymbol) <= prec {
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
	t := p.tokenStream.Consume()
	return nodes.NewNumberNode([]byte(t.Text()), nodes.INT32_NUMBER, nodes.MakeMetadata(t.LineNo(), t.LineNo()))
}

func parseFloat(p *Parser) nodes.TreeNode {
	t := p.tokenStream.Consume()
	return nodes.NewNumberNode([]byte(t.Text()), nodes.FLOAT_NUMBER, nodes.MakeMetadata(t.LineNo(), t.LineNo()))
}

func parseString(p *Parser) nodes.TreeNode {
	t := p.tokenStream.Consume()
	return nodes.NewStringNode([]byte(t.Text()), nodes.MakeMetadata(t.LineNo(), t.LineNo()))
}

func parseBoolean(p *Parser) nodes.TreeNode {
	tok := p.tokenStream.Consume()
	truth := tok.Text() == lexer.TRUE
	if truth {
		return nodes.NewBooleanNode(true, nodes.MakeMetadata(tok.LineNo(), tok.LineNo()))
	}
	if tok.Text() != lexer.FALSE {
		parsingError("Expected boolean value", tok.LineNo())
	}
	return nodes.NewBooleanNode(false, nodes.MakeMetadata(tok.LineNo(), tok.LineNo()))
}

func parseIdentifier(p *Parser) nodes.TreeNode {
	t := p.tokenStream.Consume()
	return nodes.NewIdentifierNode(t.Text(), nodes.MakeMetadata(t.LineNo(), t.LineNo()))
}

func parsePrefixOperator(p *Parser) nodes.TreeNode {
	operator := p.tokenStream.Consume()
	operand := parseExpression(p, 0)
	return nodes.NewPrePostOperatorNode(nodes.PREFIX, operator.Text(), operand, nodes.MakeMetadata(operator.LineNo(), operator.LineNo()))
}

func parseInfixOperator(p *Parser, leftNode nodes.TreeNode) nodes.TreeNode {
	operator := p.tokenStream.Consume()
	rightNode := parseExpression(p, getPrecedence(operator.Text()))
	if operator.Text() == lexer.TERN_IF {
		p.tokenStream.ConsumeOnlyIf(lexer.COLON)
		ternElse := parseExpression(p, getPrecedence(operator.Text()))
		return nodes.NewTernaryNode(leftNode, rightNode, ternElse)
	}
	return nodes.NewInfixOperatorNode(leftNode, operator.Text(), rightNode, nodes.MakeMetadata(operator.LineNo(), operator.LineNo()))
}

func parsePostfixOperator(p *Parser, leftNode nodes.TreeNode) nodes.TreeNode {
	operator := p.tokenStream.Consume()
	return nodes.NewPrePostOperatorNode(nodes.POSTFIX, operator.Text(), leftNode, nodes.MakeMetadata(operator.LineNo(), operator.LineNo()))
}

func parseFuncCallArgs(p *Parser, leftNode nodes.TreeNode) nodes.TreeNode {
	ls := p.tokenStream.ConsumeOnlyIf(lexer.OPEN_PAREN).LineNo()
	var args []nodes.TreeNode
	for p.tokenStream.HasTokens() && p.tokenStream.Current().Text() != lexer.CLOSE_PAREN {
		args = append(args, parseExpression(p, 0))
		if p.tokenStream.Current().Text() == lexer.COMMA {
			p.tokenStream.Consume()
		}
	}
	le := p.tokenStream.Consume().LineNo()
	fcNode := nodes.NewFuncCallNode(leftNode, nodes.MakeMetadata(ls, le))
	fcNode.Args = args
	return fcNode
}

func parseArrayIndex(p *Parser, leftNode nodes.TreeNode) nodes.TreeNode {
	p.tokenStream.ConsumeOnlyIf(lexer.OPEN_SQUARE).LineNo()
	indexer := parseExpression(p, 0)
	le := p.tokenStream.ConsumeOnlyIf(lexer.CLOSE_SQUARE).LineNo()
	arrIndNode := nodes.NewArrIndNode(leftNode, indexer, nodes.MakeMetadata(leftNode.Range().Start, le))
	if p.tokenStream.HasTokens() && p.tokenStream.LookOneAhead().Text() == lexer.OPEN_SQUARE {
		arrIndNode = parseArrayIndex(p, arrIndNode).(*nodes.ArrIndNode)
	}
	return arrIndNode
}

func parseArrayDeclaration(p *Parser) nodes.TreeNode {
	ls := p.tokenStream.ConsumeOnlyIf(lexer.OPEN_SQUARE).LineNo()
	dt := parseDataType(p)
	p.tokenStream.ConsumeOnlyIf(lexer.CLOSE_SQUARE)
	if p.tokenStream.Current().Text() == lexer.OPEN_SQUARE {
		p.tokenStream.Consume()
		size := parseExpression(p, 0)
		le := p.tokenStream.ConsumeOnlyIf(lexer.CLOSE_SQUARE).LineNo()
		return nodes.MakeArrayDeclarationNode(size, nil, dt, nodes.MakeMetadata(ls, le))

	} else {
		p.tokenStream.ConsumeOnlyIf(lexer.LPAREN)
		elems := make([]nodes.TreeNode, 0)
		for p.tokenStream.Current().Text() != lexer.RPAREN {
			k := parseExpression(p, 0.0)
			elems = append(elems, k)
			p.tokenStream.ConsumeIf(lexer.COMMA)
		}
		le := p.tokenStream.ConsumeOnlyIf(lexer.RPAREN).LineNo()
		str := new(bytes.Buffer)
		binary.Write(str, binary.BigEndian, int64(len(elems)))
		return nodes.MakeArrayDeclarationNode(
			nodes.NewNumberNode(
				str.Bytes(),
				nodes.INT32_NUMBER, nodes.MakeMetadata(ls, le),
			),
			elems, dt,
			nodes.MakeMetadata(ls, le),
		)
	}
}

func parseStructValue(p *Parser) nodes.TreeNode {
	ls := p.tokenStream.ConsumeOnlyIf(lexer.LPAREN).LineNo()
	var mp map[string]nodes.TreeNode = make(map[string]nodes.TreeNode)
	for p.tokenStream.Current().Text() != lexer.RPAREN {
		name := p.tokenStream.ConsumeOnlyIfType(lexer.IDENTIFIER).Text()
		p.tokenStream.ConsumeOnlyIf(lexer.COLON)
		val := parseExpression(p, 0)
		mp[name] = val
		p.tokenStream.ConsumeIf(lexer.COMMA)
	}
	le := p.tokenStream.ConsumeOnlyIf(lexer.RPAREN).LineNo()
	return nodes.MakeStructValueNode(mp, nodes.MakeMetadata(ls, le))
}
