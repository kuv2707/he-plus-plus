package parser

import (
	"fmt"
	"he++/lexer"
	"he++/parser/node_types"
	"he++/utils"
	"math"
)

func (p *Parser) ParseAST() *node_types.SourceFileNode {
	if !p.tokenStream.HasTokens() {
		parsingError("No tokens to parse!", -1)
		return nil
	}
	root := node_types.MakeSourceFileNode(p.Path)
	parseStatements(p, root)

	return root
}

func parseScope(p *Parser) node_types.TreeNode {
	t := p.tokenStream
	ls := t.ConsumeOnlyIf(lexer.LPAREN).LineNo()
	scopeNode := node_types.MakeScopeNode()
	parseStatements(p, scopeNode)
	le := t.ConsumeOnlyIf(lexer.RPAREN).LineNo()
	scopeNode.NodeMetadata = *node_types.MakeMetadata(ls, le)
	return scopeNode
}

func parseStatements(p *Parser, scope node_types.StatementsContainer) node_types.StatementsContainer {
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

func parseFunction(p *Parser) node_types.TreeNode {
	t := p.tokenStream
	ls := t.ConsumeOnlyIf(lexer.FUNCTION).LineNo()
	funcName := t.Consume()
	var argList []node_types.FuncArg
	t.ConsumeOnlyIf(lexer.OPEN_PAREN)
	for t.Current().Text() != lexer.CLOSE_PAREN {
		varName := t.Consume()
		dataType := parseDataType(p)
		argList = append(argList, node_types.FuncArg{Name: varName.Text(), DataT: dataType})
		t.ConsumeIf(lexer.COMMA)
	}
	le := t.ConsumeOnlyIf(lexer.CLOSE_PAREN).LineNo()

	funcNode := node_types.MakeFunctionNode(funcName.Text(), argList, parseDataType(p), parseScope(p).(*node_types.ScopeNode), node_types.MakeMetadata(ls, le))
	return funcNode
}

func parseReturnStatement(p *Parser) node_types.TreeNode {
	ls := p.tokenStream.ConsumeOnlyIf(lexer.RETURN).LineNo()
	return node_types.MakeReturnNode(parseExpression(p, 0), node_types.MakeMetadata(ls, ls))
}

func parseDataType(p *Parser) node_types.DataType {
	fmt.Println("PDT")
	t := p.tokenStream
	currTok := t.Current()
	if currTok.Type() == lexer.IDENTIFIER {
		// possibly the type was absent and this ident is the varname
		if t.LookOneAhead().Text() == lexer.ASSN {
			return &node_types.UnspecifiedType{DataTypeMetaData: node_types.DataTypeMetaData{TypeSize: 0, Tid: node_types.UniqueTypeId()}}
		}
		t.Consume()
		return &node_types.NamedType{
			Name:             currTok.Text(),
			DataTypeMetaData: node_types.DataTypeMetaData{TypeSize: -1, Tid: node_types.UniqueTypeId()}}
	} else if currTok.Text() == lexer.OPEN_SQUARE {
		t.Consume()
		pt := &node_types.PrefixOfType{Prefix: node_types.ArrayOf, OfType: parseDataType(p), DataTypeMetaData: node_types.DataTypeMetaData{TypeSize: node_types.POINTER_SIZE, Tid: node_types.UniqueTypeId()}}
		t.ConsumeOnlyIf(lexer.CLOSE_SQUARE)
		return pt
	} else if currTok.Text() == lexer.AMP {
		t.Consume()
		return &node_types.PrefixOfType{Prefix: node_types.PointerOf, OfType: parseDataType(p), DataTypeMetaData: node_types.DataTypeMetaData{TypeSize: node_types.POINTER_SIZE, Tid: node_types.UniqueTypeId()}}

	} else if currTok.Text() == lexer.LPAREN {
		// anonymous object type
		return parseStructType(p)
	} else if currTok.Text() == lexer.FUNCTION {
		t.Consume()
		t.ConsumeOnlyIf(lexer.OPEN_PAREN)
		argsTypes := make([]node_types.DataType, 0)
		for t.Current().Text() != lexer.CLOSE_PAREN {
			argsTypes = append(argsTypes, parseDataType(p))
			t.ConsumeIf(lexer.COMMA)
		}
		t.ConsumeOnlyIf(lexer.CLOSE_PAREN)
		retType := parseDataType(p)
		return &node_types.FuncType{ArgTypes: argsTypes, ReturnType: retType, DataTypeMetaData: node_types.DataTypeMetaData{TypeSize: node_types.POINTER_SIZE, Tid: node_types.UniqueTypeId()}}
	} else if currTok.Text() == lexer.VOID {
		t.Consume()
		return &node_types.VOID_DATATYPE
	}
	parsingError("Couldn't parse type: "+currTok.String(), currTok.LineNo())
	return nil

}

func parseVariableDeclaration(p *Parser) node_types.TreeNode {
	ls := p.tokenStream.ConsumeOnlyIf(lexer.LET).LineNo()
	dt := parseDataType(p)
	var decls []node_types.TreeNode
	decls = append(decls, parseExpression(p, 0))
	for p.tokenStream.Current().Text() == lexer.COMMA {
		p.tokenStream.Consume()
		decls = append(decls, parseExpression(p, 0))
	}
	varDec := node_types.MakeVariableDeclarationNode(decls, dt, node_types.MakeMetadata(ls, decls[len(decls)-1].Range().End))
	return varDec
}

func parseIfStatement(p *Parser) node_types.TreeNode {
	ls := p.tokenStream.ConsumeOnlyIf(lexer.IF).LineNo()
	var branches []node_types.ConditionalBranch
	ifCond := parseExpression(p, 0)
	p.tokenStream.ConsumeOnlyIf(lexer.THEN)
	ifScope := parseScope(p).(*node_types.ScopeNode)
	branches = append(branches, node_types.ConditionalBranch{Condition: ifCond, Scope: ifScope})
	le := ls
	for p.tokenStream.Current().Text() == lexer.ELSE {
		le = p.tokenStream.Consume().LineNo()
		tok := p.tokenStream.ConsumeIf(lexer.IF)
		var cond node_types.TreeNode
		if tok != nil {
			cond = parseExpression(p, 0)
			p.tokenStream.ConsumeOnlyIf(lexer.THEN)
		} else {
			cond = node_types.NewBooleanNode(true, node_types.MakeMetadata(le, le))
		}
		scp := parseScope(p).(*node_types.ScopeNode)
		branches = append(branches, node_types.ConditionalBranch{Condition: cond, Scope: scp})

	}
	return node_types.MakeIfNode(branches, node_types.MakeMetadata(ls, le))
}

func parseLoopStatement(p *Parser) node_types.TreeNode {
	if p.tokenStream.Current().Text() == lexer.FOR {
		return parseForLoop(p)
	} else {
		return parseWhileLoop(p)
	}

}

func parseWhileLoop(p *Parser) node_types.TreeNode {
	ls := p.tokenStream.ConsumeOnlyIf(lexer.WHILE).LineNo()
	le := p.tokenStream.ConsumeOnlyIf(lexer.THAT).LineNo()
	condNode := parseExpression(p, 0)
	return node_types.MakeLoopNode(&node_types.EmptyPlaceholderNode{}, condNode, &node_types.EmptyPlaceholderNode{}, parseScope(p).(*node_types.ScopeNode), node_types.MakeMetadata(ls, le))
}

func parseForLoop(p *Parser) node_types.TreeNode {
	ls := p.tokenStream.ConsumeOnlyIf(lexer.FOR).LineNo()
	varDecl := parseVariableDeclaration(p)
	p.tokenStream.ConsumeOnlyIf(lexer.SEMICOLON)
	condNode := parseExpression(p, 0)
	le := p.tokenStream.ConsumeOnlyIf(lexer.SEMICOLON).LineNo()
	updNode := parseExpression(p, 0)
	return node_types.MakeLoopNode(varDecl, condNode, updNode, parseScope(p).(*node_types.ScopeNode), node_types.MakeMetadata(ls, le))
}

func parseStructType(p *Parser) *node_types.StructType {
	p.tokenStream.ConsumeOnlyIf(lexer.LPAREN)
	strct := node_types.StructType{}
	siz := 0
	for p.tokenStream.Current().Text() != lexer.RPAREN {
		name := p.tokenStream.Consume().Text()
		dt := parseDataType(p)
		siz += dt.Size()
		strct.Fields = append(strct.Fields, node_types.StructFieldTypeInfo{Name: name, Type: dt})
		p.tokenStream.ConsumeIf(lexer.COMMA)
	}
	strct.Tid = node_types.UniqueTypeId()
	// 4 byte alignment
	strct.TypeSize = (int)(math.Ceil(float64(siz)/4)) * 4
	p.tokenStream.ConsumeOnlyIf(lexer.RPAREN)
	return &strct
}

func parseStructDefn(p *Parser) node_types.TreeNode {
	p.tokenStream.ConsumeOnlyIf(lexer.STRUCT)
	name := p.tokenStream.ConsumeOnlyIfType(lexer.IDENTIFIER).Text()
	return &node_types.StructDefnNode{Name: name, StructDef: parseStructType(p)}
}
