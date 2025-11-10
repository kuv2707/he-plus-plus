package parser

import (
	"fmt"
	"he++/lexer"
	nodes "he++/parser/node_types"
	"he++/utils"
)

func (p *Parser) ParseAST() *nodes.SourceFileNode {
	if !p.tokenStream.HasTokens() {
		parsingError("No tokens to parse!", -1)
		return nil
	}
	root := nodes.MakeSourceFileNode(p.Path)
	parseStatements(p, root)

	return root
}

func parseScope(p *Parser) nodes.TreeNode {
	t := p.tokenStream
	ls := t.ConsumeOnlyIf(lexer.LPAREN).LineNo()
	scopeNode := nodes.MakeScopeNode()
	parseStatements(p, scopeNode)
	le := t.ConsumeOnlyIf(lexer.RPAREN).LineNo()
	scopeNode.NodeMetadata = *nodes.MakeMetadata(ls, le)
	return scopeNode
}

func parseStatements(p *Parser, scope nodes.StatementsContainer) nodes.StatementsContainer {
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
				parsingError(fmt.Sprintf("Cannot parse %s", utils.Red(curr.Text())), curr.LineNo())
			}
			if !p.tokenStream.HasTokens() {
				break OUT
			}
		} else {
			scope.AddChild(parselet(p))
		}
	}
	return scope
}

func parseFunction(p *Parser) nodes.TreeNode {
	t := p.tokenStream
	ls := t.ConsumeOnlyIf(lexer.FUNCTION).LineNo()
	funcName := t.Consume()
	var argList []nodes.FuncArg
	t.ConsumeOnlyIf(lexer.OPEN_PAREN)
	for t.Current().Text() != lexer.CLOSE_PAREN {
		varName := t.Consume()
		dataType := parseDataType(p)
		argList = append(argList, nodes.FuncArg{Name: varName.Text(), DataT: dataType})
		t.ConsumeIf(lexer.COMMA)
	}
	le := t.ConsumeOnlyIf(lexer.CLOSE_PAREN).LineNo()

	funcNode := nodes.MakeFunctionNode(funcName.Text(), argList, parseDataType(p), parseScope(p).(*nodes.ScopeNode), nodes.MakeMetadata(ls, le))
	return funcNode
}

func parseReturnStatement(p *Parser) nodes.TreeNode {
	ls := p.tokenStream.ConsumeOnlyIf(lexer.RETURN).LineNo()
	return nodes.MakeReturnNode(parseExpression(p, 0), nodes.MakeMetadata(ls, ls))
}

func parseDataType(p *Parser) nodes.DataType {
	t := p.tokenStream
	currTok := t.Current()
	if currTok.Type() == lexer.IDENTIFIER {
		// possibly the type was absent and this ident is the varname
		if t.LookOneAhead().Text() == lexer.ASSN {
			return &nodes.UnspecifiedType{}
		}
		t.Consume()
		return &nodes.NamedType{Name: currTok.Text()}
	} else if currTok.Text() == lexer.OPEN_SQUARE {
		t.Consume()
		pt := &nodes.PrefixOfType{Prefix: nodes.ArrayOf, OfType: parseDataType(p)}
		t.ConsumeOnlyIf(lexer.CLOSE_SQUARE)
		return pt
	} else if currTok.Text() == lexer.AMP {
		t.Consume()
		return &nodes.PrefixOfType{Prefix: nodes.PointerOf, OfType: parseDataType(p)}

	} else if currTok.Text() == lexer.LPAREN {
		// anonymous object type
		return parseStructType(p)
	} else if currTok.Text() == lexer.FUNCTION {
		t.Consume()
		t.ConsumeOnlyIf(lexer.OPEN_PAREN)
		argsTypes := make([]nodes.DataType, 0)
		for t.Current().Text() != lexer.CLOSE_PAREN {
			argsTypes = append(argsTypes, parseDataType(p))
			t.ConsumeIf(lexer.COMMA)
		}
		t.ConsumeOnlyIf(lexer.CLOSE_PAREN)
		retType := parseDataType(p)
		return &nodes.FuncType{ArgTypes: argsTypes, ReturnType: retType}
	} else if currTok.Text() == lexer.VOID {
		t.Consume()
		return &nodes.VoidType{}
	}
	parsingError("Couldn't parse type: "+currTok.String(), currTok.LineNo())
	return nil

}

func parseVariableDeclaration(p *Parser) nodes.TreeNode {
	ls := p.tokenStream.ConsumeOnlyIf(lexer.LET).LineNo()
	dt := parseDataType(p)
	var decls []nodes.TreeNode
	decls = append(decls, parseExpression(p, 0))
	for p.tokenStream.Current().Text() == lexer.COMMA {
		p.tokenStream.Consume()
		decls = append(decls, parseExpression(p, 0))
	}
	varDec := nodes.MakeVariableDeclarationNode(decls, dt, nodes.MakeMetadata(ls, decls[len(decls)-1].Range().End))
	return varDec
}

func parseIfStatement(p *Parser) nodes.TreeNode {
	ls := p.tokenStream.ConsumeOnlyIf(lexer.IF).LineNo()
	var branches []nodes.ConditionalBranch
	ifCond := parseExpression(p, 0)
	p.tokenStream.ConsumeOnlyIf(lexer.THEN)
	ifScope := parseScope(p).(*nodes.ScopeNode)
	branches = append(branches, nodes.ConditionalBranch{Condition: ifCond, Scope: ifScope})
	le := ls
	if p.tokenStream.Current().Text() == lexer.ELSE {
		le = p.tokenStream.Consume().LineNo()
		tok := p.tokenStream.ConsumeIf(lexer.IF)
		var cond nodes.TreeNode
		if tok == nil {
			cond = parseExpression(p, 0)
		} else {
			cond = nodes.NewBooleanNode([]byte(lexer.TRUE), nodes.MakeMetadata(tok.LineNo(), tok.LineNo()))
		}
		scp := parseScope(p).(*nodes.ScopeNode)
		branches = append(branches, nodes.ConditionalBranch{Condition: cond, Scope: scp})

	}
	return nodes.MakeIfNode(branches, nodes.MakeMetadata(ls, le))
}

func parseLoopStatement(p *Parser) nodes.TreeNode {
	if p.tokenStream.Current().Text() == lexer.FOR {
		return parseForLoop(p)
	} else {
		return parseWhileLoop(p)
	}

}

func parseWhileLoop(p *Parser) nodes.TreeNode {
	ls := p.tokenStream.ConsumeOnlyIf(lexer.WHILE).LineNo()
	le := p.tokenStream.ConsumeOnlyIf(lexer.THAT).LineNo()
	condNode := parseExpression(p, 0)
	return nodes.MakeLoopNode(&nodes.EmptyPlaceholderNode{}, condNode, &nodes.EmptyPlaceholderNode{}, parseScope(p).(*nodes.ScopeNode), nodes.MakeMetadata(ls, le))
}

func parseForLoop(p *Parser) nodes.TreeNode {
	ls := p.tokenStream.ConsumeOnlyIf(lexer.FOR).LineNo()
	varDecl := parseVariableDeclaration(p)
	p.tokenStream.ConsumeOnlyIf(lexer.SEMICOLON)
	condNode := parseExpression(p, 0)
	le := p.tokenStream.ConsumeOnlyIf(lexer.SEMICOLON).LineNo()
	updNode := parseExpression(p, 0)
	return nodes.MakeLoopNode(varDecl, condNode, updNode, parseScope(p).(*nodes.ScopeNode), nodes.MakeMetadata(ls, le))
}

func parseStructType(p *Parser) *nodes.StructType {
	p.tokenStream.ConsumeOnlyIf(lexer.LPAREN)
	strct := nodes.StructType{}
	for p.tokenStream.Current().Text() != lexer.RPAREN {
		Name := p.tokenStream.Consume().Text()
		Type := parseDataType(p)
		strct.Fields = append(strct.Fields, nodes.StructFieldType{Name: Name, Type: Type})
		p.tokenStream.ConsumeIf(lexer.COMMA)
	}
	p.tokenStream.ConsumeOnlyIf(lexer.RPAREN)
	return &strct
}

func parseStructDefn(p *Parser) nodes.TreeNode {
	p.tokenStream.ConsumeOnlyIf(lexer.STRUCT)
	name := p.tokenStream.ConsumeOnlyIfType(lexer.IDENTIFIER).Text()
	return &nodes.StructDefnNode{Name: name, StructDef: parseStructType(p)}
}
