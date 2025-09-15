package parser

import (
	"fmt"
	"he++/globals"
	"he++/lexer"
	nodes "he++/parser/node_types"
	"sort"
)

func (p *Parser) ParseAST() *nodes.SourceFileNode {
	if !p.tokenStream.HasTokens() {
		parsingError("No tokens to parse!", -1)
		return nil
	}
	root := nodes.MakeSourceFileNode()
	parseStatements(p, root)

	// todo: add some postprocessing
	sort.SliceStable(root.Children, func(i, j int) bool {
		return root.Children[i].Type() == nodes.VAR_DECL
	})
	return root
}

func parseScope(p *Parser) *nodes.ScopeNode {
	t := p.tokenStream
	t.ConsumeOnlyIf(lexer.LPAREN)
	scopeNode := nodes.MakeScopeNode()
	parseStatements(p, scopeNode)
	t.ConsumeOnlyIf(lexer.RPAREN)
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
				parsingError(fmt.Sprintf("Cannot parse %s", globals.Red(curr.Text())), curr.LineNo())
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
	t.ConsumeOnlyIf(lexer.FUNCTION)
	funcName := t.Consume()
	funcNode := nodes.MakeFunctionNode(funcName.Text())
	t.ConsumeOnlyIf(lexer.OPEN_PAREN)
	for t.Current().Text() != lexer.CLOSE_PAREN {
		varName := t.Consume()
		dataType := parseDataType(p)
		funcNode.AddArg(varName.Text(), dataType)
		t.ConsumeIf(lexer.COMMA)
	}
	t.ConsumeOnlyIf(lexer.CLOSE_PAREN)
	// parse properly using parseDataType
	funcNode.ReturnType = parseDataType(p)
	// fmt.Println("dtype for " + funcNode.Name + " " + funcNode.ReturnType.Text())
	funcNode.Scope = parseScope(p)
	// fmt.Println(funcNode.Scope.String(""))
	return funcNode
}

func parseReturnStatement(p *Parser) nodes.TreeNode {
	p.tokenStream.ConsumeOnlyIf(lexer.RETURN)
	return nodes.MakeReturnNode(parseExpression(p, 0))
}

func parseDataType(p *Parser) nodes.DataType {
	t := p.tokenStream
	currTok := t.Current()
	if currTok.Type() == lexer.IDENTIFIER {
		// possibly the type was absent and this ident is the varname
		if t.LookAhead(1).Text() == lexer.ASSN {
			return &nodes.UnspecifiedType{}
		}
		t.Consume()
		return &nodes.NamedType{Name: currTok.Text()}
	} else if currTok.Text() == lexer.OPEN_SQUARE {
		t.Consume()
		t.ConsumeOnlyIf(lexer.CLOSE_SQUARE)
		return &nodes.PrefixOfType{Prefix: nodes.ArrayOf, OfType: parseDataType(p)}
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
	return &nodes.ErrorType{Message: "Couldn't parse type: " + currTok.String()}

}

func parseVariableDeclaration(p *Parser) nodes.TreeNode {
	p.tokenStream.ConsumeOnlyIf(lexer.LET)
	varDec := nodes.MakeVariableDeclarationNode()
	// todo: add type
	varDec.SetDataType(parseDataType(p))
	varDec.AddDeclaration(parseExpression(p, 0))
	for p.tokenStream.Current().Text() == lexer.COMMA {
		p.tokenStream.Consume()
		varDec.AddDeclaration(parseExpression(p, 0))
	}
	return varDec
}

func parseIfStatement(p *Parser) nodes.TreeNode {
	p.tokenStream.ConsumeOnlyIf(lexer.IF)
	ifCond := parseExpression(p, 0)
	p.tokenStream.ConsumeOnlyIf(lexer.THEN)
	ifScope := parseScope(p)
	p.tokenStream.ConsumeOnlyIf(lexer.ELSE)
	if p.tokenStream.Current().Text() == lexer.IF {
		elseScope := parseIfStatement(p)
		return nodes.MakeIfNode(ifCond, ifScope, elseScope.(*nodes.IfNode))
	}
	elseScope := parseScope(p)
	return nodes.MakeIfNode(ifCond, ifScope, elseScope)
}

func parseLoopStatement(p *Parser) nodes.TreeNode {
	if p.tokenStream.Current().Text() == lexer.FOR {
		return parseForLoop(p)
	} else {
		return parseWhileLoop(p)
	}

}

func parseWhileLoop(p *Parser) nodes.TreeNode {
	p.tokenStream.ConsumeOnlyIf(lexer.WHILE)
	p.tokenStream.ConsumeOnlyIf(lexer.THAT)
	condNode := parseExpression(p, 0)
	return nodes.MakeLoopNode(&nodes.EmptyPlaceholderNode{}, condNode, &nodes.EmptyPlaceholderNode{}, parseScope(p))
}

func parseForLoop(p *Parser) nodes.TreeNode {
	p.tokenStream.ConsumeOnlyIf(lexer.FOR)
	varDecl := parseVariableDeclaration(p)
	p.tokenStream.ConsumeOnlyIf(lexer.SEMICOLON)
	condNode := parseExpression(p, 0)
	p.tokenStream.ConsumeOnlyIf(lexer.SEMICOLON)
	updNode := parseExpression(p, 0)
	return nodes.MakeLoopNode(varDecl, condNode, updNode, parseScope(p))
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

