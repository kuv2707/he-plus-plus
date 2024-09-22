package parser

import (
	"fmt"
	"he++/globals"
	"he++/lexer"
)

func (p *Parser) ParseAST() *SourceFileNode {
	if !p.tokenStream.HasTokens() {
		parsingError("No tokens to parse!", -1)
		return nil
	}
	root := MakeSourceFileNode()
	parseStatements(p, root)
	return root
}

func parseScope(p *Parser) *ScopeNode {
	t := p.tokenStream
	t.ConsumeOnlyIf(lexer.LPAREN)
	scopeNode := MakeScopeNode()
	parseStatements(p, scopeNode)
	t.ConsumeOnlyIf(lexer.RPAREN)
	return scopeNode
}

func parseStatements(p *Parser, scope StatementsContainer) StatementsContainer {
OUT:
	for p.tokenStream.HasTokens() {
		curr := p.tokenStream.Current()
		if curr.Type() == lexer.BRACKET && curr.Text() == lexer.RPAREN {
			break OUT
		}
		parselet, exists := p.scopeParselets[curr.Text()]
		if !exists {
			expr := parseExpression(p, 0)
			if expr != nil {
				scope.AddChild(expr)
			} else {
				parsingError(fmt.Sprintf("Cannot parse %s", globals.Red(curr.Text())), curr.LineNo())
			}
			break OUT
		} else {
			scope.AddChild(parselet(p))
		}
	}
	return scope
}

func parseFunction(p *Parser) TreeNode {
	t := p.tokenStream
	t.ConsumeOnlyIf(lexer.FUNCTION)
	funcName := t.Consume()
	funcNode := MakeFunctionNode(funcName.Text())
	t.ConsumeOnlyIf(lexer.OPEN_PAREN)
	for t.Current().Text() != lexer.CLOSE_PAREN {
		varName := t.Consume()
		dataType := t.Consume()
		funcNode.AddArg(varName.Text(), dataType.Text())
		t.ConsumeIf(lexer.COMMA)
	}
	t.ConsumeOnlyIf(lexer.CLOSE_PAREN)
	funcNode.scope = parseScope(p)
	return funcNode
}

func parseReturnStatement(p *Parser) TreeNode {
	p.tokenStream.ConsumeOnlyIf(lexer.RETURN)
	return MakeReturnNode(parseExpression(p, 0))
}

func parseVariableDeclaration(p *Parser) TreeNode {
	p.tokenStream.ConsumeOnlyIf(lexer.LET)
	varDec := MakeVariableDeclarationNode()
	varDec.AddDeclaration(parseExpression(p, 0))
	for p.tokenStream.Current().Text() == lexer.COMMA {
		p.tokenStream.Consume()
		varDec.AddDeclaration(parseExpression(p, 0))
	}
	return varDec
}

func parseIfStatement(p *Parser) TreeNode {
	p.tokenStream.ConsumeOnlyIf(lexer.IF)
	ifCond := parseExpression(p, 0)
	p.tokenStream.ConsumeOnlyIf(lexer.THEN)
	ifScope := parseScope(p)
	p.tokenStream.ConsumeOnlyIf(lexer.ELSE)
	if p.tokenStream.Current().Text() == lexer.IF {
		elseScope := parseIfStatement(p)
		return MakeIfNode(ifCond, ifScope, elseScope.(*IfNode))
	}
	elseScope := parseScope(p)
	return MakeIfNode(ifCond, ifScope, elseScope)
}

func parseLoopStatement(p *Parser) TreeNode {
	if p.tokenStream.Current().Text() == lexer.FOR {
		return parseForLoop(p)
	} else {
		return parseWhileLoop(p)
	}

}

func parseWhileLoop(p *Parser) TreeNode {
	p.tokenStream.ConsumeOnlyIf(lexer.WHILE)
	p.tokenStream.ConsumeOnlyIf(lexer.THAT)
	condNode := parseExpression(p, 0)
	return MakeLoopNode(&EmptyPlaceholderNode{}, condNode, &EmptyPlaceholderNode{}, parseScope(p))
}

func parseForLoop(p *Parser) TreeNode {
	p.tokenStream.ConsumeOnlyIf(lexer.FOR)
	varDecl := parseVariableDeclaration(p)
	p.tokenStream.ConsumeOnlyIf(lexer.SEMICOLON)
	condNode := parseExpression(p, 0)
	p.tokenStream.ConsumeOnlyIf(lexer.SEMICOLON)
	updNode := parseExpression(p, 0)
	return MakeLoopNode(varDecl, condNode, updNode, parseScope(p))
}
